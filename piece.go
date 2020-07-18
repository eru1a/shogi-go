package shogi

import "fmt"

type PieceType int

const (
	promote = 8
	white   = 16
)

const (
	NO_PIECE_TYPE PieceType = iota
	FU
	KY
	KE
	GI
	KI
	KA
	HI
	OU

	TO PieceType = FU | promote
	NY PieceType = KY | promote
	NK PieceType = KE | promote
	NG PieceType = GI | promote
	UM PieceType = KA | promote
	RY PieceType = HI | promote
)

var (
	mapPieceTypeUSI = map[string]PieceType{
		"P":  FU,
		"L":  KY,
		"N":  KE,
		"S":  GI,
		"G":  KI,
		"B":  KA,
		"R":  HI,
		"K":  OU,
		"+P": TO,
		"+L": NY,
		"+N": NK,
		"+S": NG,
		"+B": UM,
		"+R": RY,
	}
	mapUSIPieceType = map[PieceType]string{
		FU: "P",
		KY: "L",
		KE: "N",
		GI: "S",
		KI: "G",
		KA: "B",
		HI: "R",
		OU: "K",
		TO: "+P",
		NY: "+L",
		NK: "+N",
		NG: "+S",
		UM: "+B",
		RY: "+R",
	}
	mapKIFPieceType = map[PieceType]string{
		FU: "歩",
		KY: "香",
		KE: "桂",
		GI: "銀",
		KI: "金",
		KA: "角",
		HI: "飛",
		OU: "王",
		TO: "と",
		NY: "杏",
		NK: "圭",
		NG: "全",
		UM: "馬",
		RY: "龍",
	}
)

func NewPieceTypeFromUSI(usi string) (PieceType, error) {
	pt, ok := mapPieceTypeUSI[usi]
	if !ok {
		return NO_PIECE_TYPE, fmt.Errorf("invalid usi: %s", usi)
	}
	return pt, nil
}

func (pt PieceType) USI() string {
	usi, ok := mapUSIPieceType[pt]
	if !ok {
		return "_"
	}
	return usi
}

func (pt PieceType) KIF() string {
	kif, ok := mapKIFPieceType[pt]
	if !ok {
		return "・"
	}
	return kif
}

func (pt PieceType) Promote() PieceType {
	switch pt {
	case NO_PIECE_TYPE, KI, OU:
		return pt
	}
	return pt | promote
}

func (pt PieceType) Demote() PieceType {
	if pt == OU {
		return pt
	}
	return pt & ^promote
}

func (pt PieceType) IsPromoted() bool {
	return pt == pt.Promote()
}

func (pt PieceType) String() string {
	return pt.USI()
}

type Piece int

const (
	NO_PIECE = Piece(NO_PIECE_TYPE)
	BFU      = Piece(FU)
	BKY      = Piece(KY)
	BKE      = Piece(KE)
	BGI      = Piece(GI)
	BKI      = Piece(KI)
	BKA      = Piece(KA)
	BHI      = Piece(HI)
	BOU      = Piece(OU)
	BTO      = Piece(TO)
	BNY      = Piece(NY)
	BNK      = Piece(NK)
	BNG      = Piece(NG)
	BUM      = Piece(UM)
	BRY      = Piece(RY)
	WFU      = Piece(FU | white)
	WKY      = Piece(KY | white)
	WKE      = Piece(KE | white)
	WGI      = Piece(GI | white)
	WKI      = Piece(KI | white)
	WKA      = Piece(KA | white)
	WHI      = Piece(HI | white)
	WOU      = Piece(OU | white)
	WTO      = Piece(TO | white)
	WNY      = Piece(NY | white)
	WNK      = Piece(NK | white)
	WNG      = Piece(NG | white)
	WUM      = Piece(UM | white)
	WRY      = Piece(RY | white)
)

var (
	mapPieceUSI = map[string]Piece{
		"P":  BFU,
		"L":  BKY,
		"N":  BKE,
		"S":  BGI,
		"G":  BKI,
		"B":  BKA,
		"R":  BHI,
		"K":  BOU,
		"+P": BTO,
		"+L": BNY,
		"+N": BNK,
		"+S": BNG,
		"+B": BUM,
		"+R": BRY,
		"p":  WFU,
		"l":  WKY,
		"n":  WKE,
		"s":  WGI,
		"g":  WKI,
		"b":  WKA,
		"r":  WHI,
		"k":  WOU,
		"+p": WTO,
		"+l": WNY,
		"+n": WNK,
		"+s": WNG,
		"+b": WUM,
		"+r": WRY,
	}
	mapUSIPiece = map[Piece]string{
		BFU: "P",
		BKY: "L",
		BKE: "N",
		BGI: "S",
		BKI: "G",
		BKA: "B",
		BHI: "R",
		BOU: "K",
		BTO: "+P",
		BNY: "+L",
		BNK: "+N",
		BNG: "+S",
		BUM: "+B",
		BRY: "+R",
		WFU: "p",
		WKY: "l",
		WKE: "n",
		WGI: "s",
		WKI: "g",
		WKA: "b",
		WHI: "r",
		WOU: "k",
		WTO: "+p",
		WNY: "+l",
		WNK: "+n",
		WNG: "+s",
		WUM: "+b",
		WRY: "+r",
	}
)

func NewPiece(pt PieceType, c Color) Piece {
	if c == White {
		return Piece(pt | white)
	}
	return Piece(pt)
}

func (p Piece) PieceType() PieceType {
	return PieceType(p & ^white)
}

func (p Piece) Color() Color {
	if p == NO_PIECE {
		return NO_COLOR
	} else if p&white == 0 {
		return Black
	}
	return White
}

func NewPieceFromUSI(usi string) (Piece, error) {
	p, ok := mapPieceUSI[usi]
	if !ok {
		return NO_PIECE, fmt.Errorf("invalid usi: %s", usi)
	}
	return p, nil
}

func (p Piece) USI() string {
	usi, ok := mapUSIPiece[p]
	if !ok {
		return "_"
	}
	return usi
}

func (p Piece) Promote() Piece {
	if p == NO_PIECE {
		return NO_PIECE
	}
	return NewPiece(p.PieceType().Promote(), p.Color())
}

func (p Piece) Demote() Piece {
	return NewPiece(p.PieceType().Demote(), p.Color())
}

func (p Piece) IsPromoted() bool {
	return p == p.Promote()
}

func (p Piece) String() string {
	return p.USI()
}

func forbiddenRank(p Piece, rank int) bool {
	switch {
	case p == BFU && rank == 0,
		p == WFU && rank == 8,
		p == BKY && rank == 0,
		p == WKY && rank == 8,
		p == BKE && rank <= 1,
		p == WKE && rank >= 7:
		return true
	}
	return false
}

func CanPromote(p Piece, from, to Square) bool {
	if p == NO_PIECE {
		return false
	}
	if from == NullSquare || to == NullSquare {
		return false
	}
	pt := p.PieceType()
	if p.IsPromoted() || pt == KI || pt == OU {
		return false
	}

	c := p.Color()

	switch {
	case c == Black && (from.rank <= 2 || to.rank <= 2):
		return true
	case c == White && (from.rank >= 6 || to.rank >= 6):
		return true
	}
	return false
}

func NeedForcePromotion(p Piece, rank int) bool {
	return forbiddenRank(p, rank)
}
