package chronos

type ColumnSprintDone struct{}

func init() {
	RegisterColumn(&ColumnSprintDone{})
}

func (p ColumnSprintDone) ID() int64 {
	return 4284730
}

func (p ColumnSprintDone) Name() string {
	return "Sprint Done"
}

func (p ColumnSprintDone) Project() int64 {
	return 1302676
}

func (p ColumnSprintDone) StandardIssueState() string {
	return "closed"
}
