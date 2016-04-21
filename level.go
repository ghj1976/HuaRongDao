package main

import (
	"log"
	"strings"
)

const (
	BlackChessManPos = '一' // 棋盘中空位的标示
)

type ChessManStatus byte // 棋子的状态枚举

// 每个关卡的信息类
type LevelInfo struct {
	Name       string             // 级别名称
	MinStepNum int                // 最小步数
	Layout     string             // 布局字符串
	MapArray   [5][4]rune         // 实时的当前地图数组
	ChessMans  map[rune]*ChessMan // 棋子集合
	StepRecord [2][]rune          // 游戏所走步数的记录， 第一个是棋子名，第二个是方向（上下左右）
}

// 游戏中的棋子类
type ChessMan struct {
	name            rune           // 棋子名称，唯一识别编号
	rect            GameRectangle  // 棋子所在位置（长方形）， 实际位置, pt 单位
	status          ChessManStatus // 棋子的状态，一共三种：可移动，不可移动，正在移动
	RelWidth        int            // 相对宽度，相对于小兵的棋子的宽度，小兵棋子宽为1.
	RelHeight       int            // 相对高度，相对于小兵的棋子的高度，小兵棋子高为1.
	RelLeftTopX     int            // 相对坐标，相对左上角的坐标位置 X 轴， 左上角最小为 0，0
	RelLeftTopY     int            // 相对坐标，相对左上角的坐标位置 Y 轴， 左上角最小为 0，0
	RelRightBottomX int            // 相对坐标，相对右下角的坐标位置 X 轴， 右下角最大为 3，4
	RelRightBottomY int            // 相对坐标，相对右下角的坐标位置 Y 轴， 右下角最大为 3，4
}

// 初始化一个关卡信息类
func InitLevel(name, layout string, minStepNum int) *LevelInfo {

	lv := &LevelInfo{}
	lv.Name = name
	lv.Layout = layout
	lv.MinStepNum = minStepNum
	// 布局信息转关卡棋子map
	lv.MapArray = layout2map(layout)
	// 把当前地图部署转化成棋子哈西map
	lv.ChessMans = chessManArray2Map(lv.MapArray)
	// 每个棋子是否可移动的判断
	lv.ComputeChessManStatus()

	// 布局校验检查代码
	// 只能有2个空格，4*5

	return lv
}

// 游戏是否成功的判断
func (lv *LevelInfo) Success() bool {
	// 每个棋子都有个当前相对位置，曹操棋子到了指定位置就成功了。
	// 然后是成功动画。

	return false
}

// 计算棋子的状态，是否可移动的计算
func (lv *LevelInfo) ComputeChessManStatus() {
	for name, cm := range lv.ChessMans {
		// 上移判断
		if lv.ChessManCanMoveUp(name) {
			cm.status = ChessManMovable
			continue
		}

		// 下移判断
		if lv.ChessManCanMoveDown(name) {
			cm.status = ChessManMovable
			continue
		}

		// 左移判断
		if lv.ChessManCanMoveLeft(name) {
			cm.status = ChessManMovable
			continue
		}

		// 右移判断
		if lv.ChessManCanMoveRight(name) {
			cm.status = ChessManMovable
			continue
		}

		cm.status = ChessManStable // 上面都无法命中
	}
}

// 指定棋子是否可 上移 判断
func (lv *LevelInfo) ChessManCanMoveUp(name rune) bool {
	cm, ok := lv.ChessMans[name]
	if !ok {
		return false
	}
	if cm.RelLeftTopY <= 0 {
		// 最上面无法再上移了，只有第二行的才可能上移。
		return false
	}
	// 上移判断
	b := true
	for i := 0; i < cm.RelWidth; i++ {
		b = b && (lv.MapArray[cm.RelLeftTopY-1][cm.RelLeftTopX+i] == BlackChessManPos)
	}
	return b
}

// 指定棋子是否可 下移 判断
func (lv *LevelInfo) ChessManCanMoveDown(name rune) bool {
	cm, ok := lv.ChessMans[name]
	if !ok {
		return false
	}
	if cm.RelRightBottomY >= 4 {
		// 最下面 是无法下移的，只有倒数第二行才有可能
		return false
	}
	// 下移判断
	b := true
	for i := 0; i < cm.RelWidth; i++ {
		b = b && (lv.MapArray[cm.RelRightBottomY+1][cm.RelRightBottomX+i] == BlackChessManPos)
	}
	return b
}

