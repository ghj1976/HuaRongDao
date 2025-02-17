package view

import (
	"fmt"
	"log"
	"math"

	"github.com/ghj1976/HuaRongDao/button"
	"github.com/ghj1976/HuaRongDao/common"
	"github.com/ghj1976/HuaRongDao/level"
	"github.com/ghj1976/HuaRongDao/model"
	"github.com/ghj1976/HuaRongDao/textures"

	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

// GameView 具体某一关卡的游戏视图类，这里只处理显示相关逻辑。
type GameView struct {
	// 每次新进入一个游戏时，这部分的变量不能被修改。
	RootViewNode *sprite.Node  // 游戏视图的根节点
	eng          sprite.Engine // 绘制页面引擎

	// 新进入一个游戏时，可以被修改的内容
	model   *model.GameModel // 游戏的模型
	winNode *sprite.Node     // 游戏过关的显示节点，这个在需要的时才显示，所以才会单独处理。

	// 不用考虑初始化，内部逻辑就会初始化的成员变量
	touchBeginPoint     common.GamePoint // touch 事件时，判断位移大小的初始位置。
	gameLevelNameTxt    string
	gameLevelNameSubTex sprite.SubTex // 游戏名称纹理
	gameStepNumTxt      string
	gameStepNumSubTex   sprite.SubTex // 当前步数纹理
}

// NewGameView 完全重新创建一个游戏对象。
func NewGameView(m *model.GameModel, eng sprite.Engine) *GameView {
	gv := GameView{}
	gv.eng = eng
	gv.RootViewNode = &sprite.Node{} // GaveView 的绘图根节点
	eng.Register(gv.RootViewNode)
	eng.SetTransform(gv.RootViewNode, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	gv.Reset(m)

	return &gv
}

// Reset 第二次再进入时的入口。
func (g *GameView) Reset(m *model.GameModel) {
	g.model = m
	// 计算每个元素最终的显示位置。
	sz := model.GetScreenSizeBlock()
	g.model.InitGameElementLength(sz)
	g.loadGameView(g.eng)
}

// 如果没加载好，则加载好 再返回显示节点。
// 如果已经加载好了， 直接返回显示节点
func (g *GameView) loadGameView(eng sprite.Engine) {

	newNode := func(fn common.ArrangerFunc) {
		n := &sprite.Node{Arranger: common.ArrangerFunc(fn)}
		eng.Register(n)
		g.RootViewNode.AppendChild(n)
	}

	err := textures.LoadGameFont("")
	if err != nil {
		log.Panicln(err)
		return
	}

	texs := textures.LoadTexturesGame(eng)

	// 绘制游戏区域背景
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texs[textures.TexGameArea])

		if g.model.TexGameAreaRectangle == nil {
			log.Println("g.model.TexGameAreaRectangle nil")
		}
		eng.SetTransform(n, g.model.TexGameAreaRectangle.ToF32Affine())
	})

	// 绘制关卡名称
	if len(g.gameLevelNameTxt) <= 0 {
		g.gameLevelNameTxt = g.model.Level.Name
		g.gameLevelNameSubTex = textures.LoadFontTextTextures(eng, g.gameLevelNameTxt)
	}
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, g.gameLevelNameSubTex)
		eng.SetTransform(n, g.model.TexLevelNameRectangle.ToF32Affine())
	})

	// 绘制关卡最佳步速、当前步速
	// 这里之前存在内存泄漏。

	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		levelStep := fmt.Sprintf("%d/%d", g.model.Level.StepNum, g.model.Level.MinStepNum)
		if levelStep != g.gameStepNumTxt {
			g.gameStepNumSubTex = textures.LoadFontTextTextures(eng, levelStep)
			g.gameStepNumTxt = levelStep
		}
		eng.SetSubTex(n, g.gameStepNumSubTex)
		eng.SetTransform(n, g.model.TexLevelStepRectangle.ToF32Affine())
	})

	// 返回按钮
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if g.model.BtnReturn.Visible {
			eng.SetSubTex(n, texs[textures.GameButtonFrame("return", g.model.BtnReturn.Status)])
			eng.SetTransform(n, g.model.BtnReturn.ToF32Affine())
		}
	})

	// 攻略按钮
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if g.model.BtnGuide.Visible {
			eng.SetSubTex(n, texs[textures.GameButtonFrame("guide", g.model.BtnGuide.Status)])
			eng.SetTransform(n, g.model.BtnGuide.ToF32Affine())
		}
	})

	// 重玩按钮
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if g.model.BtnReload.Visible {
			eng.SetSubTex(n, texs[textures.GameButtonFrame("reload", g.model.BtnReload.Status)])
			eng.SetTransform(n, g.model.BtnReload.ToF32Affine())
		}
	})

	// 绘制所有棋子
	for name := range g.model.Level.ChessMans {
		// 比较诡异， 直接使用遍历出来的内容， 在 for 循环时，指针混乱,怀疑它不是一个线程安全的，
		// 所以这里全部再赋值给一个本地变量，再根据 本地变量 cName 直接去取，避免这个问题。
		// 这里 for 循环的是指针， 但是内部又会依靠这个指针， 当 for 循环指针发生变换时，内部就会指向混乱。
		// 由于内部还要再用，所以这里需要复制一份对象，避免影响。
		cName := name
		newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
			cm := g.model.Level.ChessMans[cName] // 注意，这个必须在这里， 否则会 reset 时 指针指向之前的。
			p := textures.GameChessManFrame(cName, g.model.Level.Success, cm.RelWidth, cm.Status, t, 16)
			//			log.Println(string(cName), p, cm.rect)

			// 避免某些纹理配置错误，无法加载的问题
			eng.SetSubTex(n, texs[p])
			eng.SetTransform(n, cm.Rect.ToF32Affine())
		})
		//	log.Println(string(name))
	}

	// 通关提示图
	g.winNode = common.NewNodeNoShow(eng, func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texs[textures.TexWin])
		eng.SetTransform(n, g.model.TexWinRectangle.ToF32Affine())
	})

}

