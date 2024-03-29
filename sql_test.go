package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"github.com/team-ide/go-driver/db_dm"
	"github.com/team-ide/go-driver/db_gbase"
	"github.com/team-ide/go-driver/db_kingbase_v8r3"
	"github.com/team-ide/go-driver/db_kingbase_v8r6"
	"github.com/team-ide/go-driver/db_mysql"
	"github.com/team-ide/go-driver/db_oracle"
	"github.com/team-ide/go-driver/db_shentong"
	"github.com/team-ide/go-driver/db_sqlite3"
	"strings"
	"testing"
)

var testDialectList []*testDialect
var testDialectCache = make(map[string]*testDialect)

type testDialect struct {
	table       *dialect.TableModel
	mapping     *dialect.SqlMapping
	dialect     dialect.Dialect
	owner       *dialect.OwnerModel
	db          func() (db *sql.DB, err error)
	ownerDb     func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error)
	killSession func(owner *dialect.OwnerModel, db *sql.DB)
}

func (this_ *testDialect) init() {
	this_.table = this_.mapping.GenDemoTable()
	var err error
	this_.dialect, err = dialect.NewMappingDialect(this_.mapping)
	if err != nil {
		panic(err)
	}
}

func init() {
	appendTestDialectMysql()
	appendTestDialectSqlite()
	//appendTestDialectOracle()
	appendTestDialectShenTong()
	appendTestDialectDM()
	appendTestDialectKingBase()
}

func appendTestDialectMysql() {
	one := &testDialect{}
	testDialectCache["mysql"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingMysql()

	one.owner = &dialect.OwnerModel{
		OwnerName:             "TEST_DB",
		OwnerCharacterSetName: "utf8mb4",
	}
	one.db = func() (db *sql.DB, err error) {
		db, err = db_mysql.Open(db_mysql.GetDSN("root", "123456", "127.0.0.1", 3306, ""))
		return
	}
	one.ownerDb = func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error) {
		ownerDb, err = db_mysql.Open(db_mysql.GetDSN("root", "123456", "127.0.0.1", 3306, owner.OwnerName))
		return
	}
	one.init()
}

func appendTestDialectSqlite() {
	one := &testDialect{}
	testDialectCache["sqlite"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingSqlite()

	one.owner = &dialect.OwnerModel{
		OwnerName: "",
	}
	one.db = func() (db *sql.DB, err error) {
		db, err = db_sqlite3.Open(db_sqlite3.GetDSN("temp/sqlite.test.db"))
		return
	}
	one.ownerDb = func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error) {
		ownerDb, err = db_sqlite3.Open(db_sqlite3.GetDSN("temp/sqlite.test.db"))
		return
	}
	one.init()
}

func appendTestDialectOracle() {
	one := &testDialect{}
	testDialectCache["oracle"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingOracle()

	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB",
		OwnerPassword: "123456",
	}
	one.db = func() (db *sql.DB, err error) {
		db, err = db_oracle.Open(db_oracle.GetDSN("root", "123456", "127.0.0.1", 1521, "xe"))
		return
	}
	one.ownerDb = func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error) {
		dsn_ := db_oracle.GetDSN(owner.OwnerName, owner.OwnerPassword, "127.0.0.1", 1521, "xe")
		ownerDb, err = db_oracle.Open(dsn_)
		return
	}
	one.killSession = func(owner *dialect.OwnerModel, db *sql.DB) {

		var err error
		var list []map[string]interface{}
		list, err = worker.DoQuery(db, `SELECT * FROM V$SESSION WHERE USERNAME = '`+owner.OwnerName+`'`, nil)
		if err != nil {
			println(err)
			return
		}

		if len(list) == 0 {
			return
		} else {
			fmt.Println("session list:")
			for _, one := range list {
				bs, _ := json.Marshal(one)
				fmt.Println(string(bs))
				sid := dialect.GetStringValue(one["SID"])
				serial := dialect.GetStringValue(one["SERIAL#"])
				_, _, _, err = worker.DoExecs(db, []string{`ALTER SYSTEM KILL SESSION '` + sid + `,` + serial + `'`}, nil)
				if err != nil {
					println(err)
					continue
				}
			}
		}
	}
	one.init()
}

