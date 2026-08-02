package rating

import (
	"math"

	"jiwa-jawa/internal/player"
)

type PlayerRating struct {
	PlayerID player.ID
	Rating   float64
	Games    int
	Wins     int
	Losses   int
	Draws    int
}

type Result float64

const (
	Loss Result = 0.0
	Draw Result = 0.5
	Win  Result = 1.0
)

type System interface {
	ExpectedScore(a PlayerRating, b PlayerRating) float64
	Update(a PlayerRating, b PlayerRating, result Result) (PlayerRating, PlayerRating, error)
}

type Elo struct {
	KFactor float64
}

func NewElo(kFactor float64) *Elo {
	return &Elo{KFactor: kFactor}
}

func (e *Elo) ExpectedScore(a PlayerRating, b PlayerRating) float64 {
	return 1 / (1 + math.Pow(10, (b.Rating-a.Rating)/400))
}

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
