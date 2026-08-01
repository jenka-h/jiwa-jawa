package rating

import err "jiwa-jawa/internal/error"

// Elo implements a basic Elo-style rating system.
type Elo struct {
	KFactor float64
}

// NewElo creates an Elo rating system.
func NewElo(kFactor float64) *Elo {
	return nil
}

// ExpectedScore returns player a's expected score against player b.
func (e *Elo) ExpectedScore(a PlayerRating, b PlayerRating) float64 {
	return 0
}

// Update updates ratings after one game.
func (e *Elo) Update(a PlayerRating, b PlayerRating, result Result) (PlayerRating, PlayerRating, error) {
	return PlayerRating{}, PlayerRating{}, err.ErrNotImplemented
}
