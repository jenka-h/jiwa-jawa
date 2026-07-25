// Domain Types

type Player struct{

}

type Node struct{
	Owner string
	IsOccupied bool
}

type Position struct{
	ID int
	X int
	Y int
}

type Move struct {
	source Position
	target Position
	piece  Node
}


// Phase
type Phase uint8

const (
	PhaseNone Phase = iota + 1
	PhaseJoin
	PhaseMove
	PhasePenalty
	PhaseFinish
)

// Event

type Event struct {
	Phase Phase
	Move  Move

}
