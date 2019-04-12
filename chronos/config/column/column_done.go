package column

type ColumnDone struct{}

func init() {
	RegisterColumn(&ColumnDone{})
}

func (p ColumnDone) ID() int64 {
	return 3731528
}

func (p ColumnDone) Name() string {
	return "Done"
}

func (p ColumnDone) StandardIssueState() string {
	return "closed"
}
