package shogi

import (
	"fmt"
	"strconv"
)

type MoveData struct {
	Move
	Piece   Piece
	Capture Piece
	Same    bool
	Color   Color
	Ply     int
}

var InitialMoveData = MoveData{Move: InitialMove}
var ToryoMoveData = MoveData{Move: ToryoMove}

func NewMoveData(m Move, p *Position, before Square) MoveData {
	if m.IsInitialMove() {
		return InitialMoveData
	}
	if m.IsToryo() {
		return ToryoMoveData
	}
	if m.IsDropMove() {
		return MoveData{
			Move:  m,
			Color: p.Turn,
			Ply:   p.Ply,
		}
	}

	same := false
	if m.To == before {
		same = true
	}
	piece := p.Get(m.From)
	capture := p.Get(m.To)
	return MoveData{
		Move:    m,
		Piece:   piece,
		Capture: capture,
		Same:    same,
		Color:   p.Turn,
		Ply:     p.Ply + 1,
	}
}

func (m MoveData) KIF() string {
	if m.IsInitialMove() {
		return "開始局面"
	}
	if m.IsToryo() {
		return "投了"
	}
	to := m.To.KIF()
	if m.IsDropMove() {
		p := m.DropPieceType.KIF()
		return fmt.Sprintf("%s%s打", to, p)
	}
	from := strconv.Itoa(9-m.From.File()) + strconv.Itoa(m.From.Rank()+1)
	nari := ""
	if m.IsPromotion() {
		nari = "成"
	}
	p := m.Piece.PieceType().KIF()
	if m.Same {
		return fmt.Sprintf("同%s%s(%s)", p, nari, from)
	}
	return fmt.Sprintf("%v%v%v(%v)", to, p, nari, from)
}
