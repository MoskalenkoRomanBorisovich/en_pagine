package main

import (
	"errors"
	"strconv"
	"strings"
)

type BoardState uint32 // turn, castle, en passant, moves since capture or pawn advance

const (
	bsTurnMask    BoardState = 0b_100000_000000_000000
	bsKMask       BoardState = 0b_010000_000000_000000
	bsQMask       BoardState = 0b_001000_000000_000000
	bskMask       BoardState = 0b_000100_000000_000000
	bsqMask       BoardState = 0b_000010_000000_000000
	bsIsEnPosMask BoardState = 0b_000001_000000_000000
	bsEnPosMask   BoardState = 0b_000000_111111_000000
	bsHMovesMask  BoardState = 0b_000000_000000_111111
)

func (bs BoardState) Set_Turn(is_white bool) BoardState {
	if is_white {
		return bs | bsTurnMask
	}
	return bs &^ bsTurnMask
}

func (bs BoardState) Get_Turn() (is_white bool) {
	return bs&bsTurnMask != 0
}

func (bs BoardState) Set_K(is_K bool) BoardState {
	if is_K {
		return bs | bsKMask
	}
	return bs &^ bsKMask
}

func (bs BoardState) Get_K() bool {
	return bs&bsKMask != 0
}

func (bs BoardState) Set_Q(is_Q bool) BoardState {
	if is_Q {
		return bs | bsQMask
	}
	return bs &^ bsQMask
}

func (bs BoardState) Get_Q() bool {
	return bs&bsQMask != 0
}

func (bs BoardState) Set_k(is_k bool) BoardState {
	if is_k {
		return bs | bskMask
	}
	return bs &^ bskMask
}

func (bs BoardState) Get_k() bool {
	return bs&bskMask != 0
}

func (bs BoardState) Set_q(is_q bool) BoardState {
	if is_q {
		return bs | bsqMask
	}
	return bs &^ bsqMask
}

func (bs BoardState) Get_q() bool {
	return bs&bsqMask != 0
}

func (bs BoardState) Set_IsEnPos(is_en_pos bool) BoardState {
	if is_en_pos {
		return bs | bsIsEnPosMask
	}
	return bs &^ bsIsEnPosMask
}

func (bs BoardState) Get_IsEnPos() bool {
	return bs&bsIsEnPosMask != 0
}

func (bs BoardState) Set_EnPos(pos PiecePos) BoardState {
	return (bs &^ bsEnPosMask) | BoardState(pos)<<6
}

func (bs BoardState) Get_EnPos() PiecePos {
	return PiecePos(bs & bsEnPosMask >> 6)
}

func (bs BoardState) Set_HMoves(h_moves uint8) BoardState {
	return (bs &^ bsHMovesMask) | BoardState(h_moves)
}

func (bs BoardState) Get_HMoves() uint8 {
	return uint8(bs & bsHMovesMask)
}

func MakeInitialBoardState() BoardState {
	var res BoardState
	res = res.Set_Turn(true)
	res = res.Set_K(true)
	res = res.Set_Q(true)
	res = res.Set_k(true)
	res = res.Set_q(true)
	res = res.Set_IsEnPos(false)
	// no need to set ep passant position
	res = res.Set_HMoves(1)
	return res
}

func (bs BoardState) FEN() string {
	var res string
	if bs.Get_Turn() {
		res += "w"
	} else {
		res += "b"
	}
	res += " "

	var is_castle = false
	if bs.Get_K() {
		res += "K"
		is_castle = true
	}
	if bs.Get_Q() {
		res += "Q"
		is_castle = true
	}
	if bs.Get_k() {
		res += "k"
		is_castle = true
	}
	if bs.Get_q() {
		res += "q"
		is_castle = true
	}

	if !is_castle {
		res += "-"
	}
	res += " "

	if bs.Get_IsEnPos() {
		pos := bs.Get_EnPos()
		res += pos.String()
	} else {
		res += "-"
	}
	res += " "

	res += strconv.Itoa(int(bs.Get_HMoves()))

	return res
}

func (bs BoardState) setTurnFromFen(fen string) (BoardState, error) {
	switch fen {
	case "w":
		return bs.Set_Turn(true), nil
	case "b":
		return bs.Set_Turn(false), nil
	}
	return bs, errors.New("unknown turn type")
}

func (bs BoardState) setCastleFromFen(fen string) (BoardState, error) {
	bs = bs.Set_K(false)
	bs = bs.Set_Q(false)
	bs = bs.Set_k(false)
	bs = bs.Set_q(false)
	if fen == "-" {
		return bs, nil
	}
	error_rcv := errors.New("repeated castle value")
	for _, c := range fen {
		switch c {
		case 'K':
			if bs.Get_K() {
				return bs, error_rcv
			}
			bs = bs.Set_K(true)
		case 'Q':
			if bs.Get_Q() {
				return bs, error_rcv
			}
			bs = bs.Set_Q(true)
		case 'k':
			if bs.Get_k() {
				return bs, error_rcv
			}
			bs = bs.Set_k(true)
		case 'q':
			if bs.Get_q() {
				return bs, error_rcv
			}
			bs = bs.Set_q(true)
		}
	}
	return bs, nil
}

func (bs BoardState) setEnPosFromFen(fen string) (BoardState, error) {
	if fen == "-" {
		bs = bs.Set_EnPos(MakePiecePos(0, 0)) // for consistency
		return bs.Set_IsEnPos(false), nil
	}
	bs = bs.Set_IsEnPos(true)
	pos, er := MakePiecePosFromFEN(fen)
	if er != nil {
		return bs, er
	}
	return bs.Set_EnPos(pos), nil
}

func (bs BoardState) setHMovesFromFen(fen string) (BoardState, error) {
	hmoves, er := strconv.Atoi(fen)
	if er != nil {
		return bs, er
	}
	if hmoves < 1 || hmoves > 50 {
		return bs, errors.New("Impossible number of half moves")
	}
	return bs.Set_HMoves(uint8(hmoves)), nil
}

func MakeBoardStateFromFEN(fen string) (BoardState, error) {
	var res BoardState
	fen_params := strings.Split(fen, " ")
	if len(fen_params) != 4 {
		return res, errors.New("wrong number of parameters")
	}
	var er error

	res, er = res.setTurnFromFen(fen_params[0])
	if er != nil {
		return res, er
	}

	res, er = res.setCastleFromFen(fen_params[1])
	if er != nil {
		return res, er
	}

	res, er = res.setEnPosFromFen(fen_params[2])
	if er != nil {
		return res, er
	}

	res, er = res.setHMovesFromFen(fen_params[3])
	if er != nil {
		return res, er
	}

	return res, nil
}
