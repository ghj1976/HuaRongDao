package textures

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	_ "image/png"
	"log"

	"github.com/ghj1976/HuaRongDao/level"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

const (
	TexGameArea = iota // 游戏区域

	TexBtnReturn1 // 返回按钮 1
	TexBtnReturn2 // 返回按钮 2
	TexBtnReturn3 // 返回按钮 3
	TexBtnGuide1  // 攻略按钮 1
	TexBtnGuide2  // 攻略按钮 2
	TexBtnGuide3  // 攻略按钮 3
	TexBtnReload1 // 重玩按钮 1
	TexBtnReload2 // 重玩按钮 2
	TexBtnReload3 // 重玩按钮 3

	TexChessmanCaoCao1 // 棋子 曹操 1
	TexChessmanCaoCao2 // 棋子 曹操 2
	TexChessmanCaoCao3 // 棋子 曹操 3

	TexChessmanHGuanYu1     // 棋子 横 关羽 1
	TexChessmanHGuanYu2     // 棋子 横 关羽 2
	TexChessmanHGuanYu3     // 棋子 横 关羽 3
	TexChessmanHHuangZhong1 // 棋子 横 黄忠 1
	TexChessmanHHuangZhong2 // 棋子 横 黄忠 2
	TexChessmanHHuangZhong3 // 棋子 横 黄忠 3
	TexChessmanHMaChao1     // 棋子 横 马超 1
	TexChessmanHMaChao2     // 棋子 横 马超 2
	TexChessmanHMaChao3     // 棋子 横 马超 3
	TexChessmanHZhangFei1   // 棋子 横 张飞 1
	TexChessmanHZhangFei2   // 棋子 横 张飞 2
	TexChessmanHZhangFei3   // 棋子 横 张飞 3
	TexChessmanHZhaoYun1    // 棋子 横 赵云 1
	TexChessmanHZhaoYun2    // 棋子 横 赵云 2
	TexChessmanHZhaoYun3    // 棋子 横 赵云 3

	TexChessmanSGuanYu1     // 棋子 竖 关羽 1
	TexChessmanSGuanYu2     // 棋子 竖 关羽 2
	TexChessmanSGuanYu3     // 棋子 竖 关羽 3
	TexChessmanSHuangZhong1 // 棋子 竖 黄忠 1
	TexChessmanSHuangZhong2 // 棋子 竖 黄忠 2
	TexChessmanSHuangZhong3 // 棋子 竖 黄忠 3
	TexChessmanSMaChao1     // 棋子 竖 马超 1
	TexChessmanSMaChao2     // 棋子 竖 马超 2
	TexChessmanSMaChao3     // 棋子 竖 马超 3
	TexChessmanSZhangFei1   // 棋子 竖 张飞 1
	TexChessmanSZhangFei2   // 棋子 竖 张飞 2
	TexChessmanSZhangFei3   // 棋子 竖 张飞 3
	TexChessmanSZhaoYun1    // 棋子 竖 赵云 1
	TexChessmanSZhaoYun2    // 棋子 竖 赵云 2
	TexChessmanSZhaoYun3    // 棋子 竖 赵云 3

	TexChessmanA1 // 棋子 兵 甲 1
	TexChessmanA2 // 棋子 兵 甲 2
	TexChessmanA3 // 棋子 兵 甲 3
	TexChessmanB1 // 棋子 兵 乙 1
	TexChessmanB2 // 棋子 兵 乙 2
	TexChessmanB3 // 棋子 兵 乙 3
	TexChessmanC1 // 棋子 兵 丙 1
	TexChessmanC2 // 棋子 兵 丙 2
	TexChessmanC3 // 棋子 兵 丙 3
	TexChessmanD1 // 棋子 兵 丁 1
	TexChessmanD2 // 棋子 兵 丁 2
	TexChessmanD3 // 棋子 兵 丁 3

	TexWin // 过关时的显示内容素材
)

