// 游戏模型
package model

import (
	"log"
	"unicode/utf8"

	"github.com/ghj1976/HuaRongDao/button"
	"github.com/ghj1976/HuaRongDao/common"
	"github.com/ghj1976/HuaRongDao/level"
	"golang.org/x/mobile/event/size"

	"golang.org/x/mobile/exp/sprite/clock"
)

// todo const、var 应该都可以改成小写的

const (
	gameAreaWidth                   float32 = 4.0                                            // 游戏区域宽度，不含边框，小兵棋子的 4 倍
	gameAreaHeight                  float32 = 5.0                                            // 游戏区域高度， 不含边框和曹营
	GameAreaAndBorderWidth          float32 = gameAreaWidth + 2.0*1.0/8.0                    // 游戏区域含边框宽度应该是小兵棋子的 4.25倍
	GameAreaAndBorderAndCampsHeight float32 = gameAreaHeight + 1.0/2.0 + 1.0/8.0             // 游戏区域高度，含一个边框＋曹营的高度 应该是小兵棋子的 5.625 倍
	ScreenAreaHeight                float32 = gameAreaHeight + 1.0/2 + 1.0/8 + 1.0/2 + 3.0/8 // 屏幕区域高度，应该是小兵棋子的 6.5倍  游戏区域 ＋ 曹营 ＋ 上边框 ＋ 按钮 ＋ 问题提示区域

	Speed = 5 // 棋子移动的速度
)

// 具体某一局游戏的模型实体类
type GameModel struct {
	chessManWidth                                                  float32 // 小兵棋子的宽度或者高度 ，单位 pt
	gameAreaAndBorderAndCampsAreaX, gameAreaAndBorderAndCampsAreaY float32 // 游戏纹理1绘制区域（含边框、曹营绘制内容，纹理1对应的绘图区域）的左上角坐标，单位 pt
	gameChessManAreaX, gameChessManAreaY                           float32 // 游戏棋子会出现最左上角的位置，单位 pt

	Level *level.LevelInfo // 当前的关卡信息类

	lastCalc clock.Time // when we last calculated a frame

	preTouchChessMan  rune // 前一个被移动的棋子，棋子的连续移动，前一个是被连续移动的哪个棋子
	CurrTouchChessMan rune // 当前正在移动的棋子

	BtnReturn *button.GameBtn // 返回按钮
	BtnGuide  *button.GameBtn // 攻略按钮
	BtnReload *button.GameBtn // 重玩按钮

	TexGameAreaRectangle  *common.GameRectangle // TexGameArea 纹理位置
	TexLevelNameRectangle *common.GameRectangle // 关卡名称的 纹理位置
	TexLevelStepRectangle *common.GameRectangle // 关卡最佳步数、当前步数
	TexWinRectangle       *common.GameRectangle // 成功过关纹理的位置
}

func NewGameModel(lv *level.LevelInfo) *GameModel {
	gm := GameModel{}
	gm.Level = lv

	// 需要知道屏幕的大小，才能做这里的计算。
	gm.chessManWidth = 0.0
	gm.gameAreaAndBorderAndCampsAreaX = 0.0
	gm.gameAreaAndBorderAndCampsAreaY = 0.0
	gm.gameChessManAreaX = 0.0
	gm.gameChessManAreaY = 0.0

	// 返回按钮
	gm.BtnReturn = &button.GameBtn{Status: button.BtnNormal}
	// 攻略 按钮
	gm.BtnGuide = &button.GameBtn{Status: button.BtnNormal}
	// 重玩 按钮
	gm.BtnReload = &button.GameBtn{Status: button.BtnNormal}

	gm.TexGameAreaRectangle = &common.GameRectangle{}
	gm.TexLevelNameRectangle = &common.GameRectangle{}
	gm.TexLevelStepRectangle = &common.GameRectangle{}
	gm.TexWinRectangle = &common.GameRectangle{}

	gm.Reset()
	return &gm
}

// 当前关卡返回到第一步
func (gm *GameModel) Reset() {
	gm.Level.Reset()

	gm.CurrTouchChessMan = level.BlackChessManPos
	gm.preTouchChessMan = level.BlackChessManPos

	if gm.chessManWidth > 0 {
		gm.Level.ComputeChessManRect(gm.chessManWidth, gm.gameChessManAreaX, gm.gameChessManAreaY)
	}
}

