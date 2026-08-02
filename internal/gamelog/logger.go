package gamelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"jiwa-jawa/internal/engine"
)

type Logger struct {
	mu        sync.Mutex
	file      *os.File
	remoteURL string
	client    *http.Client
}

type Event struct {
	Time      time.Time   `json:"time"`
	Type      string      `json:"type"`
	SessionID uint64      `json:"session_id,omitempty"`
	Side      engine.Side `json:"side,omitempty"`
	Message   string      `json:"message,omitempty"`
	Move      *MoveEvent  `json:"move,omitempty"`
	Seq       uint32      `json:"seq,omitempty"`
	Error     string      `json:"error,omitempty"`
}

type MoveEvent struct {
	Source engine.PointID `json:"source"`
	Target engine.PointID `json:"target"`
}

func New(path string) (*Logger, error) {
	return NewWithRemote(path, "")
}

func NewWithRemote(path, remoteURL string) (*Logger, error) {
	if path == "" {
		path = "logs/game.jsonl"
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	remoteURL = strings.TrimRight(strings.TrimSpace(remoteURL), "/")
	return &Logger{
		file:      file,
		remoteURL: remoteURL,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (l *Logger) Log(event Event) error {
	if l == nil || l.file == nil {
		return nil
	}
	if event.Time.IsZero() {
		event.Time = time.Now()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	l.mu.Lock()
	if _, err = l.file.Write(append(data, '\n')); err == nil {
		err = l.file.Sync()
	}
	l.mu.Unlock()
	if err != nil {
		return err
	}
	if l.remoteURL == "" {
		return nil
	}

	response, err := l.client.Post(l.remoteURL+"/events", "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("raft logger: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("raft logger returned %s", response.Status)
	}
	return nil
}

func (l *Logger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
