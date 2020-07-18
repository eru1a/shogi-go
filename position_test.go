package shogi

import (
	"reflect"
	"testing"
)

func TestPositionSFEN(t *testing.T) {
	okTests := []struct {
		sfen     string
		position *Position
	}{
		{
			sfen: "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1",
			position: &Position{
				Board: &Board{
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
				Hand: handFromArray(
					[7]int{0, 0, 0, 0, 0, 0, 0},
					[7]int{0, 0, 0, 0, 0, 0, 0}),
				Turn: Black,
				Ply:  0,
			},
		},
		{
			sfen: "3g2snl/R8/2+P1ppgp1/B1pp4p/G2n1S3/2PbP1P2/KP1+lkPN1P/6S2/L+r3G2L w 3Psn2p 98",
			position: &Position{
				Board: &Board{
					{___, ___, ___, WKI, ___, ___, WGI, WKE, WKY},
					{BHI, ___, ___, ___, ___, ___, ___, ___, ___},
					{___, ___, BTO, ___, WFU, WFU, WKI, WFU, ___},
					{BKA, ___, WFU, WFU, ___, ___, ___, ___, WFU},
					{BKI, ___, ___, WKE, ___, BGI, ___, ___, ___},
					{___, ___, BFU, WKA, BFU, ___, BFU, ___, ___},
					{BOU, BFU, ___, WNY, WOU, BFU, BKE, ___, BFU},
					{___, ___, ___, ___, ___, ___, BGI, ___, ___},
					{BKY, WRY, ___, ___, ___, BKI, ___, ___, BKY},
				},
				Hand: handFromArray(
					[7]int{3, 0, 0, 0, 0, 0, 0},
					[7]int{2, 0, 1, 1, 0, 0, 0}),
				Turn: White,
				Ply:  97,
			},
		},
	}

	for _, test := range okTests {
		// sfen -> position
		p, err := NewPositionFromSFEN(test.sfen)
		if err != nil {
			t.Errorf("NewPositionFromSFEN(%v): got error %v", test.sfen, err)
		}
		if !reflect.DeepEqual(p, test.position) {
			t.Errorf("NewPositionFromSFEN(%v): want %v, got %v", test.sfen, test.position, p)
		}

		// position -> sfen
		sfen := test.position.SFEN()
		if sfen != test.sfen {
			t.Errorf("%v.SFEN(): want %v, got %v", test.position, test.sfen, sfen)
		}
	}
}

func TestPositionMove(t *testing.T) {
	tests := []struct {
		startSFEN string
		position  *Position
		moves     []Move
		expected  *Position
	}{
		{
			startSFEN: "r6n1/6gk1/P2g1sspl/+B+Sp2ppl1/3pP2Np/3P1PP2/2+b1GG1S1/5K3/7RL b NLn7p 109",
			position: &Position{
				Board: &Board{
					{WHI, ___, ___, ___, ___, ___, ___, WKE, ___},
					{___, ___, ___, ___, ___, ___, WKI, WOU, ___},
					{BFU, ___, ___, WKI, ___, WGI, WGI, WFU, WKY},
					{BRY, BNG, WFU, ___, ___, WFU, WFU, WKY, ___},
					{___, ___, ___, WFU, BFU, ___, ___, BKE, WFU},
					{___, ___, ___, BFU, ___, BFU, BFU, ___, ___},
					{___, ___, WUM, ___, BKI, BKI, ___, BGI, ___},
					{___, ___, ___, ___, ___, BOU, ___, ___, ___},
					{___, ___, ___, ___, ___, ___, ___, BHI, BKY},
				},
				Hand: handFromArray(
					[7]int{0, 1, 1, 0, 0, 0, 0},
					[7]int{7, 0, 1, 0, 0, 0, 0}),
				Turn: Black,
				Ply:  108,
			},
			moves: []Move{
				NewNormalMove(Square{0, 2}, Square{0, 1}, true),
				NewNormalMove(Square{0, 0}, Square{0, 1}, false),
				NewNormalMove(Square{7, 4}, Square{6, 2}, false),
				NewNormalMove(Square{7, 3}, Square{7, 6}, true),
				NewDropMove(GI, Square{4, 1}),
				NewDropMove(GI, Square{4, 7}),
				NewNormalMove(Square{5, 7}, Square{4, 7}, false),
			},
			expected: &Position{
				Board: &Board{
					{___, ___, ___, ___, ___, ___, ___, WKE, ___},
					{WHI, ___, ___, ___, BGI, ___, WKI, WOU, ___},
					{___, ___, ___, WKI, ___, WGI, BKE, WFU, WKY},
					{BRY, BNG, WFU, ___, ___, WFU, WFU, ___, ___},
					{___, ___, ___, WFU, BFU, ___, ___, ___, WFU},
					{___, ___, ___, BFU, ___, BFU, BFU, ___, ___},
					{___, ___, WUM, ___, BKI, BKI, ___, WNY, ___},
					{___, ___, ___, ___, BOU, ___, ___, ___, ___},
					{___, ___, ___, ___, ___, ___, ___, BHI, BKY},
				},
				Hand: handFromArray(
					[7]int{0, 1, 1, 1, 0, 0, 0},
					[7]int{8, 0, 1, 0, 0, 0, 0}),
				Turn: White,
				Ply:  115,
			},
		},
	}

	for _, test := range tests {
		for _, m := range test.moves {
			if err := test.position.Move(m); err != nil {
				t.Fatal(err)
			}
		}

		if !reflect.DeepEqual(test.position, test.expected) {
			t.Errorf("position %v\n moves %v\n want %v\n got %v",
				test.startSFEN, test.moves, test.expected, test.position)
		}
	}
}

func TestLegalMoves(t *testing.T) {
	tests := []struct {
		msg          string
		sfen         string
		legalMovesNr int
	}{
		{"initial", "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1", 30},
		{"nifu", "8k/PP7/2P6/3P5/9/5P3/6P2/7P1/8P b P 1", 18},
		{"test", "9/4B1SGL/PN2R4/1N1P5/6N2/5+R3/3+B5/9/8L b - 1", 112},
		{"oute(1)", "9/9/3rR2B1/9/8b/4s4/4K4/3N5/9 b 2P 1", 8},
		{"oute(2)", "4r4/9/3R5/7B1/9/9/9/9/4K4 b G 1", 16},
		{"floodgate", "l+S3ks1R/3g2g1+L/4pp1p1/p5p2/1KPS1P1P1/P2p1BP2/+bg2P4/1P5R1/1N7 b 3N2L5Pgs 1", 153},
		// 以下はpython-shogiから
		{"stalemate", "+R+N+SGKG+S+N+R/+B+N+SG+LG+S+N+B/P+LPP+LPP+LP/1P2P2P1/9/9/9/9/6k2 b - 200", 0},
		{"checkmate by dropping FU(1)", "kn7/9/1G7/9/9/9/9/9/9 b P 1", 76},
		{"checkmate by dropping FU(2)", "kn7/9/9/1NN6/9/9/9/9/9 b P 1", 73},
		{"check by dropping FU(1)", "k8/9/9/9/9/9/9/9/9 b P 1", 72},
		{"check by dropping FU(2)", "kn7/1n7/9/9/9/9/9/9/9 b P 1", 71},
		{"check by dropping FU(3)", "kn7/9/9/1N7/9/9/9/9/9 b P 1", 73},
		// 82歩打で相手がstalemateになるけど王手でないので打ち歩詰めではない(?)
		{"check by dropping FU(4)", "k8/9/1S7/9/9/9/9/9/9 b P 1", 81},
		{"check by dropping FU(5)", "kg7/9/1G7/9/9/9/9/9/9 b P 1", 77},
	}

	for _, test := range tests {
		p, err := NewPositionFromSFEN(test.sfen)
		if err != nil {
			t.Fatal(err)
		}
		moves := p.LegalMoves()
		if len(moves) != test.legalMovesNr {
			t.Errorf("LegalMoves[%v]: want %v, got %v\nposition %v\n moves %v",
				test.msg, test.legalMovesNr, len(moves), p, moves)
		}
	}
}

func TestIsLegalMove(t *testing.T) {
	tests := []struct {
		sfen         string
		legalMoves   []string
		illegalMoves []string
	}{
		{
			sfen:         "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1",
			legalMoves:   []string{"7g7f", "3i4h", "1i1h"},
			illegalMoves: []string{"7g6f", "8h2b", "P*5e"},
		},
		{
			sfen:         "l+P6l/9/p1p1g1k1p/4pp3/1gP4pB/2r2P2P/P3P2PK/4+r1S2/5+p2L w 2S2N3Pb2gs2nlp 1",
			legalMoves:   []string{"3c3d", "3c2c", "B*2d", "L*2d"},
			illegalMoves: []string{"3c2d", "3c4b", "4h5h", "P*2d"},
		},
	}

	for _, test := range tests {
		p, _ := NewPositionFromSFEN(test.sfen)
		for _, move := range test.legalMoves {
			m, _ := NewMoveFromUSI(move)
			if !p.IsLegalMove(m) {
				t.Errorf("(%v).IsLegalMove(%v): should be legal", test.sfen, m)
			}
		}
		for _, move := range test.illegalMoves {
			m, _ := NewMoveFromUSI(move)
			if p.IsLegalMove(m) {
				t.Errorf("(%v).IsLegalMove(%v): should be illegal", test.sfen, m)
			}
		}
	}
}

func TestIsInCheck(t *testing.T) {
	tests := []struct {
		sfen  string
		check bool
	}{
		{"1r6l/3g2kg1/3sSpn2/4P1p1p/l1Pp1P3/2Sn1B2P/1P4B2/K1gG2+r2/LN6L b N8Psp 1", true},
		{"1r6l/3g2kg1/3sSpn2/4P1p1p/l1Pp1P3/2Sn1B2P/PP4B2/K1gG2+r2/LN6L w N7Psp 1", false},
		{"l2g1p1nl/1s4k2/p2p2ppp/9/1r3G1NP/2B2P1PL/P1pP2P2/3s1SSK1/L4G3 w R4Pbg2np 1", true},
		{"l2g1p1nl/1s4k2/p2p1bppp/9/1r3G1NP/2B2P1PL/P1pP2P2/3s1SSK1/L4G3 b R4Pg2np 1", false},
		{"4k4/9/9/9/9/9/4B4/9/1r2L4 w - 1", false},
		{"4k4/9/9/1B7/9/9/9/9/1r2L4 w - 1", true},
		{"k8/9/9/LK7/9/9/9/9/9 w - 1", true},
		{"k8/n8/9/LK7/9/9/9/9/9 b - 1", true},
	}

	for _, test := range tests {
		p, err := NewPositionFromSFEN(test.sfen)
		if err != nil {
			t.Fatal(err)
		}
		if p.IsInCheck() != test.check {
			t.Errorf("%v.IsInCheck(): want %v, got %v", test.sfen, test.check, p.IsInCheck())
		}
	}
}

func TestIsCheckmate(t *testing.T) {
	tests := []struct {
		sfen      string
		checkmate bool
	}{
		{"ln3k2l/3R5/p1p4p1/2s5p/6Pn1/4P1b1P/L+pPP3s1/3s3K1/1N2+s+r1NL b B4GP7p 1", true},
		{"lR2+R2+B1/+N3kg3/pPPp4p/3spsN2/5p1K1/Pp2S3P/n1N2P2L/3P5/L8 w B2GS6Pgl 1", true},
		{"ln7/2+R6/p1pppp1+Bp/1Nn6/L1S+b5/S1k6/P1LPP3P/1GG2P1P1/1N2KGS1L w GPrs5p 1", true},
		{"ln7/2+R6/p1pppp1+Bp/1Nn6/L1S+b5/S1k6/P1PPP3P/1GG2P1P1/1N2KGS1L w GPrs5p 1", false},
		{"8k/8P/7+R1/9/9/9/9/9/9 w - 1", true},
		{"8k/8P/7R1/9/9/9/9/9/9 w - 1", false},
		{"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1", false},
	}

	for _, test := range tests {
		p, err := NewPositionFromSFEN(test.sfen)
		if err != nil {
			t.Fatal(err)
		}
		if p.IsCheckmate() != test.checkmate {
			t.Errorf("%v.IsCheckmate(): want %v, got %v", test.sfen, test.checkmate, p.IsCheckmate())
		}
	}
}
