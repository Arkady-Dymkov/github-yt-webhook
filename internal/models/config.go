package models

type YouTrackCommand struct {
	Command string `json:"command"`
	Comment string `json:"comment"`
}

type EventMapping struct {
	EventName     string                  `json:"eventName"`
	GitHubActions map[string]GitHubAction `json:"githubActions"`
}

// GitHubAction defines a mapping between a GitHub action and a YouTrack command
type GitHubAction struct {
	GitHubAction    string            `json:"githubAction"`
	YouTrackCommand []YouTrackCommand `json:"youtrackCommand"`
}

// ActionMappingConfig holds the configuration for GitHub action to YouTrack command mappings
type ActionMappingConfig struct {
	Mappings []EventMapping `json:"mappings"`
}
