package chronos

import (
	"fmt"

	"github.com/flavioltonon/go-github/github"
)

func (chronos Chronos) HandleLabelEvent(event *github.LabelEvent) error {
	fmt.Println(fmt.Sprintf("Event: Label %s has been %s", event.GetLabel().GetName(), event.GetAction()))

	return nil
}

func SetColorToLabel(name string) string {
	if name == DEADLINE_LABEL_OVERDUE {
		return "4b21c6"
	}

	return "ffffff"
}
