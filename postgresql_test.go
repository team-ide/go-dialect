package go_dialect

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_postgresql"
	"strings"
	"testing"
)

var (
	PostgresqlDb *sql.DB
)

func initPostgresql() {
	if PostgresqlDb != nil {
		return
	}
	var err error
	PostgresqlDb, err = db_postgresql.Open(db_postgresql.GetDSN("postgres", "123456", "127.0.0.1", 5432, "postgres"))
	if err != nil {
		panic(err)
	}
	return
}

func TestPostgresql(t *testing.T) {
	initPostgresql()
	databases(PostgresqlDb, dialect.Postgresql)
}

func TestPostgresqlTableCreate(t *testing.T) {
	initPostgresql()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	//testTableDelete(PostgresqlDb, dialect.Postgresql, param, "", getTable().Name)
	testTableCreate(PostgresqlDb, dialect.Postgresql, param, "", getTable())

	testColumnUpdate(PostgresqlDb, dialect.Postgresql, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(PostgresqlDb, dialect.Postgresql, param, "", getTable().Name, "detail3")
	testColumnAdd(PostgresqlDb, dialect.Postgresql, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(PostgresqlDb, dialect.Postgresql, "", getTable().Name)
}

func TestPostgresqlSql(t *testing.T) {
	initPostgresql()
	sqlInfo := loadSql("temp/sql_kinbase.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(PostgresqlDb, sqlList)
	tables(PostgresqlDb, dialect.Postgresql, "SYSTEM")
}
