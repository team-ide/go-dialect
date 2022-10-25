package go_dialect

import (
	"context"
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_sqlite3"
	"strings"
	"testing"
)

var (
	SqliteDb *sql.DB
)

func initSqlite() (dbContext context.Context) {
	if SqliteDb != nil {
		return
	}
	connStr := "temp/test_sqlite"
	var err error
	SqliteDb, err = sql.Open(db_sqlite3.GetDriverName(), connStr)
	if err != nil {
		panic(err)
	}
	return
}

func TestSqlite(t *testing.T) {
	initSqlite()
	databases(SqliteDb, dialect.Sqlite)
}

func TestSqliteTableCreate(t *testing.T) {
	initSqlite()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	//testTableDelete(SqliteDb, dialect.Sqlite, param, "", getTable().Name)
	testTableCreate(SqliteDb, dialect.Sqlite, param, "", getTable())

	testColumnUpdate(SqliteDb, dialect.Sqlite, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(SqliteDb, dialect.Sqlite, param, "", getTable().Name, "detail3")
	testColumnAdd(SqliteDb, dialect.Sqlite, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(SqliteDb, dialect.Sqlite, "", getTable().Name)
}

func TestSqliteSql(t *testing.T) {
	initSqlite()
	sqlInfo := loadSql("temp/sql_sqlite.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(SqliteDb, sqlList)
	tables(SqliteDb, dialect.Sqlite, "")
}
