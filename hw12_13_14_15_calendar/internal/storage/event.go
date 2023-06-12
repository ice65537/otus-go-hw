package storage

import "time"

type Event struct {
	ID           string
	Title        string
	DateTime     time.Time
	Duration     time.Duration
	Desc         string
	Owner        string
	NotifyBefore time.Duration
}
