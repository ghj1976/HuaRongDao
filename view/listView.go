// 列表类的视图类
package view

import (
	"log"

	"github.com/ghj1976/HuaRongDao/common"
	"github.com/ghj1976/HuaRongDao/model"
	"github.com/ghj1976/HuaRongDao/textures"

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

	// 初始化关卡信息
	for _, lev := range lv.model.Arr {

		// 需要缓存的每个关卡的纹理图
		tex := 
		log.Println("view:", lev.Name)
	}

}

// 每次绘图前，逻辑相关的操作。
func (lv *ListView) Update(now clock.Time) {
}

// 当 touch 事件发生时， 判断是按在那个游戏精灵元素上，以及对应的处理策略分支。
func (lv *ListView) Press(touchEvent touch.Event) {
}
