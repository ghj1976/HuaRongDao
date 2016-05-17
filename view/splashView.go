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
	var width, height float32 = 30.0, 10.0
	return f32.Affine{
		{(float32(sz.WidthPt) - width) / 2.0, 0, width},
		{0, (float32(sz.HeightPt) - height) / 3.0, height},
	}

}

// 如果没加载好，则加载好 再返回显示节点。
// 如果已经加载好了， 直接返回显示节点
func LoadSplashView(eng sprite.Engine) *sprite.Node {
	if !splashViewFinishInit {
		texs := textures.LoadTexturesSplash(eng)

		splashNode = common.NewNodeNoShow(eng, func(e sprite.Engine, n *sprite.Node, t clock.Time) {
			e.SetSubTex(n, texs[textures.TexSplash])
			e.SetTransform(n, getScreenSplashCenterPoint(model.ScreenSize))
		})
		splashViewFinishInit = true
	}
	return splashNode

}
