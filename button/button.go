package button

import (
	"github.com/ghj1976/HuaRongDao/common"
)

type BtnStatus byte // 按钮的状态枚举

const (
	BtnPress  BtnStatus = iota // 按钮被按下状态
	BtnNormal                  // 按钮正常状态
)

// 游戏中的按钮类
type GameBtn struct {
	Status               BtnStatus // 按钮的状态， 一共2种，按下、正常
	common.GameRectangle           // 按钮所在位置（长方形）
}
