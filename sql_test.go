package main

import (
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
	db      *sql.DB
	owner   *dialect.OwnerModel
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
	//appendTestDialectSqlite()
	//appendTestDialectOracle()
	//appendTestDialectShenTong()
	//appendTestDialectDM()
	appendTestDialectKinBase()
}

func appendTestDialectMysql() {
	one := &testDialect{}
	testDialectCache["mysql"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingMysql()

	var err error
	one.db, err = db_mysql.Open(db_mysql.GetDSN("root", "123456", "127.0.0.1", 3306, ""))
	if err != nil {
		panic(err)
	}
	one.owner = &dialect.OwnerModel{
		OwnerName:             "TEST_DB",
		OwnerCharacterSetName: "utf8mb4",
	}
	one.init()
}

func appendTestDialectSqlite() {
	one := &testDialect{}
	testDialectCache["sqlite"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingSqlite()

	var err error
	one.db, err = db_sqlite3.Open(db_sqlite3.GetDSN("temp/sqlite.test.db"))
	if err != nil {
		panic(err)
	}
	one.owner = &dialect.OwnerModel{
		OwnerName: "",
	}
	one.init()
}

func appendTestDialectOracle() {
	one := &testDialect{}
	testDialectCache["oracle"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingOracle()

	var err error
	one.db, err = db_oracle.Open(db_oracle.GetDSN("root", "123456", "127.0.0.1", 1521, "xe"))
	if err != nil {
		panic(err)
	}
	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB",
		OwnerPassword: "123456",
	}
	one.init()
}

func appendTestDialectDM() {
	one := &testDialect{}
	testDialectCache["dm"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingDM()

	var err error
	one.db, err = db_dm.Open(db_dm.GetDSN("SYSDBA", "SYSDBA", "127.0.0.1", 5236))
	if err != nil {
		panic(err)
	}
	one.db.SetMaxIdleConns(1)
	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB_USER",
		OwnerPassword: "123456789",
	}
	one.init()
}

func appendTestDialectKinBase() {
	one := &testDialect{}
	testDialectCache["kinbase"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingKinBase()

	var err error
	one.db, err = db_kingbase_v8r3.Open(db_kingbase_v8r3.GetDSN("SYSTEM", "123456", "127.0.0.1", 54321, "TEST"))
	if err != nil {
		panic(err)
	}
	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB_USER2",
		OwnerPassword: "123456",
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
	//tables, err := worker.TablesSelect(db, testDialectCache["kinbase"].dialect, nil, "information_schema")
	//if err != nil {
	//	panic(err)
	//}
	//for _, one := range tables {
	//	bs, _ := json.Marshal(one)
	//	fmt.Println(string(bs))
	//}
	sqlInfo := `
SELECT
    t2.COLUMN_NAME columnName,
    t1.TABLE_NAME tableName,
    t1.TABLE_SCHEMA ownerName,
    t1.CONSTRAINT_NAME indexName,
    t1.CONSTRAINT_TYPE indexType,
    t1.*,
    t2.*
FROM information_schema.table_constraints t1
LEFT JOIN information_schema.key_column_usage t2 
ON (t2.CONSTRAINT_NAME=t1.CONSTRAINT_NAME AND t2.TABLE_SCHEMA=t1.TABLE_SCHEMA AND t2.TABLE_NAME=t1.TABLE_NAME)
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

func appendTestDialectShenTong() {
	one := &testDialect{}
	testDialectCache["shentong"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingShenTong()

	var err error
	one.db, err = db_shentong.Open(db_shentong.GetDSN("sysdba", "szoscar55", "127.0.0.1", 2003, "OSRDB"))
	if err != nil {
		panic(err)
	}
	one.owner = &dialect.OwnerModel{
		OwnerName:     "TEST_DB",
		OwnerPassword: "123456",
	}
	one.init()
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
		if from.owner.OwnerName != "" {
			_, err := worker.OwnerCover(from.db, from.dialect, param, from.owner)
			if err != nil {
				panic(err)
			}
		}
		err := worker.TableCover(from.db, from.dialect, param, from.owner.OwnerName, from.table)
		if err != nil {
			panic(err)
		}
		table, err := worker.TableDetail(from.db, from.dialect, param, from.owner.OwnerName, from.table.TableName, false)
		if err != nil {
			panic(err)
		}
		if table == nil {
			panic("dialect [" + from.dialect.DialectType().Name + "]  ownerName [" + from.owner.OwnerName + "] tableName [" + from.table.TableName + "] is null.")
		}
		bs, err := json.Marshal(table)
		if err != nil {
			panic(err)
		}
		fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] create table success---")
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

	param := &dialect.ParamModel{}
	if to.owner.OwnerName != "" {
		_, err := worker.OwnerCover(to.db, to.dialect, param, to.owner)
		if err != nil {
			panic(err)
		}
	}
	err = worker.TableCover(to.db, to.dialect, param, to.owner.OwnerName, fromTable)
	if err != nil {
		panic(err)
	}
	table, err := worker.TableDetail(to.db, to.dialect, param, to.owner.OwnerName, fromTable.TableName, false)
	if err != nil {
		panic(err)
	}
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
