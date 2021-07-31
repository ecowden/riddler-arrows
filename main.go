package main

import (
	"fmt"
	"math/rand"
	"time"
)

type team string

const (
	riddlers       team = "Riddler Nation"
	conundrumers   team = "Conundrum Country"
	numGames            = 10000000000
	arrowsPerRound      = 3
	cScore              = 8 * arrowsPerRound
)

func main() {
	scoreboard := map[team]int{
		riddlers:     0,
		conundrumers: 0,
	}

	r := newRiddler()
	for i := 0; i < numGames; i++ {
		winner := game(r)
		scoreboard[winner] = scoreboard[winner] + 1
	}

	fmt.Printf("Riddlers:     %d\n", scoreboard[riddlers])
	fmt.Printf("Conundrumers: %d\n", scoreboard[conundrumers])
	var riddlerWinRatio float64
	riddlerWinRatio = float64(scoreboard[riddlers]) / numGames
	fmt.Printf("Riddler Win Ratio: %f\n", riddlerWinRatio)
}

func game(r riddler) (winner team) {

	rScore := 0
	for i := 0; i < arrowsPerRound; i++ {
		rScore += r.shoot()
	}

	if rScore > cScore {
		return riddlers
	} else if rScore < cScore {
		return conundrumers
	} else { // Tie. Play a new game.
		return game(r)
	}
}

type riddler struct {
	r *rand.Rand
}

func newRiddler() riddler {
	source := rand.NewSource(time.Now().UnixNano())
	return riddler{
		r: rand.New(source),
	}
}

//shoot an arrow as a member of Riddler Nation:
// > For every shot, each archer of Riddler Nation has
// > a one-third chance of hitting the bull’s-eye (i.e., earning 10 points),
// > a one-third chance of earning 9 points
// > and a one-third chance of earning 5 points.
func (r riddler) shoot() (points int) {
	n := r.r.Intn(3) // TODO oof bad naming
	if n == 0 {
		return 10
	} else if n == 1 {
		return 9
	} else {
		return 5
	}
}
