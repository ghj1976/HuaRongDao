// 列表类的视图类
package view

import (
	"fmt"
	"log"

	"github.com/ghj1976/HuaRongDao/button"
	"github.com/ghj1976/HuaRongDao/common"
	"github.com/ghj1976/HuaRongDao/model"
	"github.com/ghj1976/HuaRongDao/textures"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

type ListView struct {
	model        *model.ListModel         // 游戏列表模型类
	RootViewNode *sprite.Node             // 游戏视图的根节点
	mapTex       map[string]sprite.SubTex // 游戏列表页的关卡纹理，key 是 1-3 位置的字符串
}

// 构造一个 list 视图类
func NewListView(m *model.ListModel, eng sprite.Engine) *ListView {
	lv := &ListView{}
	lv.model = m

	// 计算每个元素最终的显示位置。
	sz := model.GetScreenSizeBlock()
	lv.model.InitListSizeAndData(sz)
	lv.loadListView(eng)

	return lv
}

// 加载列表页面
func (lv *ListView) loadListView(eng sprite.Engine) {
	lv.RootViewNode = &sprite.Node{} // View 的绘图根节点
	eng.Register(lv.RootViewNode)
	eng.SetTransform(lv.RootViewNode, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	newNode := func(fn common.ArrangerFunc) {
		n := &sprite.Node{Arranger: common.ArrangerFunc(fn)}
		eng.Register(n)
		lv.RootViewNode.AppendChild(n)
	}

	err := textures.LoadGameFont("")
	if err != nil {
		log.Panicln(err)
		return
	}

	texs := textures.LoadTexturesList(eng)

	// 上一页按钮
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if lv.model.BtnPrePage.Visible {
			eng.SetSubTex(n, texs[textures.ListButtonFrame("pre", lv.model.BtnPrePage.Status)])
			eng.SetTransform(n, lv.model.BtnPrePage.ToF32Affine())
		}
	})

	// 下一页按钮
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		if lv.model.BtnNextPage.Visible {
			eng.SetSubTex(n, texs[textures.ListButtonFrame("next", lv.model.BtnNextPage.Status)])
			eng.SetTransform(n, lv.model.BtnNextPage.ToF32Affine())
		}
	})

	// 初始化需要缓存的每个关卡的纹理图Map
	levelTexs := map[string]sprite.SubTex{}
	textures.InitListTexMap(eng, lv.model.Arr, levelTexs)
	log.Println("levelTexs len:", len(levelTexs))

	// 初始化关卡信息
	for _, lev := range lv.model.Arr {
		keyd := fmt.Sprintf("%d-%d-d", lev.RelX, lev.RelY)
		log.Println("find:", keyd)

		levv := lev // 注意，newNode 内部不能用 lev， 这样会指针指向混乱， 所以 额外用了一个局部变量。
		newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
			eng.SetSubTex(n, levelTexs[keyd])
			eng.SetTransform(n, levv.Rect.ToF32Affine())
		})

	}

}

// 每次绘图前，逻辑相关的操作。
func (lv *ListView) Update(now clock.Time) {
}

// 当 touch 事件发生时， 判断是按在那个游戏精灵元素上，以及对应的处理策略分支。
func (lv *ListView) Press(touchEvent touch.Event) {
	sz, _ := model.GetScreenSize()
	// 单位修改成 pt， 而不是 px
	gp := common.GamePoint{X: touchEvent.X / sz.PixelsPerPt, Y: touchEvent.Y / sz.PixelsPerPt}

	// 按钮 按下逻辑处理。
	if touchEvent.Type == touch.TypeBegin {
		if gp.In(lv.model.BtnNextPage.GameRectangle) {
			// 按钮被点击
			lv.model.BtnNextPage.Status = button.BtnPress
			log.Println("BtnNextPage 被按下")
			return
		} else if gp.In(lv.model.BtnPrePage.GameRectangle) {
			lv.model.BtnPrePage.Status = button.BtnPress
			log.Println("BtnPrePage 被按下")
			return
		}
	} else if touchEvent.Type == touch.TypeEnd {
		if lv.model.BtnNextPage.Status == button.BtnPress {
			// 按钮被释放
			lv.model.BtnNextPage.Status = button.BtnNormal
			log.Println("BtnNextPage 释放按下状态")
			// 后一页按钮的操作逻辑
			return
		} else if lv.model.BtnPrePage.Status == button.BtnPress {
			lv.model.BtnPrePage.Status = button.BtnNormal
			log.Println("BtnPrePage 释放按下状态")
			// 前一页按钮的操作逻辑
			return
		}

	}

	// 判断是那个关卡被点击

	// 关卡被点击的处理逻辑

}

func (lv *ListView) OnScreenSizeChange(currSZ size.Event, displayMultiple float32) {
	lv.model.OnScreenSizeChange(currSZ, displayMultiple)
}
