package main

import (
	"database/sql"
	"encoding/json"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"strings"
)

func doSync() {
	if *targetDialect == "" {
		println("请输入 同步 目标 数据库类型")
		return
	}
	if *syncOwner == "" {
		println("请输入 同步 库或表所属者")
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

	targetDb, err := getDbInfo(*targetDialect, *targetUser, *targetPassword, *targetHost, *targetPort, *targetDatabase)
	if err != nil {
		panic(err)
	}
	targetDia, err := dialect.NewDialect(*targetDialect)
	if err != nil {
		panic(err)
	}
	if targetDb == nil || targetDia == nil {
		panic("targetDialect [" + *targetDialect + "] not support")
	}

	password := *importOwnerCreatePassword
	if password == "" {
		password = *sourcePassword
	}
	var owners = getSyncOwners(*syncOwner)
	for _, owner := range owners {
		owner.Password = password
	}
	task := worker.NewTaskSync(db, dia, targetDb, targetDia,
		func(owner *worker.TaskSyncOwner) (workDb *sql.DB, err error) {
			ownerName := owner.TargetName
			if ownerName == "" {
				ownerName = owner.SourceName
			}
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
		&worker.TaskSyncParam{
			Owners:                owners,
			SyncStruct:            *syncStruct == "" || *syncStruct == "1" || *syncStruct == "true",
			SyncData:              *syncData == "" || *syncData == "1" || *syncData == "true",
			OwnerCreateIfNotExist: *syncOwnerCreateIfNotExist == "1" || *syncOwnerCreateIfNotExist == "true",
			FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
				return tableName + "_" + index.IndexName
			},
			ErrorContinue: true,
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
	println("同步成功")
}

func getSyncOwners(ownerInfoStr string) (owners []*worker.TaskSyncOwner) {
	ownerStrList := strings.Split(ownerInfoStr, ",")
	for _, ownerStr := range ownerStrList {
		ss := strings.Split(ownerStr, "=")
		if len(ss) > 1 {
			owners = append(owners, &worker.TaskSyncOwner{
				SourceName: strings.TrimSpace(ss[0]),
				TargetName: strings.TrimSpace(ss[1]),
			})
		} else if len(ss) > 0 {
			owners = append(owners, &worker.TaskSyncOwner{
				SourceName: strings.TrimSpace(ss[0]),
			})
		}
	}
	return
}
