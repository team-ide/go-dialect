package main

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_kingbase_v8r3"
	"testing"
)

var (
	KinBaseDb *sql.DB
)

func initKinBase() {
	if KinBaseDb != nil {
		return
	}
	dsn := db_kingbase_v8r3.GetDSN("SYSTEM", "123456", "127.0.0.1", 54321, "TEST")
	var err error
	KinBaseDb, err = db_kingbase_v8r3.Open(dsn)
	if err != nil {
		panic(err)
	}
	return
}

func TestKinBaseLoad(t *testing.T) {
	initKinBase()
	owners(KinBaseDb, dialect.KinBase)
}

func TestKinBaseDDL(t *testing.T) {
	initKinBase()
	testDLL(KinBaseDb, dialect.KinBase, "")
}

func TestKinBaseSql(t *testing.T) {
	initKinBase()
	sqlInfo := loadSql("temp/sql_kinbase.sql")
	testSql(KinBaseDb, dialect.KinBase, "SYSDBA", sqlInfo)
}
