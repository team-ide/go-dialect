package go_dialect

import (
	"context"
	"gitee.com/chunanyong/zorm"
	_ "github.com/mattn/go-oci8"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_shentong"
	"testing"
)

var (
	ShenTongContext context.Context
)

func initShenTong() {
	if ShenTongContext != nil {
		return
	}
	dbConfig := db_shentong.NewDataSourceConfig("SYSDBA", "szoscar55", "127.0.0.1", 2003, "OSRDB")
	dbDao, err := zorm.NewDBDao(&dbConfig)
	if err != nil {
		return
	}

	cxt := context.Background()
	ShenTongContext, err = dbDao.BindContextDBConnection(cxt)
	if err != nil {
		return
	}
	return
}

func TestShenTong(t *testing.T) {
	initShenTong()
	testDatabases(ShenTongContext, dialect.ShenTong)
}

func TestShenTongTableCreate(t *testing.T) {
	initShenTong()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(ShenTongContext, dialect.ShenTong, param, "", getTable().Name)
	testTableCreate(ShenTongContext, dialect.ShenTong, param, "", getTable())

	testColumnUpdate(ShenTongContext, dialect.ShenTong, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(ShenTongContext, dialect.ShenTong, param, "", getTable().Name, "detail3")
	testColumnAdd(ShenTongContext, dialect.ShenTong, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	testTable(ShenTongContext, dialect.ShenTong, "", getTable().Name)
}
