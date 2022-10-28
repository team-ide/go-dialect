package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
)

func TablesSelect(db *sql.DB, dia dialect.Dialect, ownerName string) (list []*dialect.TableModel, err error) {
	sqlInfo, err := dia.TablesSelectSql(ownerName)
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

func TableDetail(db *sql.DB, dia dialect.Dialect, ownerName string, tableName string) (table *dialect.TableModel, err error) {
	sqlInfo, err := dia.TableSelectSql(ownerName, tableName)
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
			model.ColumnList, e = ColumnsSelect(db, dia, ownerName, model.Name)
			if e != nil {
				model.Error = e.Error()
			} else {
				ps, e := PrimaryKeysSelect(db, dia, ownerName, model.Name)
				if e != nil {
					model.Error = e.Error()
				} else {
					model.AddPrimaryKey(ps...)
					is, e := IndexesSelect(db, dia, ownerName, model.Name)
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
