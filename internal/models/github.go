package models

import (
	"regexp"
	"strings"
)

type ActionableEvent interface {
	GetAction() string
}

type IssueExtractable interface {
	GitIssueNumber() string
	FillComment(str string) string
}

type GitHubEvent interface {
	ActionableEvent
	IssueExtractable
}

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

type PushEvent struct {
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Author  Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *Commit) GetAction() string {
	return "any"
}

func (c *Commit) GitIssueNumber() string {
	return extractTicket(c.Message)
}

func (c *Commit) FillComment(str string) string {
	replacements := map[string]string{
		"{{commit_id}}":           c.Id,
		"{{commit_message}}":      c.Message,
		"{{commit_author.name}}":  c.Author.Name,
		"{{commit_author.email}}": c.Author.Email,
	}

	return replaceMultiple(str, replacements)
}

func (pr *PullRequestEvent) GetAction() string {
	return pr.Action
}

// GitIssueNumber getTitle returns the title of the pull request
func (pr *PullRequestEvent) GitIssueNumber() string {
	return extractTicket(pr.PullRequest.Title)
}

// FillComment fillComment replaces the placeholders in the comment with the actual values
func (pr *PullRequestEvent) FillComment(str string) string {
	replacements := map[string]string{
		"{{title}}":    pr.PullRequest.Title,
		"{{html_url}}": pr.PullRequest.HtmlUrl,
	}
	return replaceMultiple(str, replacements)
}

// replaceMultiple replaces multiple strings in a given string
func replaceMultiple(str string, replacements map[string]string) string {
	for oldValue, newValue := range replacements {
		str = strings.ReplaceAll(str, oldValue, newValue)
	}
	return str
}

// extractTicket extracts a ticket ID from a pull request title
func extractTicket(title string) string {
	// Assuming ticket format like ABC-123
	re := regexp.MustCompile(`[A-Z]+-\d+`)
	ticket := re.FindString(title)
	return strings.TrimSpace(ticket)
}
