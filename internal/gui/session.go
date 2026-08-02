package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/gamelog"
	"jiwa-jawa/internal/transport/protocol"
	"jiwa-jawa/internal/transport/rudp"
)

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		ProfileID string `json:"profile_id"`
		Role      string `json:"role"`
		Peer      string `json:"peer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid session request", http.StatusBadRequest)
		return
	}
	profile, ok := s.profileByID(request.ProfileID)
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}
	role := strings.ToLower(strings.TrimSpace(request.Role))
	if role == "" {
		role = "host"
	}
	if role != "host" && role != "join" {
		http.Error(w, "role must be host or join", http.StatusBadRequest)
		return
	}
	peer := strings.TrimSpace(request.Peer)
	if peer == "" {
		peer = s.options.PeerAddr
	}
	match, err := s.newMatch(role, profile, peer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	if s.match != nil {
		s.match.cancel()
		_ = s.match.conn.Close()
		_ = s.match.logger.Close()
	}
	s.match = match
	s.mu.Unlock()

	go s.networkLoop(match)
	if role == "join" {
		join := protocol.JoinPayload{PlayerID: profile.ID, PlayerName: profile.Name, PlayerRating: profile.Rating}
		if err := s.sendJSON(match, protocol.MessageJoin, join); err != nil {
			s.appendLog(match, "System", "Error", err.Error())
		}
	}
	writeJSON(w, http.StatusOK, s.snapshot(match))
}

func (s *Server) newMatch(role string, profile Profile, peer string) (*Match, error) {
	board, err := engine.NewDefaultBoard()
	if err != nil {
		return nil, err
	}
	logger, err := gamelog.NewWithRemote(s.options.LogPath, s.options.LogServer)
	if err != nil {
		return nil, err
	}
	peerAddr := ""
	if role == "join" {
		peerAddr = peer
	}
	conn, err := rudp.NewConnection(rudp.Config{ListenAddr: s.options.ListenAddr, PeerAddr: peerAddr, SessionID: s.options.SessionID})
	if err != nil {
		_ = logger.Close()
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err := conn.Start(ctx); err != nil {
		cancel()
		_ = conn.Close()
		_ = logger.Close()
		return nil, err
	}
	localSide, remoteSide, status := engine.SideOne, engine.SideTwo, "waiting"
	if role == "join" {
		localSide, remoteSide, status = engine.SideTwo, engine.SideOne, "joining"
	}
	match := &Match{
		ctx: ctx, cancel: cancel, conn: conn, logger: logger,
		state: engine.NewGameState(board, engine.SideOne),
		role:  role, status: status, localID: profile.ID, localName: profile.Name,
		localRating: profile.Rating, localSide: localSide, remoteSide: remoteSide,
		sessionID: s.options.SessionID, startedAt: time.Now(),
	}
	message := fmt.Sprintf("%s listening on %s", role, conn.LocalAddr())
	s.appendLog(match, "System", "Start", message)
	_ = logger.Log(gamelog.Event{Type: "start", SessionID: match.sessionID, Message: message})
	return match, nil
}
