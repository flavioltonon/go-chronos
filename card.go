package chronos

import (
	"github.com/flavioltonon/go-github/github"
)

type Card struct {
	*github.ProjectCard

	PriorityLevel int
	Deadline      int // in hours
}

type CardsByPriority []Card

func (c CardsByPriority) Len() int           { return len(c) }
func (c CardsByPriority) Less(i, j int) bool { return c[i].PriorityLevel < c[j].PriorityLevel }
func (c CardsByPriority) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type CardsByDeadline []Card

func (c CardsByDeadline) Len() int           { return len(c) }
func (c CardsByDeadline) Less(i, j int) bool { return c[i].Deadline < c[j].Deadline }
func (c CardsByDeadline) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
