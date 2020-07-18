package shogi

import (
	"reflect"
	"testing"
)

func TestGame(t *testing.T) {
	tests := []struct {
		usiMoves []string
		sfens    []string
	}{
		{
			usiMoves: []string{
				"7g7f",
				"3c3d",
				"8h2b+",
				"3a2b",
				"B*5e",
			},
			sfens: []string{
				"lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1",
				"lnsgkgsnl/1r5b1/ppppppppp/9/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL w - 2",
				"lnsgkgsnl/1r5b1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL b - 3",
				"lnsgkgsnl/1r5+B1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/7R1/LNSGKGSNL w B 4",
				"lnsgkg1nl/1r5s1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/7R1/LNSGKGSNL b Bb 5",
				"lnsgkg1nl/1r5s1/pppppp1pp/6p2/4B4/2P6/PP1PPPPPP/7R1/LNSGKGSNL w b 6",
			},
		},
	}

	for _, test := range tests {
		tree := NewGameTree()
		for _, usi := range test.usiMoves {
			move, err := NewMoveFromUSI(usi)
			if err != nil {
				t.Fatal(err)
			}
			if err := tree.Move(move); err != nil {
				t.Fatal(err)
			}
		}
		sfens := []string{}
		for node := tree.Root; node != nil; node = node.Next {
			sfens = append(sfens, node.Position.SFEN())
		}
		if !reflect.DeepEqual(sfens, test.sfens) {
			t.Errorf("\nwant %v\ngot %v", test.sfens, sfens)
		}
	}

}
