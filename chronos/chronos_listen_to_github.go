package chronos

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

func (chronos *Chronos) ListenToGitHub(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		log.Println(fmt.Errorf("error validating request body: err=%s", err))
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Println(fmt.Errorf("could not parse webhook: err=%s", err))
		return
	}

	auth := github.BasicAuthTransport{
		Username: os.Getenv("CHRONOS_GITHUB_LOGIN"),
		Password: os.Getenv("CHRONOS_GITHUB_PASSWORD"),
	}

	chronos.SetClient(github.NewClient(auth.Client()))

	switch event.(type) {
	case *github.IssuesEvent:
		log.Println(chronos.HandleIssuesEvent(event))
		return
	case *github.LabelEvent:
		log.Println(chronos.HandleLabelEvent(event))
		return
	case *github.ProjectCardEvent:
		log.Println(chronos.HandleProjectCardEvent(event))
		return
	case *github.PingEvent:
		log.Println(chronos.HandlePingEvent(event))
		return
	default:
		log.Println(fmt.Errorf("unknown event type %s", github.WebHookType(r)))
		return
	}
}
