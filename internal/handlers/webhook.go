package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github-yt-webhook/internal/config"
	"github-yt-webhook/internal/models"
	"github-yt-webhook/internal/youtrack"
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
	if eventType != "pull_request" {
		log.Printf("Ignored event type: %s", eventType)
		c.String(http.StatusOK, "Event ignored")
		return
	}

	var prEvent models.PullRequestEvent
	if err := c.ShouldBindJSON(&prEvent); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if we have a mapping for this action
	mapping, exists := h.config.ActionMappings[prEvent.Action]
	if !exists {
		log.Printf("Ignored pull_request action: %s (no mapping configured)", prEvent.Action)
		c.String(http.StatusOK, "Action ignored (no mapping configured)")
		return
	}

	ticket := extractTicket(prEvent.PullRequest.Title)
	if ticket == "" {
		log.Printf("No ticket found in pull request title: %s", prEvent.PullRequest.Title)
		c.String(http.StatusOK, "No ticket found")
		return
	}

	// Send request to YouTrack to update the issue status using the configured command
	err := h.ytClient.ExecuteCommand(ticket, mapping.YouTrackCommand)
	if err != nil {
		log.Printf("Failed to update YouTrack issue %s: %v", ticket, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update YouTrack issue"})
		return
	}

	log.Printf("YouTrack issue %s updated with command '%s'", ticket, mapping.YouTrackCommand)
	c.String(http.StatusOK, "Issue "+ticket+" updated")
}

// extractTicket extracts a ticket ID from a pull request title
func extractTicket(title string) string {
	// Assuming ticket format like ABC-123
	re := regexp.MustCompile(`[A-Z]+-\d+`)
	ticket := re.FindString(title)
	return strings.TrimSpace(ticket)
}
