package gui

import (
	"context"

	"jiwa-jawa/internal/engine"
	err "jiwa-jawa/internal/error"
)

// UI describes a Dam Daman user interface.
type UI interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Render(state *engine.GameState) error
	Moves() <-chan engine.Move
	Errors() <-chan error
}

// CLI is a terminal-based Dam Daman interface.
type CLI struct{}

// NewCLI creates a CLI interface.
func NewCLI() *CLI {
	return nil
}

// Start starts the CLI input/render loop.
func (c *CLI) Start(ctx context.Context) error {
	return err.ErrNotImplemented
}

// Stop stops the CLI.
func (c *CLI) Stop(ctx context.Context) error {
	return err.ErrNotImplemented
}

// Render draws the current game state.
func (c *CLI) Render(state *engine.GameState) error {
	return err.ErrNotImplemented
}

// Moves returns parsed moves from user input.
func (c *CLI) Moves() <-chan engine.Move {
	return nil
}

// Errors returns asynchronous UI errors.
func (c *CLI) Errors() <-chan error {
	return nil
}
