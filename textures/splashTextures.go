package textures

import (
	"image"

	_ "image/png"
	"log"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite"
)

const (
	TexSplash = iota // 启动页面的图片， 没有复杂逻辑，所以直接暴露
)

// 加载 Splash 视图需要的纹理。
func LoadTexturesSplash(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("game.png")
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
		TexSplash: sprite.SubTex{t, image.Rect(0, 0, 180, 60)}}
}
