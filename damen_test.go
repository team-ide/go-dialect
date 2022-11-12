package main

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_dm"
	"testing"
)

var (
	DaMenDb      *sql.DB
	DaMenDialect dialect.Dialect
)

func initDaMen() {
	if DaMenDb != nil {
		return
	}
	dsn := db_dm.GetDSN("SYSDBA", "SYSDBA", "127.0.0.1", 5236)
	var err error
	DaMenDb, err = db_dm.Open(dsn)
	if err != nil {
		panic(err)
	}
	DaMenDialect, err = dialect.NewDialect(dialect.TypeDaMen.Name)
	if err != nil {
		panic(err)
	}
	return
}

func TestDaMenLoad(t *testing.T) {
	initDaMen()
	owners(DaMenDb, DaMenDialect)
}

func TestDaMenDDL(t *testing.T) {
	initDaMen()
	//testTableDelete(DaMenDb, dialect.DaMen, param, "", getTable().Name)
	testDLL(DaMenDb, DaMenDialect, "")
}

func TestDaMenSql(t *testing.T) {
	initDaMen()
	sqlInfo := loadSql("temp/sql_damen.sql")
	testSql(DaMenDb, DaMenDialect, "SYSDBA", sqlInfo)
}
