package engine

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func NewBoard() *Board {
	return &Board{
		Points:    map[PointID]Position{},
		Neighbors: map[PointID][]PointID{},
		Pieces:    map[PointID]Piece{},
	}
}

func (b *Board) String() string {
	ids := make([]int, 0, len(b.Points))
	for id := range b.Points {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	var out strings.Builder
	for _, rawID := range ids {
		id := PointID(rawID)
		position := b.Points[id]
		piece, occupied := b.PieceAt(id)
		if occupied {
			fmt.Fprintf(&out, "%d:(%d,%d)=%s\n", id, position.X, position.Y, piece.Owner)
			continue
		}
		fmt.Fprintf(&out, "%d:(%d,%d)=empty\n", id, position.X, position.Y)
	}
	return strings.TrimSpace(out.String())
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
	_, ok := b.Points[id]
	return ok
}

func (b *Board) PieceAt(id PointID) (Piece, bool) {
	piece, ok := b.Pieces[id]
	return piece, ok
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
	piece, ok := b.PieceAt(move.Source)
	if !ok {
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
	return slices.Contains(b.Neighbors[from], to)
}

func appendUniquePoint(points []PointID, point PointID) []PointID {
	if slices.Contains(points, point) {
		return points
	}
	return append(points, point)
}
