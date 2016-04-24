// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux

// Flappy Gopher is a simple one-button game that uses the
// mobile framework and the experimental sprite engine.
package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app.Main(func(a app.App) {
		var glctx gl.Context
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					glctx = nil
				}
			case size.Event:
				sz = e
				//				log.Println("aaa", sz)
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				onPaint(glctx, sz)
				a.Publish()
				a.Send(paint.Event{}) // keep animating
			case touch.Event:
				te = e
				game.Press(te)
				//				if down := e.Type == touch.TypeBegin; down || e.Type == touch.TypeEnd {
				//					game.Press(down)
				//				}
				// log.Println(te.Type, te.X, te.Y)

			}
		}
	})
}

var (
	sz        size.Event  // 当前屏幕尺寸信息
	te        touch.Event // 当前的触屏事件信息
	startTime = time.Now()
	images    *glutil.Images
	eng       sprite.Engine
	scene     *sprite.Node
	game      *Game
)

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	game = NewGame()

	scene = game.InitScene(eng, sz)
}

func onStop() {
	eng.Release()
	images.Release()
	game.stop()
	game = nil
	images = nil
	eng = nil
	log.Println("onStop.")
	os.Exit(-1)
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(171.0/255.0, 190.0/255.0, 62.0/255.0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	game.Update(now)           // 游戏逻辑相关更新操作
	eng.Render(scene, now, sz) // 只管绘图，不管游戏逻辑
}

func NewGame() *Game {
	var g Game
	g.reset()
	// 关卡信息
	layout := `
				黄关关赵
				黄甲乙赵
				张曹曹马
				张曹曹马
				丙一一丁`
	g.Level = InitLevel("横刀立马", layout, 0) // 不涉及具体绘图数据的计算，只做业务数据的计算初始化
	return &g
}
