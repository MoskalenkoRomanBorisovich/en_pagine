package main

type Move uint16

const (
	movePromoteMask = 0b0111_000000_000000 // piece to which the pawn is promoted to
	moveStartMask   = 0b0000_111111_000000
	moveEndMask     = 0b0000_000000_111111
)

func (m Move) SetPromote(p Piece) Move {
	return (m &^ movePromoteMask) | (Move(p) << 12)
}
func (m Move) GetPromote() Piece {
	return Piece((m & movePromoteMask) >> 12)
}

func (m Move) SetStart(p Position) Move {
	return (m &^ moveStartMask) | (Move(p) << 6)
}
func (m Move) GetStart() Position {
	return Position((m & moveStartMask) >> 6)
}

func (m Move) SetEnd(p Position) Move {
	return (m &^ moveEndMask) | Move(p)
}

func (m Move) GetEnd() Position {
	return Position(m & moveEndMask)
}
