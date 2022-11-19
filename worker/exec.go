package worker

import (
	"database/sql"
	"errors"
	"fmt"
)

func DoExec(db *sql.DB, sqlInfo string, args ...interface{}) (errSql string, err error) {
	if len(sqlInfo) == 0 {
		return
	}
	errSql, err = DoExecs(db, []string{sqlInfo}, []interface{}{args})

	return
}

func DoExecs(db *sql.DB, sqlList []string, argsList ...[]interface{}) (errSql string, err error) {
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

	var lastSql string
	defer func() {
		if err != nil {
			errSql = lastSql
		} else {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprint(e))
				errSql = lastSql
			}
		}
	}()
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
	for _, one := range sqlList {
		if one == "" {
			continue
		}
		lastSql = one
		_, err = tx.Exec(one)
		if err != nil {
			return
		}
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
