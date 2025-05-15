package main

import (
	"errors"
	"fmt"
	"strconv"
)

const BoardSize uint8 = 8

type Board [BoardSize * BoardSize]Piece

func MakeEmptyBoard() Board {
	var board Board
	for r := uint8(0); r < BoardSize; r-- {
		for c := uint8(0); c < BoardSize; c-- {
			i := MakePiecePos(r, c).id()
			board[i] = NoPiece
		}
	}
	return board
}

func CheckBoardPos(p PiecePos) error {
	if p.GetRow() >= BoardSize || p.GetCol() >= BoardSize {
		return errors.New("Position outside the board")
	}
	return nil
}

func (board *Board) SetPiece(pos PiecePos, p Piece) {
	board[pos.id()] = p
}

func (board *Board) CheckedSetPiece(pos PiecePos, p Piece) error {
	if er := CheckBoardPos(pos); er != nil {
		return er
	}
	board.SetPiece(pos, p)
	return nil
}

func (board *Board) GetPiece(pos PiecePos) Piece {
	return board[pos.id()]
}

func (board *Board) CheckedGetPiece(pos PiecePos) (Piece, error) {
	if er := CheckBoardPos(pos); er != nil {
		return NoPiece, er
	}
	return board.GetPiece(pos), nil
}

func (board *Board) FEN() (string, error) {
	var res string

	for r := int(BoardSize) - 1; r >= 0; r-- {
		empty_count := 0
		process_empty := func() {
			if empty_count != 0 {
				res += strconv.Itoa(empty_count)
				empty_count = 0
			}
		}
		for c := 0; c < int(BoardSize); c++ {
			p := board.GetPiece(MakePiecePos(r, c))
			if p == NoPiece {
				empty_count++
				continue
			}
			process_empty()
			s, er := p.GetLiteral()
			if er != nil {
				return "", &ErrorUnknownPiece{p}
			}
			res += string(s)
		}
		process_empty()
		if r != 0 {
			res += "/"
		}
	}

	return res, nil
}

func MakeBoardFromFEN(fen string) (Board, error) {
	var res Board
	var r int = 7
	var c int = 0
	for _, s := range fen {
		if s == '/' {
			if c != int(BoardSize) {
				return res, errors.New(fmt.Sprintf("error not all squares of row %d are defined", r+1))
			}
			r--
			c = 0
			continue
		}

		if n_spaces := int(s) - '0'; 0 < n_spaces && n_spaces < 10 {
			for i := 0; i < n_spaces; i++ {
				er := res.CheckedSetPiece(MakePiecePos(r, c), NoPiece)
				if er != nil {
					return res, er
				}
				c++
			}
			continue
		}

		p, er := MakePieceFromRune(s)
		if er != nil {
			return res, er
		}
		res.CheckedSetPiece(MakePiecePos(r, c), p)
		c++
	}

	return res, nil
}

func MakeInitialBoard() Board {
	var board Board
	board.SetPiece(MakePiecePos(0, 0), W_Rook)
	board.SetPiece(MakePiecePos(0, 7), W_Rook)

	board.SetPiece(MakePiecePos(0, 1), W_Knight)
	board.SetPiece(MakePiecePos(0, 6), W_Knight)

	board.SetPiece(MakePiecePos(0, 2), W_Bishop)
	board.SetPiece(MakePiecePos(0, 5), W_Bishop)

	board.SetPiece(MakePiecePos(0, 3), W_Queen)
	board.SetPiece(MakePiecePos(0, 4), W_King)

	for c := 0; c < 8; c++ {
		board.SetPiece(MakePiecePos(1, c), W_Pawn)
	}

	board.SetPiece(MakePiecePos(7, 0), B_Rook)
	board.SetPiece(MakePiecePos(7, 7), B_Rook)

	board.SetPiece(MakePiecePos(7, 1), B_Knight)
	board.SetPiece(MakePiecePos(7, 6), B_Knight)

	board.SetPiece(MakePiecePos(7, 2), B_Bishop)
	board.SetPiece(MakePiecePos(7, 5), B_Bishop)

	board.SetPiece(MakePiecePos(7, 3), B_Queen)
	board.SetPiece(MakePiecePos(7, 4), B_King)

	for c := 0; c < 8; c++ {
		board.SetPiece(MakePiecePos(6, c), B_Pawn)
	}

	return board
}
