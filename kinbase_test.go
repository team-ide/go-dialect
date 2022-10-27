package go_dialect

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_kingbase_v8r3"
	"strings"
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

func TestKinBase(t *testing.T) {
	initKinBase()
	owners(KinBaseDb, dialect.KinBase)
}

func TestKinBaseTableCreate(t *testing.T) {
	initKinBase()
	param := &dialect.GenerateParam{
		AppendOwner: true,
	}
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
	testTableDelete(KinBaseDb, dialect.KinBase, param, "", getTable().Name)
}

func TestKinBaseSql(t *testing.T) {
	initKinBase()
	sqlInfo := loadSql("temp/sql_kinbase.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(KinBaseDb, sqlList)
	tables(KinBaseDb, dialect.KinBase, "SYSTEM")
}
