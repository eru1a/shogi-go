package shogi

import "fmt"

type GameTree struct {
	Root    *GameNode
	Current *GameNode
}

type GameNode struct {
	Prev     *GameNode
	Next     *GameNode
	Position *Position
	MoveData MoveData
}

func NewGameTree() *GameTree {
	p := NewPosition()
	n := NewGameNode(nil, p, InitialMoveData)
	return &GameTree{
		Root:    n,
		Current: n,
	}
}

func NewGameTreeFromSFEN(sfen string) (*GameTree, error) {
	p, err := NewPositionFromSFEN(sfen)
	if err != nil {
		return nil, fmt.Errorf("NewGameTreeFromSFEN: %v", err)
	}
	n := NewGameNode(nil, p, InitialMoveData)
	return &GameTree{Root: n, Current: n}, nil
}

func NewGameNode(prev *GameNode, p *Position, m MoveData) *GameNode {
	n := &GameNode{
		Prev:     prev,
		Next:     nil,
		Position: p,
		MoveData: m,
	}
	if prev != nil {
		prev.Next = n
	}
	return n
}

func (t *GameTree) Move(m Move) error {
	// Nextが同じ指し手だったら単に進める
	if t.Current.Next != nil && m == t.Current.Next.MoveData.Move {
		t.Next()
		return nil
	}

	before := t.Current.MoveData.To
	moveData := NewMoveData(m, t.Current.Position, before)

	p := t.Current.Position.Clone()
	if err := p.Move(m); err != nil {
		return err
	}
	t.Current = NewGameNode(t.Current, p, moveData)
	return nil
}

// Currentを次の局面に進める。
// 成功したらtrue。
func (t *GameTree) Next() bool {
	if t.Current.Next == nil {
		return false
	}
	t.Current = t.Current.Next
	return true
}

// Currentを前の局面に進める。
// 成功したらtrue。
func (t *GameTree) Prev() bool {
	if t.Current.Prev == nil {
		return false
	}
	t.Current = t.Current.Prev
	return true
}
