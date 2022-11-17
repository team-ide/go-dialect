package worker

import (
	"context"
	"database/sql"
	"errors"
	"github.com/team-ide/go-dialect/dialect"
)

func PrimaryKeysSelect(db *sql.DB, cxt context.Context, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, tableName string, ignoreError bool) (list []*dialect.PrimaryKeyModel, err error) {
	sqlInfo, err := dia.PrimaryKeysSelectSql(param, ownerName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQueryContext(db, cxt, sqlInfo)
	if err != nil {
		err = errors.New("PrimaryKeysSelect error sql:" + sqlInfo + ",error:" + err.Error())
		return
	}
	for _, data := range dataList {
		model, e := dia.PrimaryKeyModel(data)
		if e != nil {
			if !ignoreError {
				err = e
				return
			}
			model = &dialect.PrimaryKeyModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
	}
	return
}

func IndexesSelect(db *sql.DB, cxt context.Context, dia dialect.Dialect, param *dialect.ParamModel, ownerName string, tableName string, ignoreError bool) (list []*dialect.IndexModel, err error) {
	sqlInfo, err := dia.IndexesSelectSql(param, ownerName, tableName)
	if err != nil {
		return
	}
	if sqlInfo == "" {
		return
	}
	dataList, err := DoQueryContext(db, cxt, sqlInfo)
	if err != nil {
		err = errors.New("IndexesSelect error sql:" + sqlInfo + ",error:" + err.Error())
		return
	}
	for _, data := range dataList {
		model, e := dia.IndexModel(data)
		if e != nil {
			if !ignoreError {
				err = e
				return
			}
			model = &dialect.IndexModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
	}
	return
}
