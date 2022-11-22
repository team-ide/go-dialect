package main

import (
	"encoding/json"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"strings"
)

func doExport() {
	if *fileType == "" {
		println("请输入 文件 类型")
		return
	}
	if *exportDir == "" {
		println("请输入 导出 生成文件目录")
		return
	}
	if *exportOwner == "" {
		println("请输入 导出 库或表所属者")
		return
	}
	db, err := getDbInfo(*sourceDialect, *sourceUser, *sourcePassword, *sourceHost, *sourcePort, *sourceDatabase)
	if err != nil {
		panic(err)
	}
	dia, err := dialect.NewDialect(*sourceDialect)
	if err != nil {
		panic(err)
	}
	if db == nil || dia == nil {
		panic("sourceDialect [" + *sourceDialect + "] not support")
	}

	exportDia, err := dialect.NewDialect(*exportDialect)
	if err != nil {
		panic(err)
	}
	dataSourceType := worker.GetDataSource(*fileType)
	if dataSourceType == nil {
		panic("fileType [" + *fileType + "] not support")
	}

	var owners = getExportOwners(*exportOwner)
	task := worker.NewTaskExport(db, dia, exportDia, &worker.TaskExportParam{
		Owners:          owners,
		ExportStruct:    *exportStruct == "" || *exportStruct == "1" || *exportStruct == "true",
		ExportData:      *exportData == "" || *exportData == "1" || *exportData == "true",
		AppendOwnerName: *exportAppendOwner == "1" || *exportAppendOwner == "true",
		Dir:             *exportDir,
		ExportBatchSql:  true,
		FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
			return tableName + "_" + index.IndexName
		},
		DataSourceType: dataSourceType,
		OnProgress: func(progress *worker.TaskProgress) {
			bs, err := json.Marshal(progress)
			if err != nil {
				panic(err)
			}
			println(string(bs))
		},
	})
	err = task.Start()
	if err != nil {
		panic(err)
	}
	println("导出成功")
}

func getExportOwners(ownerInfoStr string) (owners []*worker.TaskExportOwner) {
	ownerStrList := strings.Split(ownerInfoStr, ",")
	for _, ownerStr := range ownerStrList {
		ss := strings.Split(ownerStr, "=")
		if len(ss) > 1 {
			owners = append(owners, &worker.TaskExportOwner{
				SourceName: strings.TrimSpace(ss[0]),
				TargetName: strings.TrimSpace(ss[1]),
			})
		} else if len(ss) > 0 {
			owners = append(owners, &worker.TaskExportOwner{
				SourceName: strings.TrimSpace(ss[0]),
			})
		}
	}
	return
}
