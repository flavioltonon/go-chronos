package chronos

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type ChronosSetIssueDeadlineRequest struct {
	IssueNumber int
	Label       string
	Created     time.Time

	elapsedTime  float64
	nonWorkHours float64
	timer        string
	overdue      bool
	timerLabel   string
}

type ChronosSetIssueDeadlineResponse struct {
}

func (r *ChronosSetIssueDeadlineRequest) calculateElapsedTime() error {
	var (
		chronos      Chronos
		nonWorkHours float64
		weekendHours float64
		holidayHours float64
	)

	loc, _ := time.LoadLocation("America/Sao_Paulo")
	now := time.Now().In(loc)
	created := r.Created.In(loc)
	elapsedTime := int(math.Round(now.Sub(created).Hours()))

	err := chronos.GetHolidays(now.Year())
	if err != nil {
		return err
	}

	for t := 0; t < elapsedTime; t++ {
		// Check if it is Sunday
		if created.Add(time.Duration(t)*time.Hour).Weekday() == 0 {
			weekendHours++
			continue
		}

		// Check if it is Saturday
		if created.Add(time.Duration(t)*time.Hour).Weekday() == 6 {
			weekendHours++
			continue
		}

		// Check for holidays
		_, exists := chronos.holidays[created.Add(time.Duration(t)*time.Hour).Format("2006-01-02")]
		if exists {
			holidayHours++
			continue
		}

		// Check if it is a work hour
		if created.Add(time.Duration(t)*time.Hour).Hour() < 9 {
			nonWorkHours++
			continue
		}
		if created.Add(time.Duration(t)*time.Hour).Hour() >= 18 {
			nonWorkHours++
			continue
		}
	}

	r.elapsedTime = now.Sub(created).Hours() - weekendHours - holidayHours
	r.nonWorkHours = nonWorkHours

	return nil
}

func (r *ChronosSetIssueDeadlineRequest) defineTimer() error {
	var (
		deadline string
		usage    float64
	)

	timeTable := make(map[string]float64)
	timeTable["horas"] = r.elapsedTime - usage*r.nonWorkHours
	timeTable["dias"] = r.elapsedTime / 9.0

	switch r.Label {
	case "Prioridade: Baixa":
		deadline = DEADLINE_PRIORIDADE_BAIXA
		usage = 1
	case "Prioridade: Média":
		deadline = DEADLINE_PRIORIDADE_MEDIA
		usage = 1
	case "Prioridade: Alta":
		deadline = DEADLINE_PRIORIDADE_ALTA
		usage = 1
	case "Prioridade: Muito Alta":
		deadline = DEADLINE_PRIORIDADE_MUITO_ALTA
		usage = 0
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
		chronos  Chronos
		newLabel string
	)

	newLabel = fmt.Sprintf("Prazo: %s", r.timer)

	if r.overdue {
		newLabel = "Overdue"
	}

	err := chronos.GetLabel(newLabel)
	if err != nil {
		err := chronos.CreateLabel(newLabel)
		if err != nil {
			return err
		}
	}

	r.timerLabel = newLabel

	return nil
}

func (r ChronosSetIssueDeadlineRequest) removeOldTimer() error {
	var (
		err     error
		chronos Chronos
		labels  []string
	)

	err = chronos.GetLabelsFromIssue(r.IssueNumber)
	if err != nil {
		return err
	}

	for _, label := range chronos.labels {
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

	err = chronos.DeleteLabelsFromIssue(r.IssueNumber, labels)
	if err != nil {
		return err
	}

	return nil
}

func (r ChronosSetIssueDeadlineRequest) addNewTimer() error {
	var chronos Chronos

	err := chronos.AddLabelsToIssue(r.IssueNumber, []string{r.timerLabel})
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
