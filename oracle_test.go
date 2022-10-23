package go_dialect

import (
	"context"
	"fmt"
	"gitee.com/chunanyong/zorm"
	_ "github.com/mattn/go-oci8"
	"github.com/team-ide/go-dialect/dialect"
	"testing"
)

var (
	OracleContext context.Context
)

func initOracle() {
	if OracleContext != nil {
		return
	}
	//dbConfig := db_oracle.NewDataSourceConfig("root", "123456", "127.0.0.1", 1521, "xe")
	connStr := fmt.Sprintf("%s/%s@%s:%d/%s", "root", "123456", "127.0.0.1", 1521, "xe")
	dbConfig := zorm.DataSourceConfig{
		//DSN 数据库的连接字符串
		DSN: connStr,
		//数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,clickhouse,dm,kingbase,aci 和DBType对应,处理数据库有多个驱动
		DriverName: "oci8",
		//数据库类型(方言判断依据):mysql,postgresql,oracle,mssql,sqlite,clickhouse,dm,kingbase,shentong 和 DriverName 对应,处理数据库有多个驱动
		Dialect: "oracle",
	}
	dbDao, err := zorm.NewDBDao(&dbConfig)
	if err != nil {
		return
	}

	cxt := context.Background()
	OracleContext, err = dbDao.BindContextDBConnection(cxt)
	if err != nil {
		return
	}
	return
}

func TestOracle(t *testing.T) {
	initOracle()
	testDatabases(OracleContext, dialect.Oracle)
}

func TestOracleTableCreate(t *testing.T) {
	initOracle()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(OracleContext, dialect.Oracle, param, "", getTable().Name)
	testTableCreate(OracleContext, dialect.Oracle, param, "", getTable())

	testColumnUpdate(OracleContext, dialect.Oracle, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(OracleContext, dialect.Oracle, param, "", getTable().Name, "detail3")
	testColumnAdd(OracleContext, dialect.Oracle, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	testTable(OracleContext, dialect.Oracle, "", getTable().Name)
}
