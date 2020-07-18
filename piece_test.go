package shogi

import "testing"

func TestPieceTypeUSI(t *testing.T) {
	okTests := []struct {
		usi string
		pt  PieceType
	}{
		{"P", FU},
		{"L", KY},
		{"N", KE},
		{"S", GI},
		{"G", KI},
		{"B", KA},
		{"R", HI},
		{"K", OU},
		{"+P", TO},
		{"+L", NY},
		{"+N", NK},
		{"+S", NG},
		{"+B", UM},
		{"+R", RY},
	}

	for _, test := range okTests {
		// pt -> usi
		pt, err := NewPieceTypeFromUSI(test.usi)
		if err != nil {
			t.Errorf("NewPieceTypeFromUSI(%v): got error %v", test.usi, err)
		}
		if pt != test.pt {
			t.Errorf("NewPieceTypeFromUSI(%v): want %v, got %v",
				test.usi, test.pt, pt)
		}

		usi := test.pt.USI()
		if usi != test.usi {
			t.Errorf("%v.USI(): want %v, got %v",
				test.pt, test.usi, usi)
		}
	}

	ngTests := []struct {
		input string
	}{
		{""},
		{"p"},
		{"r+"},
		{"-r"},
		{"+K"},
		{"A"},
		{"PP"},
		{"P "},
		{" P"},
		{"+ P"},
		{"P +"},
		{"B+"},
	}

	for _, test := range ngTests {
		pt, err := NewPieceTypeFromUSI(test.input)
		if err == nil {
			t.Errorf("NewPieceTypeFromUSI(%v): want error, got %v", test.input, pt)
		}
	}
}

func TestPieceUSI(t *testing.T) {
	okTests := []struct {
		usi   string
		piece Piece
	}{
		{"P", BFU},
		{"L", BKY},
		{"N", BKE},
		{"S", BGI},
		{"G", BKI},
		{"B", BKA},
		{"R", BHI},
		{"K", BOU},
		{"+P", BTO},
		{"+L", BNY},
		{"+N", BNK},
		{"+S", BNG},
		{"+B", BUM},
		{"+R", BRY},
		{"p", WFU},
		{"l", WKY},
		{"n", WKE},
		{"s", WGI},
		{"g", WKI},
		{"b", WKA},
		{"r", WHI},
		{"k", WOU},
		{"+p", WTO},
		{"+l", WNY},
		{"+n", WNK},
		{"+s", WNG},
		{"+b", WUM},
		{"+r", WRY},
	}

	for _, test := range okTests {
		// usi -> piece
		piece, err := NewPieceFromUSI(test.usi)
		if err != nil {
			t.Errorf("NewPieceFromUSI(%v): got error %v", test.usi, err)
		}
		if piece != test.piece {
			t.Errorf("NewPieceFromUSI(%v): want %v, got %v",
				test.usi, test.piece, piece)
		}

		// piece -> usi
		usi := test.piece.USI()
		if usi != test.usi {
			t.Errorf("%v.USI(): want %v, got %v",
				test.piece, test.usi, usi)
		}
	}

	ngTests := []struct {
		input string
	}{
		{""},
		{"-P"},
		{"P+"},
		{"+K"},
		{"A"},
		{"PP"},
		{"P "},
		{" P"},
		{"+ P"},
		{"B+"},
	}

	for _, test := range ngTests {
		pt, err := NewPieceFromUSI(test.input)
		if err == nil {
			t.Errorf("NewPieceTypeFromUSI(%v): want error, got %v", test.input, pt)
		}
	}
}

func TestColor(t *testing.T) {
	tests := []struct {
		piece  Piece
		expect Color
	}{
		{BFU, Black},
		{BGI, Black},
		{BOU, Black},
		{BUM, Black},
		{WFU, White},
		{WGI, White},
		{WOU, White},
		{WUM, White},
		{NO_PIECE, NO_COLOR},
	}

	for _, test := range tests {
		c := test.piece.Color()
		if c != test.expect {
			t.Errorf("%v.Color() wand %v, got %v", test.piece, test.expect, c)
		}
	}
}

func TestPieceTypePromote(t *testing.T) {
	tests := []struct {
		input  PieceType
		expect PieceType
	}{
		{NO_PIECE_TYPE, NO_PIECE_TYPE},
		{FU, TO},
		{KY, NY},
		{KE, NK},
		{GI, NG},
		{KI, KI},
		{KA, UM},
		{HI, RY},
		{OU, OU},
		{TO, TO},
		{NY, NY},
		{NK, NK},
		{NG, NG},
		{UM, UM},
		{RY, RY},
	}

	for _, test := range tests {
		pt := test.input.Promote()
		if pt != test.expect {
			t.Errorf("%v.Promote(): want %v, got %v", test.input, test.expect, pt)
		}
	}
}

