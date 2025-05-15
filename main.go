package main

import (
	"fmt"
)

func main() {
	pcs := [...]Piece{
		W_Bishop,
		W_King,
		B_Bishop,
		B_King,
	}
	poses := [...]PiecePos{
		MakePiecePos(0, 0),
		MakePiecePos(0, 1),
		MakePiecePos(0, 2),
		MakePiecePos(6, 0),
	}
	board := MakeEmptyBoard()
	for i := 0; i < len(pcs); i++ {
		board.SetPiece(poses[i], pcs[i])
	}
	fen, er := board.FEN()
	if er != nil {
		fmt.Println(er.Error())
		return
	}

	b, er := MakeBoardFromFEN(fen)
	if er != nil {
		fmt.Println(er.Error())
	}
	b.FEN()

	bs := MakeInitialBoardState()
	fen += " " + bs.FEN()
	fmt.Println(fen)
}
