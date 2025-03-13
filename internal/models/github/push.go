package github

import "github-yt-webhook/internal/utils"

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

func (c *Commit) GetIssueNumberPlace() string {
	return c.Message
}

func (c *Commit) FillComment(str string) string {
	replacements := map[string]string{
		"{{commit_id}}":           c.Id,
		"{{commit_message}}":      c.Message,
		"{{commit_author.name}}":  c.Author.Name,
		"{{commit_author.email}}": c.Author.Email,
	}

	return utils.ReplaceMultiple(str, replacements)
}
