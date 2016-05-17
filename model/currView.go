package model

import (
	"golang.org/x/mobile/event/size"
)

type CurrView byte // 当前正在显示的View

const (
	SplashView CurrView = iota
	LoadingView
	ListView
	GameView
)

var (
	ScreenSize size.Event // 屏幕尺寸
)
