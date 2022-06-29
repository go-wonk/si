package testmodels

import (
	"encoding/json"
	"time"
)

type Document struct {
	Name      string    `json:"name"`
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

func (o *Document) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}
