package chronos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
)

func (h *ProjectCardEventHandler) HandleEvent(event interface{}) error {
	var (
		chronos          Chronos
		projectCardEvent = event.(*github.ProjectCardEvent)
	)

	switch projectCardEvent.GetAction() {
	case "moved":
		fmt.Println(fmt.Sprintf("Event: Project card #%d has been %s", projectCardEvent.GetProjectCard().GetID(), projectCardEvent.GetAction()))

		issueNumber, _ := strconv.Atoi(strings.Split(projectCardEvent.GetProjectCard().GetContentURL(), "/issues/")[1])
		data := ChronosUpdateSingleIssueStatusRequest{
			IssueNumber: issueNumber,
			ColumnToID:  projectCardEvent.GetProjectCard().GetColumnID(),
		}
		return chronos.UpdateSingleIssueStatus(&data)
	default:
		fmt.Println(fmt.Sprintf("Event: Project card #%d has been %s", projectCardEvent.GetProjectCard().GetID(), projectCardEvent.GetAction()))
		return nil
	}
}
