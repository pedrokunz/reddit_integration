package domain

import "time"

type Post struct {
	ID       string
	Title    string
	Author   string
	Origin   Integration
	SyncedAt time.Time
	Metadata map[string]any
}