func appendTestDialectDM() {
	one := &testDialect{}
	testDialectCache["dm"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingDM()

	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB_USER",
		OwnerPassword: "123456789",
	}
	one.db = func() (db *sql.DB, err error) {
		db, err = db_dm.Open(db_dm.GetDSN("SYSDBA", "SYSDBA", "127.0.0.1", 5236))
		return
	}
	one.ownerDb = func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error) {
		ownerDb, err = db_dm.Open(db_dm.GetDSN(owner.OwnerName, owner.OwnerPassword, "127.0.0.1", 5236))
		return
	}
	one.init()
}

func appendTestDialectKingBase() {
	one := &testDialect{}
	testDialectCache["kingbase"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingKingBase()

	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB_USER",
		OwnerPassword: "123456",
	}
	one.db = func() (db *sql.DB, err error) {
		db, err = db_kingbase_v8r3.Open(db_kingbase_v8r3.GetDSN("SYSTEM", "123456", "127.0.0.1", 54321, "TEST"))
		return
	}
	one.ownerDb = func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error) {
		dsn := db_kingbase_v8r3.GetDSN(owner.OwnerName, owner.OwnerPassword, "127.0.0.1", 54321, "TEST")
		dsn += "&search_path=" + owner.OwnerName
		ownerDb, err = db_kingbase_v8r3.Open(dsn)
		return
	}
	one.init()
}

func appendTestDialectShenTong() {
	one := &testDialect{}
	testDialectCache["shentong"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingShenTong()

	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB",
		OwnerPassword: "123456",
	}
	one.db = func() (db *sql.DB, err error) {
		db, err = db_shentong.Open(db_shentong.GetDSN("sysdba", "szoscar55", "127.0.0.1", 2003, "OSRDB"))
		return
	}
	one.ownerDb = func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error) {
		ownerDb, err = db_shentong.Open(db_shentong.GetDSN(owner.OwnerName, owner.OwnerPassword, "127.0.0.1", 2003, "OSRDB"))
		return
	}
	one.init()
}

func TestMysql(t *testing.T) {

	db, err := db_mysql.Open(db_mysql.GetDSN("root", "123456", "127.0.0.1", 3306, ""))
	if err != nil {
		panic(err)
	}
	//list, err := worker.DoQuery(db, `select * from ALL_TableS`)
	//list, err := worker.DoQuery(db, `select * from ALL_OBJECTS`)
	//list, err := worker.DoQuery(db, `select * from SYS_CLASS`)
	list, err := worker.DoQuery(db, `set global show_compatibility_56=on`, nil)

	if err != nil {
		panic(err)
	}
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}
func TestKingBase(t *testing.T) {

	db, err := db_kingbase_v8r3.Open(db_kingbase_v8r3.GetDSN("SYSTEM", "123456", "127.0.0.1", 54321, "TEST"))
	if err != nil {
		panic(err)
	}
	//list, err := worker.DoQuery(db, `select * from ALL_TableS`)
	//list, err := worker.DoQuery(db, `select * from ALL_OBJECTS`)
	//list, err := worker.DoQuery(db, `select * from SYS_CLASS`)
	//list, err := worker.DoQuery(db, `SELECT * FROM information_schema.TABLES`)
	//list, err := worker.DoQuery(db, `SELECT * FROM information_schema.schemata`)

	//list, err := worker.DoQuery(db, `DROP SCHEMA TEST_DB`)
	//owners, err := worker.OwnersSelect(db, testDialectCache["kingbase"].dialect, nil)
	//for _, owner := range owners {
	//	tables, err := worker.TablesSelect(db, testDialectCache["kingbase"].dialect, nil, owner.OwnerName)
	//	if err != nil {
	//		println(err)
	//		//panic(err)
	//		continue
	//	}
	//	for _, one := range tables {
	//		//bs, _ := json.Marshal(one)
	//		//fmt.Println(string(bs))
	//		if strings.EqualFold(one.TableName, "dba_free_space") ||
	//			strings.EqualFold(one.TableName, "SYS_FREESPACES") {
	//			//	strings.EqualFold(one.TableName, "dba_cons_columns") ||
	//			//	strings.EqualFold(one.TableName, "dba_col_privs") ||
	//			//	strings.EqualFold(one.TableName, "dba_col_comments") ||
	//			//	strings.EqualFold(one.TableName, "ALL_VIEWS") ||
	//			//	strings.EqualFold(one.TableName, "all_users") ||
	//			//	strings.EqualFold(one.TableName, "all_triggers") ||
	//			//	strings.EqualFold(one.TableName, "all_trigger_cols") {
	//			continue
	//		}
	//
	//		sqlInfo := `
	//	SELECT * from ` + owner.OwnerName + `.` + one.TableName + `
	//	`
	//		fmt.Println("table:", owner.OwnerName, ".", one.TableName)
	//		list, err := worker.DoQuery(db, sqlInfo)
	//
	//		if err != nil {
	//			println(err)
	//			//panic(err)
	//			continue
	//		}
	//		for _, one := range list {
	//			bs, _ := json.Marshal(one)
	//			str := string(bs)
	//			if strings.Contains(str, "TEST_DB_TABLE_DEMO_col_3") ||
	//				strings.Contains(str, "TEST_DB_TABLE_DEMO_col_4") {
	//				fmt.Println(string(bs))
	//			}
	//		}
	//	}
	//}

	sqlInfo := `
	SELECT
	   *
	FROM ALL_INDEXES
	WHERE TABLE_NAME='TABLE_DEMO'
		`
	list, err := worker.DoQuery(db, sqlInfo, nil)

	if err != nil {
		panic(err)
	}
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}

