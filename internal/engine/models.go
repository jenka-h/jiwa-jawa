package engine

import "fmt"

type PointID int

type Side string

const (
	SideNone Side = ""
	SideOne  Side = "player_one"
	SideTwo  Side = "player_two"
)

type Position struct {
	ID PointID
	X  int
	Y  int
}

type Piece struct {
	Owner Side
}

type Move struct {
	Source PointID
	Target PointID
}

func (m Move) String() string {
	return fmt.Sprintf("%d -> %d", m.Source, m.Target)
}

type Board struct {
	Points    map[PointID]Position
	Neighbors map[PointID][]PointID
	Pieces    map[PointID]Piece
}

type Phase uint8

const (
	PhaseNone Phase = iota
	PhaseJoin
	PhaseMove
	PhasePenalty
	PhaseFinish
)

type Event struct {
	Phase Phase
	Move  Move
}

type GameState struct {
	Board               *Board
	CurrentTurn         Side
	Phase               Phase
	MoveNumber          int
	Finished            bool
	Winner              Side
	PenaltySide         Side
	PenaltiesRemaining  int
	CaptureContinuation bool
}

type MoveResult struct {
	Move         Move
	Captured     []PointID
	NextTurn     Side
	GameOver     bool
	Winner       Side
	MustContinue bool
	DamOraMangan bool
}

type Coordinate struct{ X, Y int }

type EdgeSpec struct{ From, To Coordinate }

type TipSpec struct{ TipY, MidY, BaseY int }

type PlacementSpec struct {
	Coordinate
	Owner Side
}
