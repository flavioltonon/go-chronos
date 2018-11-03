package chronos

import (
	"context"
	"strings"
	"sync"

	"github.com/google/go-github/github"
)

type ChronosUpdateSingleIssueStatusRequest struct {
	IssueNumber int
	ProjectID   int64
	ColumnToID  int64

	columnsMap map[int64]string

	issue *github.Issue

	issueState       string
	issueStatusLabel string
}

type ChronosUpdateSingleIssueStatusResponse struct {
}

func (h *Chronos) mapProjectColumns() error {
	var (
		req = h.request.(ChronosUpdateSingleIssueStatusRequest)
		err error
	)

	projectColumns, _, err := h.client.Projects.ListProjectColumns(context.Background(), req.ProjectID, nil)
	if err != nil {
		return ErrUnableToGetProjectColumns
	}

	columnsMap := make(map[int64]string)
	for _, column := range projectColumns {
		columnsMap[column.GetID()] = column.GetName()
	}

	req.columnsMap = columnsMap
	h.request = req

	return nil
}

func (h *Chronos) getIssue() error {
	var req = h.request.(ChronosUpdateSingleIssueStatusRequest)

	issue, _, err := h.client.Issues.Get(context.Background(), OWNER, REPO, req.IssueNumber)
	if err != nil {
		return ErrUnableToGetIssue
	}

	req.issue = issue
	h.request = req

	return nil
}

func (h *Chronos) prepareStatusLabel() error {
	var req = h.request.(ChronosUpdateSingleIssueStatusRequest)

	switch req.columnsMap[req.ColumnToID] {
	case COLUMN_BACKLOG:
		req.issueState = "open"
	case COLUMN_SPRINTBACKLOG:
		req.issueState = "open"
	case COLUMN_DEPLOY:
		req.issueState = "closed"
		req.issueStatusLabel = DEPLOY_STATUS_LABEL
	case COLUMN_DONE:
		req.issueState = "closed"
	default:
		return ErrUnexpectedColumnName
	}

	color := SetColorToLabel(req.issueStatusLabel)
	newLabel := &github.Label{
		Name:  &req.issueStatusLabel,
		Color: &color,
	}

	if *newLabel.Name != "" {
		_, _, err := h.client.Issues.GetLabel(context.Background(), OWNER, REPO, req.issueStatusLabel)
		if err != nil {
			_, _, err := h.client.Issues.CreateLabel(context.Background(), OWNER, REPO, newLabel)
			if err != nil {
				return err
			}
		}
	}

	h.request = req

	return nil
}

func (h *Chronos) updateIssueStatusLabel() error {
	var (
		req             = h.request.(ChronosUpdateSingleIssueStatusRequest)
		oldStatusLabels []string
		wg              sync.WaitGroup
		err             error
	)

	for _, label := range req.issue.Labels {
		if strings.Split(label.GetName(), ": ")[0] == "Status" {
			oldStatusLabels = append(oldStatusLabels, label.GetName())
		}
	}

	for _, label := range oldStatusLabels {
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

	if req.issueStatusLabel == "" {
		return nil
	}

	go func(issueNumber int, newLabel string) {
		_, _, e := h.client.Issues.AddLabelsToIssue(context.Background(), OWNER, REPO, issueNumber, []string{newLabel})
		if e != nil {
			err = ErrUnableToAddLabelsToIssue
			return
		}
	}(req.IssueNumber, req.issueStatusLabel)

	return nil
}

func (h *Chronos) updateIssueState() error {
	var req = h.request.(ChronosUpdateSingleIssueStatusRequest)

	_, _, err := h.client.Issues.Edit(context.Background(), OWNER, REPO, req.IssueNumber, &github.IssueRequest{
		State: &req.issueState,
	})
	if err != nil {
		return ErrUnableToUpdateIssueState
	}

	return nil
}

func (h Chronos) UpdateSingleIssueStatus() error {
	var err error

	err = h.mapProjectColumns()
	if err != nil {
		return err
	}

	err = h.getIssue()
	if err != nil {
		return err
	}

	err = h.prepareStatusLabel()
	if err != nil {
		return err
	}

	err = h.updateIssueStatusLabel()
	if err != nil {
		return err
	}

	err = h.updateIssueState()
	if err != nil {
		return err
	}

	return nil
}
