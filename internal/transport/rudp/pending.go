package rudp

import (
	"time"

	"jiwa-jaw/internal/transport/protocol"
)

// PendingPacket stores metadata for a packet waiting for ACK.
type PendingPacket struct {
	Packet  protocol.Packet
	SentAt  time.Time
	Retries int
	Acked   bool
}

// PendingStore tracks packets that have been sent but not acknowledged.
type PendingStore struct{}

// NewPendingStore creates an empty pending packet store.
func NewPendingStore() *PendingStore {
	return nil
}

// Add stores a packet as pending.
func (s *PendingStore) Add(packet protocol.Packet) error {
	return ErrNotImplemented
}

// Ack marks a packet sequence as acknowledged.
func (s *PendingStore) Ack(sequence uint32) error {
	return ErrNotImplemented
}

// Get returns the pending packet for sequence.
func (s *PendingStore) Get(sequence uint32) (*PendingPacket, bool) {
	return nil, false
}

// Remove deletes a pending packet by sequence.
func (s *PendingStore) Remove(sequence uint32) error {
	return ErrNotImplemented
}

// IncrementRetry increments and returns the retry count for sequence.
func (s *PendingStore) IncrementRetry(sequence uint32) (int, error) {
	return 0, ErrNotImplemented
}

// Expired returns packets that should be retried.
func (s *PendingStore) Expired(timeout time.Duration) []PendingPacket {
	return nil
}

// Len returns the number of pending packets.
func (s *PendingStore) Len() int {
	return 0
}

// Clear removes all pending packets.
func (s *PendingStore) Clear() {}
