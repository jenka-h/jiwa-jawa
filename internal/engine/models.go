package engine

import "fmt"

type PointID int

type Player string

const (
	PlayerNone Player = ""
	PlayerOne  Player = "player_one"
	PlayerTwo  Player = "player_two"
)

type Position struct {
	ID PointID
	X  int
	Y  int
}

type Piece struct {
	Owner Player
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
	Board       *Board
	CurrentTurn Player
	Phase       Phase
	MoveNumber  int
	Finished    bool
	Winner      Player
}

type coordinate struct{ x, y int }

type edgeSpec struct{ from, to coordinate }

type tipSpec struct{ tipY, midY, baseY int }

type placementSpec struct {
	coordinate
	owner Player
}
