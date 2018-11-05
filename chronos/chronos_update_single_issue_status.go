package chronos

import (
	"context"
	"strings"
	"sync"

	"github.com/flavioltonon/go-github/github"
)

type ChronosUpdateSingleIssueStatusRequest struct {
	IssueNumber int
	ProjectID   int64
	ColumnToID  int64

	client *github.Client

	columnsMap map[int64]string

	issue *github.Issue

	issueState       string
	issueStatusLabel string
}

type ChronosUpdateSingleIssueStatusResponse struct {
}

func (h *ChronosUpdateSingleIssueStatusRequest) mapProjectColumns() error {
	var err error

	projectColumns, _, err := h.client.Projects.ListProjectColumns(context.Background(), h.ProjectID, nil)
	if err != nil {
		return ErrUnableToGetProjectColumns
	}

	columnsMap := make(map[int64]string)
	for _, column := range projectColumns {
		columnsMap[column.GetID()] = column.GetName()
	}

	h.columnsMap = columnsMap

	return nil
}

func (h *ChronosUpdateSingleIssueStatusRequest) getIssue() error {
	issue, _, err := h.client.Issues.Get(context.Background(), OWNER, REPO, h.IssueNumber)
	if err != nil {
		return ErrUnableToGetIssue
	}

	h.issue = issue

	return nil
}

func (h *ChronosUpdateSingleIssueStatusRequest) prepareStatusLabel() error {
	switch h.columnsMap[h.ColumnToID] {
	case COLUMN_BACKLOG:
		h.issueState = STANDARD_ISSUE_STATE_COLUMN_BACKLOG
	case COLUMN_SPRINT_BACKLOG:
		h.issueState = STANDARD_ISSUE_STATE_COLUMN__SPRINT_BACKLOG
	case COLUMN_DEPLOY:
		h.issueState = STANDARD_ISSUE_STATE_COLUMN_DEPLOY
		// h.issueStatusLabel = STATUS_LABEL_DEPLOY
	case COLUMN_DONE:
		h.issueState = STANDARD_ISSUE_STATE_COLUMN_DONE
	default:
		return ErrUnexpectedProjectColumnName
	}

	// color := SetColorToLabel(h.issueStatusLabel)
	// newLabel := &github.Label{
	// 	Name:  &h.issueStatusLabel,
	// 	Color: &color,
	// }

	// if *newLabel.Name != "" {
	// 	_, _, err := h.client.Issues.GetLabel(context.Background(), OWNER, REPO, h.issueStatusLabel)
	// 	if err != nil {
	// 		_, _, err := h.client.Issues.CreateLabel(context.Background(), OWNER, REPO, newLabel)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func (h *ChronosUpdateSingleIssueStatusRequest) updateIssueStatusLabel() error {
	var (
		oldStatusLabels []string
		wg              sync.WaitGroup
		err             error
	)

	for _, label := range h.issue.Labels {
		if strings.Split(label.GetName(), ": ")[0] == STATUS_LABEL_SIGNATURE {
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
		}(h.IssueNumber, label)
	}
	if err != nil {
		return err
	}

	if h.issueStatusLabel == "" {
		wg.Wait()
		return nil
	}

	wg.Add(1)
	go func(issueNumber int, newLabel string) {
		_, _, e := h.client.Issues.AddLabelsToIssue(context.Background(), OWNER, REPO, issueNumber, []string{newLabel})
		if e != nil {
			err = ErrUnableToAddLabelsToIssue
			wg.Done()
			return
		}
		wg.Done()
	}(h.IssueNumber, h.issueStatusLabel)

	wg.Wait()

	return nil
}

func (h *ChronosUpdateSingleIssueStatusRequest) updateIssueState() error {
	_, _, err := h.client.Issues.Edit(context.Background(), OWNER, REPO, h.IssueNumber, &github.IssueRequest{
		State: &h.issueState,
	})
	if err != nil {
		return ErrUnableToUpdateIssueState
	}

	return nil
}

func (h Chronos) UpdateSingleIssueStatus() error {
	var (
		req = h.request.(ChronosUpdateSingleIssueStatusRequest)
		err error
	)

	req.client = h.client

	err = req.mapProjectColumns()
	if err != nil {
		return err
	}

	err = req.getIssue()
	if err != nil {
		return err
	}

	err = req.prepareStatusLabel()
	if err != nil {
		return err
	}

	// err = req.updateIssueStatusLabel()
	// if err != nil {
	// 	return err
	// }

	err = req.updateIssueState()
	if err != nil {
		return err
	}

	return nil
}
