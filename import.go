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

	password := *importOwnerCreatePassword
	if password == "" {
		password = *sourcePassword
	}
	var owners = getImportOwners(*importOwner)
	for _, owner := range owners {
		owner.Password = password
	}

	task := worker.NewTaskImport(db, dia,
		func(owner *worker.TaskImportOwner) (workDb *sql.DB, err error) {
			ownerName := owner.Name
			ownerUsername := owner.Username
			ownerPassword := owner.Password

			if ownerPassword == "" {
				ownerPassword = password
			}

			if *sourceDialect == "sqlite" || *sourceDialect == "sqlite3" {
				workDb = db
				return
			}
			if *sourceDialect == "mysql" {
				if ownerUsername == "" {
					ownerUsername = *sourceUser
				}
				workDb, err = getDbInfo(*sourceDialect, ownerUsername, ownerPassword, *sourceHost, *sourcePort, ownerName)
				return
			}
			if ownerUsername == "" {
				ownerUsername = ownerName
			}
			workDb, err = getDbInfo(*sourceDialect, ownerName, password, *sourceHost, *sourcePort, *sourceDatabase)
			return
		},
		&worker.TaskImportParam{
			Owners:                owners,
			OwnerCreateIfNotExist: *importOwnerCreateIfNotExist == "1" || *importOwnerCreateIfNotExist == "true",
			FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
				return tableName + "_" + index.IndexName
			},
			DataSourceType: dataSourceType,
			ErrorContinue:  true,
			OnProgress: func(progress *worker.TaskProgress) {
				progress.OnError = func(err error) {
					dataBytes, _ := json.Marshal(progress)
					println("progress:" + string(dataBytes))
					println("progress error:" + err.Error())
				}
				//println(string(bs))
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
