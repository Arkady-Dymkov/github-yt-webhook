package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Port          string
	YouTrackURL   string
	YouTrackToken string
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
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
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
		"Config{Port: %s, YouTrackURL: %s, YouTrackToken: %s}",
		c.Port,
		c.YouTrackURL,
		maskedToken,
	)
}

func removeQuotes(input string) string {
	if input == "" {
		return input
	}
	return strings.Trim(input, "\"")
}
