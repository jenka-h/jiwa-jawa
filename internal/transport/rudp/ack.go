package rudp

import (
	"net"
	"time"

	"jiwa-jaw/internal/transport/protocol"
)

// SendACK sends an ACK for sequence to the configured peer.
func (c *Connection) SendACK(sequence uint32) error {
	return ErrNotImplemented
}

// SendACKTo sends an ACK for sequence to a specific address.
func (c *Connection) SendACKTo(sequence uint32, addr *net.UDPAddr) error {
	return ErrNotImplemented
}

// HandleACK processes an incoming ACK packet.
func (c *Connection) HandleACK(packet protocol.Packet) error {
	return ErrNotImplemented
}

// WaitForACK waits until sequence is acknowledged or timeout expires.
func (c *Connection) WaitForACK(sequence uint32, timeout time.Duration) error {
	return ErrNotImplemented
}
