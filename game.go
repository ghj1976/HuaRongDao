package main

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	_ "image/jpeg"
)

var (
	txtFont *truetype.Font // 游戏上显示文字时，用的字体，简单期间只用一个字体
)

type Game struct {
	lastCalc clock.Time // when we last calculated a frame

}

func (g *Game) InitScene(eng sprite.Engine, sz size.Event) *sprite.Node {
	scene := &sprite.Node{}

	err := loadFont("./assets/f1.ttf")
	if err != nil {
		log.Panicln(err)
		return scene
	}

	// back := loadTextures(eng)
	// b2 := loadTexturesB2(eng)
	b3 := loadFontTextTextures(eng, "哈哈", 46.0, color.Black, image.Rect(0, 0, 300, 150))

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

	// newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
	// 	eng.SetSubTex(n, back)
	// 	eng.SetTransform(n, f32.Affine{
	// 		{float32(sz.WidthPt), 0, 0},
	// 		{0, float32(sz.HeightPx), 0},
	// 	})
	// })

	// newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
	// 	eng.SetSubTex(n, b2)
	// 	eng.SetTransform(n, f32.Affine{
	// 		{float32(sz.WidthPt), 0, 0},
	// 		{0, 50, 0},
	// 	})
	// })

	// 测试字体
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, b3)
		eng.SetTransform(n, f32.Affine{
			{float32(sz.WidthPt), 0, 80},
			{0, 250, 100},
		})
	})
	return scene

}

// 加载纹理图
func loadTextures(eng sprite.Engine) sprite.SubTex {
	a, err := asset.Open("odd.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	m, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	return sprite.SubTex{t, image.Rect(0, 0, 128, 128)}

}

// 加载纹理图
func loadTexturesB2(eng sprite.Engine) sprite.SubTex {
	a, err := asset.Open("b2.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	m, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	return sprite.SubTex{t, image.Rect(0, 0, 267, 189)}

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
func loadFontTextTextures(eng sprite.Engine, txt string, txtSize float64, txtColor color.Gray16, rect image.Rectangle) sprite.SubTex {

	bg := image.Transparent

	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(txtFont)
	c.SetFontSize(txtSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(txtColor))
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

func (g *Game) Press(down bool) {
}

func (g *Game) Update(now clock.Time) {
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
