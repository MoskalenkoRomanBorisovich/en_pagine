package main

import "golang.org/x/exp/constraints"

func UpdateCastle(pos Position, bs *BoardState) {
	switch pos {
	// rooks move
	case MakePos(0, 0):
		*bs = bs.Set_Q(false)
	case MakePos(0, 7):
		*bs = bs.Set_K(false)
	case MakePos(7, 0):
		*bs = bs.Set_q(false)
	case MakePos(7, 7):
		*bs = bs.Set_k(false)
	// kings move
	case MakePos(0, 4):
		*bs = bs.Set_K(false)
		*bs = bs.Set_Q(false)
	case MakePos(7, 4):
		*bs = bs.Set_k(false)
		*bs = bs.Set_q(false)
	}
}

func MakeMove(move Move, board *Board, bs *BoardState) {
	start := move.GetStart()
	end := move.GetEnd()
	UpdateCastle(start, bs)
	UpdateCastle(end, bs)
	piece := board.GetPiece(start)
	switch piece.GetType() {
	case Knight:
		fallthrough
	case Bishop:
		fallthrough
	case Rook:
		fallthrough
	case Queen:
		board.SetPiece(start, NoPiece)
		board.SetPiece(end, piece)
	case Pawn:
		// TODO
	case King:
		// TODO
	}

}

func isTakenByFriend(board *Board, pos Position, is_white bool) bool {
	piece := board.GetPiece(pos)
	return piece != NoPiece && piece.IsWhite() == is_white
}

func isTaken(board *Board, pos Position) bool {
	return board.GetPiece(pos) != NoPiece
}

