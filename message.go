package tailor

import (
	"encoding/json"
	"time"
)

const dateTimeFormat = "2006-01-02T03:04:05Z0700"

// Message represents an entry read from a logfile
type Message struct {
	Time   time.Time `json:"time"`
	Source string    `json:"source"`
	Body   string    `json:"body"`
}

// MarshalJSON implements the JSON marshaler interface, used for customized timestamp formatting
func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		Time string `json:"time"`
		*Alias
	}{
		Time:  m.Time.Format(dateTimeFormat),
		Alias: (*Alias)(m),
	})
}

// NewMessage creates a new message with a UTC timestamp
func NewMessage(source, body string) *Message {
	return &Message{Time: time.Now().UTC(), Source: source, Body: body}
}
