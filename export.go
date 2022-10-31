package main

import (
	"encoding/json"
	"flag"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"strings"
)

var (
	exportType  = flag.String("exportType", "", "导出 类型：sql、excel、txt、csv")
	exportDir   = flag.String("exportDir", "", "导出 文件存储目录")
	exportOwner = flag.String("exportOwner", "", "导出 库或表所属者，多个使用“,”隔开，“x,xx=xx1”")
	exportDia   = flag.String("exportDia", "", "导出 方言：mysql、sqlite、dm、kingbase、oracle")
)

func doExport() {
	if *exportType == "" {
		println("请输入 导出 类型")
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
	db, err := getDbInfo(*dbType, *user, *password, *host, *port, *database)
	if err != nil {
		panic(err)
	}
	dia := dialect.GetDialect(*dbType)
	if db == nil || dia == nil {
		panic("dbType [" + *dbType + "] not support")
	}

	targetDialect := dialect.GetDialect(*exportDia)
	dataSourceType := worker.GetDataSource(*exportType)
	if dataSourceType == nil {
		panic("export [" + *exportType + "] not support")
	}
	var owners = getExportOwners(*exportOwner)
	task := worker.NewTaskExport(db, dia, targetDialect, &worker.TaskExportParam{
		Owners:          owners,
		ExportStructure: true,
		ExportData:      true,
		Dir:             *exportDir,
		ExportBatchSql:  true,
		FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
			return tableName + "_" + index.Name
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
