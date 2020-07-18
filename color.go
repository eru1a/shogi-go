package shogi

import "fmt"

type Color int

const (
	NO_COLOR Color = iota
	Black
	White
)

func (c Color) Inv() Color {
	switch c {
	case Black:
		return White
	case White:
		return Black
	default:
		return NO_COLOR
	}
}

func NewColorFromUSI(usi string) (Color, error) {
	switch usi {
	case "b":
		return Black, nil
	case "w":
		return White, nil
	default:
		return NO_COLOR, fmt.Errorf("invalid usi: %s", usi)
	}
}

func (c Color) USI() string {
	switch c {
	case Black:
		return "b"
	case White:
		return "w"
	default:
		return "_"
	}
}

func (c Color) String() string {
	return c.USI()
}
