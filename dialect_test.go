package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
)

func getTable() (table *dialect.TableModel) {
	table = &dialect.TableModel{
		Name:    "USER_INFO",
		Comment: "用户信息",
		ColumnList: []*dialect.ColumnModel{
			{Name: "userId", Type: "bigint", Length: 20, PrimaryKey: true},
			{Name: "name", Type: "varchar", Length: 200},
			{Name: "account", Type: "varchar", Length: 50},
			{Name: "status", Type: "int", Length: 3},
			{Name: "deleted", Type: "bit", Length: 1},
			{Name: "detail", Type: "text", Length: 500},
			{Name: "detail2", Type: "longtext", Length: 500},
			{Name: "detail3", Type: "blob", Length: 500},
			{Name: "detail4", Type: "longblob", Length: 500},
			{Name: "createDate", Type: "date", Length: 20},
			{Name: "createDate1", Type: "datetime", Length: 20},
		},
		IndexList: []*dialect.IndexModel{
			{Name: "account", Type: "UNIQUE", Columns: []string{"account"}},
		},
	}
	return
}

func testDLL(db *sql.DB, dia dialect.Dialect, ownerName string) {
	initKinBase()
	table := getTable()
	testTableCreate(db, dia, ownerName, getTable())

	testColumnUpdate(db, dia, ownerName, table.Name,
		&dialect.ColumnModel{
			Name:    "name",
			Type:    "varchar",
			Length:  500,
			Comment: "name1注释",
		}, &dialect.ColumnModel{
			Name:    "name1",
			Type:    "varchar",
			Length:  600,
			Comment: "name1注释",
		},
	)
	testColumnDelete(db, dia, ownerName, table.Name, "detail3")
	testColumnAdd(db, dia, ownerName, table.Name, &dialect.ColumnModel{
		Name:    "name2",
		Type:    "varchar",
		Length:  500,
		Comment: "name2注释",
	})
	tableDetail(db, dia, ownerName, table.Name)
	testTableDelete(db, dia, ownerName, table.Name)
}

func testSql(db *sql.DB, dia dialect.Dialect, ownerName, sqlInfo string) {
	sqlList := dia.SqlSplit(sqlInfo)
	exec(db, sqlList)
	tables(db, dia, ownerName)
}

func testOwnerCreate(db *sql.DB, dia dialect.Dialect, owner *dialect.OwnerModel) {
	sqlList, err := dia.OwnerCreateSql(owner)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + owner.Name + "] create--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()
}

func testOwnerDelete(db *sql.DB, dia dialect.Dialect, ownerName string) {
	sqlList, err := dia.OwnerDeleteSql(ownerName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()
}
func testTableCreate(db *sql.DB, dia dialect.Dialect, ownerName string, table *dialect.TableModel) {
	sqlList, err := dia.TableCreateSql(ownerName, table)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + table.Name + "] create--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testTableDelete(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) {
	sqlList, err := dia.TableDeleteSql(ownerName, tableName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnAdd(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string, column *dialect.ColumnModel) {
	sqlList, err := dia.ColumnAddSql(ownerName, tableName, column)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] column [" + column.Name + "] add--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnUpdate(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string, oldColumn *dialect.ColumnModel, newColumn *dialect.ColumnModel) {
	sqlList, err := dia.ColumnUpdateSql(ownerName, tableName, oldColumn, newColumn)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] column [" + oldColumn.Name + "] update--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnDelete(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string, columnName string) {
	sqlList, err := dia.ColumnDeleteSql(ownerName, tableName, columnName)
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
	list, err := worker.OwnersSelect(db, dia)
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
		tables(db, dia, one.Name)

	}

}

func tables(db *sql.DB, dia dialect.Dialect, ownerName string) {
	fmt.Println("--------owner [" + ownerName + "] tables--------")
	list, err := worker.TablesSelect(db, dia, ownerName)
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
		tableDetail(db, dia, ownerName, one.Name)
	}

}

func tableDetail(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) {
	fmt.Println("--------owner [" + ownerName + "] table [" + tableName + "] detail--------")
	table, err := worker.TableDetail(db, dia, ownerName, tableName)
	if err != nil {
		panic(err)
	}

	bs, _ := json.MarshalIndent(table, "", "  ")
	fmt.Printf("%s\n", bs)

}