// 指定棋子是否可 左移 判断
func (lv *LevelInfo) ChessManCanMoveLeft(name rune) bool {
	cm, ok := lv.ChessMans[name]
	if !ok {
		return false
	}
	if cm.RelLeftTopX <= 0 {
		// 最左边是无法左移的。
		return false
	}
	// 左移判断
	b := true
	for i := 0; i < cm.RelHeight; i++ {
		b = b && (lv.MapArray[cm.RelRightBottomY+i][cm.RelRightBottomX-1] == BlackChessManPos)
	}
	return b
}

// 指定棋子是否可 右移 判断
func (lv *LevelInfo) ChessManCanMoveRight(name rune) bool {
	cm, ok := lv.ChessMans[name]
	if !ok {
		return false
	}
	if cm.RelLeftTopX >= 3 {
		// 最右边是无法右移的。
		return false
	}

	// 右移判断
	b := true
	for i := 0; i < cm.RelHeight; i++ {
		b = b && (lv.MapArray[cm.RelRightBottomY+i][cm.RelRightBottomX+1] == BlackChessManPos)
	}
	return b
}

// 游戏成功与否的判断
func (lv *LevelInfo) IsSuccess() bool {
	cmc, _ := lv.ChessMans['曹']
	if cmc.status == ChessManStable && cmc.RelLeftTopX == 1 && cmc.RelLeftTopY == 3 {
		return true
	} else {
		return false
	}
}

// 把布局 string 文件转换成二维数组
func layout2map(layout string) [5][4]rune {
	layoutArray := [5][4]rune{}
	posX, posY := 0, 0
	for _, c := range layout {
		// 忽略为显示而无意义的字符

		if !strings.ContainsRune("曹关赵马黄张甲乙丙丁一", c) {
			continue
		}

		layoutArray[posY][posX] = c

		// 计算下一个棋子的相对坐标位置
		posX++
		if posX >= 4 {
			posX = 0
			posY++
		}
		if posY >= 5 {
			break
		}
	}
	return layoutArray

}

// 把二维数组转换成哈西Map棋子集合,并计算棋子准确位置
func chessManArray2Map(arr [5][4]rune) map[rune]*ChessMan {
	cmap := make(map[rune]*ChessMan, 10)

	// 为了计算方便先把 layout 变成规范的 二维数组
	for x, y := 0, 0; ; x++ {
		if x >= 4 {
			x = 0
			y++
			if y >= 5 {
				break
			}
		}

		c := arr[y][x]
		// x，y 是棋子的相对位置
		if c == BlackChessManPos {
			continue // 空位不做处理，继续下一轮处理
		}

		// 遍历每个棋子，当发现一个没记录的棋子时，先向后，向下探索出这个棋子的边界，然后记录，然后继续遍历。
		if _, ok := cmap[c]; !ok {
			// map 中 没有才需要进行处理。
			cm := &ChessMan{name: c}
			// 左上角的位置，第一次被发现，一定是发现的左上角的点
			cm.RelLeftTopX = x
			cm.RelLeftTopY = y

			// 几个初始值
			cm.RelWidth = 1
			cm.RelHeight = 1
			cm.RelRightBottomX = x
			cm.RelRightBottomY = y

			if x < 3 && arr[y][x+1] == c {
				// 判断随后一个是同样的棋子
				cm.RelWidth = 2
				cm.RelRightBottomX = x + 1
			}
			if y < 4 && arr[y+1][x] == c {
				// 判断下面一个是同样的棋子
				cm.RelHeight = 2
				cm.RelRightBottomY = y + 1
			}
			cmap[c] = cm
		}
	}
	return cmap
}

// 计算棋子的具体准确位置，
// 注意，需要在页面可绘图后才能做这个运算
func (lv *LevelInfo) ComputeChessManRect() {
	log.Println(len(lv.ChessMans))
	log.Println(lv)
	for _, cm := range lv.ChessMans {
		// 计算棋子实际该出现的位置
		cm.rect = GameRectangle{}
		cm.rect.SetGameRectangle(
			GamePoint{
				X: GameChessManAreaX + float32(cm.RelLeftTopX)*ChessManWidth,
				Y: GameChessManAreaY + float32(cm.RelLeftTopY)*ChessManWidth},
			float32(cm.RelWidth)*ChessManWidth,
			float32(cm.RelHeight)*ChessManWidth)

		log.Println(cm)
		log.Println(cm.rect)
	}
}
