package rudp

import (
	"context"
	"net"

	"jiwa-jaw/internal/transport/protocol"
)

// Connection represents one UDP-based reliable transport endpoint.
type Connection struct {
	conn *net.UDPConn
	peer *net.UDPAddr

	sessionID uint64

	seq      *SequenceGenerator
	receiver *Receiver
	pending  *PendingStore
	window   *FlowWindow
	tracker  *SequenceTracker

	done chan struct{}
}

// NewConnection creates a new RUDP connection bound to cfg.ListenAddr.
func NewConnection(cfg Config) (*Connection, error) {
	return nil, ErrNotImplemented
}

// Start starts background workers, such as receive loop and retransmission loop.
func (c *Connection) Start(ctx context.Context) error {
	return ErrNotImplemented
}

// Shutdown stops background workers and closes the underlying socket.
func (c *Connection) Shutdown(ctx context.Context) error {
	return ErrNotImplemented
}

// Close closes the underlying connection immediately.
func (c *Connection) Close() error {
	return ErrNotImplemented
}

// Done returns a channel that is closed when the connection shuts down.
func (c *Connection) Done() <-chan struct{} {
	return nil
}

// LocalAddr returns the local UDP address.
func (c *Connection) LocalAddr() net.Addr {
	return nil
}

// PeerAddr returns the configured peer UDP address.
func (c *Connection) PeerAddr() *net.UDPAddr {
	return nil
}

// SetPeer sets the remote peer address.
func (c *Connection) SetPeer(addr string) error {
	return ErrNotImplemented
}

// SendPacket sends one packet without reliability handling.
func (c *Connection) SendPacket(packet protocol.Packet) error {
	return ErrNotImplemented
}

// SendPacketTo sends one packet to a specific address.
func (c *Connection) SendPacketTo(packet protocol.Packet, addr *net.UDPAddr) error {
	return ErrNotImplemented
}

// ReceivePacket receives and decodes one packet synchronously.
func (c *Connection) ReceivePacket() (protocol.Packet, *net.UDPAddr, error) {
	return protocol.Packet{}, nil, ErrNotImplemented
}

// ReceivePacketContext receives and decodes one packet with cancellation support.
func (c *Connection) ReceivePacketContext(ctx context.Context) (protocol.Packet, *net.UDPAddr, error) {
	return protocol.Packet{}, nil, ErrNotImplemented
}

// Incoming returns the receive channel for packets accepted by the RUDP layer.
func (c *Connection) Incoming() <-chan ReceivedPacket {
	return nil
}

// Errors returns asynchronous connection errors.
func (c *Connection) Errors() <-chan error {
	return nil
}

// SendReliable sends one packet and waits for its ACK.
func (c *Connection) SendReliable(packet protocol.Packet) error {
	return ErrNotImplemented
}