// 计算每个元素的具体位置。
func (gm *GameModel) InitGameElementLength(sz size.Event) {
	if sz.HeightPt == 0 {
		return
	}
	log.Println("屏幕尺寸：", sz)
	// 计算棋子兵应该的高度或长度。
	ch := float32(sz.HeightPt) / ScreenAreaHeight
	cw := float32(sz.WidthPt) / GameAreaAndBorderWidth
	if cw < ch {
		gm.chessManWidth = cw
		gm.gameAreaAndBorderAndCampsAreaX = 0.0
		gm.gameAreaAndBorderAndCampsAreaY = float32(sz.HeightPt) - gm.chessManWidth*GameAreaAndBorderAndCampsHeight
	} else {
		gm.chessManWidth = ch
		gm.gameAreaAndBorderAndCampsAreaX = (float32(sz.WidthPt) - gm.chessManWidth*GameAreaAndBorderWidth) / 2
		gm.gameAreaAndBorderAndCampsAreaY = float32(sz.HeightPt) - gm.chessManWidth*GameAreaAndBorderAndCampsHeight
	}
	// 棋子可以出现的最左上角位置。
	gm.gameChessManAreaX = gm.gameAreaAndBorderAndCampsAreaX + gm.chessManWidth*1.0/8
	gm.gameChessManAreaY = gm.gameAreaAndBorderAndCampsAreaY + gm.chessManWidth*1.0/8

	log.Println("棋子 兵 宽度:", gm.chessManWidth)

	// 返回按钮位置信息
	gm.BtnReturn.SetGameRectangle(
		common.GamePoint{
			X: (gm.gameAreaAndBorderAndCampsAreaX + gm.chessManWidth*3/8),
			Y: (gm.chessManWidth * 3 / 8),
		},
		gm.chessManWidth,
		(gm.chessManWidth / 2))
	//log.Println(game.btnReturn)
	// 攻略按钮位置信息
	gm.BtnGuide.SetGameRectangle(
		common.GamePoint{
			X: gm.gameAreaAndBorderAndCampsAreaX + gm.chessManWidth*13/8,
			Y: gm.chessManWidth * 3 / 8,
		},
		gm.chessManWidth,
		gm.chessManWidth/2)
	// 重玩 按钮位置信息
	gm.BtnReload.SetGameRectangle(
		common.GamePoint{
			X: gm.gameAreaAndBorderAndCampsAreaX + gm.chessManWidth*23/8,
			Y: gm.chessManWidth * 3 / 8,
		},
		gm.chessManWidth,
		gm.chessManWidth/2)

	// 计算每个棋子的准确绘图位置
	gm.Level.ComputeChessManRect(gm.chessManWidth, gm.gameChessManAreaX, gm.gameChessManAreaY)

	// 游戏背景区域的位置
	gm.TexGameAreaRectangle.SetGameRectangle(
		common.GamePoint{
			X: gm.gameAreaAndBorderAndCampsAreaX,
			Y: gm.gameAreaAndBorderAndCampsAreaY,
		},
		gm.chessManWidth*GameAreaAndBorderWidth,
		gm.chessManWidth*GameAreaAndBorderAndCampsHeight,
	)

	// 关卡名称的位置
	gm.TexLevelNameRectangle.SetGameRectangle(
		common.GamePoint{
			X: gm.gameAreaAndBorderAndCampsAreaX + gm.chessManWidth/2,
			Y: 0,
		},
		gm.chessManWidth*5/16*float32(utf8.RuneCountInString(gm.Level.Name)),
		gm.chessManWidth*5/16, // 适当压缩， 彻底撑满是 3/8 ，压缩 1/16
	)

	// 关卡步数信息的位置
	gm.TexLevelStepRectangle.SetGameRectangle(
		common.GamePoint{
			X: gm.gameAreaAndBorderAndCampsAreaX + 3*gm.chessManWidth,
			Y: 0,
		},

		gm.chessManWidth*5/16*2.5,
		gm.chessManWidth*5/16, // 适当压缩， 彻底撑满是 3/8 ，压缩 1/16
	)

	gm.TexWinRectangle.SetGameRectangle(
		common.GamePoint{
			X: gm.gameChessManAreaX,
			Y: gm.gameChessManAreaY + gm.chessManWidth*1.5,
		},
		gm.chessManWidth*4,
		gm.chessManWidth*2,
	)
}

