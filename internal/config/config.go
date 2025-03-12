package config

import (
	"encoding/json"
	"fmt"
	"log
	"os"
	"strings"

	"github-yt-webhook/internal/models"
)

// Config holds all configuration for the application
type Config struct {
	Port          string
	YouTrackURL   string
	YouTrackToken string
	EventMapping  map[string]models.EventMapping
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
		Port:          port,
		YouTrackURL:   ytURL,
		YouTrackToken: ytToken,
		EventMapping:  make(map[string]models.EventMapping),
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

	var mappingConfig struct {
		Mappings []struct {
			EventName     string                `json:"eventName"`
			GitHubActions []models.GitHubAction `json:"githubActions"`
		} `json:"mappings"`
	}

	if err := json.Unmarshal(data, &mappingConfig); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	c.EventMapping = make(map[string]models.EventMapping)

	for _, mapping := range mappingConfig.Mappings {
		githubActionsMap := make(map[string]models.GitHubAction)
		for _, action := range mapping.GitHubActions {
			githubActionsMap[action.GitHubAction] = action
		}

		c.EventMapping[mapping.EventName] = models.EventMapping{
			EventName:     mapping.EventName,
			GitHubActions: githubActionsMap,
		}
	}

	// Log the loaded structure
	jsonData, err := json.MarshalIndent(c.EventMapping, "", "  ")
	if err != nil {
		log.Printf("Error marshaling EventMapping for logging: %v", err)
	} else {
		log.Printf("Loaded EventMapping:\n%s", jsonData)
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
		"Config{Port: %s, YouTrackURL: %s, YouTrackToken: %s, EventMapping: %d}",
		c.Port,
		c.YouTrackURL,
		maskedToken,
		len(c.EventMapping),
	)
}

func removeQuotes(input string) string {
	if input == "" {
		return input
	}
	return strings.Trim(input, "\"")
}
