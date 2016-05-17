package common

import (
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

// 每个精灵多一个需要判断是否自己被点击、被拖动，所以多传一个参数touch.Event
type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	a(e, n, t)
}

var (
	NewNodeNoShow = func(eng sprite.Engine, fn arrangerFunc) *sprite.Node {
		n := &sprite.Node{Arranger: arrangerFunc(fn)}
		eng.Register(n)
		return n
	}
)
