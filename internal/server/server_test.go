package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github-yt-webhook/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		wantError bool
	}{
		{
			name: "Valid config",
			config: &config.Config{
				Port:          "8080",
				YouTrackURL:   "https://youtrack.example.com",
				YouTrackToken: "token123",
			},
			wantError: false,
		},
		{
			name:      "Nil config",
			config:    nil,
			wantError: true,
		},
		{
			name: "Invalid YouTrack URL",
			config: &config.Config{
				Port:          "8080",
				YouTrackURL:   "",
				YouTrackToken: "token123",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, err := New(tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("New() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && srv == nil {
				t.Errorf("New() returned nil server when no error was expected")
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Create a valid config
	cfg := &config.Config{
		Port:          "8080",
		YouTrackURL:   "https://youtrack.example.com",
		YouTrackToken: "token123",
	}

	// Create server
	srv, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create a test HTTP server
	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	// Make a request to the health endpoint
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}