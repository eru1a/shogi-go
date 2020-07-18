package shogi

import (
	"errors"
	"fmt"
)

type Square struct {
	file, rank int
}

var NullSquare = Square{-1, -1}

func valid(file, rank int) bool {
	return 0 <= file && file < 9 && 0 <= rank && rank < 9
}

func NewSquare(file, rank int) (Square, error) {
	if !valid(file, rank) {
		return NullSquare,
			fmt.Errorf("out of bounds: (%d, %d)", file, rank)
	}
	return Square{file, rank}, nil
}

func (s Square) IsNull() bool {
	return s == NullSquare
}

func (s Square) File() int {
	return s.file
}

func (s Square) Rank() int {
	return s.rank
}

func (s Square) Add(file, rank int) (Square, error) {
	if s.IsNull() {
		return NullSquare, errors.New("attempted to add to the NullSquare")
	}
	newSq, err := NewSquare(s.file+file, s.rank+rank)
	if err != nil {
		return NullSquare, fmt.Errorf("Square(%v) Add(%d, %d): %v", s, file, rank, err)
	}
	return newSq, nil
}

func NewSquareFromUSI(usi string) (Square, error) {
	if len(usi) != 2 {
		return NullSquare,
			fmt.Errorf("should be at least two letters: %v", usi)
	}
	if usi[0] < '1' || usi[0] > '9' {
		return NullSquare,
			fmt.Errorf("first letter should be '0' <= c <= '8': %v", usi)
	}
	if usi[1] < 'a' || usi[1] > 'i' {
		return NullSquare,
			fmt.Errorf("second letter should be 'a' <= c <= 'h': %v", usi)
	}
	file := 8 - int(usi[0]-'1')
	rank := int(usi[1] - 'a')
	return NewSquare(file, rank)
}

func (s Square) USI() string {
	return fmt.Sprintf("%d%c", 9-s.file, s.rank+'a')
}

var (
	fileKifMap = map[int]string{
		0: "９",
		1: "８",
		2: "７",
		3: "６",
		4: "５",
		5: "４",
		6: "３",
		7: "２",
		8: "１",
	}
	rankKifMap = map[int]string{
		0: "一",
		1: "二",
		2: "三",
		3: "四",
		4: "五",
		5: "六",
		6: "七",
		7: "八",
		8: "九",
	}
)

func (s Square) KIF() string {
	if s.IsNull() {
		return "００"
	}
	return fileKifMap[s.file] + rankKifMap[s.rank]
}
