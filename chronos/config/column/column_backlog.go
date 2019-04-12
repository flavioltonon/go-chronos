package column

type ColumnBacklog struct{}

func init() {
	RegisterColumn(&ColumnBacklog{})
}

func (p ColumnBacklog) ID() int64 {
	return 3719175
}

func (p ColumnBacklog) Name() string {
	return "Backlog"
}

func (p ColumnBacklog) StandardIssueState() string {
	return "open"
}
