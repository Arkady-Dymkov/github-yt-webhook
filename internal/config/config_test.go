package config

import (
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
	}{
		{
			name: "Valid config",
			config: Config{
				Port:          "8080",
				YouTrackURL:   "https://youtrack.example.com",
				YouTrackToken: "token123",
			},
			wantError: false,
		},
		{
			name: "Missing YouTrack URL",
			config: Config{
				Port:          "8080",
				YouTrackURL:   "",
				YouTrackToken: "token123",
			},
			wantError: true,
		},
		{
			name: "Missing YouTrack token",
			config: Config{
				Port:          "8080",
				YouTrackURL:   "https://youtrack.example.com",
				YouTrackToken: "",
			},
			wantError: true,
		},
		{
			name: "Invalid YouTrack URL format",
			config: Config{
				Port:          "8080",
				YouTrackURL:   "youtrack.example.com", // Missing http:// or https://
				YouTrackToken: "token123",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Save original environment variables
	origYTURL := os.Getenv("YOUTRACK_TEST_URL")
	origYTToken := os.Getenv("YOUTRACK_TEST_TOKEN")
	origPort := os.Getenv("PORT")

	// Restore environment variables after test
	defer func() {
		os.Setenv("YOUTRACK_TEST_URL", origYTURL)
		os.Setenv("YOUTRACK_TEST_TOKEN", origYTToken)
		os.Setenv("PORT", origPort)
	}()

	tests := []struct {
		name      string
		envVars   map[string]string
		wantError bool
		checkPort string
	}{
		{
			name: "Valid configuration",
			envVars: map[string]string{
				"YOUTRACK_TEST_URL":   "https://youtrack.example.com",
				"YOUTRACK_TEST_TOKEN": "token123",
				"PORT":                "9090",
			},
			wantError: false,
			checkPort: "9090",
		},
		{
			name: "Default port",
			envVars: map[string]string{
				"YOUTRACK_TEST_URL":   "https://youtrack.example.com",
				"YOUTRACK_TEST_TOKEN": "token123",
				"PORT":                "",
			},
			wantError: false,
			checkPort: "8080", // Default port
		},
		{
			name: "Missing YouTrack URL",
			envVars: map[string]string{
				"YOUTRACK_TEST_URL":   "",
				"YOUTRACK_TEST_TOKEN": "token123",
			},
			wantError: true,
		},
		{
			name: "Missing YouTrack token",
			envVars: map[string]string{
				"YOUTRACK_TEST_URL":   "https://youtrack.example.com",
				"YOUTRACK_TEST_TOKEN": "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for this test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Call Load
			cfg, err := Load()
			if (err != nil) != tt.wantError {
				t.Errorf("Load() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// Check port if no error expected
			if !tt.wantError && cfg.Port != tt.checkPort {
				t.Errorf("Load() port = %v, want %v", cfg.Port, tt.checkPort)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name: "Config with token",
			config: Config{
				Port:          "8080",
				YouTrackURL:   "https://youtrack.example.com",
				YouTrackToken: "abcdefghijklmnop",
			},
			want: "Config{Port: 8080, YouTrackURL: https://youtrack.example.com, YouTrackToken: abcd****}",
		},
		{
			name: "Config with short token",
			config: Config{
				Port:          "8080",
				YouTrackURL:   "https://youtrack.example.com",
				YouTrackToken: "abc",
			},
			want: "Config{Port: 8080, YouTrackURL: https://youtrack.example.com, YouTrackToken: ****}",
		},
		{
			name: "Config without token",
			config: Config{
				Port:          "8080",
				YouTrackURL:   "https://youtrack.example.com",
				YouTrackToken: "",
			},
			want: "Config{Port: 8080, YouTrackURL: https://youtrack.example.com, YouTrackToken: ****}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
