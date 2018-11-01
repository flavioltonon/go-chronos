package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-resty/resty"
	"github.com/google/go-github/github"
)

const (
	GITHUB_API_URL = "https://api.github.com"
	OWNER          = "flavioltonon"
	REPO           = "go-chronos"
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
	var repo Repo

	event := h.event.(*github.IssuesEvent)
	switch event.GetAction() {
	case "opened":
		fmt.Println("Event: Issue", event.GetAction())
	case "labeled":
		fmt.Println(fmt.Sprintf("Event: Issue %d has been %s", event.GetIssue().GetNumber(), event.GetAction()))
		fmt.Println("Label:", event.GetLabel().GetName())
	case "unlabeled":
		fmt.Println(fmt.Sprintf("Event: Issue %d has been %s", event.GetIssue().GetNumber(), event.GetAction()))
		fmt.Println("Label:", event.GetLabel().GetName())
		query := ""
		repo.GetIssues(query)
	default:
		fmt.Println("Event: Issue", event.GetAction())
		return nil
	}
	return nil
}

func (h *Chronos) HandleProjectCardEvent() error {
	event := h.event.(*github.ProjectCardEvent)
	switch event.GetAction() {
	case "moved":
		fmt.Println(fmt.Sprintf("Event: Project card %d has been %s to column %d", event.GetProjectCard().GetID(), event.GetAction(), event.GetProjectCard().GetColumnID()))
	default:
		fmt.Println(fmt.Sprintf("Event: Project card", event.GetAction()))
		return nil
	}
	return nil
}

func (h *Chronos) HandlePingEvent() error {
	event := h.event.(*github.PingEvent)
	fmt.Println("Event: Ping received.\nZen message:", event.GetZen())
	return nil
}

type Issues struct {
	ID     int
	Number int
	Labels []Label
}

type Label struct {
	ID   int
	Name string
}

type Repo struct {
	name   *string
	issues *Issues
}

func NewRepo(name string) Repo {
	return Repo{
		name: &name,
	}
}

func (h Repo) Name() string {
	return *h.name
}

func (h Repo) Issues() Issues {
	return *h.issues
}

func (h Repo) GetIssues(query string) error {
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/issues", GITHUB_API_URL, OWNER, REPO))
	fullURL, _ := u.Parse(query)
	resp, err := resty.R().
		SetBasicAuth(os.Getenv("CHRONOS_GITHUB_LOGIN"), os.Getenv("CHRONOS_GITHUB_PASSWORD")).
		Get(fullURL.String())
	if err != nil {
		log.Println(err)
		return err
	}

	var issues interface{}
	err = json.Unmarshal(resp.Body(), &issues)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("%+v", issues))

	return nil
}
