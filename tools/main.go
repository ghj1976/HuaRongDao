package main

import (
	"flag"
	"log"

	"github.com/ghj1976/HuaRongDao/db"
	"github.com/ghj1976/HuaRongDao/level"
)

var (
	cmd = flag.String("c", "draw", "需要执行的命令，默认 draw 绘制效果图")
)

func main() {
	flag.Parse() // 读取命令参数

	// 待处理的数据
	levelArr := level.InitData()

	switch *cmd {
	case "draw":
		drawListDemo(levelArr, 10, 4, "001.png")
	case "db":
		db.UpdateArrToDB(levelArr, "../assets", true)
	default:
		drawListDemo(levelArr, 10, 4, "001.png")

	}

	log.Println("finish")
}
