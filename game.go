package main

import (
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"image"
	"log"

	_ "image/jpeg"
)

type Game struct {
	lastCalc clock.Time // when we last calculated a frame
}

func (g *Game) InitScene(eng sprite.Engine, sz size.Event) *sprite.Node {

	back := loadTextures(eng)
	b2 := loadTexturesB2(eng)

	scene := &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	newNode := func(fn arrangerFunc) {
		n := &sprite.Node{Arranger: arrangerFunc(fn)}
		eng.Register(n)
		scene.AppendChild(n)
	}

	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, back)
		eng.SetTransform(n, f32.Affine{
			{float32(sz.WidthPt), 0, 0},
			{0, float32(sz.HeightPx), 0},
		})
	})

	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		eng.SetSubTex(n, b2)
		eng.SetTransform(n, f32.Affine{
			{float32(sz.WidthPt), 0, 0},
			{0, 50, 0},
		})
	})

	return scene

}

// 加载纹理图
func loadTextures(eng sprite.Engine) sprite.SubTex {
	a, err := asset.Open("odd.jpg")
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

	return sprite.SubTex{t, image.Rect(0, 0, 128, 128)}

}

// 加载纹理图
func loadTexturesB2(eng sprite.Engine) sprite.SubTex {
	a, err := asset.Open("b2.jpeg")
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

	return sprite.SubTex{t, image.Rect(0, 0, 267, 189)}

}

func (g *Game) reset() {
}

func (g *Game) Press(down bool) {
}

func (g *Game) Update(now clock.Time) {
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
