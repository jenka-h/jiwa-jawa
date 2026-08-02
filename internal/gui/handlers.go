package gui

import (
	"encoding/json"
	"fmt"
	"net/http"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/transport/protocol"
)

func (s *Server) currentMatch() *Match {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.match
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	match := s.currentMatch()
	if match == nil {
		writeJSON(w, http.StatusOK, struct {
			Active bool `json:"active"`
		}{Active: false})
		return
	}
	writeJSON(w, http.StatusOK, s.snapshot(match))
}

func (s *Server) handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		Source int `json:"source"`
		Target int `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid move request", http.StatusBadRequest)
		return
	}
	match := s.currentMatch()
	if match == nil {
		http.Error(w, "no active session", http.StatusBadRequest)
		return
	}
	match.mu.Lock()
	if match.status != "ready" {
		match.mu.Unlock()
		http.Error(w, "match is not ready yet", http.StatusBadRequest)
		return
	}
	if match.state.CurrentTurn != match.localSide {
		match.mu.Unlock()
		http.Error(w, "not your turn", http.StatusBadRequest)
		return
	}
	move := engine.Move{Source: engine.PointID(request.Source), Target: engine.PointID(request.Target)}
	result, err := match.state.ApplyMoveResult(move)
	localName, localID := match.localName, match.localID
	match.mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload := protocol.MovePayload{PlayerID: localID, Source: request.Source, Target: request.Target}
	if err := s.sendJSON(match, protocol.MessageMove, payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	s.logMove(match, "Move", localName, move, result)
	if result.GameOver {
		match.mu.Lock()
		match.status = "finished"
		match.mu.Unlock()
		s.applyRating(match)
	}
	writeJSON(w, http.StatusOK, s.snapshot(match))
}

func (s *Server) handlePenalty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		Target int `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid penalty request", http.StatusBadRequest)
		return
	}
	match := s.currentMatch()
	if match == nil {
		http.Error(w, "no active session", http.StatusBadRequest)
		return
	}
	match.mu.Lock()
	if match.status != "ready" {
		match.mu.Unlock()
		http.Error(w, "match is not ready yet", http.StatusBadRequest)
		return
	}
	err := match.state.ApplyPenalty(match.localSide, engine.PointID(request.Target))
	localID, localName := match.localID, match.localName
	finished := match.state.Finished
	match.mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload := protocol.PenaltyPayload{PlayerID: localID, Target: request.Target}
	if err := s.sendJSON(match, protocol.MessagePenalty, payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	s.logPenalty(match, localName, engine.PointID(request.Target))
	if finished {
		match.mu.Lock()
		match.status = "finished"
		match.mu.Unlock()
		s.applyRating(match)
	}
	writeJSON(w, http.StatusOK, s.snapshot(match))
}

func (s *Server) handleEndTurn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	match := s.currentMatch()
	if match == nil {
		http.Error(w, "no active session", http.StatusBadRequest)
		return
	}
	match.mu.Lock()
	err := match.state.EndTurn(match.localSide)
	localID, localName := match.localID, match.localName
	match.mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.sendJSON(match, protocol.MessageEndTurn, protocol.EndTurnPayload{PlayerID: localID}); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	s.logEndTurn(match, localName)
	writeJSON(w, http.StatusOK, s.snapshot(match))
}

func (s *Server) handleFinish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	match := s.currentMatch()
	if match == nil {
		http.Error(w, "no active session", http.StatusBadRequest)
		return
	}
	match.mu.Lock()
	winner := engine.Opponent(match.localSide)
	match.state.Finish(winner)
	match.status = "finished"
	match.mu.Unlock()
	s.applyRating(match)
	_ = s.sendJSON(match, protocol.MessageFinish, protocol.FinishPayload{WinnerID: string(winner), Reason: "surrender"})
	s.appendLog(match, "System", "Surrender", fmt.Sprintf("winner: %s", winner))
	writeJSON(w, http.StatusOK, s.snapshot(match))
}

func (s *Server) handleBoard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	board, err := engine.NewDefaultBoard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, boardSnapshot(board))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
