package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/rating"
)

func (s *Server) loadProfiles() error {
	data, err := os.ReadFile(s.profilePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read profiles: %w", err)
	}
	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return fmt.Errorf("decode profiles: %w", err)
	}
	s.profiles = profiles
	for _, profile := range profiles {
		var number int
		if _, err := fmt.Sscanf(profile.ID, "p%d", &number); err == nil && number >= s.nextProfileID {
			s.nextProfileID = number + 1
		}
	}
	return nil
}

func (s *Server) saveProfilesLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.profilePath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.profiles, "", "  ")
	if err != nil {
		return err
	}
	temporary := s.profilePath + ".tmp"
	if err := os.WriteFile(temporary, append(data, '\n'), 0o644); err != nil {
		return err
	}
	if err := os.Rename(temporary, s.profilePath); err != nil {
		_ = os.Remove(temporary)
		return err
	}
	return nil
}

func normalizedRating(value float64) float64 {
	if value <= 0 {
		return 1200
	}
	return value
}

func (s *Server) applyRating(match *Match) {
	match.mu.Lock()
	if match.ratingApplied || !match.state.Finished || match.state.Winner == engine.SideNone {
		match.mu.Unlock()
		return
	}
	match.ratingApplied = true
	local := rating.PlayerRating{Rating: normalizedRating(match.localRating)}
	remote := rating.PlayerRating{Rating: normalizedRating(match.remoteRating)}
	result := rating.Loss
	if match.state.Winner == match.localSide {
		result = rating.Win
	}
	updated, _, err := rating.NewElo(32).Update(local, remote, result)
	if err != nil {
		match.mu.Unlock()
		return
	}
	match.localRating = updated.Rating
	localID := match.localID
	match.mu.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.profiles {
		if s.profiles[i].ID == localID {
			s.profiles[i].Rating = updated.Rating
			_ = s.saveProfilesLocked()
			return
		}
	}
}

func (s *Server) handleProfiles(w http.ResponseWriter, r *http.Request) {
	if s.profileErr != nil {
		http.Error(w, s.profileErr.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		s.writeProfiles(w)
	case http.MethodPost:
		var request struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid profile request", http.StatusBadRequest)
			return
		}
		name := strings.TrimSpace(request.Name)
		if name == "" {
			http.Error(w, "profile name is required", http.StatusBadRequest)
			return
		}
		s.mu.Lock()
		profile := Profile{ID: fmt.Sprintf("p%d", s.nextProfileID), Name: name, Rating: 1200}
		s.nextProfileID++
		s.profiles = append(s.profiles, profile)
		if err := s.saveProfilesLocked(); err != nil {
			s.profiles = s.profiles[:len(s.profiles)-1]
			s.nextProfileID--
			s.mu.Unlock()
			http.Error(w, "save profile: "+err.Error(), http.StatusInternalServerError)
			return
		}
		s.mu.Unlock()
		writeJSON(w, http.StatusCreated, profile)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) writeProfiles(w http.ResponseWriter) {
	s.mu.Lock()
	profiles := append([]Profile(nil), s.profiles...)
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, profiles)
}

func (s *Server) profileByID(id string) (Profile, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, profile := range s.profiles {
		if profile.ID == id {
			return profile, true
		}
	}
	return Profile{}, false
}
