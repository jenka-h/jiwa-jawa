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
		return 0, fmt.Errorf(
			"no default board point at (%d, %d)",
			x,
			y,
		)
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
		{2, -2},
		{1, -1},
		{2, -1},
		{3, -1},
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			coords = append(coords, Coordinate{x, y})
		}
	}

	coords = append(
		coords,
		Coordinate{1, 5},
		Coordinate{2, 5},
		Coordinate{3, 5},
		Coordinate{2, 6},
	)

	points := make([]Position, len(coords))
	for i, coord := range coords {
		points[i] = Position{
			ID: PointID(i + 1),
			X:  coord.X,
			Y:  coord.Y,
		}
	}

	return points
}

func BuildDefaultEdges() []EdgeSpec {
	points := make(map[Coordinate]struct{}, len(defaultPoints))
	for _, point := range defaultPoints {
		points[Coordinate{point.X, point.Y}] = struct{}{}
	}

	edges := make(map[EdgeSpec]struct{})

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

		edges[EdgeSpec{a, b}] = struct{}{}
	}

	for point := range points {
		add(point, Coordinate{point.X + 1, point.Y})
		add(point, Coordinate{point.X, point.Y + 1})
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if (x+y)%2 != 0 {
				continue
			}

			add(
				Coordinate{x, y},
				Coordinate{x + 1, y + 1},
			)

			add(
				Coordinate{x + 1, y},
				Coordinate{x, y + 1},
			)
		}
	}

	addTip := func(tipY, middleY, baseY int) {
		add(Coordinate{2, tipY}, Coordinate{1, middleY})
		add(Coordinate{2, tipY}, Coordinate{2, middleY})
		add(Coordinate{2, tipY}, Coordinate{3, middleY})

		add(Coordinate{1, middleY}, Coordinate{1, baseY})
		add(Coordinate{1, middleY}, Coordinate{2, baseY})

		add(Coordinate{2, middleY}, Coordinate{1, baseY})
		add(Coordinate{2, middleY}, Coordinate{2, baseY})
		add(Coordinate{2, middleY}, Coordinate{3, baseY})

		add(Coordinate{3, middleY}, Coordinate{2, baseY})
		add(Coordinate{3, middleY}, Coordinate{3, baseY})
	}

	addTip(-2, -1, 0)
	addTip(6, 5, 4)

	result := make([]EdgeSpec, 0, len(edges))
	for edge := range edges {
		result = append(result, edge)
	}

	return result
}

func BuildDefaultPlacements() []PlacementSpec {
	var placements []PlacementSpec

	add := func(owner Side, coords ...Coordinate) {
		for _, coord := range coords {
			placements = append(
				placements,
				PlacementSpec{
					Coordinate: coord,
					Owner:      owner,
				},
			)
		}
	}

	addRows := func(owner Side, rows ...int) {
		for _, y := range rows {
			for x := 0; x < 5; x++ {
				add(owner, Coordinate{x, y})
			}
		}
	}

	addRows(SideOne, 3, 4)
	add(
		SideOne,
		Coordinate{1, 5},
		Coordinate{2, 5},
		Coordinate{3, 5},
		Coordinate{2, 6},
		Coordinate{3, 2},
		Coordinate{4, 2},
	)

	addRows(SideTwo, 0, 1)
	add(
		SideTwo,
		Coordinate{1, -1},
		Coordinate{2, -1},
		Coordinate{3, -1},
		Coordinate{2, -2},
		Coordinate{0, 2},
		Coordinate{1, 2},
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
