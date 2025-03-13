package github

import (
	"github-yt-webhook/internal/utils"
)

// PullRequestEvent represents a GitHub pull request event
type PullRequestEvent struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Title   string `json:"title"`
	HtmlUrl string `json:"html_url"`
	// Add other fields as needed
}

func (pr *PullRequestEvent) GetAction() string {
	return pr.Action
}

// GetIssueNumberPlace getTitle returns the title of the pull request
func (pr *PullRequestEvent) GetIssueNumberPlace() string {
	return pr.PullRequest.Title
}

// FillComment fillComment replaces the placeholders in the comment with the actual values
func (pr *PullRequestEvent) FillComment(str string) string {
	replacements := map[string]string{
		"{{title}}":    pr.PullRequest.Title,
		"{{html_url}}": pr.PullRequest.HtmlUrl,
	}
	return utils.ReplaceMultiple(str, replacements)
}
