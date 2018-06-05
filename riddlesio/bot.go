package main

import (
	"math/rand"
	"time"
	"sync"
	"github.com/jpastoor/godeeplearning/game"
)

type Agent interface {
	selectMove(state game.GameState) game.Move
}

type RandomBot struct {
}

func (agent *RandomBot) selectMove(state game.GameState) game.Move {
	// Choose a random valid move that preserves our own eyes.
	var lock sync.Mutex
	var candidates []game.Point

	var wg sync.WaitGroup

	for row := 1; row <= state.Board.NumRows; row++ {
		for col := 1; col <= state.Board.NumCols; col++ {
			candidate := game.Point{row, col}

			wg.Add(1)
			go func() {
				defer wg.Done()

				if state.IsMoveValid(game.Play(candidate)) && !state.Board.IsPointAnEye(candidate, state.NextPlayer) {
					lock.Lock()
					candidates = append(candidates, candidate)
					lock.Unlock()
				}
			}()
		}
	}

	wg.Wait()

	if len(candidates) == 0 {
		return game.Pass()
	} else {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s) // initialize local pseudorandom generator
		randomMove := r.Intn(len(candidates))
		return game.Play(candidates[randomMove])
	}
}
