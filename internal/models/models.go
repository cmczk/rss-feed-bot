package models

import "time"

type RSSource struct {
	ID        int64
	Name      string
	FeedURL   string
	CreatedAt time.Time
}

// Material as part of RSS feed
type Item struct {
	Title        string
	Categories   []string
	Link         string
	Date         time.Time
	Summary      string
	RSSourceName string
}

// Material as part of business logic
type Material struct {
	ID          int64
	RSSourceID  int64
	Title       string
	Link        string
	Summary     string
	PublishedAt time.Time
	CreatedAt   time.Time
	PostedAt    time.Time
}
