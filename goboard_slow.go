package main

import (
	"fmt"
	"reflect"
	"errors"
)

type Move struct {
	Point    Point
	IsPass   bool
	IsResign bool
	IsPlay   bool
}

func Play(point Point) Move {
	return Move{
		Point:  point,
		IsPlay: true,
	}
}

func Pass() Move {
	return Move{
		IsPass: true,
	}
}

func Resign() Move {
	return Move{
		IsResign: true,
	}
}

type GoString struct {
	Player    Player
	Stones    map[int]map[int]bool
	Liberties map[int]map[int]bool
}

func (gs GoString) Copy() GoString {
	newPlayer := gs.Player

	newLiberties := make(map[int]map[int]bool)
	for row, cols := range gs.Liberties {
		newLibertiesRow := make(map[int]bool)
		for col, val := range cols {
			newLibertiesRow[col] = val
		}
		newLiberties[row] = newLibertiesRow
	}

	newStones := make(map[int]map[int]bool)
	for row, cols := range gs.Stones {
		newStonesRow := make(map[int]bool)
		for col, val := range cols {
			newStonesRow[col] = val
		}
		newStones[row] = newStonesRow
	}

	return GoString{
		Player:    newPlayer,
		Liberties: newLiberties,
		Stones:    newStones,
	}
}

func (gs GoString) removeLiberty(point Point) {
	delete(gs.Liberties[point.Row], point.Col)

	if len(gs.Liberties[point.Row]) == 0 {
		delete(gs.Liberties, point.Row)
	}
}

func (gs GoString) withoutLiberty(point Point) GoString {
	newGs := gs.Copy()

	delete(newGs.Liberties[point.Row], point.Col)

	if len(newGs.Liberties[point.Row]) == 0 {
		delete(newGs.Liberties, point.Row)
	}

	return newGs
}

func (gs GoString) addLiberty(point Point) {
	if _, exists := gs.Liberties[point.Row]; !exists {
		gs.Liberties[point.Row] = make(map[int]bool)
	}

	gs.Liberties[point.Row][point.Col] = true
}

func (gs GoString) withLiberty(point Point) GoString {
	newGs := gs.Copy()

	if _, exists := newGs.Liberties[point.Row]; !exists {
		newGs.Liberties[point.Row] = make(map[int]bool)
	}

	newGs.Liberties[point.Row][point.Col] = true

	return newGs
}

func (gs GoString) remove(set map[int]map[int]bool, row, col int) {
	delete(set[row], col)

	if len(set[row]) == 0 {
		delete(set, row)
	}
}

func (gs GoString) add(set map[int]map[int]bool, row, col int) {
	if _, exists := set[row]; !exists {
		set[row] = make(map[int]bool)
	}

	set[row][col] = true
}

func (gs GoString) union(a, b map[int]map[int]bool) map[int]map[int]bool {

	newset := make(map[int]map[int]bool)

	for _, rows := range []map[int]map[int]bool{a, b} {
		for row, cols := range rows {
			for col, _ := range cols {
				gs.add(newset, row, col)
			}
		}
	}

	return newset
}

func (gs GoString) diff(a, b map[int]map[int]bool) map[int]map[int]bool {

	newset := make(map[int]map[int]bool)

	for row, cols := range a {
		for col, _ := range cols {
			if _, exists := b[row][col]; !exists {
				gs.add(newset, row, col)
			}
		}
	}

	return newset
}

func (gs GoString) mergeStrings(other GoString) GoString {
	unionStones := gs.union(gs.Stones, other.Stones)
	unionLiberties := gs.union(gs.Liberties, other.Liberties)
	unionLiberties = gs.diff(unionLiberties, unionStones)
	return GoString{
		Player:    gs.Player,
		Stones:    unionStones,
		Liberties: unionLiberties,
	}
}

func (gs GoString) equals(other GoString) bool {
	return gs.Player == other.Player && reflect.DeepEqual(gs.Stones, other.Stones)
}

type Board struct {
	NumRows int
	NumCols int
	Grid    []GoString
	Hash    uint64
}

func (b *Board) Copy() Board {
	var newGrid []GoString
	for _, gs := range b.Grid {
		newGrid = append(newGrid, gs.Copy())
	}

	return Board{
		NumRows: b.NumRows,
		NumCols: b.NumCols,
		Grid:    newGrid,
		Hash:    b.Hash,
	}
}

