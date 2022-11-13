package main

import (
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"testing"
)

var testDialectList []*testDialect
var testDialectCache = make(map[string]*testDialect)

type testDialect struct {
	table   *dialect.TableModel
	mapping *dialect.SqlMapping
	dialect dialect.Dialect
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
}

func appendTestDialectMysql() {
	one := &testDialect{}
	testDialectCache["mysql"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingMysql()
	one.init()
}

func appendTestDialectSqlite() {
	one := &testDialect{}
	testDialectCache["sqlite"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingSqlite()
	one.init()
}

func appendTestDialectOracle() {
	one := &testDialect{}
	testDialectCache["oracle"] = one
	testDialectList = append(testDialectList, one)
	one.mapping = dialect.NewMappingOracle()
	one.init()
}

func TestAllTable(t *testing.T) {
	for _, one := range testDialectList {
		sqlList, err := one.dialect.TableCreateSql(nil, "", one.table)
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
			toTable(one, to)
		}
	}
}

func toTable(from *testDialect, to *testDialect) {
	sqlList, err := to.dialect.TableCreateSql(nil, "", from.table)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----dialect [" + from.dialect.DialectType().Name + "] to dialect [" + to.dialect.DialectType().Name + "] create table sql---")
	for _, sqlOne := range sqlList {
		fmt.Println(sqlOne, ";")
	}
}

func TestToTable(t *testing.T) {
	toTable(testDialectCache["oracle"], testDialectCache["mysql"])
}
