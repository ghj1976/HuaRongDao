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
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

// CurrView 枚举对象，会出现的几个视图
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

	currView CurrView       // 当前是哪个视图
	gv       *view.GameView // 当前的游戏视图
	listv    *view.ListView // 当前列表视图

	rwMutex *sync.RWMutex // 读写锁

	listvPage int // 列表页当前在第几页
)

// Init 初始化， 在 onStart 中完成，
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

// load 异步协程加载
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

	// 游戏列表页面了，
	lm := model.NewListModel(1)
	listv = view.NewListView(lm, eng)

	//	// 这里简单期间， 加载具体一个游戏。
	//	lv := level.NewLevelInfo(1, "横刀立马", 81, "经典布局",
	//		`	赵曹曹马
	//			赵曹曹马
	//			黄关关张
	//			黄甲乙张
	//			丙一一丁
	//			`, level.LevelNotPass)
	//	gm := model.NewGameModel(lv)

	//	gv = view.NewGameView(gm, eng)

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
	gameScene.AppendChild(listv.RootViewNode)
	rwMutex.Unlock()

	currView = currListView
	log.Println("游戏列表 页加载完成。")

	// 进入死循环，接收chan消息，看是否要切换 游戏页 还是游戏列表页。
	for {
		select {
		case switchView := <-view.SwitchingChan:

			if switchView.SwitchingType == view.List2Game {
				listViewLoadGameView(switchView.ListCurrPage, switchView.Level)
			} else {
				gameViewReturnListView()
			}

		}

	}
}

// Update 更新绘图信息
func Update(now clock.Time) {
	// 把 update 透传给需要的当前视图
	if currView == currGameView {
		gv.Update(now)
	} else if currView == currListView {
		listv.Update(now)
	} else {

	}
}

// Press 更新拖动事件
func Press(touchEvent touch.Event) {
	// 把 touch事件 透传
	if currView == currGameView {
		gv.Press(touchEvent)
	} else if currView == currListView {
		listv.Press(touchEvent)
	} else {

	}
}

// ScreenSizeChange 屏幕尺寸发生变化
func ScreenSizeChange(sz size.Event) {
	model.InitScreenSize(sz) // 记录尺寸变化
	if currView == currGameView {
		// gv
	} else if currView == currListView {
		listv.OnScreenSizeChange(sz, model.GetDisplayMultiple())
	} else {

	}
}

// listViewLoadGameView 从游戏列表页加载游戏视图
// listCurrPage 列表页当前在哪页，用于返回时返回该页
// level 要加载的关卡信息
func listViewLoadGameView(listCurrPage int, level *level.LevelInfo) {
	listvPage = listCurrPage

	// 显示进度中页面
	rwMutex.Lock()
	if listv != nil && listv.RootViewNode != nil {
		//listv.ClearRootViewChildNodes()
		gameScene.RemoveChild(listv.RootViewNode)
	}
	if splashViewNode != nil {
		gameScene.AppendChild(splashViewNode)
	}
	rwMutex.Unlock()

	gm := model.NewGameModel(level) // 关卡逻辑对象
	if gv == nil {
		gv = view.NewGameView(gm, eng)
	} else {
		gv.Reset(gm)
	}

	// 切换成列表页
	rwMutex.Lock()
	if splashViewNode != nil {
		gameScene.RemoveChild(splashViewNode)
	}

	gameScene.AppendChild(gv.RootViewNode)
	rwMutex.Unlock()

	// 切换成功，更新指向
	currView = currGameView
	log.Println("游戏页加载完成。")
}

// gameViewReturnListView 从游戏页面返回游戏列表
func gameViewReturnListView() {
	// 显示进度中页面
	rwMutex.Lock()
	if gv != nil && gv.RootViewNode != nil {
		gameScene.RemoveChild(gv.RootViewNode)
	}
	if splashViewNode != nil {
		gameScene.AppendChild(splashViewNode)
	}
	rwMutex.Unlock()

	// 由于后续会涉及到过关数据的变化，所以重新构造一个列表对象
	lm := model.NewListModel(listvPage)
	listv.Reset(lm)

	// 切换成列表页
	rwMutex.Lock()
	if splashViewNode != nil {
		gameScene.RemoveChild(splashViewNode)
	}

	gameScene.AppendChild(listv.RootViewNode)
	rwMutex.Unlock()

	// 切换成功，更新指向
	currView = currListView
	log.Println("游戏列表 页加载完成。")
}
