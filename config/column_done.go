package config

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

func (p ColumnDone) Project() int64 {
	return 1908642
}

func (p ColumnDone) StandardIssueState() string {
	return "closed"
}
