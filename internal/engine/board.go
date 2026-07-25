import "fmt"

type Board struct {
	Points map[PointID]Position
	Neighbors map[PointID][]PointID // Adjacency list
	Pieces map[PointID]Piece // Stores only occupied positions.
}

func NewBoard() *Board {
	return &Board{
		Points: make(map[PointID]Position),
		Neighbors: make(map[PointID][]PointID),
		Pieces: make(map[PointID]Piece),
	}
}

func (b *Board) String() string {
	return "cemas kau dek"
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

func (b *Board) AreConnected(from, to PointID) bool {
	for _, neighbor := range b.Neighbors[from] {
		if neighbor == to {
			return true
		}
	}

	return false
}
