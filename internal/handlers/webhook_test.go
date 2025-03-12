package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github-yt-webhook/internal/models"
)

// Mock YouTrack client for testing
type mockYouTrackClient struct {
	updateCalled bool
	ticket       string
	shouldError  bool
}

func (m *mockYouTrackClient) UpdateIssueStatus(ticket string) error {
	m.updateCalled = true
	m.ticket = ticket
	if m.shouldError {
		return errors.New("mock error")
	}
	return nil
}

func TestHandleGitHubWebhook(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		eventType      string
		body           models.PullRequestEvent
		mockShouldErr  bool
		expectedStatus int
		checkUpdate    bool
		expectedTicket string
	}{
		{
			name:           "Valid pull request opened",
			eventType:      "pull_request",
			body:           models.PullRequestEvent{Action: "opened", PullRequest: models.PullRequest{Title: "TEST-123: Test PR"}},
			mockShouldErr:  false,
			expectedStatus: http.StatusOK,
			checkUpdate:    true,
			expectedTicket: "TEST-123",
		},
		{
			name:           "Non-pull request event",
			eventType:      "push",
			body:           models.PullRequestEvent{},
			mockShouldErr:  false,
			expectedStatus: http.StatusOK,
			checkUpdate:    false,
		},
		{
			name:           "Pull request not opened",
			eventType:      "pull_request",
			body:           models.PullRequestEvent{Action: "closed"},
			mockShouldErr:  false,
			expectedStatus: http.StatusOK,
			checkUpdate:    false,
		},
		{
			name:           "No ticket in title",
			eventType:      "pull_request",
			body:           models.PullRequestEvent{Action: "opened", PullRequest: models.PullRequest{Title: "Test PR without ticket"}},
			mockShouldErr:  false,
			expectedStatus: http.StatusOK,
			checkUpdate:    false,
		},
		{
			name:           "YouTrack update error",
			eventType:      "pull_request",
			body:           models.PullRequestEvent{Action: "opened", PullRequest: models.PullRequest{Title: "TEST-123: Test PR"}},
			mockShouldErr:  true,
			expectedStatus: http.StatusInternalServerError,
			checkUpdate:    true,
			expectedTicket: "TEST-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &mockYouTrackClient{shouldError: tt.mockShouldErr}

			// Create handler with mock client
			handler := NewWebhookHandler(mockClient)

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request
			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest("POST", "/webhook", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("X-GitHub-Event", tt.eventType)
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			handler.HandleGitHubWebhook(c)

			// Check response status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check if update was called
			if tt.checkUpdate && !mockClient.updateCalled {
				t.Errorf("Expected YouTrack update to be called, but it wasn't")
			} else if !tt.checkUpdate && mockClient.updateCalled {
				t.Errorf("Expected YouTrack update not to be called, but it was")
			}

			// Check ticket
			if tt.checkUpdate && mockClient.ticket != tt.expectedTicket {
				t.Errorf("Expected ticket %s, got %s", tt.expectedTicket, mockClient.ticket)
			}
		})
	}
}

func TestExtractTicket(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "Standard format",
			title:    "ABC-123: Fix bug",
			expected: "ABC-123",
		},
		{
			name:     "No ticket",
			title:    "Fix bug",
			expected: "",
		},
		{
			name:     "Ticket at end",
			title:    "Fix bug ABC-123",
			expected: "ABC-123",
		},
		{
			name:     "Multiple tickets",
			title:    "ABC-123 DEF-456: Fix bug",
			expected: "ABC-123",
		},
		{
			name:     "Invalid format",
			title:    "abc-123: Fix bug",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTicket(tt.title)
			if result != tt.expected {
				t.Errorf("extractTicket(%s) = %s, expected %s", tt.title, result, tt.expected)
			}
		})
	}
}
