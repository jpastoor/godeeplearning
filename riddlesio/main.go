package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"github.com/jpastoor/godeeplearning/game"
	"math/rand"
	"time"
	"log"
)

func main() {
	//f, err := os.OpenFile("testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()

	log.SetOutput(bufio.NewWriter(os.Stderr))

	reader := bufio.NewReader(os.Stdin)

	var board game.Board
	settings := &Settings{}
	var legalMoves []game.Point

	var round string
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
			parts[i] = strings.Trim(p," \n\r")
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
						board = game.Board{
							NumRows: settings.fieldHeight,
							NumCols: settings.fieldWidth,
							Grid:    []game.GoString{},
							Hash:    0,
						}

						round = parts[3]
						legalMoves = []game.Point{}
					}

					if parts[2] == "field" {
						stones := strings.Split(parts[3], ",")
						for i, stone := range stones {
							row := (i / 19) + 1
							col := (i % 19) + 1

							point := game.Point{Row: row, Col: col}

							switch stone {
							case ".":
								legalMoves = append(legalMoves, point)
							case "0":
								board.PlaceStone(game.PlayerBlack, point)
							case "1":
								board.PlaceStone(game.PlayerWhite, point)
							}
						}
					}
				}
			}
		case "action":
			{
				if len(legalMoves) > 0 {
					s := rand.NewSource(time.Now().UnixNano())
					r := rand.New(s) // initialize local pseudorandom generator
					randomMoveNr := r.Intn(len(legalMoves))
					randomMr := legalMoves[randomMoveNr]
					log.Printf("Round %s - Legal moves: %d. Chose %d, which is [%d,%d]", round, len(legalMoves), randomMoveNr, randomMr.Row-1, randomMr.Col-1)

					fmt.Printf("place_move %d %d\n", randomMr.Col-1, randomMr.Row-1)
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
