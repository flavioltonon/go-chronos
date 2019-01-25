package chronos

import (
	"context"
	"encoding/json"
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
	priorityLabel    string
	deadlineLabel    string
	elapsedTime      int
	nonWorkHours     int
	deadline         string
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
	h.deadlineLabel = ""
	h.priorityLabel = ""

	for _, label := range h.issue.Labels {
		if strings.Split(label.GetName(), ": ")[0] == DEADLINE_LABEL_SIGNATURE {
			h.deadlineLabel = label.GetName()
		}
		if strings.Split(label.GetName(), ": ")[0] == PRIORITY_LABEL_SIGNATURE {
			h.priorityLabel = label.GetName()
		}
	}

	if h.priorityLabel == "" {
		return ErrNothingToUpdate
	}

	return nil
}

func (h *ChronosUpdateIssuesDeadlinesRequest) defineNewDeadline() error {
	var (
		deadline           string
		deducer            int
		deduceNonWorkHours bool
	)

	timeTable := make(map[string]int)
	timeTable[DEADLINE_TYPE_HOURS] = h.elapsedTime - deducer*h.nonWorkHours
	timeTable[DEADLINE_TYPE_DAYS] = (h.elapsedTime - deducer*h.nonWorkHours) / (WORK_HOURS_FINAL - WORK_HOURS_INITIAL)

	switch h.priorityLabel {
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

	deadlineTime, _ := strconv.Atoi(strings.Split(deadline, " ")[1])
	deadlineType := strings.Split(deadline, " ")[2]
	if deduceNonWorkHours && deadlineTime-timeTable[deadlineType] < 1 {
		deadlineType = DEADLINE_TYPE_HOURS
		deadlineTime = deadlineTime * (WORK_HOURS_FINAL - WORK_HOURS_INITIAL)
	}

	if timeTable[deadlineType] > deadlineTime {
		h.overdue = true
	}

	h.deadline = strconv.Itoa(deadlineTime-timeTable[deadlineType]) + " " + deadlineType

	return nil
}

func (h *ChronosUpdateIssuesDeadlinesRequest) prepareDeadlineLabel() error {
	var labelName string

	labelName = DEADLINE_LABEL_SIGNATURE + ": " + h.deadline
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
			if strings.Split(label.GetName(), " ")[2] == DEADLINE_TYPE_DAYS || strings.Split(label.GetName(), " ")[2] == DEADLINE_TYPE_HOURS {
				labelsNames = append(labelsNames, label.GetName())
			}
			continue
		}
		if label.GetName() == DEADLINE_LABEL_OVERDUE {
			labelsNames = append(labelsNames, label.GetName())
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
