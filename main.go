package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type team int

const (
	riddlers = iota
	conundrumers

	arrowsPerRound = 3
	cScore         = 8 * arrowsPerRound
)

var (
	games   int64
	workers int
)

func main() {
	flag.Int64Var(&games, "g", 10000, "number of games to simulate (default: 10000)")
	flag.IntVar(&workers, "w", runtime.NumCPU(), "number of concurrent workers (default: number of CPUs)")
	flag.Parse()
	start := time.Now()
	var rWins, cWins int64

	// Create workers
	bar := pb.Start64(games)
	gamesPerWorker := games / int64(workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	rWinsCh := make(chan int64, workers)
	cWinsCh := make(chan int64, workers)
	for i := 0; i < workers-1; i++ {
		go worker(gamesPerWorker, rWinsCh, cWinsCh, &wg, bar)
		time.Sleep(2 * time.Nanosecond) // cheap way to ensure each worker has a unique seed
	}
	// Schedule last worker with remaining games to avoid rounding errors
	remainingGames := games - (gamesPerWorker * (int64(workers) - 1))
	go worker(remainingGames, rWinsCh, cWinsCh, &wg, bar)

	// Read results
	go func() {
		for n := range rWinsCh {
			rWins += n
		}
	}()
	go func() {
		for n := range cWinsCh {
			cWins += n
		}
	}()

	// Wait for completion
	wg.Wait()
	close(rWinsCh)
	close(cWinsCh)

	// Print results
	bar.Finish()
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Riddlers:     %d\n", rWins)
	fmt.Printf("Conundrumers: %d\n", cWins)
	var riddlerWinRatio float64
	riddlerWinRatio = float64(rWins) / float64(games)
	fmt.Printf("Riddler Win Ratio: %f\n", riddlerWinRatio)
}

func worker(n int64, rWinsCh chan<- int64, cWinsCh chan<- int64, wg *sync.WaitGroup, bar *pb.ProgressBar) {
	r := newRiddler()
	var rWins, cWins int64
	var i int64
	for i = 0; i < n; i++ {
		winner := game(r)
		if winner == riddlers {
			rWins++
		} else {
			cWins++
		}
		bar.Increment()
	}
	rWinsCh <- rWins
	cWinsCh <- cWins
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
