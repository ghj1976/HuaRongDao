package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

const (
	GameAreaWidth                   float32 = 4.0                                            // 游戏区域宽度，不含边框，小兵棋子的 4 倍
	GameAreaHeight                  float32 = 5.0                                            // 游戏区域高度， 不含边框和曹营
	GameAreaAndBorderWidth          float32 = GameAreaWidth + 2.0*1.0/8.0                    // 游戏区域含边框宽度应该是小兵棋子的 4.25倍
	GameAreaAndBorderAndCampsHeight float32 = GameAreaHeight + 1.0/2.0 + 1.0/8.0             // 游戏区域高度，含一个边框＋曹营的高度 应该是小兵棋子的 5.625 倍
	ScreenAreaHeight                float32 = GameAreaHeight + 1.0/2 + 1.0/8 + 1.0/2 + 3.0/8 // 屏幕区域高度，应该是小兵棋子的 6.5倍  游戏区域 ＋ 曹营 ＋ 上边框 ＋ 按钮 ＋ 问题提示区域

	Speed = 2 // 棋子移动的速度
)

var (
	ChessManWidth                                                  float32   // 小兵棋子的宽度或者高度 ，单位 pt
	GameAreaAndBorderAndCampsAreaX, GameAreaAndBorderAndCampsAreaY float32   // 游戏纹理1绘制区域（含边框、曹营绘制内容，纹理1对应的绘图区域）的左上角坐标，单位 pt
	GameChessManAreaX, GameChessManAreaY                           float32   // 游戏棋子会出现最左上角的位置，单位 pt
	touchBeginPoint                                                GamePoint // touch 事件时，判断位移大小的初始位置。

	winNode *sprite.Node // 游戏过关的显示节点，这个在需要的时才显示，所以才会单独处理。
)

type Game struct {
	Level *LevelInfo // 当前的关卡信息类

	lastCalc clock.Time // when we last calculated a frame

	btnReturn *GameBtn // 返回按钮
	btnGuide  *GameBtn // 攻略按钮
	btnReload *GameBtn // 重玩按钮

	CurrTouchChessMan rune // 当前正在移动的棋子
}

