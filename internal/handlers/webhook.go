package handlers

import (
	"log"
	"net/http"
	"regexp"
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
	// Check if the event is a pull_request event
	eventType := c.GetHeader("X-GitHub-Event")

	actionsMapping, exists := h.config.EventMapping[eventType]
	if !exists {
		log.Printf("Ignored event type: %s", eventType)
		c.String(http.StatusOK, "Event ignored")
		return
	}

	event, err := bindGithubEvent(eventType, c)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if we have a mapping for this action
	mapping, exists := actionsMapping.GitHubActions[event.GetAction()]
	if !exists {
		log.Printf("Ignored pull_request action: %s (no mapping configured)", event.GetAction())
		c.String(http.StatusOK, "Action ignored (no mapping configured)")
		return
	}

	// Send request to YouTrack to update the issue status using the configured command
	err = h.ytClient.ExecuteCommands(event, mapping.YouTrackCommand)
	if err != nil {
		log.Printf("Failed to complete hook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete hook: " + err.Error()})
		return
	}

	log.Printf("YouTrack issue %s updated with command '%s'", event.GitIssueNumber(), mapping.YouTrackCommand)

	c.String(http.StatusOK, "Issue "+event.GitIssueNumber()+" updated")
}

// extractTicket extracts a ticket ID from a pull request title
func extractTicket(title string) string {
	// Assuming ticket format like ABC-123
	re := regexp.MustCompile(`[A-Z]+-\d+`)
	ticket := re.FindString(title)
	return strings.TrimSpace(ticket)
}

func bindGithubEvent(eventType string, context *gin.Context) (models.GitHubEvent, error) {
	var event models.GitHubEvent
	switch eventType {
	case "pull_request":
		var prEvent models.PullRequestEvent
		if err := context.ShouldBindJSON(&prEvent); err != nil {
			return nil, err
		}
		event = &prEvent
	default:
		return nil, nil
	}
	return event, nil
}
