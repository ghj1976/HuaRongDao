package model

import (
	"github.com/ghj1976/HuaRongDao/level"
)

// 是不是空白位置
func ChessManIsBlank(chessman rune) bool {
	if chessman == level.BlackChessManPos {
		return true
	} else {
		return false
	}
}
