package chronos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) HandleProjectCardEvent(event *github.ProjectCardEvent) error {
	fmt.Println(fmt.Sprintf("Event: Project card #%d has been %s", event.GetProjectCard().GetID(), event.GetAction()))

	switch event.GetAction() {
	case "moved":
		projectID, _ := strconv.ParseInt(strings.Split(event.GetProjectCard().GetProjectURL(), "/projects/")[1], 10, 64)
		issueNumber, _ := strconv.Atoi(strings.Split(event.GetProjectCard().GetContentURL(), "/issues/")[1])

		chronos.SetRequest(ChronosUpdateSingleIssueStateRequest{
			IssueNumber: issueNumber,
			ProjectID:   projectID,
			ColumnToID:  event.GetProjectCard().GetColumnID(),
		})

		return chronos.UpdateSingleIssueState()
	default:
		return nil
	}
}
