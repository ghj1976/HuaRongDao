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
	txtFont *truetype.Font // 游戏上显示文字时，用的字体，简单期间只用一个字体

	c  *freetype.Context // 绘制文字缓存的对象
	pt fixed.Point26_6
	bg *image.Uniform

	// 通过缓存避免重复的建大量对象
	lastFontText string        // 最后一个字体图片纹理的名字
	lastSubTex   sprite.SubTex // 最后一个字体图片纹理
)

func ReleaseFont() {
	txtFont = nil
}

// 初始化时候加载字体
func LoadGameFont() error {
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
	return nil
}

func InitFontText(txtSize float64, fontColor color.RGBA, rect image.Rectangle) {
	bg = image.Transparent

	c = freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(txtFont)
	c.SetFontSize(txtSize)
	uniform := image.NewUniform(fontColor)
	c.SetSrc(uniform)
	pt = freetype.Pt(10, 10+int(c.PointToFixed(txtSize)>>6))

	return
}

// 加载制定大小、颜色字体的文字
// 假设字体只有一行，而且使用的是默认字体
func LoadFontTextTextures(eng sprite.Engine, txt string, rect image.Rectangle) sprite.SubTex {
	// 这里之前有内存泄漏，目前修改为缓存最后一副画，只有需要的时候才画纹理。
	if lastFontText == txt {
		return lastSubTex
	} else {

		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)

		_, err := c.DrawString(txt, pt)
		if err != nil {
			log.Fatal(err)
		}

		t, err := eng.LoadTexture(rgba)
		if err != nil {
			log.Fatal(err)
		}

		lastFontText = txt
		lastSubTex = sprite.SubTex{t, rect}
		return lastSubTex
	}
}
