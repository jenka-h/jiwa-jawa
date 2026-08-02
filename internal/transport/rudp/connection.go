package rudp

import (
	"context"
	"net"
	"sync"
	"time"

	err "jiwa-jawa/internal/error"
	"jiwa-jawa/internal/transport/protocol"
)

type ReceivedPacket struct {
	Packet protocol.Packet
	Addr   *net.UDPAddr
}

type Connection struct {
	conn *net.UDPConn
	peer *net.UDPAddr

	sessionID uint64
	nextSeq   uint32

	pending map[uint32]chan struct{}
	seen    map[uint32]struct{}

	incoming chan ReceivedPacket
	errors   chan error
	done     chan struct{}

	timeout    time.Duration
	maxRetries int
	bufferSize int

	mu        sync.Mutex
	closeOnce sync.Once
}

func NewConnection(cfg Config) (*Connection, error) {
	listenAddr, err := net.ResolveUDPAddr("udp", cfg.ListenAddr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	bufferSize := cfg.BufferSize
	if bufferSize <= 0 {
		bufferSize = DefaultBufferSize
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}

	conn := &Connection{
		conn:       udpConn,
		sessionID:  cfg.SessionID,
		nextSeq:    1,
		pending:    make(map[uint32]chan struct{}),
		seen:       make(map[uint32]struct{}),
		incoming:   make(chan ReceivedPacket, bufferSize),
		errors:     make(chan error, bufferSize),
		done:       make(chan struct{}),
		timeout:    timeout,
		maxRetries: maxRetries,
		bufferSize: bufferSize,
	}

	if cfg.PeerAddr != "" {
		if err := conn.SetPeer(cfg.PeerAddr); err != nil {
			_ = udpConn.Close()
			return nil, err
		}
	}

	return conn, nil
}

func (c *Connection) Start(ctx context.Context) error {
	go c.receiveLoop(ctx)
	return nil
}

func (c *Connection) Close() error {
	var closeErr error

	c.closeOnce.Do(func() {
		close(c.done)

		c.mu.Lock()
		for sequence := range c.pending {
			delete(c.pending, sequence)
		}
		c.mu.Unlock()

		closeErr = c.conn.Close()
	})

	return closeErr
}

func (c *Connection) Done() <-chan struct{} {
	return c.done
}

func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connection) PeerAddr() *net.UDPAddr {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.peer
}

func (c *Connection) SetSessionID(sessionID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessionID = sessionID
}

func (c *Connection) SetPeer(addr string) error {
	peer, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.peer = peer
	c.mu.Unlock()

	return nil
}

func (c *Connection) SendPacket(packet protocol.Packet) error {
	c.mu.Lock()
	peer := c.peer
	c.mu.Unlock()

	if peer == nil {
		return err.ErrPeerNotSet
	}

	return c.sendPacketTo(packet, peer)
}

func (c *Connection) SendReliable(packet protocol.Packet) error {
	c.mu.Lock()
	peer := c.peer
	if peer == nil {
		c.mu.Unlock()
		return err.ErrPeerNotSet
	}

	sequence := c.nextSequenceLocked()
	packet.Header.Sequence = sequence
	packet.Header.SessionID = c.sessionID
	packet.Header.Flags |= protocol.FlagReliable

	acked := make(chan struct{})
	c.pending[sequence] = acked
	c.mu.Unlock()

	defer c.removePending(sequence)

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if err := c.sendPacketTo(packet, peer); err != nil {
			return err
		}

		select {
		case <-acked:
			return nil
		case <-time.After(c.timeout):
		case <-c.done:
			return err.ErrClosed
		}
	}

	return err.ErrMaxRetries
}

func (c *Connection) ReceivePacketContext(ctx context.Context) (protocol.Packet, *net.UDPAddr, error) {
	buffer := make([]byte, c.bufferSize)

	for {
		select {
		case <-ctx.Done():
			return protocol.Packet{}, nil, ctx.Err()
		case <-c.done:
			return protocol.Packet{}, nil, err.ErrClosed
		default:
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		n, addr, readErr := c.conn.ReadFromUDP(buffer)
		if readErr != nil {
			if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
				continue
			}
			return protocol.Packet{}, nil, readErr
		}

		packet, decodeErr := protocol.DecodePacket(buffer[:n])
		if decodeErr != nil {
			return protocol.Packet{}, nil, decodeErr
		}

		if packet.IsACK() {
			c.handleACK(packet)
			continue
		}

		if packet.IsReliable() {
			_ = c.sendACKTo(packet.Header.Sequence, addr)
			if c.markDuplicate(packet.Header.Sequence) {
				continue
			}
		}

		return packet, addr, nil
	}
}

func (c *Connection) Incoming() <-chan ReceivedPacket {
	return c.incoming
}

func (c *Connection) Errors() <-chan error {
	return c.errors
}

func (c *Connection) receiveLoop(ctx context.Context) {
	for {
		packet, addr, err := c.ReceivePacketContext(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case <-c.done:
				return
			default:
				c.report(err)
				continue
			}
		}

		select {
		case c.incoming <- ReceivedPacket{Packet: packet, Addr: addr}:
		case <-ctx.Done():
			return
		case <-c.done:
			return
		}
	}
}

func (c *Connection) sendPacketTo(packet protocol.Packet, addr *net.UDPAddr) error {
	data, err := protocol.EncodePacket(packet)
	if err != nil {
		return err
	}

	_, err = c.conn.WriteToUDP(data, addr)
	return err
}

func (c *Connection) sendACKTo(sequence uint32, addr *net.UDPAddr) error {
	return c.sendPacketTo(protocol.NewACK(c.sessionID, sequence), addr)
}

func (c *Connection) handleACK(packet protocol.Packet) {
	c.mu.Lock()
	acked, ok := c.pending[packet.Header.Sequence]
	if ok {
		delete(c.pending, packet.Header.Sequence)
	}
	c.mu.Unlock()

	if ok {
		close(acked)
	}
}

func (c *Connection) markDuplicate(sequence uint32) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.seen[sequence]; ok {
		return true
	}

	c.seen[sequence] = struct{}{}
	return false
}

func (c *Connection) removePending(sequence uint32) {
	c.mu.Lock()
	delete(c.pending, sequence)
	c.mu.Unlock()
}

func (c *Connection) nextSequenceLocked() uint32 {
	sequence := c.nextSeq
	c.nextSeq++
	if c.nextSeq == 0 {
		c.nextSeq = 1
	}
	return sequence
}

func (c *Connection) report(reportErr error) {
	select {
	case c.errors <- reportErr:
	default:
	}
}
