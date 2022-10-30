package worker

import (
	"github.com/team-ide/go-dialect/dialect"
	"strings"
)

func InsertDataListSql(dia dialect.Dialect, ownerName string, tableName string, columnList []*dialect.ColumnModel, dataList []map[string]interface{}) (sqlList []string, batchSqlList []string, err error) {

	var batchSqlCache = make(map[string]string)
	var batchSqlIndexCache = make(map[string]int)
	var columnNames []string
	for _, one := range columnList {
		columnNames = append(columnNames, one.Name)
	}
	for _, data := range dataList {
		var columnList_ []string
		var values = "("
		for _, column := range columnList {
			str := dia.PackValue(column, data[column.Name])
			if strings.EqualFold(str, "null") {
				continue
			}
			columnList_ = append(columnList_, column.Name)
			values += str + ", "
		}
		values = strings.TrimSuffix(values, ", ")
		values += ")"

		insertSqlInfo := "INSERT INTO "
		if ownerName != "" {
			insertSqlInfo += dia.PackOwner(ownerName) + "."
		}
		insertSqlInfo += dia.PackTable(tableName)
		insertSqlInfo += " ("
		insertSqlInfo += dia.PackColumns(columnList_)
		insertSqlInfo += ") VALUES "

		sqlList = append(sqlList, insertSqlInfo+values)

		key := strings.Join(columnList_, ",")
		find, ok := batchSqlCache[key]
		if ok {
			find += ",\n" + values
			batchSqlCache[key] = find
			batchSqlList[batchSqlIndexCache[key]] = find
		} else {
			find = insertSqlInfo + "\n" + values
			batchSqlIndexCache[key] = len(batchSqlCache)
			batchSqlCache[key] = find
			batchSqlList = append(batchSqlList, find)
		}
	}

	return
}
