package go_dialect

import (
	"context"
	"gitee.com/chunanyong/zorm"
	_ "github.com/mattn/go-oci8"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_dm"
	"testing"
)

var (
	DaMenContext context.Context
)

func initDaMen() {
	if DaMenContext != nil {
		return
	}
	dbConfig := db_dm.NewDataSourceConfig("SYSDBA", "SYSDBA", "127.0.0.1", 5236)
	dbDao, err := zorm.NewDBDao(&dbConfig)
	if err != nil {
		return
	}

	cxt := context.Background()
	DaMenContext, err = dbDao.BindContextDBConnection(cxt)
	if err != nil {
		return
	}
	return
}

func TestDaMen(t *testing.T) {
	initDaMen()
	testDatabases(DaMenContext, dialect.DaMen)
}

func TestDaMenTableCreate(t *testing.T) {
	initDaMen()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(DaMenContext, dialect.DaMen, param, "", getTable().Name)
	testTableCreate(DaMenContext, dialect.DaMen, param, "", getTable())

	testColumnUpdate(DaMenContext, dialect.DaMen, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(DaMenContext, dialect.DaMen, param, "", getTable().Name, "detail3")
	testColumnAdd(DaMenContext, dialect.DaMen, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	testTable(DaMenContext, dialect.DaMen, "", getTable().Name)
}
