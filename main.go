// +build darwin linux
// 华容道
package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/ghj1976/HuaRongDao/model"
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

	"runtime/pprof"
)

var (
	sz        size.Event // 当前屏幕尺寸信息
	startTime = time.Now()
	images    *glutil.Images
	eng       sprite.Engine
	//scene     *sprite.Node
	//game      *Game

	OpenProf = flag.Bool("prof", false, "是否启用性能跟踪，默认不启用。")
	f        *os.File // 性能跟踪写的文件
)

func main() {
	flag.Parse()

	// 性能监控部分开始代码 	log.Println(*OpenProf)
	if *OpenProf {
		f, err := os.OpenFile("./tmp/cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		// 注意， 这里的 defer 不会被执行到， 所以结束pprof的方法在另外的地方。
	}

	rand.Seed(time.Now().UnixNano())

	app.Main(func(a app.App) {

		var glctx gl.Context
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageAlive) {
				case lifecycle.CrossOn:
					log.Println("onCreate")
					onCreate()
				case lifecycle.CrossOff:
					log.Println("onDestroy")
					glctx = nil
					onDestroy()
				}
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					log.Println("onStart")
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					log.Println("onStop")
					onStop()
					if runtime.GOOS != "android" && runtime.GOOS != "ios" {
						glctx = nil
						onDestroy()
						os.Exit(-1) // 桌面版本，直接退出,跳到onDestroy。
					}
				}

			case size.Event:
				model.InitScreenSize(e)
				sz = e
				//				log.Println("屏幕：", sz)
				//				if game != nil {
				//					game.InitGameElementLength(sz)
				//				}
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				onPaint(glctx, sz)
				a.Publish()
				a.Send(paint.Event{}) // keep animating

			case touch.Event:
				Press(e)
			}
		}
		log.Println("end 123")
	})
}

func onCreate() {
	//game = NewGame()
}

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	log.Println("112")
	Init(eng)

	//game.InitGameElementLength(sz)
	//scene = game.InitScene(eng, sz)
}

func onStop() {

	//					log.Println(*OpenProf)
	// 性能跟踪
	if *OpenProf {
		// 手工退出时，写内存
		fm, err := os.OpenFile("./tmp/mem.out", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("开始写 HeapProfile")
		pprof.WriteHeapProfile(fm)
		log.Println("写 HeapProfile 完成。")
		fm.Close()

		// 结束性能跟踪，之前的 defer 不会被触发。
		pprof.StopCPUProfile()
		f.Close()
	}

}

func onDestroy() {
	eng.Release()
	images.Release()
	//game.stop()
	//game = nil
	images = nil
	eng = nil
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(171.0/255.0, 190.0/255.0, 62.0/255.0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	//game.Update(now)
	Update(now)                    // 游戏逻辑相关更新操作
	eng.Render(gameScene, now, sz) // 只管绘图，不管游戏逻辑
}
