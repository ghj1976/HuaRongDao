package main

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"

	_ "image/png"
)

const (
	GameAreaWidth                   float32 = 4.0                                            // 游戏区域宽度，不含边框，小兵棋子的 4 倍
	GameAreaHeight                  float32 = 5.0                                            // 游戏区域高度， 不含边框和曹营
	GameAreaAndBorderWidth          float32 = GameAreaWidth + 2.0*1.0/8.0                    // 游戏区域含边框宽度应该是小兵棋子的 4.25倍
	GameAreaAndBorderAndCampsHeight float32 = GameAreaHeight + 1.0/2.0 + 1.0/8.0             // 游戏区域高度，含一个边框＋曹营的高度 应该是小兵棋子的 5.625 倍
	ScreenAreaHeight                float32 = GameAreaHeight + 1.0/2 + 1.0/8 + 1.0/2 + 3.0/8 // 屏幕区域高度，应该是小兵棋子的 6.5倍  游戏区域 ＋ 曹营 ＋ 上边框 ＋ 按钮 ＋ 问题提示区域
)

var (
	txtFont                                                        *truetype.Font // 游戏上显示文字时，用的字体，简单期间只用一个字体
	ChessManWidth                                                  float32        // 小兵棋子的宽度或者高度 ，单位 pt
	GameAreaAndBorderAndCampsAreaX, GameAreaAndBorderAndCampsAreaY float32        // 游戏纹理1绘制区域（含边框、曹营绘制内容，纹理1对应的绘图区域）的左上角坐标，单位 pt
	GameChessManAreaX, GameChessManAreaY                           float32        // 游戏棋子会出现最左上角的位置，单位 pt
)

type BtnStatus byte // 按钮的状态枚举

// 游戏中的按钮类
type GameBtn struct {
	status        BtnStatus // 按钮的状态， 一共2种，按下、正常
	GameRectangle           // 按钮所在位置（长方形）
}

type Game struct {
	Level *LevelInfo // 当前的关卡信息类

	lastCalc clock.Time // when we last calculated a frame

	btnReturn *GameBtn // 返回按钮
	btnGuide  *GameBtn // 攻略按钮
	btnReload *GameBtn // 重玩按钮

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

	//	log.Println("aaa:", ChessManWidth)
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
			p := chessManFrame(cName, cm.RelWidth, cm.status, t, 16)
			log.Println(string(cName), p, cm.rect)

			// 避免某些纹理配置错误，无法加载的问题
			eng.SetSubTex(n, texs[p])
			eng.SetTransform(n, f32.Affine{
				{cm.rect.Width, 0, cm.rect.LeftTop.X},
				{0, cm.rect.Height, cm.rect.LeftTop.Y},
			})
		})
		log.Println(string(name))
	}

	return scene
}

const (
	texGameArea = iota // 游戏区域

	texBtnReturn1 // 返回按钮 1
	texBtnReturn2 // 返回按钮 2
	texBtnReturn3 // 返回按钮 3
	texBtnGuide1  // 攻略按钮 1
	texBtnGuide2  // 攻略按钮 2
	texBtnGuide3  // 攻略按钮 3
	texBtnReload1 // 重玩按钮 1
	texBtnReload2 // 重玩按钮 2
	texBtnReload3 // 重玩按钮 3

	texChessmanCaoCao1 // 棋子 曹操 1
	texChessmanCaoCao2 // 棋子 曹操 2
	texChessmanCaoCao3 // 棋子 曹操 3

	texChessmanHGuanYu1     // 棋子 横 关羽 1
	texChessmanHGuanYu2     // 棋子 横 关羽 2
	texChessmanHGuanYu3     // 棋子 横 关羽 3
	texChessmanHHuangZhong1 // 棋子 横 黄忠 1
	texChessmanHHuangZhong2 // 棋子 横 黄忠 2
	texChessmanHHuangZhong3 // 棋子 横 黄忠 3
	texChessmanHMaChao1     // 棋子 横 马超 1
	texChessmanHMaChao2     // 棋子 横 马超 2
	texChessmanHMaChao3     // 棋子 横 马超 3
	texChessmanHZhangFei1   // 棋子 横 张飞 1
	texChessmanHZhangFei2   // 棋子 横 张飞 2
	texChessmanHZhangFei3   // 棋子 横 张飞 3
	texChessmanHZhaoYun1    // 棋子 横 赵云 1
	texChessmanHZhaoYun2    // 棋子 横 赵云 2
	texChessmanHZhaoYun3    // 棋子 横 赵云 3

	texChessmanSGuanYu1     // 棋子 竖 关羽 1
	texChessmanSGuanYu2     // 棋子 竖 关羽 2
	texChessmanSGuanYu3     // 棋子 竖 关羽 3
	texChessmanSHuangZhong1 // 棋子 竖 黄忠 1
	texChessmanSHuangZhong2 // 棋子 竖 黄忠 2
	texChessmanSHuangZhong3 // 棋子 竖 黄忠 3
	texChessmanSMaChao1     // 棋子 竖 马超 1
	texChessmanSMaChao2     // 棋子 竖 马超 2
	texChessmanSMaChao3     // 棋子 竖 马超 3
	texChessmanSZhangFei1   // 棋子 竖 张飞 1
	texChessmanSZhangFei2   // 棋子 竖 张飞 2
	texChessmanSZhangFei3   // 棋子 竖 张飞 3
	texChessmanSZhaoYun1    // 棋子 竖 赵云 1
	texChessmanSZhaoYun2    // 棋子 竖 赵云 2
	texChessmanSZhaoYun3    // 棋子 竖 赵云 3

	texChessmanA1 // 棋子 兵 甲 1
	texChessmanA2 // 棋子 兵 甲 2
	texChessmanA3 // 棋子 兵 甲 3
	texChessmanB1 // 棋子 兵 乙 1
	texChessmanB2 // 棋子 兵 乙 2
	texChessmanB3 // 棋子 兵 乙 3
	texChessmanC1 // 棋子 兵 丙 1
	texChessmanC2 // 棋子 兵 丙 2
	texChessmanC3 // 棋子 兵 丙 3
	texChessmanD1 // 棋子 兵 丁 1
	texChessmanD2 // 棋子 兵 丁 2
	texChessmanD3 // 棋子 兵 丁 3

	// 以下为状态
	ChessManStable  // 棋子不可移动状态
	ChessManMovable // 棋子可移动状态
	ChessManMoving  // 棋子正在移动中

	BtnPress  // 按钮被按下状态
	BtnNormal // 按钮正常状态
)

