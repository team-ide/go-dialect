package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
)

func OwnersSelect(db *sql.DB, dia dialect.Dialect) (list []*dialect.OwnerModel, err error) {
	sqlInfo, err := dia.OwnersSelectSql()
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
