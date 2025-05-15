package main

import (
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"
)

type PiecePos uint8

func MakePiecePos[T constraints.Integer](row, col T) PiecePos {
	return PiecePos((row << 3) + col)
}

func (pos PiecePos) GetCol() uint8 {
	return uint8(pos & 0b111)
}

func (pos PiecePos) SetCol(col uint8) PiecePos {
	return pos&(0b111_000) | PiecePos(col)
}

func (pos PiecePos) GetRow() uint8 {
	return uint8((pos >> 3) & 0b111)
}

func (pos PiecePos) SetRow(row uint8) PiecePos {
	return pos&(0b000_111) | PiecePos(row<<3)
}

func (pos PiecePos) id() uint8 {
	return uint8(pos)
}

func (pos PiecePos) String() string {
	return fmt.Sprintf("%c%d", 'a'+pos.GetCol(), pos.GetRow()+1)
}

func MakePiecePosFromFEN(fen string) (PiecePos, error) {
	er := errors.New("incorrect position format")
	if len(fen) != 2 {
		return PiecePos(0), er
	}
	c := rune(fen[0]) - 'a'
	if c < 0 || 7 < c {
		return PiecePos(0), er
	}
	r := rune(fen[1]) - '1'
	if r < 0 || 7 < r {
		return PiecePos(0), er
	}
	return MakePiecePos(r, c), nil
}
