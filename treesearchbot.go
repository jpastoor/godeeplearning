package main

import (
	"math/rand"
	"time"
)

type TreeSearchBot struct {
	maxWidth int
	maxDepth int
}

func (bot TreeSearchBot) selectMove(state GameState) Move {
	legalMoves := state.LegalMoves()

	if len(legalMoves) == 0 {
		return Pass()
	}

	bestSoFar := -99999
	var bestMove Move

	randomMoves := getRandomMoves(bot.maxWidth, legalMoves)

	for _, randomMove := range randomMoves {
		nextState, err := state.ApplyMove(state.NextPlayer, Play(randomMove))
		if err != nil {
			panic(err)
		}
		// Find out our opponent's best result from that position.
		opponentBestResult := bot.bestResult(nextState.Copy(), bot.maxWidth, bot.maxDepth, bot.CaptureDiff)
		ourScore := opponentBestResult * -1

		if ourScore > bestSoFar {
			bestSoFar = ourScore
			bestMove = Play(randomMove)
		}
	}

	return bestMove
}

/**
Find the best result that next_player can get from this game state.
 */
func (bot TreeSearchBot) bestResult(state GameState, maxWidth int, maxDepth int, evalFunc EvalFunc) int {
	if state.IsOver() {
		if state.Winner().equals(state.NextPlayer) {
			return 99999
		} else {
			return -99999
		}
	}

	// We have reached our maximum search depth. Use our heuristic to
	// decide how good this sequence is.
	if maxDepth == 0 {
		return evalFunc(state)
	}

	bestSoFar := -99999

	// Select some random moves!
	legalMoves := state.LegalMoves()
	if len(legalMoves) == 0 {
		return 0
	}

	randomMoves := getRandomMoves(maxWidth, legalMoves)

	for _, randomMove := range randomMoves {
		nextState, err := state.ApplyMove(state.NextPlayer, Play(randomMove))
		if err != nil {
			panic(err)
		}

		// Find out our opponent's best result from that position.
		opponentBestResult := bot.bestResult(nextState.Copy(), maxWidth, maxDepth-1, evalFunc)
		ourResult := opponentBestResult * -1

		if ourResult > bestSoFar {
			bestSoFar = ourResult
		}
	}

	return bestSoFar
}

func getRandomMoves(maxWidth int, legalMoves []Point) []Point {
	maxRandomMoves := maxWidth
	if maxWidth > len(legalMoves) {
		maxRandomMoves = len(legalMoves)
	}
	randomMoves := make([]Point, maxRandomMoves)


	for i := 0; i < maxRandomMoves; i++ {

		r := rand.New(rand.NewSource(time.Now().Unix()*int64(i)))
		randIndex := r.Intn(len(legalMoves))
		randomMoves[i] = legalMoves[randIndex]
	}
	return randomMoves
}

/**
Calculate the difference between the number of black stones and
white stones on the board. This will be the same as the difference
in the number of captures, unless one player passes early.

Returns the difference from the perspective of the next player to
play.
If it's black's move, we return (black stones) - (white stones).
If it's white's move, we return (white stones) - (black stones).
 */
func (bot TreeSearchBot) CaptureDiff(state GameState) int {
	return state.Board.CaptureDiff(state.NextPlayer)
}

type EvalFunc func(state GameState) int

func (b Board) CaptureDiff(nextPlayer Player) int {
	blackStones := 0
	whiteStones := 0

	for row := 1; row <= b.NumRows; row++ {
		for col := 1; col <= b.NumCols; col++ {
			if player, exists := b.get(Point{row, col}); exists {
				if player.equals(PlayerBlack) {
					blackStones++
				} else {
					whiteStones++
				}
			}
		}
	}

	diff := blackStones - whiteStones
	if nextPlayer == PlayerBlack {
		return diff
	}

	return diff * -1
}
