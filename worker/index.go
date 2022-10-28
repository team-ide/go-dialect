package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
)

func PrimaryKeysSelect(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) (list []*dialect.PrimaryKeyModel, err error) {
	sqlInfo, err := dia.PrimaryKeysSelectSql(ownerName, tableName)
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

func IndexesSelect(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) (list []*dialect.IndexModel, err error) {
	sqlInfo, err := dia.IndexesSelectSql(ownerName, tableName)
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
