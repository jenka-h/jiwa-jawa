package rating

import "math"

// Elo implements a basic Elo-style rating system.
type Elo struct {
	KFactor float64
}

// NewElo creates an Elo rating system.
func NewElo(kFactor float64) *Elo {
	return &Elo{KFactor: kFactor}
}

// ExpectedScore returns player a's expected score against player b.
func (e *Elo) ExpectedScore(a PlayerRating, b PlayerRating) float64 {
	return 1 / (1 + math.Pow(10, (b.Rating-a.Rating)/400))
}

// Update updates ratings after one game.
func (e *Elo) Update(a PlayerRating, b PlayerRating, result Result) (PlayerRating, PlayerRating, error) {
	expectedA := e.ExpectedScore(a, b)
	expectedB := e.ExpectedScore(b, a)

	scoreA := float64(result)
	scoreB := 1 - scoreA

	a.Rating = a.Rating + e.KFactor*(scoreA-expectedA)
	b.Rating = b.Rating + e.KFactor*(scoreB-expectedB)

	a.Games++
	b.Games++

	switch result {
	case Win:
		a.Wins++
		b.Losses++
	case Draw:
		a.Draws++
		b.Draws++
	case Loss:
		a.Losses++
		b.Wins++
	}

	return a, b, nil
}
