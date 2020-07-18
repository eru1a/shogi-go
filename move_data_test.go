package shogi

import "testing"

func newMoveData(usi string, sfen string, before Square) MoveData {
	p, _ := NewPositionFromSFEN(sfen)
	m, _ := NewMoveFromUSI(usi)
	return NewMoveData(m, p, before)
}

func TestMoveDataKIF(t *testing.T) {
	tests := []struct {
		kif      string
		moveData MoveData
	}{
		{
			kif:      "開始局面",
			moveData: InitialMoveData,
		},
		{
			kif:      "投了",
			moveData: ToryoMoveData,
		},
		{
			kif:      "７六歩(77)",
			moveData: newMoveData("7g7f", "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1", NullSquare),
		},
		{
			kif:      "２二角成(88)",
			moveData: newMoveData("8h2b+", "lnsgkgsnl/1r5b1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL b - 3", Square{6, 3}),
		},
		{
			kif:      "同銀(31)",
			moveData: newMoveData("3a2b", "lnsgkgsnl/1r5+B1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/7R1/LNSGKGSNL w B 4", Square{7, 1}),
		},
		{
			kif:      "５五角打",
			moveData: newMoveData("B*5e", "lnsgkg1nl/1r5s1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/7R1/LNSGKGSNL b Bb 5", Square{7, 1}),
		},
	}

	for _, test := range tests {
		kif := test.moveData.KIF()
		if kif != test.kif {
			t.Errorf("%v.KIF(): want %v, got %v", test.moveData, test.kif, kif)
		}
	}
}
