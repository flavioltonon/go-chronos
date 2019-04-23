package chronos

type ColumnBacklog struct{}

func init() {
	RegisterColumn(&ColumnBacklog{})
}

func (p ColumnBacklog) ID() int64 {
	return 3493247
}

func (p ColumnBacklog) Name() string {
	return "Backlog"
}

func (p ColumnBacklog) Project() int64 {
	return 1302676
}

func (p ColumnBacklog) StandardIssueState() string {
	return "open"
}
