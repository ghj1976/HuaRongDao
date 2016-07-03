// 列表页面需要的纹理、图片初始化
// 绘图技术请参考： http://www.cnblogs.com/ghj1976/p/3443638.html
// 颜色对照表： http://tool.oschina.net/commons?type=3
package textures

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/ghj1976/HuaRongDao/button"
	"github.com/ghj1976/HuaRongDao/level"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite"
)

var (
	// 常用的几种颜色
	bg1Color  color.RGBA = color.RGBA{255, 165, 0, 255}   // 最外面外框的颜色
	bg11Color color.RGBA = color.RGBA{255, 0, 0, 255}     // 最外面外框的颜色
	bg2Color  color.RGBA = color.RGBA{102, 205, 170, 255} // 游戏区域的外框颜色
	// 三种棋子的颜色
	blueColor  color.RGBA = color.RGBA{0, 0, 255, 255}
	redColor   color.RGBA = color.RGBA{255, 0, 0, 255}
	greenColor color.RGBA = color.RGBA{0, 255, 0, 255}

	// 写最小步数文字的字体颜色
	fontColor11 color.RGBA = color.RGBA{32, 178, 170, 255}
	fontColor12 color.RGBA = color.RGBA{255, 0, 0, 255}
	fontColor13 color.RGBA = color.RGBA{0, 0, 255, 255}
)

// 计算绘图的大小
func GetBounds(d int) (borderWidth, chessManWidth, areaWidth, areaHeight int) {
	borderWidth = 3 * d
	chessManWidth = borderWidth*6 + 2*d
	areaWidth = 4*borderWidth + 2*borderWidth + 4*chessManWidth
	areaHeight = 3*chessManWidth + 2*borderWidth + 5*chessManWidth
	return
}

// 绘制出指定关卡数组的纹理图
func InitListTexMap(eng sprite.Engine, levelArr []*level.LevelInfo, levelMap map[string]sprite.SubTex) {

	for _, le := range levelArr {
		keyd := fmt.Sprintf("%d-%d-d", le.RelX, le.RelY) // 默认显示的纹理
		log.Println("tex:", le.Name, keyd)
		if _, ok := levelMap[keyd]; !ok {
			m1 := LevelRGBA(2, le, false)
			t1, err := eng.LoadTexture(m1)
			if err != nil {
				log.Fatal(err)
			}
			levelMap[keyd] = sprite.SubTex{t1, m1.Rect}
		}

		keys := fmt.Sprintf("%d-%d-s", le.RelX, le.RelY) // 选中状态的纹理
		log.Println("tex:", le.Name, keys)
		if _, ok := levelMap[keys]; !ok {
			m2 := LevelRGBA(2, le, true)
			t2, err := eng.LoadTexture(m2)
			if err != nil {
				log.Fatal(err)
			}
			levelMap[keys] = sprite.SubTex{t2, m2.Rect}
		}
	}
}

// 绘出每个布局的缩略效果图
// check 是否选中的
func LevelRGBA(d int, le *level.LevelInfo, check bool) *image.RGBA {
	borderWidth, chessManWidth, areaWidth, areaHeight := GetBounds(d)

	// 绘图区域创建
	m := image.NewRGBA(image.Rect(0, 0, areaWidth, areaHeight))

	// 画所有区域的外框, 不同状态的关卡，背景颜色不一样。
	if check {
		draw.Draw(m, m.Bounds(), &image.Uniform{bg11Color}, image.ZP, draw.Src)
	} else {
		draw.Draw(m, m.Bounds(), &image.Uniform{bg1Color}, image.ZP, draw.Src)
	}
	topy1 := 3 * chessManWidth / 2
	// 画游戏区域外框
	draw.Draw(m,
		image.Rect(2*borderWidth,
			topy1,
			areaWidth-2*borderWidth,
			areaHeight-topy1),
		&image.Uniform{bg2Color}, image.ZP, draw.Src)

	var currColor *color.RGBA
	// 画每个棋子
	for _, cm := range le.ChessMans {
		if cm.RelHeight == 2 && cm.RelWidth == 2 {
			currColor = &redColor
		} else if cm.RelHeight == 1 && cm.RelWidth == 1 {
			currColor = &greenColor
		} else {
			currColor = &blueColor
		}
		log.Println("chessman draw:", cm.Name, currColor)

		draw.Draw(m,
			image.Rect(3*borderWidth+chessManWidth*cm.RelLeftTopX+d,
				topy1+borderWidth+chessManWidth*cm.RelLeftTopY+d,
				3*borderWidth+chessManWidth*(cm.RelRightBottomX+1)-d,
				topy1+borderWidth+chessManWidth*(cm.RelRightBottomY+1)-d),
			&image.Uniform{currColor}, image.ZP, draw.Src)

	}

	txtColor := color.RGBA{0, 0, 255, 255} // RGBA, 不透明 A 为 255

	// 写关卡名称
	cpt1 := image.Point{X: areaWidth / 2, Y: 40}
	DrawString(m, 35.0, txtColor, cpt1, le.Name)

	var fontColor1 color.RGBA
	if le.LevelStatus == level.LevelPass {
		fontColor1 = fontColor12
	} else if le.LevelStatus == level.LevelBestPass {
		fontColor1 = fontColor13
	} else {
		le.LevelStatus = level.LevelNotPass
		fontColor1 = fontColor11
	}
	// 写最小步数
	cpt2 := image.Point{X: areaWidth / 2, Y: areaHeight - 15}

	var stepStr string
	if le.LevelStatus == level.LevelPass {
		stepStr = fmt.Sprintf("过关：%d/%d", 0, le.MinStepNum)
	} else if le.LevelStatus == level.LevelBestPass {
		stepStr = fmt.Sprintf("最佳：%d/%d", 0, le.MinStepNum)
	} else {
		stepStr = fmt.Sprintf("%d/%d", 0, le.MinStepNum)
	}
	DrawString(m, 35.0, fontColor1, cpt2, stepStr)

	return m
}

const (
	texBtnPre1  = iota // 前一页按钮 1
	texBtnPre2         // 前一页按钮 2
	texBtnNext1        // 下一页按钮 1
	texBtnNext2        // 下一页按钮 2
)

// 加载纹理图,多张纹理
func LoadTexturesList(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("list.png")
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
		texBtnPre1:  sprite.SubTex{t, image.Rect(0, 0, 110, 110)},
		texBtnPre2:  sprite.SubTex{t, image.Rect(110, 0, 220, 110)},
		texBtnNext1: sprite.SubTex{t, image.Rect(220, 0, 330, 110)},
		texBtnNext2: sprite.SubTex{t, image.Rect(330, 0, 440, 110)},
	}

}

// 获得游戏页面的按钮用哪个纹理图。
func ListButtonFrame(name string, status button.BtnStatus) int {
	switch name {
	case "pre":
		if status == button.BtnPress {
			return texBtnPre2
		} else {
			return texBtnPre1
		}
	case "next":
		if status == button.BtnPress {
			return texBtnNext2
		} else {
			return texBtnNext1
		}

	default:
		return 0
	}
}