func TestKingBaseSchema(t *testing.T) {
	schema := "TEST_DB"

	dsn := db_kingbase_v8r6.GetDSN("TEST_DB", "123456", "127.0.0.1", 54321, "TEST")
	dsn += "&search_path=" + schema
	db, err := db_kingbase_v8r6.Open(dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(2)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sqlInfo := `
show superuser_reserved_connections;
		`
	//errsSql, err := worker.DoExec(db, []string{`select * from "TEST_DB"."TABLE_DEMO1"`, sqlInfo})
	//if err != nil {
	//	fmt.Println("errsSql:", errsSql)
	//	panic(err)
	//}
	list, err := worker.DoQuery(db, sqlInfo, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("list:", len(list))
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}

func TestDm(t *testing.T) {
	//var a dm.DmClob
	db, err := db_dm.Open(db_dm.GetDSN("SYSDBA", "SYSDBA", "127.0.0.1", 5236))
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(2)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sqlInfo := `
SELECT 
	*
FROM VRV_JOB.JOB_EXECUTOR_LOG
`
	list, err := worker.DoQuery(db, sqlInfo, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("list:", len(list))
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}
func TestOracle(t *testing.T) {

	db, err := db_oracle.Open(db_oracle.GetDSN("root", "123456", "127.0.0.1", 1521, "xe"))
	if err != nil {
		panic(err)
	}
	//list, err := worker.DoQuery(db, `select * FROM ALL_CONS_COLUMNS`)
	list, err := worker.DoQuery(db, `SELECT * FROM V$SESSION`, nil)
	if err != nil {
		panic(err)
	}
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}
func TestGBase(t *testing.T) {
	dsn := `DRIVER=com.gbasebt.jdbc.Driver;NEWCODESET=UTF8,zh_cn.UTF8,57372;DB_LOCALE=zh_cn.57372;DELIMIDENT=Y;CLIENT_LOCALE=zh_cn.57372;ServerName=gbase01;host=127.0.0.1;service=9088;uid=gbasedbt;pwd=GBase123;DATABASE=VRV_JOB1;`
	db, err := db_gbase.Open(dsn)
	if err != nil {
		panic(err)
	}

	dia, err := dialect.NewDialect("GBase")
	if err != nil {
		panic(err)
	}
	owners, err := worker.OwnersSelect(db, dia, nil)
	if err != nil {
		panic(err)
	}
	for _, one := range owners {
		if !strings.EqualFold(one.OwnerName, "im_dbconfig") {
			continue
		}
		bs, _ := json.Marshal(one)
		fmt.Println("owner:", string(bs))

		tables, err := worker.TablesDetail(db, dia, nil, one.OwnerName, false)
		if err != nil {
			panic(err)
		}
		for _, one := range tables {
			bs, _ := json.Marshal(one)
			fmt.Println("table:", string(bs))
		}
	}

}
func TestAllTableSql(t *testing.T) {
	for _, one := range testDialectList {
		sqlList, err := one.dialect.TableCreateSql(nil, one.owner.OwnerName, one.table)
		if err != nil {
			panic(err)
		}
		fmt.Println("-----dialect [" + one.dialect.DialectType().Name + "] create table sql---")
		for _, sqlOne := range sqlList {
			fmt.Println(sqlOne, ";")
		}

		for _, to := range testDialectList {
			if to == one {
				continue
			}
			toTableSql(one, to)
		}
	}
}

func toTableSql(from *testDialect, to *testDialect) {
	sqlList, err := to.dialect.TableCreateSql(nil, to.owner.OwnerName, from.table)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create table sql---")
	for _, sqlOne := range sqlList {
		fmt.Println(sqlOne, ";")
	}
}

func TestToTableSql(t *testing.T) {
	toTableSql(testDialectCache["oracle"], testDialectCache["mysql"])
}

func TestAllSql(t *testing.T) {
	param := &dialect.ParamModel{}

	var err error
	var bs []byte
	//var aa godror.Number
	//aa.String()
	for _, from := range testDialectList {
		fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] create table---")

		bs, err = json.Marshal(from.table)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bs))

		db, err := from.db()
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		if from.owner.OwnerName != "" {
			fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] owner [" + from.owner.OwnerName + "] cover---")
			_, err = worker.OwnerCover(db, from.dialect, param, from.owner)
			if err != nil {
				panic(err)
			}
			fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] owner [" + from.owner.OwnerName + "] success---")
		}
		dialectDb, err := from.ownerDb(from.owner)
		if err != nil {
			panic(err)
		}
		err = dialectDb.Ping()
		if err != nil {
			panic(err)
		}
		err = worker.TableCover(dialectDb, from.dialect, param, from.owner.OwnerName, from.table)
		if err != nil {
			panic(err)
		}
		table, err := worker.TableDetail(db, from.dialect, param, from.owner.OwnerName, from.table.TableName, false)
		if err != nil {
			panic(err)
		}
		if table == nil {
			panic("dialect [" + from.dialect.DialectType().Name + "]  ownerName [" + from.owner.OwnerName + "] tableName [" + from.table.TableName + "] is null.")
		}

		var dataList []map[string]interface{}

		//var a dm.DmBlob
		var data = make(map[string]interface{})
		for i, column := range from.table.ColumnList {
			info, _ := from.dialect.GetColumnTypeInfo(column)

			if info.IsNumber {
				data[column.ColumnName] = 1
			} else if info.IsString {
				data[column.ColumnName] = "s" + fmt.Sprint(i)
			} else if info.IsBytes {
				if info.Name != "BFILE" {
					data[column.ColumnName] = []byte{1}
				}
			} else if info.IsBoolean {
				data[column.ColumnName] = true
			} else if info.IsDateTime {
				//data[column.ColumnName] = time.Now().Local()
				//if info.Name == "YEAR" {
				//	data[column.ColumnName] = 22
				//}
			}

		}
		dataList = append(dataList, data)

		_, _, batchSqlList, batchValuesList, err := from.dialect.DataListInsertSql(nil, from.owner.OwnerName, from.table.TableName, from.table.ColumnList, dataList)
		if err != nil {
			panic(err)
		}
		_, errSql, errArgs, err := worker.DoExecs(dialectDb, batchSqlList, batchValuesList)
		if err != nil {
			fmt.Println("error sql :", errSql)
			fmt.Println("error args:", errArgs)
			panic(err)
		}

		selectSql := "select * from "
		selectSql += from.dialect.OwnerTablePack(nil, from.owner.OwnerName, from.table.TableName)
		list, err := worker.DoQuery(dialectDb, selectSql, nil)
		if err != nil {
			panic(err)
		}
		for _, one := range dataList {
			bs, _ := json.Marshal(one)
			fmt.Println("insert data:", string(bs))
		}
		for _, one := range list {
			bs, _ := json.Marshal(one)
			fmt.Println(string(bs))
			fmt.Println("select data:", string(bs))
		}
		_ = dialectDb.Close()
		if from.killSession != nil {
			from.killSession(from.owner, db)
		}
		_ = db.Close()

		//for columnIndex, fromColumn := range from.table.ColumnList {
		//	savedColumn := table.ColumnList[columnIndex]
		//	if fromColumn.ColumnName != savedColumn.ColumnName ||
		//		//fromColumn.ColumnDataType != savedColumn.ColumnDataType ||
		//		fromColumn.ColumnLength != savedColumn.ColumnLength ||
		//		fromColumn.ColumnPrecision != savedColumn.ColumnPrecision ||
		//		fromColumn.ColumnScale != savedColumn.ColumnScale {
		//		fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] table column not eq---")
		//		bs, err = json.Marshal(fromColumn)
		//		if err != nil {
		//			panic(err)
		//		}
		//		fmt.Println("fromColumn:", string(bs))
		//		bs, err = json.Marshal(savedColumn)
		//		if err != nil {
		//			panic(err)
		//		}
		//		fmt.Println("savedColumn:", string(bs))
		//		if fromColumn.ColumnLength == 0 &&
		//			fromColumn.ColumnPrecision == 0 &&
		//			fromColumn.ColumnScale == 0 {
		//			continue
		//		} else if fromColumn.ColumnLength != savedColumn.ColumnLength &&
		//			fromColumn.ColumnPrecision == savedColumn.ColumnPrecision &&
		//			fromColumn.ColumnScale == savedColumn.ColumnScale {
		//			continue
		//		}
		//		//panic("字段不一致")
		//	}
		//}

		//bs, err = json.Marshal(table)
		//if err != nil {
		//	panic(err)
		//}
		fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] create table result---")
		//fmt.Println(string(bs))

		//for _, to := range testDialectList {
		//if from == to {
		//	continue
		//}
		//	fromTableToTableSql(from, table, to)
		//}
	}
}

