package gui

import (
	"context"
	"sync"
	"time"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/gamelog"
	"jiwa-jawa/internal/transport/rudp"
)

type Server struct {
	assetsDir string
	options   Options

	mu            sync.Mutex
	nextProfileID int
	profilePath   string
	profileErr    error
	profiles      []Profile
	match         *Match
}

type Options struct {
	ListenAddr  string
	PeerAddr    string
	SessionID   uint64
	LogPath     string
	LogServer   string
	ProfilePath string
}

type Match struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	conn   *rudp.Connection
	logger *gamelog.Logger
	state  *engine.GameState

	role          string
	status        string
	localID       string
	localName     string
	localRating   float64
	remoteID      string
	remoteName    string
	remoteRating  float64
	ratingApplied bool
	localSide     engine.Side
	remoteSide    engine.Side
	sessionID     uint64
	startedAt     time.Time
	logs          []LogItem
}

type Profile struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}

type PageData struct {
	Players  []PlayerView
	Board    BoardView
	Logs     []LogItem
	Duration string
	TurnText string
	Elo      EloView
}

type PlayerView struct {
	Name     string
	Alias    string
	Side     engine.Side
	Piece    string
	Captured int
	Rating   float64
}

type BoardView struct {
	Width  int
	Height int
	Points []PointView
	Edges  []EdgeView
}

type PointView struct {
	ID    engine.PointID
	X     int
	Y     int
	Piece string
	Side  engine.Side
}

type EdgeView struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type LogItem struct {
	No     int    `json:"no"`
	Player string `json:"player"`
	Action string `json:"action"`
	Move   string `json:"move"`
}

type EloView struct {
	PlayerOneBefore float64
	PlayerTwoBefore float64
	PlayerOneAfter  float64
	PlayerTwoAfter  float64
}
