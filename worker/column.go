package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
)

func ColumnsSelect(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) (list []*dialect.ColumnModel, err error) {
	sqlInfo, err := dia.ColumnsSelectSql(ownerName, tableName)
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
