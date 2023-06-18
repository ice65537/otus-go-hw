package storage

import (
	"fmt"
	"time"
)

type Event struct {
	ID           string
	Title        string
	StartDt      time.Time
	StopDt       time.Time
	Desc         string
	Owner        string
	NotifyBefore time.Duration
}

func (evt Event) String() string {
	return fmt.Sprintf("%s[%s]%s", evt.Owner, evt.StartDt.Format(time.RFC3339), evt.Title)
}
