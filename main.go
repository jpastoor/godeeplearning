package main

import (
	"fmt"
	"time"
	"github.com/jpastoor/godeeplearning/game"
)

func main() {

	boardSize := 9

	players := map[game.Player]Agent{
		game.PlayerBlack: &RandomBot{},
		game.PlayerWhite: &RandomBot{},
	}

	myGame := game.NewGame(boardSize)
	var err error
	moveNr := 1
	for !myGame.IsOver() {
		time.Sleep(1 * time.Second)
		myGame.Print()
		player := myGame.NextPlayer

		nextMove := players[player].selectMove(myGame)

		myGame, err = myGame.ApplyMove(player, nextMove)

		fmt.Printf("-- Move %d --\n", moveNr)
		moveNr++

		if err != nil {
			fmt.Printf("Error during move! %s\n", err)
			break
		}
	}

	fmt.Println("Game is over")
}
