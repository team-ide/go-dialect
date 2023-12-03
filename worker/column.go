package worker

import (
	"database/sql"
	"errors"
	"github.com/team-ide/go-dialect/dialect"
	"strings"
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
		errStr := err.Error()
		if dia.DialectType() == dialect.TypeMysql && strings.Contains(errStr, "Unknown column 'DATETIME_PRECISION'") {
			sqlInfo = `SELECT
    COLUMN_NAME columnName,
    COLUMN_COMMENT columnComment,
    COLUMN_DEFAULT columnDefault,
    EXTRA columnExtra,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    CHARACTER_SET_NAME columnCharacterSetName,
    IS_NULLABLE isNullable,
    DATA_TYPE columnDataType,
    COLUMN_TYPE columnType,
    NUMERIC_PRECISION NUMERIC_PRECISION,
    NUMERIC_SCALE NUMERIC_SCALE,
    CHARACTER_MAXIMUM_LENGTH CHARACTER_MAXIMUM_LENGTH
FROM information_schema.columns
WHERE TABLE_SCHEMA='` + ownerName + `'
  AND TABLE_NAME='` + tableName + `'`
			dataList, err = DoQuery(db, sqlInfo, nil)
		}
	}
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
			column.ColumnAfterColumn = last.ColumnName
		}
		last = column
	}
	return
}
