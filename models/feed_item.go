package models

import "time"

type FeedItem struct {
	ProjectPath string
	Start       time.Time
	End         time.Time
	Duration    time.Duration
	Message     string
}
