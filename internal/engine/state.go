package engine

import "fmt"

func NewGameState(board *Board, firstTurn Side) *GameState {
	return &GameState{Board: board, CurrentTurn: firstTurn, Phase: PhaseMove}
}

func (s *GameState) ApplyMove(move Move) error {
	_, err := s.ApplyMoveResult(move)
	return err
}

func (s *GameState) ApplyMoveResult(move Move) (MoveResult, error) {
	if s == nil {
		return MoveResult{}, fmt.Errorf("game state is nil")
	}
	if s.Finished {
		return MoveResult{}, fmt.Errorf("game is already finished")
	}
	if s.Phase != PhaseMove {
		return MoveResult{}, fmt.Errorf("game is not in move phase")
	}

	movingSide := s.CurrentTurn
	if s.CaptureContinuation {
		if _, ok := captureTarget(s.Board, move.Source, move.Target, movingSide); !ok {
			return MoveResult{}, fmt.Errorf("continuation move must capture with an eligible piece or end the turn")
		}
	}
	result, err := ApplyMoveWithResult(s.Board, movingSide, move)
	if err != nil {
		return MoveResult{}, err
	}

	s.MoveNumber++
	if result.GameOver {
		s.Finish(result.Winner)
		result.NextTurn = SideNone
		return result, nil
	}

	s.CurrentTurn = result.NextTurn
	s.CaptureContinuation = result.MustContinue
	if result.DamOraMangan {
		s.Phase = PhasePenalty
		s.PenaltySide = s.CurrentTurn
		s.PenaltiesRemaining = min(3, CountPieces(s.Board, movingSide))
	}
	return result, nil
}

func (s *GameState) EndTurn(side Side) error {
	if s == nil {
		return fmt.Errorf("game state is nil")
	}
	if s.Finished {
		return fmt.Errorf("game is already finished")
	}
	if s.Phase != PhaseMove || !s.CaptureContinuation {
		return fmt.Errorf("there is no optional capture continuation to end")
	}
	if side != s.CurrentTurn {
		return fmt.Errorf("not %s's turn", side)
	}
	s.CaptureContinuation = false
	s.CurrentTurn = Opponent(side)
	return nil
}

func (s *GameState) ApplyPenalty(side Side, target PointID) error {
	if s == nil {
		return fmt.Errorf("game state is nil")
	}
	if s.Finished {
		return fmt.Errorf("game is already finished")
	}
	if s.Phase != PhasePenalty || s.PenaltiesRemaining <= 0 {
		return fmt.Errorf("game is not in Dam Ora Mangan penalty phase")
	}
	if side != s.CurrentTurn || side != s.PenaltySide {
		return fmt.Errorf("penalty belongs to %s", s.PenaltySide)
	}
	piece, ok := s.Board.PieceAt(target)
	if !ok {
		return fmt.Errorf("penalty target %d is empty", target)
	}
	offender := Opponent(side)
	if piece.Owner != offender {
		return fmt.Errorf("penalty target %d must belong to %s", target, offender)
	}
	delete(s.Board.Pieces, target)
	s.PenaltiesRemaining--
	if CountPieces(s.Board, offender) == 0 {
		s.Finish(side)
		return nil
	}
	if s.PenaltiesRemaining == 0 {
		s.Phase = PhaseMove
		s.PenaltySide = SideNone
	}
	return nil
}

func (s *GameState) Finish(winner Side) {
	if s == nil {
		return
	}
	s.Finished = true
	s.Winner = winner
	s.Phase = PhaseFinish
	s.PenaltySide = SideNone
	s.PenaltiesRemaining = 0
	s.CaptureContinuation = false
}

func NextSide(side Side) Side {
	return Opponent(side)
}
