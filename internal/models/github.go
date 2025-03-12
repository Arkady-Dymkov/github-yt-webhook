package models

// PullRequestEvent represents a GitHub pull request event
type PullRequestEvent struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Title string `json:"title"`
	// Add other fields as needed
}