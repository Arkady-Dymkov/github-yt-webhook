# Build Instructions

## Prerequisites
- Go (version 1.20 or later) must be installed and available in your PATH.
- Git (optional, for version control).

## Build and Run
1. Clone the repository:
   git clone https://github.com/yourusername/github-webhook-youtrack.git
2. Change directory to the project folder:
   cd github-webhook-youtrack
3. Build the project:
   go build -o github-webhook-youtrack main.go router.go youtrack.go
4. Run the application:
   ./github-webhook-youtrack

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

## YouTrack Integration
To enable updating YouTrack issues, set the following environment variables:
- YT_URL: The base URL of your YouTrack instance (for example, "http://youtrack.example.com").
- YT_TOKEN: (Optional) The authentication token to access YouTrack API.

## Running Tests
To run tests, execute:
   go test -v