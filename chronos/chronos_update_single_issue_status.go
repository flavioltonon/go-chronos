package chronos

type ChronosUpdateSingleIssueStatusRequest struct {
	IssueNumber int
	ColumnToID  int64

	columnsMap map[int64]string
}

type ChronosUpdateSingleIssueStatusResponse struct {
}

func (h *ChronosUpdateSingleIssueStatusRequest) mapProjectColumns() error {
	// var (
	// 	chronos Chronos
	// 	err     error
	// )

	// err = chronos.GetProjectColumns()
	// if err != nil {
	// 	return err
	// }

	// columnsMap := make(map[int64]string)
	// for _, column := range chronos.projectColumns {
	// 	columnsMap[column.ID] = column.Name
	// }

	// h.columnsMap = columnsMap

	return nil
}

func (h Chronos) UpdateSingleIssueStatus(req *ChronosUpdateSingleIssueStatusRequest) error {
	var err error

	err = req.mapProjectColumns()
	if err != nil {
		return err
	}

	return nil
}
