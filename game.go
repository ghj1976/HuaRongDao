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
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"

	_ "image/png"
)

const (
	GameAreaWidth    = 4.25 // 游戏区域宽度应该是小兵棋子的 4.25倍
	ScreenAreaHeight = 6.5  // 屏幕区域高度，应该是小兵棋子的 6.5倍
	GameAreaHeight   = 5.75 // 游戏区域高度，应该是小兵棋子的 6.5倍
)

var (
	txtFont              *truetype.Font // 游戏上显示文字时，用的字体，简单期间只用一个字体
	ChessManWidth        float32        // 小兵棋子的宽度或者高度 ，单位 pt
	GameAreaX, GameAreaY float32        // 游戏区域的左上角坐标，单位 pt
)

type Game struct {
	lastCalc clock.Time // when we last calculated a frame

}

func (g *Game) InitScene(eng sprite.Engine, sz size.Event) *sprite.Node {

	// 计算棋子兵应该的高度或长度。
	ch := float32(sz.HeightPt) / ScreenAreaHeight
	cw := float32(sz.WidthPt) / GameAreaWidth
	if cw < ch {
		ChessManWidth = cw
		GameAreaX = 0.0
		GameAreaY = float32(sz.HeightPt) - ChessManWidth*GameAreaHeight
	} else {
		ChessManWidth = ch
		GameAreaX = (float32(sz.WidthPt) - ChessManWidth*GameAreaWidth) / 2
		GameAreaY = float32(sz.HeightPt) - ChessManWidth*GameAreaHeight
	}

	scene := &sprite.Node{}

	err := loadFont("./assets/f1.ttf")
	if err != nil {
		log.Panicln(err)
		return scene
	}

	texs := loadTextures(eng)

	txtColor := color.RGBA{227, 16, 205, 1}
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
			{ChessManWidth * GameAreaWidth, 0, GameAreaX},
			{0, ChessManWidth * GameAreaHeight, GameAreaY},
		})
	})

	// 绘制关卡名称
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texLevelName)
		eng.SetTransform(n, f32.Affine{
			{ChessManWidth * 1.5, 0, GameAreaX + ChessManWidth/2},
			{0, ChessManWidth * 3 / 8, 0},
		})
	})
	// 绘制关卡最佳步速、当前步速
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texLevelStep)
		eng.SetTransform(n, f32.Affine{
			{ChessManWidth * 1.5, 0, GameAreaX + 3*ChessManWidth},
			{0, ChessManWidth * 3 / 8, 0},
		})
	})

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
		texGameArea: sprite.SubTex{t, image.Rect(0, 0, 1020, 1350)},
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
func (g *Game) Press(down bool) {
}

func (g *Game) Update(now clock.Time) {
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
