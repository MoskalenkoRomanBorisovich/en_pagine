package main

/*
functions NextXXXMove generate moves by id in order
first call should always be with id = 0, and consecutive calls should use returned id
if function returns false, then it is finished, and currently returned move is invalid

*/

type MoveStencil [2]int8

var SentryStencil = MoveStencil{0, 0} // padding place holder

func applyStencil(pos Position, st MoveStencil) (r int8, c int8) {
	return pos.GetRow() + st[0], pos.GetCol() + st[1]
}

// Rook

var RookStencils = [...]MoveStencil{
	{0, 1},
	{0, 2},
	{0, 3},
	{0, 4},
	{0, 5},
	{0, 6},
	{0, 7},
	SentryStencil,

	{1, 0},
	{2, 0},
	{3, 0},
	{4, 0},
	{5, 0},
	{6, 0},
	{7, 0},
	SentryStencil,

	{0, -1},
	{0, -2},
	{0, -3},
	{0, -4},
	{0, -5},
	{0, -6},
	{0, -7},
	SentryStencil,

	{-1, 0},
	{-2, 0},
	{-3, 0},
	{-4, 0},
	{-5, 0},
	{-6, 0},
	{-7, 0},
	SentryStencil,
}

var BishopStencils = [...]MoveStencil{
	{1, 1},
	{2, 2},
	{3, 3},
	{4, 4},
	{5, 5},
	{6, 6},
	{7, 7},
	SentryStencil,

	{-1, 1},
	{-2, 2},
	{-3, 3},
	{-4, 4},
	{-5, 5},
	{-6, 6},
	{-7, 7},
	SentryStencil,

	{1, -1},
	{2, -2},
	{3, -3},
	{4, -4},
	{5, -5},
	{6, -6},
	{7, -7},
	SentryStencil,

	{-1, -1},
	{-2, -2},
	{-3, -3},
	{-4, -4},
	{-5, -5},
	{-6, -6},
	{-7, -7},
	SentryStencil,
}

var QueenStencils = [...]MoveStencil{
	// straights
	{0, 1},
	{0, 2},
	{0, 3},
	{0, 4},
	{0, 5},
	{0, 6},
	{0, 7},
	SentryStencil,

	{1, 0},
	{2, 0},
	{3, 0},
	{4, 0},
	{5, 0},
	{6, 0},
	{7, 0},
	SentryStencil,

	{0, -1},
	{0, -2},
	{0, -3},
	{0, -4},
	{0, -5},
	{0, -6},
	{0, -7},
	SentryStencil,

	{-1, 0},
	{-2, 0},
	{-3, 0},
	{-4, 0},
	{-5, 0},
	{-6, 0},
	{-7, 0},
	SentryStencil,

	// diagonals
	{1, 1},
	{2, 2},
	{3, 3},
	{4, 4},
	{5, 5},
	{6, 6},
	{7, 7},
	SentryStencil,

	{-1, 1},
	{-2, 2},
	{-3, 3},
	{-4, 4},
	{-5, 5},
	{-6, 6},
	{-7, 7},
	SentryStencil,

	{1, -1},
	{2, -2},
	{3, -3},
	{4, -4},
	{5, -5},
	{6, -6},
	{7, -7},
	SentryStencil,

	{-1, -1},
	{-2, -2},
	{-3, -3},
	{-4, -4},
	{-5, -5},
	{-6, -6},
	{-7, -7},
	SentryStencil,
}

func NextSlidingPieceMove(pos Position, is_white bool, id uint, board *Board, stencil_arr []MoveStencil) (move Move, next_id uint, is_finished bool) {
	const stencil_mask = 0b111
	const step = stencil_mask + 1
	var res Move
	for id < uint(len(stencil_arr)) {
		if (id & stencil_mask) == stencil_mask {
			id++
			continue
		}
		stencil := stencil_arr[id]
		re, ce := applyStencil(pos, stencil)
		if !CheckBoardPos(re, ce) {
			id = (id &^ stencil_mask) + step
			continue
		}
		end_pos := MakePos(re, ce)
		end_piece := board.GetPiece(end_pos)
		if end_piece.GetType() != NoPiece { // line is blocked here
			id = (id &^ stencil_mask) + step // go to next direction
			if end_piece.IsWhite() == is_white {
				continue
			}
		} else {
			id++
		}
		res.SetStart(pos)
		res.SetEnd(end_pos)
		return res, id, false
	}
	return res, id, true
}

