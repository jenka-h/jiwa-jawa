package app

import (
	"time"

	"jiwa-jawa/internal/engine"
)

// EventType identifies application events that can be rendered or logged.
type EventType string

const (
	EventJoined       EventType = "joined"
	EventAccepted     EventType = "accepted"
	EventMoveSent     EventType = "move_sent"
	EventMoveReceived EventType = "move_received"
	EventMoveApplied  EventType = "move_applied"
	EventGameFinished EventType = "game_finished"
	EventDisconnected EventType = "disconnected"
	EventError        EventType = "error"
)

// Event is a high-level Dam Daman application event.
type Event struct {
	Type      EventType
	PlayerID  string
	Move      *engine.Move
	Message   string
	Timestamp time.Time
}

// EventHandler handles application events.
type EventHandler interface {
	HandleEvent(event Event) error
}
