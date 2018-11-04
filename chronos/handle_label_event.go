package chronos

import (
	"fmt"
	"os"

	"github.com/google/go-github/github"
)

func (h *LabelEventHandler) HandleEvent(event interface{}) error {
	var chronos Chronos

	auth := github.BasicAuthTransport{
		Username: os.Getenv("CHRONOS_GITHUB_LOGIN"),
		Password: os.Getenv("CHRONOS_GITHUB_PASSWORD"),
	}
	chronos.client = github.NewClient(auth.Client())

	labelEvent := event.(*github.LabelEvent)

	fmt.Println(fmt.Sprintf("Event: Label %s has been %s", labelEvent.GetLabel().GetName(), labelEvent.GetAction()))

	return nil
}

func SetColorToLabel(name string) string {
	if name == DEADLINE_LABEL_OVERDUE {
		return "4b21c6"
	}

	return "ffffff"
}
