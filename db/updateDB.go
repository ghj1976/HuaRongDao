// 所有跟数据库有关的更新操作封装。
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghj1976/HuaRongDao/level"

	_ "github.com/mattn/go-sqlite3"
)

// 更新一个关卡信息到DB
// li 需要更新的数据
// dir 数据库文件的相对目录，默认是 assets
// isUpdate 如果数据重复时，是否要替换成新的？
func UpdateToDB(li *level.LevelInfo, dir string, isUpdate bool) {
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

	stmt, err := db.Prepare("insert into level(ID,Name,Layout,MinStepNum,Class,CaoPos,HNum,SuccessStepRecord,MyStepNum) values(?,?,?,?,?,?,?,?,?) ")
	if err != nil {
		log.Println("db insert Prepare Err:", err)
		os.Exit(-1)
		return
	}

	needUpdate := false // 需要做更新操作

	_, err = stmt.Exec(li.ID, li.Name, li.Layout, li.MinStepNum, li.Class, li.CaoPos, li.HNum, li.StepRecord, li.StepNum)
	if err != nil {
		log.Printf("%d  %s", li.ID, err.Error())
		if err.Error() == "UNIQUE constraint failed: game.ID" {
			// 数据库中已经有这条记录了。
			if isUpdate {
				needUpdate = true
			}
		}
	}

	if needUpdate {
		// 做数据库更新操作。更新操作的逻辑，至少不是默认值，就做update，以最新的为准。
		sql := "update level set "
		sql += fmt.Sprintf(" MinStepNum = %d ", li.MinStepNum)
		sql += fmt.Sprintf(" , CaoPos = %d ", li.CaoPos)
		sql += fmt.Sprintf(" , HNum = %d ", li.HNum)
		sql += fmt.Sprintf(" , MyStepNum = %d ", li.StepNum)

		if len(strings.TrimSpace(li.Name)) > 0 {
			sql += fmt.Sprintf(", Name = \"%s\" ", strings.TrimSpace(li.Name))
		}
		if len(strings.TrimSpace(li.Layout)) > 0 {
			sql += fmt.Sprintf(", Layout = \"%s\" ", strings.TrimSpace(li.Layout))
		}
		if len(strings.TrimSpace(li.Class)) > 0 {
			sql += fmt.Sprintf(", Class = \"%s\" ", strings.TrimSpace(li.Class))
		}
		if len(strings.TrimSpace(li.StepRecord)) > 0 {
			sql += fmt.Sprintf(", SuccessStepRecord = \"%s\" ", strings.TrimSpace(li.StepRecord))
		}
		sql += fmt.Sprintf(" where ID = %d ", li.ID)

		_, err = db.Exec(sql)
		if err != nil {
			log.Println("db update Err:", err)
		} else {
			log.Println("shop update db: ", li.ID)
		}
	}
}

// 更新一组关卡信息到DB，中间发生异常时，跳过继续
func UpdateArrToDB(liArr []level.LevelInfo, dir string, isUpdate bool) {
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

	stmt, err := db.Prepare("insert into level(ID,Name,Layout,MinStepNum,Class,CaoPos,HNum,SuccessStepRecord,MyStepNum) values(?,?,?,?,?,?,?,?,?) ")
	if err != nil {
		log.Println("db insert Prepare Err:", err)
		os.Exit(-1)
		return
	}

	for _, li := range liArr {
		needUpdate := false // 需要做更新操作

		_, err = stmt.Exec(li.ID, li.Name, li.Layout, li.MinStepNum, li.Class, li.CaoPos, li.HNum, li.StepRecord, li.StepNum)
		if err != nil {
			log.Printf("%d  %s", li.ID, err.Error())
			if err.Error() == "UNIQUE constraint failed: game.ID" {
				// 数据库中已经有这条记录了。
				if isUpdate {
					needUpdate = true
				}
			}
		}

		if needUpdate {
			// 做数据库更新操作。更新操作的逻辑，至少不是默认值，就做update，以最新的为准。
			sql := "update level set "
			sql += fmt.Sprintf(" MinStepNum = %d ", li.MinStepNum)
			sql += fmt.Sprintf(" , CaoPos = %d ", li.CaoPos)
			sql += fmt.Sprintf(" , HNum = %d ", li.HNum)
			sql += fmt.Sprintf(" , MyStepNum = %d ", li.StepNum)

			if len(strings.TrimSpace(li.Name)) > 0 {
				sql += fmt.Sprintf(", Name = \"%s\" ", strings.TrimSpace(li.Name))
			}
			if len(strings.TrimSpace(li.Layout)) > 0 {
				sql += fmt.Sprintf(", Layout = \"%s\" ", strings.TrimSpace(li.Layout))
			}
			if len(strings.TrimSpace(li.Class)) > 0 {
				sql += fmt.Sprintf(", Class = \"%s\" ", strings.TrimSpace(li.Class))
			}
			if len(strings.TrimSpace(li.StepRecord)) > 0 {
				sql += fmt.Sprintf(", SuccessStepRecord = \"%s\" ", strings.TrimSpace(li.StepRecord))
			}
			sql += fmt.Sprintf(" where ID = %d ", li.ID)

			_, err = db.Exec(sql)
			if err != nil {
				log.Println("db update Err:", err)
			} else {
				log.Println("shop update db: ", li.ID)
			}
		}
	}
}
