package protocol

import (
	"encoding/json"
	"fmt"
	"time"
)

type MessageType uint8

const (
	MessageJoin   MessageType = iota + 1 // sebagai SYN
	MessageAccept                        // sebagai SYN-ACK
	MessageMove                          // for every move

	MessageACK

	MessageFinish
	MessageLeave

	// state recovery + disconnect
	MessageHeartbeat
	MessageStateRequest
	MessageStateSnapshot
)

func (t MessageType) IsValid() bool {
	return t >= MessageJoin && t <= MessageStateSnapshot
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
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

type Envelope struct {
	MessageID   string          `json:"message_id"`
	SessionID   string          `json:"session_id"`
	SenderID    string          `json:"sender_id"`
	Sequence    uint64          `json:"sequence"`
	MessageType MessageType     `json:"message_type"`
	AckFor      string          `json:"ack_for,omitempty"`
	Timestamp   int64           `json:"timestamp"`
	Payload     json.RawMessage `json:"payload,omitempty"`
}

func NewEnvelope(
	messageID string,
	sessionID string,
	senderID string,
	sequence uint64,
	messageType MessageType,
	payload json.RawMessage,
	ackFor string,
) Envelope {
	return Envelope{
		MessageID:   messageID,
		SessionID:   sessionID,
		SenderID:    senderID,
		Sequence:    sequence,
		MessageType: messageType,
		AckFor:      ackFor,
		Timestamp:   time.Now().UnixMilli(),
		Payload:     payload,
	}
}
