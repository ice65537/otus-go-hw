package storage

import (
	"fmt"
	"time"
)

type Event struct {
	ID           string
	Title        string
	DateTime     time.Time
	Duration     time.Duration
	Desc         string
	Owner        string
	NotifyBefore time.Duration
}

func (evt Event) String() string {
	return fmt.Sprintf("[%s]<%s>%s", evt.ID, evt.DateTime.Format("RFC822"), evt.Title)
}
