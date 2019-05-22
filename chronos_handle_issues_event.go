package chronos

import (
	"fmt"
	"log"
	"regexp"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) HandleIssuesEvent(event *github.IssuesEvent) error {
	log.Println(fmt.Sprintf("Event: Issue #%d has been %s", event.GetIssue().GetNumber(), event.GetAction()))

	switch event.GetAction() {
	case "labeled":
		log.Println("Label:", event.GetLabel().GetName())

		if regexp.MustCompile(DEADLINE_LABEL_SIGNATURE).MatchString(event.GetLabel().GetName()) == false &&
			regexp.MustCompile(PRIORITY_LABEL_SIGNATURE).MatchString(event.GetLabel().GetName()) == false {
			return nil
		}

		// If an issue gets an priority label, it should have its deadline updated
		if _, exists := chronos.priorities[event.GetLabel().GetID()]; exists {
			chronos.SetRequest(ChronosUpdateSingleIssueDeadlineRequest{
				IssueNumber: event.GetIssue().GetNumber(),
				LabelID:     event.GetLabel().GetID(),
				LabelName:   event.GetLabel().GetName(),
				Created:     event.GetIssue().GetCreatedAt(),
			})
			return chronos.UpdateSingleIssueDeadline()
		}

		if event.GetSender().GetID() != chronos.UserID() {
			// Remove deadline label added by human user
			chronos.SetRequest(ChronosUnlabelIssueRequest{
				IssueNumber: event.GetIssue().GetNumber(),
				LabelName:   event.GetLabel().GetName(),
			})
			return chronos.UnlabelIssue()
		}
	case "unlabeled":
		log.Println("Label:", event.GetLabel().GetName())

		if regexp.MustCompile(DEADLINE_LABEL_SIGNATURE).MatchString(event.GetLabel().GetName()) == false &&
			regexp.MustCompile(PRIORITY_LABEL_SIGNATURE).MatchString(event.GetLabel().GetName()) == false {
			return nil
		}

		if event.GetSender().GetID() != chronos.UserID() {
			chronos.SetRequest(ChronosRelabelIssueRequest{
				IssueNumber: event.GetIssue().GetNumber(),
				LabelName:   event.GetLabel().GetName(),
			})
			return chronos.RelabelIssue()
		}
	case "closed":
		// Closed by human user
		if event.GetSender().GetID() != chronos.UserID() {
			chronos.SetRequest(ChronosUpdateIssueRequest{
				IssueNumber: event.GetIssue().GetNumber(),
				IssueState:  "open",
			})
			return chronos.UpdateIssue()
		}
	case "reopened":
		// Reopened by Chronos
		if event.GetSender().GetID() == chronos.UserID() {
			for _, label := range event.GetIssue().Labels {
				if _, exists := chronos.priorities[label.GetID()]; exists {
					chronos.SetRequest(ChronosUpdateSingleIssueDeadlineRequest{
						IssueNumber: event.GetIssue().GetNumber(),
						LabelID:     label.GetID(),
						LabelName:   label.GetName(),
						Created:     event.GetIssue().GetCreatedAt(),
					})
					return chronos.UpdateSingleIssueDeadline()
				}
			}
		}

		// Reopened by human user
		if event.GetSender().GetID() != chronos.UserID() {
			chronos.SetRequest(ChronosUpdateIssueRequest{
				IssueNumber: event.GetIssue().GetNumber(),
				IssueState:  "closed",
			})
			return chronos.UpdateIssue()
		}
	}

	return nil
}
