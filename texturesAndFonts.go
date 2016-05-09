package main

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
	"golang.org/x/mobile/exp/sprite/clock"
)

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

	texWin // 过关时的显示内容素材
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

		texWin: sprite.SubTex{t, image.Rect(726, 4150, 726+480, 4150+240)},
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

// 初始化时候加载字体
func loadGameFont() error {
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

func initFontText(txtSize float64, fontColor color.RGBA, rect image.Rectangle) {
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
func loadFontTextTextures(eng sprite.Engine, txt string, rect image.Rectangle) sprite.SubTex {
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
func chessManFrame(name rune, success bool, relWidth int, status ChessManStatus, t, d clock.Time) int {

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
	if success { // 游戏结束，所有棋子都不能动了。
		return frames[0]
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
