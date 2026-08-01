package app

// Config contains runtime options for one Dam Daman client.
type Config struct {
	PlayerID   string
	PlayerName string
	ListenAddr string
	PeerAddr   string
	SessionID  uint64
	LogPath    string
}
