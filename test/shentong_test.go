package test

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_shentong"
	"testing"
)

var (
	ShenTongDb      *sql.DB
	ShenTongDialect dialect.Dialect
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
	ShenTongDialect, err = dialect.NewDialect(dialect.TypeShenTong.Name)
	if err != nil {
		panic(err)
	}
	return
}

func TestShenTongLoad(t *testing.T) {
	initShenTong()
	owners(ShenTongDb, ShenTongDialect)
}

func TestShenTongDDL(t *testing.T) {
	initShenTong()
	testDLL(ShenTongDb, ShenTongDialect, "")
}

func TestShenTongSql(t *testing.T) {
	initShenTong()
	sqlInfo := loadSql("temp/sql_shentong.sql")
	testSql(ShenTongDb, ShenTongDialect, "SYSDBA", sqlInfo)
}
