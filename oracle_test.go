package go_dialect

import (
	"database/sql"
	//_ "github.com/mattn/go-oci8"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_oracle"
	"testing"
)

var (
	OracleDb *sql.DB
)

func initOracle() {
	if OracleDb != nil {
		return
	}
	dsn := db_oracle.GetDSN("root", "123456", "127.0.0.1", 1521, "xe")
	var err error
	OracleDb, err = db_oracle.Open(dsn)
	if err != nil {
		panic(err)
	}

	return
}

func TestOracleLoad(t *testing.T) {
	initOracle()
	owners(OracleDb, dialect.Oracle)
}

func TestOracleDDL(t *testing.T) {
	initOracle()
	testDLL(OracleDb, dialect.Oracle, "")
}

func TestOracleSql(t *testing.T) {
	initOracle()
	sqlInfo := loadSql("temp/sql_oracle.sql")
	testSql(OracleDb, dialect.Oracle, "ROOT", sqlInfo)
}
