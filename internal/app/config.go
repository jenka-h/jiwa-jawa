package app

import "jiwa-jawa/internal/player"

type Config struct {
	Player     player.Player
	ListenAddr string
	PeerAddr   string
	SessionID  uint64
	LogPath    string
}
