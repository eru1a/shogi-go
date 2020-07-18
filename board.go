package shogi

import (
	"fmt"
	"strconv"
	"strings"
)

type Board [9][9]Piece

func (b *Board) Set(s Square, p Piece) {
	b[s.rank][s.file] = p
}

func (b *Board) Get(s Square) Piece {
	return b[s.rank][s.file]
}

func NewBoardFromSFEN(sfen string) (*Board, error) {
	var b Board

	rank := 0
	file := 0
	promotion := false
	for i := 0; i < len(sfen); i++ {
		if sfen[i] == '/' {
			if file != 9 {
				return nil, fmt.Errorf("file != 9: %v, %v", sfen, file)
			}
			file = 0
			rank++
			continue
		}
		if sfen[i] == '+' {
			promotion = true
			continue
		}
		if file > 8 {
			return nil, fmt.Errorf("file is out of bounds: %v, %v, %v-th", sfen, file, i)
		}
		if rank > 8 {
			return nil, fmt.Errorf("rank is out of bounds: %v, %v, %v-th", sfen, rank, i)
		}
		if '1' <= sfen[i] && sfen[i] <= '9' {
			if i+1 < len(sfen) && ('0' <= sfen[i+1] && sfen[i+1] <= '9') {
				return nil, fmt.Errorf("n > 9: %v", sfen)
			}
			file += int(sfen[i] - '0')
			continue
		}
		p, err := NewPieceFromUSI(string(sfen[i]))
		if err != nil {
			return nil, fmt.Errorf("%d-th char is invalid: %v", i, sfen)
		}
		if promotion {
			p = p.Promote()
			promotion = false
		}
		b[rank][file] = p
		file++
	}
	if file != 9 || rank != 8 {
		return nil, fmt.Errorf("file != 9 || rank != 8: %v, %v, %v", sfen, file, rank)
	}

	return &b, nil
}

func (b *Board) SFEN() string {
	var usi strings.Builder

	nBlank := 0
	for rank := 0; rank < 9; rank++ {
		for file := 0; file < 9; file++ {
			p := b[rank][file]
			if p == NO_PIECE {
				nBlank++
				continue
			}
			if nBlank > 0 {
				usi.WriteString(strconv.Itoa(nBlank))
				nBlank = 0
			}
			usi.WriteString(p.USI())
		}
		if nBlank > 0 {
			usi.WriteString(strconv.Itoa(nBlank))
			nBlank = 0
		}
		if rank != 8 {
			usi.WriteString("/")
		}
	}

	return usi.String()
}

func (b *Board) string() string {
	var s strings.Builder
	for rank := 0; rank < 9; rank++ {
		for file := 0; file < 9; file++ {
			p := b[rank][file]
			if p == NO_PIECE {
				s.WriteString(" .")
			} else {
				s.WriteString(fmt.Sprintf("%2s", p.USI()))
			}
		}
		s.WriteString("\n")
	}
	return s.String()
}

func (b *Board) String() string {
	return b.SFEN()
}
