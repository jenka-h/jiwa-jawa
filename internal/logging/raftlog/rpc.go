package raftlog

// RequestVoteRequest is the Raft election request.
type RequestVoteRequest struct {
	Term         uint64
	CandidateID  uint64
	LastLogIndex uint64
	LastLogTerm  uint64
}

// RequestVoteResponse is the Raft election response.
type RequestVoteResponse struct {
	Term        uint64
	VoteGranted bool
}

// AppendEntriesRequest is used for heartbeats and log replication.
type AppendEntriesRequest struct {
	Term         uint64
	LeaderID     uint64
	PrevLogIndex uint64
	PrevLogTerm  uint64
	Entries      []LogEntry
	LeaderCommit uint64
}

// AppendEntriesResponse is the AppendEntries response.
type AppendEntriesResponse struct {
	Term    uint64
	Success bool
}

// HandleRequestVote handles an incoming RequestVote RPC.
func (n *Node) HandleRequestVote(request RequestVoteRequest) (RequestVoteResponse, error) {
	return RequestVoteResponse{}, ErrNotImplemented
}

// HandleAppendEntries handles an incoming AppendEntries RPC.
func (n *Node) HandleAppendEntries(request AppendEntriesRequest) (AppendEntriesResponse, error) {
	return AppendEntriesResponse{}, ErrNotImplemented
}
