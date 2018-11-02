package github

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
)

func (h *IssuesEventHandler) HandleEvent(event interface{}) error {
	var chronos Chronos

	issuesEvent := event.(*github.IssuesEvent)
	switch issuesEvent.GetAction() {
	case "opened":
		fmt.Println("Event: Issue", issuesEvent.GetAction())
	case "labeled":
		fmt.Println(fmt.Sprintf("Event: Issue %d has been %s", issuesEvent.GetIssue().GetID(), issuesEvent.GetAction()))
		fmt.Println("Label:", issuesEvent.GetLabel().GetName())

		data := ChronosSetIssueDeadlineRequest{
			IssueNumber: strconv.Itoa(issuesEvent.GetIssue().GetNumber()),
			Label:       issuesEvent.GetLabel().GetName(),
			Created:     issuesEvent.GetIssue().GetCreatedAt(),
		}

		if strings.Split(issuesEvent.GetLabel().GetName(), ": ")[0] == "Prioridade" {
			return chronos.SetIssueDeadline(&data)
		}
	case "unlabeled":
		fmt.Println(fmt.Sprintf("Event: Issue %d has been %s", issuesEvent.GetIssue().GetID(), issuesEvent.GetAction()))
		fmt.Println("Label:", issuesEvent.GetLabel().GetName())
	default:
		fmt.Println("Event: Issue", issuesEvent.GetAction())
		return nil
	}
	return nil
}
