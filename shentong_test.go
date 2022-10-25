package go_dialect

import (
	"database/sql"
	"fmt"
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
	connStr := fmt.Sprintf("%s/%s@%s:%d/%s", "SYSDBA", "szoscar55", "127.0.0.1", 2003, "OSRDB")
	var err error
	ShenTongDb, err = sql.Open(db_shentong.GetDriverName(), connStr)
	if err != nil {
		panic(err)
	}
	return
}

func TestShenTong(t *testing.T) {
	initShenTong()
	databases(ShenTongDb, dialect.ShenTong)
}

func TestShenTongTableCreate(t *testing.T) {
	initShenTong()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(ShenTongDb, dialect.ShenTong, param, "", getTable().Name)
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
}

func TestShenTongSql(t *testing.T) {
	initShenTong()
	sqlInfo := loadSql("sql_shentong.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(ShenTongDb, sqlList)
	tables(ShenTongDb, dialect.ShenTong, "SYSDBA")
}
