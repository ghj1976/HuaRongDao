package main

import (
	"testing"
)

func Test_Layout2Map(t *testing.T) {
	layout := `	张曹曹马
				张曹曹马
				黄关关赵
				黄甲乙赵
				丙一一丁`

	li := InitLevel("横刀立马", layout, 0)

	cmc := li.ChessMans['曹']
	if cmc.RelLeftTopX != 1 ||
		cmc.RelLeftTopY != 0 ||
		cmc.RelRightBottomX != 2 ||
		cmc.RelRightBottomY != 1 ||
		cmc.RelHeight != 2 ||
		cmc.RelWidth != 2 {
		t.Error("曹操的棋子位置计算错误", cmc)
	}

	if cmc.status != ChessManStable {
		t.Error("曹操初始可移动状态计算错误", cmc)
	}

	cmA := li.ChessMans['甲']
	if cmA.status != ChessManMovable {
		t.Error("甲兵初始可移动状态计算错误", cmA)
	}

	cmB := li.ChessMans['乙']
	if cmB.status != ChessManMovable {
		t.Error("乙兵初始可移动状态计算错误", cmB)
	}

	cmC := li.ChessMans['丙']
	if cmC.status != ChessManMovable {
		t.Error("丙兵初始可移动状态计算错误", cmC)
	}

	cmD := li.ChessMans['丁']
	if cmD.status != ChessManMovable {
		t.Error("丁兵初始可移动状态计算错误", cmD)
	}

}
