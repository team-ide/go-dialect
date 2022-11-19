package worker

import (
	"database/sql"
	"errors"
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
		err = errors.New("OwnersSelect error sql:" + sqlInfo + ",error:" + err.Error())
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
		err = errors.New("OwnerSelect error sql:" + sqlInfo + ",error:" + err.Error())
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
	_, errorSql, _, err := DoExecs(db, sqlList)
	if err != nil {
		err = errors.New("OwnerCreate error sql:" + errorSql + ",error:" + err.Error())
		return
	}
	created = true
	return
}

func OwnerDelete(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, ownerName string) (deleted bool, err error) {
	sqlList, err := dia.OwnerDeleteSql(param, ownerName)
	if err != nil {
		return
	}
	if len(sqlList) == 0 {
		return
	}
	_, errorSql, _, err := DoExecs(db, sqlList)
	if err != nil {
		err = errors.New("OwnerDelete error sql:" + errorSql + ",error:" + err.Error())
		return
	}
	deleted = true
	return
}

// OwnerCover 库或表所属者 覆盖，如果 库或表所属者 已经存在，则删除后 再创建
func OwnerCover(db *sql.DB, dia dialect.Dialect, param *dialect.ParamModel, owner *dialect.OwnerModel) (success bool, err error) {
	if owner.OwnerName == "" {
		return
	}
	find, err := OwnerSelect(db, dia, param, owner.OwnerName)
	if err != nil {
		return
	}
	if find != nil {
		_, err = OwnerDelete(db, dia, param, owner.OwnerName)
		if err != nil {
			return
		}
	}
	success, err = OwnerCreate(db, dia, param, owner)
	if err != nil {
		return
	}
	return
}
