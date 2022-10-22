package go_dialect

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"github.com/team-ide/go-dialect/dialect"
)

func init() {
	zorm.FuncPrintSQL = func(ctx context.Context, sqlstr string, args []interface{}, execSQLMillis int64) {

	}
	zorm.FuncLogError = func(ctx context.Context, err error) {

	}
	zorm.FuncLogPanic = func(ctx context.Context, err error) {

	}
}
func getTable() (table *dialect.TableModel) {
	table = &dialect.TableModel{
		Name:    "USER_INFO",
		Comment: "用户信息",
		ColumnList: []*dialect.ColumnModel{
			{Name: "userId", Type: "bigint", Length: 20, PrimaryKey: true},
			{Name: "name", Type: "varchar", Length: 200},
			{Name: "account", Type: "varchar", Length: 50},
		},
		IndexList: []*dialect.IndexModel{
			{Name: "account", Type: "UNIQUE", Columns: []string{"account"}},
		},
	}
	return
}
func testDatabaseCreate(dbContext context.Context, dialect2 dialect.Dialect, param *dialect.GenerateParam, database *dialect.DatabaseModel) {
	sqlList, err := dialect2.DatabaseCreateSql(param, database)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + database.Name + "] create--------")
	testUpdate(dbContext, sqlList)
	fmt.Println()
	fmt.Println()
}
func testDatabaseDelete(dbContext context.Context, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string) {
	sqlList, err := dialect2.DatabaseDeleteSql(param, databaseName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] delete--------")
	testUpdate(dbContext, sqlList)
	fmt.Println()
	fmt.Println()
}
func testTableCreate(dbContext context.Context, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string, table *dialect.TableModel) {
	sqlList, err := dialect2.TableCreateSql(param, databaseName, table)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] table [" + table.Name + "] create--------")
	testUpdate(dbContext, sqlList)
	fmt.Println()
	fmt.Println()

}
func testTableDelete(dbContext context.Context, dialect2 dialect.Dialect, param *dialect.GenerateParam, databaseName string, tableName string) {
	sqlList, err := dialect2.TableDeleteSql(param, databaseName, tableName)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] delete--------")
	testUpdate(dbContext, sqlList)
	fmt.Println()
	fmt.Println()

}

func testUpdate(dbContext context.Context, sqlList []string) {

	_, err := zorm.Transaction(dbContext, func(ctx context.Context) (res interface{}, err error) {
		for _, one := range sqlList {
			finder := zorm.NewFinder()
			finder.InjectionCheck = false
			finder.Append(one)

			fmt.Printf("%s\n", one)
			_, err = zorm.UpdateFinder(ctx, finder)
		}
		return
	})
	if err != nil {
		panic(err)
	}

}
func testDatabases(dbContext context.Context, dialect2 dialect.Dialect) {
	sql, err := dialect2.DatabasesSelectSql()
	if err != nil {
		panic(err)
	}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)

	list, err := queryList(dbContext, finder)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------databases--------")
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Printf("data:%s\n", bs)

		model, err := dialect2.DatabaseModel(one)
		if err != nil {
			panic(err)
		}
		bs, _ = json.Marshal(model)
		fmt.Printf("model:%s\n\n\n", bs)
		testTables(dbContext, dialect2, model.Name)
	}

}

func testTables(dbContext context.Context, dialect2 dialect.Dialect, databaseName string) {
	sql, err := dialect2.TablesSelectSql(databaseName)
	if err != nil {
		panic(err)
	}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)

	list, err := queryList(dbContext, finder)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------database [" + databaseName + "] tables--------")
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Printf("data:%s\n", bs)

		model, err := dialect2.TableModel(one)
		if err != nil {
			panic(err)
		}

		cs := testColumns(dbContext, dialect2, databaseName, model.Name)
		model.ColumnList = cs
		fmt.Println()
		fmt.Println()
		pks := testPrimaryKeys(dbContext, dialect2, databaseName, model.Name)
		model.AddPrimaryKey(pks...)
		fmt.Println()
		fmt.Println()
		is := testIndexes(dbContext, dialect2, databaseName, model.Name)
		model.AddIndex(is...)
		fmt.Println()
		fmt.Println()

		bs, _ = json.MarshalIndent(model, "", "  ")
		fmt.Printf("table:%s\n\n\n", bs)
	}

}

func testColumns(dbContext context.Context, dialect2 dialect.Dialect, databaseName string, tableName string) (res []*dialect.ColumnModel) {
	sql, err := dialect2.ColumnsSelectSql(databaseName, tableName)
	if err != nil {
		panic(err)
	}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)

	list, err := queryList(dbContext, finder)
	if err != nil {
		panic(err)
	}
	if len(list) > 0 {
		fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] columns--------")
		for _, one := range list {
			bs, _ := json.Marshal(one)
			fmt.Printf("data:%s\n", bs)

			model, err := dialect2.ColumnModel(one)
			if err != nil {
				panic(err)
			}
			bs, _ = json.Marshal(model)
			fmt.Printf("model:%s\n", bs)
			res = append(res, model)
		}
	}
	return

}

func testPrimaryKeys(dbContext context.Context, dialect2 dialect.Dialect, databaseName string, tableName string) (res []*dialect.PrimaryKeyModel) {
	sql, err := dialect2.PrimaryKeysSelectSql(databaseName, tableName)
	if err != nil {
		panic(err)
	}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)

	list, err := queryList(dbContext, finder)
	if err != nil {
		panic(err)
	}

	if len(list) > 0 {
		fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] primaryKeys--------")
		for _, one := range list {
			bs, _ := json.Marshal(one)
			fmt.Printf("data:%s\n", bs)

			model, err := dialect2.PrimaryKeyModel(one)
			if err != nil {
				panic(err)
			}
			bs, _ = json.Marshal(model)
			fmt.Printf("model:%s\n", bs)
			res = append(res, model)
		}
	}
	return
}

func testIndexes(dbContext context.Context, dialect2 dialect.Dialect, databaseName string, tableName string) (res []*dialect.IndexModel) {
	sql, err := dialect2.IndexesSelectSql(databaseName, tableName)
	if err != nil {
		panic(err)
	}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)

	list, err := queryList(dbContext, finder)
	if err != nil {
		panic(err)
	}

	if len(list) > 0 {
		fmt.Println("--------database [" + databaseName + "] table [" + tableName + "] indexes--------")
		for _, one := range list {
			bs, _ := json.Marshal(one)
			fmt.Printf("data:%s\n", bs)

			model, err := dialect2.IndexModel(one)
			if err != nil {
				panic(err)
			}
			bs, _ = json.Marshal(model)
			fmt.Printf("model:%s\n", bs)
			res = append(res, model)
		}
	}
	return

}

func queryList(dbContext context.Context, finder *zorm.Finder) (list []map[string]interface{}, err error) {

	list, err = zorm.QueryMap(dbContext, finder, nil)
	if err != nil {
		return
	}

	return
}
