package raftlog

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/hashicorp/raft"

	"jiwa-jawa/internal/gamelog"
)

type HTTPServer struct{ store *Store }

func NewHTTPHandler(store *Store) http.Handler {
	server := &HTTPServer{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.health)
	mux.HandleFunc("/events", server.events)
	mux.HandleFunc("/cluster/join", server.join)
	return mux
}

func (s *HTTPServer) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"state":  s.store.State(),
		"leader": s.store.Leader(),
	})
}

func (s *HTTPServer) events(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var sessionID uint64
		if raw := r.URL.Query().Get("session_id"); raw != "" {
			parsed, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				http.Error(w, "invalid session_id", http.StatusBadRequest)
				return
			}
			sessionID = parsed
		}
		writeJSON(w, http.StatusOK, s.store.Events(sessionID))
	case http.MethodPost:
		if !s.store.IsLeader() {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "not raft leader", "leader": s.store.Leader()})
			return
		}
		var event gamelog.Event
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&event); err != nil {
			http.Error(w, "invalid event: "+err.Error(), http.StatusBadRequest)
			return
		}
		if event.Type == "" {
			http.Error(w, "event type is required", http.StatusBadRequest)
			return
		}
		if err := s.store.Append(event); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusCreated, event)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *HTTPServer) join(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		ID      string `json:"id"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&request); err != nil || request.ID == "" || request.Address == "" {
		http.Error(w, "id and address are required", http.StatusBadRequest)
		return
	}
	if err := s.store.Join(request.ID, request.Address); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, raft.ErrNotLeader) {
			status = http.StatusConflict
		}
		writeJSON(w, status, map[string]string{"error": err.Error(), "leader": s.store.Leader()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "joined"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
