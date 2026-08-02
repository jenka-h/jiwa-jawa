package gui

import (
	"fmt"
	"sort"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/rating"
)

func demoPageData() PageData {
	board, err := engine.NewDefaultBoard()
	if err != nil {
		panic(err)
	}
	elo := rating.NewElo(32)
	p1, p2 := rating.PlayerRating{Rating: 1248}, rating.PlayerRating{Rating: 1312}
	p1After, p2After, _ := elo.Update(p1, p2, rating.Win)
	return PageData{
		Players: []PlayerView{
			{Name: "P1 (Dark)", Alias: "Host", Side: engine.SideOne, Piece: "/assets/p1.png", Rating: p1.Rating},
			{Name: "P2 (Gold)", Alias: "Guest", Side: engine.SideTwo, Piece: "/assets/p2.png", Rating: p2.Rating},
		},
		Board:    buildBoardView(board),
		Logs:     []LogItem{{No: 1, Player: "System", Action: "Ready", Move: "Create or join a game"}},
		Duration: "00:00:00", TurnText: "Player 1's Turn",
		Elo: EloView{PlayerOneBefore: p1.Rating, PlayerTwoBefore: p2.Rating, PlayerOneAfter: p1After.Rating, PlayerTwoAfter: p2After.Rating},
	}
}

func buildBoardView(board *engine.Board) BoardView {
	const width, height, padX, padY = 900, 520, 70, 60
	coords := map[engine.PointID]PointView{}
	ids := make([]int, 0, len(board.Points))
	for id := range board.Points {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	for _, rawID := range ids {
		id := engine.PointID(rawID)
		position := board.Points[id]
		x, y := project(position, width, height, padX, padY)
		pieceName, side := "", engine.SideNone
		if piece, ok := board.PieceAt(id); ok {
			side = piece.Owner
			if side == engine.SideOne {
				pieceName = "/assets/p1.png"
			} else if side == engine.SideTwo {
				pieceName = "/assets/p2.png"
			}
		}
		coords[id] = PointView{ID: id, X: x, Y: y, Piece: pieceName, Side: side}
	}
	points := make([]PointView, 0, len(coords))
	for _, rawID := range ids {
		points = append(points, coords[engine.PointID(rawID)])
	}
	seen := map[string]struct{}{}
	edges := []EdgeView{}
	for from, neighbors := range board.Neighbors {
		for _, to := range neighbors {
			key := edgeKey(from, to)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			a, b := coords[from], coords[to]
			edges = append(edges, EdgeView{X1: a.X, Y1: a.Y, X2: b.X, Y2: b.Y})
		}
	}
	return BoardView{Width: width, Height: height, Points: points, Edges: edges}
}

func project(position engine.Position, width, height, padX, padY int) (int, int) {
	normalizedX := float64(position.X) / 8
	normalizedY := float64(position.Y) / 4
	return padX + int(normalizedX*float64(width-2*padX)), padY + int(normalizedY*float64(height-2*padY))
}

func edgeKey(a, b engine.PointID) string {
	if a > b {
		a, b = b, a
	}
	return fmt.Sprintf("%d:%d", a, b)
}
