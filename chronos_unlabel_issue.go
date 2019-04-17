package chronos

import (
	"context"
	"os"

	"github.com/flavioltonon/go-github/github"
)

type ChronosUnlabelIssueRequest struct {
	IssueNumber int
	LabelName   string

	client *github.Client
}

type ChronosUnlabelIssueResponse struct {
}

func (r ChronosUnlabelIssueRequest) removeLabel() error {
	_, err := r.client.Issues.RemoveLabelForIssue(
		context.Background(),
		os.Getenv("GITHUB_REPOSITORY_OWNER"),
		os.Getenv("GITHUB_REPOSITORY_NAME"),
		r.IssueNumber, r.LabelName,
	)
	return err
}

func (h Chronos) UnlabelIssue() error {
	var req = h.request.(ChronosUnlabelIssueRequest)

	req.client = h.client

	return req.removeLabel()
}
