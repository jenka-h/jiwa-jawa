package engine

import (
	"fmt"
	"sort"
)

func ValidateMove(board *Board, side Side, move Move) error {
	_, err := classifyMove(board, side, move)
	return err
}

func ApplyMove(board *Board, side Side, move Move) error {
	_, err := ApplyMoveWithResult(board, side, move)
	return err
}

func ApplyMoveWithResult(board *Board, side Side, move Move) (MoveResult, error) {
	captureAvailable := HasCapture(board, side)
	captured, err := classifyMove(board, side, move)
	if err != nil {
		return MoveResult{}, err
	}

	piece, _ := board.PieceAt(move.Source)
	delete(board.Pieces, move.Source)
	for _, point := range captured {
		delete(board.Pieces, point)
	}
	board.Pieces[move.Target] = piece

	winner := SideNone
	gameOver := false
	opponent := Opponent(side)
	if CountPieces(board, opponent) == 0 {
		winner = side
		gameOver = true
	}

	mustContinue := false
	nextTurn := opponent
	if len(captured) > 0 && !gameOver && HasCapture(board, side) {
		mustContinue = true
		nextTurn = side
	}

	return MoveResult{
		Move:         move,
		Captured:     captured,
		NextTurn:     nextTurn,
		GameOver:     gameOver,
		Winner:       winner,
		MustContinue: mustContinue,
		DamOraMangan: captureAvailable && len(captured) == 0,
	}, nil
}

func HasCapture(board *Board, side Side) bool {
	if board == nil || side == SideNone {
		return false
	}
	for point, piece := range board.Pieces {
		if piece.Owner == side && HasCaptureFrom(board, point) {
			return true
		}
	}
	return false
}

func CaptureSources(board *Board, side Side) []PointID {
	if board == nil || side == SideNone {
		return nil
	}
	result := make([]PointID, 0)
	for point, piece := range board.Pieces {
		if piece.Owner == side && HasCaptureFrom(board, point) {
			result = append(result, point)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func HasCaptureFrom(board *Board, source PointID) bool {
	if board == nil {
		return false
	}
	piece, ok := board.PieceAt(source)
	if !ok || piece.Owner == SideNone {
		return false
	}
	return len(captureTargets(board, source, piece.Owner)) > 0
}

func Opponent(side Side) Side {
	switch side {
	case SideOne:
		return SideTwo
	case SideTwo:
		return SideOne
	default:
		return SideNone
	}
}

func CountPieces(board *Board, side Side) int {
	if board == nil {
		return 0
	}
	count := 0
	for _, piece := range board.Pieces {
		if piece.Owner == side {
			count++
		}
	}
	return count
}

func classifyMove(board *Board, side Side, move Move) ([]PointID, error) {
	if board == nil {
		return nil, fmt.Errorf("board is nil")
	}
	if side == SideNone {
		return nil, fmt.Errorf("side is required")
	}
	if !board.HasPoint(move.Source) {
		return nil, fmt.Errorf("source point %d does not exist", move.Source)
	}
	if !board.HasPoint(move.Target) {
		return nil, fmt.Errorf("target point %d does not exist", move.Target)
	}

	piece, ok := board.PieceAt(move.Source)
	if !ok {
		return nil, fmt.Errorf("source point %d is empty", move.Source)
	}
	if piece.Owner != side {
		return nil, fmt.Errorf("piece at source point %d belongs to %s, not %s", move.Source, piece.Owner, side)
	}
	if !board.IsEmpty(move.Target) {
		return nil, fmt.Errorf("target point %d is occupied", move.Target)
	}

	if captured, ok := captureTarget(board, move.Source, move.Target, side); ok {
		return []PointID{captured}, nil
	}

	if !board.AreConnected(move.Source, move.Target) {
		return nil, fmt.Errorf("source point %d is not connected to target point %d", move.Source, move.Target)
	}

	return nil, nil
}

func captureTargets(board *Board, source PointID, side Side) map[PointID]PointID {
	result := map[PointID]PointID{}
	sourcePos, ok := board.Points[source]
	if !ok {
		return result
	}
	for _, middle := range board.Neighbors[source] {
		middlePiece, occupied := board.PieceAt(middle)
		if !occupied || middlePiece.Owner != Opponent(side) {
			continue
		}
		middlePos := board.Points[middle]
		targetX := middlePos.X + (middlePos.X - sourcePos.X)
		targetY := middlePos.Y + (middlePos.Y - sourcePos.Y)
		target, ok := pointAt(board, targetX, targetY)
		if !ok || !board.AreConnected(middle, target) || !board.IsEmpty(target) {
			continue
		}
		result[target] = middle
	}
	return result
}

func captureTarget(board *Board, source, target PointID, side Side) (PointID, bool) {
	captured, ok := captureTargets(board, source, side)[target]
	return captured, ok
}

func pointAt(board *Board, x, y int) (PointID, bool) {
	for id, position := range board.Points {
		if position.X == x && position.Y == y {
			return id, true
		}
	}
	return 0, false
}
