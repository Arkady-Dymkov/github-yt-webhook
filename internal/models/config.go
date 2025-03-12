package models

// ActionMapping defines a mapping between a GitHub action and a YouTrack command
type ActionMapping struct {
	GitHubAction    string          `json:"githubAction"`
	YouTrackCommand YouTrackCommand `json:"youtrackCommand"`
}

type YouTrackCommand struct {
	Command string
	Comment string
}

// ActionMappingConfig holds the configuration for GitHub action to YouTrack command mappings
type ActionMappingConfig struct {
	Mappings []ActionMapping `json:"mappings"`
}
