package go_dialect

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_dm"
	"strings"
	"testing"
)

var (
	DaMenDb *sql.DB
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
	return
}

func TestDaMen(t *testing.T) {
	initDaMen()
	owners(DaMenDb, dialect.DaMen)
}

func TestDaMenTableCreate(t *testing.T) {
	initDaMen()
	param := &dialect.GenerateParam{
		AppendOwner: true,
	}
	//testTableDelete(DaMenDb, dialect.DaMen, param, "", getTable().Name)
	testTableCreate(DaMenDb, dialect.DaMen, param, "", getTable())

	testColumnUpdate(DaMenDb, dialect.DaMen, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(DaMenDb, dialect.DaMen, param, "", getTable().Name, "detail3")
	testColumnAdd(DaMenDb, dialect.DaMen, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(DaMenDb, dialect.DaMen, "", getTable().Name)
	testTableDelete(DaMenDb, dialect.DaMen, param, "", getTable().Name)
}

func TestDaMenSql(t *testing.T) {
	initDaMen()
	sqlInfo := loadSql("temp/sql_damen.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(DaMenDb, sqlList)
	tables(DaMenDb, dialect.DaMen, "SYSDBA")
}
