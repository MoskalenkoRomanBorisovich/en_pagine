package main

import (
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"
)

type Position uint8

func MakePos[T constraints.Integer](row, col T) Position {
	return Position((row << 3) + col)
}

func (pos Position) GetCol() int8 {
	return int8(pos & 0b111)
}

func (pos Position) SetCol(col uint8) Position {
	return pos&(0b111_000) | Position(col)
}

func (pos Position) GetRow() int8 {
	return int8((pos >> 3) & 0b111)
}

func (pos Position) SetRow(row uint8) Position {
	return pos&(0b000_111) | Position(row<<3)
}

func (pos Position) id() uint8 {
	return uint8(pos)
}

func (pos Position) String() string {
	return fmt.Sprintf("%c%d", 'a'+pos.GetCol(), pos.GetRow()+1)
}

func MakePiecePosFromFEN(fen string) (Position, error) {
	er := errors.New("incorrect position format")
	if len(fen) != 2 {
		return Position(0), er
	}
	c := rune(fen[0]) - 'a'
	if c < 0 || 7 < c {
		return Position(0), er
	}
	r := rune(fen[1]) - '1'
	if r < 0 || 7 < r {
		return Position(0), er
	}
	return MakePos(r, c), nil
}
