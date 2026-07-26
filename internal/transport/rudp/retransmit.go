package rudp

import (
	"context"
	"time"
)

// Retransmitter scans pending packets and retries expired packets.
type Retransmitter struct{}

// NewRetransmitter creates a retransmission worker.
func NewRetransmitter(timeout time.Duration, maxRetries int) *Retransmitter {
	return nil
}

// Run starts the retransmission loop.
func (r *Retransmitter) Run(ctx context.Context, conn *Connection, pending *PendingStore) error {
	return ErrNotImplemented
}

// RetryExpired retries packets that have exceeded timeout.
func (r *Retransmitter) RetryExpired(conn *Connection, pending *PendingStore) error {
	return ErrNotImplemented
}
