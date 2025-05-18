package main

import (
	"testing"
)

func assert_equal[T comparable](a, b T, t *testing.T) {
	if a != b {
		t.Error("values are not equal")
	}
}

func assert_er(er error, t *testing.T) {
	if er != nil {
		t.Error(er.Error())
	}
}

func TestInitialBoardState(t *testing.T) {
	bs := MakeInitialBoardState()
	if !bs.Get_Turn() {
		t.Error("wrong turn")
	}

	if !bs.Get_K() {
		t.Error("wrong castle")
	}
	if !bs.Get_Q() {
		t.Error("wrong castle")
	}
	if !bs.Get_k() {
		t.Error("wrong castle")
	}
	if !bs.Get_q() {
		t.Error("wrong castle")
	}

	if bs.Get_IsEnPos() {
		t.Error("wrong en passant")
	}

	if bs.Get_HMoves() != 1 {
		t.Error("wrong half move count")
	}
}

func TestBoardState(t *testing.T) {
	for _, is_w := range [2]bool{true, false} {
		for _, K := range [2]bool{true, false} {
			for _, Q := range [2]bool{true, false} {
				for _, k := range [2]bool{true, false} {
					for _, q := range [2]bool{true, false} {
						for _, is_en := range [2]bool{true, false} {
							for _, en_pos := range [3]Position{MakePos(1, 2), MakePos(7, 0), MakePos(7, 7)} {
								for _, hmove := range [3]uint8{1, 13, 50} {
									var bs BoardState
									bs = bs.Set_Turn(is_w)
									bs = bs.Set_K(K)
									bs = bs.Set_Q(Q)
									bs = bs.Set_k(k)
									bs = bs.Set_q(q)
									bs = bs.Set_IsEnPos(is_en)
									if is_en {
										bs = bs.Set_EnPos(en_pos)
									}
									bs = bs.Set_HMoves(hmove)

									assert_equal(bs.Get_Turn(), is_w, t)
									assert_equal(bs.Get_K(), K, t)
									assert_equal(bs.Get_Q(), Q, t)
									assert_equal(bs.Get_k(), k, t)
									assert_equal(bs.Get_q(), q, t)
									assert_equal(bs.Get_IsEnPos(), is_en, t)
									if is_en {
										assert_equal(bs.Get_EnPos(), en_pos, t)
									}
									assert_equal(bs.Get_HMoves(), hmove, t)

									fen := bs.FEN()
									bs_copy, er := MakeBoardStateFromFEN(fen)
									assert_er(er, t)
									assert_equal(bs, bs_copy, t)
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestBoardStateFEN(t *testing.T) {
	bs := MakeInitialBoardState()
	var fen string
	fen = bs.FEN()
	assert_equal(fen, "w KQkq - 1", t)
	bs_copy, er := MakeBoardStateFromFEN(fen)
	assert_er(er, t)
	assert_equal(bs, bs_copy, t)

	bs = bs.Set_Q(false)
	bs = bs.Set_IsEnPos(true)
	bs = bs.Set_EnPos(MakePos(1, 2))
	fen = bs.FEN()
	assert_equal(fen, "w Kkq c2 1", t)
	bs_copy, er = MakeBoardStateFromFEN(fen)
	assert_er(er, t)
	assert_equal(bs, bs_copy, t)

	bs = bs.Set_Turn(false)
	bs = bs.Set_HMoves(3)
	bs = bs.Set_K(false)
	bs = bs.Set_k(false)
	bs = bs.Set_q(false)
	fen = bs.FEN()
	assert_equal(fen, "b - c2 3", t)
	bs_copy, er = MakeBoardStateFromFEN(fen)
	assert_er(er, t)
	assert_equal(bs, bs_copy, t)
}
