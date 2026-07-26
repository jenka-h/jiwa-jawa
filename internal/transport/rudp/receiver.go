package rudp

import (
	"context"
	"net"

	"jiwa-jaw/internal/transport/protocol"
)

// ReceivedPacket is a packet accepted from the network with sender metadata.
type ReceivedPacket struct {
	Packet protocol.Packet
	Addr   *net.UDPAddr
}

// Receiver owns receive and error channels for asynchronous packet handling.
type Receiver struct{}

// NewReceiver creates a receiver with buffered channels.
func NewReceiver(bufferSize int) *Receiver {
	return nil
}

// Packets returns the accepted packet channel.
func (r *Receiver) Packets() <-chan ReceivedPacket {
	return nil
}

// Errors returns the asynchronous error channel.
func (r *Receiver) Errors() <-chan error {
	return nil
}

// Run starts the receive loop.
func (r *Receiver) Run(ctx context.Context, conn *Connection) error {
	return ErrNotImplemented
}

// Deliver sends a packet to the receive channel.
func (r *Receiver) Deliver(packet ReceivedPacket) error {
	return ErrNotImplemented
}

// Report sends an error to the error channel.
func (r *Receiver) Report(err error) {}

// Close closes receiver channels.
func (r *Receiver) Close() error {
	return ErrNotImplemented
}
