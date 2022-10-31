package main

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_postgresql"
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

func TestPostgresqlLoad(t *testing.T) {
	initPostgresql()
	owners(PostgresqlDb, dialect.Postgresql)
}

func TestPostgresqlDDL(t *testing.T) {
	initPostgresql()
	//testTableDelete(PostgresqlDb, dialect.Postgresql, param, "", getTable().Name)
	testDLL(PostgresqlDb, dialect.Postgresql, "")
}

func TestPostgresqlSql(t *testing.T) {
	initPostgresql()
	sqlInfo := loadSql("temp/sql_kinbase.sql")
	testSql(PostgresqlDb, dialect.Postgresql, "ROOT", sqlInfo)
}
