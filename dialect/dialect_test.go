package dialect

import (
	"fmt"
	"testing"
)

func TestInsertDataList(t *testing.T) {
	dia, err := NewDialect("oracle")
	if err != nil {
		return
	}
	param := &ParamModel{}
	var columnList = []*ColumnModel{
		{ColumnName: "name"},
		{ColumnName: "age"},
		{ColumnName: "account"},
	}
	var dataList = []map[string]interface{}{
		{"name": "名称1", "age": 1},
		{"name": "名称2", "age": 11},
		{"name": "名称3", "age": 11, "account": "name3"},
	}
	sqlList, valuesList, batchSqlList, batchValuesList, err := dia.DataListInsertSql(param, "TEST_DB", "USER_INFO", columnList, dataList)
	fmt.Println("--------sql list--------")
	for index := range sqlList {
		fmt.Println("sql:", sqlList[index])
		fmt.Println("values:", valuesList[index])
	}
	fmt.Println("--------batch sql list--------")
	for index := range batchSqlList {
		fmt.Println("batchSq:", batchSqlList[index])
		fmt.Println("batchValues:", batchValuesList[index])
	}
}