// 加载纹理图,多张纹理
func LoadTextures(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("game.png")
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

	return []sprite.SubTex{
		TexGameArea:   sprite.SubTex{t, image.Rect(0, 0, 1020, 1350)},
		TexBtnReturn1: sprite.SubTex{t, image.Rect(1100, 0, 1100+240, 120)},
		TexBtnReturn2: sprite.SubTex{t, image.Rect(1100, 125, 1100+240, 125+120)},
		TexBtnReturn3: sprite.SubTex{t, image.Rect(1100, 250, 1100+240, 250+120)},
		TexBtnGuide1:  sprite.SubTex{t, image.Rect(1100, 375, 1100+240, 375+120)},
		TexBtnGuide2:  sprite.SubTex{t, image.Rect(1100, 500, 1100+240, 500+120)},
		TexBtnGuide3:  sprite.SubTex{t, image.Rect(1100, 625, 1100+240, 625+120)},
		TexBtnReload1: sprite.SubTex{t, image.Rect(1100, 750, 1100+240, 750+120)},
		TexBtnReload2: sprite.SubTex{t, image.Rect(1100, 875, 1100+240, 875+120)},
		TexBtnReload3: sprite.SubTex{t, image.Rect(1100, 1000, 1100+240, 1000+120)},

		TexChessmanCaoCao1: sprite.SubTex{t, image.Rect(0, 1400, 0+480, 1400+480)},
		TexChessmanCaoCao2: sprite.SubTex{t, image.Rect(484, 1400, 484+480, 1400+480)},
		TexChessmanCaoCao3: sprite.SubTex{t, image.Rect(968, 1400, 968+480, 1400+480)},

		TexChessmanHGuanYu1:     sprite.SubTex{t, image.Rect(0, 1900, 0+480, 1900+240)},
		TexChessmanHGuanYu2:     sprite.SubTex{t, image.Rect(484, 1900, 484+480, 1900+240)},
		TexChessmanHGuanYu3:     sprite.SubTex{t, image.Rect(968, 1900, 960+480, 1900+240)},
		TexChessmanHHuangZhong1: sprite.SubTex{t, image.Rect(0, 2150, 0+480, 2150+240)},
		TexChessmanHHuangZhong2: sprite.SubTex{t, image.Rect(484, 2150, 484+480, 2150+240)},
		TexChessmanHHuangZhong3: sprite.SubTex{t, image.Rect(968, 2150, 960+480, 2150+240)},
		TexChessmanHMaChao1:     sprite.SubTex{t, image.Rect(0, 2400, 0+480, 2400+240)},
		TexChessmanHMaChao2:     sprite.SubTex{t, image.Rect(484, 2400, 484+480, 2400+240)},
		TexChessmanHMaChao3:     sprite.SubTex{t, image.Rect(968, 2400, 968+480, 2400+240)},
		TexChessmanHZhangFei1:   sprite.SubTex{t, image.Rect(0, 2650, 0+480, 2650+240)},
		TexChessmanHZhangFei2:   sprite.SubTex{t, image.Rect(484, 2650, 484+480, 2650+240)},
		TexChessmanHZhangFei3:   sprite.SubTex{t, image.Rect(968, 2650, 968+480, 2650+240)},
		TexChessmanHZhaoYun1:    sprite.SubTex{t, image.Rect(0, 2900, 0+480, 2900+240)},
		TexChessmanHZhaoYun2:    sprite.SubTex{t, image.Rect(484, 2900, 484+480, 2900+240)},
		TexChessmanHZhaoYun3:    sprite.SubTex{t, image.Rect(968, 2900, 968+480, 2900+240)},

		TexChessmanSGuanYu1:     sprite.SubTex{t, image.Rect(0, 3150, 0+240, 3150+480)},
		TexChessmanSGuanYu2:     sprite.SubTex{t, image.Rect(242, 3150, 242+240, 3150+480)},
		TexChessmanSGuanYu3:     sprite.SubTex{t, image.Rect(484, 3150, 484+240, 3150+480)},
		TexChessmanSHuangZhong1: sprite.SubTex{t, image.Rect(726, 3150, 726+240, 3150+480)},
		TexChessmanSHuangZhong2: sprite.SubTex{t, image.Rect(968, 3150, 968+240, 3150+480)},
		TexChessmanSHuangZhong3: sprite.SubTex{t, image.Rect(1210, 3150, 1210+240, 3150+480)},
		TexChessmanSMaChao1:     sprite.SubTex{t, image.Rect(0, 3650, 0+240, 3650+480)},
		TexChessmanSMaChao2:     sprite.SubTex{t, image.Rect(242, 3650, 242+240, 3650+480)},
		TexChessmanSMaChao3:     sprite.SubTex{t, image.Rect(484, 3650, 484+240, 3650+480)},
		TexChessmanSZhangFei1:   sprite.SubTex{t, image.Rect(726, 3650, 726+240, 3650+480)},
		TexChessmanSZhangFei2:   sprite.SubTex{t, image.Rect(968, 3650, 968+240, 3650+480)},
		TexChessmanSZhangFei3:   sprite.SubTex{t, image.Rect(1210, 3650, 1210+240, 3650+480)},
		TexChessmanSZhaoYun1:    sprite.SubTex{t, image.Rect(0, 4150, 0+240, 4150+480)},
		TexChessmanSZhaoYun2:    sprite.SubTex{t, image.Rect(242, 4150, 242+240, 4150+480)},
		TexChessmanSZhaoYun3:    sprite.SubTex{t, image.Rect(484, 4150, 484+240, 4150+480)},

		TexChessmanA1: sprite.SubTex{t, image.Rect(0, 4650, 0+240, 4650+240)},
		TexChessmanA2: sprite.SubTex{t, image.Rect(242, 4650, 242+240, 4650+240)},
		TexChessmanA3: sprite.SubTex{t, image.Rect(484, 4650, 484+240, 4650+240)},
		TexChessmanB1: sprite.SubTex{t, image.Rect(726, 4650, 726+240, 4650+240)},
		TexChessmanB2: sprite.SubTex{t, image.Rect(968, 4650, 968+240, 4650+240)},
		TexChessmanB3: sprite.SubTex{t, image.Rect(1210, 4650, 1210+240, 4650+240)},
		TexChessmanC1: sprite.SubTex{t, image.Rect(0, 4900, 0+240, 4900+240)},
		TexChessmanC2: sprite.SubTex{t, image.Rect(242, 4900, 242+240, 4900+240)},
		TexChessmanC3: sprite.SubTex{t, image.Rect(484, 4900, 484+240, 4900+240)},
		TexChessmanD1: sprite.SubTex{t, image.Rect(726, 4900, 726+240, 4900+240)},
		TexChessmanD2: sprite.SubTex{t, image.Rect(968, 4900, 968+240, 4900+240)},
		TexChessmanD3: sprite.SubTex{t, image.Rect(1210, 4900, 1210+240, 4900+240)},

		TexWin: sprite.SubTex{t, image.Rect(726, 4150, 726+480, 4150+240)},
	}

}

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

