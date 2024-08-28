package quiz

import (
	"backendProject/internal/spotify"
	"time"
)

type Quiz struct {
	Track     spotify.Track `json:"track"`
	CreatedAt time.Time     `json:"created_at"`
}