// Press 当 touch 事件发生时， 判断是按在那个游戏精灵元素上，以及对应的处理策略分支。
func (g *GameView) Press(touchEvent touch.Event) {
	sz, _ := model.GetScreenSize()
	// 单位修改成 pt， 而不是 px
	gp := common.GamePoint{X: touchEvent.X / sz.PixelsPerPt, Y: touchEvent.Y / sz.PixelsPerPt}
	// 按钮 按下逻辑处理。
	if touchEvent.Type == touch.TypeBegin {
		if gp.In(g.model.BtnReturn.GameRectangle) {
			// 返回按钮被点击
			g.model.BtnReturn.Status = button.BtnPress
			log.Println("btnReturn 被按下")
			return
		} else if gp.In(g.model.BtnGuide.GameRectangle) {
			g.model.BtnGuide.Status = button.BtnPress
			return
		} else if gp.In(g.model.BtnReload.GameRectangle) {
			g.model.BtnReload.Status = button.BtnPress
			return
		}
	} else if touchEvent.Type == touch.TypeEnd {
		if g.model.BtnReturn.Status == button.BtnPress {
			// 返回按钮被释放
			g.model.BtnReturn.Status = button.BtnNormal
			log.Println("btnReturn 释放按下状态")
			// 返回按钮的操作逻辑

			ReturnListView()
			return
		} else if g.model.BtnGuide.Status == button.BtnPress {
			g.model.BtnGuide.Status = button.BtnNormal
			// 攻略按钮的操作逻辑
			return
		} else if g.model.BtnReload.Status == button.BtnPress {
			g.model.BtnReload.Status = button.BtnNormal
			// 重玩按钮的操作逻辑
			g.reset()
			return
		}

	}

	// 关卡成功结束后，不需要再处理棋子的移动事件
	if g.model.Level.Success {
		return
	}

	if touchEvent.Type == touch.TypeBegin && model.ChessManIsBlank(g.model.CurrTouchChessMan) {
		// 寻找是哪个棋子被按下了。
		for name := range g.model.Level.ChessMans {
			cName := name
			cm := g.model.Level.ChessMans[cName]
			// 需要记录开始移动点的位置
			if cm.Status == level.ChessManMovable {
				if gp.In(cm.Rect) {
					g.model.CurrTouchChessMan = cName
					g.touchBeginPoint = gp

					log.Println("当前移动棋子：", string(g.model.CurrTouchChessMan), "当前棋子状态：",
						g.model.Level.ChessMans[g.model.CurrTouchChessMan])

					break
				}
			}

		}
	}

	// 棋子按下拖动逻辑处理。当移动距离超过一定距离，才触发移动，避免灵敏度太高
	if g.model.CurrTouchChessMan != level.BlackChessManPos {
		cm := g.model.Level.ChessMans[g.model.CurrTouchChessMan]
		// 移动距离超过一定距离，才触发移动动画
		if touchEvent.Type == touch.TypeMove && cm.Status == level.ChessManMovable {
			moveX := g.touchBeginPoint.X - gp.X
			moveY := g.touchBeginPoint.Y - gp.Y
			absMoveX := math.Abs(float64(moveX))
			absMoveY := math.Abs(float64(moveY))
			//			log.Println("位移距离", moveX, moveY)
			if absMoveX > 6 || absMoveY > 6 { // 移动距离超过一定距离，才触发移动
				if absMoveX > absMoveY && moveX > 0 {
					// 向左移动
					if g.model.Level.ChessManCanMoveLeft(g.model.CurrTouchChessMan) {
						cm.Status = level.ChessManMoving
						cm.MoveXDirection = -1
						cm.MoveYDirection = 0
					}
				} else if absMoveX > absMoveY && moveX < 0 {
					// 向右移动
					if g.model.Level.ChessManCanMoveRight(g.model.CurrTouchChessMan) {
						cm.Status = level.ChessManMoving
						cm.MoveXDirection = 1
						cm.MoveYDirection = 0
					}
				} else if absMoveX < absMoveY && moveY > 0 {
					// 向上移动
					if g.model.Level.ChessManCanMoveUp(g.model.CurrTouchChessMan) {
						cm.Status = level.ChessManMoving
						cm.MoveXDirection = 0
						cm.MoveYDirection = -1
					}
				} else if absMoveX < absMoveY && moveY < 0 {
					// 向下移动
					if g.model.Level.ChessManCanMoveDown(g.model.CurrTouchChessMan) {
						cm.Status = level.ChessManMoving
						cm.MoveXDirection = 0
						cm.MoveYDirection = 1
					}
				}
			}
		}
	}

}

// Update 每次绘图前，逻辑相关的操作。
func (g *GameView) Update(now clock.Time) {

	g.model.Update(now)
	if g.model.Level.IsSuccess() {
		if g.winNode.Parent == nil {
			g.RootViewNode.AppendChild(g.winNode) // 显示成功节点
		}
	}
}

// Stop 被暂停离开时，保存相关的操作
func (g *GameView) Stop() {
}

// Destroy 彻底销毁前的释放操作。
func (g *GameView) Destroy() {
}

// reset 返回到本关卡的第一步
func (g *GameView) reset() {
	if g.model.Level.IsSuccess() {
		if g.winNode.Parent != nil {
			g.RootViewNode.RemoveChild(g.winNode)
		}
	}
	g.model.Reset()
}

// 游戏结束，释放资源，退出游戏
func (g *GameView) stop() {
	textures.ReleaseFont()
}
