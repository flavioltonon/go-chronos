package chronos

import (
	"context"

	"github.com/google/go-github/github"
)

type ChronosUpdateIssueRequest struct {
	IssueNumber int
	IssueState  string

	client *github.Client
}

type ChronosUpdateIssueResponse struct {
}

func (r ChronosUpdateIssueRequest) updateIssue() error {
	_, _, err := r.client.Issues.Edit(context.Background(), GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME, r.IssueNumber, &github.IssueRequest{
		State: &r.IssueState,
	})
	return err
}

func (h Chronos) UpdateIssue() error {
	var req = h.request.(ChronosUpdateIssueRequest)

	req.client = h.client

	return req.updateIssue()
}
