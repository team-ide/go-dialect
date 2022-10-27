package go_dialect

import (
	"database/sql"
	"strings"

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

func TestOracle(t *testing.T) {
	initOracle()
	owners(OracleDb, dialect.Oracle)
}

func TestOracleTableCreate(t *testing.T) {
	initOracle()
	param := &dialect.GenerateParam{
		AppendOwner: true,
	}
	testTableCreate(OracleDb, dialect.Oracle, param, "", getTable())

	testColumnUpdate(OracleDb, dialect.Oracle, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(OracleDb, dialect.Oracle, param, "", getTable().Name, "detail3")
	testColumnAdd(OracleDb, dialect.Oracle, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(OracleDb, dialect.Oracle, "", getTable().Name)
	testTableDelete(OracleDb, dialect.Oracle, param, "", getTable().Name)
}

func TestOracleSql(t *testing.T) {
	initOracle()
	sqlInfo := loadSql("temp/sql_oracle.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(OracleDb, sqlList)
	tables(OracleDb, dialect.Oracle, "ROOT")
}
