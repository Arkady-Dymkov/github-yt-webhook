# GitHub YouTrack Webhook Postman Collection

This directory contains Postman collection and environment files for testing the GitHub YouTrack webhook integration.

## Files

- `github-webhook-youtrack.postman_collection.json`: Postman collection with test requests
- `github-webhook-youtrack.postman_environment.json`: Postman environment variables

## How to Use

1. Import the collection and environment into Postman:
    - Open Postman
    - Click "Import" button
    - Select both JSON files
    - Click "Import"

2. Select the "GitHub YouTrack Webhook - Local" environment from the dropdown in the top-right corner

3. Run the application locally:

```bash
go build -o github-yt-webhook
export YOUTRACK_TEST_URL=https://youtrack.example.com
export YOUTRACK_TEST_TOKEN=your-api-token
./github-yt-webhook
```

4. Use the collection to test different webhook scenarios:
    - Health Check: Verify the service is running
    - Pull Request Opened: Test a valid PR with ticket ID
    - Pull Request Closed: Test a closed PR with ticket ID
    - Pull Request Reopened: Test a reopened PR with ticket ID
    - Pull Request No Ticket: Test a PR without a ticket ID
    - Push Event: Test that push events are ignored
    - Invalid JSON: Test error handling for malformed requests

## Test Scenarios

**Health Check**

- Expected Result: 200 OK with status information

**Pull Request Opened**

- Expected Result: 200 OK with message "Issue ABC-123 updated"
- YouTrack issue ABC-123 should be updated to "In Progress" status

**Pull Request Closed**

- Expected Result: 200 OK with message "Issue ABC-123 updated"
- YouTrack issue ABC-123 should be updated to "Done" status

**Pull Request Reopened**

- Expected Result: 200 OK with message "Issue ABC-123 updated"
- YouTrack issue ABC-123 should be updated to "In Progress" status

**Pull Request No Ticket**

- Expected Result: 200 OK with message "No ticket found"
- No changes should be made to YouTrack issues

**Push Event**

- Expected Result: 200 OK with message "Event ignored"
- No changes should be made to YouTrack issues

**Invalid JSON**

- Expected Result: 400 Bad Request with error message
- No changes should be made to YouTrack issues

## Troubleshooting

**If you encounter issues:**

1. Check that the application is running and accessible at http://localhost:8080
2. Verify your YouTrack URL and token are correctly set
3. Check the application logs for detailed error messages
4. Ensure your action_mappings.json file is properly configured