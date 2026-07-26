package rudp

// SequenceGenerator produces packet sequence numbers.
type SequenceGenerator struct{}

// NewSequenceGenerator creates a sequence generator starting at initial.
func NewSequenceGenerator(initial uint32) *SequenceGenerator {
	return nil
}

// Next returns the next sequence number.
func (g *SequenceGenerator) Next() uint32 {
	return 0
}

// Current returns the current sequence number without advancing it.
func (g *SequenceGenerator) Current() uint32 {
	return 0
}

// Reset sets the next sequence number.
func (g *SequenceGenerator) Reset(next uint32) {}

// SequenceTracker tracks expected inbound sequence numbers.
type SequenceTracker struct{}

// NewSequenceTracker creates an inbound sequence tracker.
func NewSequenceTracker(expected uint32) *SequenceTracker {
	return nil
}

// Expected returns the next expected inbound sequence.
func (t *SequenceTracker) Expected() uint32 {
	return 0
}

// Accept checks whether sequence is acceptable and advances tracker state if needed.
func (t *SequenceTracker) Accept(sequence uint32) error {
	return ErrNotImplemented
}

// Reset sets the next expected inbound sequence.
func (t *SequenceTracker) Reset(expected uint32) {}
