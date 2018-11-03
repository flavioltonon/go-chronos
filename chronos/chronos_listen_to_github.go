package chronos

import (
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

func (h Chronos) ListenToGitHub(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch event.(type) {
	case *github.IssuesEvent:
		var issuesEventHandler IssuesEventHandler
		err := issuesEventHandler.HandleEvent(event)
		if err != nil {
			log.Println(err)
			return
		}
	case *github.LabelEvent:
		var labelEventHandler LabelEventHandler
		err := labelEventHandler.HandleEvent(event)
		if err != nil {
			log.Println(err)
			return
		}
	case *github.ProjectCardEvent:
		var projectCardEventHandler ProjectCardEventHandler
		err := projectCardEventHandler.HandleEvent(event)
		if err != nil {
			log.Println(err)
			return
		}
	case *github.PingEvent:
		var pingEventHandler PingEventHandler
		err := pingEventHandler.HandleEvent(event)
		if err != nil {
			log.Println(err)
			return
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
