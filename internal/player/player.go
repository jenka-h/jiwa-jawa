package player

// ID uniquely identifies one real player/user across sessions, logs, and ratings.
type ID string

// Player stores stable player identity information.
type Player struct {
	ID   ID
	Name string
}
