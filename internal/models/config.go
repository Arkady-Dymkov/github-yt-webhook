package models

type YouTrackCommand struct {
	Command string `json:"command"`
	Comment string `json:"comment"`
}

// ActionMapping defines a mapping between a GitHub action and a YouTrack command
type ActionMapping struct {
	GitHubAction    string            `json:"githubAction"`
	YouTrackCommand []YouTrackCommand `json:"youtrackCommand"`
}

// ActionMappingConfig holds the configuration for GitHub action to YouTrack command mappings
type ActionMappingConfig struct {
	Mappings []ActionMapping `json:"mappings"`
}
