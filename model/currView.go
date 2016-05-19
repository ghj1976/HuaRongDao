package model

import (
	"golang.org/x/mobile/event/size"
)

var (
	screenSize size.Event // 屏幕尺寸
	chSZ       chan bool  // 是否已经正确的获得屏幕尺寸了
)

// 初始化屏幕尺寸
func InitScreenSize(sz size.Event) {
	if sz.HeightPt == 0 {
		return
	}
	screenSize = sz
	//log.Println("2334：", sz)
	// 多次发送也不怕，每次都 构造一个新的。
	chSZ = make(chan bool, 1)
	chSZ <- true
	//log.Println("2335：", sz)
}

// 获取屏幕尺寸， has 标示是否已经被成功设置了,不阻塞的方法
func GetScreenSize() (sz size.Event, has bool) {
	if screenSize.HeightPt == 0 {
		return screenSize, false
	} else {
		return screenSize, true
	}
}

// 阻塞获得屏幕尺寸，如果无法获得，就一直阻塞中
func GetScreenSizeBlock() size.Event {
	if screenSize.HeightPt == 0 {
		<-chSZ
	}
	return screenSize
}
