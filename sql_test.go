package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"github.com/team-ide/go-driver/db_dm"
	"github.com/team-ide/go-driver/db_kingbase_v8r3"
	"github.com/team-ide/go-driver/db_mysql"
	"github.com/team-ide/go-driver/db_oracle"
	"github.com/team-ide/go-driver/db_shentong"
	"github.com/team-ide/go-driver/db_sqlite3"
	"testing"
)

var testDialectList []*testDialect
var testDialectCache = make(map[string]*testDialect)

type testDialect struct {
	table   *dialect.TableModel
	mapping *dialect.SqlMapping
	dialect dialect.Dialect
	owner   *dialect.OwnerModel
	db      func() (db *sql.DB, err error)
	ownerDb func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error)
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
	appendTestDialectOracle()
	appendTestDialectShenTong()
	appendTestDialectDM()
	//appendTestDialectKinBase()
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
		ownerDb, err = db_oracle.Open(db_oracle.GetDSN(owner.OwnerName, owner.OwnerPassword, "127.0.0.1", 1521, "xe"))
		return
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

func appendTestDialectKinBase() {
	one := &testDialect{}
	testDialectCache["kinbase"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingKinBase()

	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB_USER2",
		OwnerPassword: "123456",
	}
	one.db = func() (db *sql.DB, err error) {
		db, err = db_kingbase_v8r3.Open(db_kingbase_v8r3.GetDSN("SYSTEM", "123456", "127.0.0.1", 54321, "TEST"))
		return
	}
	one.ownerDb = func(owner *dialect.OwnerModel) (ownerDb *sql.DB, err error) {
		ownerDb, err = db_kingbase_v8r3.Open(db_kingbase_v8r3.GetDSN(owner.OwnerName, owner.OwnerPassword, "127.0.0.1", 54321, "TEST"))
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
	list, err := worker.DoQuery(db, `set global show_compatibility_56=on`)

	if err != nil {
		panic(err)
	}
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}
func TestKinBase(t *testing.T) {

	db, err := db_kingbase_v8r3.Open(db_kingbase_v8r3.GetDSN("SYSTEM", "123456", "127.0.0.1", 54321, "TEST"))
	if err != nil {
		panic(err)
	}
	//list, err := worker.DoQuery(db, `select * from ALL_TableS`)
	//list, err := worker.DoQuery(db, `select * from ALL_OBJECTS`)
	//list, err := worker.DoQuery(db, `select * from SYS_CLASS`)
	//list, err := worker.DoQuery(db, `SELECT * FROM information_schema.TABLES`)
	//list, err := worker.DoQuery(db, `SELECT * FROM information_schema.schemata`)

	//list, err := worker.DoQuery(db, `DROP SCHEMA TEST_DB_USER2`)
	//owners, err := worker.OwnersSelect(db, testDialectCache["kinbase"].dialect, nil)
	//for _, owner := range owners {
	//	tables, err := worker.TablesSelect(db, testDialectCache["kinbase"].dialect, nil, owner.OwnerName)
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
	//			if strings.Contains(str, "TEST_DB_USER2_TABLE_DEMO_col_3") ||
	//				strings.Contains(str, "TEST_DB_USER2_TABLE_DEMO_col_4") {
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
	list, err := worker.DoQuery(db, sqlInfo)

	if err != nil {
		panic(err)
	}
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}

func TestKinBaseSchema(t *testing.T) {
	schema := "TEST_DB_USER2"

	dsn := db_kingbase_v8r3.GetDSN("TEST_DB_USER2", "123456", "127.0.0.1", 54321, "TEST")
	dsn += "&search_path=" + schema
	db, err := db_kingbase_v8r3.Open(dsn)
	if err != nil {
		panic(err)
	}

	sqlInfo := `
	SELECT
	   *
	FROM TABLE_DEMO
		`
	list, err := worker.DoQuery(db, sqlInfo)

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
	list, err := worker.DoQuery(db, `select * FROM ALL_CONSTRAINTS`)

	if err != nil {
		panic(err)
	}
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
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
	for _, from := range testDialectList {
		fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] create table---")

		bs, err := json.Marshal(from.table)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bs))

		db, err := from.db()
		if err != nil {
			panic(err)
		}
		dbContext := context.Background()
		err = db.PingContext(dbContext)
		if err != nil {
			panic(err)
		}
		if from.owner.OwnerName != "" {
			fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] owner [" + from.owner.OwnerName + "] cover---")
			_, err := worker.OwnerCover(db, dbContext, from.dialect, param, from.owner)
			if err != nil {
				panic(err)
			}
			fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] owner [" + from.owner.OwnerName + "] success---")
		}
		_ = db.Close()
		dialectDb, err := from.ownerDb(from.owner)
		if err != nil {
			panic(err)
		}
		dialectDbContext := context.Background()
		err = dialectDb.PingContext(dialectDbContext)
		if err != nil {
			panic(err)
		}
		err = worker.TableCover(dialectDb, dialectDbContext, from.dialect, param, from.owner.OwnerName, from.table)
		if err != nil {
			panic(err)
		}
		table, err := worker.TableDetail(dialectDb, dialectDbContext, from.dialect, param, from.owner.OwnerName, from.table.TableName, false)
		if err != nil {
			panic(err)
		}
		_ = dialectDb.Close()
		if table == nil {
			panic("dialect [" + from.dialect.DialectType().Name + "]  ownerName [" + from.owner.OwnerName + "] tableName [" + from.table.TableName + "] is null.")
		}
		bs, err = json.Marshal(table)
		if err != nil {
			panic(err)
		}
		fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] create table result---")
		fmt.Println(string(bs))

		for _, to := range testDialectList {
			//if from == to {
			//	continue
			//}
			fromTableToTableSql(from, table, to)
		}
	}
}

func fromTableToTableSql(from *testDialect, fromTable *dialect.TableModel, to *testDialect) {
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create table---")

	bs, err := json.Marshal(fromTable)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create from table---")
	fmt.Println(string(bs))

	toDb, err := to.db()
	if err != nil {
		panic(err)
	}
	toDbContext := context.Background()
	err = toDb.PingContext(toDbContext)
	if err != nil {
		panic(err)
	}

	param := &dialect.ParamModel{}
	if to.owner.OwnerName != "" {
		_, err := worker.OwnerCover(toDb, toDbContext, to.dialect, param, to.owner)
		if err != nil {
			panic(err)
		}
	}
	_ = toDb.Close()
	dialectDb, err := to.ownerDb(to.owner)
	if err != nil {
		panic(err)
	}
	dialectDbContext := context.Background()
	err = dialectDb.PingContext(dialectDbContext)
	if err != nil {
		panic(err)
	}
	err = worker.TableCover(dialectDb, dialectDbContext, to.dialect, param, to.owner.OwnerName, fromTable)
	if err != nil {
		panic(err)
	}
	table, err := worker.TableDetail(dialectDb, dialectDbContext, to.dialect, param, to.owner.OwnerName, fromTable.TableName, false)
	if err != nil {
		panic(err)
	}
	_ = dialectDb.Close()
	if table == nil {
		panic("dialect [" + from.dialect.DialectType().Name + "]  ownerName [" + from.owner.OwnerName + "] tableName [" + from.table.TableName + "] is null.")
	}
	bs, err = json.Marshal(table)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create to table---")
	fmt.Println(string(bs))
}