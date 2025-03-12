package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github-yt-webhook/internal/config"
	"github-yt-webhook/internal/handlers"
	"github-yt-webhook/internal/youtrack"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *config.Config
	server *http.Server
}

// New creates a new server instance
func New(config *config.Config) (*Server, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Create YouTrack client
	ytClient, err := youtrack.NewClient(config.YouTrackURL, config.YouTrackToken)
	if err != nil {
		return nil, err
	}

	// Create webhook handler
	webhookHandler := handlers.NewWebhookHandler(ytClient)

	// Set up Gin router
	router := gin.Default()

	// Register routes
	router.POST("/webhook", webhookHandler.HandleGitHubWebhook)

	// Add a health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Create server
	server := &Server{
		router: router,
		config: config,
		server: &http.Server{
			Addr:    ":" + config.Port,
			Handler: router,
		},
	}

	return server, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Server starting at port %s", s.config.Port)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
