package chronos

import (
	"context"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/github"
)

type ChronosUpdateSingleIssueDeadlineRequest struct {
	IssueNumber int
	LabelName   string
	Created     time.Time

	elapsedTime  float64
	nonWorkHours float64
	timer        string
	overdue      bool
	timerLabel   string
}

type ChronosUpdateSingleIssueDeadlineResponse struct {
}

func (h *Chronos) calculateElapsedTime() error {
	var (
		req          = h.request.(ChronosUpdateSingleIssueDeadlineRequest)
		nonWorkHours float64
		weekendHours float64
		holidayHours float64
	)

	loc, _ := time.LoadLocation(STANDARD_TIME_LOCATION)
	now := time.Now().In(loc)
	created := req.Created.In(loc)
	elapsedTime := int(math.Round(now.Sub(created).Hours()))

	holidays, err := h.GetHolidays(now.Year())
	if err != nil {
		return err
	}

	for t := 0; t < elapsedTime; t++ {
		// Check if it is Sunday
		if created.Add(time.Duration(t)*time.Hour).Weekday() == 0 {
			weekendHours++
			continue
		}

		// Check if it is Saturday
		if created.Add(time.Duration(t)*time.Hour).Weekday() == 6 {
			weekendHours++
			continue
		}

		// Check for holidays
		_, exists := holidays[created.Add(time.Duration(t)*time.Hour).Format("2006-01-02")]
		if exists {
			holidayHours++
			continue
		}

		// Check if it is a work hour
		if created.Add(time.Duration(t)*time.Hour).Hour() < WORK_HOURS_INITIAL {
			nonWorkHours++
			continue
		}
		if created.Add(time.Duration(t)*time.Hour).Hour() >= WORK_HOURS_FINAL {
			nonWorkHours++
			continue
		}
	}

	req.elapsedTime = now.Sub(created).Hours() - weekendHours - holidayHours
	req.nonWorkHours = nonWorkHours
	h.request = req

	return nil
}

func (h *Chronos) defineNewDeadline() error {
	var (
		req                = h.request.(ChronosUpdateSingleIssueDeadlineRequest)
		deadline           string
		deducer            float64
		deduceNonWorkHours bool
	)

	timeTable := make(map[string]float64)
	timeTable[DEADLINE_TYPE_HOURS] = req.elapsedTime - deducer*req.nonWorkHours
	timeTable[DEADLINE_TYPE_DAYS] = (req.elapsedTime - deducer*req.nonWorkHours) / (WORK_HOURS_FINAL - WORK_HOURS_INITIAL)

	switch req.LabelName {
	case PRIORITY_LABEL_PRIORIDADE_BAIXA:
		deadline = DEADLINE_LABEL_PRIORIDADE_BAIXA
		deduceNonWorkHours = true
	case PRIORITY_LABEL_PRIORIDADE_MEDIA:
		deadline = DEADLINE_LABEL_PRIORIDADE_MEDIA
		deduceNonWorkHours = true
	case PRIORITY_LABEL_PRIORIDADE_ALTA:
		deadline = DEADLINE_LABEL_PRIORIDADE_ALTA
		deduceNonWorkHours = true
	case PRIORITY_LABEL_PRIORIDADE_MUITO_ALTA:
		deadline = DEADLINE_LABEL_PRIORIDADE_MUITO_ALTA
		deduceNonWorkHours = false
	default:
		return ErrUnableToDefineTimer
	}

	if deduceNonWorkHours {
		deducer = 1
	} else {
		deducer = 0
	}

	deadlineTime, _ := strconv.ParseFloat(strings.Split(deadline, " ")[1], 64)
	deadlineType := strings.Split(deadline, " ")[1]
	if deduceNonWorkHours && deadlineTime-timeTable[deadlineType] < 1 {
		deadlineType = DEADLINE_TYPE_HOURS
		deadlineTime = deadlineTime * (WORK_HOURS_FINAL - WORK_HOURS_INITIAL)
	}

	if timeTable[deadlineType] > deadlineTime {
		req.overdue = true
	}

	req.timer = strconv.FormatFloat(deadlineTime-math.Round(timeTable[deadlineType]), 'f', -1, 64) + " " + deadlineType
	h.request = req

	return nil
}

func (h *Chronos) prepareDeadlineLabel() error {
	var (
		req       = h.request.(ChronosUpdateSingleIssueDeadlineRequest)
		labelName string
	)

	labelName = DEADLINE_LABEL_SIGNATURE + ": " + req.timer
	if req.overdue {
		labelName = DEADLINE_LABEL_OVERDUE
	}

	color := SetColorToLabel(labelName)
	newLabel := &github.Label{
		Name:  &labelName,
		Color: &color,
	}

	_, _, err := h.client.Issues.GetLabel(context.Background(), OWNER, REPO, labelName)
	if err != nil {
		_, _, err := h.client.Issues.CreateLabel(context.Background(), OWNER, REPO, newLabel)
		if err != nil {
			return err
		}
	}

	req.timerLabel = labelName
	h.request = req

	return nil
}

func (h Chronos) updateDeadlineLabel() error {
	var (
		wg          sync.WaitGroup
		req         = h.request.(ChronosUpdateSingleIssueDeadlineRequest)
		labelsNames []string
	)

	labels, _, err := h.client.Issues.ListLabelsByIssue(context.Background(), OWNER, REPO, req.IssueNumber, nil)
	if err != nil {
		return err
	}

	for _, label := range labels {
		if strings.Split(label.GetName(), ": ")[0] == DEADLINE_LABEL_SIGNATURE {
			if strings.Split(label.GetName(), " ")[2] == DEADLINE_TYPE_DAYS || strings.Split(label.GetName(), " ")[2] == DEADLINE_TYPE_HOURS {
				labelsNames = append(labelsNames, label.GetName())
			}
			continue
		}
		if label.GetName() == DEADLINE_LABEL_OVERDUE {
			if label.GetName() != req.LabelName {
				labelsNames = append(labelsNames, label.GetName())
			}
			continue
		}
		if strings.Split(label.GetName(), ": ")[0] == PRIORITY_LABEL_SIGNATURE {
			if label.GetName() != req.LabelName {
				labelsNames = append(labelsNames, label.GetName())
			}
			continue
		}
	}

	for _, label := range labelsNames {
		wg.Add(1)
		go func(issueNumber int, label string) {
			_, e := h.client.Issues.RemoveLabelForIssue(context.Background(), OWNER, REPO, issueNumber, label)
			if e != nil {
				err = ErrUnableToDeleteLabelsFromIssue
				wg.Done()
				return
			}
			wg.Done()
		}(req.IssueNumber, label)
	}
	if err != nil {
		return err
	}

	wg.Wait()

	_, _, e := h.client.Issues.AddLabelsToIssue(context.Background(), OWNER, REPO, req.IssueNumber, []string{req.timerLabel})
	if e != nil {
		return ErrUnableToAddLabelsToIssue
	}

	return err
}

func (h Chronos) UpdateSingleIssueDeadline() error {
	var err error

	err = h.calculateElapsedTime()
	if err != nil {
		return err
	}

	err = h.defineNewDeadline()
	if err != nil {
		return err
	}

	err = h.prepareDeadlineLabel()
	if err != nil {
		return err
	}

	err = h.updateDeadlineLabel()
	if err != nil {
		return err
	}

	return nil
}
