package model

import (
	"log"

	"golang.org/x/mobile/event/size"
)

var (
	firstScreenSize size.Event // 第一次成功计算出的屏幕尺寸。
	currScreenSize  size.Event // 当前屏幕尺寸
	chSZ            chan bool  // 是否已经正确的获得屏幕尺寸了
	displayMultiple float32    // 显示倍数，来回切换屏幕时，屏幕分辨率会发生变化，为了确保游戏不发生变化，所以需要这个显示倍数
)

/*
+ 同样的电脑，返回的屏幕尺寸 PT 值 不一样
	ScreenSize: {400 400 351.22pt 351.22pt 1.1388888 0}
	ScreenSize: {400 400 252.63pt 252.63pt 1.5833334 0}
	通过model中定义了一个 displayMultiple， 在每次 ScreenSizeChange 变化是，重算位置时使用。
*/

// 初始化屏幕尺寸
func InitScreenSize(sz size.Event) {
	log.Println("ScreenSize:", sz)
	if sz.HeightPt == 0 {
		return
	}
	currScreenSize = sz
	if firstScreenSize.HeightPt == 0 {
		firstScreenSize = sz
		displayMultiple = 1.0
	} else {
		displayMultiple = firstScreenSize.PixelsPerPt / currScreenSize.PixelsPerPt
	}
	log.Println("firstScreenSize:", firstScreenSize, "currScreenSize:", currScreenSize, "displayMultiple:", displayMultiple)

	//log.Println("2334：", sz)
	// 多次发送也不怕，每次都 构造一个新的。
	chSZ = make(chan bool, 1)
	chSZ <- true
	//log.Println("2335：", sz)
}

// 获取屏幕尺寸， has 标示是否已经被成功设置了,不阻塞的方法
func GetScreenSize() (sz size.Event, has bool) {
	if currScreenSize.HeightPt == 0 {
		return currScreenSize, false
	} else {
		return currScreenSize, true
	}
}

// 阻塞获得屏幕尺寸，如果无法获得，就一直阻塞中
func GetScreenSizeBlock() size.Event {
	if currScreenSize.HeightPt == 0 {
		<-chSZ
	}
	return currScreenSize
}

// 返回显示倍数
func GetDisplayMultiple() float32 {
	return displayMultiple
}
