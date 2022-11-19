package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func DoExec(db *sql.DB, sqlInfo string, args ...interface{}) (result sql.Result, err error) {
	if len(sqlInfo) == 0 {
		return
	}
	resultList, _, _, err := DoExecs(db, []string{sqlInfo}, []interface{}{args})
	if err != nil {
		return
	}
	if len(resultList) > 0 {
		result = resultList[0]
	}
	return
}

func DoExecs(db *sql.DB, sqlList []string, argsList ...[]interface{}) (resultList []sql.Result, errSql string, errArgs []interface{}, err error) {
	sqlListSize := len(sqlList)
	if sqlListSize == 0 {
		return
	}
	if len(argsList) == 0 {
		argsList = make([][]interface{}, sqlListSize)
	}
	argsListSize := len(argsList)
	if sqlListSize != argsListSize {
		err = errors.New(fmt.Sprintf("sqlList size is [%d] but argsList size is [%d]", sqlListSize, argsListSize))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	var result sql.Result
	for i := 0; i < sqlListSize; i++ {
		sqlInfo := sqlList[i]
		args := argsList[i]
		if strings.TrimSpace(sqlInfo) == "" {
			continue
		}
		result, err = tx.Exec(sqlInfo, args...)
		if err != nil {
			errSql = sqlInfo
			errArgs = args
			return
		}
		resultList = append(resultList, result)
	}

	return
}

func DoQuery(db *sql.DB, sqlInfo string, args ...interface{}) (list []map[string]interface{}, err error) {
	_, list, err = DoQueryWithColumnTypes(db, sqlInfo, args)
	if err != nil {
		return
	}
	return
}

func DoQueryWithColumnTypes(db *sql.DB, sqlInfo string, args ...interface{}) (columnTypes []*sql.ColumnType, list []map[string]interface{}, err error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	columnTypes, err = rows.ColumnTypes()
	if err != nil {
		return
	}
	cache := GetSqlValueCache(columnTypes) //临时存储每行数据
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]interface{})
		for index, data := range cache {
			item[columnTypes[index].Name()] = GetSqlValue(columnTypes[index], data)
		}
		list = append(list, item)
	}

	return
}
