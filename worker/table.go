package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
)

func TablesSelect(db *sql.DB, dia dialect.Dialect, databaseName string) (list []*dialect.TableModel, err error) {
	sqlInfo, err := dia.TablesSelectSql(databaseName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo)
	if err != nil {
		return
	}
	for _, data := range dataList {
		model, e := dia.TableModel(data)
		if e != nil {
			model = &dialect.TableModel{
				Error: e.Error(),
			}
		}

		list = append(list, model)
	}
	return
}

func TableDetail(db *sql.DB, dia dialect.Dialect, databaseName string, tableName string) (table *dialect.TableModel, err error) {
	sqlInfo, err := dia.TableSelectSql(databaseName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo)
	if err != nil {
		return
	}
	if len(dataList) > 0 {
		model, e := dia.TableModel(dataList[0])
		if e != nil {
			model = &dialect.TableModel{
				Error: e.Error(),
			}
		} else {
			model.ColumnList, e = ColumnsSelect(db, dia, databaseName, model.Name)
			if e != nil {
				model.Error = e.Error()
			} else {
				ps, e := PrimaryKeysSelect(db, dia, databaseName, model.Name)
				if e != nil {
					model.Error = e.Error()
				} else {
					model.AddPrimaryKey(ps...)
					is, e := IndexesSelect(db, dia, databaseName, model.Name)
					if e != nil {
						model.Error = e.Error()
					} else {
						model.AddIndex(is...)
					}
				}
			}
		}
		table = model
	}

	return
}

func ColumnsSelect(db *sql.DB, dia dialect.Dialect, databaseName string, tableName string) (list []*dialect.ColumnModel, err error) {
	sqlInfo, err := dia.ColumnsSelectSql(databaseName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo)
	if err != nil {
		return
	}
	for _, data := range dataList {
		model, e := dia.ColumnModel(data)
		if e != nil {
			model = &dialect.ColumnModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
	}
	return
}

func PrimaryKeysSelect(db *sql.DB, dia dialect.Dialect, databaseName string, tableName string) (list []*dialect.PrimaryKeyModel, err error) {
	sqlInfo, err := dia.PrimaryKeysSelectSql(databaseName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo)
	if err != nil {
		return
	}
	for _, data := range dataList {
		model, e := dia.PrimaryKeyModel(data)
		if e != nil {
			model = &dialect.PrimaryKeyModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
	}
	return
}

func IndexesSelect(db *sql.DB, dia dialect.Dialect, databaseName string, tableName string) (list []*dialect.IndexModel, err error) {
	sqlInfo, err := dia.IndexesSelectSql(databaseName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo)
	if err != nil {
		return
	}
	for _, data := range dataList {
		model, e := dia.IndexModel(data)
		if e != nil {
			model = &dialect.IndexModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
	}
	return
}
