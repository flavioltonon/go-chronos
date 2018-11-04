package chronos

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
)

func (h *ProjectCardEventHandler) HandleEvent(event interface{}) error {
	var (
		chronos          Chronos
		projectCardEvent = event.(*github.ProjectCardEvent)
	)

	auth := github.BasicAuthTransport{
		Username: os.Getenv("CHRONOS_GITHUB_LOGIN"),
		Password: os.Getenv("CHRONOS_GITHUB_PASSWORD"),
	}
	chronos.client = github.NewClient(auth.Client())

	fmt.Println(fmt.Sprintf("Event: Project card #%d has been %s", projectCardEvent.GetProjectCard().GetID(), projectCardEvent.GetAction()))

	switch projectCardEvent.GetAction() {
	case "moved":
		issueNumber, _ := strconv.Atoi(strings.Split(projectCardEvent.GetProjectCard().GetContentURL(), "/issues/")[1])
		projectID, _ := strconv.ParseInt(strings.Split(projectCardEvent.GetProjectCard().GetProjectURL(), "/projects/")[1], 10, 64)

		chronos.SetRequest(ChronosUpdateSingleIssueStatusRequest{
			IssueNumber: issueNumber,
			ProjectID:   projectID,
			ColumnToID:  projectCardEvent.GetProjectCard().GetColumnID(),
		})

		return chronos.UpdateSingleIssueStatus()
	default:
		return nil
	}
}
