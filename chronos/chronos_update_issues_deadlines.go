package chronos

import (
	"context"
	"encoding/json"
	"flavioltonon/go-chronos/chronos/config/priority"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flavioltonon/go-github/github"
	"github.com/go-resty/resty"
)

type ChronosUpdateIssuesDeadlinesRequest struct {
	client *github.Client

	issue  *github.Issue
	issues []*github.Issue

	holidays         Holidays
	priority         priority.Priority
	currentDeadline  string
	elapsedTime      int
	nonWorkHours     int
	newDeadline      priority.Deadline
	newDeadlineLabel string
	overdue          bool
}

type ChronosUpdateIssuesDeadlinesResponse struct {
}

func (h *ChronosUpdateIssuesDeadlinesRequest) getHolidays() error {
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

func (h *ChronosUpdateIssuesDeadlinesRequest) getRepoIssues() error {
	issues, _, err := h.client.Issues.ListByRepo(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, &github.IssueListByRepoOptions{
		State: "open",
	})
	if err != nil {
		return ErrUnableToGetIssuesFromRepo
	}

	h.issues = issues

	return nil
}

func (h *ChronosUpdateIssuesDeadlinesRequest) calculateElapsedTime() error {
	var (
		nonWorkHours int
		weekendHours int
		holidayHours int
	)

	loc, _ := time.LoadLocation(STANDARD_TIME_LOCATION)
	now := time.Now().In(loc)
	created := h.issue.GetCreatedAt().In(loc)

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

func (h *ChronosUpdateIssuesDeadlinesRequest) findLabels() error {
	var hasPriorityLabel bool

	for _, label := range h.issue.Labels {
		if strings.Split(label.GetName(), ": ")[0] == DEADLINE_LABEL_SIGNATURE {
			h.currentDeadline = label.GetName()
			continue
		}
		if _, exists := priority.Priorities()[label.GetID()]; exists {
			h.priority = priority.Priorities()[label.GetID()]
			hasPriorityLabel = true
		}
	}
	if false == hasPriorityLabel {
		return ErrNothingToUpdate
	}

	return nil
}

func (h *ChronosUpdateIssuesDeadlinesRequest) defineNewDeadline() error {
	var (
		deadline = h.priority.Deadline()
		t        = h.elapsedTime
	)

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

func (h *ChronosUpdateIssuesDeadlinesRequest) prepareDeadlineLabel() error {
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

	h.newDeadlineLabel = labelName

	return nil
}

func (h *ChronosUpdateIssuesDeadlinesRequest) updateDeadlineLabel() error {
	var (
		wg          sync.WaitGroup
		labelsNames []string
	)

	labels, _, err := h.client.Issues.ListLabelsByIssue(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, h.issue.GetNumber(), nil)
	if err != nil {
		return err
	}

	for _, label := range labels {
		if strings.Split(label.GetName(), ": ")[0] == DEADLINE_LABEL_SIGNATURE {
			if label.GetName() == DEADLINE_LABEL_OVERDUE {
				labelsNames = append(labelsNames, label.GetName())
				continue
			}
			if strings.Split(label.GetName(), " ")[2] == DEADLINE_TYPE_DAYS || strings.Split(label.GetName(), " ")[2] == DEADLINE_TYPE_HOURS {
				labelsNames = append(labelsNames, label.GetName())
			}
			continue
		}
	}

	for _, label := range labelsNames {
		wg.Add(1)
		go func(issueNumber int, label string) {
			_, e := h.client.Issues.RemoveLabelForIssue(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, issueNumber, label)
			if e != nil {
				err = ErrUnableToDeleteLabelsFromIssue
				wg.Done()
				return
			}
			wg.Done()
		}(h.issue.GetNumber(), label)
	}
	if err != nil {
		return err
	}

	wg.Wait()

	_, _, e := h.client.Issues.AddLabelsToIssue(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, h.issue.GetNumber(), []string{h.newDeadlineLabel})
	if e != nil {
		return ErrUnableToAddLabelsToIssue
	}

	return err
}

func (h *Chronos) UpdateIssuesDeadlines() error {
	var (
		req = h.request.(ChronosUpdateIssuesDeadlinesRequest)
		err error
	)

	req.client = h.client

	err = req.getHolidays()
	if err != nil {
		return err
	}

	err = req.getRepoIssues()
	if err != nil {
		return err
	}

	for i, issue := range req.issues {
		req.issue = issue

		log.Println(fmt.Sprintf("Updating issue %d out of %d...", i+1, len(req.issues)))

		err = req.calculateElapsedTime()
		if err != nil {
			continue
		}

		err = req.findLabels()
		if err != nil {
			continue
		}

		err = req.defineNewDeadline()
		if err != nil {
			continue
		}

		err = req.prepareDeadlineLabel()
		if err != nil {
			continue
		}

		err = req.updateDeadlineLabel()
		if err != nil {
			continue
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}
