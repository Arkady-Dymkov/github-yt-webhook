package youtrack

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		token     string
		wantError bool
	}{
		{
			name:      "Valid client",
			baseURL:   "https://youtrack.example.com",
			token:     "test-token",
			wantError: false,
		},
		{
			name:      "Missing URL",
			baseURL:   "",
			token:     "test-token",
			wantError: true,
		},
		{
			name:      "Missing token",
			baseURL:   "https://youtrack.example.com",
			token:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.token)
			if (err != nil) != tt.wantError {
				t.Errorf("NewClient() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && client == nil {
				t.Errorf("NewClient() returned nil client when no error was expected")
			}
		})
	}
}

func TestUpdateIssueStatus(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check request path
		if r.URL.Path != "/api/commands" {
			t.Errorf("Expected path /api/commands, got %s", r.URL.Path)
		}

		// Check authorization header
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test update issue status
	err = client.UpdateIssueStatus("TEST-123", "Duplicate")
	if err != nil {
		t.Errorf("UpdateIssueStatus() error = %v", err)
	}

	// Test with empty ticket
	err = client.UpdateIssueStatus("", "")
	if err == nil {
		t.Errorf("UpdateIssueStatus() with empty ticket should return error")
	}
}
