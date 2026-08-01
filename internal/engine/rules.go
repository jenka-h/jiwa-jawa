package engine

import "fmt"

func ValidateMove(board *Board, side Side, move Move) error {
	if board == nil {
		return fmt.Errorf("board is nil")
	}
	if side == SideNone {
		return fmt.Errorf("side is required")
	}
	if !board.HasPoint(move.Source) {
		return fmt.Errorf("source point %d does not exist", move.Source)
	}
	if !board.HasPoint(move.Target) {
		return fmt.Errorf("target point %d does not exist", move.Target)
	}
	if !board.AreConnected(move.Source, move.Target) {
		return fmt.Errorf("source point %d is not connected to target point %d", move.Source, move.Target)
	}

	piece, ok := board.PieceAt(move.Source)
	if !ok {
		return fmt.Errorf("source point %d is empty", move.Source)
	}
	if piece.Owner != side {
		return fmt.Errorf("piece at source point %d belongs to %s, not %s", move.Source, piece.Owner, side)
	}
	if !board.IsEmpty(move.Target) {
		return fmt.Errorf("target point %d is occupied", move.Target)
	}

	return nil
}

func ApplyMove(board *Board, side Side, move Move) error {
	if err := ValidateMove(board, side, move); err != nil {
		return err
	}
	return board.MovePiece(move)
}
