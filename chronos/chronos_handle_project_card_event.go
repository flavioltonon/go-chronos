package chronos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) HandleProjectCardEvent(event interface{}) error {
	var projectCardEvent = event.(*github.ProjectCardEvent)

	fmt.Println(fmt.Sprintf("Event: Project card #%d has been %s", projectCardEvent.GetProjectCard().GetID(), projectCardEvent.GetAction()))

	switch projectCardEvent.GetAction() {
	case "moved":
		projectID, _ := strconv.ParseInt(strings.Split(projectCardEvent.GetProjectCard().GetProjectURL(), "/projects/")[1], 10, 64)
		issueNumber, _ := strconv.Atoi(strings.Split(projectCardEvent.GetProjectCard().GetContentURL(), "/issues/")[1])

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
