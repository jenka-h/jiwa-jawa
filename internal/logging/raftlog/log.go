package raftlog

import err "jiwa-jawa/internal/error"

// LogEntry is one replicated Raft log command.
type LogEntry struct {
	Index   uint64
	Term    uint64
	Command []byte
}

// Append appends a command through Raft consensus.
func (n *Node) Append(command []byte) (LogEntry, error) {
	return LogEntry{}, err.ErrNotImplemented
}

// Entries returns the current known log entries.
func (n *Node) Entries() []LogEntry {
	return nil
}
