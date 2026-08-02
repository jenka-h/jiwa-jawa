package protocol

import "fmt"

type MessageType uint8

const (
	MessageJoin MessageType = iota + 1
	MessageAccept
	MessageMove
	MessageACK
	MessageFinish
	MessageLeave
	MessageHeartbeat
	MessageStateRequest
	MessageStateSnapshot
	MessagePenalty
	MessageEndTurn
)

func (t MessageType) IsValid() bool {
	return t >= MessageJoin && t <= MessageEndTurn
}

func (t MessageType) String() string {
	switch t {
	case MessageJoin:
		return "JOIN"
	case MessageAccept:
		return "ACCEPT"
	case MessageMove:
		return "MOVE"
	case MessageACK:
		return "ACK"
	case MessageFinish:
		return "FINISH"
	case MessageLeave:
		return "LEAVE"
	case MessageHeartbeat:
		return "HEARTBEAT"
	case MessageStateRequest:
		return "STATE_REQUEST"
	case MessageStateSnapshot:
		return "STATE_SNAPSHOT"
	case MessagePenalty:
		return "PENALTY"
	case MessageEndTurn:
		return "END_TURN"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}
