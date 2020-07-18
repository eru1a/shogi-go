package shogi

import (
	"reflect"
	"testing"
)

var ___ = NO_PIECE

func TestBoardUSI(t *testing.T) {
	okTests := []struct {
		sfen  string
		board *Board
	}{
		{
			sfen: "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL",
			board: &Board{
				{WKY, WKE, WGI, WKI, WOU, WKI, WGI, WKE, WKY},
				{___, WHI, ___, ___, ___, ___, ___, WKA, ___},
				{WFU, WFU, WFU, WFU, WFU, WFU, WFU, WFU, WFU},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{BFU, BFU, BFU, BFU, BFU, BFU, BFU, BFU, BFU},
				{___, BKA, ___, ___, ___, ___, ___, BHI, ___},
				{BKY, BKE, BGI, BKI, BOU, BKI, BGI, BKE, BKY},
			},
		},
		{
			sfen: "9/9/9/9/9/9/9/9/9",
			board: &Board{
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
			},
		},
		{
			sfen: "l4n3/3+P5/p3p4/6Sp1/4SG2k/P3P3p/4+bPK+r1/+p5P2/9",
			board: &Board{
				{WKY, ___, ___, ___, ___, WKE, ___, ___, ___},
				{___, ___, ___, BTO, ___, ___, ___, ___, ___},
				{WFU, ___, ___, ___, WFU, ___, ___, ___, ___},
				{___, ___, ___, ___, ___, ___, BGI, WFU, ___},
				{___, ___, ___, ___, BGI, BKI, ___, ___, WOU},
				{BFU, ___, ___, ___, BFU, ___, ___, ___, WFU},
				{___, ___, ___, ___, WUM, BFU, BOU, WRY, ___},
				{WTO, ___, ___, ___, ___, ___, BFU, ___, ___},
				{___, ___, ___, ___, ___, ___, ___, ___, ___},
			},
		},
	}

	for _, test := range okTests {
		b, err := NewBoardFromSFEN(test.sfen)
		if err != nil {
			t.Errorf("NewBoardFromSFEN(%v): erorr %v", test.sfen, err)
		}
		if !reflect.DeepEqual(b, test.board) {
			t.Errorf("NewBoardFromSFEN(%v): want %v, got %v", test.sfen, test.board, b)
		}

		sfen := test.board.SFEN()
		if sfen != test.sfen {
			t.Errorf("%v.SFEN(): want %v, got %v", test.board, test.sfen, sfen)
		}
	}

	ngTests := []struct {
		sfen string
	}{
		{""},
		{"9"},
		{"////////"},
		{"lnsgkgsnl/1r5b1/ppppApppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL"},
		{"lnsgkgsnl/1r5b1/pppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL"},
		{"lnsgkgsnl/1r5b1/pppp0ppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL"},
		{"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/"},
		{"9/9/9/9/8/9/9/9/9"},
		{"9/9/9/9/10/9/9/9/9"},
	}

	for _, test := range ngTests {
		b, err := NewBoardFromSFEN(test.sfen)
		if err == nil {
			t.Errorf("NewBoardFromSFEN(%v): want erorr, got %v", test.sfen, b)
		}
	}
}
