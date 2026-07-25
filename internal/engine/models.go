package engine

import "fmt"

// PointID identifies a playable point/node on the board graph.
type PointID int

// Player identifies a player in the game.
type Player string

const (
	PlayerNone Player = ""
	PlayerOne  Player = "player_one"
	PlayerTwo  Player = "player_two"
)

// Position stores display coordinates for a point on the board.
type Position struct {
	ID PointID
	X  int
	Y  int
}

// Piece represents one piece owned by a player.
type Piece struct {
	Owner Player
}

// Move represents a move from one board point to another.
type Move struct {
	Source PointID
	Target PointID
}

func (m Move) String() string {
	return fmt.Sprintf("%d -> %d", m.Source, m.Target)
}

// Phase describes the current game phase.
type Phase uint8

const (
	PhaseNone Phase = iota
	PhaseJoin
	PhaseMove
	PhasePenalty
	PhaseFinish
)

// Event describes an engine-level event that can be logged or sent to the app layer.
type Event struct {
	Phase Phase
	Move  Move
}