func NextRookMove(pos Position, is_white bool, id uint, board *Board) (move Move, next_id uint, is_finished bool) {
	return NextSlidingPieceMove(pos, is_white, id, board, RookStencils[:])
}

// Bishop

func NextBishopMove(pos Position, is_white bool, id uint, board *Board) (move Move, next_id uint, is_finished bool) {
	return NextSlidingPieceMove(pos, is_white, id, board, BishopStencils[:])
}

// Queen

func NextQueenMove(pos Position, is_white bool, id uint, board *Board) (move Move, next_id uint, is_finished bool) {
	return NextSlidingPieceMove(pos, is_white, id, board, QueenStencils[:])
}

// Knight

var KnightStencils = [...]MoveStencil{
	{1, 2},
	{2, 1},
	{2, -1},
	{1, -2},
	{-1, -2},
	{-2, -1},
	{-2, 1},
	{-1, 2},
}

func NextKnightMove(pos Position, is_white bool, id uint, board *Board) (move Move, next_id uint, is_finished bool) {
	var res Move
	for id < uint(len(KnightStencils)) {
		id++
		stencil := KnightStencils[id]
		re, ce := applyStencil(pos, stencil)
		if !CheckBoardPos(re, ce) {
			continue
		}
		end_pos := MakePos(re, ce)
		if isTakenByFriend(board, end_pos, is_white) {
			continue
		}
		res.SetStart(pos)
		res.SetEnd(end_pos)
		return res, id, false
	}

	return res, id, true
}

var KingStencils = [...]MoveStencil{
	// Moves
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
	{-1, -1},
	{0, -1},
	{1, -1},
}

func NextKingMove(pos Position, is_white bool, id uint, board *Board) (move Move, next_id uint, is_finished bool) {
	var res Move
	for id < uint(len(KingStencils)) {
		id++
		stencil := KingStencils[id]
		re, ce := applyStencil(pos, stencil)
		if !CheckBoardPos(re, ce) {
			continue
		}
		end_pos := MakePos(re, ce)
		if isTakenByFriend(board, end_pos, is_white) {
			continue
		}
		res.SetStart(pos)
		res.SetEnd(end_pos)
		return res, id, false
	}

	return res, id, true
}

func NextPawnMove(pos Position, is_white bool, id uint, board *Board) (move Move, next_id uint, is_finished bool) {
	var res Move
	res.SetStart(pos)
	var direction int8
	if is_white {
		direction = 1
	} else {
		direction = -1
	}
	for id < 4 {
		id_cur := id
		id++

		ce := pos.GetCol()
		re := pos.GetRow()

		switch id_cur {
		case 0: // move by 1
			re += direction
		case 1: // move by 2
			re += direction * 2
		case 2: // capture left
			ce--
			re += direction
		case 3: // capture right
			ce++
			re += direction
		}

		if !CheckBoardPos(re, ce) {
			continue
		}

		end_pos := MakePos(re, ce)
		switch id_cur {
		case 0: // move by 1
			fallthrough
		case 1: // move by 2
			if board.GetPiece(end_pos) != NoPiece {
				continue
			}
		case 2: // capture left
			fallthrough
		case 3: // capture right
			if p := board.GetPiece(end_pos); p == NoPiece || p.IsWhite() == is_white {
				continue
			}
		}

		res.SetStart(pos)
		res.SetEnd(end_pos)
		return res, id, false
	}
	return res, id, true
}

func NextEnPassantMove(enPos Position, id uint, board *Board, is_white bool) (move Move, next_id uint, is_finished bool) {
	var res Move

	rs_ := enPos.GetRow()
	if is_white {
		rs_--
	} else {
		rs_++
	}
	cs_ := enPos.GetCol()

	{ // check if captured piece is there
		if !CheckBoardPos(rs_, cs_) {
			id = 2
			return res, id, true
		}
		capture_pose := MakePos(rs_, cs_)
		capture_piece := board.GetPiece(capture_pose)
		if (is_white && capture_piece != B_Pawn) ||
			(!is_white && capture_piece != W_Pawn) {
			id = 2
			return res, id, true
		}
	}

	for id < 2 {
		rs := rs_
		cs := cs_
		switch id {
		case 0:
			cs++
		case 1:
			cs--
		}
		id++
		if !CheckBoardPos(rs, cs) {
			continue
		}

		start_pos := MakePos(rs, cs)
		start_piece := board.GetPiece(start_pos)
		if (is_white && start_piece != W_Pawn) ||
			(!is_white && start_piece != B_Pawn) {
			continue
		}

		res.SetStart(start_pos)
		res.SetEnd(enPos)
		return res, id, false
	}
	return res, id, true
}

