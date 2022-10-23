package go_dialect

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

func testDatabaseCreate(db *sql.DB, dialect2 dialect.Dialect, param *dialect.GenerateParam, database *dialect.DatabaseModel) {
	sqlList, err := dialect2.DatabaseCreateSql(param, database)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + database.Name + "] create--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()
}

func testDatabaseDelete(db *sql.DB, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string) {
	sqlList, err := dialect2.DatabaseDeleteSql(param, databaseName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()
}
func testTableCreate(db *sql.DB, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string, table *dialect.TableModel) {
	sqlList, err := dialect2.TableCreateSql(param, databaseName, table)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] table [" + table.Name + "] create--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testTableDelete(db *sql.DB, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string, tableName string) {
	sqlList, err := dialect2.TableDeleteSql(param, databaseName, tableName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnAdd(db *sql.DB, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string, tableName string, column *dialect.ColumnModel) {
	sqlList, err := dialect2.ColumnAddSql(param, databaseName, tableName, column)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] column [" + column.Name + "] add--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnUpdate(db *sql.DB, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string, tableName string, column *dialect.ColumnModel) {
	sqlList, err := dialect2.ColumnUpdateSql(param, databaseName, tableName, column)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] column [" + column.Name + "] update--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}
func testColumnDelete(db *sql.DB, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string, tableName string, columnName string) {
	sqlList, err := dialect2.ColumnDeleteSql(param, databaseName, tableName, columnName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] column [" + columnName + "] delete--------")
	exec(db, sqlList)
	fmt.Println()
	fmt.Println()

}

func exec(db *sql.DB, sqlList []string) {
	if len(sqlList) == 0 {
		return
	}
	for _, one := range sqlList {

		fmt.Printf("%s\n", one)
		_, err := db.Exec(one)
		if err != nil {
			fmt.Println("error sql:" + one)
			panic(err)
			return
		}

	}

}

func databases(db *sql.DB, dia dialect.Dialect) {
	fmt.Println("--------databases--------")
	list, err := worker.DatabasesSelect(db, dia)
	if err != nil {
		panic(err)
	}
	for _, one := range list {
		if one.Error != "" {
			println("database error:" + one.Error)
			continue
		}

		bs, _ := json.Marshal(one)
		fmt.Printf("%s\n", bs)
		tables(db, dia, one.Name)

	}

}

func tables(db *sql.DB, dia dialect.Dialect, databaseName string) {
	fmt.Println("--------database [" + databaseName + "] tables--------")
	list, err := worker.TablesSelect(db, dia, databaseName)
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
		tableDetail(db, dia, databaseName, one.Name)
	}

}

func tableDetail(db *sql.DB, dia dialect.Dialect, databaseName string, tableName string) {
	fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] detail--------")
	table, err := worker.TableDetail(db, dia, databaseName, tableName)
	if err != nil {
		panic(err)
	}

	bs, _ := json.MarshalIndent(table, "", "  ")
	fmt.Printf("%s\n", bs)

}
