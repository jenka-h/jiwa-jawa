package engine

import "fmt"

// GameState keeps the mutable engine state for one game.
type GameState struct {
	Board       *Board
	CurrentTurn Player
	Phase       Phase
	MoveNumber  int
	Finished    bool
	Winner      Player
}

func NewGameState(board *Board, firstTurn Player) *GameState {
	return &GameState{
		Board:       board,
		CurrentTurn: firstTurn,
		Phase:       PhaseMove,
	}
}

func (s *GameState) ApplyMove(move Move) error {
	if s == nil {
		return fmt.Errorf("game state is nil")
	}
	if s.Finished {
		return fmt.Errorf("game is already finished")
	}
	if s.Phase != PhaseMove {
		return fmt.Errorf("game is not in move phase")
	}

	if err := ApplyMove(s.Board, s.CurrentTurn, move); err != nil {
		return err
	}

	s.MoveNumber++
	s.CurrentTurn = NextPlayer(s.CurrentTurn)
	return nil
}

func (s *GameState) Finish(winner Player) {
	if s == nil {
		return
	}

	s.Finished = true
	s.Winner = winner
	s.Phase = PhaseFinish
}

func NextPlayer(player Player) Player {
	switch player {
	case PlayerOne:
		return PlayerTwo
	case PlayerTwo:
		return PlayerOne
	default:
		return PlayerNone
	}
}
