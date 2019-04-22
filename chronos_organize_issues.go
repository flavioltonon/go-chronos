package chronos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/flavioltonon/go-github/github"
)

type ChronosOrganizeIssuesRequest struct {
	Option string

	client *github.Client
	issues []*github.Issue
}

type ChronosOrganizeIssuesResponse struct{}

var organizationStandardOptions = []string{
	"priority",
	"deadline",
}

func (r *ChronosOrganizeIssuesRequest) preCondition() error {
	for _, option := range organizationStandardOptions {
		if r.Option == option {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("organization option not available:", r.Option))
}

func (r *ChronosOrganizeIssuesRequest) organize(cards []*github.ProjectCard) error {
	var unorganizedCards = make([]Card, 0)

	for _, card := range cards {
		var (
			newCard Card
			ddl     int
		)

		newCard.ProjectCard = card

		issueNumber, _ := strconv.Atoi(strings.Split(card.GetContentURL(), "/issues/")[1])

		labels, _, err := r.client.Issues.ListLabelsByIssue(
			context.Background(),
			os.Getenv("GITHUB_REPOSITORY_OWNER"),
			os.Getenv("GITHUB_REPOSITORY_NAME"),
			issueNumber,
			&github.ListOptions{},
		)
		if err != nil {
			return ErrUnableToGetIssueLabels
		}

		for _, label := range labels {
			if regexp.MustCompile(DEADLINE_LABEL_SIGNATURE).MatchString(label.GetName()) {
				if label.GetName() != DEADLINE_LABEL_OVERDUE {
					deadline := strings.TrimPrefix(strings.Split(label.GetName(), ":")[1], " ")
					split := strings.Split(deadline, " ")
					ddl, _ = strconv.Atoi(split[0])
					unit := split[1]

					if unit == DEADLINE_TYPE_DAYS {
						ddl *= 24
					}
				} else {
					ddl = 0
				}

				newCard.Deadline = ddl
			}

			if regexp.MustCompile(PRIORITY_LABEL_SIGNATURE).MatchString(label.GetName()) {
				if _, exists := Priorities()[label.GetID()]; !exists {
					log.Println("priority label not registered:", label.GetName())
					continue
				}

				newCard.PriorityLevel = Priorities()[label.GetID()].Level()
			}
		}

		unorganizedCards = append(unorganizedCards, newCard)
	}

	switch r.Option {
	case "priority":
		var cardsByPriority = CardsByPriority(unorganizedCards)

		sort.Sort(cardsByPriority)

		for _, card := range cardsByPriority {
			r.client.Projects.MoveProjectCard(context.Background(), card.GetID(), &github.ProjectCardMoveOptions{
				Position: "bottom",
				ColumnID: card.GetColumnID(),
			})
		}
	case "deadline":
		var cardsByDeadline = CardsByDeadline(unorganizedCards)

		sort.Sort(cardsByDeadline)

		for _, card := range cardsByDeadline {
			r.client.Projects.MoveProjectCard(context.Background(), card.GetID(), &github.ProjectCardMoveOptions{
				Position: "bottom",
				ColumnID: card.GetColumnID(),
			})
		}
	}

	return nil
}

func (h *Chronos) OrganizeIssues() error {
	var req = h.request.(ChronosOrganizeIssuesRequest)

	req.client = h.client

	err := req.preCondition()
	if err != nil {
		return err
	}

	var columns = Columns()
	for _, column := range columns {
		cards, _, err := h.client.Projects.ListProjectCards(
			context.Background(),
			column.ID(),
			&github.ProjectCardListOptions{},
		)
		if err != nil {
			log.Println(fmt.Sprintf("unable to get project column %v cards", column.ID()))
			continue
		}

		err = req.organize(cards)
		if err != nil {
			log.Fatal(fmt.Sprintf("unable to organize column %v cards", column.ID()))
		}
	}

	return nil
}
