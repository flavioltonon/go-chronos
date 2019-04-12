package column

type ColumnSprintBacklog struct{}

func init() {
	RegisterColumn(&ColumnSprintBacklog{})
}

func (p ColumnSprintBacklog) ID() int64 {
	return 3719176
}

func (p ColumnSprintBacklog) Name() string {
	return "Sprint backlog"
}

func (p ColumnSprintBacklog) StandardIssueState() string {
	return "open"
}
