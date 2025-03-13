package clients

import (
	"github-yt-webhook/internal/clients/youtrack"
	"github-yt-webhook/internal/config"
)

type Clients struct {
	YouTrackClient *youtrack.Client
}

func CreateClients(c *config.Config) *Clients {
	return &Clients{
		YouTrackClient: youtrack.NewClient(c.YouTrackURL, c.YouTrackToken),
	}
}
