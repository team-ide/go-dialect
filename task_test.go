package go_dialect

import (
	"database/sql"
	"encoding/json"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"time"
)

func exportSql(db *sql.DB, dia dialect.Dialect, targetDia dialect.Dialect, owners []*worker.TaskExportOwner) {
	task := worker.NewTaskExport(db, dia, targetDia, &worker.TaskExportParam{
		Owners:          owners,
		ExportStructure: true,
		ExportData:      true,
		Dir:             "temp/export",
		ExportBatchSql:  true,
		FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
			return tableName + "_" + index.Name
		},
		DataSourceType: worker.DataSourceTypeSql,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
	time.Sleep(1 * time.Second)
}
func importSql(db *sql.DB, dia dialect.Dialect, owners []*worker.TaskImportOwner) {
	task := worker.NewTaskImport(db, dia, &worker.TaskImportParam{
		Owners: owners,
		FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
			return tableName + "_" + index.Name
		},
		DataSourceType: worker.DataSourceTypeSql,
	})
	err := task.Start()
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(task)
	println(string(bs))
	time.Sleep(1 * time.Second)
}
