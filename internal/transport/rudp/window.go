package rudp

// FlowWindow controls how many reliable packets may be in flight.
// Start with size 1 for stop-and-wait behavior.
type FlowWindow struct{}

// NewFlowWindow creates a flow-control window.
func NewFlowWindow(size int) *FlowWindow {
	return nil
}

// Size returns the configured window size.
func (w *FlowWindow) Size() int {
	return 0
}

// Available reports whether another packet can be sent.
func (w *FlowWindow) Available() bool {
	return false
}

// Acquire reserves one window slot.
func (w *FlowWindow) Acquire() error {
	return ErrNotImplemented
}

// Release frees one window slot.
func (w *FlowWindow) Release() error {
	return ErrNotImplemented
}

// Reset clears all in-flight state.
func (w *FlowWindow) Reset() {}
