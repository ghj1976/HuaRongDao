package view

import (
	"github.com/ghj1976/HuaRongDao/common"
	"github.com/ghj1976/HuaRongDao/model"
	"github.com/ghj1976/HuaRongDao/textures"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

var (
	splashViewFinishInit bool         = false // 是否完成了初始化
	splashNode           *sprite.Node         // 外部要显示的视图节点
)

// 获得 Splash 的显示位置
func getScreenSplashCenterPoint(sz size.Event) f32.Affine {
	// 计算相对大小
	var height float32 = float32(sz.HeightPt) / 6.0
	var width float32 = float32(sz.WidthPt) * 2.0 / 3.0
	if width > 3.0*height {
		width = 3.0 * height
	} else {
		height = width / 3.0
	}
	return f32.Affine{
		{width, 0, (float32(sz.WidthPt) - width) / 2.0},
		{0, height, (float32(sz.HeightPt) - height) / 4.0},
	}

}

// 如果没加载好，则加载好 再返回显示节点。
// 如果已经加载好了， 直接返回显示节点
func LoadSplashView(eng sprite.Engine) *sprite.Node {
	if !splashViewFinishInit {
		texs := textures.LoadTexturesSplash(eng)

		splashNode = common.NewNodeNoShow(eng, func(e sprite.Engine, n *sprite.Node, t clock.Time) {
			e.SetSubTex(n, texs[textures.TexSplash])

			// 不管屏幕尺寸获得有无，都可以提前准备显示元素了
			sz, b := model.GetScreenSize()
			if b {
				//				log.Println("Splash Node 刷新 ok")
				e.SetTransform(n, getScreenSplashCenterPoint(sz))
			}

		})
		splashViewFinishInit = true
	}
	return splashNode

}
