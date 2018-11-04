package chronos

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty"
	"github.com/google/go-github/github"
)

type ChronosUpdateSingleIssueDeadlineRequest struct {
	IssueNumber int
	LabelName   string
	Created     time.Time

	client *github.Client

	holidays     Holidays
	elapsedTime  float64
	nonWorkHours float64
	timer        string
	overdue      bool
	timerLabel   string
}

type ChronosUpdateSingleIssueDeadlineResponse struct {
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) getHolidays() error {
	loc, _ := time.LoadLocation(STANDARD_TIME_LOCATION)
	now := time.Now().In(loc)

	query := map[string]string{
		"country": "BR",
		"year":    strconv.Itoa(now.Year()),
	}

	resp, err := resty.R().SetQueryParams(query).Get(HOLIDAY_API_URL)
	if err != nil {
		return ErrUnableToSendGetHolidaysRequest
	}

	var res holidaysResponse
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return ErrUnableToUnmarshalGetHolidaysResponse
	}

	h.holidays = res.Holidays

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) calculateElapsedTime() error {
	var (
		nonWorkHours float64
		weekendHours float64
		holidayHours float64
	)

	loc, _ := time.LoadLocation(STANDARD_TIME_LOCATION)
	now := time.Now().In(loc)
	created := h.Created.In(loc)
	elapsedTime := int(math.Round(now.Sub(created).Hours()))

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
		_, exists := h.holidays[created.Add(time.Duration(t)*time.Hour).Format("2006-01-02")]
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

	h.elapsedTime = now.Sub(created).Hours() - weekendHours - holidayHours
	h.nonWorkHours = nonWorkHours

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) defineNewDeadline() error {
	var (
		deadline           string
		deducer            float64
		deduceNonWorkHours bool
	)

	timeTable := make(map[string]float64)
	timeTable[DEADLINE_TYPE_HOURS] = h.elapsedTime - deducer*h.nonWorkHours
	timeTable[DEADLINE_TYPE_DAYS] = (h.elapsedTime - deducer*h.nonWorkHours) / (WORK_HOURS_FINAL - WORK_HOURS_INITIAL)

	switch h.LabelName {
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
	deadlineType := strings.Split(deadline, " ")[2]
	if deduceNonWorkHours && deadlineTime-timeTable[deadlineType] < 1 {
		deadlineType = DEADLINE_TYPE_HOURS
		deadlineTime = deadlineTime * (WORK_HOURS_FINAL - WORK_HOURS_INITIAL)
	}

	if timeTable[deadlineType] > deadlineTime {
		h.overdue = true
	}

	h.timer = strconv.FormatFloat(deadlineTime-math.Round(timeTable[deadlineType]), 'f', -1, 64) + " " + deadlineType

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) prepareDeadlineLabel() error {
	var (
		labelName string
	)

	labelName = DEADLINE_LABEL_SIGNATURE + ": " + h.timer
	if h.overdue {
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

	h.timerLabel = labelName

	return nil
}

func (h ChronosUpdateSingleIssueDeadlineRequest) updateDeadlineLabel() error {
	var (
		wg          sync.WaitGroup
		labelsNames []string
	)

	labels, _, err := h.client.Issues.ListLabelsByIssue(context.Background(), OWNER, REPO, h.IssueNumber, nil)
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
			if label.GetName() != h.LabelName {
				labelsNames = append(labelsNames, label.GetName())
			}
			continue
		}
		if strings.Split(label.GetName(), ": ")[0] == PRIORITY_LABEL_SIGNATURE {
			if label.GetName() != h.LabelName {
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
		}(h.IssueNumber, label)
	}
	if err != nil {
		return err
	}

	wg.Wait()

	_, _, e := h.client.Issues.AddLabelsToIssue(context.Background(), OWNER, REPO, h.IssueNumber, []string{h.timerLabel})
	if e != nil {
		return ErrUnableToAddLabelsToIssue
	}

	return err
}

func (h Chronos) UpdateSingleIssueDeadline() error {
	var (
		req = h.request.(ChronosUpdateSingleIssueDeadlineRequest)
		err error
	)

	req.client = h.client

	err = req.getHolidays()
	if err != nil {
		return err
	}

	err = req.calculateElapsedTime()
	if err != nil {
		return err
	}

	err = req.defineNewDeadline()
	if err != nil {
		return err
	}

	err = req.prepareDeadlineLabel()
	if err != nil {
		return err
	}

	err = req.updateDeadlineLabel()
	if err != nil {
		return err
	}

	return nil
}
