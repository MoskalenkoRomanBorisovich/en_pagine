package main

import "testing"

func TestBoard(t *testing.T) {
	board := MakeInitialBoard()
	const initial_state string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	fen, er := board.FEN()
	assert_er(er, t)
	assert_equal(fen, initial_state, t)

	board_copy, er := MakeBoardFromFEN(fen)
	assert_er(er, t)
	assert_equal(board, board_copy, t)
}
