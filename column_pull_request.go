package chronos

type ColumnPullRequest struct{}

func init() {
	RegisterColumn(&ColumnPullRequest{})
}

func (p ColumnPullRequest) ID() int64 {
	return 2244744
}

func (p ColumnPullRequest) Name() string {
	return "Pull Request"
}

func (p ColumnPullRequest) Project() int64 {
	return 1302676
}

func (p ColumnPullRequest) StandardIssueState() string {
	return "open"
}
