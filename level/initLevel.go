package level

func InitData() *[]LevelInfo {
	levelArr := []LevelInfo{
		LevelInfo{
			ID:         1,
			Name:       "横刀立马",
			Class:      "经典布局",
			Difficult:  81,
			MinStepNum: 81,
			CaoPos:     1,
			HNum:       1,
			Layout: `	赵曹曹马
						赵曹曹马
						黄关关张
						黄甲乙张
						丙一一丁
			`,
		},
		LevelInfo{
			ID:         2,
			Name:       "指挥若定",
			Class:      "经典布局",
			Difficult:  70,
			MinStepNum: 70,
			CaoPos:     1,
			HNum:       1,
			Layout: `	赵曹曹马
						赵曹曹马
						甲关关乙
						黄丙丁张
						黄一一张
			`,
		},
		LevelInfo{
			ID:         3,
			Name:       "将拥曹营",
			Class:      "经典布局",
			Difficult:  72,
			MinStepNum: 72,
			CaoPos:     1,
			HNum:       1,
			Layout: `	一曹曹一
						赵曹曹马
						赵黄张马
						甲黄张乙
						关关丙丁
			`,
		},
		LevelInfo{
			ID:         4,
			Name:       "齐头并进",
			Class:      "经典布局",
			Difficult:  60,
			MinStepNum: 60,
			CaoPos:     1,
			HNum:       1,
			Layout: `	赵曹曹马
						赵曹曹马
						甲乙丙丁
						黄关关张
						黄一一张
			`,
		},
		LevelInfo{
			ID:         5,
			Name:       "兵分三路",
			Class:      "经典布局",
			Difficult:  72,
			MinStepNum: 72,
			CaoPos:     1,
			HNum:       1,
			Layout: `	甲曹曹乙
						赵曹曹马
						赵关关马
						黄丙丁张
						黄一一张
			`,
		},
		LevelInfo{
			ID:         6,
			Name:       "雨声淅沥",
			Class:      "经典布局",
			Difficult:  47,
			MinStepNum: 47,
			CaoPos:     1,
			HNum:       1,
			Layout: `	赵曹曹甲
						赵曹曹乙
						黄关关马
						黄张一马
						丙张一丁
			`,
		},
		LevelInfo{
			ID:         7,
			Name:       "左右布兵",
			Class:      "经典布局",
			Difficult:  54,
			MinStepNum: 54,
			CaoPos:     1,
			HNum:       1,
			Layout: `	甲曹曹丙
						乙曹曹丁
						黄赵张马
						黄赵张马
						一关关一
			`,
		},
		LevelInfo{
			ID:         8,
			Name:       "桃花园中",
			Class:      "经典布局",
			Difficult:  70,
			MinStepNum: 70,
			CaoPos:     1,
			HNum:       1,
			Layout: `	甲曹曹丙
						黄曹曹马
						黄赵张马
						乙赵张丁
						一关关一
			`,
		},
		LevelInfo{
			ID:         9,
			Name:       "一路进军",
			Class:      "经典布局",
			Difficult:  58,
			MinStepNum: 58,
			CaoPos:     1,
			HNum:       1,
			Layout: `	黄曹曹甲
						黄曹曹乙
						赵张马丙
						赵张马丁
						一关关一
			`,
		},
		LevelInfo{
			ID:         10,
			Name:       "一路顺风",
			Class:      "经典布局",
			Difficult:  39,
			MinStepNum: 39,
			CaoPos:     1,
			HNum:       1,
			Layout: `	黄曹曹甲
						黄曹曹乙
						赵关关马
						赵丙张马
						一丁张一
			`,
		},
		LevelInfo{
			ID:         11,
			Name:       "围而不歼",
			Class:      "经典布局",
			Difficult:  62,
			MinStepNum: 62,
			CaoPos:     1,
			HNum:       1,
			Layout: `	黄曹曹甲
						黄曹曹乙
						赵关关丙
						赵马张丁
						一马张一
			`,
		},
		LevelInfo{
			ID:         12,
			Name:       "捷足先登",
			Class:      "经典布局",
			Difficult:  32,
			MinStepNum: 32,
			CaoPos:     1,
			HNum:       1,
			Layout: `	丙曹曹甲
						丁曹曹乙
						一关关一
						赵马张黄
						赵马张黄
			`,
		},
		LevelInfo{
			ID:         13,
			Name:       "插翅难飞",
			Class:      "经典布局",
			Difficult:  62,
			MinStepNum: 62,
			CaoPos:     1,
			HNum:       2,
			Layout: `	马曹曹甲
						马曹曹乙
						张张丙丁
						赵关关黄
						赵一一黄
			`,
		},
		LevelInfo{
			ID:         14,
			Name:       "守口如瓶一",
			Class:      "经典布局",
			Difficult:  81,
			MinStepNum: 81,
			CaoPos:     1,
			HNum:       2,
			Layout: `	马曹曹黄
						马曹曹黄
						甲赵一丙
						乙赵一丁
						张张关关
			`,
		},
		LevelInfo{
			ID:         15,
			Name:       "守口如瓶二",
			Class:      "经典布局",
			Difficult:  99,
			MinStepNum: 99,
			CaoPos:     1,
			HNum:       2,
			Layout: `	甲曹曹丙
						马曹曹黄
						马赵一黄
						乙赵一丁
						张张关关
			`,
		},
		LevelInfo{
			ID:         16,
			Name:       "双将挡路",
			Class:      "经典布局",
			Difficult:  73,
			MinStepNum: 73,
			CaoPos:     1,
			HNum:       2,
			Layout: `	马曹曹甲
						马曹曹乙
						赵张张黄
						赵关关黄
						丙一一丁
			`,
		},
		LevelInfo{
			ID:         17,
			Name:       "横马当关",
			Class:      "经典布局",
			Difficult:  83,
			MinStepNum: 83,
			CaoPos:     1,
			HNum:       2,
			Layout: `	马曹曹黄
						马曹曹黄
						张张关关
						甲赵一丙
						乙赵一丁
			`,
		},
		LevelInfo{
			ID:         18,
			Name:       "层层设防一",
			Class:      "经典布局",
			Difficult:  102,
			MinStepNum: 102,
			CaoPos:     1,
			HNum:       3,
			Layout: `	马曹曹黄
						马曹曹黄
						甲张张丙
						乙关关丁
						一赵赵一
			`,
		},
		LevelInfo{
			ID:         19,
			Name:       "层层设防二",
			Class:      "经典布局",
			Difficult:  120,
			MinStepNum: 120,
			CaoPos:     1,
			HNum:       3,
			Layout: `	甲曹曹丙
						马曹曹黄
						马张张黄
						乙关关丁
						一赵赵一
			`,
		},
		LevelInfo{
			ID:         20,
			Name:       "兵挡将阻",
			Class:      "经典布局",
			Difficult:  87,
			MinStepNum: 87,
			CaoPos:     1,
			HNum:       3,
			Layout: `	甲曹曹黄
						乙曹曹黄
						马张张丙
						马关关丁
						一赵赵一
			`,
		},
		LevelInfo{
			ID:         21,
			Name:       "堵塞要道",
			Class:      "经典布局",
			Difficult:  40,
			MinStepNum: 40,
			CaoPos:     1,
			HNum:       3,
			Layout: `	甲曹曹丙
						乙曹曹丁
						马黄张张
						马黄关关
						一赵赵一
			`,
		},
		LevelInfo{
			ID:         22,
			Name:       "瓮中之鳖",
			Class:      "经典布局",
			Difficult:  103,
			MinStepNum: 103,
			CaoPos:     1,
			HNum:       3,
			Layout: `	马曹曹黄
						马曹曹黄
						张张赵赵
						甲关关丙
						乙一一丁
			`,
		},
		LevelInfo{
			ID:         23,
			Name:       "层峦叠嶂",
			Class:      "经典布局",
			Difficult:  98,
			MinStepNum: 98,
			CaoPos:     1,
			HNum:       3,
			Layout: `	马曹曹黄
						马曹曹黄
						甲张张丙
						关关赵赵
						乙一一丁
			`,
		},
		LevelInfo{
			ID:         24,
			Name:       "水泄不通",
			Class:      "经典布局",
			Difficult:  79,
			MinStepNum: 79,
			CaoPos:     1,
			HNum:       4,
			Layout: `	马曹曹甲
						马曹曹乙
						张张黄黄
						关关赵赵
						丙一一丁
			`,
		},
		LevelInfo{
			ID:         25,
			Name:       "四路进兵",
			Class:      "经典布局",
			Difficult:  77,
			MinStepNum: 77,
			CaoPos:     1,
			HNum:       4,
			Layout: `	丙曹曹甲
						丁曹曹乙
						马一张张
						马一黄黄
						关关赵赵
			`,
		},
		LevelInfo{
			ID:         26,
			Name:       "入地无门",
			Class:      "经典布局",
			Difficult:  87,
			MinStepNum: 87,
			CaoPos:     1,
			HNum:       4,
			Layout: `	马曹曹甲
						马曹曹乙
						丙张张丁
						黄黄赵赵
						一关关一
			`,
		},
		LevelInfo{
			ID:         27,
			Name:       "勇闯五关",
			Class:      "经典布局",
			Difficult:  34,
			MinStepNum: 34,
			CaoPos:     1,
			HNum:       5,
			Layout: `	丙曹曹甲
						丁曹曹乙
						张张马马
						黄黄赵赵
						一关关一
			`,
		},

		LevelInfo{
			ID:         39,
			Name:       "四面楚歌",
			Class:      "经典布局",
			Difficult:  56,
			MinStepNum: 56,
			CaoPos:     2,
			HNum:       1,
			Layout: `	张甲乙马
						张曹曹马
						丁曹曹赵
						黄关关赵
						黄一一丙
			`,
		},
		LevelInfo{
			ID:         40,
			Name:       "前呼后拥",
			Class:      "经典布局",
			Difficult:  22,
			MinStepNum: 22,
			CaoPos:     1,
			HNum:       5,
			Layout: `	甲乙曹曹
						张张曹曹
						马马赵赵
						黄黄关关
						一一丙丁
			`,
		},
		LevelInfo{
			ID:         41,
			Name:       "兵临曹营",
			Class:      "经典布局",
			Difficult:  34,
			MinStepNum: 34,
			CaoPos:     1,
			HNum:       1,
			Layout: `	甲曹曹乙
						丙曹曹丁
						张关关马
						张黄赵马
						一黄赵一
			`,
		},
		LevelInfo{
			ID:         42,
			Name:       "五将逼宫",
			Class:      "经典布局",
			Difficult:  36,
			MinStepNum: 36,
			CaoPos:     2,
			HNum:       3,
			Layout: `	张张赵赵
						黄曹曹马
						黄曹曹马
						甲关关乙
						丙一一丁
			`,
		},

		LevelInfo{
			ID:         46,
			Name:       "前挡后阻",
			Class:      "经典布局",
			Difficult:  42,
			MinStepNum: 42,
			CaoPos:     1,
			HNum:       2,
			Layout: `	曹曹张张
						曹曹马甲
						黄赵马乙
						黄赵丙丁
						一关关一
			`,
		},
		LevelInfo{
			ID:         47,
			Name:       "近在咫尺",
			Class:      "经典布局",
			Difficult:  98,
			MinStepNum: 98,
			CaoPos:     4,
			HNum:       2,
			Layout: `	甲张马黄
						乙张马黄
						赵赵丙丁
						关关曹曹
						一一曹曹
			`,
		},
		LevelInfo{
			ID:         48,
			Name:       "走投无路",
			Class:      "经典布局",
			Difficult:  -1,
			MinStepNum: -1,
			CaoPos:     1,
			HNum:       0,
			Layout: `	张曹曹马
						张曹曹马
						赵关甲黄
						赵关乙黄
						丙一一丁
			`,
		},
		LevelInfo{
			ID:         49,
			Name:       "小燕出巢",
			Class:      "经典布局",
			Difficult:  103,
			MinStepNum: 103,
			CaoPos:     1,
			HNum:       3,
			Layout: `	张曹曹马
						张曹曹马
						赵赵黄黄
						甲关关乙
						丙一一丁
			`,
		},
		LevelInfo{
			ID:         50,
			Name:       "比翼横空",
			Class:      "经典布局",
			Difficult:  28,
			MinStepNum: 28,
			CaoPos:     1,
			HNum:       4,
			Layout: `	张张曹曹
						马马曹曹
						赵赵黄黄
						甲一乙关
						丙一丁关
			`,
		},
		LevelInfo{
			ID:         51,
			Name:       "夹道藏兵",
			Class:      "经典布局",
			Difficult:  75,
			MinStepNum: 75,
			CaoPos:     1,
			HNum:       4,
			Layout: `	曹曹甲张
						曹曹乙张
						赵赵黄黄
						马马关关
						丙一一丁
			`,
		},
		LevelInfo{
			ID:         52,
			Name:       "屯兵东路",
			Class:      "经典布局",
			Difficult:  71,
			MinStepNum: 71,
			CaoPos:     1,
			HNum:       1,
			Layout: `	曹曹张黄
						曹曹张黄
						关关甲乙
						赵马丙丁
						赵马一一
			`,
		},
		LevelInfo{
			ID:         53,
			Name:       "四将连关",
			Class:      "经典布局",
			Difficult:  39,
			MinStepNum: 39,
			CaoPos:     1,
			HNum:       3,
			Layout: `	曹曹关关
						曹曹张张
						赵马黄黄
						赵马甲乙
						丙一一丁
			`,
		},

		LevelInfo{
			ID:         55,
			Name:       "峰回路转",
			Class:      "经典布局",
			Difficult:  138,
			MinStepNum: 138,
			CaoPos:     2,
			HNum:       2,
			Layout: `	甲乙丙赵
						曹曹张赵
						曹曹张黄
						一关关黄
						一丁马马
						
			`,
		},
		// 结束位
	}

	return &levelArr
}
