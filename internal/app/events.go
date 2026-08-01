package app

import (
	"time"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/player"
)

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

type Event struct {
	Type      EventType
	PlayerID  player.ID
	Move      *engine.Move
	Message   string
	Timestamp time.Time
}

type EventHandler interface {
	HandleEvent(event Event) error
}
