package chronos

import (
	"context"
	"encoding/json"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flavioltonon/go-github/github"
	"github.com/go-resty/resty"
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
	case PRIORITY_LABEL_PRIORITY_LOW:
		deadline = DEADLINE_LABEL_PRIORITY_LOW
		deduceNonWorkHours = true
	case PRIORITY_LABEL_PRIORITY_MEDIUM:
		deadline = DEADLINE_LABEL_PRIORITY_MEDIUM
		deduceNonWorkHours = true
	case PRIORITY_LABEL_PRIORITY_HIGH:
		deadline = DEADLINE_LABEL_PRIORITY_HIGH
		deduceNonWorkHours = true
	case PRIORITY_LABEL_PRIORITY_VERY_HIGH:
		deadline = DEADLINE_LABEL_PRIORITY_VERY_HIGH
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

	_, _, err := h.client.Issues.GetLabel(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, labelName)
	if err != nil {
		_, _, err := h.client.Issues.CreateLabel(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, newLabel)
		if err != nil {
			return err
		}
	}

	h.timerLabel = labelName

	return nil
}

func (h ChronosUpdateSingleIssueDeadlineRequest) updateDeadlineLabel() error {
	var labelsNames = make([]string, 0)

	labels, _, err := h.client.Issues.ListLabelsByIssue(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, h.IssueNumber, nil)
	if err != nil {
		return err
	}

	labelsNames = append(labelsNames, h.LabelName)
	labelsNames = append(labelsNames, h.timerLabel)

	for _, label := range labels {
		if regexp.MustCompile(DEADLINE_LABEL_SIGNATURE).MatchString(label.GetName()) {
			continue
		}

		if label.GetName() == DEADLINE_LABEL_OVERDUE {
			continue
		}

		if regexp.MustCompile(PRIORITY_LABEL_SIGNATURE).MatchString(label.GetName()) {
			continue
		}

		labelsNames = append(labelsNames, label.GetName())
	}

	_, _, e := h.client.Issues.ReplaceLabelsForIssue(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, h.IssueNumber, labelsNames)
	if e != nil {
		return ErrUnableToReplaceLabelsFromIssue
	}
	if err != nil {
		return err
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
