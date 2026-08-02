package err

import "errors"

var (
	ErrPacketTooShort     = errors.New("packet too short")
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrInvalidACK         = errors.New("invalid ACK packet")

	ErrClosed     = errors.New("rudp connection is closed")
	ErrMaxRetries = errors.New("max retries exceeded")
	ErrPeerNotSet = errors.New("peer address is not set")

	ErrInvalidPlayerID   = errors.New("invalid player ID")
	ErrInvalidPlayerName = errors.New("invalid player name")
)
