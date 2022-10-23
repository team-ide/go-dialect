package go_dialect

import (
	"context"
	"gitee.com/chunanyong/zorm"
	_ "github.com/mattn/go-oci8"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_kingbase_v8r3"
	"testing"
)

var (
	KinBaseContext context.Context
)

func initKinBase() {
	if KinBaseContext != nil {
		return
	}
	dbConfig := db_kingbase_v8r3.NewDataSourceConfig("SYSTEM", "123456", "127.0.0.1", 54321, "TEST")
	dbDao, err := zorm.NewDBDao(&dbConfig)
	if err != nil {
		return
	}

	cxt := context.Background()
	KinBaseContext, err = dbDao.BindContextDBConnection(cxt)
	if err != nil {
		return
	}
	return
}

func TestKinBase(t *testing.T) {
	initKinBase()
	testDatabases(KinBaseContext, dialect.KinBase)
}

func TestKinBaseTableCreate(t *testing.T) {
	initKinBase()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(KinBaseContext, dialect.KinBase, param, "", getTable().Name)
	testTableCreate(KinBaseContext, dialect.KinBase, param, "", getTable())

	testColumnUpdate(KinBaseContext, dialect.KinBase, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(KinBaseContext, dialect.KinBase, param, "", getTable().Name, "detail3")
	testColumnAdd(KinBaseContext, dialect.KinBase, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	testTable(KinBaseContext, dialect.KinBase, "", getTable().Name)
}