// 每个棋子应该用那个纹理来绘制
// name 棋子的名称， success 游戏是否结束
// relWidth 决定棋子是横版 还是竖版
// status 棋子目前的状态
// t 当前的时间，用于轮播，  d 多长时间轮播一次
func ChessManFrame(name rune, success bool, relWidth int, status level.ChessManStatus, t, d clock.Time) int {

	var frames []int // 那个棋子的判断
	switch name {
	case '曹':
		frames = []int{TexChessmanCaoCao1, TexChessmanCaoCao2, TexChessmanCaoCao3}
	case '甲':
		frames = []int{TexChessmanA1, TexChessmanA2, TexChessmanA3}
	case '乙':
		frames = []int{TexChessmanB1, TexChessmanB2, TexChessmanB3}
	case '丙':
		frames = []int{TexChessmanC1, TexChessmanC2, TexChessmanC3}
	case '丁':
		frames = []int{TexChessmanD1, TexChessmanD2, TexChessmanD3}
	case '关':
		if relWidth == 1 {
			frames = []int{TexChessmanSGuanYu1, TexChessmanSGuanYu2, TexChessmanSGuanYu3}
		} else {
			frames = []int{TexChessmanHGuanYu1, TexChessmanHGuanYu2, TexChessmanHGuanYu3}
		}
	case '张':
		if relWidth == 1 {
			frames = []int{TexChessmanSZhangFei1, TexChessmanSZhangFei2, TexChessmanSZhangFei3}
		} else {
			frames = []int{TexChessmanHZhangFei1, TexChessmanHZhangFei2, TexChessmanHZhangFei3}
		}
	case '赵':
		if relWidth == 1 {
			frames = []int{TexChessmanSZhaoYun1, TexChessmanSZhaoYun2, TexChessmanSZhaoYun3}
		} else {
			frames = []int{TexChessmanHZhaoYun1, TexChessmanHZhaoYun2, TexChessmanHZhaoYun3}
		}
	case '马':
		if relWidth == 1 {
			frames = []int{TexChessmanSMaChao1, TexChessmanSMaChao2, TexChessmanSMaChao3}
		} else {
			frames = []int{TexChessmanHMaChao1, TexChessmanHMaChao2, TexChessmanHMaChao3}
		}
	case '黄':
		if relWidth == 1 {
			frames = []int{TexChessmanSHuangZhong1, TexChessmanSHuangZhong2, TexChessmanSHuangZhong3}
		} else {
			frames = []int{TexChessmanHHuangZhong1, TexChessmanHHuangZhong2, TexChessmanHHuangZhong3}
		}
	}
	if success { // 游戏结束，所有棋子都不能动了。
		return frames[0]
	}

	if status == level.ChessManStable {
		// 不可移动
		return frames[0]
	} else if status == level.ChessManMovable {
		// 可移动
		a := int(d) * 2
		b := (int(t) % a) / int(d)
		if b == 0 {
			return frames[0]
		} else {
			return frames[1]
		}
	} else if status == level.ChessManMoving {
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
