// 列表页面需要的纹理、图片初始化
// 绘图技术请参考： http://www.cnblogs.com/ghj1976/p/3443638.html
// 颜色对照表： http://tool.oschina.net/commons?type=3
package textures

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/ghj1976/HuaRongDao/level"
)

var (
	// 常用的几种颜色
	bgColor    color.RGBA = color.RGBA{102, 205, 170, 255}
	blueColor  color.RGBA = color.RGBA{0, 0, 255, 255}
	redColor   color.RGBA = color.RGBA{255, 0, 0, 255}
	greenColor color.RGBA = color.RGBA{0, 255, 0, 255}
)

// 计算绘图的大小
func GetBounds(d int) (borderWidth, chessManWidth, areaWidth, areaHeight int) {
	borderWidth = 3 * d
	chessManWidth = borderWidth*6 + 2*d
	areaWidth = 2*borderWidth + 4*chessManWidth
	areaHeight = 2*borderWidth + 5*chessManWidth
	return
}

// 绘出每个布局的缩略效果图
func LevelRGBA(d int, level *level.LevelInfo) *image.RGBA {
	borderWidth, chessManWidth, areaWidth, areaHeight := GetBounds(d)

	// 绘图区域创建
	m := image.NewRGBA(image.Rect(0, 0, areaWidth, areaHeight))

	// 画总体外框
	draw.Draw(m, m.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)

	var currColor *color.RGBA
	// 画每个棋子
	for _, cm := range level.ChessMans {

		if cm.RelHeight == 2 && cm.RelWidth == 2 {
			currColor = &redColor
		} else if cm.RelHeight == 1 && cm.RelWidth == 1 {
			currColor = &greenColor
		} else {
			currColor = &blueColor
		}

		draw.Draw(m,
			image.Rect(borderWidth+chessManWidth*cm.RelLeftTopX+d,
				borderWidth+chessManWidth*cm.RelLeftTopY+d,
				borderWidth+chessManWidth*(cm.RelRightBottomX+1)-d,
				borderWidth+chessManWidth*(cm.RelRightBottomY+1)-d),
			&image.Uniform{currColor}, image.ZP, draw.Src)

	}
	return m
}
