package main

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	"testing"
)

func TestSqlParseGen(t *testing.T) {
	err := dataTypeParse(`数据库SQL.xlsx`, "dialect/mapping.sql.go")
	if err != nil {
		panic(err)
	}
}

type sqlDatabaseModel struct {
	Name    string
	sqlList []*sqlDatabaseSqlModel
}

type sqlDatabaseSqlModel struct {
	Name    string
	Sql     string
	Comment string
}

func sqlParse(path string, outPath string) (err error) {
	xlsxFForRead, err := xlsx.OpenFile(path)
	if err != nil {
		err = errors.New("excel [" + path + "] open error, " + err.Error())
		return
	}
	sheets := xlsxFForRead.Sheets

	var databases []*sqlDatabaseModel

	for _, sheet := range sheets {
		database := &sqlDatabaseModel{}
		database.Name = sheet.Name

		var titles []string

		var RowMergeEnd = -1
		var RowMergeCell = -1
		var RowMergeValue string
		for rowIndex, row := range sheet.Rows {

			if rowIndex == 0 {
				for _, cell := range row.Cells {
					title := cell.Value
					title = strings.TrimSpace(title)
					titles = append(titles, title)
				}
				continue
			}
			var dataType = map[string]string{}
			for cellIndex, cell := range row.Cells {
				if cellIndex >= len(titles) {
					break
				}
				title := titles[cellIndex]
				if title == "" {
					continue
				}
				value := cell.Value
				value = strings.TrimSpace(value)
				if cell.VMerge > 0 {
					RowMergeCell = cellIndex
					RowMergeEnd = rowIndex + cell.VMerge
					RowMergeValue = value
				}
				if RowMergeCell == cellIndex {
					if rowIndex <= RowMergeEnd {
						value = RowMergeValue
					} else {
						RowMergeEnd = -1
						RowMergeValue = ""
					}
				}
				dataType[title] = value
			}
			if dataType["名称"] == "" {
				continue
			}

			sqlDatabaseSql := &sqlDatabaseSqlModel{}
			sqlDatabaseSql.Name = dataType["名称"]
			sqlDatabaseSql.Sql = dataType["SQL"]
			sqlDatabaseSql.Comment = dataType["说明"]
			database.sqlList = append(database.sqlList, sqlDatabaseSql)
		}

		databases = append(databases, database)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return
	}
	_, err = outFile.WriteString(`package dialect

import "strings"

`)
	if err != nil {
		return
	}
	for _, one := range databases {
		fmt.Println("-------- database [" + one.Name + "] start --------")

		var code string
		code += "// " + one.Name + " 数据库 SQL" + "\n"
		funcName := ""
		if strings.EqualFold(one.Name, "Mysql") {
			funcName = "appendMysqlSql"
		} else if strings.EqualFold(one.Name, "Oracle") {
			funcName = "appendOracleSql"
		} else if strings.EqualFold(one.Name, "达梦") {
			funcName = "appendDmSql"
		} else if strings.EqualFold(one.Name, "金仓") {
			funcName = "appendKingBaseSql"
		} else if strings.EqualFold(one.Name, "神通") {
			funcName = "appendShenTongSql"
		} else if strings.EqualFold(one.Name, "Sqlite") {
			funcName = "appendSqliteSql"
		} else if strings.EqualFold(one.Name, "GBase") {
			funcName = "appendGBaseSql"
		} else if strings.EqualFold(one.Name, "Postgresql") {
			funcName = "appendPostgresqlSql"
		} else if strings.EqualFold(one.Name, "DB2") {
			funcName = "appendDb2Sql"
		}
		code += "func " + funcName + "Init(){" + "\n"
		//for _, sqlModel := range one.sqlList {
		//
		//}
		code += "}" + "\n\n"
		fmt.Println(code)
		_, err = outFile.WriteString(code)
		if err != nil {
			return
		}
		fmt.Println("-------- database [" + one.Name + "] end --------")
	}
	return
}
