package main

import "fmt"

type Piece byte

const (
	NoPiece Piece = 0
	Pawn    Piece = 1
	Knight  Piece = 2
	Bishop  Piece = 3
	Rook    Piece = 4
	Queen   Piece = 5
	King    Piece = 6

	White Piece = 8

	B_Pawn   = Pawn
	B_Knight = Knight
	B_Bishop = Bishop
	B_Rook   = Rook
	B_Queen  = Queen
	B_King   = King

	W_Pawn   = Pawn | White
	W_Knight = Knight | White
	W_Bishop = Bishop | White
	W_Rook   = Rook | White
	W_Queen  = Queen | White
	W_King   = King | White
)

type pieceLiteral struct {
	piece   Piece
	literal rune
}

var pieceLieterals = [...]pieceLiteral{
	{NoPiece, ' '},

	{B_Pawn, 'p'},
	{B_Knight, 'n'},
	{B_Bishop, 'b'},
	{B_Rook, 'r'},
	{B_Queen, 'q'},
	{B_King, 'k'},

	{W_Pawn, 'P'},
	{W_Knight, 'N'},
	{W_Bishop, 'B'},
	{W_Rook, 'R'},
	{W_Queen, 'Q'},
	{W_King, 'K'},
}

func MakePieceFromRune(r rune) (Piece, error) {
	for _, pl := range pieceLieterals {
		if pl.literal == r {
			return pl.piece, nil
		}
	}
	return NoPiece, &ErrorUnknownPieceLiteral{r}
}

func (p Piece) GetLiteral() (rune, error) {
	for _, pl := range pieceLieterals {
		if pl.piece == p {
			return pl.literal, nil
		}
	}
	return ' ', &ErrorUnknownPiece{p}
}

func (p Piece) String() string {
	l, er := p.GetLiteral()
	if er != nil {
		return "?"
	}
	return string(l)
}

func (p Piece) IsWhite() bool {
	return p&White != 0
}

func (p Piece) IsBlack() bool {
	return !p.IsWhite()
}

func (p Piece) GetType() Piece {
	return p & 7
}

type ErrorUnknownPiece struct {
	p Piece
}

func (e *ErrorUnknownPiece) Error() string {
	return fmt.Sprintf("unknown piece type: %d", e)
}

type ErrorUnknownPieceLiteral struct {
	literal rune
}

func (e *ErrorUnknownPieceLiteral) Error() string {
	return fmt.Sprintf("unknown piece literal: %c", e.literal)
}
