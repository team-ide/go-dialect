package go_dialect

import (
	"database/sql"
	"encoding/json"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"github.com/team-ide/go-driver/db_mysql"
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

func TestMysqlLoad(t *testing.T) {
	initMysql()
	owners(MysqlDb, dialect.Mysql)
}

func TestMysqlDDL(t *testing.T) {
	initMysql()
	owner := &dialect.OwnerModel{
		Name: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, dialect.Mysql, owner.Name)
	testOwnerCreate(MysqlDb, dialect.Mysql, owner)

	testDLL(MysqlDb, dialect.Mysql, owner.Name)
}

func TestMysqlSql(t *testing.T) {
	initMysql()
	sqlInfo := loadSql("sql_mysql.sql")
	owner := &dialect.OwnerModel{
		Name: "TEST_DB",
	}
	testOwnerDelete(MysqlDb, dialect.Mysql, owner.Name)
	testOwnerCreate(MysqlDb, dialect.Mysql, owner)
	sqlInfo = "use " + owner.Name + ";\n" + sqlInfo

	testSql(MysqlDb, dialect.Mysql, owner.Name, sqlInfo)
}

func TestMysqlSyncMysql(t *testing.T) {
	initMysql()
	task := worker.NewTaskSync(MysqlDb, dialect.Mysql, MysqlDb, dialect.Mysql, &worker.TaskSyncParam{
		Owners: []*worker.TaskSyncOwner{
			{SourceName: "information_schema", TargetName: "XXX1"},
			{SourceName: "mysql", TargetName: "XXX2"},
			{SourceName: "performance_schema", TargetName: "XXX3"},
		},
		SyncStructure: true,
		SyncData:      true,
		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
		//	return tableName + "_" + index.Name
		//},
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}

func TestMysqlSyncSqlite(t *testing.T) {
	initMysql()
	initSqlite()
	task := worker.NewTaskSync(MysqlDb, dialect.Mysql, SqliteDb, dialect.Sqlite, &worker.TaskSyncParam{
		Owners: []*worker.TaskSyncOwner{
			{SourceName: "mysql", TargetName: "main"},
		},
		SyncStructure: true,
		SyncData:      true,
		FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
			return tableName + "_" + index.Name
		},
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}
