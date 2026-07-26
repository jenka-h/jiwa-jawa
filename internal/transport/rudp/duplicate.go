package rudp

// DuplicateTracker records received packet sequences so duplicates can be suppressed.
type DuplicateTracker struct{}

// NewDuplicateTracker creates a duplicate tracker.
func NewDuplicateTracker() *DuplicateTracker {
	return nil
}

// Seen reports whether sequence has already been received.
func (t *DuplicateTracker) Seen(sequence uint32) bool {
	return false
}

// Mark records sequence as received.
func (t *DuplicateTracker) Mark(sequence uint32) error {
	return ErrNotImplemented
}

// Forget removes sequence from the tracker.
func (t *DuplicateTracker) Forget(sequence uint32) error {
	return ErrNotImplemented
}

// Reset clears all duplicate tracking state.
func (t *DuplicateTracker) Reset() {}
