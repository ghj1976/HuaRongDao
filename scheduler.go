// 调度器, 这里是一个 MVC 的模型，
// 调度器相当于 C，
// 每个需要显示的页面是 View
// 业务模型，游戏逻辑 是 model
package main

import (
	"log"
	"sync"

	"github.com/ghj1976/HuaRongDao/level"
	"github.com/ghj1976/HuaRongDao/model"
	"github.com/ghj1976/HuaRongDao/view"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

type CurrView byte // 当前正在显示的View

const (
	currNoView CurrView = iota
	currSplashView
	currLoadingView
	currListView
	currGameView
)

var (
	gameScene       *sprite.Node // 游戏的绘图根节点， 不同 view 的都是绘制在这个下面的。
	splashViewNode  *sprite.Node
	loadingViewNode *sprite.Node
	currView        CurrView // 当前是哪个视图

	gv *view.GameView // 当前的游戏视图

	rwMutex *sync.RWMutex // 读写锁
)

// 初始化， 在 onStart 中完成，
// 手机上再次打开时，如果没有被回收，也会再次进入这里。
func Init(eng sprite.Engine) {
	currView = currNoView
	rwMutex = new(sync.RWMutex)

	gameScene = &sprite.Node{}
	eng.Register(gameScene)
	eng.SetTransform(gameScene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	// 开协程 加载 loading 页
	go load(eng)
}

func load(eng sprite.Engine) {

	// 加载 启动页
	splashViewNode = view.LoadSplashView(eng)

	rwMutex.Lock()
	gameScene.AppendChild(splashViewNode)
	rwMutex.Unlock()

	currView = currSplashView
	log.Println("启动页加载完成。")

	loadingViewNode = view.LoadLoadingView(eng)

	rwMutex.Lock()
	gameScene.AppendChild(loadingViewNode)
	rwMutex.Unlock()

	currView = currLoadingView
	log.Println("Loading 页加载完成。")

	// 可以开协程加载 游戏列表页面了，
	// 这里简单期间， 加载具体一个游戏。

	lv := level.NewLevelInfo(1, "横刀立马", 81, "经典布局",
		`	赵曹曹马
			赵曹曹马
			黄关关张
			黄甲乙张
			丙一一丁
			`, level.LevelNotPass)
	gm := model.NewGameModel(lv)

	gv = view.NewGameView(gm, eng)

	if gameScene == nil {
		log.Println("gameScene nil")
		return
	}

	rwMutex.Lock()
	if splashViewNode != nil {
		gameScene.RemoveChild(splashViewNode)
	}
	if loadingViewNode != nil {
		gameScene.RemoveChild(loadingViewNode)
	}
	rwMutex.Unlock()

	rwMutex.Lock()
	gameScene.AppendChild(gv.GameViewNode)
	rwMutex.Unlock()

	currView = currGameView
	log.Println("游戏 页加载完成。")
}

// 更新绘图信息
func Update(now clock.Time) {
	// 把 update 透传给需要的当前视图
	if currView == currGameView {
		gv.Update(now)
	} else {

	}
}

// 更新拖动事件
func Press(touchEvent touch.Event) {
	// 把  透传
	if currView == currGameView {
		gv.Press(touchEvent)
	} else {

	}
}
