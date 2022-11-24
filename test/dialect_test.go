package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
)

func getTable() (table *dialect.TableModel) {
	table = &dialect.TableModel{
		TableName:    "USER_INFO",
		TableComment: "用户信息",
		ColumnList: []*dialect.ColumnModel{
			{ColumnName: "userId", ColumnType: "bigint", ColumnLength: 20, PrimaryKey: true},
			{ColumnName: "name", ColumnType: "varchar", ColumnLength: 200},
			{ColumnName: "account", ColumnType: "varchar", ColumnLength: 50},
			{ColumnName: "status", ColumnType: "int", ColumnLength: 3},
			{ColumnName: "deleted", ColumnType: "bit", ColumnLength: 1},
			{ColumnName: "detail", ColumnType: "text", ColumnLength: 500},
			{ColumnName: "detail2", ColumnType: "longtext", ColumnLength: 500},
			{ColumnName: "detail3", ColumnType: "blob", ColumnLength: 500},
			{ColumnName: "detail4", ColumnType: "longblob", ColumnLength: 500},
			{ColumnName: "createDate", ColumnType: "date", ColumnLength: 20},
			{ColumnName: "createDate1", ColumnType: "datetime", ColumnLength: 20},
		},
		IndexList: []*dialect.IndexModel{
			{Name: "account", Type: "UNIQUE", Columns: []string{"account"}},
		},
	}
	return
}

func testDLL(db *sql.DB, dia dialect.Dialect, ownerName string) {
	initKingBase()
	table := getTable()
	testTableCreate(db, dia, ownerName, getTable())

	testColumnUpdate(db, dia, ownerName, table.TableName,
		&dialect.ColumnModel{
			ColumnName:    "name",
			ColumnType:    "varchar",
			ColumnLength:  500,
			ColumnComment: "name1注释",
		}, &dialect.ColumnModel{
			ColumnName:    "name1",
			ColumnType:    "varchar",
			ColumnLength:  600,
			ColumnComment: "name1注释",
		},
	)
	testColumnDelete(db, dia, ownerName, table.TableName, "detail3")
	testColumnAdd(db, dia, ownerName, table.TableName, &dialect.ColumnModel{
		ColumnName:    "name2",
		ColumnType:    "varchar",
		ColumnLength:  500,
		ColumnComment: "name2注释",
	})
	tableDetail(db, dia, ownerName, table.TableName)
	testTableDelete(db, dia, ownerName, table.TableName)
}

func testSql(db *sql.DB, dia dialect.Dialect, ownerName, sqlInfo string) {
	sqlList := dia.SqlSplit(sqlInfo)
	exec(db, sqlList)
	tables(db, dia, ownerName)
}

func testOwnerCreate(db *sql.DB, dia dialect.Dialect, owner *dialect.OwnerModel) {
	sqlList, err := dia.OwnerCreateSql(nil, owner)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + owner.OwnerName + "] create--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()
}

func testOwnerDelete(db *sql.DB, dia dialect.Dialect, ownerName string) {
	sqlList, err := dia.OwnerDeleteSql(nil, ownerName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()
}
func testTableCreate(db *sql.DB, dia dialect.Dialect, ownerName string, table *dialect.TableModel) {
	sqlList, err := dia.TableCreateSql(nil, ownerName, table)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + table.TableName + "] create--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testTableDelete(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) {
	sqlList, err := dia.TableDeleteSql(nil, ownerName, tableName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnAdd(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string, column *dialect.ColumnModel) {
	sqlList, err := dia.ColumnAddSql(nil, ownerName, tableName, column)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] column [" + column.ColumnName + "] add--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnUpdate(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string, oldColumn *dialect.ColumnModel, newColumn *dialect.ColumnModel) {
	sqlList, err := dia.ColumnUpdateSql(nil, ownerName, tableName, oldColumn, newColumn)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] column [" + oldColumn.ColumnName + "] update--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnDelete(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string, columnName string) {
	sqlList, err := dia.ColumnDeleteSql(nil, ownerName, tableName, columnName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] column [" + columnName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}

func exec(db *sql.DB, sqlList []string) {
	if len(sqlList) == 0 {
		return
	}
	for _, one := range sqlList {
		if one == "" {
			continue
		}
		fmt.Printf("%s\n", one)
		_, err := db.Exec(one)
		if err != nil {
			fmt.Println("error sql:" + one)
			panic(err)
			return
		}

	}

}

func owners(db *sql.DB, dia dialect.Dialect) {
	fmt.Println("--------owners--------")
	list, err := worker.OwnersSelect(db, dia, nil)
	if err != nil {
		panic(err)
	}
	for _, one := range list {
		if one.Error != "" {
			println("owner error:" + one.Error)
			continue
		}

		bs, _ := json.Marshal(one)
		fmt.Printf("%s\n", bs)
		tables(db, dia, one.OwnerName)

	}

}

func tables(db *sql.DB, dia dialect.Dialect, ownerName string) {
	fmt.Println("--------owner [" + ownerName + "] tables--------")
	list, err := worker.TablesSelect(db, dia, nil, ownerName)
	if err != nil {
		panic(err)
	}
	for _, one := range list {
		if one.Error != "" {
			println("table error:" + one.Error)
			continue
		}

		bs, _ := json.Marshal(one)
		fmt.Printf("%s\n", bs)
		tableDetail(db, dia, ownerName, one.OwnerName)
	}

}

func tableDetail(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) {
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] detail--------")
	table, err := worker.TableDetail(db, dia, nil, ownerName, tableName)
	if err != nil {
		panic(err)
	}

	bs, _ := json.MarshalIndent(table, "", "  ")
	fmt.Printf("%s\n", bs)

}
