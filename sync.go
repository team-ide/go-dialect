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
	dia := dialect.GetDialect(*sourceDialect)
	if db == nil || dia == nil {
		panic("sourceDialect [" + *sourceDialect + "] not support")
	}

	targetDb, err := getDbInfo(*targetDialect, *targetUser, *targetPassword, *targetHost, *targetPort, *targetDatabase)
	if err != nil {
		panic(err)
	}
	targetDia := dialect.GetDialect(*targetDialect)
	if targetDb == nil || targetDia == nil {
		panic("targetDialect [" + *targetDialect + "] not support")
	}

	password := *importOwnerCreatePassword
	if password == "" {
		password = *sourcePassword
	}
	var owners = getSyncOwners(*syncOwner)
	task := worker.NewTaskSync(db, dia, targetDb, targetDia,
		func(ownerName string) (db *sql.DB, err error) {
			changeSql, _ := dia.OwnerChangeSql(ownerName)
			if changeSql != "" {
				db, err = getDbInfo(*targetDialect, *targetUser, password, *targetHost, *targetPort, ownerName)
			} else {
				db, err = getDbInfo(*targetDialect, ownerName, password, *targetHost, *targetPort, *targetDatabase)
			}
			return
		},
		&worker.TaskSyncParam{
			Owners:                owners,
			SyncStructure:         true,
			OwnerCreateIfNotExist: *syncOwnerCreateIfNotExist == "1" || *syncOwnerCreateIfNotExist == "true",
			OwnerCreatePassword:   password,
			FormatIndexName: func(ownerName string, tableName string, index *dialect.IndexModel) string {
				return tableName + "_" + index.Name
			},
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
