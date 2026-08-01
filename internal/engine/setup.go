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
		from := MustPointIDByPosition(edge.from.x, edge.from.y)
		to := MustPointIDByPosition(edge.to.x, edge.to.y)

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
		id := MustPointIDByPosition(placement.x, placement.y)

		if err := board.PlacePiece(
			id,
			Piece{Owner: placement.owner},
		); err != nil {
			return nil, err
		}
	}

	return board, nil
}

func PointIDByPosition(x, y int) (PointID, error) {
	id, ok := pointByPosition[coordinate{x, y}]
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
	coords := []coordinate{
		{2, -2},
		{1, -1},
		{2, -1},
		{3, -1},
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			coords = append(coords, coordinate{x, y})
		}
	}

	coords = append(
		coords,
		coordinate{1, 5},
		coordinate{2, 5},
		coordinate{3, 5},
		coordinate{2, 6},
	)

	points := make([]Position, len(coords))
	for i, coord := range coords {
		points[i] = Position{
			ID: PointID(i + 1),
			X:  coord.x,
			Y:  coord.y,
		}
	}

	return points
}

func BuildDefaultEdges() []edgeSpec {
	points := make(map[coordinate]struct{}, len(defaultPoints))
	for _, point := range defaultPoints {
		points[coordinate{point.X, point.Y}] = struct{}{}
	}

	edges := make(map[edgeSpec]struct{})

	add := func(a, b coordinate) {
		if _, ok := points[a]; !ok {
			return
		}

		if _, ok := points[b]; !ok {
			return
		}

		if LessCoordinate(b, a) {
			a, b = b, a
		}

		edges[edgeSpec{a, b}] = struct{}{}
	}

	for point := range points {
		add(point, coordinate{point.x + 1, point.y})
		add(point, coordinate{point.x, point.y + 1})
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if (x+y)%2 != 0 {
				continue
			}

			add(
				coordinate{x, y},
				coordinate{x + 1, y + 1},
			)

			add(
				coordinate{x + 1, y},
				coordinate{x, y + 1},
			)
		}
	}

	addTip := func(tipY, middleY, baseY int) {
		add(coordinate{2, tipY}, coordinate{1, middleY})
		add(coordinate{2, tipY}, coordinate{2, middleY})
		add(coordinate{2, tipY}, coordinate{3, middleY})

		add(coordinate{1, middleY}, coordinate{1, baseY})
		add(coordinate{1, middleY}, coordinate{2, baseY})

		add(coordinate{2, middleY}, coordinate{1, baseY})
		add(coordinate{2, middleY}, coordinate{2, baseY})
		add(coordinate{2, middleY}, coordinate{3, baseY})

		add(coordinate{3, middleY}, coordinate{2, baseY})
		add(coordinate{3, middleY}, coordinate{3, baseY})
	}

	addTip(-2, -1, 0)
	addTip(6, 5, 4)

	result := make([]edgeSpec, 0, len(edges))
	for edge := range edges {
		result = append(result, edge)
	}

	return result
}

func BuildDefaultPlacements() []placementSpec {
	var placements []placementSpec

	add := func(owner Player, coords ...coordinate) {
		for _, coord := range coords {
			placements = append(
				placements,
				placementSpec{
					coordinate: coord,
					owner:      owner,
				},
			)
		}
	}

	addRows := func(owner Player, rows ...int) {
		for _, y := range rows {
			for x := 0; x < 5; x++ {
				add(owner, coordinate{x, y})
			}
		}
	}

	addRows(PlayerOne, 3, 4)
	add(
		PlayerOne,
		coordinate{1, 5},
		coordinate{2, 5},
		coordinate{3, 5},
		coordinate{2, 6},
		coordinate{3, 2},
		coordinate{4, 2},
	)

	addRows(PlayerTwo, 0, 1)
	add(
		PlayerTwo,
		coordinate{1, -1},
		coordinate{2, -1},
		coordinate{3, -1},
		coordinate{2, -2},
		coordinate{0, 2},
		coordinate{1, 2},
	)

	return placements
}

func IndexPoints(points []Position) map[coordinate]PointID {
	index := make(map[coordinate]PointID, len(points))

	for _, point := range points {
		index[coordinate{point.X, point.Y}] = point.ID
	}

	return index
}

func LessCoordinate(a, b coordinate) bool {
	return a.y < b.y || a.y == b.y && a.x < b.x
}