func Abs[T constraints.Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// return true if position is under attack, is_white - determines color of the piece at this position
func IsPosUnderAttack(pos Position, board *Board, is_white bool) bool {
	// TODO
	return false
}

func IsPawnMovePossible(move Move, board *Board, is_white bool, is_en bool, en_pos Position) bool {
	start := move.GetStart()
	end := move.GetEnd()
	rs := start.GetRow()
	cs := start.GetCol()
	re := end.GetRow()
	ce := end.GetCol()
	var dir int8
	if is_white {
		dir = 1
	} else {
		dir = -1
	}
	switch {
	case ce == cs && re == rs+dir: // move forward 1
		if board.GetPiece(end) != NoPiece {
			return false
		}
		return true
	case ce == cs+1 || ce == cs-1 && re == rs+dir: // capture
		end_p := board.GetPiece(end)
		if end_p != NoPiece {
			if end_p.IsWhite() == is_white {
				return false
			}
		} else {
			// en passant
			if !is_en {
				return false
			}
			if en_pos != end {
				return false
			}
		}
	case ce == cs && re == rs+dir*2: // move forward 2
		if board.GetPiece(end) != NoPiece {
			return false
		}
		if board.GetPiece(MakePos(rs+dir, cs)) != NoPiece {
			return false
		}
	default:
		return false
	}

	if move.GetPromote() != NoPiece {
		if (is_white && re != 7) || (!is_white && re != 0) {
			return false
		}
	}

	return true
}

func IsRookMovePossible(move Move, board *Board, is_white bool) bool {
	start := move.GetStart()
	end := move.GetEnd()
	rs := start.GetRow()
	cs := start.GetCol()
	re := end.GetRow()
	ce := end.GetCol()
	if start == end {
		return false
	}

	if isTakenByFriend(board, end, is_white) {
		return false
	}

	switch {
	case cs == ce:
		var dir int8 = 1
		if rs > re {
			dir = -1
		}
		for ri := rs + dir; ri != re; ri += dir {
			if isTaken(board, MakePos(ri, cs)) {
				return false
			}
		}
	case rs == re:
		var dir int8 = 1
		if cs > ce {
			dir = -1
		}
		for ci := cs + dir; ci != ce; ci += dir {
			if isTaken(board, MakePos(rs, ci)) {
				return false
			}
		}
	default:
		return false
	}

	return true
}

func IsKnightMovePossible(move Move, board *Board, is_white bool) bool {
	start := move.GetStart()
	end := move.GetEnd()
	rs := start.GetRow()
	cs := start.GetCol()
	re := end.GetRow()
	ce := end.GetCol()
	if isTakenByFriend(board, end, is_white) {
		return false
	}
	rd := Abs(rs - re)
	cd := Abs(cs - ce)
	if !((cd == 2 && rd == 1) || (cd == 1 && rd == 2)) {
		return false
	}
	return true
}

func IsBishopMovePossible(move Move, board *Board, is_white bool) bool {
	start := move.GetStart()
	end := move.GetEnd()
	rs := start.GetRow()
	cs := start.GetCol()
	re := end.GetRow()
	ce := end.GetCol()
	if start == end {
		return false
	}
	dist := Abs(rs - re)
	if dist != Abs(cs-ce) {
		return false
	}
	var dir_r int8 = 1
	if rs > re {
		dir_r = -1
	}
	var dir_c int8 = 1
	if cs > ce {
		dir_c = -1
	}
	ri := rs
	ci := cs
	for i := int8(0); i < dist-1; i++ {
		ri += dir_r
		ci += dir_c
		if isTaken(board, MakePos(ri, ci)) {
			return false
		}
	}

	return true
}

func IsQueenMovePossible(move Move, board *Board, is_white bool) bool {
	return IsRookMovePossible(move, board, is_white) || IsBishopMovePossible(move, board, is_white)
}

func IsKingMovePossible(move Move, board *Board, bs BoardState, is_white bool) bool {
	start := move.GetStart()
	end := move.GetEnd()
	rs := start.GetRow()
	cs := start.GetCol()
	re := end.GetRow()
	ce := end.GetCol()
	if start == end {
		return false
	}
	if isTakenByFriend(board, end, is_white) {
		return false
	}
	if Abs(rs-re) < 2 && Abs(cs-ce) < 2 {
		return true
	}
	// Check castles
	var r int8
	if is_white {
		r = 0
	} else {
		r = 7
	}

	if start != MakePos(r, 4) {
		return false
	}

	switch {
	case end == MakePos(r, 6):
		if (is_white && !bs.Get_K()) || (!is_white && !bs.Get_k()) {
			return false
		}
		{
			rook_pice := board.GetPiece(MakePos(r, 7))
			if is_white && rook_pice != W_Rook || !is_white && rook_pice != B_Rook {
				return false
			}
		}
		for ci := int8(4); ci < 7; ci++ {
			if IsPosUnderAttack(MakePos(r, ci), board, is_white) {
				return false
			}
		}
	case end == MakePos(r, 2):
		if (is_white && !bs.Get_Q()) || (!is_white && !bs.Get_q()) {
			return false
		}
		{
			rook_pice := board.GetPiece(MakePos(r, 0))
			if is_white && rook_pice != W_Rook || !is_white && rook_pice != B_Rook {
				return false
			}
		}
		for ci := int8(2); ci <= 4; ci++ {
			if IsPosUnderAttack(MakePos(r, ci), board, is_white) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func IsMovePossible(move Move, board *Board, bs BoardState) bool {
	piece := board.GetPiece(move.GetStart())
	is_white := bs.Get_Turn()

	if is_white != piece.IsWhite() {
		return false
	}

	var is_possible bool
	switch piece.GetType() {
	case Pawn:
		is_possible = IsPawnMovePossible(move, board, is_white, bs.Get_IsEnPos(), bs.Get_EnPos())
	case Knight:
		is_possible = IsKnightMovePossible(move, board, is_white)
	case Bishop:
		is_possible = IsBishopMovePossible(move, board, is_white)
	case Rook:
		is_possible = IsRookMovePossible(move, board, is_white)
	case Queen:
		is_possible = IsQueenMovePossible(move, board, is_white)
	case King:
		is_possible = IsKingMovePossible(move, board, bs, is_white)
	default:
		is_possible = false
	}

	if !is_possible {
		return false
	}

	var king Piece
	if is_white {
		king = W_King
	} else {
		king = B_King
	}

	for i, piece := range board {
		if piece != king {
			continue
		}
		if IsPosUnderAttack(Position(i), board, is_white) {
			is_possible = false
		}
		break
	}

	return false
}
