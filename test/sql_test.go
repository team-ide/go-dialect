package test

import (
	"os"
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

//func TestSqlParse(t *testing.T) {
//	var err error
//	var convertParser *worker.ConvertParser
//
//	srcSql := loadSql(`temp/sql_test.sql`)
//
//	convertParser = worker.NewConvertParser(srcSql, dialect.Mysql)
//	err = convertParser.Parse()
//	if err != nil {
//		panic(err)
//	}
//	saveSql(convertParser.GetDestSql(), "temp/sql_mysql.sql")
//
//	convertParser = worker.NewConvertParser(srcSql, dialect.Oracle)
//	err = convertParser.Parse()
//	if err != nil {
//		panic(err)
//	}
//	saveSql(convertParser.GetDestSql(), "temp/sql_oracle.sql")
//
//	convertParser = worker.NewConvertParser(srcSql, dialect.ShenTong)
//	err = convertParser.Parse()
//	if err != nil {
//		panic(err)
//	}
//	saveSql(convertParser.GetDestSql(), "temp/sql_shentong.sql")
//
//	convertParser = worker.NewConvertParser(srcSql, dialect.KinBase)
//	err = convertParser.Parse()
//	if err != nil {
//		panic(err)
//	}
//	saveSql(convertParser.GetDestSql(), "temp/sql_kinbase.sql")
//
//	convertParser = worker.NewConvertParser(srcSql, dialect.DM)
//	err = convertParser.Parse()
//	if err != nil {
//		panic(err)
//	}
//	saveSql(convertParser.GetDestSql(), "temp/sql_dm.sql")
//
//	convertParser = worker.NewConvertParser(srcSql, dialect.Sqlite)
//	err = convertParser.Parse()
//	if err != nil {
//		panic(err)
//	}
//	saveSql(convertParser.GetDestSql(), "temp/sql_sqlite.sql")
//
//}
