// 改变自go官方下面源码
// https://github.com/golang/go/blob/master/src/image/geom.go
// 修改的目的是之前点是 int 类型的，现在要改成 float32 类型的
// type Point struct {
//	X, Y int
// }
// 郭红俊 2016-04-19
package main

// 一个具体的点位置
type GamePoint struct {
	X, Y float32
}

// 长方形位置
type GameRectangle struct {
	LeftTop     GamePoint // 左上角的点坐标
	Width       float32   // 长方形的宽度
	Height      float32   // 长方形的高度
	RightBottom GamePoint // 右下角的点坐标
}

// 参数赋值
// 请使用这个赋值函数，因为参数赋值涉及到逻辑运算。
func (gr *GameRectangle) SetGameRectangle(lt GamePoint, w, h float32) {
	gr.LeftTop = lt
	gr.Width = w
	gr.Height = h
	gr.RightBottom = GamePoint{
		X: lt.X + w,
		Y: lt.Y + h,
	}
}

// 点p是否在r这个区域中？
func (p GamePoint) In(r GameRectangle) bool {
	return r.LeftTop.X <= p.X &&
		p.X <= r.RightBottom.X &&
		r.LeftTop.Y <= p.Y &&
		p.Y <= r.RightBottom.Y
}
