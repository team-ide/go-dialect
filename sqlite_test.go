package go_dialect

import (
	"context"
	"gitee.com/chunanyong/zorm"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_sqlite3"
	"testing"
)

func getSqliteConfig() (dbContext context.Context) {
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
	dbContext, err = dbDao.BindContextDBConnection(cxt)
	if err != nil {
		return
	}
	return
}

func TestSqlite(t *testing.T) {
	testDatabases(getSqliteConfig(), dialect.Sqlite)
}

func TestSqliteTableCreate(t *testing.T) {
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	database := &dialect.DatabaseModel{
		Name: "TEST_DB",
	}
	testDatabaseDelete(getSqliteConfig(), dialect.Sqlite, param, database.Name)
	testDatabaseCreate(getSqliteConfig(), dialect.Sqlite, param, database)
	testTableDelete(getSqliteConfig(), dialect.Sqlite, param, database.Name, getTable().Name)
	testTableCreate(getSqliteConfig(), dialect.Sqlite, param, database.Name, getTable())

	testColumnUpdate(getSqliteConfig(), dialect.Sqlite, param, database.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(getSqliteConfig(), dialect.Sqlite, param, database.Name, getTable().Name, "detail3")
	testColumnAdd(getSqliteConfig(), dialect.Sqlite, param, database.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	testTables(getSqliteConfig(), dialect.Sqlite, database.Name)
}
