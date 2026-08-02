package raftlog

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"

	"jiwa-jawa/internal/gamelog"
)

type command struct {
	Event gamelog.Event `json:"event"`
}

type Store struct {
	raft *raft.Raft
	fsm  *eventFSM
}

type eventFSM struct {
	mu     sync.RWMutex
	events []gamelog.Event
}

func Open(nodeID, raftAddress, dataDir string, bootstrap bool) (*Store, error) {
	if nodeID == "" || raftAddress == "" || dataDir == "" {
		return nil, fmt.Errorf("node id, raft address, and data directory are required")
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)
	config.SnapshotInterval = 20 * time.Second
	config.SnapshotThreshold = 64

	address, err := net.ResolveTCPAddr("tcp", raftAddress)
	if err != nil {
		return nil, fmt.Errorf("resolve raft address: %w", err)
	}
	transport, err := raft.NewTCPTransport(raftAddress, address, 5, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("create raft transport: %w", err)
	}
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft.db"))
	if err != nil {
		_ = transport.Close()
		return nil, fmt.Errorf("open raft database: %w", err)
	}
	snapshots, err := raft.NewFileSnapshotStore(dataDir, 2, os.Stderr)
	if err != nil {
		_ = transport.Close()
		return nil, fmt.Errorf("open snapshot store: %w", err)
	}

	fsm := &eventFSM{}
	r, err := raft.NewRaft(config, fsm, logStore, logStore, snapshots, transport)
	if err != nil {
		_ = transport.Close()
		return nil, fmt.Errorf("start raft: %w", err)
	}
	if bootstrap {
		future := r.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: config.LocalID, Address: transport.LocalAddr()}}})
		if err := future.Error(); err != nil && err != raft.ErrCantBootstrap {
			_ = r.Shutdown().Error()
			return nil, fmt.Errorf("bootstrap cluster: %w", err)
		}
	}
	return &Store{raft: r, fsm: fsm}, nil
}

func (s *Store) Append(event gamelog.Event) error {
	if event.Time.IsZero() {
		event.Time = time.Now().UTC()
	}
	data, err := json.Marshal(command{Event: event})
	if err != nil {
		return err
	}
	if err := s.raft.Apply(data, 10*time.Second).Error(); err != nil {
		return fmt.Errorf("commit event: %w", err)
	}
	return nil
}

func (s *Store) Events(sessionID uint64) []gamelog.Event {
	return s.fsm.eventsFor(sessionID)
}

func (s *Store) Join(id, address string) error {
	if s.raft.State() != raft.Leader {
		return raft.ErrNotLeader
	}
	configuration := s.raft.GetConfiguration()
	if err := configuration.Error(); err != nil {
		return err
	}
	for _, server := range configuration.Configuration().Servers {
		if server.ID == raft.ServerID(id) && server.Address == raft.ServerAddress(address) {
			return nil
		}
		if server.ID == raft.ServerID(id) || server.Address == raft.ServerAddress(address) {
			if err := s.raft.RemoveServer(server.ID, 0, 10*time.Second).Error(); err != nil {
				return err
			}
		}
	}
	return s.raft.AddVoter(raft.ServerID(id), raft.ServerAddress(address), 0, 10*time.Second).Error()
}

func (s *Store) Leader() string { return string(s.raft.Leader()) }
func (s *Store) State() string  { return s.raft.State().String() }
func (s *Store) IsLeader() bool { return s.raft.State() == raft.Leader }
func (s *Store) Close() error   { return s.raft.Shutdown().Error() }

func (f *eventFSM) Apply(entry *raft.Log) any {
	var cmd command
	if err := json.Unmarshal(entry.Data, &cmd); err != nil {
		return err
	}
	f.mu.Lock()
	f.events = append(f.events, cmd.Event)
	f.mu.Unlock()
	return nil
}

func (f *eventFSM) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.RLock()
	copyOfEvents := append([]gamelog.Event(nil), f.events...)
	f.mu.RUnlock()
	return eventSnapshot{events: copyOfEvents}, nil
}

func (f *eventFSM) Restore(reader io.ReadCloser) error {
	defer reader.Close()
	var events []gamelog.Event
	if err := json.NewDecoder(reader).Decode(&events); err != nil {
		return err
	}
	f.mu.Lock()
	f.events = events
	f.mu.Unlock()
	return nil
}

func (f *eventFSM) eventsFor(sessionID uint64) []gamelog.Event {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if sessionID == 0 {
		return append([]gamelog.Event(nil), f.events...)
	}
	result := make([]gamelog.Event, 0)
	for _, event := range f.events {
		if event.SessionID == sessionID {
			result = append(result, event)
		}
	}
	return result
}

type eventSnapshot struct{ events []gamelog.Event }

func (s eventSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := json.NewEncoder(sink).Encode(s.events); err != nil {
		_ = sink.Cancel()
		return err
	}
	return sink.Close()
}
func (eventSnapshot) Release() {}
