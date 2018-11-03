package chronos

import (
	"strings"
	"sync"
)

func (h *Chronos) getIssues(query map[string]string) error {
	return h.GetIssues(query)
}

func (h *Chronos) updateIssuesDeadlineLabels() error {
	var wg sync.WaitGroup

	for _, issue := range h.issues {

		wg.Add(1)

		go func(issue Issue) {
			var (
				err    error
				labels []string
			)

			for _, label := range issue.Labels {
				if strings.Split(label.Name, ": ")[0] == "Prioridade" {
					labels = append(labels, label.Name)
				}
			}

			err = h.DeleteLabelsFromIssue(issue.Number, labels)
			if err != nil {
				wg.Done()
				return
			}

			err = h.AddLabelsToIssue(issue.Number, labels)
			if err != nil {
				wg.Done()
				return
			}

			wg.Done()

		}(issue)
	}

	wg.Wait()

	return nil
}

func (h Chronos) UpdateIssuesDeadlines() error {
	var err error

	query := map[string]string{
		"state": "open",
	}
	err = h.getIssues(query)
	if err != nil {
		return err
	}

	err = h.updateIssuesDeadlineLabels()
	if err != nil {
		return err
	}

	return nil
}
