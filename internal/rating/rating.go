package rating

// PlayerRating stores rating information for a Dam Daman player.
type PlayerRating struct {
	PlayerID string
	Rating   float64
	Games    int
	Wins     int
	Losses   int
	Draws    int
}

// Result describes the result of one game for rating calculation.
type Result float64

const (
	Loss Result = 0.0
	Draw Result = 0.5
	Win  Result = 1.0
)

// System defines rating behavior.
type System interface {
	ExpectedScore(a PlayerRating, b PlayerRating) float64
	Update(a PlayerRating, b PlayerRating, result Result) (PlayerRating, PlayerRating, error)
}
