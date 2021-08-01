package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type team int

const (
	riddlers = iota
	conundrumers

	arrowsPerRound = 3
	cScore         = 8 * arrowsPerRound
)

var (
	numGames int64
	workers  int
)

func main() {
	flag.Int64Var(&numGames, "n", 10000, "number of games to simulate (default: 10000)")
	flag.IntVar(&workers, "w", runtime.NumCPU(), "number of concurrent workers (default: number of CPUs)")
	flag.Parse()

	fmt.Printf("Simulating %d games using %d workers...\n", numGames, workers)

	scoreboard := map[team]int{
		riddlers:     0,
		conundrumers: 0,
	}

	// Create workers
	results := make(chan team)
	var wg sync.WaitGroup
	gamesPerWorker := numGames / int64(workers) // TODO recipe for rounding errors
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(gamesPerWorker, results, &wg)
	}

	go func() { // Close the results channel when all workers are done
		wg.Wait()
		close(results)
	}()

	for winner := range results { // Read results
		scoreboard[winner] = scoreboard[winner] + 1
	}

	// Print results
	fmt.Printf("Riddlers:     %d\n", scoreboard[riddlers])
	fmt.Printf("Conundrumers: %d\n", scoreboard[conundrumers])
	var riddlerWinRatio float64
	riddlerWinRatio = float64(scoreboard[riddlers]) / float64(numGames)
	fmt.Printf("Riddler Win Ratio: %f\n", riddlerWinRatio)
}

func worker(n int64, results chan<- team, wg *sync.WaitGroup) {
	r := newRiddler()
	var i int64
	for i = 0; i < n; i++ {
		results <- game(r)
	}
	wg.Done()
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
