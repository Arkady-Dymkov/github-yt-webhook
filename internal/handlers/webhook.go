package handlers

import (
	"log"
	"net/http"
	"strings"

	"github-yt-webhook/internal/config"
	"github-yt-webhook/internal/models"
	"github-yt-webhook/internal/youtrack"
	"github.com/gin-gonic/gin"
)

// WebhookHandler handles GitHub webhook requests
type WebhookHandler struct {
	ytClient youtrack.Client
	config   *config.Config
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(ytClient youtrack.Client, config *config.Config) *WebhookHandler {
	return &WebhookHandler{
		ytClient: ytClient,
		config:   config,
	}
}

// HandleGitHubWebhook handles GitHub webhook events
func (h *WebhookHandler) HandleGitHubWebhook(c *gin.Context) {
	// Check if the events is a pull_request events
	eventType := c.GetHeader("X-GitHub-Event")

	actionsMapping, exists := h.config.EventMapping[eventType]
	if !exists {
		log.Printf("Ignored events type: %s", eventType)
		c.String(http.StatusOK, "Event ignored")
		return
	}

	events, err := bindGithubEvent(eventType, c)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	h.processEvents(events, actionsMapping.GitHubActions)

	c.String(http.StatusOK, "Issues "+strings.Join(collectsIssuesNumbers(events), ", ")+" updated")
}

func (h *WebhookHandler) processEvents(events []models.GitHubEvent, actionsMapping map[string]models.GitHubAction) {
	for _, event := range events {
		mapping, exists := actionsMapping[event.GetAction()]
		if !exists {
			log.Printf("Ignored pull_request action: %s (no mapping configured)", event.GetAction())
			continue
		}

		// Send request to YouTrack to update the issue status using the configured command
		err := h.ytClient.ExecuteCommands(event, mapping.YouTrackCommand)
		if err != nil {
			log.Printf("Failed to complete hook: %v", err)
			continue
		}

		log.Printf("YouTrack issue %s updated with command '%s'", event.GitIssueNumber(), mapping.YouTrackCommand)
	}
}

func collectsIssuesNumbers(events []models.GitHubEvent) []string {
	var issues []string
	for _, event := range events {
		issues = append(issues, event.GitIssueNumber())
	}
	return issues
}

func bindGithubEvent(eventType string, context *gin.Context) ([]models.GitHubEvent, error) {
	var events []models.GitHubEvent
	switch eventType {
	case "pull_request":
		var prEvent models.PullRequestEvent
		if err := context.ShouldBindJSON(&prEvent); err != nil {
			return nil, err
		}
		events = append(events, &prEvent)
	case "push":
		var pushEvent models.PushEvent
		if err := context.ShouldBindJSON(&pushEvent); err != nil {
			return nil, err
		}
		for _, commit := range pushEvent.Commits {
			events = append(events, &commit)
		}

	default:
		return nil, nil
	}
	return events, nil
}
