package main

type Agent interface {
	selectMove(state GameState) Move
}
