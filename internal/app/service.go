package app

import (
	"context"

	"jiwa-jaw/internal/engine"
)

// Service coordinates Dam Daman game logic, transport, UI, and logging.
type Service struct {
	config Config
	state  *engine.GameState
}

// NewService creates an application service.
func NewService(config Config) (*Service, error) {
	return nil, ErrNotImplemented
}

// Run starts the client application loop.
func (s *Service) Run(ctx context.Context) error {
	return ErrNotImplemented
}

// Stop shuts down the application service.
func (s *Service) Stop(ctx context.Context) error {
	return ErrNotImplemented
}

// Join starts or requests a Dam Daman session with a peer.
func (s *Service) Join(ctx context.Context) error {
	return ErrNotImplemented
}

// ApplyLocalMove validates, applies, sends, and logs a local move.
func (s *Service) ApplyLocalMove(move engine.Move) error {
	return ErrNotImplemented
}

// ApplyRemoteMove validates, applies, and logs a remote move.
func (s *Service) ApplyRemoteMove(move engine.Move) error {
	return ErrNotImplemented
}