func (g *Game) InitScene(eng sprite.Engine, sz size.Event) *sprite.Node {

	// log.Println(GameAreaAndBorderWidth, GameAreaAndBorderAndCampsHeight, ScreenAreaHeight)
	// 计算棋子兵应该的高度或长度。
	ch := float32(sz.HeightPt) / ScreenAreaHeight
	cw := float32(sz.WidthPt) / GameAreaAndBorderWidth
	if cw < ch {
		ChessManWidth = cw
		GameAreaAndBorderAndCampsAreaX = 0.0
		GameAreaAndBorderAndCampsAreaY = float32(sz.HeightPt) - ChessManWidth*GameAreaAndBorderAndCampsHeight
	} else {
		ChessManWidth = ch
		GameAreaAndBorderAndCampsAreaX = (float32(sz.WidthPt) - ChessManWidth*GameAreaAndBorderWidth) / 2
		GameAreaAndBorderAndCampsAreaY = float32(sz.HeightPt) - ChessManWidth*GameAreaAndBorderAndCampsHeight
	}
	// 棋子可以出现的最左上角位置。
	GameChessManAreaX = GameAreaAndBorderAndCampsAreaX + ChessManWidth*1.0/8
	GameChessManAreaY = GameAreaAndBorderAndCampsAreaY + ChessManWidth*1.0/8

	log.Println("棋子 兵 宽度:", ChessManWidth)
	scene := &sprite.Node{}

	err := loadFont("./assets/f1.ttf")
	if err != nil {
		log.Panicln(err)
		return scene
	}

	texs := loadTextures(eng)

	txtColor := color.RGBA{227, 16, 205, 255} // RGBA, 不透明 A 为 255
	texLevelName := loadFontTextTextures(eng, "横刀立马", 40.0, txtColor, image.Rect(0, 0, 240, 60))
	texLevelStep := loadFontTextTextures(eng, "0/0", 40.0, txtColor, image.Rect(0, 0, 240, 60))

	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	newNode := func(fn arrangerFunc) {
		n := &sprite.Node{Arranger: arrangerFunc(fn)}
		eng.Register(n)
		scene.AppendChild(n)
	}

	newNodeNoShow := func(fn arrangerFunc) *sprite.Node {
		n := &sprite.Node{Arranger: arrangerFunc(fn)}
		eng.Register(n)
		return n
	}

	// 绘制游戏区域背景
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texs[texGameArea])
		eng.SetTransform(n, f32.Affine{
			{ChessManWidth * GameAreaAndBorderWidth, 0, GameAreaAndBorderAndCampsAreaX},
			{0, ChessManWidth * GameAreaAndBorderAndCampsHeight, GameAreaAndBorderAndCampsAreaY},
		})

	})

	// 绘制关卡名称
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texLevelName)
		eng.SetTransform(n, f32.Affine{
			{ChessManWidth * 1.5, 0, GameAreaAndBorderAndCampsAreaX + ChessManWidth/2},
			{0, ChessManWidth * 3 / 8, 0},
		})

	})
	// 绘制关卡最佳步速、当前步速
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texLevelStep)
		eng.SetTransform(n, f32.Affine{
			{ChessManWidth * 1.5, 0, GameAreaAndBorderAndCampsAreaX + 3*ChessManWidth},
			{0, ChessManWidth * 3 / 8, 0},
		})

	})

	// 返回按钮
	game.btnReturn = &GameBtn{status: BtnNormal}
	// 位置信息
	game.btnReturn.SetGameRectangle(
		GamePoint{
			X: (GameAreaAndBorderAndCampsAreaX + ChessManWidth*3/8),
			Y: (ChessManWidth * 3 / 8),
		},
		ChessManWidth,
		(ChessManWidth / 2))
	//log.Println(game.btnReturn)
	// 绘图
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if game.btnReturn.status == BtnNormal {
			eng.SetSubTex(n, texs[texBtnReturn1])
		} else {
			//log.Println("ReDraw BtnReturn Press ")
			eng.SetSubTex(n, texs[texBtnReturn3])
		}
		eng.SetTransform(n, f32.Affine{
			{game.btnReturn.Width, 0, game.btnReturn.LeftTop.X},
			{0, game.btnReturn.Height, game.btnReturn.LeftTop.Y},
		})

	})

	// 攻略 按钮
	game.btnGuide = &GameBtn{status: BtnNormal}
	game.btnGuide.SetGameRectangle(
		GamePoint{
			X: GameAreaAndBorderAndCampsAreaX + ChessManWidth*13/8,
			Y: ChessManWidth * 3 / 8,
		},
		ChessManWidth,
		ChessManWidth/2)
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if game.btnGuide.status == BtnNormal {
			eng.SetSubTex(n, texs[texBtnGuide1])
		} else {
			eng.SetSubTex(n, texs[texBtnGuide3])
		}
		eng.SetTransform(n, f32.Affine{
			{game.btnGuide.Width, 0, game.btnGuide.LeftTop.X},
			{0, game.btnGuide.Height, game.btnGuide.LeftTop.Y},
		})

	})
	// 重玩 按钮
	game.btnReload = &GameBtn{status: BtnNormal}
	game.btnReload.SetGameRectangle(
		GamePoint{
			X: GameAreaAndBorderAndCampsAreaX + ChessManWidth*23/8,
			Y: ChessManWidth * 3 / 8,
		},
		ChessManWidth,
		ChessManWidth/2)
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if game.btnReload.status == BtnNormal {
			eng.SetSubTex(n, texs[texBtnReload1])
		} else {
			eng.SetSubTex(n, texs[texBtnReload3])
		}
		eng.SetTransform(n, f32.Affine{
			{game.btnReload.Width, 0, game.btnReload.LeftTop.X},
			{0, game.btnReload.Height, game.btnReload.LeftTop.Y},
		})

	})

	// 计算每个棋子的准确绘图位置
	game.Level.ComputeChessManRect()

	// 绘制所有棋子
	for name, _ := range game.Level.ChessMans {
		// 比较诡异， 直接使用遍历出来的内容， 在 for 循环时，指针混乱,怀疑它不是一个线程安全的，
		// 所以这里全部再赋值给一个本地变量，再根据 本地变量 cName 直接去取，避免这个问题。
		// 这里 for 循环的是指针， 但是内部又会依靠这个指针， 当 for 循环指针发生变换时，内部就会指向混乱。
		// 由于内部还要再用，所以这里需要复制一份对象，避免影响。
		cName := name
		cm := game.Level.ChessMans[cName]
		newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
			p := chessManFrame(cName, g.Level.Success, cm.RelWidth, cm.status, t, 16)
			//			log.Println(string(cName), p, cm.rect)

			// 避免某些纹理配置错误，无法加载的问题
			eng.SetSubTex(n, texs[p])

			eng.SetTransform(n, f32.Affine{
				{cm.rect.Width, 0, cm.rect.LeftTop.X},
				{0, cm.rect.Height, cm.rect.LeftTop.Y},
			})
		})
		//	log.Println(string(name))
	}

	winNode = newNodeNoShow(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texs[texWin])
		eng.SetTransform(n, f32.Affine{
			{ChessManWidth * 4, 0, GameChessManAreaX},
			{0, ChessManWidth * 2, GameChessManAreaY + ChessManWidth*1.5},
		})

	})

	return scene
}

