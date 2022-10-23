package go_dialect

import (
	"context"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_mysql"
	"testing"
)

var (
	MysqlContext context.Context
)

func initMysql() {
	if MysqlContext != nil {
		return
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", "root", "123456", "127.0.0.1", 3306, "")
	dbConfig := zorm.DataSourceConfig{
		//DSN 数据库的连接字符串
		DSN: connStr,
		//数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,clickhouse,dm,kingbase,aci 和DBType对应,处理数据库有多个驱动
		DriverName: db_mysql.GetDriverName(),
		//数据库类型(方言判断依据):mysql,postgresql,oracle,mssql,sqlite,clickhouse,dm,kingbase,shentong 和 DriverName 对应,处理数据库有多个驱动
		Dialect: db_mysql.GetDialect(),
	}
	dbDao, err := zorm.NewDBDao(&dbConfig)
	if err != nil {
		return
	}

	cxt := context.Background()
	MysqlContext, err = dbDao.BindContextDBConnection(cxt)
	if err != nil {
		return
	}
	return
}

func TestMysql(t *testing.T) {
	initMysql()
	testDatabases(MysqlContext, dialect.Mysql)
}

func TestMysqlTableCreate(t *testing.T) {
	initMysql()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	database := &dialect.DatabaseModel{
		Name: "TEST_DB",
	}
	testDatabaseDelete(MysqlContext, dialect.Mysql, param, database.Name)
	testDatabaseCreate(MysqlContext, dialect.Mysql, param, database)
	testTableDelete(MysqlContext, dialect.Mysql, param, database.Name, getTable().Name)
	testTableCreate(MysqlContext, dialect.Mysql, param, database.Name, getTable())

	testColumnUpdate(MysqlContext, dialect.Mysql, param, database.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(MysqlContext, dialect.Mysql, param, database.Name, getTable().Name, "detail3")
	testColumnAdd(MysqlContext, dialect.Mysql, param, database.Name, getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	testTables(MysqlContext, dialect.Mysql, database.Name)
}
