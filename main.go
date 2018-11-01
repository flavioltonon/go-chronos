package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var chronos Chronos

	payload, err := github.ValidatePayload(r, []byte("=Xwn.cj7"))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	chronos.event, err = github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch chronos.event.(type) {
	case *github.IssuesEvent:
		err := chronos.HandleIssueEvent()
		if err != nil {
			return
		}
	case *github.ProjectCardEvent:
		err := chronos.HandleProjectCardEvent()
		if err != nil {
			return
		}
	case *github.PingEvent:
		err := chronos.HandlePingEvent()
		if err != nil {
			return
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func main() {
	log.Println("server started")
	http.HandleFunc("/chronos", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Chronos struct {
	event interface{}
}

func (h *Chronos) HandleIssueEvent() error {
	event := h.event.(*github.IssuesEvent)
	fmt.Println("got webhook payload: ", *event.Action)
	return nil
}

func (h *Chronos) HandleProjectCardEvent() error {
	event := h.event.(*github.ProjectCardEvent)
	switch *event.Action {
	case "moved":
		fmt.Println("got webhook payload: ", *event.Action)
		fmt.Println("who: ", *event.ProjectCard.URL)
		fmt.Println("from: ", *event.Changes.Note.From)
		fmt.Println("to: ", *event.ProjectCard.ColumnID)
	default:
		return nil
	}
	return nil
}

func (h *Chronos) HandlePingEvent() error {
	event := h.event.(*github.PingEvent)
	fmt.Println("got webhook payload: ", *event.Zen)
	return nil
}
