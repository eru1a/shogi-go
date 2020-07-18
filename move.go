package shogi

import (
	"fmt"
)

type MoveKind uint8

const (
	NullMoveKind MoveKind = iota
	InitialMoveKind
	NormalMoveKind
	DropMoveKind
	ToryoMoveKind
)

type Move struct {
	Kind          MoveKind
	From          Square
	To            Square
	Promotion     bool
	DropPieceType PieceType
}

var NullMove = Move{Kind: NullMoveKind, From: NullSquare, To: NullSquare}
var InitialMove = Move{Kind: InitialMoveKind, From: NullSquare, To: NullSquare}
var ToryoMove = Move{Kind: ToryoMoveKind, From: NullSquare, To: NullSquare}

func NewNormalMove(from, to Square, promotion bool) Move {
	return Move{
		Kind:          NormalMoveKind,
		From:          from,
		To:            to,
		Promotion:     promotion,
		DropPieceType: NO_PIECE_TYPE,
	}
}

func NewDropMove(dropPieceType PieceType, to Square) Move {
	return Move{
		Kind:          DropMoveKind,
		From:          NullSquare,
		To:            to,
		Promotion:     false,
		DropPieceType: dropPieceType,
	}
}

func (m Move) IsNullMove() bool {
	return m.Kind == NullMoveKind
}

func (m Move) IsInitialMove() bool {
	return m.Kind == InitialMoveKind
}

func (m Move) IsToryo() bool {
	return m.Kind == ToryoMoveKind
}

func (m Move) IsNormalMove() bool {
	return m.Kind == NormalMoveKind
}

func (m Move) IsDropMove() bool {
	return m.Kind == DropMoveKind
}

func (m Move) IsSpecialMove() bool {
	switch m.Kind {
	case NormalMoveKind, DropMoveKind:
		return false
	default:
		return true
	}
}

func (m Move) IsPromotion() bool {
	return m.Promotion
}

func NewMoveFromUSI(usi string) (Move, error) {
	if !(len(usi) == 4 || len(usi) == 5) {
		return NullMove, fmt.Errorf("length should be at least 4 or 5: %s", usi)
	}
	if usi[1] == '*' {
		return newDropMoveFromUSI(usi)
	}
	return newNormalMoveFromUSI(usi)
}

func newNormalMoveFromUSI(usi string) (Move, error) {
	from, err := NewSquareFromUSI(usi[0:2])
	if err != nil {
		return NullMove, fmt.Errorf("from usi: %v", err)
	}
	to, err := NewSquareFromUSI(usi[2:4])
	if err != nil {
		return NullMove, fmt.Errorf("to usi: %v", err)
	}
	promotion := false
	if len(usi) == 5 {
		if usi[4] != '+' {
			return NullMove,
				fmt.Errorf("fifth char should be '+': %s", usi)
		}
		promotion = true
	}
	return NewNormalMove(from, to, promotion), nil
}

func newDropMoveFromUSI(usi string) (Move, error) {
	if usi[1] != '*' {
		return NullMove, fmt.Errorf("second char should be '*': %s", usi)
	}

	pt, err := NewPieceTypeFromUSI(usi[0:1])
	if err != nil {
		return NullMove, fmt.Errorf("piece type: %v", err)
	}
	to, err := NewSquareFromUSI(usi[2:4])
	if err != nil {
		return NullMove, fmt.Errorf("to usi: %v", err)
	}
	return NewDropMove(pt, to), nil
}

func (m Move) USI() string {
	if m.IsNormalMove() {
		return m.normalMoveUSI()
	} else if m.IsDropMove() {
		return m.dropMoveUSI()
	}
	return "NULL_MOVE"
}

func (m Move) normalMoveUSI() string {
	promotion := ""
	if m.Promotion {
		promotion = "+"
	}
	return fmt.Sprintf("%s%s%s", m.From.USI(), m.To.USI(), promotion)
}

func (m Move) dropMoveUSI() string {
	return fmt.Sprintf("%s*%s", m.DropPieceType.USI(), m.To.USI())
}

func (m Move) String() string {
	return m.USI()
}
