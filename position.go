package shogi

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: 千日手

type Position struct {
	Board *Board
	Turn  Color
	Hand  *Hand
	// 手数。これは今の局面が何手目かを指す。初期局面は0。
	// SFENは次の手が何手目かを指すので1手ずれることに注意。
	Ply int
}

func NewPosition() *Position {
	p, _ := NewPositionFromSFEN("lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1")
	return p
}

func (p *Position) Clone() *Position {
	clone, _ := NewPositionFromSFEN(p.SFEN())
	return clone
}

func NewPositionFromSFEN(sfen string) (*Position, error) {
	slice := strings.Split(sfen, " ")
	if len(slice) != 4 {
		return nil, fmt.Errorf("sfen length should be 4: %v", sfen)
	}

	board, err := NewBoardFromSFEN(slice[0])
	if err != nil {
		return nil, fmt.Errorf("invalid board sfen: %v", err)
	}

	turn, err := NewColorFromUSI(slice[1])
	if err != nil {
		return nil, fmt.Errorf("invalid color sfen: %v", err)
	}

	hand, err := NewHandFromSFEN(slice[2])
	if err != nil {
		return nil, fmt.Errorf("invalid hand sfen: %v", err)
	}

	ply, err := strconv.Atoi(slice[3])
	if err != nil {
		return nil, fmt.Errorf("invalid turn sfen: %v", err)
	}

	return &Position{board, turn, hand, ply - 1}, nil
}

func (p *Position) SFEN() string {
	return fmt.Sprintf("%s %s %s %d", p.Board.SFEN(), p.Turn.USI(), p.Hand.SFEN(), p.Ply+1)
}

func (p *Position) Get(s Square) Piece {
	return p.Board[s.rank][s.file]
}

func (p *Position) Set(s Square, piece Piece) {
	p.Board[s.rank][s.file] = piece
}

func (p *Position) HandAdd(pt PieceType, c Color) error {
	return p.Hand.Add(pt, c)
}

func (p *Position) HandRemove(pt PieceType, c Color) error {
	return p.Hand.Remove(pt, c)
}

func (p *Position) HandGet(pt PieceType, c Color) (int, error) {
	return p.Hand.Get(pt, c)
}

