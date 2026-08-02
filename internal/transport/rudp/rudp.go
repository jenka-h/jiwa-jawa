package rudp

import "time"

const (
	DefaultTimeout    = 500 * time.Millisecond
	DefaultMaxRetries = 5
	DefaultBufferSize = 65535
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
