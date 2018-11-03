package chronos

import (
	"context"
	"strings"
	"sync"

	"github.com/google/go-github/github"
)

func (h *Chronos) getIssues() error {
	issues, _, err := h.client.Issues.ListByRepo(context.Background(), OWNER, REPO, &github.IssueListByRepoOptions{
		State: "open",
	})
	if err != nil {
		return ErrUnableToGetIssuesFromRepo
	}

	h.issues = issues

	return nil
}

func (h *Chronos) updateIssuesDeadlineLabels() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	for _, issue := range h.issues {
		var (
			labels []string
		)

		for _, label := range issue.Labels {
			if strings.Split(label.GetName(), ": ")[0] == "Prioridade" {
				labels = append(labels, label.GetName())
			}
		}

		for _, label := range labels {
			wg.Add(1)
			go func(issue *github.Issue, label string) {
				_, e := h.client.Issues.RemoveLabelForIssue(context.Background(), OWNER, REPO, issue.GetNumber(), label)
				if e != nil {
					err = ErrUnableToDeleteLabelsFromIssue
					wg.Done()
					return
				}

				_, _, e = h.client.Issues.AddLabelsToIssue(context.Background(), OWNER, REPO, issue.GetNumber(), []string{label})
				if e != nil {
					err = ErrUnableToAddLabelsToIssue
					wg.Done()
					return
				}
				wg.Done()
			}(issue, label)
		}
	}

	if err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func (h *Chronos) UpdateIssuesDeadlines() error {
	var err error

	err = h.getIssues()
	if err != nil {
		return err
	}

	err = h.updateIssuesDeadlineLabels()
	if err != nil {
		return err
	}

	return nil
}
