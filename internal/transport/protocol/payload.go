package protocol

type JoinPayload struct {
	PlayerID     string  `json:"player_id"`
	PlayerName   string  `json:"player_name"`
	PlayerRating float64 `json:"player_rating,omitempty"`
}

type AcceptPayload struct {
	PlayerID     string  `json:"player_id"`
	PlayerName   string  `json:"player_name"`
	PlayerRating float64 `json:"player_rating,omitempty"`
	SessionID    uint64  `json:"session_id"`
}

type MovePayload struct {
	PlayerID string `json:"player_id"`
	Source   int    `json:"source"`
	Target   int    `json:"target"`
}

type EndTurnPayload struct {
	PlayerID string `json:"player_id"`
}

type PenaltyPayload struct {
	PlayerID string `json:"player_id"`
	Target   int    `json:"target"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FinishPayload struct {
	WinnerID string `json:"winner_id,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type StateSnapshotPayload struct {
	StateHash string `json:"state_hash,omitempty"`
	Data      []byte `json:"data,omitempty"`
}
