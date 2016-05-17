package textures

import (
	"image"

	_ "image/png"
	"log"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

const (
	texLoading01 = iota // loading 动画第一帧
	texLoading02
	texLoading03
	texLoading04
	texLoading05
	texLoading06
	texLoading07
	texLoading08
)

var (
	loadingFrames = []int{texLoading01, texLoading02, texLoading03, texLoading04, texLoading05, texLoading06, texLoading07, texLoading08}
)

// 加载 Splash 视图需要的纹理。
func LoadTexturesLoading(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("loading.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	m, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	return []sprite.SubTex{
		texLoading01: sprite.SubTex{t, image.Rect(0, 0, 768, 288)},
		texLoading02: sprite.SubTex{t, image.Rect(770, 0, 1538, 288)},
		texLoading03: sprite.SubTex{t, image.Rect(0, 300, 768, 588)},
		texLoading04: sprite.SubTex{t, image.Rect(770, 300, 1538, 588)},
		texLoading05: sprite.SubTex{t, image.Rect(0, 600, 768, 888)},
		texLoading06: sprite.SubTex{t, image.Rect(770, 600, 1538, 888)},
		texLoading07: sprite.SubTex{t, image.Rect(0, 900, 768, 1188)},
		texLoading08: sprite.SubTex{t, image.Rect(770, 900, 1538, 1188)},
	}
}

// 显示那副图
func LoadingFrame(t, d clock.Time) int {
	total := int(d) * len(loadingFrames)
	return loadingFrames[(int(t)%total)/int(d)]
}
