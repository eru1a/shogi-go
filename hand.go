package shogi

import (
	"fmt"
	"strconv"
	"strings"
)

type Hand struct {
	Black map[PieceType]int
	White map[PieceType]int
}

func NewHand() *Hand {
	h := &Hand{
		Black: make(map[PieceType]int),
		White: make(map[PieceType]int),
	}
	for pt := FU; pt <= HI; pt++ {
		h.Black[pt] = 0
		h.White[pt] = 0
	}
	return h
}

func (h *Hand) add(pt PieceType, c Color, n int) error {
	if c != Black && c != White {
		return fmt.Errorf("color should be black or white: %v", c)
	}

	handSide := h.Black
	if c == White {
		handSide = h.White
	}
	num, ok := handSide[pt]
	if !ok {
		return fmt.Errorf("piece type should be FU <= pt <= HI: %v", pt)
	}
	if num+n < 0 {
		return fmt.Errorf("%v's hand[%s]+%d < 0: %v", c, pt, n, h)
	}
	handSide[pt] += n
	return nil
}

func (h *Hand) Add(pt PieceType, c Color) error {
	return h.add(pt, c, 1)
}

func (h *Hand) Remove(pt PieceType, c Color) error {
	return h.add(pt, c, -1)
}

func (h *Hand) Get(pt PieceType, c Color) (int, error) {
	if c != Black && c != White {
		return -1, fmt.Errorf("color should be black or white: %v", c)
	}

	handSide := h.Black
	if c == White {
		handSide = h.White
	}
	n, ok := handSide[pt]
	if !ok {
		return -1, fmt.Errorf("piece type should be FU <= pt <= HI: %v", pt)
	}
	return n, nil
}

func NewHandFromSFEN(sfen string) (*Hand, error) {
	if len(sfen) == 0 {
		return nil, fmt.Errorf("hand sfen length should not be 0: %s", sfen)
	}

	if sfen == "-" {
		return NewHand(), nil
	}

	h := NewHand()

	isNumber := func(b byte) bool {
		return '0' <= b && b <= '9'
	}
	isAlpha := func(b byte) bool {
		return ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z')
	}

	num := 0
	for i := 0; i < len(sfen); i++ {
		switch {
		case isNumber(sfen[i]):
			num = num*10 + int(sfen[i]-'0')
		case isAlpha(sfen[i]):
			p, err := NewPieceFromUSI(string(sfen[i]))
			if err != nil {
				return nil, fmt.Errorf("invalid usi piece: %v", err)
			}
			n := 1
			if num != 0 {
				n = num
			}
			if err := h.add(p.PieceType(), p.Color(), n); err != nil {
				return nil, err
			}
			num = 0
		default:
			return nil, fmt.Errorf("unknown char: %s, %c", sfen, sfen[i])
		}
	}

	if num != 0 {
		return nil, fmt.Errorf("num > 0: %s", sfen)
	}
	return h, nil
}

func (h *Hand) SFEN() string {
	var b strings.Builder
	var w strings.Builder

	// 飛,角,金,銀,桂,香,歩の順が正しい?
	for pt := HI; pt >= FU; pt-- {
		if h.Black[pt] > 0 {
			if h.Black[pt] > 1 {
				b.WriteString(strconv.Itoa(h.Black[pt]))
			}
			b.WriteString(NewPiece(pt, Black).USI())
		}

		if h.White[pt] > 0 {
			if h.White[pt] > 1 {
				w.WriteString(strconv.Itoa(h.White[pt]))
			}
			w.WriteString(NewPiece(pt, White).USI())
		}
	}

	sfen := b.String() + w.String()
	if sfen == "" {
		return "-"
	}
	return sfen
}