func (b *Board) PlaceStone(player Player, point Point) error {
	if !b.isOnGrid(point) {
		return fmt.Errorf("Point %v not on Grid", point)
	}

	if _, exists := b.get(point); exists {
		return fmt.Errorf("Point %v already played on Grid", point)
	}

	var adjacentSameColor []GoString
	var adjacentOtherColor []GoString
	liberties := make(map[int]map[int]bool)

	// Loop over all possible neighbors
	for _, neighbor := range point.Neighbors() {
		if !b.isOnGrid(neighbor) {
			continue;
		}

		neighborString, neighborStringExists := b.getGoString(neighbor)
		if !neighborStringExists {
			if _, exists := liberties[neighbor.Row]; !exists {
				liberties[neighbor.Row] = make(map[int]bool)
			}

			liberties[neighbor.Row][neighbor.Col] = true
		} else if neighborString.Player.equals(player) {
			// When the neighbor string is of the current Player
			adjacentSameColor = b.appendGoStringUnique(adjacentSameColor, neighborString)
		} else {
			// When the neighbor string is of the other Player
			adjacentOtherColor = b.appendGoStringUnique(adjacentOtherColor, neighborString)
		}
	}

	newString := GoString{Player: player, Stones: map[int]map[int]bool{point.Row: {point.Col: true}}, Liberties: liberties}

	// Merge any adjacent strings of the same color
	for _, gs := range adjacentSameColor {
		newString = newString.mergeStrings(gs)

		// Remove the old string from the Grid
		var removeIndex int
		for gridIndex, gridString := range b.Grid {
			if gridString.equals(gs) {
				removeIndex = gridIndex
				break
			}
		}

		b.Grid = b.removeGoString(b.Grid, removeIndex)
	}

	b.Grid = append(b.Grid, newString)
	b.Hash ^= hashes[player.isBlack][point.Row][point.Col]

	// Reduce Liberties of any adjacent strings of the opposite color
	for _, gs := range adjacentOtherColor {
		gs.removeLiberty(point)
	}

	// If any opposite color strings now have zero Liberties, remove them
	for _, gs := range adjacentOtherColor {
		if len(gs.Liberties) == 0 {
			b.removeString(gs)
		}
	}

	return nil
}

/**
We need to keep in mind that removing a string results in Liberties
 */
func (b *Board) removeString(gs GoString) {
	// Increase Liberties of neighbors
	for row, cols := range gs.Stones {
		for col, _ := range cols {
			point := Point{row, col}
			for _, neighbor := range point.Neighbors() {
				neighborStr, exists := b.getGoString(neighbor);
				if !exists {
					continue
				}

				if !neighborStr.equals(gs) {
					neighborStr.addLiberty(point)
				}
			}
		}
	}

	// Remove string from Grid
	for i, other := range b.Grid {
		if other.equals(gs) {

			for row, cols := range other.Stones {
				for col, _ := range cols {
					b.Hash ^= hashes[other.Player.isBlack][row][col]
				}
			}

			b.Grid = b.removeGoString(b.Grid, i)
			break
		}
	}
}

