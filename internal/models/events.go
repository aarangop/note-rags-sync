package models

import "time"

type FileChangedEvent struct {
	EventType string    `json:"event_type"`
	FilePath  string    `json:"file_path"`
	Content   string    `json:"content,omitempty"`
	Timestamp time.Time `json:timestamp`
}
