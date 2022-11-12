package main

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_mysql"
	"testing"
)

var (
	MysqlDb      *sql.DB
	MysqlDialect dialect.Dialect
)

func initMysql() {
	if MysqlDb != nil {
		return
	}
	dsn := db_mysql.GetDSN("root", "123456", "localhost", 3306, "")
	var err error
	MysqlDb, err = db_mysql.Open(dsn)
	if err != nil {
		panic(err)
	}
	MysqlDialect, err = dialect.NewDialect(dialect.TypeMysql.Name)
	if err != nil {
		panic(err)
	}
	return
}

func TestMysqlLoad(t *testing.T) {
	initMysql()
	owners(MysqlDb, MysqlDialect)
}

func TestMysqlDDL(t *testing.T) {
	initMysql()
	owner := &dialect.OwnerModel{
		OwnerName: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, MysqlDialect, owner.OwnerName)
	testOwnerCreate(MysqlDb, MysqlDialect, owner)

	testDLL(MysqlDb, MysqlDialect, owner.OwnerName)
}

func TestMysqlSql(t *testing.T) {
	initMysql()
	sqlInfo := loadSql("sql_mysql.sql")
	owner := &dialect.OwnerModel{
		OwnerName: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, MysqlDialect, owner.OwnerName)
	testOwnerCreate(MysqlDb, MysqlDialect, owner)
	sqlInfo = "use " + owner.OwnerName + ";\n" + sqlInfo

	testSql(MysqlDb, MysqlDialect, owner.OwnerName, sqlInfo)
}
