package engine

import "fmt"

func NewGameState(board *Board, firstTurn Side) *GameState {
	return &GameState{Board: board, CurrentTurn: firstTurn, Phase: PhaseMove}
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
	s.CurrentTurn = NextSide(s.CurrentTurn)
	return nil
}

func (s *GameState) Finish(winner Side) {
	if s == nil {
		return
	}
	s.Finished = true
	s.Winner = winner
	s.Phase = PhaseFinish
}

func NextSide(side Side) Side {
	switch side {
	case SideOne:
		return SideTwo
	case SideTwo:
		return SideOne
	default:
		return SideNone
	}
}
