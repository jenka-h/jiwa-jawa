package engine

import "fmt"

var (
	defaultPoints     = BuildDefaultPoints()
	pointByPosition   = IndexPoints(defaultPoints)
	defaultEdges      = BuildDefaultEdges()
	defaultPlacements = BuildDefaultPlacements()
)

func NewEmptyBoard() (*Board, error) {
	board := NewBoard()
	for _, point := range defaultPoints {
		board.AddPoint(point.ID, point.X, point.Y)
	}
	for _, edge := range defaultEdges {
		from := MustPointIDByPosition(edge.From.X, edge.From.Y)
		to := MustPointIDByPosition(edge.To.X, edge.To.Y)
		if err := board.Connect(from, to); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func NewDefaultBoard() (*Board, error) {
	board, err := NewEmptyBoard()
	if err != nil {
		return nil, err
	}
	for _, placement := range defaultPlacements {
		id := MustPointIDByPosition(placement.X, placement.Y)
		if err := board.PlacePiece(id, Piece{Owner: placement.Owner}); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func PointIDByPosition(x, y int) (PointID, error) {
	id, ok := pointByPosition[Coordinate{x, y}]
	if !ok {
		return 0, fmt.Errorf("no default board point at (%d, %d)", x, y)
	}
	return id, nil
}

func MustPointIDByPosition(x, y int) PointID {
	id, err := PointIDByPosition(x, y)
	if err != nil {
		panic(err)
	}
	return id
}

func BuildDefaultPoints() []Position {
	coords := []Coordinate{
		// left home triangle
		{0, 1}, {0, 2}, {0, 3}, {1, 1}, {1, 2}, {1, 3},
	}

	// main board: two left columns, five center column points, two right columns
	for _, x := range []int{2, 3, 4, 5, 6} {
		for y := 0; y <= 4; y++ {
			coords = append(coords, Coordinate{x, y})
		}
	}

	coords = append(coords,
		// right home triangle
		Coordinate{7, 1}, Coordinate{7, 2}, Coordinate{7, 3},
		Coordinate{8, 1}, Coordinate{8, 2}, Coordinate{8, 3},
	)

	points := make([]Position, len(coords))
	for i, coord := range coords {
		points[i] = Position{ID: PointID(i + 1), X: coord.X, Y: coord.Y}
	}
	return points
}

func BuildDefaultEdges() []EdgeSpec {
	points := map[Coordinate]struct{}{}
	for _, point := range defaultPoints {
		points[Coordinate{point.X, point.Y}] = struct{}{}
	}

	edges := map[EdgeSpec]struct{}{}
	add := func(a, b Coordinate) {
		if _, ok := points[a]; !ok {
			return
		}
		if _, ok := points[b]; !ok {
			return
		}
		if LessCoordinate(b, a) {
			a, b = b, a
		}
		edges[EdgeSpec{From: a, To: b}] = struct{}{}
	}

	// horizontal and vertical lines in the main body
	for _, x := range []int{2, 3, 4, 5, 6} {
		for y := 0; y < 4; y++ {
			add(Coordinate{x, y}, Coordinate{x, y + 1})
		}
	}
	for y := 0; y <= 4; y++ {
		for x := 2; x < 6; x++ {
			add(Coordinate{x, y}, Coordinate{x + 1, y})
		}
	}

	// Diagonals pass through the center node of each 2x2 square. Keep
	// every node as an explicit graph step so diagonal captures can jump
	// source -> opponent -> landing and continue from that landing point.
	for _, x := range []int{2, 4} {
		for _, y := range []int{0, 2} {
			center := Coordinate{x + 1, y + 1}
			add(Coordinate{x, y}, center)
			add(center, Coordinate{x + 2, y + 2})
			add(Coordinate{x + 2, y}, center)
			add(center, Coordinate{x, y + 2})
		}
	}

	// left triangle
	add(Coordinate{0, 1}, Coordinate{0, 2})
	add(Coordinate{0, 2}, Coordinate{0, 3})
	add(Coordinate{0, 1}, Coordinate{1, 1})
	add(Coordinate{0, 2}, Coordinate{1, 2})
	add(Coordinate{0, 3}, Coordinate{1, 3})
	add(Coordinate{1, 1}, Coordinate{1, 2})
	add(Coordinate{1, 2}, Coordinate{1, 3})
	add(Coordinate{1, 1}, Coordinate{2, 2})
	add(Coordinate{1, 2}, Coordinate{2, 2})
	add(Coordinate{1, 3}, Coordinate{2, 2})

	// right triangle
	add(Coordinate{6, 2}, Coordinate{7, 1})
	add(Coordinate{6, 2}, Coordinate{7, 2})
	add(Coordinate{6, 2}, Coordinate{7, 3})
	add(Coordinate{7, 1}, Coordinate{7, 2})
	add(Coordinate{7, 2}, Coordinate{7, 3})
	add(Coordinate{7, 1}, Coordinate{8, 1})
	add(Coordinate{7, 2}, Coordinate{8, 2})
	add(Coordinate{7, 3}, Coordinate{8, 3})
	add(Coordinate{8, 1}, Coordinate{8, 2})
	add(Coordinate{8, 2}, Coordinate{8, 3})

	result := make([]EdgeSpec, 0, len(edges))
	for edge := range edges {
		result = append(result, edge)
	}
	return result
}

func BuildDefaultPlacements() []PlacementSpec {
	placements := []PlacementSpec{}
	add := func(owner Side, coords ...Coordinate) {
		for _, coord := range coords {
			placements = append(placements, PlacementSpec{Coordinate: coord, Owner: owner})
		}
	}

	// SideOne starts on the left side.
	add(SideOne,
		Coordinate{0, 1}, Coordinate{0, 2}, Coordinate{0, 3},
		Coordinate{1, 1}, Coordinate{1, 2}, Coordinate{1, 3},
	)
	for _, x := range []int{2, 3} {
		for y := 0; y <= 4; y++ {
			add(SideOne, Coordinate{x, y})
		}
	}

	// SideTwo starts on the right side.
	for _, x := range []int{5, 6} {
		for y := 0; y <= 4; y++ {
			add(SideTwo, Coordinate{x, y})
		}
	}
	add(SideTwo,
		Coordinate{7, 1}, Coordinate{7, 2}, Coordinate{7, 3},
		Coordinate{8, 1}, Coordinate{8, 2}, Coordinate{8, 3},
	)

	return placements
}

func IndexPoints(points []Position) map[Coordinate]PointID {
	index := make(map[Coordinate]PointID, len(points))
	for _, point := range points {
		index[Coordinate{point.X, point.Y}] = point.ID
	}
	return index
}

func LessCoordinate(a, b Coordinate) bool {
	return a.Y < b.Y || a.Y == b.Y && a.X < b.X
}
