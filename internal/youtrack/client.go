package youtrack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Client interface for YouTrack operations
type Client interface {
	UpdateIssueStatus(ticket string) error
}

// HTTPClient is the implementation of the Client interface
type HTTPClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewClient creates a new YouTrack client
func NewClient(baseURL, token string) (*HTTPClient, error) {
	if baseURL == "" {
		return nil, errors.New("YouTrack URL is required")
	}
	if token == "" {
		return nil, errors.New("YouTrack token is required")
	}

	// Ensure baseURL doesn't end with a slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &HTTPClient{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// UpdateIssueStatus updates the status of a YouTrack issue to "In Review"
func (c *HTTPClient) UpdateIssueStatus(ticket string) error {
	if ticket == "" {
		return errors.New("ticket ID is required")
	}

	// Prepare the command to update the issue status
	commandURL := fmt.Sprintf("%s/api/commands", c.baseURL)
	commandData := map[string]string{
		"query":   ticket + " Duplicate",
		"comment": "Status updated by GitHub webhook",
	}

	// Convert command data to JSON
	jsonData, err := json.Marshal(commandData)
	if err != nil {
		return fmt.Errorf("failed to marshal command data: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", commandURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("YouTrack API returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
