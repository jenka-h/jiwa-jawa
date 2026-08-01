package session

import (
	"context"
	"time"

	err "jiwa-jawa/internal/error"
	"jiwa-jawa/internal/player"
)

// State describes the lifecycle of a Dam Daman network session.
type State uint8

const (
	StateIdle State = iota
	StateJoining
	StateConnected
	StateRecovering
	StateClosing
	StateClosed
)

// Participant describes one network participant in a Dam Daman session.
type Participant struct {
	Player  player.Player
	Address string
}

// Session stores local and peer metadata for one game session.
type Session struct {
	ID        uint64
	Local     Participant
	Peer      Participant
	State     State
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New creates a new session descriptor.
func New(id uint64, local Participant) *Session {
	return nil
}

// Join starts session negotiation with a peer.
func (s *Session) Join(ctx context.Context, peer Participant) error {
	return err.ErrNotImplemented
}

// Accept accepts a peer join request.
func (s *Session) Accept(peer Participant) error {
	return err.ErrNotImplemented
}

// MarkConnected marks the session as connected.
func (s *Session) MarkConnected() error {
	return err.ErrNotImplemented
}

// MarkRecovering marks the session as recovering after heartbeat/state mismatch.
func (s *Session) MarkRecovering() error {
	return err.ErrNotImplemented
}

// Close closes the session.
func (s *Session) Close() error {
	return err.ErrNotImplemented
}
