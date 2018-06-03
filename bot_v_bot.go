package main

import (
	"fmt"
	"strings"
)

const COLS = "ABCDEFGHJKLMNOPQRST"

/**
COLS = 'ABCDEFGHJKLMNOPQRST'
STONE_TO_CHAR = {
    None: '.',
    gotypes.Player.black: 'x',
    gotypes.Player.white: 'o',
}


def print_move(player, move):
    if move.is_pass:
        move_str = 'passes'
    elif move.is_resign:
        move_str = 'resigns'
    else:
        move_str = '%s%d' % (COLS[move.Point.Col - 1], move.Point.Row)
    print('%s %s' % (player, move_str))


def print_board(Board):
    for Row in range(Board.num_rows, 0, -1):
        line = []
        for Col in range(1, Board.num_cols + 1):
            stone = Board.get(gotypes.Point(Row=Row, Col=Col))
            line.append(STONE_TO_CHAR[stone])
        print('%d %s' % (Row, ''.join(line)))
    print('  ' + COLS[:Board.num_cols])
 */

func (b Board) print(printGoStrings bool) {

	fmt.Printf("  %s\n", COLS[0:b.NumCols])
	for row := 1; row <= b.NumRows; row++ {
		colStr := ""
		for col := 1; col <= b.NumCols; col++ {
			stone, exists := b.get(Point{row, col})
			if !exists {
				colStr += " "
			} else if stone.isBlack {
				colStr += "x"
			} else {
				colStr += "o"
			}
		}

		fmt.Printf("%d|%s\n", row, colStr)
	}

	if printGoStrings {
		fmt.Printf("\nGo Strings (%d)\n", len(b.Grid))
		for _, gs := range b.Grid {
			gs.print()
		}
	}
}

func (state GameState) print() {

	if state.PreviousState != nil {
		fmt.Printf("%s %s\n", state.NextPlayer.other(), state.LastMove.String())
		state.Board.print(true)
	} else {
		fmt.Println("Empty Board")
	}
}

func (move Move) String() string {

	if move.IsPass {
		return "passes"
	}
	if move.IsResign {
		return "resigns"
	}

	return fmt.Sprintf("plays [%d,%d]", move.Point.Row, move.Point.Col)
}

func (gs GoString) print() {
	libCoords := gs.stoneCoords(gs.Liberties)
	stoneCoords := gs.stoneCoords(gs.Stones)
	fmt.Printf("{Player: %s, Liberties (%d): %s, Stones (%d): %s}\n", gs.Player, len(libCoords), strings.Join(libCoords, ", "), len(stoneCoords), strings.Join(stoneCoords, ", "))
}

func (gs GoString) stoneCoords(input map[int]map[int]bool) []string {
	var stoneStr []string
	for row, cols := range input {
		for col, _ := range cols {
			stoneStr = append(stoneStr, fmt.Sprintf("[%d,%d]", row, col))
		}
	}

	return stoneStr
}

func (p Player) String() string {
	if p.isBlack {
		return "Black"
	} else {
		return "White"
	}
}
