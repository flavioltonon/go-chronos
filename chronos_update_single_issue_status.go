package chronos

import (
	"context"
	"flavioltonon/go-chronos/config/column"
	"os"

	"github.com/flavioltonon/go-github/github"
)

type ChronosUpdateSingleIssueStateRequest struct {
	IssueNumber int
	ProjectID   int64
	ColumnToID  int64

	client *github.Client

	issue *github.Issue
}

type ChronosUpdateSingleIssueStateResponse struct{}

func (h *ChronosUpdateSingleIssueStateRequest) updateIssueState() error {
	var columns = column.Columns()

	if _, exists := columns[h.ColumnToID]; !exists {
		return ErrUnexpectedProjectColumnName
	}

	issueState := columns[h.ColumnToID].StandardIssueState()

	_, _, err := h.client.Issues.Edit(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		h.IssueNumber,
		&github.IssueRequest{
			State: &issueState,
		})
	if err != nil {
		return ErrUnableToUpdateIssueState
	}

	return nil
}

func (h Chronos) UpdateSingleIssueState() error {
	var (
		req = h.request.(ChronosUpdateSingleIssueStateRequest)
		err error
	)

	req.client = h.client

	err = req.updateIssueState()
	if err != nil {
		return err
	}

	return nil
}
