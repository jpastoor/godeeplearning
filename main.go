package main

import (
	"log"
	"bufio"
	"os"
	"strings"
	"fmt"
)

//func main() {
//
//	boardSize := 3
//
//	players := map[Player]Agent{
//		PlayerBlack: &TreeSearchBot{
//			5,
//			5,
//		},
//		PlayerWhite: &RandomBot{},
//	}
//
//	scoreBoard := map[Player]int{
//		PlayerBlack: 0,
//		PlayerWhite: 0,
//	}
//
//	maxGames := 100
//
//	for gameNr := 1; gameNr <= maxGames; gameNr++ {
//
//		state := NewGame(boardSize)
//		var err error
//		moveNr := 1
//		for !state.IsOver() {
//			//time.Sleep(1 * time.Second)
//		//	state.Board.print(false)
//			player := state.NextPlayer
//
//			nextMove := players[player].selectMove(state)
//
//			state, err = state.ApplyMove(player, nextMove)
//
//			//fmt.Printf("-- Move %d --\n", moveNr)
//			moveNr++
//
//			if err != nil {
//				fmt.Printf("Error during move! %s\n", err)
//				break
//			}
//		}
//
//		winner := state.Winner()
//		scoreBoard[*winner]++
//		fmt.Printf("Game %d is over. Winner: %s\n", gameNr, winner)
//	}
//
//	fmt.Printf("Scoreboard: %v\n", scoreBoard)
//}

func main() {
	//f, err := os.OpenFile("testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()

	log.SetOutput(bufio.NewWriter(os.Stderr))

	reader := bufio.NewReader(os.Stdin)

	var board Board
	settings := &Settings{}
	var legalMoves []Point

	bot := TreeSearchBot{
		maxWidth: 3,
		maxDepth: 3,
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, " ")
		for i, p := range parts {
			parts[i] = strings.Trim(p, " \n\r")
		}

		switch parts[0] {
		case "settings":
			{
				err := settings.update(parts[1], parts[2])
				if err != nil {
					log.Printf("Could update settings %v", err)
					panic(err)
				}
			}
		case "update":
			{
				if parts[1] == "game" {
					if parts[2] == "round" {
						board = Board{
							NumRows: settings.fieldHeight,
							NumCols: settings.fieldWidth,
							Grid:    []GoString{},
							Hash:    0,
						}

						legalMoves = []Point{}
					}

					if parts[2] == "field" {
						stones := strings.Split(parts[3], ",")
						for i, stone := range stones {
							row := (i / 19) + 1
							col := (i % 19) + 1

							point := Point{Row: row, Col: col}

							switch stone {
							case ".":
								legalMoves = append(legalMoves, point)
							case "0":
								board.PlaceStone(PlayerBlack, point)
							case "1":
								board.PlaceStone(PlayerWhite, point)
							}
						}
					}
				}
			}
		case "action":
			{
				player := PlayerWhite
				if settings.yourBotId == 0 {
					player = PlayerBlack
				}

				state := GameState{
					NextPlayer:     player,
					Board:          board,
					PreviousStates: nil,
					PreviousState:  nil,
				}

				log.Println("Starting move selection...")
				move := bot.selectMove(state)

				if (move.IsPlay) {
					fmt.Printf("place_move %d %d\n", move.Point.Col-1, move.Point.Row-1)
				} else {
					fmt.Printf("pass\n", board)
				}
			}
		default:
			{
				log.Printf("Unsupported command: %s", line)
				panic(err)
			}
		}
	}
}
