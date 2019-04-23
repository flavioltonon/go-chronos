package chronos

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/flavioltonon/go-github/github"
)

type ChronosUpdateIssuesDeadlinesRequest struct {
	client *github.Client

	issues []*github.Issue
}

type ChronosUpdateIssuesDeadlinesResponse struct {
}

func (h *ChronosUpdateIssuesDeadlinesRequest) getRepoIssues() error {
	var issues = make([]*github.Issue, 0)

	var lastPage = 1
	for page := 1; page <= lastPage; page++ {
		i, resp, err := h.client.Issues.ListByRepo(
			context.Background(),
			os.Getenv("GITHUB_REPOSITORY_OWNER"),
			os.Getenv("GITHUB_REPOSITORY_NAME"),
			&github.IssueListByRepoOptions{
				State: "open",
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: 30,
				},
			},
		)
		if err != nil {
			return ErrUnableToGetIssuesFromRepo
		}

		lastPage = resp.LastPage

		issues = append(issues, i...)
	}

	h.issues = issues

	return nil
}

func (h *Chronos) UpdateIssuesDeadlines() error {
	var (
		req = h.request.(ChronosUpdateIssuesDeadlinesRequest)
		err error
	)

	req.client = h.client

	err = req.getRepoIssues()
	if err != nil {
		return err
	}

	for i, issue := range req.issues {
		log.Println(fmt.Sprintf("Updating issue %d out of %d...", i+1, len(req.issues)))

		for _, label := range issue.Labels {
			if regexp.MustCompile(PRIORITY_LABEL_SIGNATURE).MatchString(label.GetName()) {
				h.SetRequest(ChronosUpdateSingleIssueDeadlineRequest{
					IssueNumber: issue.GetNumber(),
					LabelID:     label.GetID(),
					LabelName:   label.GetName(),
					Created:     issue.GetCreatedAt(),
				})

				err := h.UpdateSingleIssueDeadline()
				if err != nil {
					log.Println(fmt.Sprintf("failed to update issue #%d deadline", issue.GetNumber()))
				}

				// Sleep to avoid over-requesting to GitHub's API; minimum sleeping period necessary: 5 seconds
				time.Sleep(7 * time.Second)

				break
			}
		}

	}

	return nil
}
