package shogi

import "testing"

func TestNewSquare(t *testing.T) {
	okTests := []struct {
		file   int
		rank   int
		expect Square
	}{
		{0, 0, Square{0, 0}},
		{4, 6, Square{4, 6}},
		{8, 8, Square{8, 8}},
	}

	for _, test := range okTests {
		square, err := NewSquare(test.file, test.rank)
		if err != nil {
			t.Errorf("NewSquare(%d,%d): error %v",
				test.file, test.rank, err)
		}
		if square != test.expect {
			t.Errorf("NewSquare(%d,%d): want %v, got %v",
				test.file, test.rank, test.expect, square)
		}
	}

	ngTests := []struct {
		file int
		rank int
	}{
		{-1, -1},
		{0, 9},
		{9, 5},
		{10, 10},
	}

	for _, test := range ngTests {
		square, err := NewSquare(test.file, test.rank)
		if err == nil {
			t.Errorf("NewSquare(%d,%d): want error, got %v",
				test.file, test.rank, square)
		}
	}
}

func TestSquareAdd(t *testing.T) {
	okTests := []struct {
		square Square
		file   int
		rank   int
		expect Square
	}{
		{Square{0, 0}, 0, 0, Square{0, 0}},
		{Square{0, 0}, 1, 1, Square{1, 1}},
		{Square{5, 5}, 3, -4, Square{8, 1}},
		{Square{8, 8}, -2, -3, Square{6, 5}},
	}

	for _, test := range okTests {
		square, err := test.square.Add(test.file, test.rank)
		if err != nil {
			t.Errorf("%v.Add(%d,%d): error %v", test.square, test.file, test.rank, err)
		}
		if square != test.expect {
			t.Errorf("%v.Add(%d,%d): want %v, got %v",
				test.square, test.file, test.rank, test.expect, square)
		}
	}

	ngTests := []struct {
		square Square
		file   int
		rank   int
	}{
		{Square{0, 0}, 1, -1},
		{Square{5, 5}, 5, 5},
		{Square{8, 8}, -2, 3},
	}

	for _, test := range ngTests {
		square, err := test.square.Add(test.file, test.rank)
		if err == nil {
			t.Errorf("%v.Add(%d,%d): want error, got %v",
				test.square, test.file, test.rank, square)
		}
	}
}

func TestSquareUSI(t *testing.T) {
	tests := []struct {
		usi    string
		square Square
	}{
		{"1a", Square{8, 0}},
		{"9a", Square{0, 0}},
		{"7g", Square{2, 6}},
		{"9i", Square{0, 8}},
	}

	for _, test := range tests {
		square, err := NewSquareFromUSI(test.usi)
		if err != nil {
			t.Errorf("NewSquareFromUSI(%s): got error %v", test.usi, err)
		}
		if square != test.square {
			t.Errorf("NewSquareFromUSI(%s): want %v, got %v", test.usi, test.square, square)
		}

		usi := test.square.USI()
		if usi != test.usi {
			t.Errorf("%v.USI(): want %v, got %v", test.square, test.usi, usi)
		}
	}
}

func TestSquareKIF(t *testing.T) {
	tests := []struct {
		kif    string
		square Square
	}{
		{"７六", Square{2, 5}},
		{"５九", Square{4, 8}},
		{"１一", Square{8, 0}},
	}

	for _, test := range tests {
		kif := test.square.KIF()
		if kif != test.kif {
			t.Errorf("%v.KIF(): want %v, got %v", test.square, test.kif, kif)
		}
	}
}
