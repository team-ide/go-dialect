package main

import (
	"database/sql"
	"encoding/json"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"strings"
)

func doImport() {
	if *fileType == "" {
		println("请输入 文件 类型")
		return
	}
	if *importOwner == "" {
		println("请输入 库或表所属者名称")
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

	dataSourceType := worker.GetDataSource(*fileType)
	if dataSourceType == nil {
		panic("fileType [" + *fileType + "] not support")
	}
	var owners = getImportOwners(*importOwner)

	password := *importOwnerCreatePassword
	if password == "" {
		password = *sourcePassword
	}
	task := worker.NewTaskImport(db, dia,
		func(ownerName string) (workDb *sql.DB, err error) {

			if *sourceDialect == "sqlite" || *sourceDialect == "sqlite3" {
				workDb = db
				return
			}
			if *sourceDialect == "mysql" {
				workDb, err = getDbInfo(*sourceDialect, *sourceUser, password, *sourceHost, *sourcePort, ownerName)
				return
			}
			workDb, err = getDbInfo(*sourceDialect, ownerName, password, *sourceHost, *sourcePort, *sourceDatabase)
			return
		},
		&worker.TaskImportParam{
			Owners:                      owners,
			ImportOwnerCreateIfNotExist: *importOwnerCreateIfNotExist == "1" || *importOwnerCreateIfNotExist == "true",
			ImportOwnerCreatePassword:   password,
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
