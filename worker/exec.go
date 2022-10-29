package worker

import (
	"database/sql"
	"errors"
	"fmt"
)

func DoExec(db *sql.DB, sqlList []string) (errSql string, err error) {
	if len(sqlList) == 0 {
		return
	}

	var lastSql string
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			errSql = lastSql
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
		_, err = db.Exec(one)
		if err != nil {
			return
		}
	}

	return
}

func DoQuery(db *sql.DB, sqlInfo string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()
	cache := GetSqlValueCache(columnTypes) //临时存储每行数据
	var list []map[string]interface{}      //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for index, data := range cache {
			item[columns[index]] = GetSqlValue(columnTypes[index], data)
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil
}
