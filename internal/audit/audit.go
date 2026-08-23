package audit

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Event struct {
	Actor     string    `json:"actor"`
	Tenant    string    `json:"tenant"`
	Action    string    `json:"action"`
	ObjectID  string    `json:"object_id,omitempty"`
	Outcome   string    `json:"outcome"`
	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"ts"`
}

type Logger struct{}

func New() *Logger { return &Logger{} }

func (l *Logger) Emit(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return
	}
	log.New(os.Stdout, "audit ", log.LstdFlags).Print(string(b))
}
