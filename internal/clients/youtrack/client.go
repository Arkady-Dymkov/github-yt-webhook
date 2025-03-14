package youtrack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github-yt-webhook/internal/models/github"
	"net/http"
	"strings"

	"github-yt-webhook/internal/models"
	"github-yt-webhook/internal/utils"
)

// Client interface for YouTrack operations
type Client interface {
	ExecuteCommands(issueExtractable github.IssueExtractable, commands []models.YouTrackCommand) error
}

// HTTPClient is the implementation of the Client interface
type HTTPClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// IssueReference represents a reference to a YouTrack issue
type IssueReference struct {
	IDReadable string `json:"idReadable"`
}

// CommandRequest represents the request structure for YouTrack commands
type CommandRequest struct {
	Query   string           `json:"query"`
	Issues  []IssueReference `json:"issues"`
	Comment string           `json:"comment,omitempty"`
}

// NewClient creates a new YouTrack client
func NewClient(baseURL, token string) (Client, error) {
	if baseURL == "" {
		return nil, errors.New("YouTrack URL is required")
	}
	if token == "" {
		return nil, errors.New("YouTrack token is required")
	}

	// Ensure baseURL doesn't end with a slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	httpClient := &HTTPClient{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{},
	}

	return httpClient, nil
}

func (c *HTTPClient) ExecuteCommands(issueExtractable github.IssueExtractable, commands []models.YouTrackCommand) error {
	for _, command := range commands {
		var err = c.ExecuteCommand(issueExtractable, command)
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecuteCommand executes a YouTrack command on an issue
func (c *HTTPClient) ExecuteCommand(issueExtractable github.IssueExtractable, command models.YouTrackCommand) error {

	ticket := ExtractTicket(issueExtractable.GetIssueNumberPlace())
	if ticket == "" {
		return errors.New("ticket ID is required")
	}

	if command.Command == "" {
		return errors.New("command is not specified")
	}

	// Prepare the command to update the issue status
	commandURL := fmt.Sprintf("%s/api/commands", c.baseURL)

	// Create the command request with the correct structure
	commandData := CommandRequest{
		Query: command.Command,
		Issues: []IssueReference{
			{IDReadable: ticket},
		},
		Comment: issueExtractable.FillComment(command.Comment),
	}

	// Convert command data to JSON
	jsonData, err := json.Marshal(commandData)
	if err != nil {
		return fmt.Errorf("failed to marshal command data: %w", err)
	}

	// Debug log - print the request body
	utils.Debugf("YouTrack API Request Body: %s", string(jsonData))

	// Create request
	req, err := http.NewRequest("POST", commandURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Debug log - print the complete request details
	utils.Debugf("YouTrack API Request: %s %s", req.Method, req.URL.String())
	utils.Debugf("YouTrack API Request Headers: Content-Type: %s, Authorization: Bearer %s***",
		req.Header.Get("Content-Type"),
		c.token[:4]) // Only show first 4 chars of token for security

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Debug log - print the response status
	utils.Infof("YouTrack API Response Status: %d %s", resp.StatusCode, resp.Status)

	// Read and log response body for debugging
	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		utils.Infof("Error reading response body: %v", err)
	} else {
		utils.Debugf("YouTrack API Response Body: %s", responseBody.String())
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("YouTrack API returned non-OK status: %d, body: %s",
			resp.StatusCode, responseBody.String())
	}

	return nil
}