func (g *Game) reset() {
	g.CurrTouchChessMan = BlackChessManPos
}

// 游戏结束，释放资源，退出游戏
func (g *Game) stop() {
	txtFont = nil
}

// 当 touch 事件发生时， 判断是按在那个游戏精灵元素上，以及对应的处理策略分支。
func (g *Game) Press(touchEvent touch.Event) {
	// 单位修改成 pt， 而不是 px
	gp := GamePoint{X: touchEvent.X / sz.PixelsPerPt, Y: touchEvent.Y / sz.PixelsPerPt}

	// 按钮 按下逻辑处理。
	if touchEvent.Type == touch.TypeBegin {
		if gp.In(g.btnReturn.GameRectangle) {
			// 返回按钮被点击
			g.btnReturn.status = BtnPress
			//log.Println("btnReturn 被按下")
			return
		} else if gp.In(g.btnGuide.GameRectangle) {
			g.btnGuide.status = BtnPress
			return
		} else if gp.In(g.btnReload.GameRectangle) {
			g.btnReload.status = BtnPress
			return
		}
	} else if touchEvent.Type == touch.TypeEnd {
		if g.btnReturn.status == BtnPress {
			// 返回按钮被释放
			g.btnReturn.status = BtnNormal
			log.Println("btnReturn 释放按下状态")
			// 返回按钮的操作逻辑
			return
		} else if g.btnGuide.status == BtnPress {
			g.btnGuide.status = BtnNormal
			// 攻略按钮的操作逻辑
			return
		} else if g.btnReload.status == BtnPress {
			g.btnReload.status = BtnNormal
			// 重玩按钮的操作逻辑
			return
		}

	}
	//	if touchEvent.Type == touch.TypeEnd {
	//		g.CurrTouchChessMan = BlackChessManPos
	//	}

	// 关卡成功结束后，不需要再处理棋子的移动事件
	if game.Level.Success {
		return
	}

	if touchEvent.Type == touch.TypeBegin && game.CurrTouchChessMan == BlackChessManPos {
		// 寻找是哪个棋子被按下了。
		for name, _ := range game.Level.ChessMans {
			cName := name
			cm := game.Level.ChessMans[cName]
			// 需要记录开始移动点的位置
			if cm.status == ChessManMovable {
				if gp.In(cm.rect) {
					game.CurrTouchChessMan = cName
					touchBeginPoint = gp

					log.Println("当前移动棋子：", string(game.CurrTouchChessMan), "当前棋子状态：", game.Level.ChessMans[game.CurrTouchChessMan])

					break
				}
			}

		}
	}

	// 棋子按下拖动逻辑处理。当移动距离超过一定距离，才触发移动，避免灵敏度太高
	if game.CurrTouchChessMan != BlackChessManPos {
		cm := game.Level.ChessMans[game.CurrTouchChessMan]
		// 移动距离超过一定距离，才触发移动动画
		if touchEvent.Type == touch.TypeMove && cm.status == ChessManMovable {
			moveX := touchBeginPoint.X - gp.X
			moveY := touchBeginPoint.Y - gp.Y
			absMoveX := math.Abs(float64(moveX))
			absMoveY := math.Abs(float64(moveY))
			//			log.Println("位移距离", moveX, moveY)
			if absMoveX > 6 || absMoveY > 6 { // 移动距离超过一定距离，才触发移动
				if absMoveX > absMoveY && moveX > 0 {
					// 向左移动
					if game.Level.ChessManCanMoveLeft(game.CurrTouchChessMan) {
						cm.status = ChessManMoving
						cm.MoveXDirection = -1
						cm.MoveYDirection = 0
					}
				} else if absMoveX > absMoveY && moveX < 0 {
					// 向右移动
					if game.Level.ChessManCanMoveRight(game.CurrTouchChessMan) {
						cm.status = ChessManMoving
						cm.MoveXDirection = 1
						cm.MoveYDirection = 0
					}
				} else if absMoveX < absMoveY && moveY > 0 {
					// 向上移动
					if game.Level.ChessManCanMoveUp(game.CurrTouchChessMan) {
						cm.status = ChessManMoving
						cm.MoveXDirection = 0
						cm.MoveYDirection = -1
					}
				} else if absMoveX < absMoveY && moveY < 0 {
					// 向下移动
					if game.Level.ChessManCanMoveDown(game.CurrTouchChessMan) {
						cm.status = ChessManMoving
						cm.MoveXDirection = 0
						cm.MoveYDirection = 1
					}
				}
			}
		}
	}

}

