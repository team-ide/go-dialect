package go_dialect

import (
	"context"
	"gitee.com/chunanyong/zorm"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_sqlite3"
	"testing"
)

var (
	SqliteContext context.Context
)

func initSqlite() (dbContext context.Context) {
	if SqliteContext != nil {
		return
	}
	connStr := "test_sqlite"
	dbConfig := zorm.DataSourceConfig{
		//DSN 数据库的连接字符串
		DSN: connStr,
		//数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,clickhouse,dm,kingbase,aci 和DBType对应,处理数据库有多个驱动
		DriverName: db_sqlite3.GetDriverName(),
		//数据库类型(方言判断依据):mysql,postgresql,oracle,mssql,sqlite,clickhouse,dm,kingbase,shentong 和 DriverName 对应,处理数据库有多个驱动
		Dialect: db_sqlite3.GetDialect(),
	}
	dbDao, err := zorm.NewDBDao(&dbConfig)
	if err != nil {
		return
	}

	cxt := context.Background()
	SqliteContext, err = dbDao.BindContextDBConnection(cxt)
	if err != nil {
		return
	}
	return
}

func TestSqlite(t *testing.T) {
	initSqlite()
	testDatabases(SqliteContext, dialect.Sqlite)
}

func TestSqliteTableCreate(t *testing.T) {
	initSqlite()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(SqliteContext, dialect.Sqlite, param, "", getTable().Name)
	testTableCreate(SqliteContext, dialect.Sqlite, param, "", getTable())

	testColumnUpdate(SqliteContext, dialect.Sqlite, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(SqliteContext, dialect.Sqlite, param, "", getTable().Name, "detail3")
	testColumnAdd(SqliteContext, dialect.Sqlite, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	testTables(SqliteContext, dialect.Sqlite, "")
}
