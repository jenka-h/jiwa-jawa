package engine

import "fmt"

// ValidateMove checks whether a player can make a move on the current board.
func ValidateMove(board *Board, player Player, move Move) error {
	if board == nil {
		return fmt.Errorf("board is nil")
	}
	if player == PlayerNone {
		return fmt.Errorf("player is required")
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

	piece, occupied := board.PieceAt(move.Source)
	if !occupied {
		return fmt.Errorf("source point %d is empty", move.Source)
	}
	if piece.Owner != player {
		return fmt.Errorf("piece at source point %d belongs to %s, not %s", move.Source, piece.Owner, player)
	}
	if !board.IsEmpty(move.Target) {
		return fmt.Errorf("target point %d is occupied", move.Target)
	}

	return nil
}

// ApplyMove validates and applies a move to the board.
func ApplyMove(board *Board, player Player, move Move) error {
	if err := ValidateMove(board, player, move); err != nil {
		return err
	}

	return board.MovePiece(move)
}
