# GitHub YouTrack Webhook

A Golang application that updates YouTrack issues based on GitHub pull request events.

## Features

- Listens for GitHub webhook events
- Updates YouTrack issue status when a pull request is opened, closed, or reopened
- Extracts ticket ID from pull request title (format: ABC-123)
- Configurable action mappings between GitHub events and YouTrack commands
- Built with Gin framework for high performance
- Comprehensive logging for debugging and monitoring

## Requirements

- Go 1.20 or higher
- Docker (optional, for containerized deployment)

## Configuration

The application uses environment variables for configuration:

- `PORT`: HTTP server port (default: 8080)
- `YOUTRACK_TEST_URL`: YouTrack API URL
- `YOUTRACK_TEST_TOKEN`: YouTrack API token
- `ACTION_MAPPINGS_CONFIG`: Path to action mappings config file (default: "action_mappings.json")

### Action Mappings

The application uses a JSON configuration file to map GitHub pull request actions to YouTrack commands.
Create an `action_mappings.json` file with the following structure:

```json
{
    "mappings": [
        {
            "eventName": "pull_request",
            "githubActions": [
                {
                    "githubAction": "opened",
                    "youtrackCommand": [
                        {
                            "command": "In Progress",
                            "comment": "PR opened - status updated by GitHub webhook"
                        },
                        {
                            "command": "Comment",
                            "comment": "Check commit here: {{html_url}}"
                        }
                    ]
                },
                {
                    "githubAction": "closed",
                    "youtrackCommand": [
                        {
                            "command": "Done",
                            "comment": "PR closed - status updated by GitHub webhook"
                        }
                    ]
                },
                {
                    "githubAction": "reopened",
                    "youtrackCommand": [
                        {
                            "command": "In Progress",
                            "comment": "PR reopened - status updated by GitHub webhook"
                        }
                    ]
                }
            ]
        },
        {
            "eventName": "push",
            "githubActions": [
                {
                    "githubAction": "any",
                    "youtrackCommand": [
                        {
                            "command": "Comment",
                            "comment": "New push made: \n by: {{commit_author.name}}\n commit message: {{commit_message}}\n commit hash: {{commit_id}}"
                        }
                    ]
                }
            ]
        }
    ]
}
```