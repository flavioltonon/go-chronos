package chronos

type ColumnDone struct{}

func init() {
	RegisterColumn(&ColumnDone{})
}

func (p ColumnDone) ID() int64 {
	return 2244760
}

func (p ColumnDone) Name() string {
	return "Done"
}

func (p ColumnDone) Project() int64 {
	return 1302676
}

func (p ColumnDone) StandardIssueState() string {
	return "closed"
}
