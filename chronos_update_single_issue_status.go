package chronos

import (
	"context"
	"fmt"
	"log"
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

func (h *ChronosUpdateSingleIssueStateRequest) validate() error {
	if h.IssueNumber == 0 {
		return fmt.Errorf("invalid issue number: %v", h.IssueNumber)
	}
	if h.ProjectID == 0 {
		return fmt.Errorf("invalid project ID: %v", h.ProjectID)
	}
	if h.ColumnToID == 0 {
		return fmt.Errorf("invalid column-to ID: %v", h.ColumnToID)
	}

	return nil
}

func (h *ChronosUpdateSingleIssueStateRequest) updateIssueState() error {
	var columns = Columns()

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
		log.Println("failed to edit issue:", err)
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

	err = req.validate()
	if err != nil {
		log.Println("failed to validate request:", h.request)
		return err
	}

	err = req.updateIssueState()
	if err != nil {
		log.Println("failed to update issue state:", h.request)
		return err
	}

	return nil
}
