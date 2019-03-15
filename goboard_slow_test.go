package main

import (
	"testing"
	"reflect"
	"fmt"
)

func TestGoStringMergeStrings(t *testing.T) {

	string1 := map[int]map[int]bool{4: {3: true}}
	string2 := map[int]map[int]bool{4: {4: true}, 5: {4: true}}
	string3 := map[int]map[int]bool{4: {3: true, 4: true}, 5: {4: true}}

	gs1 := GoString{}
	expected3 := gs1.union(string1, string2)

	if !reflect.DeepEqual(expected3, string3) {
		t.Fatalf("Expected %v, but got %v", expected3, string3)
	}

	expected1 := gs1.diff(string3, string2)
	if !reflect.DeepEqual(expected1, string1) {
		t.Fatalf("Expected %v, but got %v", expected1, string1)
	}

	expected2 := gs1.diff(string3, string1)
	if !reflect.DeepEqual(expected2, string2) {
		t.Fatalf("Expected %v, but got %v", expected2, string2)
	}
}

func TestBoard_PlaceStone(t *testing.T) {
	board := Board{
		Grid:    []GoString{},
		NumRows: 5,
		NumCols: 7,
	}

	// As figure 3.2
	board.PlaceStone(PlayerBlack, Point{3, 6})
	board.PlaceStone(PlayerBlack, Point{4, 5})
	board.PlaceStone(PlayerWhite, Point{2, 3})
	board.PlaceStone(PlayerBlack, Point{3, 3})
	board.PlaceStone(PlayerWhite, Point{2, 4})
	board.PlaceStone(PlayerBlack, Point{3, 4}) // TODO This should merge??
	board.PlaceStone(PlayerWhite, Point{3, 2})
	board.PlaceStone(PlayerBlack, Point{4, 3})
	board.PlaceStone(PlayerWhite, Point{4, 2})
	board.PlaceStone(PlayerBlack, Point{4, 5})
	board.PlaceStone(PlayerWhite, Point{3, 5})
	board.PlaceStone(PlayerBlack, Point{4, 6})
	board.PlaceStone(PlayerWhite, Point{4, 4})
	board.PlaceStone(PlayerBlack, Point{5, 4})

	if _, exists := board.get(Point{4, 4}); exists {
		t.Fatal("Did not expect stone in 4/4")
	}
}

func TestGameState(t *testing.T) {
	state := NewGame(3)

	state.
		applyOrFail(t, PlayerBlack, Play(Point{1, 1})).
		applyOrFail(t, PlayerWhite, Play(Point{1, 2})).
		applyOrFail(t, PlayerBlack, Play(Point{2, 1})).
		applyOrFail(t, PlayerWhite, Play(Point{2, 2})).
		applyOrFail(t, PlayerBlack, Play(Point{1, 3})).
		applyOrFail(t, PlayerWhite, Play(Point{3, 1})).
		applyOrFail(t, PlayerBlack, Pass()).
		applyOrFail(t, PlayerWhite, Pass())
}

func TestGameStateSelfCapture(t *testing.T) {
	state := NewGame(3)

	state = state.
		applyOrFail(t, PlayerBlack, Play(Point{1, 2})).
		applyOrFail(t, PlayerWhite, Play(Point{3, 1})).
		applyOrFail(t, PlayerBlack, Play(Point{2, 1}))

	if !state.IsMoveSelfCapture(PlayerWhite, Play(Point{1, 1})) {
		t.Fatalf("Did not detect self capture")
	}
}

func TestGameStateViolateKo(t *testing.T) {
	state := NewGame(6)

	state = state.
		applyOrFail(t, PlayerBlack, Play(Point{2, 3})).
		applyOrFail(t, PlayerWhite, Play(Point{2, 4})).
		applyOrFail(t, PlayerBlack, Play(Point{3, 2})).
		applyOrFail(t, PlayerWhite, Play(Point{3, 5})).
		applyOrFail(t, PlayerBlack, Play(Point{4, 3})).
		applyOrFail(t, PlayerWhite, Play(Point{4, 4})).
		applyOrFail(t, PlayerBlack, Play(Point{3, 4})).
		applyOrFail(t, PlayerWhite, Play(Point{3, 3}))

	if !state.DoesMoveViolateKo(PlayerBlack, Play(Point{3, 4})) {
		t.Fatalf("Did not detect KO")
	}
}

func (state GameState) applyOrFail(t *testing.T, player Player, move Move) GameState {
	state, err := state.ApplyMove(player, move)
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}

	state.Print()
	fmt.Println("")
	return state
}

func BenchmarkRandomBotSelectMove(b *testing.B) {
	boardSize := 19

	players := map[Player]Agent{
		PlayerBlack: &RandomBot{},
		PlayerWhite: &RandomBot{},
	}

	game := NewGame(boardSize)

	for n := 0; n < b.N; n++ {
		players[PlayerBlack].selectMove(game)
	}
}

func BenchmarkRandomBots(b *testing.B) {
	for n := 0; n < b.N; n++ {
		boardSize := 19

		players := map[Player]Agent{
			PlayerBlack: &RandomBot{},
			PlayerWhite: &RandomBot{},
		}

		game := NewGame(boardSize)
		var err error
		moveNr := 1
		for !game.IsOver() {
			player := game.NextPlayer
			nextMove := players[player].selectMove(game)
			game, err = game.ApplyMove(player, nextMove)
			moveNr++

			if err != nil {
				b.Errorf("Error during move! %s\n", err)
				break
			}
		}
	}
}
