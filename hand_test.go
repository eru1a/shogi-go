package shogi

import (
	"reflect"
	"testing"
)

func TestHand(t *testing.T) {
	h := NewHand()

	for i := 1; i <= 20; i++ {
		h.Add(FU, Black)
		n, err := h.Get(FU, Black)
		if err != nil {
			t.Errorf("h.Get(FU, Black) error %v", err)
		}
		if i != n {
			t.Errorf("h.Get(FU, Black) want %v, got %v", i, n)
		}
	}
	for i := 19; i >= 0; i-- {
		h.Remove(FU, Black)
		n, err := h.Get(FU, Black)
		if err != nil {
			t.Errorf("h.Get(FU, Black) error %v", err)
		}
		if i != n {
			t.Errorf("h.Get(FU, Black) want %v, got %v", i, n)
		}
	}

	h.Add(GI, White)
	h.Add(GI, White)
	if err := h.Remove(GI, White); err != nil {
		t.Errorf("h.Remove(GI, White) error %v", err)
	}
	if err := h.Remove(GI, White); err != nil {
		t.Errorf("h.Remove(GI, White) error %v", err)
	}
	if err := h.Remove(GI, White); err == nil {
		t.Errorf("h.Remove(GI, White) want error")
	}

	ngTests := []struct {
		pt PieceType
		c  Color
	}{
		{NO_PIECE_TYPE, NO_COLOR},
		{NO_PIECE_TYPE, Black},
		{NO_PIECE_TYPE, White},
		{FU, NO_COLOR},
		{OU, Black},
		{TO, White},
		{RY, Black},
	}

	for _, test := range ngTests {
		if err := h.Add(test.pt, test.c); err == nil {
			t.Errorf("h.Add(%v, %v) want error", test.pt, test.c)
		}
		if err := h.Remove(test.pt, test.c); err == nil {
			t.Errorf("h.Remove(%v, %v) want error", test.pt, test.c)
		}
		if n, err := h.Get(test.pt, test.c); err == nil {
			t.Errorf("h.Got(%v, %v) want error, got %v", test.pt, test.c, n)
		}
	}
}

func handFromArray(b, w [7]int) *Hand {
	h := NewHand()
	for pt := FU; pt <= HI; pt++ {
		h.Black[pt] += b[pt-1]
		h.White[pt] += w[pt-1]
	}
	return h
}

func TestHandSFEN(t *testing.T) {
	okTests := []struct {
		sfen string
		hand *Hand
	}{
		{"-", NewHand()},
		{"P", handFromArray([7]int{1, 0, 0, 0, 0, 0, 0}, [7]int{0, 0, 0, 0, 0, 0, 0})},
		{"2P3p", handFromArray([7]int{2, 0, 0, 0, 0, 0, 0}, [7]int{3, 0, 0, 0, 0, 0, 0})},
		{"S2Pb3p", handFromArray([7]int{2, 0, 0, 1, 0, 0, 0}, [7]int{3, 0, 0, 0, 0, 1, 0})},
		{"4S4N4L4P2r2b4g20p", handFromArray([7]int{4, 4, 4, 4, 0, 0, 0}, [7]int{20, 0, 0, 0, 4, 2, 2})},
	}

	for _, test := range okTests {
		// sfen -> hand
		hand, err := NewHandFromSFEN(test.sfen)
		if err != nil {
			t.Errorf("NewHandFromSFEN(%s): got error %v", test.sfen, err)
		}
		if !reflect.DeepEqual(hand, test.hand) {
			t.Errorf("NewHandFromSFEN(%s): want %v, got %v", test.sfen, test.hand, hand)
		}

		// hand -> sfen
		sfen := test.hand.SFEN()
		if sfen != test.sfen {
			t.Errorf("%v.SFEN(): want %v, got %v", test.hand, test.sfen, sfen)
		}
	}

	ngTests := []struct {
		sfen string
	}{
		{""},
		{"A"},
		{"1"},
		// {"K"},
		{"+P"},
		{"P20"},
		{"P p"},
		// 重複もエラーにすべき
		// {"PP"},
	}

	for _, test := range ngTests {
		hand, err := NewHandFromSFEN(test.sfen)
		if err == nil {
			t.Errorf("NewHandFromSFEN(%s): want error, got %v", test.sfen, hand)
		}
	}
}
