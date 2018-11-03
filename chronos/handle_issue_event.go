package chronos

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

func (h *IssuesEventHandler) HandleEvent(event interface{}) error {
	var chronos Chronos

	auth := github.BasicAuthTransport{
		Username: os.Getenv("CHRONOS_GITHUB_LOGIN"),
		Password: os.Getenv("CHRONOS_GITHUB_PASSWORD"),
	}
	chronos.client = github.NewClient(auth.Client())

	issuesEvent := event.(*github.IssuesEvent)
	switch issuesEvent.GetAction() {
	case "opened":
		fmt.Println(fmt.Sprintf("Event: Issue #%d has been %s", issuesEvent.GetIssue().GetNumber(), issuesEvent.GetAction()))
	case "labeled":
		fmt.Println(fmt.Sprintf("Event: Issue #%d has been %s", issuesEvent.GetIssue().GetNumber(), issuesEvent.GetAction()))
		fmt.Println("Label:", issuesEvent.GetLabel().GetName())

		chronos.SetRequest(ChronosUpdateSingleIssueDeadlineRequest{
			IssueNumber: issuesEvent.GetIssue().GetNumber(),
			LabelName:   issuesEvent.GetLabel().GetName(),
			Created:     issuesEvent.GetIssue().GetCreatedAt(),
		})

		if strings.Split(issuesEvent.GetLabel().GetName(), ": ")[0] == "Prioridade" {
			return chronos.UpdateSingleIssueDeadline()
		}
	case "unlabeled":
		fmt.Println(fmt.Sprintf("Event: Issue #%d has been %s", issuesEvent.GetIssue().GetNumber(), issuesEvent.GetAction()))
		fmt.Println("Label:", issuesEvent.GetLabel().GetName())
	default:
		fmt.Println("Event: Issue", issuesEvent.GetAction())
		return nil
	}
	return nil
}
