package rudp

import "net"

// Peer describes a remote RUDP endpoint.
type Peer struct {
	ID      string
	Name    string
	Address *net.UDPAddr
}

// NewPeer creates a peer descriptor from a UDP address string.
func NewPeer(id string, name string, addr string) (*Peer, error) {
	return nil, ErrNotImplemented
}
