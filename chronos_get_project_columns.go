package chronos

import (
	"context"

	"github.com/flavioltonon/go-github/github"
)

type ChronosGetProjectColumnsRequest struct {
	ProjectID int64

	client *github.Client

	columns map[int64]string

	issue *github.Issue

	issueState       string
	issueStatusLabel string
}

type ChronosGetProjectColumnsResponse struct {
	Columns map[int64]string `json:"columns"`
}

func (h *ChronosGetProjectColumnsRequest) getProjectColumns() error {
	projectColumns, _, err := h.client.Projects.ListProjectColumns(context.Background(), h.ProjectID, nil)
	if err != nil {
		return ErrUnableToGetProjectColumns
	}

	columns := make(map[int64]string)
	for _, column := range projectColumns {
		columns[column.GetID()] = column.GetName()
	}

	h.columns = columns

	return nil
}

func (h Chronos) GetProjectColumns() (resp ChronosGetProjectColumnsResponse, err error) {
	var req = h.request.(ChronosGetProjectColumnsRequest)

	req.client = h.client

	err = req.getProjectColumns()
	resp.Columns = req.columns

	return resp, err
}
