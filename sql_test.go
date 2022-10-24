package go_dialect

import (
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/parser"
	"os"
	"testing"
)

func loadSql() (srcSql string) {
	bs, err := os.ReadFile(`sql_test.sql`)
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
	var convertParser *parser.ConvertParser

	srcSql := loadSql()

	convertParser = parser.NewConvertParser(srcSql, dialect.Mysql, &dialect.GenerateParam{
		DatabasePackingCharacter: "`",
		TablePackingCharacter:    "`",
		ColumnPackingCharacter:   "`",
	})
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "sql_mysql.sql")

	convertParser = parser.NewConvertParser(srcSql, dialect.Oracle, &dialect.GenerateParam{
		DatabasePackingCharacter: "\"",
		TablePackingCharacter:    "\"",
		ColumnPackingCharacter:   "\"",
	})
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "sql_oracle.sql")

	convertParser = parser.NewConvertParser(srcSql, dialect.ShenTong, &dialect.GenerateParam{
		DatabasePackingCharacter: "\"",
		TablePackingCharacter:    "\"",
		ColumnPackingCharacter:   "\"",
	})
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "sql_shentong.sql")

	convertParser = parser.NewConvertParser(srcSql, dialect.KinBase, &dialect.GenerateParam{
		DatabasePackingCharacter: "\"",
		TablePackingCharacter:    "\"",
		ColumnPackingCharacter:   "\"",
	})
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "sql_kinbase.sql")

	convertParser = parser.NewConvertParser(srcSql, dialect.DaMen, &dialect.GenerateParam{
		DatabasePackingCharacter: "\"",
		TablePackingCharacter:    "\"",
		ColumnPackingCharacter:   "\"",
	})
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "sql_damen.sql")

	convertParser = parser.NewConvertParser(srcSql, dialect.Sqlite, &dialect.GenerateParam{
		DatabasePackingCharacter: "\"",
		TablePackingCharacter:    "\"",
		ColumnPackingCharacter:   "\"",
	})
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "sql_sqlite.sql")

}
