package worker

import "database/sql"

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
