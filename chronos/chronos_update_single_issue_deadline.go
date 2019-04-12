package chronos

import (
	"context"
	"encoding/json"
	"flavioltonon/go-chronos/chronos/config/priority"
	"math"
	"os"
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

	holidays    Holidays
	elapsedTime int
	priority    priority.Priority
	newDeadline priority.Deadline
	timer       string
	overdue     bool
	timerLabel  string
}

type ChronosUpdateSingleIssueDeadlineResponse struct {
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) getHolidays() error {
	r, err := resty.R().SetQueryParams(map[string]string{
		"country": "BR",
		"year":    strconv.Itoa(time.Now().Local().Year()),
	}).Get(os.Getenv("HOLIDAY_API_URL"))
	if err != nil {
		return ErrUnableToSendGetHolidaysRequest
	}
	var resp holidaysResponse

	err = json.Unmarshal(r.Body(), &resp)
	if err != nil {
		return ErrUnableToUnmarshalGetHolidaysResponse
	}

	h.holidays = resp.Holidays

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) calculateElapsedTime() error {
	var (
		weekendHours int
		holidayHours int
	)

	created := h.Created.Local()
	now := time.Now().Local()

	// Calculates the difference in hours between the issue creation date and the current time (round-up)
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
	}

	h.elapsedTime = hoursElapsed - weekendHours - holidayHours

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) defineNewDeadline() error {
	p, exists := priority.NewPriority(h.LabelID)
	if false == exists {
		return ErrPriorityNotRegistered
	}

	deadline := p.Deadline()
	baseTime := deadline.Duration
	if deadline.Unit == DEADLINE_TYPE_DAYS {
		baseTime *= 24
	}

	delta := baseTime - h.elapsedTime
	if delta <= 0 {
		h.overdue = true
		return nil
	}
	if delta <= 24 {
		h.newDeadline = priority.Deadline{
			Duration: delta,
			Unit:     DEADLINE_TYPE_HOURS,
		}
		return nil
	}

	if deadline.Unit == DEADLINE_TYPE_DAYS {
		delta /= 24
	}

	h.newDeadline = priority.Deadline{
		Duration: delta,
		Unit:     deadline.Unit,
	}

	return nil
}

func (h *ChronosUpdateSingleIssueDeadlineRequest) prepareDeadlineLabel() error {
	var (
		labelName     = DEADLINE_LABEL_OVERDUE
		labelDuration = strconv.Itoa(h.newDeadline.Duration)
		labelUnit     = h.newDeadline.Unit
	)

	if "1" == labelDuration {
		labelUnit = labelUnit[:len(labelUnit)-1]
	}
	if false == h.overdue {
		labelName = DEADLINE_LABEL_SIGNATURE + ": " + labelDuration + " " + labelUnit
	}

	color := SetColorToLabel(labelName)
	newLabel := &github.Label{
		Name:  &labelName,
		Color: &color,
	}

	_, _, err := h.client.Issues.GetLabel(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		labelName,
	)
	if err != nil {
		_, _, err := h.client.Issues.CreateLabel(
			context.Background(),
			os.Getenv("GITHUB_REPOSITORY_OWNER"),
			os.Getenv("GITHUB_REPOSITORY_NAME"),
			newLabel,
		)
		if err != nil {
			return err
		}
	}

	h.timerLabel = labelName

	return nil
}

func (h ChronosUpdateSingleIssueDeadlineRequest) updateDeadlineLabel() error {
	var labelsNames = make([]string, 0)

	labels, _, err := h.client.Issues.ListLabelsByIssue(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		h.IssueNumber,
		nil,
	)
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

	_, _, e := h.client.Issues.ReplaceLabelsForIssue(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		h.IssueNumber,
		labelsNames,
	)
	if e != nil {
		return ErrUnableToReplaceLabelsFromIssue
	}

	return nil
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
