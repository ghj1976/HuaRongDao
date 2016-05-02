package level

type LevelBaseInfo struct {
	ID         int    // 编号
	Name       string // 名字
	Layout     string // 布局字符串
	MinStepNum int    // 最小步数
	Class      string // 分类
	Difficult  int    // 难度， 数字越大越难， 就是已知的最小步数
	CaoPos     int    // 曹操所在位置 1线 | 2线 | 3线 | 4线
	HNum       int    // 横棋子个数分类: 0横 | 1横 | 2横 | 3横 | 4横 | 5横
}