func (p *Position) Move(m Move) error {
	if p.IsLegalMove(m) {
		if err := p.move(m); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("illegal move: %v, %v", p, m)
}

func (p *Position) move(m Move) error {
	if m.IsNormalMove() {
		if err := p.normalMove(m); err != nil {
			return err
		}
	} else if m.IsDropMove() {
		if err := p.dropMove(m); err != nil {
			return err
		}
	}

	p.Turn = p.Turn.Inv()
	p.Ply++

	return nil
}

func (p *Position) normalMove(m Move) error {
	moved := p.Get(m.From)
	if moved == NO_PIECE {
		return fmt.Errorf("attempted to move a none piece: %v, %v", p, m)
	} else if moved.Color() != p.Turn {
		return fmt.Errorf("attempted to move an enemy piece: %v, %v", p, m)
	}
	if captured := p.Get(m.To); captured != NO_PIECE {
		if captured.Color() == p.Turn {
			return fmt.Errorf("attempted to capture an own piece: %v, %v", p, m)
		}
		if err := p.HandAdd(captured.PieceType().Demote(), p.Turn); err != nil {
			return fmt.Errorf("attempted to add an invalid piece to hand: %v", err)
		}
	}
	if m.Promotion {
		if moved.IsPromoted() {
			return fmt.Errorf("attempted to promote a promoted piece: %v, %v", p, m)
		}
		moved = moved.Promote()
	}
	p.Set(m.From, NO_PIECE)
	p.Set(m.To, moved)

	return nil
}

func (p *Position) dropMove(m Move) error {
	if p.Get(m.To) != NO_PIECE {
		return fmt.Errorf("attempted to drop a piece to nonempty square: %v, %v", p, m)
	}
	if err := p.HandRemove(m.DropPieceType, p.Turn); err != nil {
		return fmt.Errorf("attempted to remove an invalid piece from hand: %v", err)
	}
	p.Set(m.To, NewPiece(m.DropPieceType, p.Turn))

	return nil
}

func (p *Position) String() string {
	return p.SFEN()
}

func (p *Position) IsLegalMove(m Move) bool {
	for _, move := range p.LegalMoves() {
		if m == move {
			return true
		}
	}
	return false
}

// 王手放置、打ち歩詰め
func (p *Position) isForbiddenMove(m Move) bool {
	clone := p.Clone()
	clone.move(m)
	// 動かした局面でこちら側が王手なら非合法手
	if clone.isInCheckByColor(p.Turn) {
		return false
	}
	// 打ち歩詰め
	if m.IsDropMove() && m.DropPieceType == FU && clone.IsInCheck() && clone.IsCheckmate() {
		return false
	}
	return true
}

func (p *Position) LegalMoves() []Move {
	pseudo := p.pseudoLegalMoves(p.Turn)
	moves := []Move{}
	for _, move := range pseudo {
		if p.isForbiddenMove(move) {
			moves = append(moves, move)
		}
	}
	return moves
}

type offsets [][2]int

var pieceOffsets = map[Piece]offsets{
	BFU: {{0, -1}},
	BKY: {{0, -1}},
	BKE: {{-1, -2}, {1, -2}},
	BGI: {{-1, -1}, {0, -1}, {1, -1}, {-1, 1}, {1, 1}},
	BKI: {{-1, 0}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {0, 1}},
	BKA: {{-1, -1}, {1, -1}, {1, 1}, {-1, 1}},
	BHI: {{-1, 0}, {0, -1}, {1, 0}, {0, 1}},
	BOU: {{-1, 0}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}},
	WFU: {{0, 1}},
	WKY: {{0, 1}},
	WKE: {{-1, 2}, {1, 2}},
	WGI: {{-1, 1}, {0, 1}, {1, 1}, {-1, -1}, {1, -1}},
	WKI: {{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {0, -1}},
	WKA: {{-1, 1}, {1, 1}, {1, -1}, {-1, -1}},
	WHI: {{-1, 0}, {0, 1}, {1, 0}, {0, -1}},
	WOU: {{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}},
}

func canPromote(from, to Square, piece Piece) bool {
	if piece.IsPromoted() {
		return false
	}
	if piece.Color() == Black {
		return from.rank <= 2 || to.rank <= 2
	} else {
		return from.rank >= 6 || to.rank >= 6
	}
}

// 疑似合法手を生成。
// 具体的には王手放置や打ち歩詰めが含まれる。
func (p *Position) pseudoLegalMoves(c Color) []Move {
	moves := []Move{}
	for rank := 0; rank < 9; rank++ {
		for file := 0; file < 9; file++ {
			s := Square{file, rank}
			piece := p.Get(s)
			if piece == NO_PIECE || piece.Color() != c {
				continue
			}
			switch piece.PieceType() {
			case FU, KE, GI, KI, OU:
				moves = append(moves, p.generateHoppingMoves(s, c, pieceOffsets[piece])...)
			case TO, NY, NK, NG:
				moves = append(moves, p.generateHoppingMoves(s, c, pieceOffsets[NewPiece(KI, c)])...)
			case KY, KA, HI:
				moves = append(moves, p.generateSlidingMoves(s, c, pieceOffsets[piece])...)
			case UM:
				moves = append(moves, p.generateSlidingMoves(s, c, pieceOffsets[BKA])...)
				moves = append(moves, p.generateHoppingMoves(s, c, pieceOffsets[BHI])...)
			case RY:
				moves = append(moves, p.generateSlidingMoves(s, c, pieceOffsets[BHI])...)
				moves = append(moves, p.generateHoppingMoves(s, c, pieceOffsets[BKA])...)
			default:
				panic(fmt.Sprintf("invalid piece: %v, %v, %v", p, s, piece))
			}
		}
	}
	moves = append(moves, p.generateDroppingMoves(c)...)
	return moves
}

func (p *Position) generateHoppingMoves(from Square, c Color, offsets offsets) []Move {
	moves := []Move{}
	fromPiece := p.Get(from)
	for _, offset := range offsets {
		to, err := from.Add(offset[0], offset[1])
		if err != nil {
			continue
		}
		toPiece := p.Get(to)
		if toPiece == NO_PIECE || toPiece.Color() != c {
			if !forbiddenRank(fromPiece, to.rank) {
				moves = append(moves, NewNormalMove(from, to, false))
			}
			if canPromote(from, to, fromPiece) {
				moves = append(moves, NewNormalMove(from, to, true))
			}
		}
	}
	return moves
}

func (p *Position) generateSlidingMoves(from Square, c Color, offsets offsets) []Move {
	moves := []Move{}
	fromPiece := p.Get(from)
	for _, offset := range offsets {
		to := from
		var err error
	loop:
		for {
			to, err = to.Add(offset[0], offset[1])
			if err != nil {
				break loop
			}
			toPiece := p.Get(to)
			switch {
			case toPiece == NO_PIECE:
				if !forbiddenRank(fromPiece, to.rank) {
					moves = append(moves, NewNormalMove(from, to, false))
				}
				if canPromote(from, to, fromPiece) {
					moves = append(moves, NewNormalMove(from, to, true))
				}
			case toPiece.Color() != c:
				if !forbiddenRank(fromPiece, to.rank) {
					moves = append(moves, NewNormalMove(from, to, false))
				}
				if canPromote(from, to, fromPiece) {
					moves = append(moves, NewNormalMove(from, to, true))
				}
				break loop
			case toPiece.Color() == c:
				break loop
			}
		}
	}
	return moves
}

func (p *Position) generateDroppingMoves(c Color) []Move {
	var isNifu [9]bool
	for rank := 0; rank < 9; rank++ {
		for file := 0; file < 9; file++ {
			piece := p.Get(Square{file, rank})
			if piece.PieceType() == FU && piece.Color() == c {
				isNifu[file] = true
			}
		}
	}

	moves := []Move{}
	for rank := 0; rank < 9; rank++ {
		for file := 0; file < 9; file++ {
			for pt := FU; pt <= HI; pt++ {
				if n, _ := p.HandGet(pt, c); n == 0 {
					continue
				}
				if forbiddenRank(NewPiece(pt, c), rank) {
					continue
				}
				if pt == FU && isNifu[file] {
					continue
				}
				s := Square{file, rank}
				if p.Get(s) != NO_PIECE {
					continue
				}
				moves = append(moves, NewDropMove(pt, s))
			}
		}
	}

	return moves
}

func (p *Position) findKing(c Color) (Square, bool) {
	for rank := 0; rank < 9; rank++ {
		for file := 0; file < 9; file++ {
			s := Square{file, rank}
			piece := p.Get(s)
			if piece.PieceType() == OU && piece.Color() == c {
				return s, true
			}
		}
	}
	return NullSquare, false
}

func (p *Position) isInCheckByColor(c Color) bool {
	s, ok := p.findKing(c)
	if !ok {
		return false
	}
	moves := p.pseudoLegalMoves(c.Inv())
	for _, move := range moves {
		if move.To == s {
			return true
		}
	}
	return false
}

func (p *Position) IsInCheck() bool {
	return p.isInCheckByColor(p.Turn)
}

func (p *Position) IsCheckmate() bool {
	return len(p.LegalMoves()) == 0
}