// 每次绘图前，棋子逻辑相关的操作。
func (g *Game) Update(now clock.Time) {
	// 棋子是否移动到位置的判断
	// 到位后需要调整棋子的状态，以便其他地方处理逻辑的调整。
	if g.CurrTouchChessMan != BlackChessManPos {
		CurrCM, ok := g.Level.ChessMans[g.CurrTouchChessMan] // 找到当前被 touch 的棋子
		if ok {
			if CurrCM.status == ChessManMoving { // 移动状态才需要考虑移动
				// 移动后的位置
				CurrCM.rect.LeftTop.X = CurrCM.rect.LeftTop.X + Speed*float32(CurrCM.MoveXDirection)
				CurrCM.rect.LeftTop.Y = CurrCM.rect.LeftTop.Y + Speed*float32(CurrCM.MoveYDirection)
				BorderX1 := GameChessManAreaX + ChessManWidth*float32(CurrCM.RelLeftTopX)
				BorderX2 := GameChessManAreaX + ChessManWidth*float32(CurrCM.RelLeftTopX+CurrCM.MoveXDirection)
				BorderY1 := GameChessManAreaY + ChessManWidth*float32(CurrCM.RelLeftTopY)
				BorderY2 := GameChessManAreaY + ChessManWidth*float32(CurrCM.RelLeftTopY+CurrCM.MoveYDirection)
				//log.Println("移动动画判断:", CurrCM.rect.LeftTop, BorderX1, BorderX2, BorderY1, BorderY2)
				CurrCM.rect.RightBottom.X = CurrCM.rect.LeftTop.X + CurrCM.rect.Width
				CurrCM.rect.RightBottom.Y = CurrCM.rect.LeftTop.Y + CurrCM.rect.Height

				// 移动后超过边界，复原到边界值, 完成对应移动，相应参数变换
				if !PointInRange(CurrCM.rect.LeftTop.X, BorderX1, BorderX2) ||
					!PointInRange(CurrCM.rect.LeftTop.Y, BorderY1, BorderY2) {
					// 完成移动

					// 精确位置的变换完成
					CurrCM.rect.LeftTop.X = BorderX2
					CurrCM.rect.LeftTop.Y = BorderY2
					CurrCM.rect.RightBottom.X = CurrCM.rect.LeftTop.X + CurrCM.rect.Width
					CurrCM.rect.RightBottom.Y = CurrCM.rect.LeftTop.Y + CurrCM.rect.Height
					// 棋盘原先属于自己的区域全部清空
					g.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelLeftTopX] = BlackChessManPos
					g.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelRightBottomX] = BlackChessManPos
					g.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelLeftTopX] = BlackChessManPos
					g.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelRightBottomX] = BlackChessManPos
					// 计算新的相对位置
					CurrCM.RelLeftTopX = CurrCM.RelLeftTopX + CurrCM.MoveXDirection
					CurrCM.RelLeftTopY = CurrCM.RelLeftTopY + CurrCM.MoveYDirection
					CurrCM.RelRightBottomX = CurrCM.RelRightBottomX + CurrCM.MoveXDirection
					CurrCM.RelRightBottomY = CurrCM.RelRightBottomY + CurrCM.MoveYDirection
					// 新的位置赋值
					g.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelLeftTopX] = g.CurrTouchChessMan
					g.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelRightBottomX] = g.CurrTouchChessMan
					g.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelLeftTopX] = g.CurrTouchChessMan
					g.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelRightBottomX] = g.CurrTouchChessMan

					// 移动过程中的变量复原
					CurrCM.MoveXDirection = 0
					CurrCM.MoveYDirection = 0
					CurrCM.status = ChessManMovable
					g.CurrTouchChessMan = BlackChessManPos // 复原当前选择棋子

					log.Println("移动后棋子状态：", CurrCM)

					// 重算棋盘的可移动状态。
					g.Level.ComputeChessManStatus()

					if g.Level.IsSuccess() {
						log.Println("成功")
						scene.AppendChild(winNode) // 显示成功节点
						g.Level.Success = true
						return
					}
				}

			}
		}
	}

	// 游戏是否成功的判断
}

// 每个精灵多一个需要判断是否自己被点击、被拖动，所以多传一个参数touch.Event
type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	a(e, n, t)
}
