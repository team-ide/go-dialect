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

func TestMysqlExportSql(t *testing.T) {
	initMysql()
	task := worker.NewTaskExport(MysqlDb, dialect.Mysql, dialect.Mysql, &worker.TaskExportParam{
		Owners: []*worker.TaskExportOwner{
			{SourceName: "information_schema", TargetName: "XXX1"},
			{SourceName: "mysql", TargetName: "XXX2"},
			{SourceName: "performance_schema", TargetName: "XXX3"},
		},
		ExportStructure: true,
		ExportData:      true,
		Dir:             "temp/export",
		ExportBatchSql:  true,
		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
		//	return tableName + "_" + index.Name
		//},

		DataSourceType: worker.DataSourceTypeSql,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}
func TestMysqlImportSql(t *testing.T) {
	initMysql()
	task := worker.NewTaskImport(MysqlDb, dialect.Mysql, &worker.TaskImportParam{
		Owners: []*worker.TaskImportOwner{
			{Name: "XXX1", Path: "temp/export/XXX1.sql"},
			{Name: "XXX2", Path: "temp/export/XXX2.sql"},
			{Name: "XXX3", Path: "temp/export/XXX3.sql"},
		},
		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
		//	return tableName + "_" + index.Name
		//},
		DataSourceType: worker.DataSourceTypeSql,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}
func TestMysqlExportStructure(t *testing.T) {
	initMysql()
	task := worker.NewTaskExport(MysqlDb, dialect.Mysql, dialect.Mysql, &worker.TaskExportParam{
		Owners: []*worker.TaskExportOwner{
			{SourceName: "information_schema", TargetName: "XXX1"},
			{SourceName: "mysql", TargetName: "XXX2"},
			{SourceName: "performance_schema", TargetName: "XXX3"},
		},
		ExportStructure: true,
		//ExportData:      true,
		Dir:            "temp/export",
		ExportBatchSql: true,
		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
		//	return tableName + "_" + index.Name
		//},

		DataSourceType: worker.DataSourceTypeSql,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}
func TestMysqlImportStructure(t *testing.T) {
	initMysql()
	task := worker.NewTaskImport(MysqlDb, dialect.Mysql, &worker.TaskImportParam{
		Owners: []*worker.TaskImportOwner{
			{Name: "XXX1", Path: "temp/export/XXX1.sql"},
			{Name: "XXX2", Path: "temp/export/XXX2.sql"},
			{Name: "XXX3", Path: "temp/export/XXX3.sql"},
		},
		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
		//	return tableName + "_" + index.Name
		//},
		DataSourceType: worker.DataSourceTypeSql,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}
func TestMysqlExportData(t *testing.T) {
	initMysql()
	task := worker.NewTaskExport(MysqlDb, dialect.Mysql, dialect.Mysql, &worker.TaskExportParam{
		Owners: []*worker.TaskExportOwner{
			{SourceName: "information_schema", TargetName: "XXX1"},
			{SourceName: "mysql", TargetName: "XXX2"},
			{SourceName: "performance_schema", TargetName: "XXX3"},
		},
		ExportStructure: true,
		ExportData:      true,
		Dir:             "temp/export",
		ExportBatchSql:  true,
		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
		//	return tableName + "_" + index.Name
		//},
		DataSourceType: worker.DataSourceTypeExcel,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}

func TestMysqlImportData(t *testing.T) {
	initMysql()
	task := worker.NewTaskImport(MysqlDb, dialect.Mysql, &worker.TaskImportParam{
		Owners: []*worker.TaskImportOwner{
			{Name: "XXX1", Dir: "temp/export/XXX1"},
			{Name: "XXX2", Dir: "temp/export/XXX2"},
			{Name: "XXX3", Dir: "temp/export/XXX3"},
		},
		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
		//	return tableName + "_" + index.Name
		//},
		DataSourceType: worker.DataSourceTypeExcel,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
}