func IsPositionUnderAttack(pos Position, board *Board, is_white bool) bool {
	{ // bishop, rook, queen
		var id uint = 0
		var finished = false
		var move Move
		for {
			move, id, finished = NextQueenMove(pos, !is_white, id, board)
			if finished {
				break
			}
			end_pos := move.GetEnd()
			ep := board.GetPiece(end_pos)
			if (is_white && (ep == B_Bishop || ep == B_Rook || ep == B_Queen)) ||
				(!is_white && (ep == W_Bishop || ep == W_Rook || ep == W_Queen)) {
				return true
			}
		}
	}
	{ // knight
		var id uint = 0
		var finished = false
		var move Move
		for {
			move, id, finished = NextKnightMove(pos, !is_white, id, board)
			if finished {
				break
			}
			end_pos := move.GetEnd()
			ep := board.GetPiece(end_pos)
			if (is_white && ep == B_Knight) ||
				(!is_white && ep == W_Knight) {
				return true
			}
		}
	}
	{ // pawn
		var id uint = 2 // check only captures
		var finished = false
		var move Move
		for {
			move, id, finished = NextPawnMove(pos, is_white, id, board)
			if finished {
				break
			}
			end_pos := move.GetEnd()
			ep := board.GetPiece(end_pos)
			if (is_white && ep == B_Pawn) ||
				(!is_white && ep == W_Pawn) {
				return true
			}
		}
	}
	{ // king
		var id uint = 0
		var finished = false
		var move Move
		for {
			move, id, finished = NextKingMove(pos, !is_white, id, board)
			if finished {
				break
			}
			end_pos := move.GetEnd()
			ep := board.GetPiece(end_pos)
			if (is_white && ep == B_King) ||
				(!is_white && ep == W_King) {
				return true
			}
		}
	}
	return false
}

type CastleType uint8

const (
	CastleType_King  CastleType = 0
	CastleType_Queen CastleType = 1
)

/*
castle types : 0 = king, 1 = queen
*/

func NextCastleMove(id uint, board *Board, bs BoardState) (castle_type CastleType, next_id uint, is_finished bool) {
	is_white := bs.Get_Turn()
	var r int8
	if is_white {
		r = 0
	} else {
		r = 7
	}
	for id < 2 {
		id_cur := id
		id++
		switch id_cur {
		case 0: // King side
			if is_white && !bs.Get_K() || !is_white && !bs.Get_k() {
				continue
			}
			p_f := MakePos(r, 5)
			p_g := MakePos(r, 6)
			if board.GetPiece(p_f).GetType() != NoPiece ||
				board.GetPiece(p_g).GetType() != NoPiece {
				continue
			}
			if IsPosUnderAttack(MakePos(r, 4), board, is_white) { // is king in check
				continue
			}
			if IsPositionUnderAttack(p_f, board, is_white) || // is passing squares in check
				IsPositionUnderAttack(p_g, board, is_white) {
				continue
			}
			return 0, id, false
		case 1:
			if is_white && !bs.Get_Q() || !is_white && !bs.Get_q() {
				continue
			}
			p_d := MakePos(r, 3)
			p_c := MakePos(r, 2)
			if board.GetPiece(p_c).GetType() != NoPiece ||
				board.GetPiece(p_d).GetType() != NoPiece {
				continue
			}
			if IsPosUnderAttack(MakePos(r, 4), board, is_white) { // is king in check
				continue
			}
			if IsPositionUnderAttack(p_c, board, is_white) || // is passing squares in check
				IsPositionUnderAttack(p_d, board, is_white) {
				continue
			}
			return 1, id, true
		}
	}
	return 0, id, true
}

func NextPromotionMove(id uint, is_white bool) (Piece, uint, bool) {
	var piece Piece
	switch id {
	case 0:
		piece = Queen
	case 1:
		piece = Knight
	case 2:
		piece = Bishop
	case 3:
		piece = Rook
	default:
		return NoPiece, id, true
	}
	id++
	if is_white {
		piece = piece | White
	}
	return piece, id, false
}
