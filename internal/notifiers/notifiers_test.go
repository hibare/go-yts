package notifiers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hibare/go-yts/internal/config"
	"github.com/hibare/go-yts/internal/db"
)

func TestDiscord(t *testing.T) {
	ctx := context.Background()

	// Create test server to capture requests
	var receivedRequest *http.Request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRequest = r
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	tests := []struct {
		name          string
		movies        []db.Movies
		enabled       bool
		expectRequest bool
	}{
		{
			name:          "No movies",
			movies:        []db.Movies{},
			enabled:       true,
			expectRequest: false,
		},
		{
			name: "With movies",
			movies: []db.Movies{
				{
					Title:      "Test Movie",
					Year:       2023,
					Link:       "http://test.com",
					CoverImage: "http://test.com/image.jpg",
				},
			},
			enabled:       true,
			expectRequest: true,
		},
		{
			name: "Disabled notifier",
			movies: []db.Movies{
				{
					Title:      "Test Movie",
					Year:       2023,
					Link:       "http://test.com",
					CoverImage: "http://test.com/image.jpg",
				},
			},
			enabled:       false,
			expectRequest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset received request
			receivedRequest = nil

			// Set test environment
			os.Setenv("GO_YTS_NOTIFIER_DISCORD_ENABLED", fmt.Sprintf("%t", tt.enabled))
			os.Setenv("GO_YTS_NOTIFIER_DISCORD_WEBHOOK", ts.URL)
			config.LoadConfig()

			Discord(ctx, tt.movies)

			// Verify request expectations
			if tt.expectRequest && receivedRequest == nil {
				t.Error("Expected HTTP request, but none was made")
			} else if !tt.expectRequest && receivedRequest != nil {
				t.Error("Expected no HTTP request, but one was made")
			}
		})
	}
}
