# GitHub YouTrack Webhook

A Golang application that updates YouTrack issues based on GitHub pull request events.

## Features


- Listens for GitHub webhook events
- Updates YouTrack issue status to "In Review" when a pull request is opened
- Extracts ticket ID from pull request title (format: ABC-123)
- Built with Gin framework for high performance

## Requirements

- Go 1.20 or higher
- Docker (optional, for containerized deployment)

## Configuration

The application uses environment variables for configuration:

- `PORT`: HTTP server port (default: 8080)
- `YOUTRACK_TEST_URL`: YouTrack API URL
- `YOUTRACK_TEST_TOKEN`: YouTrack API token

## Development

### Building

```bash
# Build the application
make build

# Run tests
make test
```
# Running
```bash
# Set required environment variables
export YOUTRACK_TEST_URL=https://youtrack.example.com
export YOUTRACK_TEST_TOKEN=your-api-token

# Run the application
./github-webhook-youtrack
```

# Docker Deployment

**Building the Docker image**
```bash
make docker-build
```
**Run with docker**
```bash
make docker-run
```
**Running docker compose**
```bash
# Start the application
make docker-up

# Stop the application
make docker-down
```
