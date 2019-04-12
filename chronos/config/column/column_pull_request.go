package column

type ColumnPullRequest struct{}

func init() {
	RegisterColumn(&ColumnPullRequest{})
}

func (p ColumnPullRequest) ID() int64 {
	return 3719177
}

func (p ColumnPullRequest) Name() string {
	return "Pull Request"
}

func (p ColumnPullRequest) StandardIssueState() string {
	return "open"
}
