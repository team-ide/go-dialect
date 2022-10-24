package go_dialect

import (
	"database/sql"
	"fmt"
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
	connStr := fmt.Sprintf("user='%s' password='%s' host=%s port=%d dbname=%s sslmode=disable", "SYSTEM", "123456", "127.0.0.1", 54321, "TEST")
	var err error
	KinBaseDb, err = sql.Open(db_kingbase_v8r3.GetDriverName(), connStr)
	if err != nil {
		panic(err)
	}
	return
}

func TestKinBase(t *testing.T) {
	initKinBase()
	databases(KinBaseDb, dialect.KinBase)
}

func TestKinBaseTableCreate(t *testing.T) {
	initKinBase()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(KinBaseDb, dialect.KinBase, param, "", getTable().Name)
	testTableCreate(KinBaseDb, dialect.KinBase, param, "", getTable())

	testColumnUpdate(KinBaseDb, dialect.KinBase, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(KinBaseDb, dialect.KinBase, param, "", getTable().Name, "detail3")
	testColumnAdd(KinBaseDb, dialect.KinBase, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(KinBaseDb, dialect.KinBase, "", getTable().Name)
}
