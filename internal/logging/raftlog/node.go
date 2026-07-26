package raftlog

import "context"

// NodeState is the Raft node role.
type NodeState uint8

const (
	Follower NodeState = iota
	Candidate
	Leader
)

// Node represents one Raft logger node.
type Node struct {
	ID    uint64
	State NodeState
}

// Config contains Raft logger node settings.
type Config struct {
	ID         uint64
	ListenAddr string
	Peers      map[uint64]string
	DataDir    string
}

// NewNode creates a Raft logger node.
func NewNode(config Config) (*Node, error) {
	return nil, ErrNotImplemented
}

// Start starts Raft timers and RPC handling.
func (n *Node) Start(ctx context.Context) error {
	return ErrNotImplemented
}

// Stop stops the Raft node.
func (n *Node) Stop(ctx context.Context) error {
	return ErrNotImplemented
}

// IsLeader reports whether this node is currently leader.
func (n *Node) IsLeader() bool {
	return false
}
