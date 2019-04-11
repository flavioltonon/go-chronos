package chronos

import (
	"context"
	"encoding/json"
	"flavioltonon/go-chronos/chronos/config/priority"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/flavioltonon/go-github/github"
	"github.com/go-resty/resty"
)

type ChronosUpdateSingleIssueDeadlineRequest struct {
	IssueNumber int
	LabelID     int64
	LabelName   string
	Created     time.Time

	client *github.Client

	holidays     Holidays
	elapsedTime  int
	priority     priority.Priority
	newDeadline  priority.Deadline
	nonWorkHours int
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
		nonWorkHours int
		weekendHours int
		holidayHours int
	)

	loc, _ := time.LoadLocation(STANDARD_TIME_LOCATION)
	now := time.Now().In(loc)
	created := h.Created.In(loc)

	hoursElapsed := int(math.Round(now.Sub(created).Hours()))

	for t := 0; t < hoursElapsed; t++ {
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

	h.elapsedTime = hoursElapsed - weekendHours - holidayHours
	h.nonWorkHours = nonWorkHours

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) defineNewDeadline() error {
	var (
		deadline priority.Deadline
		t        = h.elapsedTime
	)

	p, exists := priority.NewPriority(h.LabelID)
	if false == exists {
		return ErrInvalidPriority
	}
	deadline = p.Deadline()

	if deadline.DeduceNonWorkHours {
		t -= h.nonWorkHours
	}

	if t < 0 {
		h.overdue = true
		return nil
	}

	if t <= 24 {
		h.newDeadline = priority.Deadline{
			Duration: t,
			Unit:     DEADLINE_TYPE_HOURS,
		}
		return nil
	}

	if deadline.Unit == DEADLINE_TYPE_DAYS {
		t /= WORK_HOURS_FINAL - WORK_HOURS_INITIAL
	}

	h.newDeadline = priority.Deadline{
		Duration: t,
		Unit:     deadline.Unit,
	}

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) prepareDeadlineLabel() error {
	var labelName = DEADLINE_LABEL_OVERDUE

	if false == h.overdue {
		labelName = DEADLINE_LABEL_SIGNATURE + ": " + strconv.Itoa(h.newDeadline.Duration) + " " + h.newDeadline.Unit
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
