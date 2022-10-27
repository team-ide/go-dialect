package go_dialect

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_mysql"
	"strings"
	"testing"
)

var (
	MysqlDb *sql.DB
)

func initMysql() {
	if MysqlDb != nil {
		return
	}
	dsn := db_mysql.GetDSN("root", "123456", "localhost", 3306, "")
	var err error
	MysqlDb, err = db_mysql.Open(dsn)
	if err != nil {
		panic(err)
	}
	return
}

func TestMysql(t *testing.T) {
	initMysql()
	owners(MysqlDb, dialect.Mysql)
}

func TestMysqlTableCreate(t *testing.T) {
	initMysql()
	param := &dialect.GenerateParam{
		AppendOwner: true,
	}
	owner := &dialect.OwnerModel{
		Name: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, dialect.Mysql, param, owner.Name)
	testOwnerCreate(MysqlDb, dialect.Mysql, param, owner)
	testTableDelete(MysqlDb, dialect.Mysql, param, owner.Name, getTable().Name)
	testTableCreate(MysqlDb, dialect.Mysql, param, owner.Name, getTable())

	testColumnUpdate(MysqlDb, dialect.Mysql, param, owner.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(MysqlDb, dialect.Mysql, param, owner.Name, getTable().Name, "detail3")
	testColumnAdd(MysqlDb, dialect.Mysql, param, owner.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(MysqlDb, dialect.Mysql, owner.Name, getTable().Name)
}

func TestMysqlSql(t *testing.T) {
	initMysql()
	sqlInfo := loadSql("sql_mysql.sql")
	param := &dialect.GenerateParam{
		AppendOwner: true,
	}
	owner := &dialect.OwnerModel{
		Name: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, dialect.Mysql, param, owner.Name)
	testOwnerCreate(MysqlDb, dialect.Mysql, param, owner)
	sqlInfo = "use " + owner.Name + ";\n" + sqlInfo

	sqlList := strings.Split(sqlInfo, ";\n")
	exec(MysqlDb, sqlList)
}
