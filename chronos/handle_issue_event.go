package chronos

import (
	"fmt"
	"log"
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

	log.Println(fmt.Sprintf("Event: Issue #%d has been %s", issuesEvent.GetIssue().GetNumber(), issuesEvent.GetAction()))

	switch issuesEvent.GetAction() {
	case "opened":
		return nil
	case "labeled":
		log.Println("Label:", issuesEvent.GetLabel().GetName())

		if strings.Split(issuesEvent.GetLabel().GetName(), ": ")[0] == PRIORITY_LABEL_SIGNATURE {
			chronos.SetRequest(ChronosUpdateSingleIssueDeadlineRequest{
				IssueNumber: issuesEvent.GetIssue().GetNumber(),
				LabelName:   issuesEvent.GetLabel().GetName(),
				Created:     issuesEvent.GetIssue().GetCreatedAt(),
			})

			return chronos.UpdateSingleIssueDeadline()
		}

		return nil
	case "unlabeled":
		log.Println("Label:", issuesEvent.GetLabel().GetName())
		return nil
	default:
		return nil
	}
}
