package worker

import (
	"database/sql"
	"errors"
	"github.com/team-ide/go-dialect/dialect"
)

func ColumnsSelect(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, tableName string, ignoreError bool) (list []*dialect.ColumnModel, err error) {
	sqlInfo, err := dia.ColumnsSelectSql(param, ownerName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQuery(db, sqlInfo, nil)
	if err != nil {
		err = errors.New("ColumnsSelect error sql:" + sqlInfo + ",error:" + err.Error())
		return
	}
	for _, data := range dataList {
		model, e := dia.ColumnModel(data)
		if e != nil {
			if !ignoreError {
				err = e
				return
			}
			model = &dialect.ColumnModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
	}
	var last *dialect.ColumnModel
	for _, column := range list {
		if last != nil {
			column.ColumnAfterColumn = column.ColumnName
		}
		last = column
	}
	return
}
