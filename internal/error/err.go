package err

import (
	"errors"
	"fmt"
)

func PrintError(msg string) {
	fmt.Printf("[Error] %s\n", msg)
}

func PrintUsageError(msg string) {
	fmt.Printf("[Error] %s\n", msg)
}

func PrintErrorf(format string, args ...interface{}) {
	fmt.Printf("[Error] "+format+"\n", args...)
}

/*
pakai method errors.Is untuk pengecekan spesifik nantinya
*/
var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrPacketTooShort     = errors.New("packet too short")
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrInvalidACK         = errors.New("invalid ACK packet")
	ErrNotImplemented     = errors.New("not implemented")

	ErrClosed         = errors.New("rudp connection is closed")
	ErrAckTimeout     = errors.New("ack timeout")
	ErrMaxRetries     = errors.New("max retries exceeded")
	ErrPeerNotSet     = errors.New("peer address is not set")
	ErrDuplicate      = errors.New("duplicate packet")
	ErrUnexpectedPeer = errors.New("packet received from unexpected peer")
)
