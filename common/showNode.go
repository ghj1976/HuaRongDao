package common

import (
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

type ArrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a ArrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	a(e, n, t)
}

var (
	NewNodeNoShow = func(eng sprite.Engine, fn ArrangerFunc) *sprite.Node {
		n := &sprite.Node{Arranger: ArrangerFunc(fn)}
		eng.Register(n)
		return n
	}
)
