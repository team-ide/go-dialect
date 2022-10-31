package main

import (
	"encoding/json"
	"flag"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"strings"
)

var (
	importType  = flag.String("importType", "", "导入 类型：sql、excel、txt、csv")
	importOwner = flag.String("importOwner", "", "导入 库或表所属者，多个使用“,”隔开，指定文件或目录，“xx=data/xx”")
)

func doImport() {
	if *importType == "" {
		println("请输入 导出类型")
		return
	}
	if *importOwner == "" {
		println("请输入 库或表所属者名称")
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

	dataSourceType := worker.GetDataSource(*importType)
	if dataSourceType == nil {
		panic("import [" + *importType + "] not support")
	}
	var owners = getImportOwners(*importOwner)
	task := worker.NewTaskImport(db, dia, &worker.TaskImportParam{
		Owners: owners,
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
	println("导入成功")
}

func getImportOwners(ownerInfoStr string) (owners []*worker.TaskImportOwner) {
	ownerStrList := strings.Split(ownerInfoStr, ",")
	for _, ownerStr := range ownerStrList {
		ss := strings.Split(ownerStr, "=")
		if len(ss) > 1 {
			owners = append(owners, &worker.TaskImportOwner{
				Name: strings.TrimSpace(ss[0]),
				Path: strings.TrimSpace(ss[1]),
			})
		} else if len(ss) > 0 {
			owners = append(owners, &worker.TaskImportOwner{
				Name: strings.TrimSpace(ss[0]),
			})
		}
	}
	return
}
