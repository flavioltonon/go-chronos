package chronos

import (
	"context"
	"os"

	"github.com/flavioltonon/go-github/github"
)

type ChronosRelabelIssueRequest struct {
	IssueNumber int
	LabelName   string

	client *github.Client
}

type ChronosRelabelIssueResponse struct {
}

func (r ChronosRelabelIssueRequest) readdLabel() error {
	_, _, err := r.client.Issues.AddLabelsToIssue(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		r.IssueNumber, []string{r.LabelName},
	)
	return err
}

func (h Chronos) RelabelIssue() error {
	var req = h.request.(ChronosRelabelIssueRequest)

	req.client = h.client

	return req.readdLabel()
}
