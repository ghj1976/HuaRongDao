package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/ghj1976/HuaRongDao/level"
)

func main() {
	// 待处理的数据
	levelArr := level.InitData()

	// 绘图区域创建
	m := image.NewRGBA(image.Rect(0, 0, 640, 480))

	// 常用的几种颜色
	bg := color.RGBA{171, 190, 62, 255}
	blue := color.RGBA{0, 0, 255, 255}
	//	red := color.RGBA{255, 0, 0, 255}
	//	green := color.RGBA{0, 255, 0, 255}

	// 绘图
	draw.Draw(m, m.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)
	draw.Draw(m, image.Rect(100, 100, 200, 300), &image.Uniform{blue}, image.ZP, draw.Src)

	// 保存到文件
	imgfile, _ := os.OpenFile("001.png", os.O_RDWR|os.O_CREATE, 0660)
	defer imgfile.Close()
	png.Encode(imgfile, m)
}
