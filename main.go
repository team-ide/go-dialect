package main

import (
	"database/sql"
	"flag"
	"github.com/team-ide/go-driver/db_dm"
	"github.com/team-ide/go-driver/db_kingbase_v8r6"
	"github.com/team-ide/go-driver/db_mysql"
	"github.com/team-ide/go-driver/db_oracle"
	"github.com/team-ide/go-driver/db_postgresql"
	"github.com/team-ide/go-driver/db_shentong"
	"github.com/team-ide/go-driver/db_sqlite3"
	"strings"
)

var (
	do       = flag.String("do", "", "操作：export(导出)、import(导入)sync(同步)")
	dbType   = flag.String("t", "", "数据库类型：mysql、sqlite3、dm、kingbase、oracle")
	host     = flag.String("h", "", "数据库Host")
	port     = flag.Int("p", 0, "数据库Port")
	user     = flag.String("u", "", "数据库登录用户")
	password = flag.String("P", "", "数据库登录密码")
	database = flag.String("d", "", "数据库模式名称")
)

func main() {
	flag.Parse()
	//flag.PrintDefaults()
	if *do == "" {
		println("请输入 操作（do）")
		return
	}

	switch strings.ToLower(*do) {
	case "export":
		doExport()
		break
	case "import":
		doImport()
		break
	case "sync":
		doSync()
		break
	default:
		panic("do [" + *do + "] not support")
	}
}

func getDbInfo(dbType string, user string, password string, host string, port int, database string) (db *sql.DB, err error) {
	switch strings.ToLower(dbType) {
	case "mysql":
		dsn := db_mysql.GetDSN(user, password, host, port, database)
		db, err = db_mysql.Open(dsn)
		break
	case "sqlite", "sqlite3":
		dsn := db_sqlite3.GetDSN(database)
		db, err = db_sqlite3.Open(dsn)
		break
	case "damen", "dm":
		dsn := db_dm.GetDSN(user, password, host, port)
		db, err = db_dm.Open(dsn)
		break
	case "kingbase", "kb":
		dsn := db_kingbase_v8r6.GetDSN(user, password, host, port, database)
		db, err = db_kingbase_v8r6.Open(dsn)
		break
	case "oracle":
		dsn := db_oracle.GetDSN(user, password, host, port, database)
		db, err = db_oracle.Open(dsn)
		break
	case "shentong", "st":
		dsn := db_shentong.GetDSN(user, password, host, port, database)
		db, err = db_shentong.Open(dsn)
		break
	case "postgresql", "ps":
		dsn := db_postgresql.GetDSN(user, password, host, port, database)
		db, err = db_postgresql.Open(dsn)
		break
	}
	return
}
