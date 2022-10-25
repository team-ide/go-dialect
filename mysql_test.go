package go_dialect

import (
	"database/sql"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_mysql"
	"strings"
	"testing"
)

var (
	MysqlDb *sql.DB
)

func initMysql() {
	if MysqlDb != nil {
		return
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", "root", "123456", "localhost", 3306, "")
	var err error
	MysqlDb, err = sql.Open(db_mysql.GetDriverName(), connStr)
	if err != nil {
		panic(err)
	}
	return
}

func TestMysql(t *testing.T) {
	initMysql()
	databases(MysqlDb, dialect.Mysql)
}

func TestMysqlTableCreate(t *testing.T) {
	initMysql()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	database := &dialect.DatabaseModel{
		Name: "TEST_DB",
	}
	testDatabaseDelete(MysqlDb, dialect.Mysql, param, database.Name)
	testDatabaseCreate(MysqlDb, dialect.Mysql, param, database)
	testTableDelete(MysqlDb, dialect.Mysql, param, database.Name, getTable().Name)
	testTableCreate(MysqlDb, dialect.Mysql, param, database.Name, getTable())

	testColumnUpdate(MysqlDb, dialect.Mysql, param, database.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(MysqlDb, dialect.Mysql, param, database.Name, getTable().Name, "detail3")
	testColumnAdd(MysqlDb, dialect.Mysql, param, database.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(MysqlDb, dialect.Mysql, database.Name, getTable().Name)
}

func TestMysqlSql(t *testing.T) {
	initMysql()
	sqlInfo := loadSql("sql_mysql.sql")
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	database := &dialect.DatabaseModel{
		Name: "TEST_DB",
	}
	testDatabaseDelete(MysqlDb, dialect.Mysql, param, database.Name)
	testDatabaseCreate(MysqlDb, dialect.Mysql, param, database)
	sqlInfo = "use " + database.Name + ";\n" + sqlInfo

	sqlList := strings.Split(sqlInfo, ";\n")
	exec(MysqlDb, sqlList)
}
