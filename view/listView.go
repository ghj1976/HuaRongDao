package view // 列表类的视图类

import (
	"fmt"
	"log"

	"github.com/ghj1976/HuaRongDao/button"
	"github.com/ghj1976/HuaRongDao/common"
	"github.com/ghj1976/HuaRongDao/level"
	"github.com/ghj1976/HuaRongDao/model"
	"github.com/ghj1976/HuaRongDao/textures"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

// ListView 列表视图类
type ListView struct {
	// 列表页 游戏页 切换时  需要保留的对象
	eng          sprite.Engine // 绘制页面引擎
	RootViewNode *sprite.Node  // 游戏视图的根节点

	// 列表页 游戏页 切换时 不需要保留的对象
	model   *model.ListModel // 游戏列表模型类
	btnNext *sprite.Node     // 下一页按钮
	btnPre  *sprite.Node     // 前一页按钮

	// 上一页，下一页 切换时， 不需要保留的对象
	mapTex     map[string]sprite.SubTex // 游戏列表页的关卡纹理，key 是 1-3 位置的字符串
	levelNodes []*sprite.Node           // 所有关卡的node列表
}

// NewListView 构造一个 list 视图类
func NewListView(m *model.ListModel, eng sprite.Engine) *ListView {
	// 只负责构造 切换页面时也需要保留的内容。
	lv := &ListView{}
	lv.eng = eng
	lv.RootViewNode = &sprite.Node{} // View 的绘图根节点
	eng.Register(lv.RootViewNode)
	eng.SetTransform(lv.RootViewNode, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	lv.Reset(m)
	return lv
}

// Reset 重新加载，由于有些对象无法释放，返回时，用此方法
// 重新加载时，下面对象的指针不能变。
// lv.RootViewNode
func (lv *ListView) Reset(m *model.ListModel) {
	// 跟游戏页切换时不保留，但是 上下一页时保留的对象初始化

	// 数据实体对象构造
	lv.model = m
	sz := model.GetScreenSizeBlock() // 计算每个元素最终的显示位置。
	lv.model.InitListSizeAndData(sz)

	// 上一页、下一页 按钮构造
	texs := textures.LoadTexturesList(lv.eng)

	// 前一页按钮
	lv.btnPre = common.NewNodeNoShow(lv.eng, func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texs[textures.ListButtonFrame("pre", lv.model.BtnPrePage.Status)])
		eng.SetTransform(n, lv.model.BtnPrePage.ToF32Affine())
	})

	// 下一页按钮
	lv.btnNext = common.NewNodeNoShow(lv.eng, func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, texs[textures.ListButtonFrame("next", lv.model.BtnNextPage.Status)])
		eng.SetTransform(n, lv.model.BtnNextPage.ToF32Affine())
	})

	lv.refreshView(lv.eng)
}

// refreshView 分页时刷新页面
func (lv *ListView) refreshView(eng sprite.Engine) {
	// 分页时，刷新页面视图

	// 上一页按钮
	if lv.model.BtnPrePage.Visible && lv.btnPre.Parent == nil {
		lv.RootViewNode.AppendChild(lv.btnPre)
	} else if !lv.model.BtnPrePage.Visible && lv.btnPre.Parent != nil {
		lv.RootViewNode.RemoveChild(lv.btnPre)
	}

	// 下一页按钮
	if lv.model.BtnNextPage.Visible && lv.btnNext.Parent == nil {
		lv.RootViewNode.AppendChild(lv.btnNext)
	} else if !lv.model.BtnNextPage.Visible && lv.btnNext.Parent != nil {
		lv.RootViewNode.RemoveChild(lv.btnNext)
	}

	err := textures.LoadGameFont("")
	if err != nil {
		log.Panicln(err)
		return
	}

	newNode := func(fn common.ArrangerFunc) *sprite.Node {
		n := &sprite.Node{Arranger: common.ArrangerFunc(fn)}
		eng.Register(n)
		lv.RootViewNode.AppendChild(n)
		return n
	}

	// 初始化需要缓存的每个关卡的纹理图Map
	levelTexs := map[string]sprite.SubTex{}
	textures.InitListTexMap(eng, lv.model.Arr, levelTexs)
	log.Println("levelTexs len:", len(levelTexs))

	lll := len(lv.model.Arr)

	log.Println("page len:", lll)

	// 清除之前的数据
	if len(lv.levelNodes) > 0 {
		for _, nn := range lv.levelNodes {
			if nn.Parent != nil {
				lv.RootViewNode.RemoveChild(nn)
			}
		}
	}

	i := 0
	lv.levelNodes = make([]*sprite.Node, lll)
	// 初始化关卡信息
	for _, lev := range lv.model.Arr {
		keyd := fmt.Sprintf("%d-%d-d", lev.RelX, lev.RelY)
		log.Println("find:", keyd, lev.Name)

		levv := lev // 注意，newNode 内部不能用 lev， 这样会指针指向混乱， 所以 额外用了一个局部变量。
		lv.levelNodes[i] = newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
			eng.SetSubTex(n, levelTexs[keyd])
			eng.SetTransform(n, levv.Rect.ToF32Affine())
		})
		i++

	}

}

// Update 每次绘图前，逻辑相关的操作。
func (lv *ListView) Update(now clock.Time) {

}

// Press 当 touch 事件发生时， 判断是按在那个游戏精灵元素上，以及对应的处理策略分支。
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

			lv.NextPage()

			return
		} else if lv.model.BtnPrePage.Status == button.BtnPress {
			lv.model.BtnPrePage.Status = button.BtnNormal
			log.Println("BtnPrePage 释放按下状态")
			// 前一页按钮的操作逻辑

			lv.PrePage()

			return
		}

	}

	var currSelectLevel *level.LevelInfo
	currSelectLevel = nil
	// 判断是那个关卡被点击
	if touchEvent.Type == touch.TypeBegin {
		for _, level := range lv.model.Arr {
			if gp.In(level.Rect) {
				currSelectLevel = level
				break
			}
		}
	}

	// 关卡被点击的处理逻辑
	if currSelectLevel != nil {
		log.Println(currSelectLevel)
		LoadGameView(lv.model.GetCurrPageNum(), currSelectLevel)
	}
}

// OnScreenSizeChange 屏幕分辨率发生变化时，相关操作
func (lv *ListView) OnScreenSizeChange(currSZ size.Event, displayMultiple float32) {
	// 更新数据实体类的数据
	lv.model.OnScreenSizeChange(currSZ, displayMultiple)
}

// NextPage 下一页
func (lv *ListView) NextPage() {
	lv.model.NextPage()    // 数据变化
	lv.refreshView(lv.eng) // 界面变化
}

// PrePage 上一页
func (lv *ListView) PrePage() {
	lv.model.PrePage()     // 数据变化
	lv.refreshView(lv.eng) // 界面变化
}

// ClearRootViewChildNodes 清理掉 RootViewNode 的所有子节点
func (lv *ListView) ClearRootViewChildNodes() {
	// 前一页
	if lv.btnPre.Parent != nil {
		lv.RootViewNode.RemoveChild(lv.btnPre)
	}

	// 下一页按钮
	if lv.btnNext.Parent != nil {
		lv.RootViewNode.RemoveChild(lv.btnNext)
	}

	// 清除之前的数据
	if len(lv.levelNodes) > 0 {
		for _, nn := range lv.levelNodes {
			if nn.Parent != nil {
				lv.RootViewNode.RemoveChild(nn)
			}
		}
	}
}
