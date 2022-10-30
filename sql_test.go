package go_dialect

import (
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"os"
	"testing"
)

func loadSql(name string) (srcSql string) {
	bs, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	srcSql = string(bs)
	return
}

func saveSql(destSql string, name string) {
	err := os.WriteFile(name, []byte(destSql), 0777)
	if err != nil {
		panic(err)
	}
	return
}

func TestSqlParse(t *testing.T) {
	var err error
	var convertParser *worker.ConvertParser

	srcSql := loadSql(`temp/sql_test.sql`)

	convertParser = worker.NewConvertParser(srcSql, dialect.Mysql)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_mysql.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.Oracle)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_oracle.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.ShenTong)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_shentong.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.KinBase)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_kinbase.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.DaMen)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_damen.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.Sqlite)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_sqlite.sql")

}

func TestSqlSplit(t *testing.T) {

	var sqlInfo = `
select * from aa where 1=1 and a='''x;x'';'
;
select * from aa where 1=1

`
	sqlList := worker.SplitSqlList(sqlInfo)
	for _, sqlOne := range sqlList {
		fmt.Println("-------sql one start--------")
		fmt.Println(sqlOne)
		fmt.Println("-------sql one end--------")
	}

}
