package worker

import (
	"database/sql"
	"errors"
	"github.com/team-ide/go-dialect/dialect"
)

func TablesSelect(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string) (list []*dialect.TableModel, err error) {
	sqlInfo, err := dia.TablesSelectSql(param, ownerName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo, nil)
	if err != nil {
		err = errors.New("TablesSelect error sql:" + sqlInfo + ",error:" + err.Error())
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

func TableSelect(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, tableName string, ignoreError bool) (one *dialect.TableModel, err error) {
	sqlInfo, err := dia.TableSelectSql(param, ownerName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo, nil)
	if err != nil {
		err = errors.New("TableSelect error sql:" + sqlInfo + ",error:" + err.Error())
		return
	}
	for _, data := range dataList {
		model, e := dia.TableModel(data)
		if e != nil {
			if !ignoreError {
				err = e
				return
			}
			model = &dialect.TableModel{
				Error: e.Error(),
			}
		}
		one = model
		return
	}
	return
}

func TablesDetail(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, ignoreError bool) (list []*dialect.TableModel, err error) {
	sqlInfo, err := dia.TablesSelectSql(param, ownerName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo, nil)
	if err != nil {
		err = errors.New("TablesSelect error sql:" + sqlInfo + ",error:" + err.Error())
		return
	}
	for _, data := range dataList {
		model, e := dia.TableModel(data)
		if e != nil {
			if !ignoreError {
				err = e
				return
			}
			model = &dialect.TableModel{
				Error: e.Error(),
			}
		}

		err = appendTableDetail(db, dia, param, ownerName, model, ignoreError)
		if err != nil {
			return
		}
		list = append(list, model)
	}
	return
}
func TableDetail(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, tableName string, ignoreError bool) (table *dialect.TableModel, err error) {
	sqlInfo, err := dia.TableSelectSql(param, ownerName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo, nil)
	if err != nil {
		err = errors.New("TableDetail error sql:" + sqlInfo + ",error:" + err.Error())
		return
	}
	if len(dataList) > 0 {
		model, e := dia.TableModel(dataList[0])
		if e != nil {
			if !ignoreError {
				err = e
				return
			}
			model = &dialect.TableModel{
				Error: e.Error(),
			}
		}
		table = model
		err = appendTableDetail(db, dia, param, ownerName, table, ignoreError)
	}
	return
}

func appendTableDetail(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, table *dialect.TableModel, ignoreError bool) (err error) {

	var e error
	table.ColumnList, e = ColumnsSelect(db, dia, param, ownerName, table.TableName, ignoreError)
	if e != nil {
		if !ignoreError {
			err = e
			return
		}
		table.Error += e.Error()
	} else {
		ps, e := PrimaryKeysSelect(db, dia, param, ownerName, table.TableName, ignoreError)
		if e != nil {
			if !ignoreError {
				err = e
				return
			}
			table.Error += e.Error()
		} else {
			table.AddPrimaryKey(ps...)
			is, e := IndexesSelect(db, dia, param, ownerName, table.TableName, ignoreError)
			if e != nil {
				if !ignoreError {
					err = e
					return
				}
				table.Error += e.Error()
			} else {
				table.AddIndex(is...)
			}
		}
	}

	return
}

func TableCreate(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, tableDetail *dialect.TableModel) (err error) {
	sqlList, err := dia.TableCreateSql(param, ownerName, tableDetail)
	if err != nil {
		return
	}
	if len(sqlList) == 0 {
		return
	}
	_, errorSql, _, err := DoExecs(db, sqlList, nil)
	if err != nil {
		err = errors.New("TableCreate error sql:" + errorSql + ",error:" + err.Error())
		return
	}
	return
}

func TableUpdate(db *sql.DB, oldDia dialect.Dialect, oldTableDetail *dialect.TableModel, newDia dialect.Dialect, newTableDetail *dialect.TableModel) (err error) {

	return
}

func TableDelete(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, tableName string) (err error) {
	sqlList, err := dia.TableDeleteSql(param, ownerName, tableName)
	if err != nil {
		return
	}
	if len(sqlList) == 0 {
		return
	}
	_, errorSql, _, err := DoExecs(db, sqlList, nil)
	if err != nil {
		err = errors.New("TableDelete error sql:" + errorSql + ",error:" + err.Error())
		return
	}
	return
}

// TableCover 表 覆盖，如果 表 已经存在，则删除后 再创建
func TableCover(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, table *dialect.TableModel) (err error) {
	find, err := TableSelect(db, dia, param, ownerName, table.TableName, true)
	if err != nil {
		return
	}
	if find != nil {
		err = TableDelete(db, dia, param, ownerName, table.TableName)
		if err != nil {
			return
		}
	}
	err = TableCreate(db, dia, param, ownerName, table)
	if err != nil {
		return
	}
	return
}
