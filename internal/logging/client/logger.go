package client

import (
	"context"
	"time"
)

// EventType identifies log event categories for Dam Daman.
type EventType string

const (
	EventGameStarted EventType = "game_started"
	EventMove        EventType = "move"
	EventNetwork     EventType = "network"
	EventRetry       EventType = "retry"
	EventGameEnded   EventType = "game_ended"
	EventError       EventType = "error"
)

// Event is one local or remote log entry.
type Event struct {
	Type      EventType `json:"type"`
	SessionID uint64    `json:"session_id"`
	PlayerID  string    `json:"player_id,omitempty"`
	Message   string    `json:"message,omitempty"`
	Data      any       `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Logger writes Dam Daman events to local storage or a remote log service.
type Logger interface {
	Log(ctx context.Context, event Event) error
	Close() error
}

// FileLogger is a JSONL local logger.
type FileLogger struct{}

// NewFileLogger creates a local JSONL logger.
func NewFileLogger(path string) (*FileLogger, error) {
	return nil, ErrNotImplemented
}

// Log writes one event.
func (l *FileLogger) Log(ctx context.Context, event Event) error {
	return ErrNotImplemented
}

// Close closes the logger.
func (l *FileLogger) Close() error {
	return ErrNotImplemented
}