func TestPieceTypeDemote(t *testing.T) {
	tests := []struct {
		input  PieceType
		expect PieceType
	}{
		{NO_PIECE_TYPE, NO_PIECE_TYPE},
		{FU, FU},
		{KY, KY},
		{KE, KE},
		{GI, GI},
		{KI, KI},
		{KA, KA},
		{HI, HI},
		{OU, OU},
		{TO, FU},
		{NY, KY},
		{NK, KE},
		{NG, GI},
		{UM, KA},
		{RY, HI},
	}

	for _, test := range tests {
		pt := test.input.Demote()
		if pt != test.expect {
			t.Errorf("%v.Demote(): want %v, got %v", test.input, test.expect, pt)
		}
	}
}

func TestPiecePromote(t *testing.T) {
	tests := []struct {
		input  Piece
		expect Piece
	}{
		{NO_PIECE, NO_PIECE},
		{BFU, BTO},
		{BKY, BNY},
		{BKE, BNK},
		{BGI, BNG},
		{BKI, BKI},
		{BKA, BUM},
		{BHI, BRY},
		{BOU, BOU},
		{BTO, BTO},
		{BNY, BNY},
		{BNK, BNK},
		{BNG, BNG},
		{BUM, BUM},
		{BRY, BRY},
		{WFU, WTO},
		{WKY, WNY},
		{WKE, WNK},
		{WGI, WNG},
		{WKI, WKI},
		{WKA, WUM},
		{WHI, WRY},
		{WOU, WOU},
		{WTO, WTO},
		{WNY, WNY},
		{WNK, WNK},
		{WNG, WNG},
		{WUM, WUM},
		{WRY, WRY},
	}

	for _, test := range tests {
		pt := test.input.Promote()
		if pt != test.expect {
			t.Errorf("%v.Promote(): want %v, got %v", test.input, test.expect, pt)
		}
	}
}

func TestPieceDemote(t *testing.T) {
	tests := []struct {
		input  Piece
		expect Piece
	}{
		{NO_PIECE, NO_PIECE},
		{BFU, BFU},
		{BKY, BKY},
		{BKE, BKE},
		{BGI, BGI},
		{BKI, BKI},
		{BKA, BKA},
		{BHI, BHI},
		{BOU, BOU},
		{BTO, BFU},
		{BNY, BKY},
		{BNK, BKE},
		{BNG, BGI},
		{BUM, BKA},
		{BRY, BHI},
		{WFU, WFU},
		{WKY, WKY},
		{WKE, WKE},
		{WGI, WGI},
		{WKI, WKI},
		{WKA, WKA},
		{WHI, WHI},
		{WOU, WOU},
		{WTO, WFU},
		{WNY, WKY},
		{WNK, WKE},
		{WNG, WGI},
		{WUM, WKA},
		{WRY, WHI},
	}

	for _, test := range tests {
		pt := test.input.Demote()
		if pt != test.expect {
			t.Errorf("%v.Demote(): want %v, got %v", test.input, test.expect, pt)
		}
	}
}

func TestCanPromote(t *testing.T) {
	tests := []struct {
		from, to Square
		p        Piece
		expect   bool
	}{
		{
			from:   Square{1, 7},
			to:     Square{7, 1},
			p:      NO_PIECE,
			expect: false,
		},
		{
			from:   Square{1, 7},
			to:     Square{7, 1},
			p:      BKA,
			expect: true,
		},
		{
			from:   Square{4, 2},
			to:     Square{5, 3},
			p:      BGI,
			expect: true,
		},
		{
			from:   Square{4, 2},
			to:     Square{5, 2},
			p:      BKI,
			expect: false,
		},
		{
			from:   Square{4, 2},
			to:     Square{3, 0},
			p:      BKE,
			expect: true,
		},
		{
			from:   Square{1, 5},
			to:     Square{1, 6},
			p:      WFU,
			expect: true,
		},
		{
			from:   Square{4, 2},
			to:     Square{5, 3},
			p:      WGI,
			expect: false,
		},
	}

	for _, test := range tests {
		got := CanPromote(test.p, test.from, test.to)
		if got != test.expect {
			t.Errorf("CanPromote(%v, %v, %v): want %v, got %v",
				test.from, test.to, test.p, test.expect, got)
		}
	}
}

func TestNeedForcePromotion(t *testing.T) {
	tests := []struct {
		p      Piece
		rank   int
		expect bool
	}{
		{
			p:      BFU,
			rank:   0,
			expect: true,
		},
		{
			p:      BFU,
			rank:   1,
			expect: false,
		},
		{
			p:      BKY,
			rank:   0,
			expect: true,
		},
		{
			p:      BKY,
			rank:   1,
			expect: false,
		},
		{
			p:      BKE,
			rank:   0,
			expect: true,
		},
		{
			p:      BKE,
			rank:   1,
			expect: true,
		},
		{
			p:      BKE,
			rank:   2,
			expect: false,
		},
		{
			p:      WFU,
			rank:   8,
			expect: true,
		},
		{
			p:      WKE,
			rank:   7,
			expect: true,
		},
		{
			p:      WKE,
			rank:   6,
			expect: false,
		},
	}

	for _, test := range tests {
		got := NeedForcePromotion(test.p, test.rank)
		if got != test.expect {
			t.Errorf("NeedForcePromotion(%v, %v): want %v, got %v",
				test.p, test.rank, test.expect, got)
		}
	}
}
