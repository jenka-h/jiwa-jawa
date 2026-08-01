package protocol

// JoinPayload carries player metadata for MessageJoin.
type JoinPayload struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

// AcceptPayload carries session metadata for MessageAccept.
type AcceptPayload struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	SessionID  uint64 `json:"session_id"`
}

// MovePayload carries one Dam Daman move.
type MovePayload struct {
	PlayerID string `json:"player_id"`
	Source   int    `json:"source"`
	Target   int    `json:"target"`
}

// ErrorPayload carries an error code and user-facing message.
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// FinishPayload carries game finish information.
type FinishPayload struct {
	WinnerID string `json:"winner_id,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// StateSnapshotPayload carries serialized recovery state.
type StateSnapshotPayload struct {
	StateHash string `json:"state_hash,omitempty"`
	Data      []byte `json:"data,omitempty"`
}
