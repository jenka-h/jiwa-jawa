package rudp

import (
	"errors"
	"time"
)

const (
	DefaultTimeout    = 500 * time.Millisecond
	DefaultMaxRetries = 5
	DefaultBufferSize = 65535
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrClosed         = errors.New("rudp connection is closed")
	ErrAckTimeout     = errors.New("ack timeout")
	ErrMaxRetries     = errors.New("max retries exceeded")
	ErrPeerNotSet     = errors.New("peer address is not set")
	ErrDuplicate      = errors.New("duplicate packet")
	ErrUnexpectedPeer = errors.New("packet received from unexpected peer")
)

// Config contains basic RUDP connection settings.
type Config struct {
	ListenAddr string
	PeerAddr   string
	SessionID  uint64

	Timeout    time.Duration
	MaxRetries int
	BufferSize int
}
