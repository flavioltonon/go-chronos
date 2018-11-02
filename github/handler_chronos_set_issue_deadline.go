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
	var nonWorkHours float64

	now := time.Now()
	_, offset := now.Zone()
	nowUTC := now.Add(time.Duration(offset) * time.Second).UTC()
	elapsedTime := int(math.Round(nowUTC.Sub(r.Created).Hours()))

	for t := 0; t < elapsedTime; t++ {
		if r.Created.Add(time.Duration(t)*time.Hour).Weekday() == 0 { // Sunday
			nonWorkHours++
			continue
		}

		if r.Created.Add(time.Duration(t)*time.Hour).Weekday() == 6 { // Saturday
			nonWorkHours++
			continue
		}

		if r.Created.Add(time.Duration(t)*time.Hour).Hour() < 12 {
			nonWorkHours++
			continue
		}
		if r.Created.Add(time.Duration(t)*time.Hour).Hour() >= 21 {
			nonWorkHours++
			continue
		}
	}

	r.elapsedTime = nowUTC.Sub(r.Created).Hours() - nonWorkHours

	return nil
}

func (r *ChronosSetIssueDeadlineRequest) defineTimer() error {
	var deadline string

	timeTable := make(map[string]float64)
	timeTable["horas"] = r.elapsedTime
	timeTable["dias"] = r.elapsedTime / 24.0

	switch r.Label {
	case "Prioridade: Baixa":
		deadline = DEADLINE_PRIORIDADE_BAIXA
	case "Prioridade: Média":
		deadline = DEADLINE_PRIORIDADE_MEDIA
	case "Prioridade: Alta":
		deadline = DEADLINE_PRIORIDADE_ALTA
	case "Prioridade: Muito Alta":
		deadline = DEADLINE_PRIORIDADE_MUITO_ALTA
	default:
		return ErrUnableToDefineTimer
	}

	deadlineTime, _ := strconv.ParseFloat(strings.Split(deadline, " ")[0], 64)
	deadlineType := strings.Split(deadline, " ")[1]
	if timeTable[deadlineType] > deadlineTime {
		r.overdue = true
	}
	r.timer = strconv.FormatFloat(deadlineTime-math.Round(timeTable[deadlineType]), 'f', -1, 64) + " " + deadlineType

	return nil
}

func (r *ChronosSetIssueDeadlineRequest) createLabel() error {
	var (
		repo     Repo
		newLabel string
	)

	newLabel = fmt.Sprintf("Prazo: %s", r.timer)

	if r.overdue {
		newLabel = "Overdue"
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
