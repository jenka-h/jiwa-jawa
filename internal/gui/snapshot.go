package gui

import (
	"sort"
	"time"

	"jiwa-jawa/internal/engine"
)

type matchSnapshot struct {
	Active             bool             `json:"active"`
	Role               string           `json:"role"`
	Status             string           `json:"status"`
	SessionID          uint64           `json:"session_id"`
	LocalSide          engine.Side      `json:"local_side"`
	RemoteSide         engine.Side      `json:"remote_side"`
	LocalName          string           `json:"local_name"`
	RemoteName         string           `json:"remote_name"`
	LocalRating        float64          `json:"local_rating"`
	RemoteRating       float64          `json:"remote_rating"`
	Turn               engine.Side      `json:"turn"`
	MoveNumber         int              `json:"move_number"`
	CanEndTurn         bool             `json:"can_end_turn"`
	CaptureSources     []engine.PointID `json:"capture_sources,omitempty"`
	PenaltyActive      bool             `json:"penalty_active"`
	PenaltySide        engine.Side      `json:"penalty_side"`
	PenaltiesRemaining int              `json:"penalties_remaining"`
	Finished           bool             `json:"finished"`
	Winner             engine.Side      `json:"winner"`
	Duration           int64            `json:"duration_seconds"`
	Board              any              `json:"board"`
	Logs               []LogItem        `json:"logs"`
}

func (s *Server) snapshot(match *Match) matchSnapshot {
	match.mu.Lock()
	defer match.mu.Unlock()
	return matchSnapshot{
		Active: true, Role: match.role, Status: match.status, SessionID: match.sessionID,
		LocalSide: match.localSide, RemoteSide: match.remoteSide,
		LocalName: match.localName, RemoteName: match.remoteName,
		LocalRating: match.localRating, RemoteRating: match.remoteRating,
		Turn: match.state.CurrentTurn, MoveNumber: match.state.MoveNumber,
		CanEndTurn:     match.state.CaptureContinuation,
		CaptureSources: engine.CaptureSources(match.state.Board, match.state.CurrentTurn),
		PenaltyActive:  match.state.Phase == engine.PhasePenalty,
		PenaltySide:    match.state.PenaltySide, PenaltiesRemaining: match.state.PenaltiesRemaining,
		Finished: match.state.Finished, Winner: match.state.Winner,
		Duration: int64(time.Since(match.startedAt).Seconds()),
		Board:    boardSnapshot(match.state.Board), Logs: append([]LogItem(nil), match.logs...),
	}
}

func boardSnapshot(board *engine.Board) any {
	type point struct {
		ID    engine.PointID `json:"id"`
		X     int            `json:"x"`
		Y     int            `json:"y"`
		Owner engine.Side    `json:"owner,omitempty"`
	}
	type edge struct {
		From engine.PointID `json:"from"`
		To   engine.PointID `json:"to"`
	}
	ids := make([]int, 0, len(board.Points))
	for id := range board.Points {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	points := make([]point, 0, len(ids))
	for _, rawID := range ids {
		id := engine.PointID(rawID)
		position := board.Points[id]
		entry := point{ID: id, X: position.X, Y: position.Y}
		if piece, ok := board.PieceAt(id); ok {
			entry.Owner = piece.Owner
		}
		points = append(points, entry)
	}
	seen := map[string]struct{}{}
	edges := []edge{}
	for from, neighbors := range board.Neighbors {
		for _, to := range neighbors {
			key := edgeKey(from, to)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			edges = append(edges, edge{From: from, To: to})
		}
	}
	return struct {
		PointCount int     `json:"point_count"`
		PieceCount int     `json:"piece_count"`
		EdgeCount  int     `json:"edge_count"`
		Points     []point `json:"points"`
		Edges      []edge  `json:"edges"`
	}{len(board.Points), len(board.Pieces), len(edges), points, edges}
}
