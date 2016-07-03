// 列表页的模型类
package model

import (
	"log"

	"github.com/ghj1976/HuaRongDao/button"
	"github.com/ghj1976/HuaRongDao/common"
	"github.com/ghj1976/HuaRongDao/db"
	"github.com/ghj1976/HuaRongDao/level"
	"golang.org/x/mobile/event/size"
)

// 列表实体类
type ListModel struct {
	BtnPrePage  *button.GameBtn    // 前一页按钮
	BtnNextPage *button.GameBtn    // 下一页按钮
	Arr         []*level.LevelInfo // 当前页显示的数据，注意这些是有位置信息的。 指针的数组，才能保障可以修改这里的值。

	currPage        int     // 当前所在页面
	pageSize        int     // 页面尺寸
	horizontalSize  int     // 水平方向可以放几个列表项
	horizontalSpace float32 // 水平方向的间距
	verticalSize    int     // 垂直方向可以放几个列表项
	verticalSpace   float32 // 垂直方向的间距
}

const (
	buttonAreaHeight = 50.0 // 按钮区域的高度
	buttonHeight     = 48.0 // 按钮的高度

	levelAreaHeight = 120.0
	levelAreaWidth  = 72.0
)

// 构造函数
func NewListModel() *ListModel {
	lm := &ListModel{}

	lm.BtnPrePage = &button.GameBtn{Visible: false, Status: button.BtnNormal}
	lm.BtnNextPage = &button.GameBtn{Visible: false, Status: button.BtnNormal}

	lm.currPage = 1

	return lm
}

// 知道屏幕尺寸大小基础下，计算出每个元素的位置和水平、垂直方向的个数。
func (lm *ListModel) InitListSizeAndData(sz size.Event) {
	lm.BtnPrePage.Visible = false
	lm.BtnNextPage.Visible = false

	lm.BtnPrePage.SetGameRectangle(
		common.GamePoint{
			X: 1.0,
			Y: 1.0,
		},
		buttonHeight,
		buttonHeight)

	lm.BtnNextPage.SetGameRectangle(
		common.GamePoint{
			X: float32(sz.WidthPt) - buttonAreaHeight + 1.0,
			Y: 1.0,
		},
		buttonHeight,
		buttonHeight)

	lm.verticalSize = int((float32(sz.HeightPt) - buttonAreaHeight) / levelAreaHeight)
	lm.horizontalSize = int(float32(sz.WidthPt) / levelAreaWidth)
	log.Println("布局：", lm.horizontalSize, "*", lm.verticalSize)
	lm.pageSize = lm.verticalSize * lm.horizontalSize

	lm.verticalSpace = (float32(sz.HeightPt) - float32(buttonAreaHeight) - levelAreaHeight*float32(lm.verticalSize)) / float32(lm.verticalSize+1)
	lm.horizontalSpace = (float32(sz.WidthPt) - levelAreaWidth*float32(lm.horizontalSize)) / float32(lm.horizontalSize+1)

	// 必须在确定页面尺寸大小后，才知道需要初始化多少数据。
	arr, hasPrePage, hasNextPage := db.ReadPage(lm.currPage, lm.pageSize, "")
	for _, levv := range arr {
		levv.Reset()
	}
	lm.Arr = arr

	// 计算每个关卡的具体位置
	var x, y int
	x = 0
	y = 0
	for _, lev := range lm.Arr {
		lev.Rect = common.GameRectangle{}
		lev.RelX = x
		lev.RelY = y
		lev.Rect.SetGameRectangle(
			common.GamePoint{
				X: float32(x)*(lm.horizontalSpace+levelAreaWidth) + lm.horizontalSpace,
				Y: float32(buttonAreaHeight) + float32(y)*(lm.verticalSpace+levelAreaHeight) + lm.verticalSpace,
			},
			levelAreaWidth,
			levelAreaHeight)
		if x >= lm.horizontalSize-1 {
			x = 0
			y++
		} else {
			x++
		}

		log.Printf("init:%s-%d-%d", lev.Name, lev.RelX, lev.RelY)

		if y >= lm.verticalSize {
			break
		}
	}

	//	for _, lev := range lm.Arr {
	//		log.Println("model:", lev.Name)
	//	}

	// 准备好了，可以开始启用按钮功能了
	lm.BtnPrePage.Visible = hasPrePage
	lm.BtnNextPage.Visible = hasNextPage

}

// 下一页的数据准备
func (lm *ListModel) NextPage() {

}

// 前一页的数据准备
func (lm *ListModel) PrePage() {

}
