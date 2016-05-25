package textures

import (
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"os"

	_ "image/png"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite"
)

var (
	txtFont     *truetype.Font         // 游戏上显示文字时，用的字体，简单期间只用一个字体
	hasLoadFont bool           = false // 是否已经加载了字体

)

func ReleaseFont() {
	txtFont = nil
	hasLoadFont = false
}

// 初始化时候加载字体
func LoadGameFont(fileName string) error {
	if !hasLoadFont { // 确保只加载一次，由于是顺序执行，这里就没有用锁。
		var f io.Reader
		if len(fileName) <= 0 {
			file, err := asset.Open("f1.ttf")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			f = file
		} else {
			log.Println(fileName)
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			f = file

		}

		fontBytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			return err
		}

		txtFont, err = truetype.Parse(fontBytes)
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

// 指定图上绘制文字
// rgba 为需要在这个图上绘制
// txtSize 为字体尺寸
// fontColor 为字体颜色
// cpt 为需要绘制文字的位置， 为需要绘制的下面中心点的位置
// txt 为需要绘制的文字
func DrawString(rgba *image.RGBA,
	txtSize float64, fontColor color.RGBA,
	cpt image.Point, txt string) {

	d := &font.Drawer{
		Dst: rgba,                        // 背景图
		Src: image.NewUniform(fontColor), // 字体颜色
		Face: truetype.NewFace(txtFont, &truetype.Options{
			Size:    txtSize,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}
	// 绘制位置
	d.Dot = fixed.Point26_6{
		X: fixed.I(cpt.X) - d.MeasureString(txt)/2,
		Y: fixed.I(cpt.Y),
	}

	d.DrawString(txt)
}

// 根据字体大小和文字长度，自动绘制一个对应大小的图。
// txtSize 字体大小
// fontColor 要绘制的字体颜色
// txt 要绘制的文字
func DrawStringRGBA(txtSize float64, fontColor color.RGBA, txt string) *image.RGBA {

	d := &font.Drawer{
		Src: image.NewUniform(fontColor), // 字体颜色
		Face: truetype.NewFace(txtFont, &truetype.Options{
			Size:    txtSize,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}
	re := d.MeasureString(txt)
	rect := image.Rect(0, 0, re.Ceil(), int(txtSize))
	log.Println(txt, "大小", rect)
	rgba := image.NewRGBA(rect)
	d.Dst = rgba

	d.Dot = fixed.Point26_6{
		X: fixed.I(0),
		Y: fixed.I(rect.Max.Y),
	}
	d.DrawString(txt)
	return rgba
}

// 加载字体纹理
func LoadFontTextTextures(eng sprite.Engine, txt string) sprite.SubTex {

	txtColor := color.RGBA{227, 16, 205, 255} // RGBA, 不透明 A 为 255
	rgba := DrawStringRGBA(40.0, txtColor, txt)

	t, err := eng.LoadTexture(rgba)
	if err != nil {
		log.Fatal(err)
	}

	lastSubTex := sprite.SubTex{t, rgba.Bounds()}
	return lastSubTex

}
