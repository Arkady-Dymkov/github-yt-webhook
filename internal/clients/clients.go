package clients

import (
	"github-yt-webhook/internal/clients/youtrack"
	"github-yt-webhook/internal/config"
)

type Clients struct {
	YouTrackClient *youtrack.Client
}

func CreateClients(c *config.Config) (*Clients, error) {

	yt, err := youtrack.NewClient(c.YouTrackURL, c.YouTrackToken)
	if err != nil {
		return nil, err
	}

	return &Clients{
		YouTrackClient: &yt,
	}, nil
}
