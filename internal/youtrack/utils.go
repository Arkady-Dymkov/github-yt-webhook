package youtrack

import (
	"regexp"
	"strings"
)

// ExtractTicket extracts a ticket ID from a pull request title
func ExtractTicket(title string) string {
	// Assuming ticket format like ABC-123
	re := regexp.MustCompile(`[A-Z]+-\d+`)
	ticket := re.FindString(title)
	return strings.TrimSpace(ticket)
}
