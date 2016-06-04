// 列表页的模型类
package model

import (
	"log"

	"github.com/ghj1976/HuaRongDao/button"
	"github.com/ghj1976/HuaRongDao/level"
	"golang.org/x/mobile/event/size"
)

// 列表实体类
type ListModel struct {
	BtnPrePage  *button.GameBtn   // 前一页按钮
	BtnNextPage *button.GameBtn   // 下一页按钮
	Arr         []level.LevelInfo // 当前页显示的数据，注意这些是有位置信息的。

	currPage       int // 当前所在页面
	pageSize       int // 页面尺寸
	horizontalSize int // 水平方向可以放几个列表项
	verticalSize   int // 垂直方向可以放几个列表项
}

const (
	buttonAreaHeight = 50.0
	levelAreaHeight  = 120.0
	levelAreaWidth   = 72.0
)

// 构造函数
func NewListModel() *ListModel {
	lm := &ListModel{}

	lm.BtnPrePage = &button.GameBtn{Visible: false, Status: button.BtnNormal}
	lm.BtnNextPage = &button.GameBtn{Visible: false, Status: button.BtnNormal}

	return lm
}

// 知道屏幕尺寸大小基础下，计算出每个元素的位置和水平、垂直方向的个数。
func (lm *ListModel) InitListSize(sz size.Event) {
	lm.BtnPrePage.Visible = false
	lm.BtnNextPage.Visible = false

	lm.verticalSize = int((float32(sz.HeightPt) - buttonAreaHeight) / levelAreaHeight)
	lm.horizontalSize = int(float32(sz.WidthPt) / levelAreaWidth)
	log.Println("布局：", lm.horizontalSize, "*", lm.verticalSize)

	// 准备好了，可以开始启用按钮功能了
	lm.BtnPrePage.Visible = true
	lm.BtnNextPage.Visible = true
}

// 下一页的数据准备
func (lm *ListModel) NextPage() {

}

// 前一页的数据准备
func (lm *ListModel) PrePage() {

}