func (b *Board) removeGoString(s []GoString, i int) []GoString {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (b *Board) appendGoStringUnique(list []GoString, gs GoString) []GoString {
	for _, other := range list {
		if gs.equals(other) {
			return list
		}
	}

	return append(list, gs)
}

func (b *Board) isOnGrid(point Point) bool {
	return 1 <= point.Row && point.Row <= b.NumRows && 1 <= point.Col && point.Col <= b.NumCols
}

func (b *Board) get(point Point) (player Player, exists bool) {

	for _, gs := range b.Grid {
		if _, exists := gs.Stones[point.Row][point.Col]; exists {
			return gs.Player, true
		}
	}

	return Player{}, false
}

func (b *Board) getGoString(point Point) (gs GoString, exists bool) {

	for _, gs := range b.Grid {
		if _, exists := gs.Stones[point.Row][point.Col]; exists {
			return gs, true
		}
	}

	return GoString{}, false
}

type GameState struct {
	Board          Board
	NextPlayer     Player
	PreviousState  *GameState
	PreviousStates []Situation
	LastMove       Move
}

func (state GameState) Copy() GameState {
	newBoard := state.Board.Copy()
	newPlayer := state.NextPlayer
	newMove := state.LastMove

	return GameState{
		Board:          newBoard,
		NextPlayer:     newPlayer,
		PreviousState:  state.PreviousState,
		PreviousStates: state.PreviousStates,
		LastMove:       newMove,
	}
}

func (state GameState) ApplyMove(player Player, move Move) (GameState, error) {
	if player != state.NextPlayer {
		return GameState{}, fmt.Errorf("Expected other Player move")
	}

	if err := state.IsMoveValid2(move); err != nil {
		return GameState{}, err
	}

	nextBoard := state.Board.Copy()

	prevSituation := Situation{state.NextPlayer, state.Board.Hash}

	if move.IsPlay {
		nextBoard.PlaceStone(player, move.Point)
	}

	newPrevSituations := append(state.PreviousStates, prevSituation)
	return GameState{Board: nextBoard, NextPlayer: player.other(), PreviousState: &state, PreviousStates: newPrevSituations, LastMove: move}, nil
}

func (state GameState) LegalMoves() []Point {
	var legalMoves []Point

	for row := 1; row <= state.Board.NumRows; row++ {
		for col := 1; col <= state.Board.NumCols; col++ {
			candidate := Point{row, col}

			if state.IsMoveValid(Play(candidate)) && !state.Board.IsPointAnEye(candidate, state.NextPlayer) {
				legalMoves = append(legalMoves, candidate)
			}
		}
	}

	return legalMoves
}

func NewGame(boardSize int) GameState {
	board := Board{NumRows: boardSize, NumCols: boardSize, Grid: []GoString{}}
	return GameState{
		Board:          board,
		NextPlayer:     PlayerBlack,
		PreviousStates: []Situation{},
	}
}

func (state GameState) IsOver() bool {
	lastMove := state.LastMove
	if &lastMove == nil {
		return false
	}

	if state.PreviousState == nil {
		return false
	}

	if lastMove.IsResign {
		return true
	}

	secondLastMove := state.PreviousState.LastMove
	if &secondLastMove == nil {
		return false
	}

	return lastMove.IsPass && secondLastMove.IsPass
}

func (state GameState) IsMoveSelfCapture(player Player, move Move) bool {
	if !move.IsPlay {
		return false
	}

	// TODO Seems that this method alters something deep down gamestate!! (at least removes liberties?)
	// maybe it happens in placestone and not here

	// Dry-Attempt to play it and see if the stone has any liberties
	nextBoard := state.Board.Copy()
	nextBoard.PlaceStone(player, move.Point)
	newString, _ := nextBoard.getGoString(move.Point)
	return len(newString.Liberties) == 0
}

func (state GameState) DoesMoveViolateKo(player Player, move Move) bool {
	if !move.IsPlay {
		return false
	}

	nextBoard := state.Board.Copy()
	nextBoard.PlaceStone(player, move.Point)

	nextSituation := Situation{NextPlayer: player.other(), Hash: nextBoard.Hash}

	for _, pastState := range state.PreviousStates {
		if pastState.NextPlayer == nextSituation.NextPlayer && pastState.Hash == nextSituation.Hash {
			return true
		}
	}

	return false
}

func (state GameState) situation() Situation {
	return Situation{NextPlayer: state.NextPlayer, Hash: state.Board.Hash}
}

func (state GameState) IsMoveValid(move Move) bool {
	if state.IsOver() {
		return false
	}

	if move.IsPass || move.IsResign {
		return true
	}

	_, exists := state.Board.get(move.Point)
	if exists {
		return false
	}

	return !state.IsMoveSelfCapture(state.NextPlayer, move) && !state.DoesMoveViolateKo(state.NextPlayer, move)
}

func (state GameState) IsMoveValid2(move Move) error {
	if state.IsOver() {
		return errors.New("Game is over")
	}

	if move.IsPass || move.IsResign {
		return nil
	}

	_, exists := state.Board.get(move.Point)
	if exists {
		return errors.New("Already a stone in place")
	}

	if state.IsMoveSelfCapture(state.NextPlayer, move) {
		return errors.New("Move is self capture")
	}

	if state.DoesMoveViolateKo(state.NextPlayer, move) {
		return errors.New("Move violates KO")
	}

	return nil
}

func (state GameState) Winner() *Player {
	if !state.IsOver() {
		return nil
	}

	if state.LastMove.IsResign {
		return &state.NextPlayer
	}

	// TODO Very incorrect scoring!! :D
	if state.Board.CaptureDiff(PlayerBlack) >= 0 {
		return &PlayerBlack
	} else {
		return &PlayerWhite;
	}
}

type Situation struct {
	NextPlayer Player
	Hash       uint64
}