// 每次绘图前，逻辑相关的操作。
func (gm *GameModel) Update(now clock.Time) {
	// 棋子是否移动到位置的判断
	// 到位后需要调整棋子的状态，以便其他地方处理逻辑的调整。
	if gm.CurrTouchChessMan != level.BlackChessManPos {
		CurrCM, ok := gm.Level.ChessMans[gm.CurrTouchChessMan] // 找到当前被 touch 的棋子
		if ok {
			if CurrCM.Status == level.ChessManMoving { // 移动状态才需要考虑移动
				// 移动后的位置
				CurrCM.Rect.LeftTop.X = CurrCM.Rect.LeftTop.X + Speed*float32(CurrCM.MoveXDirection)
				CurrCM.Rect.LeftTop.Y = CurrCM.Rect.LeftTop.Y + Speed*float32(CurrCM.MoveYDirection)
				BorderX1 := gm.gameChessManAreaX + gm.chessManWidth*float32(CurrCM.RelLeftTopX)
				BorderX2 := gm.gameChessManAreaX + gm.chessManWidth*float32(CurrCM.RelLeftTopX+CurrCM.MoveXDirection)
				BorderY1 := gm.gameChessManAreaY + gm.chessManWidth*float32(CurrCM.RelLeftTopY)
				BorderY2 := gm.gameChessManAreaY + gm.chessManWidth*float32(CurrCM.RelLeftTopY+CurrCM.MoveYDirection)
				//log.Println("移动动画判断:", CurrCM.rect.LeftTop, BorderX1, BorderX2, BorderY1, BorderY2)
				CurrCM.Rect.RightBottom.X = CurrCM.Rect.LeftTop.X + CurrCM.Rect.Width
				CurrCM.Rect.RightBottom.Y = CurrCM.Rect.LeftTop.Y + CurrCM.Rect.Height

				// 移动后超过边界，复原到边界值, 完成对应移动，相应参数变换
				if !common.PointInRange(CurrCM.Rect.LeftTop.X, BorderX1, BorderX2) ||
					!common.PointInRange(CurrCM.Rect.LeftTop.Y, BorderY1, BorderY2) {
					// 完成移动

					// 游戏详细步数记录
					var fx rune // 方向
					if CurrCM.MoveXDirection < 0 {
						fx = '左'
					} else if CurrCM.MoveXDirection > 0 {
						fx = '右'
					} else if CurrCM.MoveYDirection < 0 {
						fx = '上'
					} else if CurrCM.MoveYDirection > 0 {
						fx = '下'
					}

					// 精确位置的变换完成
					CurrCM.Rect.LeftTop.X = BorderX2
					CurrCM.Rect.LeftTop.Y = BorderY2
					CurrCM.Rect.RightBottom.X = CurrCM.Rect.LeftTop.X + CurrCM.Rect.Width
					CurrCM.Rect.RightBottom.Y = CurrCM.Rect.LeftTop.Y + CurrCM.Rect.Height
					// 棋盘原先属于自己的区域全部清空
					gm.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelLeftTopX] = level.BlackChessManPos
					gm.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelRightBottomX] = level.BlackChessManPos
					gm.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelLeftTopX] = level.BlackChessManPos
					gm.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelRightBottomX] = level.BlackChessManPos
					// 计算新的相对位置
					CurrCM.RelLeftTopX = CurrCM.RelLeftTopX + CurrCM.MoveXDirection
					CurrCM.RelLeftTopY = CurrCM.RelLeftTopY + CurrCM.MoveYDirection
					CurrCM.RelRightBottomX = CurrCM.RelRightBottomX + CurrCM.MoveXDirection
					CurrCM.RelRightBottomY = CurrCM.RelRightBottomY + CurrCM.MoveYDirection
					// 新的位置赋值
					gm.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelLeftTopX] = gm.CurrTouchChessMan
					gm.Level.MapArray[CurrCM.RelLeftTopY][CurrCM.RelRightBottomX] = gm.CurrTouchChessMan
					gm.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelLeftTopX] = gm.CurrTouchChessMan
					gm.Level.MapArray[CurrCM.RelRightBottomY][CurrCM.RelRightBottomX] = gm.CurrTouchChessMan

					// 移动过程中的变量复原
					CurrCM.MoveXDirection = 0
					CurrCM.MoveYDirection = 0
					CurrCM.Status = level.ChessManMovable
					//					log.Println("移动后棋子状态：", CurrCM)

					gm.Level.StepRecord += string(gm.CurrTouchChessMan) + string(fx)

					if gm.preTouchChessMan != gm.CurrTouchChessMan {
						gm.Level.StepNum++
						gm.preTouchChessMan = gm.CurrTouchChessMan
					}
					gm.CurrTouchChessMan = level.BlackChessManPos // 复原当前选择棋子

					// 重算棋盘的可移动状态。
					gm.Level.ComputeChessManStatus()
					// 游戏是否成功的判断
					if gm.Level.IsSuccess() {
						log.Println("成功")
						log.Println(gm.Level.StepRecord)
						gm.Level.Success = true
						return
					}

				}

			}
		}
	}
}
