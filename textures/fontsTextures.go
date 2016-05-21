package textures

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	_ "image/png"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite"
)

var (
	txtFont     *truetype.Font         // 游戏上显示文字时，用的字体，简单期间只用一个字体
	hasLoadFont bool           = false // 是否已经加载了字体

	c  *freetype.Context // 绘制文字缓存的对象
	pt fixed.Point26_6
	bg *image.Uniform
)

func ReleaseFont() {
	txtFont = nil
	hasLoadFont = false
}

// 初始化时候加载字体
func LoadGameFont() error {
	if !hasLoadFont { // 确保只加载一次，由于是顺序执行，这里就没有用锁。
		a, err := asset.Open("f1.ttf")
		if err != nil {
			log.Fatal(err)
		}
		defer a.Close()

		fontBytes, err := ioutil.ReadAll(a)
		if err != nil {
			log.Println(err)
			return err
		}
		txtFont, err = freetype.ParseFont(fontBytes)
		if err != nil {
			log.Println(err)
			return err
		}
		hasLoadFont = true
		return nil
	} else {
		return nil
	}
}

// 初始化字体绘制相关参数
func InitFont() {
	bg = image.Transparent

	c = freetype.NewContext()
	c.SetFont(txtFont)
	c.SetDPI(72)
}

// 加载制定大小、颜色字体的文字
// 假设字体只有一行，而且使用的是默认字体
func LoadFontTextTextures(eng sprite.Engine, txt string, rect image.Rectangle) sprite.SubTex {
	// 这里之前有内存泄漏，目前修改为缓存最后一副画，只有需要的时候才画纹理。
	rgba := image.NewRGBA(rect)
	txtColor := color.RGBA{227, 16, 205, 255} // RGBA, 不透明 A 为 255
	DrawString(rgba, 40.0, txtColor, 10, 10, txt)

	t, err := eng.LoadTexture(rgba)
	if err != nil {
		log.Fatal(err)
	}

	lastSubTex := sprite.SubTex{t, rect}
	return lastSubTex

}

// 指定图上绘制文字
func DrawString(rgba *image.RGBA,
	txtSize float64, fontColor color.RGBA,
	leftTopX, leftTopY int, txt string) {

	c.SetFontSize(txtSize)                 // 字号
	uniform := image.NewUniform(fontColor) // 字体颜色
	c.SetSrc(uniform)
	pt = freetype.Pt(leftTopX, leftTopY+int(c.PointToFixed(txtSize)>>6)) // 位置

	// 要绘制的指定图
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)

	// 写字
	_, err := c.DrawString(txt, pt)
	if err != nil {
		log.Fatal(err)
	}
}
