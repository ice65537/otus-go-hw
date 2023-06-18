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
	return fmt.Sprintf("<%s>[%s]%s", evt.ID, evt.StartDt.Format("RFC822"), evt.Title)
}

func Dt2int(date time.Time) int {
	return date.Year()*1000 + date.YearDay()
}

func Dt2string(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
}
