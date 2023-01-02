package models

import (
	"github.com/google/uuid"
	"time"
)

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PhotoUrl  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}
