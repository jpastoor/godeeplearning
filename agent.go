package main

import (
	"math/rand"
	"time"
	"sync"
)

type Agent interface {
	selectMove(state GameState) Move
}

type RandomBot struct {
}

func (agent *RandomBot) selectMove(state GameState) Move {
	// Choose a random valid move that preserves our own eyes.

	/**
	candidates = []
        for r in range(1, game_state.Board.num_rows + 1):
            for c in range(1, game_state.Board.num_cols + 1):
                candidate = Point(Row=r, Col=c)
                if game_state.is_valid_move(Move.play(candidate)) and \
                        not is_point_an_eye(game_state.Board,
                                            candidate,
                                            game_state.next_player):
                    candidates.append(candidate)

	 */

	var lock sync.Mutex
	var candidates []Point

	var wg sync.WaitGroup

	for row := 1; row <= state.Board.NumRows; row++ {
		for col := 1; col <= state.Board.NumCols; col++ {
			candidate := Point{row, col}

			wg.Add(1)
			go func() {
				defer wg.Done()

				if state.isMoveValid(Play(candidate)) && !state.Board.isPointAnEye(candidate, state.NextPlayer) {
					lock.Lock()
					candidates = append(candidates, candidate)
					lock.Unlock()
				}
			}()
		}
	}

	wg.Wait()

	if len(candidates) == 0 {
		return Pass()
	} else {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s) // initialize local pseudorandom generator
		randomMove := r.Intn(len(candidates))
		return Play(candidates[randomMove])
	}

}
