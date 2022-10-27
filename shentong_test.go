package go_dialect

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_shentong"
	"strings"
	"testing"
)

var (
	ShenTongDb *sql.DB
)

func initShenTong() {
	if ShenTongDb != nil {
		return
	}
	dsn := db_shentong.GetDSN("SYSDBA", "szoscar55", "127.0.0.1", 2003, "OSRDB")
	var err error
	ShenTongDb, err = db_shentong.Open(dsn)
	if err != nil {
		panic(err)
	}
	return
}

func TestShenTong(t *testing.T) {
	initShenTong()
	owners(ShenTongDb, dialect.ShenTong)
}

func TestShenTongTableCreate(t *testing.T) {
	initShenTong()
	param := &dialect.GenerateParam{
		AppendOwner: true,
	}
	testTableCreate(ShenTongDb, dialect.ShenTong, param, "", getTable())

	testColumnUpdate(ShenTongDb, dialect.ShenTong, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(ShenTongDb, dialect.ShenTong, param, "", getTable().Name, "detail3")
	testColumnAdd(ShenTongDb, dialect.ShenTong, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(ShenTongDb, dialect.ShenTong, "", getTable().Name)
	testTableDelete(ShenTongDb, dialect.ShenTong, param, "", getTable().Name)
}

func TestShenTongSql(t *testing.T) {
	initShenTong()
	sqlInfo := loadSql("temp/sql_shentong.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(ShenTongDb, sqlList)
	tables(ShenTongDb, dialect.ShenTong, "SYSDBA")
}