// 加载纹理图,多张纹理
func loadTextures(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	m, _, err := image.Decode(a)
	if err != nil {
		log.Println("2")
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	return []sprite.SubTex{
		texGameArea:   sprite.SubTex{t, image.Rect(0, 0, 1020, 1350)},
		texBtnReturn1: sprite.SubTex{t, image.Rect(1100, 0, 1100+240, 120)},
		texBtnReturn2: sprite.SubTex{t, image.Rect(1100, 125, 1100+240, 125+120)},
		texBtnReturn3: sprite.SubTex{t, image.Rect(1100, 250, 1100+240, 250+120)},
		texBtnGuide1:  sprite.SubTex{t, image.Rect(1100, 375, 1100+240, 375+120)},
		texBtnGuide2:  sprite.SubTex{t, image.Rect(1100, 500, 1100+240, 500+120)},
		texBtnGuide3:  sprite.SubTex{t, image.Rect(1100, 625, 1100+240, 625+120)},
		texBtnReload1: sprite.SubTex{t, image.Rect(1100, 750, 1100+240, 750+120)},
		texBtnReload2: sprite.SubTex{t, image.Rect(1100, 875, 1100+240, 875+120)},
		texBtnReload3: sprite.SubTex{t, image.Rect(1100, 1000, 1100+240, 1000+120)},

		texChessmanCaoCao1: sprite.SubTex{t, image.Rect(0, 1400, 0+480, 1400+480)},
		texChessmanCaoCao2: sprite.SubTex{t, image.Rect(484, 1400, 484+480, 1400+480)},
		texChessmanCaoCao3: sprite.SubTex{t, image.Rect(968, 1400, 968+480, 1400+480)},

		texChessmanHGuanYu1:     sprite.SubTex{t, image.Rect(0, 1900, 0+480, 1900+240)},
		texChessmanHGuanYu2:     sprite.SubTex{t, image.Rect(484, 1900, 484+480, 1900+240)},
		texChessmanHGuanYu3:     sprite.SubTex{t, image.Rect(968, 1900, 960+480, 1900+240)},
		texChessmanHHuangZhong1: sprite.SubTex{t, image.Rect(0, 2150, 0+480, 2150+240)},
		texChessmanHHuangZhong2: sprite.SubTex{t, image.Rect(484, 2150, 484+480, 2150+240)},
		texChessmanHHuangZhong3: sprite.SubTex{t, image.Rect(968, 2150, 960+480, 2150+240)},
		texChessmanHMaChao1:     sprite.SubTex{t, image.Rect(0, 2400, 0+480, 2400+240)},
		texChessmanHMaChao2:     sprite.SubTex{t, image.Rect(484, 2400, 484+480, 2400+240)},
		texChessmanHMaChao3:     sprite.SubTex{t, image.Rect(968, 2400, 968+480, 2400+240)},
		texChessmanHZhangFei1:   sprite.SubTex{t, image.Rect(0, 2650, 0+480, 2650+240)},
		texChessmanHZhangFei2:   sprite.SubTex{t, image.Rect(484, 2650, 484+480, 2650+240)},
		texChessmanHZhangFei3:   sprite.SubTex{t, image.Rect(968, 2650, 968+480, 2650+240)},
		texChessmanHZhaoYun1:    sprite.SubTex{t, image.Rect(0, 2900, 0+480, 2900+240)},
		texChessmanHZhaoYun2:    sprite.SubTex{t, image.Rect(484, 2900, 484+480, 2900+240)},
		texChessmanHZhaoYun3:    sprite.SubTex{t, image.Rect(968, 2900, 968+480, 2900+240)},

		texChessmanSGuanYu1:     sprite.SubTex{t, image.Rect(0, 3150, 0+240, 3150+480)},
		texChessmanSGuanYu2:     sprite.SubTex{t, image.Rect(242, 3150, 242+240, 3150+480)},
		texChessmanSGuanYu3:     sprite.SubTex{t, image.Rect(484, 3150, 484+240, 3150+480)},
		texChessmanSHuangZhong1: sprite.SubTex{t, image.Rect(726, 3150, 726+240, 3150+480)},
		texChessmanSHuangZhong2: sprite.SubTex{t, image.Rect(968, 3150, 968+240, 3150+480)},
		texChessmanSHuangZhong3: sprite.SubTex{t, image.Rect(1209, 3150, 1209+240, 3150+480)},
		texChessmanSMaChao1:     sprite.SubTex{t, image.Rect(0, 3650, 0+240, 3650+480)},
		texChessmanSMaChao2:     sprite.SubTex{t, image.Rect(242, 3650, 242+240, 3650+480)},
		texChessmanSMaChao3:     sprite.SubTex{t, image.Rect(484, 3650, 484+240, 3650+480)},
		texChessmanSZhangFei1:   sprite.SubTex{t, image.Rect(726, 3650, 726+240, 3650+480)},
		texChessmanSZhangFei2:   sprite.SubTex{t, image.Rect(968, 3650, 968+240, 3650+480)},
		texChessmanSZhangFei3:   sprite.SubTex{t, image.Rect(1209, 3650, 1209+240, 3650+480)},
		texChessmanSZhaoYun1:    sprite.SubTex{t, image.Rect(0, 4150, 0+240, 4150+480)},
		texChessmanSZhaoYun2:    sprite.SubTex{t, image.Rect(242, 4150, 242+240, 4150+480)},
		texChessmanSZhaoYun3:    sprite.SubTex{t, image.Rect(484, 4150, 484+240, 4150+480)},

		texChessmanA1: sprite.SubTex{t, image.Rect(0, 4650, 0+240, 4650+240)},
		texChessmanA2: sprite.SubTex{t, image.Rect(242, 4650, 242+240, 4650+240)},
		texChessmanA3: sprite.SubTex{t, image.Rect(484, 4650, 484+240, 4650+240)},
		texChessmanB1: sprite.SubTex{t, image.Rect(726, 4650, 726+240, 4650+240)},
		texChessmanB2: sprite.SubTex{t, image.Rect(968, 4650, 968+240, 4650+240)},
		texChessmanB3: sprite.SubTex{t, image.Rect(1209, 4650, 1209+240, 4650+240)},
		texChessmanC1: sprite.SubTex{t, image.Rect(0, 5000, 0+240, 5000+240)},
		texChessmanC2: sprite.SubTex{t, image.Rect(242, 5000, 242+240, 5000+240)},
		texChessmanC3: sprite.SubTex{t, image.Rect(484, 5000, 484+240, 5000+240)},
		texChessmanD1: sprite.SubTex{t, image.Rect(726, 5000, 726+240, 5000+240)},
		texChessmanD2: sprite.SubTex{t, image.Rect(968, 5000, 968+240, 5000+240)},
		texChessmanD3: sprite.SubTex{t, image.Rect(1209, 5000, 1209+240, 5000+240)},
	}

}

// 初始化时候加载字体
func loadFont(fontFileName string) error {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	txtFont, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 加载制定大小、颜色字体的文字
// 假设字体只有一行，而且使用的是默认字体
func loadFontTextTextures(eng sprite.Engine, txt string, txtSize float64, fontColor color.RGBA, rect image.Rectangle) sprite.SubTex {

	bg := image.Transparent

	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(txtFont)
	c.SetFontSize(txtSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	uniform := image.NewUniform(fontColor)
	c.SetSrc(uniform)
	//c.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(txtSize)>>6))

	_, err := c.DrawString(txt, pt)

	m := rgba.SubImage(rect)

	// txtFont

	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	return sprite.SubTex{t, image.Rect(0, 0, rgba.Bounds().Size().X, rgba.Bounds().Size().Y)}

}

func (g *Game) reset() {

}

// 游戏结束，释放资源，退出游戏
func (g *Game) stop() {
	txtFont = nil
}

// 当 touch 事件发生时， 判断是按在那个游戏精灵元素上，以及对应的处理策略分支。
func (g *Game) Press(touchEvent touch.Event) {
	// 单位修改成 pt， 而不是 px
	gp := GamePoint{X: touchEvent.X / sz.PixelsPerPt, Y: touchEvent.Y / sz.PixelsPerPt}

	if touchEvent.Type == touch.TypeBegin {
		if gp.In(g.btnReturn.GameRectangle) {
			// 返回按钮被点击
			g.btnReturn.status = BtnPress
			//log.Println("btnReturn 被按下")
		} else if gp.In(g.btnGuide.GameRectangle) {
			g.btnGuide.status = BtnPress
		} else if gp.In(g.btnReload.GameRectangle) {
			g.btnReload.status = BtnPress
		}
	} else if touchEvent.Type == touch.TypeEnd {
		if g.btnReturn.status == BtnPress {
			// 返回按钮被释放
			g.btnReturn.status = BtnNormal
			log.Println("btnReturn 释放按下状态")
			// 返回按钮的操作逻辑
		} else if g.btnGuide.status == BtnPress {
			g.btnGuide.status = BtnNormal
			// 攻略按钮的操作逻辑
		} else if g.btnReload.status == BtnPress {
			g.btnReload.status = BtnNormal
			// 重玩按钮的操作逻辑
		}

	}

}

// 每个棋子应该用那个纹理来绘制
// name 棋子的名称， relWidth 决定棋子是横版 还是竖版
// status 棋子目前的状态
// t 当前的时间，用于轮播，  d 多长时间轮播一次
func chessManFrame(name rune, relWidth int, status ChessManStatus, t, d clock.Time) int {

	var frames []int // 那个棋子的判断
	switch name {
	case '曹':
		frames = []int{texChessmanCaoCao1, texChessmanCaoCao2, texChessmanCaoCao3}
	case '甲':
		frames = []int{texChessmanA1, texChessmanA2, texChessmanA3}
	case '乙':
		frames = []int{texChessmanB1, texChessmanB2, texChessmanB3}
	case '丙':
		frames = []int{texChessmanC1, texChessmanC2, texChessmanC3}
	case '丁':
		frames = []int{texChessmanD1, texChessmanD2, texChessmanD3}
	case '关':
		if relWidth == 1 {
			frames = []int{texChessmanSGuanYu1, texChessmanSGuanYu2, texChessmanSGuanYu3}
		} else {
			frames = []int{texChessmanHGuanYu1, texChessmanHGuanYu2, texChessmanHGuanYu3}
		}
	case '张':
		if relWidth == 1 {
			frames = []int{texChessmanSZhangFei1, texChessmanSZhangFei2, texChessmanSZhangFei3}
		} else {
			frames = []int{texChessmanHZhangFei1, texChessmanHZhangFei2, texChessmanHZhangFei3}
		}
	case '赵':
		if relWidth == 1 {
			frames = []int{texChessmanSZhaoYun1, texChessmanSZhaoYun2, texChessmanSZhaoYun3}
		} else {
			frames = []int{texChessmanHZhaoYun1, texChessmanHZhaoYun2, texChessmanHZhaoYun3}
		}
	case '马':
		if relWidth == 1 {
			frames = []int{texChessmanSMaChao1, texChessmanSMaChao2, texChessmanSMaChao3}
		} else {
			frames = []int{texChessmanHMaChao1, texChessmanHMaChao2, texChessmanHMaChao3}
		}
	case '黄':
		if relWidth == 1 {
			frames = []int{texChessmanSHuangZhong1, texChessmanSHuangZhong2, texChessmanSHuangZhong3}
		} else {
			frames = []int{texChessmanHHuangZhong1, texChessmanHHuangZhong2, texChessmanHHuangZhong3}
		}
	}

	if status == ChessManStable {
		// 不可移动
		return frames[0]
	} else if status == ChessManMovable {
		// 可移动
		a := int(d) * 2
		b := (int(t) % a) / int(d)
		if b == 0 {
			return frames[0]
		} else {
			return frames[1]
		}
	} else if status == ChessManMoving {
		// 正在移动
		a := int(d) * 2
		b := (int(t) % a) / int(d)
		if b == 0 {
			return frames[0]
		} else {
			return frames[2]
		}
	} else {
		return frames[0]
	}
}

func (g *Game) Update(now clock.Time) {
}

// 每个精灵多一个需要判断是否自己被点击、被拖动，所以多传一个参数touch.Event
type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	a(e, n, t)
}
