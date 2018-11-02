package github

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type ChronosSetIssueDeadlineRequest struct {
	IssueNumber string
	Label       string
	Created     time.Time

	elapsedTime float64
	timer       string
	overdue     bool
	timerLabel  string
}

type ChronosSetIssueDeadlineResponse struct {
}

func (r *ChronosSetIssueDeadlineRequest) calculateElapsedTime() error {
	r.elapsedTime = time.Now().Sub(r.Created).Hours()
	return nil
}

func (r *ChronosSetIssueDeadlineRequest) defineTimer() error {
	var daysElapsed = math.Round(r.elapsedTime / 24.0)
	var hoursElapsed = math.Round(r.elapsedTime)
	switch r.Label {
	case "Prioridade: Baixa":
		if daysElapsed > 60 {
			r.overdue = true
		}
		r.timer = strconv.FormatFloat(60-daysElapsed, 'f', -1, 64) + " dias"
	case "Prioridade: Média":
		if daysElapsed > 15 {
			r.overdue = true
		}
		r.timer = strconv.FormatFloat(15-daysElapsed, 'f', -1, 64) + " dias"
	case "Prioridade: Alta":
		if daysElapsed > 3 {
			r.overdue = true
		}
		r.timer = strconv.FormatFloat(3-daysElapsed, 'f', -1, 64) + " dias"
	case "Prioridade: Muito Alta":
		if hoursElapsed > 24 {
			r.overdue = true
		}
		r.timer = strconv.FormatFloat(24-hoursElapsed, 'f', -1, 64) + " horas"
	default:
		return ErrUnableToDefineTimer
	}
	return nil
}

func (r *ChronosSetIssueDeadlineRequest) createLabel() error {
	var (
		repo     Repo
		newLabel string
	)

	if r.overdue {
		newLabel = "Overdue"
	} else {
		newLabel = fmt.Sprintf("Prazo: %s", r.timer)
	}
	err := repo.GetLabel(newLabel)
	if err != nil {
		err := repo.CreateLabel(newLabel)
		if err != nil {
			return err
		}
	}

	r.timerLabel = newLabel

	return nil
}

func (r ChronosSetIssueDeadlineRequest) removeOldTimer() error {
	var (
		err    error
		repo   Repo
		labels []string
	)

	err = repo.GetLabelsFromIssue(r.IssueNumber)
	if err != nil {
		return err
	}

	for _, label := range repo.labels {
		if strings.Split(label.Name, ": ")[0] == "Prazo" {
			if strings.Split(label.Name, " ")[2] == "dias" || strings.Split(label.Name, " ")[2] == "horas" {
				labels = append(labels, label.Name)
			}
		}
		if strings.Split(label.Name, ": ")[0] == "Prioridade" {
			if label.Name != r.Label {
				labels = append(labels, label.Name)
			}
		}
	}

	err = repo.DeleteLabelsFromIssue(r.IssueNumber, labels)
	if err != nil {
		return err
	}

	return nil
}

func (r ChronosSetIssueDeadlineRequest) addNewTimer() error {
	var repo Repo

	err := repo.AddLabelsToIssue(r.IssueNumber, []string{r.timerLabel})
	if err != nil {
		return err
	}

	return nil
}

func (h Chronos) SetIssueDeadline(req *ChronosSetIssueDeadlineRequest) error {
	var err error

	err = req.calculateElapsedTime()
	if err != nil {
		return err
	}

	err = req.defineTimer()
	if err != nil {
		return err
	}

	err = req.createLabel()
	if err != nil {
		return err
	}

	err = req.removeOldTimer()
	if err != nil {
		return err
	}

	err = req.addNewTimer()
	if err != nil {
		return err
	}

	return nil
}
