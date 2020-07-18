package engine

import (
	"reflect"
	"testing"
)

func TestNewUSIInfo(t *testing.T) {
	tests := []struct {
		line string
		info USIInfo
	}{
		{
			line: "info multipv 1 score cp 40 depth 64 pv 7g7f 3c3d 2g2f 8c8d 2f2e 4a3b 6i7h 8d8e (100.00%)",
			info: USIInfo{
				MultiPv: 1,
				IsCp:    true,
				Depth:   64,
				ScoreCp: 40,
				Pv:      []string{"7g7f", "3c3d", "2g2f", "8c8d", "2f2e", "4a3b", "6i7h", "8d8e"},
			},
		},
		{
			line: "info depth 17 seldepth 35 score cp -9236 upperbound nodes 9089473 nps 852111 hashfull 932 time 10667 pv 6h7g 5g6i+",
			info: USIInfo{
				MultiPv:    1,
				IsCp:       true,
				Depth:      17,
				SelDepth:   35,
				ScoreCp:    -9236,
				Upperbound: true,
				Nodes:      9089473,
				Nps:        852111,
				HashFull:   932,
				Time:       10667,
				Pv:         []string{"6h7g", "5g6i+"},
			},
		},
		{
			line: "info depth 23 seldepth 16 score mate 15 multipv 3 nodes 8034018 nps 936911 hashfull 995 time 8575 pv N*3c 2a3c G*5a",
			info: USIInfo{
				MultiPv:   3,
				IsMate:    true,
				Depth:     23,
				SelDepth:  16,
				ScoreMate: 15,
				Nodes:     8034018,
				Nps:       936911,
				HashFull:  995,
				Time:      8575,
				Pv:        []string{"N*3c", "2a3c", "G*5a"},
			},
		},
	}

	for _, test := range tests {
		info, err := NewUSIInfo(test.line)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(info, test.info) {
			t.Errorf("NewUSIInfo(%s):\n\twant\t%v\n\tgot\t%v", test.line, test.info, info)
		}
	}
}

func TestNewUSIBestMove(t *testing.T) {
	tests := []struct {
		line     string
		bestmove USIBestMove
	}{
		{
			line: "bestmove 7g7f",
			bestmove: USIBestMove{
				BestMove: "7g7f",
			},
		},
		{
			line: "bestmove 7g7f ponder 3c3d",
			bestmove: USIBestMove{
				BestMove: "7g7f",
				Ponder:   "3c3d",
			},
		},
	}

	for _, test := range tests {
		bestmove, err := NewUSIBestMove(test.line)
		if err != nil {
			panic(err)
		}
		if bestmove != test.bestmove {
			t.Errorf("NewUSIBestMove(%s):\n\twant\t%v\n\tgot\t%v", test.line, test.bestmove, bestmove)
		}
	}
}
