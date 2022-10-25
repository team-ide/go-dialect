package go_dialect

import (
	"database/sql"
	"fmt"
	"strings"

	//_ "github.com/mattn/go-oci8"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_oracle"
	"testing"
)

var (
	OracleDb *sql.DB
)

func initOracle() {
	if OracleDb != nil {
		return
	}
	//connStr := fmt.Sprintf("%s/%s@%s:%d/%s", "root", "123456", "127.0.0.1", 1521, "xe")
	connStr := fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, "root", "123456", "127.0.0.1", 1521, "xe")
	var err error
	OracleDb, err = sql.Open(db_oracle.GetDriverName(), connStr)
	if err != nil {
		panic(err)
	}

	//config := zorm.DataSourceConfig{
	//	DSN:        connStr,
	//	DriverName: "oci8",
	//	Dialect:    "oracle",
	//}
	//dao, err := zorm.NewDBDao(&config)
	//if err != nil {
	//	panic(err)
	//}
	//sqlInfo := "SELECT t.* from all_tab_columns t WHERE 1=1 AND TABLE_NAME='USER_INFO'"
	//finder := zorm.NewFinder()
	//finder.InjectionCheck = false
	//finder.Append(sqlInfo)
	//
	//ctx, err := dao.BindContextDBConnection(context.Background())
	//res, err := zorm.QueryMap(ctx, finder, nil)
	//if err != nil {
	//	panic(err)
	//}
	//for _, one := range res {
	//	for key, value := range one {
	//		if value == nil {
	//			continue
	//		}
	//		println("key [" + key + "] value [" + dialect.GetStringValue(value) + "] valueType [" + reflect.TypeOf(value).String() + "]")
	//	}
	//}
	//panic("stopped")

	return
}

func TestOracle(t *testing.T) {
	initOracle()
	databases(OracleDb, dialect.Oracle)
}

func TestOracleTableCreate(t *testing.T) {
	initOracle()
	param := &dialect.GenerateParam{
		AppendDatabase: true,
	}
	testTableDelete(OracleDb, dialect.Oracle, param, "", getTable().Name)
	testTableCreate(OracleDb, dialect.Oracle, param, "", getTable())

	testColumnUpdate(OracleDb, dialect.Oracle, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name1",
		Type:    "varchar",
		Length:  500,
		Comment: "name1注释",
		OldName: "name",
	})
	testColumnDelete(OracleDb, dialect.Oracle, param, "", getTable().Name, "detail3")
	testColumnAdd(OracleDb, dialect.Oracle, param, "", getTable().Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(OracleDb, dialect.Oracle, "", getTable().Name)
}

func TestOracleSql(t *testing.T) {
	initOracle()
	sqlInfo := loadSql("temp/sql_oracle.sql")
	sqlList := strings.Split(sqlInfo, ";\n")
	exec(OracleDb, sqlList)
	tables(OracleDb, dialect.Oracle, "ROOT")
}
