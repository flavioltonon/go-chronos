package chronos

import (
	"flavioltonon/go-chronos/chronos/config/priority"
	"fmt"
	"log"
	"regexp"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) HandleIssuesEvent(event interface{}) error {
	var issuesEvent = event.(*github.IssuesEvent)

	log.Println(fmt.Sprintf("Event: Issue #%d has been %s", issuesEvent.GetIssue().GetNumber(), issuesEvent.GetAction()))

	switch issuesEvent.GetAction() {
	case "labeled":
		log.Println("Label:", issuesEvent.GetLabel().GetName())

		// If an issue gets an priority label, it should have its deadline updated
		for _, priority := range priority.Priorities() {
			if issuesEvent.GetLabel().GetID() == priority.ID() {
				chronos.SetRequest(ChronosUpdateSingleIssueDeadlineRequest{
					IssueNumber: issuesEvent.GetIssue().GetNumber(),
					LabelID:     issuesEvent.GetLabel().GetID(),
					LabelName:   issuesEvent.GetLabel().GetName(),
					Created:     issuesEvent.GetIssue().GetCreatedAt(),
				})
				return chronos.UpdateSingleIssueDeadline()
			}
		}

		if regexp.MustCompile(DEADLINE_LABEL_SIGNATURE).MatchString(issuesEvent.GetLabel().GetName()) == false {
			return nil
		}

		if issuesEvent.GetSender().GetID() != chronos.UserID() {
			// Remove deadline label added by human user
			chronos.SetRequest(ChronosUnlabelIssueRequest{
				IssueNumber: issuesEvent.GetIssue().GetNumber(),
				LabelName:   issuesEvent.GetLabel().GetName(),
			})
			return chronos.UnlabelIssue()
		}

		return nil
	case "unlabeled":
		log.Println("Label:", issuesEvent.GetLabel().GetName())

		if regexp.MustCompile(DEADLINE_LABEL_SIGNATURE).MatchString(issuesEvent.GetLabel().GetName()) == false &&
			regexp.MustCompile(PRIORITY_LABEL_SIGNATURE).MatchString(issuesEvent.GetLabel().GetName()) == false {
			return nil
		}

		if issuesEvent.GetSender().GetID() != chronos.UserID() {
			chronos.SetRequest(ChronosRelabelIssueRequest{
				IssueNumber: issuesEvent.GetIssue().GetNumber(),
				LabelName:   issuesEvent.GetLabel().GetName(),
			})
			return chronos.RelabelIssue()
		}

		return nil
	case "closed":
		// Closed by human user
		if issuesEvent.GetSender().GetID() != chronos.UserID() {
			chronos.SetRequest(ChronosUpdateIssueRequest{
				IssueNumber: issuesEvent.GetIssue().GetNumber(),
				IssueState:  "open",
			})
			return chronos.UpdateIssue()
		}

		return nil
	case "reopened":
		// Reopened by Chronos
		if issuesEvent.GetSender().GetID() == chronos.UserID() {
			for _, label := range issuesEvent.GetIssue().Labels {
				for _, priority := range priority.Priorities() {
					if label.GetID() == priority.ID() {
						chronos.SetRequest(ChronosUpdateSingleIssueDeadlineRequest{
							IssueNumber: issuesEvent.GetIssue().GetNumber(),
							LabelName:   label.GetName(),
							Created:     issuesEvent.GetIssue().GetCreatedAt(),
						})
						return chronos.UpdateSingleIssueDeadline()
					}
				}
			}
		}

		// Reopened by human user
		if issuesEvent.GetSender().GetID() != chronos.UserID() {
			if issuesEvent.GetSender().GetID() != chronos.UserID() {
				chronos.SetRequest(ChronosUpdateIssueRequest{
					IssueNumber: issuesEvent.GetIssue().GetNumber(),
					IssueState:  "closed",
				})
				return chronos.UpdateIssue()
			}
		}

		return nil
	default:
		return nil
	}
}
