package main

import (
	"fmt"
	"time"
)

func main() {

	boardSize := 9

	players := map[Player]Agent{
		PlayerBlack: &RandomBot{},
		PlayerWhite: &RandomBot{},
	}

	game := NewGame(boardSize)
	var err error
	moveNr := 1
	for !game.isOver() {
		time.Sleep(1 * time.Second)
		game.print()
		player := game.NextPlayer

		nextMove := players[player].selectMove(game)

		game, err = game.applyMove(player, nextMove)

		fmt.Printf("-- Move %d --\n", moveNr)
		moveNr++

		if err != nil {
			fmt.Printf("Error during move! %s\n", err)
			break
		}
	}

	fmt.Println("Game is over")
}
