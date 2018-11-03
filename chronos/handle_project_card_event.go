package chronos

import (
	"fmt"

	"github.com/google/go-github/github"
)

func (h *ProjectCardEventHandler) HandleEvent(event interface{}) error {
	projectCardEvent := event.(*github.ProjectCardEvent)
	switch projectCardEvent.GetAction() {
	case "moved":
		fmt.Println(fmt.Sprintf("Event: Project card %d has been %s to column %d", projectCardEvent.GetProjectCard().GetID(), projectCardEvent.GetAction(), projectCardEvent.GetProjectCard().GetColumnID()))
	default:
		fmt.Println(fmt.Sprintf("Event: Project card %d has been %s", projectCardEvent.GetProjectCard().GetID(), projectCardEvent.GetAction()))
		return nil
	}
	return nil
}
