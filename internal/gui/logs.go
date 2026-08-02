package gui

import (
	"fmt"
	"strings"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/gamelog"
)

func (s *Server) logMove(match *Match, action, player string, move engine.Move, result engine.MoveResult) {
	message := move.String()
	if len(result.Captured) > 0 {
		message = fmt.Sprintf("%s captures %v", move, result.Captured)
	}
	if result.MustContinue {
		message += "; continue"
	}
	if result.DamOraMangan {
		message += "; Dam Ora Mangan — opponent may remove three pieces"
	}
	s.appendLog(match, player, action, message)
	_ = match.logger.Log(gamelog.Event{Type: strings.ToLower(action), SessionID: match.sessionID, Move: &gamelog.MoveEvent{Source: move.Source, Target: move.Target}, Message: message})
}

func (s *Server) logEndTurn(match *Match, player string) {
	s.appendLog(match, player, "End Turn", "declined optional chained capture")
	match.mu.Lock()
	sessionID := match.sessionID
	match.mu.Unlock()
	_ = match.logger.Log(gamelog.Event{Type: "end_turn", SessionID: sessionID, Message: "declined optional chained capture"})
}

func (s *Server) logPenalty(match *Match, player string, target engine.PointID) {
	message := fmt.Sprintf("removed piece at %d", target)
	s.appendLog(match, player, "Dam Ora Mangan", message)
	match.mu.Lock()
	side, sessionID := match.state.CurrentTurn, match.sessionID
	match.mu.Unlock()
	_ = match.logger.Log(gamelog.Event{Type: "dam_ora_mangan", SessionID: sessionID, Side: side, Message: message})
}

func (s *Server) appendLog(match *Match, player, action, move string) {
	match.mu.Lock()
	defer match.mu.Unlock()
	match.logs = append(match.logs, LogItem{No: len(match.logs) + 1, Player: player, Action: action, Move: move})
	if len(match.logs) > 80 {
		match.logs = match.logs[len(match.logs)-80:]
		for i := range match.logs {
			match.logs[i].No = i + 1
		}
	}
}
