package github

type ActionableEvent interface {
	GetAction() string
}

type IssueExtractable interface {
	GetIssueNumberPlace() string
	FillComment(str string) string
}

type GitHubEvent interface {
	ActionableEvent
	IssueExtractable
}
