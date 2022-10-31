package main

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_mysql"
	"testing"
)

var (
	MysqlDb *sql.DB
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
	return
}

func TestMysqlLoad(t *testing.T) {
	initMysql()
	owners(MysqlDb, dialect.Mysql)
}

func TestMysqlDDL(t *testing.T) {
	initMysql()
	owner := &dialect.OwnerModel{
		Name: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, dialect.Mysql, owner.Name)
	testOwnerCreate(MysqlDb, dialect.Mysql, owner)

	testDLL(MysqlDb, dialect.Mysql, owner.Name)
}

func TestMysqlSql(t *testing.T) {
	initMysql()
	sqlInfo := loadSql("sql_mysql.sql")
	owner := &dialect.OwnerModel{
		Name: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, dialect.Mysql, owner.Name)
	testOwnerCreate(MysqlDb, dialect.Mysql, owner)
	sqlInfo = "use " + owner.Name + ";\n" + sqlInfo

	testSql(MysqlDb, dialect.Mysql, owner.Name, sqlInfo)
}
