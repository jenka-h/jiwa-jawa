package gui

import (
	"html/template"
	"net/http"
	"time"
)

func NewServer(assetsDir string) *Server {
	return NewServerWithOptions(assetsDir, Options{})
}

func NewServerWithOptions(assetsDir string, options Options) *Server {
	if options.ListenAddr == "" {
		options.ListenAddr = ":9001"
	}
	if options.PeerAddr == "" {
		options.PeerAddr = "127.0.0.1:9002"
	}
	if options.SessionID == 0 {
		options.SessionID = uint64(time.Now().UnixNano())
	}
	if options.LogPath == "" {
		options.LogPath = "logs/gui-game.jsonl"
	}
	if options.ProfilePath == "" {
		options.ProfilePath = "data/profiles.json"
	}
	server := &Server{
		assetsDir:     assetsDir,
		options:       options,
		nextProfileID: 1,
		profilePath:   options.ProfilePath,
		profiles:      []Profile{},
	}
	server.profileErr = server.loadProfiles()
	return server
}

func (s *Server) Handler() (http.Handler, error) {
	page, err := template.New("game").Parse(pageHTML)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(s.assetsDir))))
	mux.HandleFunc("/api/profiles", s.handleProfiles)
	mux.HandleFunc("/api/session", s.handleSession)
	mux.HandleFunc("/api/state", s.handleState)
	mux.HandleFunc("/api/move", s.handleMove)
	mux.HandleFunc("/api/penalty", s.handlePenalty)
	mux.HandleFunc("/api/end-turn", s.handleEndTurn)
	mux.HandleFunc("/api/finish", s.handleFinish)
	mux.HandleFunc("/api/board", s.handleBoard)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if err := page.Execute(w, demoPageData()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return mux, nil
}

func (s *Server) ListenAndServe(addr string) error {
	handler, err := s.Handler()
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, handler)
}
