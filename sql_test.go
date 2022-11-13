package main

import (
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"testing"
)

var (
	tableMysql   *dialect.TableModel
	mappingMysql = dialect.NewMappingMysql()
	dialectMysql dialect.Dialect
)

func init() {
	var err error
	dialectMysql, err = dialect.NewMappingDialect(mappingMysql)
	if err != nil {
		panic(err)
	}
	tableMysql = &dialect.TableModel{}
	tableMysql.TableName = "TEST_TB_1"
	tableMysql.ColumnList = mappingMysql.GenColumns()
	for _, column := range tableMysql.ColumnList {
		if len(tableMysql.PrimaryKeys) > 3 {
			continue
		}
		tableMysql.PrimaryKeys = append(tableMysql.PrimaryKeys, column.ColumnName)
	}
	for i, column := range tableMysql.ColumnList {
		if len(tableMysql.IndexList) > 3 {
			continue
		}
		indexType := "INDEX"
		if i%2 == 0 {
			indexType = "UNIQUE"
		}
		tableMysql.AddIndex(&dialect.IndexModel{
			ColumnName: column.ColumnName,
			IndexType:  indexType,
		})
	}
}

func TestMysqlTable(t *testing.T) {
	sqlList, err := dialectMysql.TableCreateSql(nil, "", tableMysql)
	if err != nil {
		panic(err)
	}
	for _, sqlOne := range sqlList {
		fmt.Println(sqlOne)
	}
}
