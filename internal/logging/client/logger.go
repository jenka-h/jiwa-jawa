package client

import (
	"context"
	"time"

	err "jiwa-jawa/internal/error"
	"jiwa-jawa/internal/player"
)

type EventType string

const (
	EventGameStarted EventType = "game_started"
	EventMove        EventType = "move"
	EventNetwork     EventType = "network"
	EventRetry       EventType = "retry"
	EventGameEnded   EventType = "game_ended"
	EventError       EventType = "error"
)

type Event struct {
	Type      EventType `json:"type"`
	SessionID uint64    `json:"session_id"`
	PlayerID  player.ID `json:"player_id,omitempty"`
	Message   string    `json:"message,omitempty"`
	Data      any       `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Logger interface {
	Log(ctx context.Context, event Event) error
	Close() error
}

type FileLogger struct{}

func NewFileLogger(path string) (*FileLogger, error) {
	return nil, err.ErrNotImplemented
}

func (l *FileLogger) Log(ctx context.Context, event Event) error {
	return err.ErrNotImplemented
}

func (l *FileLogger) Close() error {
	return err.ErrNotImplemented
}
