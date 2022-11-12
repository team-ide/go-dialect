package main

import (
	"context"
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_sqlite3"
	"testing"
)

var (
	SqliteDb      *sql.DB
	SqliteDialect dialect.Dialect
)

func initSqlite() (dbContext context.Context) {
	if SqliteDb != nil {
		return
	}
	dsn := db_sqlite3.GetDSN("temp/test_sqlite")
	var err error
	SqliteDb, err = db_sqlite3.Open(dsn)
	if err != nil {
		panic(err)
	}
	SqliteDialect, err = dialect.NewDialect(dialect.TypeSqlite.Name)
	if err != nil {
		panic(err)
	}
	return
}

func TestSqliteLoad(t *testing.T) {
	initSqlite()
	owners(SqliteDb, SqliteDialect)
}

func TestSqliteDDL(t *testing.T) {
	initSqlite()
	testDLL(SqliteDb, SqliteDialect, "")
}

func TestSqliteSql(t *testing.T) {
	initSqlite()
	sqlInfo := loadSql("temp/sql_sqlite.sql")
	testSql(SqliteDb, SqliteDialect, "SYSDBA", sqlInfo)
}
