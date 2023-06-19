package storage

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	ID           string    `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	StartDt      time.Time `json:"startdt,omitempty"`
	StopDt       time.Time `json:"stopdt,omitempty"`
	Desc         string    `json:"desc,omitempty"`
	Owner        string    `json:"owner,omitempty"`
	NotifyBefore int       `json:"notifybefore,omitempty"`
}

func (evt Event) String() string {
	return fmt.Sprintf("%s[%s]%s", evt.Owner, evt.StartDt.Format(time.RFC3339), evt.Title)
}

func Unmarshal(data []byte) (Event, error) {
	var evt Event
	err := json.Unmarshal(data, &evt)
	return evt, err
}

func Marshall(events []Event) ([]byte, error) {
	data, err := json.Marshal(&events)
	return data, err
}

/*curl localhost:1234/event/new -H User:ice65588 \
-d '{"title":"x","startdt":"2023-01-10T19:00:00.000Z","stopdt":"2023-01-10T19:00:00.000Z","notifybefore":600}'*/
