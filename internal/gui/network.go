package gui

import (
	"encoding/json"
	"fmt"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/gamelog"
	"jiwa-jawa/internal/transport/protocol"
	"jiwa-jawa/internal/transport/rudp"
)

func (s *Server) networkLoop(match *Match) {
	for {
		select {
		case <-match.ctx.Done():
			return
		case err := <-match.conn.Errors():
			s.appendLog(match, "Network", "Error", err.Error())
			_ = match.logger.Log(gamelog.Event{Type: "transport_error", SessionID: match.sessionID, Error: err.Error()})
		case received := <-match.conn.Incoming():
			s.handlePacket(match, received)
		}
	}
}

func (s *Server) handlePacket(match *Match, received rudp.ReceivedPacket) {
	s.mu.Lock()
	current := s.match == match
	s.mu.Unlock()
	if !current {
		return
	}

	switch received.Packet.Header.MessageType {
	case protocol.MessageJoin:
		var payload protocol.JoinPayload
		if err := json.Unmarshal(received.Packet.Payload, &payload); err != nil {
			s.appendLog(match, "Network", "Error", err.Error())
			return
		}
		_ = match.conn.SetPeer(received.Addr.String())
		match.mu.Lock()
		match.remoteID, match.remoteName = payload.PlayerID, payload.PlayerName
		match.remoteRating, match.status = normalizedRating(payload.PlayerRating), "ready"
		match.mu.Unlock()
		s.appendLog(match, "System", "Join", fmt.Sprintf("%s connected", payload.PlayerName))
		accept := protocol.AcceptPayload{PlayerID: match.localID, PlayerName: match.localName, PlayerRating: match.localRating, SessionID: match.sessionID}
		if err := s.sendJSON(match, protocol.MessageAccept, accept); err != nil {
			s.appendLog(match, "Network", "Error", err.Error())
		}
	case protocol.MessageAccept:
		var payload protocol.AcceptPayload
		if err := json.Unmarshal(received.Packet.Payload, &payload); err != nil {
			s.appendLog(match, "Network", "Error", err.Error())
			return
		}
		match.mu.Lock()
		match.remoteID, match.remoteName = payload.PlayerID, payload.PlayerName
		match.remoteRating, match.sessionID, match.status = normalizedRating(payload.PlayerRating), payload.SessionID, "ready"
		match.mu.Unlock()
		match.conn.SetSessionID(payload.SessionID)
		s.appendLog(match, "System", "Accept", fmt.Sprintf("connected to %s", payload.PlayerName))
	case protocol.MessageMove:
		s.handlePeerMove(match, received.Packet.Payload)
	case protocol.MessagePenalty:
		s.handlePeerPenalty(match, received.Packet.Payload)
	case protocol.MessageEndTurn:
		s.handlePeerEndTurn(match, received.Packet.Payload)
	case protocol.MessageLeave:
		match.mu.Lock()
		match.status = "finished"
		match.state.Finish(match.localSide)
		match.mu.Unlock()
		s.applyRating(match)
		s.appendLog(match, "System", "Leave", "opponent left")
	case protocol.MessageFinish:
		var payload protocol.FinishPayload
		_ = json.Unmarshal(received.Packet.Payload, &payload)
		winner := engine.Side(payload.WinnerID)
		if winner == engine.SideNone {
			winner = match.remoteSide
		}
		match.mu.Lock()
		match.status = "finished"
		match.state.Finish(winner)
		match.mu.Unlock()
		s.applyRating(match)
		s.appendLog(match, "System", "Finish", payload.Reason)
	}
}

func (s *Server) handlePeerMove(match *Match, data []byte) {
	var payload protocol.MovePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		s.appendLog(match, "Network", "Error", err.Error())
		return
	}
	move := engine.Move{Source: engine.PointID(payload.Source), Target: engine.PointID(payload.Target)}
	match.mu.Lock()
	result, err := match.state.ApplyMoveResult(move)
	remoteName := match.remoteName
	match.mu.Unlock()
	if err != nil {
		s.appendLog(match, "Move", "Rejected", err.Error())
		return
	}
	s.logMove(match, "Move", remoteName, move, result)
	if result.GameOver {
		match.mu.Lock()
		match.status = "finished"
		match.mu.Unlock()
		s.applyRating(match)
	}
}

func (s *Server) handlePeerEndTurn(match *Match, data []byte) {
	var payload protocol.EndTurnPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		s.appendLog(match, "Network", "Error", err.Error())
		return
	}
	match.mu.Lock()
	err := match.state.EndTurn(match.remoteSide)
	remoteName := match.remoteName
	match.mu.Unlock()
	if err != nil {
		s.appendLog(match, "End Turn", "Rejected", err.Error())
		return
	}
	s.logEndTurn(match, remoteName)
}

func (s *Server) handlePeerPenalty(match *Match, data []byte) {
	var payload protocol.PenaltyPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		s.appendLog(match, "Network", "Error", err.Error())
		return
	}
	match.mu.Lock()
	err := match.state.ApplyPenalty(match.remoteSide, engine.PointID(payload.Target))
	remoteName := match.remoteName
	finished := match.state.Finished
	match.mu.Unlock()
	if err != nil {
		s.appendLog(match, "Dam Ora Mangan", "Rejected", err.Error())
		return
	}
	s.logPenalty(match, remoteName, engine.PointID(payload.Target))
	if finished {
		match.mu.Lock()
		match.status = "finished"
		match.mu.Unlock()
		s.applyRating(match)
	}
}

func (s *Server) sendJSON(match *Match, messageType protocol.MessageType, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	packet := protocol.NewPacket(match.sessionID, 0, messageType, data, true)
	if err := match.conn.SendReliable(packet); err != nil {
		_ = match.logger.Log(gamelog.Event{Type: "send_error", SessionID: match.sessionID, Message: messageType.String(), Error: err.Error()})
		return err
	}
	_ = match.logger.Log(gamelog.Event{Type: "send", SessionID: match.sessionID, Message: messageType.String()})
	return nil
}
