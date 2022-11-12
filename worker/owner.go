package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
)

func OwnersSelect(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel) (list []*dialect.OwnerModel, err error) {
	sqlInfo, err := dia.OwnersSelectSql(param)
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
		model, e := dia.OwnerModel(data)
		if e != nil {
			model = &dialect.OwnerModel{
				Error: e.Error(),
			}
		}
		list = append(list, model)
	}
	return
}

func OwnerSelect(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string) (one *dialect.OwnerModel, err error) {
	sqlInfo, err := dia.OwnerSelectSql(param, ownerName)
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
		model, e := dia.OwnerModel(data)
		if e != nil {
			model = &dialect.OwnerModel{
				Error: e.Error(),
			}
		}
		one = model
		break
	}
	return
}

func OwnerCreate(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, owner *dialect.OwnerModel) (created bool, err error) {
	sqlList, err := dia.OwnerCreateSql(param, owner)
	if err != nil {
		return
	}
	if len(sqlList) == 0 {
		return
	}
	_, err = DoExec(db, sqlList)
	if err != nil {
		return
	}
	created = true
	return
}
