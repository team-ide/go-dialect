package main

//func TestMysqlSyncMysql(t *testing.T) {
//	initMysql()
//	task := worker.NewTaskSync(MysqlDb, MysqlDialect, MysqlDb, MysqlDialect, func(ownerName string) (db *sql.DB, err error) {
//
//		return
//	},
//		&worker.TaskSyncParam{
//			Owners: []*worker.TaskSyncOwner{
//				{SourceName: "information_schema", TargetName: "XXX1"},
//				{SourceName: "mysql", TargetName: "XXX2"},
//				{SourceName: "performance_schema", TargetName: "XXX3"},
//			},
//			SyncStruct: true,
//			SyncData:   true,
//			//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//			//	return tableName + "_" + index.Name
//			//},
//		})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//}
//
//func TestMysqlSyncSqlite(t *testing.T) {
//	initMysql()
//	initSqlite()
//	task := worker.NewTaskSync(MysqlDb, MysqlDialect, SqliteDb, SqliteDialect, func(ownerName string) (db *sql.DB, err error) {
//
//		return
//	},
//		&worker.TaskSyncParam{
//			Owners: []*worker.TaskSyncOwner{
//				{SourceName: "mysql", TargetName: "main"},
//			},
//			SyncStruct: true,
//			SyncData:   true,
//			FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//				return tableName + "_" + index.Name
//			},
//		})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//}
//
//func TestMysqlExportSql(t *testing.T) {
//	initMysql()
//	_, _ = worker.DoExec(MysqlDb, []string{"set global show_compatibility_56=on"})
//	owners := []*worker.TaskExportOwner{
//		//{SourceName: "information_schema", TargetName: "main"},
//		{SourceName: "mysql", TargetName: "main"},
//		//{SourceName: "performance_schema", TargetName: "main"},
//	}
//	exportSql(MysqlDb, MysqlDialect, dialect.Sqlite, owners)
//}
//func TestMysqlImportSql(t *testing.T) {
//	initMysql()
//	owners := []*worker.TaskImportOwner{
//		{Name: "information_schema", Path: "temp/export/XXX1.sql"},
//		{Name: "mysql", Path: "temp/export/XXX2.sql"},
//		{Name: "performance_schema", Path: "temp/export/XXX3.sql"},
//	}
//	importSql(MysqlDb, MysqlDialect, owners)
//}
//func TestSqliteImportSql(t *testing.T) {
//	initSqlite()
//	owners := []*worker.TaskImportOwner{
//		{Name: "main", Path: "temp/export/main.sql"},
//	}
//	importSql(SqliteDb, dialect.Sqlite, owners)
//}
//func TestMysqlExportStructure(t *testing.T) {
//	initMysql()
//	task := worker.NewTaskExport(MysqlDb, MysqlDialect, MysqlDialect, &worker.TaskExportParam{
//		Owners: []*worker.TaskExportOwner{
//			{SourceName: "information_schema", TargetName: "XXX1"},
//			{SourceName: "mysql", TargetName: "XXX2"},
//			{SourceName: "performance_schema", TargetName: "XXX3"},
//		},
//		ExportStructure: true,
//		//ExportData:      true,
//		Dir:            "temp/export",
//		ExportBatchSql: true,
//		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//		//	return tableName + "_" + index.Name
//		//},
//
//		DataSourceType: worker.DataSourceTypeSql,
//	})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//}
//func TestMysqlImportStructure(t *testing.T) {
//	initMysql()
//	task := worker.NewTaskImport(MysqlDb, MysqlDialect, &worker.TaskImportParam{
//		Owners: []*worker.TaskImportOwner{
//			{Name: "XXX1", Path: "temp/export/XXX1.sql"},
//			{Name: "XXX2", Path: "temp/export/XXX2.sql"},
//			{Name: "XXX3", Path: "temp/export/XXX3.sql"},
//		},
//		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//		//	return tableName + "_" + index.Name
//		//},
//		DataSourceType: worker.DataSourceTypeSql,
//	})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//}
//func TestMysqlExportData(t *testing.T) {
//	initMysql()
//	task := worker.NewTaskExport(MysqlDb, MysqlDialect, MysqlDialect, &worker.TaskExportParam{
//		Owners: []*worker.TaskExportOwner{
//			{SourceName: "information_schema", TargetName: "XXX1"},
//			{SourceName: "mysql", TargetName: "XXX2"},
//			{SourceName: "performance_schema", TargetName: "XXX3"},
//		},
//		ExportStructure: true,
//		ExportData:      true,
//		Dir:             "temp/export",
//		ExportBatchSql:  true,
//		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//		//	return tableName + "_" + index.Name
//		//},
//		DataSourceType: worker.DataSourceTypeExcel,
//	})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//}
//
//func TestMysqlImportData(t *testing.T) {
//	initMysql()
//	task := worker.NewTaskImport(MysqlDb, MysqlDialect, &worker.TaskImportParam{
//		Owners: []*worker.TaskImportOwner{
//			{Name: "XXX1", Dir: "temp/export/XXX1"},
//			{Name: "XXX2", Dir: "temp/export/XXX2"},
//			{Name: "XXX3", Dir: "temp/export/XXX3"},
//		},
//		//FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//		//	return tableName + "_" + index.Name
//		//},
//		DataSourceType: worker.DataSourceTypeExcel,
//	})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//}
//
//func exportSql(db *sql.DB, dia dialect.Dialect, targetDia dialect.Dialect, owners []*worker.TaskExportOwner) {
//	task := worker.NewTaskExport(db, dia, targetDia, &worker.TaskExportParam{
//		Owners:          owners,
//		ExportStructure: true,
//		ExportData:      true,
//		Dir:             "temp/export",
//		ExportBatchSql:  true,
//		FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//			return tableName + "_" + index.Name
//		},
//		DataSourceType: worker.DataSourceTypeSql,
//	})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//	time.Sleep(1 * time.Second)
//}
//func importSql(db *sql.DB, dia dialect.Dialect, owners []*worker.TaskImportOwner) {
//	task := worker.NewTaskImport(db, dia, &worker.TaskImportParam{
//		Owners: owners,
//		FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
//			return tableName + "_" + index.Name
//		},
//		DataSourceType: worker.DataSourceTypeSql,
//	})
//	err := task.Start()
//	if err != nil {
//		panic(err)
//	}
//	bs, _ := json.Marshal(task)
//	println(string(bs))
//	time.Sleep(1 * time.Second)
//}
