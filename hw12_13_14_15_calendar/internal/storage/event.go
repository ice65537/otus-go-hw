package storage

import (
	"fmt"
	"time"
)

type Event struct {
	ID           string
	Title        string
	StartDT      time.Time
	StopDT       time.Time
	Desc         string
	Owner        string
	NotifyBefore time.Duration
}

func (evt Event) String() string {
	return fmt.Sprintf("<%s>[%s]%s", evt.ID, evt.StartDT.Format("RFC822"), evt.Title)
}
