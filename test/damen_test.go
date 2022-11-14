package test

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_dm"
	"testing"
)

var (
	DMDb      *sql.DB
	DMDialect dialect.Dialect
)

func initDM() {
	if DMDb != nil {
		return
	}
	dsn := db_dm.GetDSN("SYSDBA", "SYSDBA", "127.0.0.1", 5236)
	var err error
	DMDb, err = db_dm.Open(dsn)
	if err != nil {
		panic(err)
	}
	DMDialect, err = dialect.NewDialect(dialect.TypeDM.Name)
	if err != nil {
		panic(err)
	}
	return
}

func TestDMLoad(t *testing.T) {
	initDM()
	owners(DMDb, DMDialect)
}

func TestDMDDL(t *testing.T) {
	initDM()
	//testTableDelete(DMDb, dialect.DM, param, "", getTable().Name)
	testDLL(DMDb, DMDialect, "")
}

func TestDMSql(t *testing.T) {
	initDM()
	sqlInfo := loadSql("temp/sql_dm.sql")
	testSql(DMDb, DMDialect, "SYSDBA", sqlInfo)
}
