package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
)

func DatabasesSelect(db *sql.DB, dia dialect.Dialect) (list []*dialect.DatabaseModel, err error) {
	sqlInfo, err := dia.DatabasesSelectSql()
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo)
	if err != nil {
		return
	}
	for _, data := range dataList {
		model, e := dia.DatabaseModel(data)
		if e != nil {
			model = &dialect.DatabaseModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
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