func fromTableToTableSql(from *testDialect, fromTable *dialect.TableModel, to *testDialect) {
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create table---")

	for _, column := range fromTable.ColumnList {
		column.ColumnDefault = ""
		if to.dialect.DialectType() == dialect.TypeMysql {
			info, _ := to.dialect.GetColumnTypeInfo(column)
			if info != nil {
				if info.Name == "TIMESTAMP" {
					column.ColumnDefault = "current_timestamp"
				}
			}
		}
	}
	bs, err := json.Marshal(fromTable)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create from table---")
	fmt.Println(string(bs))

	db, err := to.db()
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	var dialectDb *sql.DB
	param := &dialect.ParamModel{}
	if to.owner.OwnerName != "" {
		_, err = worker.OwnerCover(db, to.dialect, param, to.owner)
		if err != nil {
			panic(err)
		}
	}
	dialectDb, err = to.ownerDb(to.owner)
	if err != nil {
		panic(err)
	}
	err = dialectDb.Ping()
	if err != nil {
		panic(err)
	}
	err = worker.TableCover(dialectDb, to.dialect, param, to.owner.OwnerName, fromTable)
	if err != nil {
		panic(err)
	}
	table, err := worker.TableDetail(db, to.dialect, param, to.owner.OwnerName, fromTable.TableName, false)
	if err != nil {
		panic(err)
	}
	if table == nil {
		panic("dialect [" + from.dialect.DialectType().Name + "]  ownerName [" + from.owner.OwnerName + "] tableName [" + from.table.TableName + "] is null.")
	}
	_ = dialectDb.Close()
	if from.killSession != nil {
		from.killSession(from.owner, db)
	}
	_ = db.Close()
	bs, err = json.Marshal(table)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create to table---")
	fmt.Println(string(bs))
}
