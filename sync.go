package main

import (
	"encoding/json"
	"flag"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"strings"
)

var (
	syncDbType   = flag.String("syncDbType", "", "同步 目标 数据库类型：mysql、sqlite3、dm、kingbase、oracle")
	syncHost     = flag.String("syncHost", "", "同步 目标 数据库Host")
	syncPort     = flag.Int("syncPort", 0, "同步 目标 数据库Port")
	syncUser     = flag.String("syncUser", "", "同步 目标 数据库登录用户")
	syncPassword = flag.String("syncPassword", "", "同步 目标 数据库登录密码")
	syncDatabase = flag.String("syncDatabase", "", "同步 目标 数据库模式名称")
	syncOwner    = flag.String("syncOwner", "", "同步 库或表所属者，多个使用“,”隔开，“x,xx=xx1”")
)

func doSync() {
	if *syncDbType == "" {
		println("请输入 同步 目标 数据库类型")
		return
	}
	if *syncOwner == "" {
		println("请输入 同步 库或表所属者")
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

	targetDb, err := getDbInfo(*syncDbType, *syncUser, *syncPassword, *syncHost, *syncPort, *syncDatabase)
	if err != nil {
		panic(err)
	}
	targetDia := dialect.GetDialect(*syncDbType)
	if targetDb == nil || targetDia == nil {
		panic("syncDbType [" + *syncDbType + "] not support")
	}

	var owners = getSyncOwners(*syncOwner)
	task := worker.NewTaskSync(db, dia, targetDb, targetDia, &worker.TaskSyncParam{
		Owners:        owners,
		SyncStructure: true,
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
