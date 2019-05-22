package chronos

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func (h *Chronos) SetupSettings() error {
	var (
		projects   = make(map[int64]Project, 0)
		columns    = make(map[int64]Column, 0)
		priorities = make(map[int64]Priority, 0)
	)

	// Open our jsonFile
	file, err := os.Open("settings.json")
	if err != nil {
		panic("failed to read settings from settings.json")
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var settings Settings

	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		panic("failed to unmarshal settings")
	}

	for _, setting := range settings {
		newBytes, _ := json.Marshal(setting)

		tmp := struct {
			Kind string          `json:"kind"`
			Data json.RawMessage `json:"data"`
		}{}

		err = json.Unmarshal(newBytes, &tmp)
		if err != nil {
			panic("failed to unmarshal setting")
		}

		switch tmp.Kind {
		case "project":
			var newProject Project

			err = json.Unmarshal(tmp.Data, &newProject)
			if err != nil {
				panic("failed to unmarshal project")
			}

			projects[newProject.ID] = newProject

			h.projects = projects

		case "column":
			var newColumn Column

			err = json.Unmarshal(tmp.Data, &newColumn)
			if err != nil {
				panic("failed to unmarshal column")
			}

			columns[newColumn.ID] = newColumn

			h.columns = columns

		case "priority":
			var newPriority Priority

			err = json.Unmarshal(tmp.Data, &newPriority)
			if err != nil {
				panic("failed to unmarshal priority")
			}

			priorities[newPriority.ID] = newPriority

			h.priorities = priorities
		}
	}

	return nil
}
