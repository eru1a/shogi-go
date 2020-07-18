package shogi

import "testing"

func TestMoveUSI(t *testing.T) {
	okTests := []struct {
		usi  string
		move Move
	}{
		{"7g7f", NewNormalMove(Square{2, 6}, Square{2, 5}, false)},
		{"8h2b+", NewNormalMove(Square{1, 7}, Square{7, 1}, true)},
		{"G*5h", NewDropMove(KI, Square{4, 7})},
	}

	for _, test := range okTests {
		// usi -> move
		move, err := NewMoveFromUSI(test.usi)
		if err != nil {
			t.Errorf("NewMoveFromUSI(%s): want %v, got error", test.usi, test.move)
		}
		if move != test.move {
			t.Errorf("NewMoveFromUSI(%s): want %v, got %v", test.usi, test.move, move)
		}

		// move -> usi
		usi := test.move.USI()
		if usi != test.usi {
			t.Errorf("%v.USI(): want %v, got %v", test.move, test.usi, usi)
		}
	}

	ngTests := []struct {
		usi string
	}{
		{""},
		{"7g"},
		{"P"},
		{"7g0f"},
		{"7g7f*"},
		{"10g7f*"},
		{"G+5h"},
		{"G*5j"},
		{"A*5h"},
	}

	for _, test := range ngTests {
		move, err := NewMoveFromUSI(test.usi)
		if err == nil {
			t.Errorf("NewMoveFromUSI(%s): want error, but %v", test.usi, move)
		}
	}
}
