package engine

import (
	"fmt"
	"sort"
	"strings"
)

// Board models the catur jawa board as a graph.
// Points are playable locations, Neighbors defines legal one-step movement,
// and Pieces stores only occupied points.
type Board struct {
	Points    map[PointID]Position
	Neighbors map[PointID][]PointID
	Pieces    map[PointID]Piece
}

func NewBoard() *Board {
	return &Board{
		Points:    make(map[PointID]Position),
		Neighbors: make(map[PointID][]PointID),
		Pieces:    make(map[PointID]Piece),
	}
}

func (b *Board) String() string {
	if b == nil {
		return "<nil board>"
	}

	ids := make([]int, 0, len(b.Points))
	for id := range b.Points {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	var builder strings.Builder
	for _, rawID := range ids {
		id := PointID(rawID)
		position := b.Points[id]
		piece, occupied := b.PieceAt(id)

		if occupied {
			fmt.Fprintf(&builder, "%d(%d,%d): %s\n", id, position.X, position.Y, piece.Owner)
		} else {
			fmt.Fprintf(&builder, "%d(%d,%d): empty\n", id, position.X, position.Y)
		}
	}

	return strings.TrimRight(builder.String(), "\n")
}

func (b *Board) AddPoint(id PointID, x, y int) {
	b.Points[id] = Position{ID: id, X: x, Y: y}
}

func (b *Board) Connect(a, c PointID) error {
	if !b.HasPoint(a) {
		return fmt.Errorf("point %d does not exist", a)
	}
	if !b.HasPoint(c) {
		return fmt.Errorf("point %d does not exist", c)
	}

	b.Neighbors[a] = appendUniquePoint(b.Neighbors[a], c)
	b.Neighbors[c] = appendUniquePoint(b.Neighbors[c], a)
	return nil
}

func (b *Board) HasPoint(id PointID) bool {
	_, exists := b.Points[id]
	return exists
}

func (b *Board) PieceAt(id PointID) (Piece, bool) {
	piece, occupied := b.Pieces[id]
	return piece, occupied
}

func (b *Board) IsEmpty(id PointID) bool {
	_, occupied := b.Pieces[id]
	return !occupied
}

func (b *Board) PlacePiece(id PointID, piece Piece) error {
	if !b.HasPoint(id) {
		return fmt.Errorf("point %d does not exist", id)
	}
	if !b.IsEmpty(id) {
		return fmt.Errorf("point %d is already occupied", id)
	}

	b.Pieces[id] = piece
	return nil
}

func (b *Board) MovePiece(move Move) error {
	piece, occupied := b.PieceAt(move.Source)
	if !occupied {
		return fmt.Errorf("source point %d is empty", move.Source)
	}
	if !b.IsEmpty(move.Target) {
		return fmt.Errorf("target point %d is occupied", move.Target)
	}

	delete(b.Pieces, move.Source)
	b.Pieces[move.Target] = piece
	return nil
}

func (b *Board) AreConnected(from, to PointID) bool {
	for _, neighbor := range b.Neighbors[from] {
		if neighbor == to {
			return true
		}
	}

	return false
}

func appendUniquePoint(points []PointID, point PointID) []PointID {
	for _, existing := range points {
		if existing == point {
			return points
		}
	}

	return append(points, point)
}
