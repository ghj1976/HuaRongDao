// 所有跟数据库有关的读操作封装。
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ghj1976/HuaRongDao/level"

	_ "github.com/mattn/go-sqlite3"
)

// 分页读取数据
// pageNum 第几页， 第一页为1
// pageSize 每页几条
// dir 数据库所在的相对目录， 默认是 assets 子目录
// arr 返回的数据
// hasPrePage 是否有前一页
// hasNextPage 是否有后一页
func ReadPage(pageNum, pageSize int, dir string) (arr []level.LevelInfo, hasPrePage, hasNextPage bool) {
	arr = []level.LevelInfo{}
	if pageNum <= 1 {
		hasPrePage = false
	} else {
		hasPrePage = true
	}
	hasNextPage = false

	if len(dir) <= 0 {
		dir = "assets"
	}
	dbFileName := filepath.Join(dir, "game.db")
	log.Println("db:", dbFileName)

	// 打开DB文件
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Println("Open DB File Err:", err)
		os.Exit(-1)
		return
	}
	defer db.Close()

	// 后都多取一个，这样可以判断后面是否还有。
	ps := pageSize + 1
	offset := (pageNum - 1) * pageSize
	sql := fmt.Sprintf("select ID,Name,Layout,MinStepNum,Class,CaoPos,HNum,SuccessStepRecord,MyStepNum from level order by ID limit %d offset %d ", ps, offset)
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("db select err:", err)
		os.Exit(-1)
	}

	var i int
	i = 0
	for rows.Next() {
		i++
		if i > pageSize { // 能多取一个，代表有下一页。
			hasNextPage = true
			break
		}
		le := level.LevelInfo{}
		err = rows.Scan(&le.ID, &le.Name, &le.Layout,
			&le.MinStepNum, &le.Class, &le.CaoPos,
			&le.HNum, &le.StepRecord, &le.StepNum)
		if err != nil {
			log.Println("db select Scan err:", err)
			os.Exit(-1)
		}
		arr = append(arr, le)
	}
	return
}
