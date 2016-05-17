// 调度器, 这里是一个 MVC 的模型，
// 调度器相当于 C，
// 每个需要显示的页面是 View
// 业务模型，游戏逻辑 是 model
package main

import (
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

var (
	cmdChan = make(chan string)
)

// 初始化， 在 onStart 中完成，
// 手机上再次打开时，如果没有被回收，也会再次进入这里。
func Init() *sprite.Node {
	scene := &sprite.Node{}
	// 加载 启动页

	// 开协程 加载 loading 页
	return scene
}

func SelectChan() {

	go func() {
		var cmd string
		for {
			// 一直在等 chan 消息
			select {
			case cmd = <-cmdChan:
				{
					switch cmd {
					case "initFinish":

					case "":
					}
				}

			}
		}
	}()
}

// 加载指定view
func LoadView(view string) {

}

// 更新绘图信息
func /*(g *Game)*/ Update(now clock.Time) {
}
