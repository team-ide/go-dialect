package dialect

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (this_ *mappingDialect) OwnersSelectSql(param *ParamModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.OwnersSelect, param)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) OwnerSelectSql(param *ParamModel, ownerName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.OwnerSelect, param,
		map[string]string{
			"ownerName": ownerName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) OwnerModel(data map[string]interface{}) (owner *OwnerModel, err error) {
	if data == nil {
		return
	}
	owner = &OwnerModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) OwnerCreateSql(param *ParamModel, owner *OwnerModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.OwnerCreate, param, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) OwnerDeleteSql(param *ParamModel, ownerName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.OwnerDelete, param,
		map[string]string{
			"ownerName": ownerName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TablesSelectSql(param *ParamModel, ownerName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.TablesSelect, param,
		map[string]string{
			"ownerName": ownerName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) TableSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.TableSelect, param,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	if data == nil {
		return
	}
	table = &TableModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, table)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TableCreateSql(param *ParamModel, ownerName string, table *TableModel) (sqlList []string, err error) {
	var tableCreateColumnContent string

	var tableCreateColumnSql string
	for i, column := range table.ColumnList {
		tableCreateColumnSql, err = this_.TableCreateColumnSql(param, column)
		if err != nil {
			return
		}
		tableCreateColumnContent += "    " + tableCreateColumnSql
		if i < len(table.ColumnList)-1 {
			tableCreateColumnContent += ",\n"
		}
	}
	var tableCreatePrimaryKeyContent string
	if len(table.PrimaryKeys) > 0 {
		tableCreatePrimaryKeyContent, err = this_.TableCreatePrimaryKeySql(param, table.PrimaryKeys)
		if len(strings.TrimSpace(tableCreatePrimaryKeyContent)) > 0 {
			tableCreatePrimaryKeyContent = "    " + tableCreatePrimaryKeyContent
			if len(strings.TrimSpace(tableCreateColumnContent)) > 0 {
				tableCreateColumnContent += ","
			}
		}
	}
	if err != nil {
		return
	}
	sqlList, err = this_.FormatSql(this_.TableCreate, param,
		table,
		map[string]string{
			"ownerName":                    ownerName,
			"tableCreateColumnContent":     tableCreateColumnContent,
			"tableCreatePrimaryKeyContent": tableCreatePrimaryKeyContent,
		},
	)
	if err != nil {
		return
	}
	var indexAddSqlList []string
	for _, index := range table.IndexList {
		indexAddSqlList, err = this_.IndexAddSql(param, ownerName, table.TableName, index)
		if err != nil {
			return
		}
		sqlList = append(sqlList, indexAddSqlList...)
	}
	return
}

func (this_ *mappingDialect) TableCreateColumnSql(param *ParamModel, column *ColumnModel) (sqlInfo string, err error) {
	columnTypePack, err := this_.ColumnTypePack(column)
	if err != nil {
		return
	}
	columnDefaultPack, err := this_.ColumnDefaultPack(param, column)
	if err != nil {
		return
	}
	sqlList, err := this_.FormatSql(this_.TableCreateColumn, param,
		column,
		map[string]string{
			"columnTypePack":    columnTypePack,
			"columnDefaultPack": columnDefaultPack,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) TableCreatePrimaryKeySql(param *ParamModel, primaryKeys []string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.TableCreatePrimaryKey, param,
		map[string]interface{}{
			"primaryKeys": primaryKeys,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}
func (this_ *mappingDialect) TableRenameSql(param *ParamModel, ownerName string, tableName string, newTableName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableRename, param,
		map[string]string{
			"ownerName":    ownerName,
			"tableName":    tableName,
			"newTableName": newTableName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TableDeleteSql(param *ParamModel, ownerName string, tableName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableDelete, param,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TableCommentSql(param *ParamModel, ownerName string, tableName string, tableComment string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableComment, param,
		map[string]string{
			"ownerName":    ownerName,
			"tableName":    tableName,
			"tableComment": tableComment,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnsSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.ColumnsSelect, param,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) ColumnSelectSql(param *ParamModel, ownerName string, tableName string, columnName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.ColumnSelect, param,
		map[string]string{
			"ownerName":  ownerName,
			"tableName":  tableName,
			"columnName": columnName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) ColumnModel(data map[string]interface{}) (column *ColumnModel, err error) {
	if data == nil {
		return
	}
	column = &ColumnModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, column)
	if err != nil {
		return
	}

	if data["isNullable"] != nil {
		isNullable, ok := data["isNullable"].(string)
		if ok {
			if strings.EqualFold(isNullable, "no") {
				column.ColumnNotNull = true
			}
		}
	}
	if GetStringValue(data["isNotNull"]) == "1" {
		column.ColumnNotNull = true
	}
	var columnType string
	if data["columnType"] != nil {
		columnType = data["columnType"].(string)
	}
	if column.ColumnDataType == "" {
		if strings.Contains(columnType, "(") {
			column.ColumnDataType = columnType[:strings.Index(columnType, "(")]
		} else {
			column.ColumnDataType = columnType
		}
	}

	columnTypeInfo, err := this_.GetColumnTypeInfo(column.ColumnDataType)
	if err != nil {
		//fmt.Println(data)
		return
	}
	if columnTypeInfo.FullColumnByColumnType != nil {
		err = columnTypeInfo.FullColumnByColumnType(columnType, column)
		if err != nil {
			return
		}
	} else {
		if strings.Contains(columnType, "(") {
			lengthStr := columnType[strings.Index(columnType, "(")+1 : strings.Index(columnType, ")")]
			if strings.Contains(lengthStr, ",") {
				column.ColumnLength, _ = strconv.Atoi(lengthStr[0:strings.Index(lengthStr, ",")])
				column.ColumnDecimal, _ = strconv.Atoi(lengthStr[strings.Index(lengthStr, ",")+1:])
			} else {
				column.ColumnLength, _ = strconv.Atoi(lengthStr)
			}
		}
	}
	return
}

func (this_ *mappingDialect) ColumnAddSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {

	columnTypePack, err := this_.ColumnTypePack(column)
	if err != nil {
		return
	}
	columnDefaultPack, err := this_.ColumnDefaultPack(param, column)
	if err != nil {
		return
	}
	sqlList, err = this_.FormatSql(this_.ColumnAdd, param,
		column,
		map[string]string{
			"ownerName":         ownerName,
			"tableName":         tableName,
			"columnTypePack":    columnTypePack,
			"columnDefaultPack": columnDefaultPack,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnUpdateSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel, newColumn *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnUpdate, param,
		column,
		newColumn,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnDeleteSql(param *ParamModel, ownerName string, tableName string, columnName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnDelete, param,
		map[string]string{
			"ownerName":  ownerName,
			"tableName":  tableName,
			"columnName": columnName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnRenameSql(param *ParamModel, ownerName string, tableName string, columnName string, newColumnName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnRename, param,
		map[string]string{
			"ownerName":     ownerName,
			"tableName":     tableName,
			"columnName":    columnName,
			"newColumnName": newColumnName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnCommentSql(param *ParamModel, ownerName string, tableName string, columnName string, columnComment string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnComment, param,
		map[string]string{
			"ownerName":     ownerName,
			"tableName":     tableName,
			"columnName":    columnName,
			"columnComment": columnComment,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) PrimaryKeysSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.PrimaryKeysSelect, param,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
	if data == nil {
		return
	}
	primaryKey = &PrimaryKeyModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, primaryKey)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) PrimaryKeyAddSql(param *ParamModel, ownerName string, tableName string, columnNames []string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.PrimaryKeyAdd, param,
		map[string]interface{}{
			"ownerName":   ownerName,
			"tableName":   tableName,
			"columnNames": columnNames,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) PrimaryKeyDeleteSql(param *ParamModel, ownerName string, tableName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.PrimaryKeyDelete, param,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) IndexesSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.IndexesSelect, param,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
	if data == nil {
		return
	}
	index = &IndexModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, index)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) IndexAddSql(param *ParamModel, ownerName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.IndexAdd, param,
		index,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
		},
	)
	if err != nil {
		return
	}
	fmt.Println("index add sql:", sqlList)
	return
}

func (this_ *mappingDialect) IndexDeleteSql(param *ParamModel, ownerName string, tableName string, indexName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.IndexDelete, param,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
			"indexName": indexName,
		},
	)
	if err != nil {
		return
	}
	return
}
