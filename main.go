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
	do = flag.String("do", "", "操作：export(导出)、import(导入)、sync(同步)")

	sourceDialect  = flag.String("sourceDialect", "", "源 数据库 方言 mysql、sqlite3、damen、kingbase、oracle、shentong")
	sourceHost     = flag.String("sourceHost", "", "源 数据库 host")
	sourcePort     = flag.Int("sourcePort", 0, "源 数据库 port")
	sourceUser     = flag.String("sourceUser", "", "源 数据库 user")
	sourcePassword = flag.String("sourcePassword", "", "源 数据库 password")
	sourceDatabase = flag.String("sourceDatabase", "", "源 数据库 连接库（库名、用户名、SID）")

	fileType    = flag.String("fileType", "", "文件 类型：sql、excel、txt、csv(sql将导出单个文件，其它每个表导出一个文件)")
	fileDialect = flag.String("fileDialect", "", "文件 数据库 方言")
	skipOwner   = flag.String("skipOwner", "", "忽略库名，多个使用“,”隔开")
	skipTable   = flag.String("skipTable", "", "忽略表名，多个使用“,”隔开")

	exportDialect     = flag.String("exportDialect", "", "导出 数据库 方言")
	exportDir         = flag.String("exportDir", "", "导出 文件存储目录")
	exportOwner       = flag.String("exportOwner", "", "导出 库（库名、表拥有者），默认全部，多个使用“,”隔开")
	exportTable       = flag.String("exportTable", "", "导出 表，默认全部，多个使用“,”隔开")
	exportStruct      = flag.String("exportStruct", "", "导出 结构体，默认true，适用于导出类型为sql、excel")
	exportData        = flag.String("exportData", "", "导出 数据，默认true")
	exportAppendOwner = flag.String("exportAppendOwner", "", "sql文件类型的sql拼接 连接库（库名、用户名），拼接原库名或重命名后的库名")

	importOwner                 = flag.String("importOwner", "", "导入 库（库名、表拥有者），并指定文件路径，多个使用“,”隔开")
	importOwnerCreateIfNotExist = flag.String("importOwnerCreateIfNotExist", "", "导入 库如果不存在，则创建")
	importOwnerCreatePassword   = flag.String("importOwnerCreatePassword", "", "导入 库创建的密码，只有库为所属者有效，默认为sourcePassword，如：oracle等数据库")

	targetDialect             = flag.String("targetDialect", "", "目标 数据库 方言 mysql、sqlite3、damen、kingbase、oracle、shentong")
	targetHost                = flag.String("targetHost", "", "目标 数据库 host")
	targetPort                = flag.Int("targetPort", 0, "目标 数据库 port")
	targetUser                = flag.String("targetUser", "", "目标 数据库 user")
	targetPassword            = flag.String("targetPassword", "", "目标 数据库 password")
	targetDatabase            = flag.String("targetDatabase", "", "目标 数据库 连接库（库名、用户名、SID）")
	syncOwner                 = flag.String("syncOwner", "", "同步 库（库名、表拥有者），默认全部，多个使用“,”隔开")
	syncOwnerCreateIfNotExist = flag.String("syncOwnerCreateIfNotExist", "", "同步 库如果不存在，则创建")
	syncOwnerCreatePassword   = flag.String("syncOwnerCreatePassword", "", "同步 库创建的密码，只有库为所属者有效，默认为targetPassword，如：oracle等数据库")
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
