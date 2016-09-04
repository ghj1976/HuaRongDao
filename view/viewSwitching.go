package view

import "github.com/ghj1976/HuaRongDao/level"

var (
	// SwitchingChan 视图切换通讯 chan
	SwitchingChan = make(chan ViewSwitching)
)

// ViewSwitchingType 视图切换类型枚举，只有2种
type ViewSwitchingType byte

const (
	// List2Game 从列表页到游戏页
	List2Game ViewSwitchingType = iota
	// Game2List 从游戏页到列表页
	Game2List
)

// ViewSwitching 视图切换参数实体类，用于 chan 消息传递
type ViewSwitching struct {
	SwitchingType ViewSwitchingType // 切换的类型
	ListCurrPage  int               // List2Game 时， list 所在的页面
	Level         *level.LevelInfo  // List2Game 时，需要传递的那个游戏
}

// LoadGameView 加载指定的关口
func LoadGameView(listCurrPage int, level *level.LevelInfo) {
	sw := ViewSwitching{
		SwitchingType: List2Game,
		ListCurrPage:  listCurrPage,
		Level:         level,
	}

	SwitchingChan <- sw
}

// ReturnListView 返回游戏列表页
func ReturnListView() {
	sw := ViewSwitching{
		SwitchingType: Game2List,
	}
	SwitchingChan <- sw
}
