package chronos

import (
	"fmt"

	"github.com/google/go-github/github"
)

func (chronos Chronos) HandleLabelEvent(event interface{}) error {
	var labelEvent = event.(*github.LabelEvent)

	fmt.Println(fmt.Sprintf("Event: Label %s has been %s", labelEvent.GetLabel().GetName(), labelEvent.GetAction()))

	return nil
}

func SetColorToLabel(name string) string {
	if name == DEADLINE_LABEL_OVERDUE {
		return "4b21c6"
	}

	return "ffffff"
}
