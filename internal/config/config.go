package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github-yt-webhook/internal/models"
)

// Config holds all configuration for the application
type Config struct {
	Port           string
	YouTrackURL    string
	YouTrackToken  string
	ActionMappings map[string]models.ActionMapping
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ytURL := removeQuotes(os.Getenv("YOUTRACK_TEST_URL"))
	ytToken := removeQuotes(os.Getenv("YOUTRACK_TEST_TOKEN"))

	config := &Config{
		Port:           port,
		YouTrackURL:    ytURL,
		YouTrackToken:  ytToken,
		ActionMappings: make(map[string]models.ActionMapping),
	}

	// Load action mappings from config file
	if err := config.loadActionMappings(); err != nil {
		return nil, fmt.Errorf("failed to load action mappings: %w", err)
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// loadActionMappings loads action mappings from a JSON configuration file
func (c *Config) loadActionMappings() error {

	// Check if config file exists
	configPath := os.Getenv("ACTION_MAPPINGS_CONFIG")
	if configPath == "" {
		configPath = "action_mappings.json"
	}

	// If file doesn't exist, use default mapping
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file doesn't exists or cannot be read: %w", err)
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	var mappingConfig models.ActionMappingConfig
	if err := json.Unmarshal(data, &mappingConfig); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Convert to map for easier lookup
	for _, mapping := range mappingConfig.Mappings {
		c.ActionMappings[mapping.GitHubAction] = mapping
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	var missingVars []string

	if c.YouTrackURL == "" {
		missingVars = append(missingVars, "YOUTRACK_TEST_URL")
	} else {
		// Basic URL validation
		if !strings.HasPrefix(c.YouTrackURL, "http://") && !strings.HasPrefix(c.YouTrackURL, "https://") {
			return fmt.Errorf("YOUTRACK_TEST_URL must start with http:// or https://")
		}
	}

	if c.YouTrackToken == "" {
		missingVars = append(missingVars, "YOUTRACK_TEST_TOKEN")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	return nil
}

// String returns a string representation of the config (with sensitive data masked)
func (c *Config) String() string {
	// Mask the token for security
	maskedToken := "****"
	if c.YouTrackToken != "" {
		// Show first 4 characters if token is long enough
		if len(c.YouTrackToken) > 8 {
			maskedToken = c.YouTrackToken[:4] + "****"
		}
	}

	return fmt.Sprintf(
		"Config{Port: %s, YouTrackURL: %s, YouTrackToken: %s, ActionMappings: %d}",
		c.Port,
		c.YouTrackURL,
		maskedToken,
		len(c.ActionMappings),
	)
}

func removeQuotes(input string) string {
	if input == "" {
		return input
	}
	return strings.Trim(input, "\"")
}
