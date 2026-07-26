package app

import "jiwa-jaw/internal/engine"

// Config contains runtime options for one Dam Daman client.
type Config struct {
	PlayerID   string
	PlayerName string
	ListenAddr string
	PeerAddr   string
	SessionID  uint64
	LogPath    string
}

// JoinDTO is sent when a player wants to join/start a session.
type JoinDTO struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

// AcceptDTO is sent when a join request is accepted.
type AcceptDTO struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	SessionID  uint64 `json:"session_id"`
}

// MoveDTO is the network-safe representation of a Dam Daman move.
type MoveDTO struct {
	PlayerID string         `json:"player_id"`
	Source   engine.PointID `json:"source"`
	Target   engine.PointID `json:"target"`
}

// ErrorDTO describes an application-level error response.
type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
