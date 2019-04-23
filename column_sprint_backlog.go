package chronos

type ColumnSprintBacklog struct{}

func init() {
	RegisterColumn(&ColumnSprintBacklog{})
}

func (p ColumnSprintBacklog) ID() int64 {
	return 2272504
}

func (p ColumnSprintBacklog) Name() string {
	return "Sprint backlog"
}

func (p ColumnSprintBacklog) Project() int64 {
	return 1302676
}

func (p ColumnSprintBacklog) StandardIssueState() string {
	return "open"
}
