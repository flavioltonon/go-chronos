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
	switch labelEvent.GetAction() {
	case "created":
		fmt.Println(fmt.Sprintf("Event: New label %s has been %s", labelEvent.GetLabel().GetName(), labelEvent.GetAction()))
	default:
		fmt.Println(fmt.Sprintf("Event: Label %s has been %s", labelEvent.GetLabel().GetName(), labelEvent.GetAction()))
		return nil
	}
	return nil
}
