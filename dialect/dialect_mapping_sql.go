package dialect

import (
	"encoding/json"
	"errors"
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
	for _, column := range table.ColumnList {
		if column.PrimaryKey {
			if StringsIndex(table.PrimaryKeys, column.ColumnName) < 0 {
				table.PrimaryKeys = append(table.PrimaryKeys, column.ColumnName)
			}
		}
	}
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

	var sqlList_ []string

	if table.TableComment != "" {
		sqlList_, err = this_.TableCommentSql(param, ownerName, table.TableName, table.TableComment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	if !this_.TableCreateColumnHasComment {
		for _, column := range table.ColumnList {
			if column.ColumnComment != "" {
				sqlList_, err = this_.ColumnCommentSql(param, ownerName, table.TableName, column.ColumnName, column.ColumnComment)
				if err != nil {
					return
				}
				sqlList = append(sqlList, sqlList_...)
			}
		}
	}

	for _, index := range table.IndexList {
		sqlList_, err = this_.IndexAddSql(param, ownerName, table.TableName, index)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
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

	var isNullable string
	if data["isNullable"] != nil {
		isNullable = GetStringValue(data["isNullable"])
	}
	if data["ISNULLABLE"] != nil {
		isNullable = GetStringValue(data["ISNULLABLE"])
	}
	if isNullable != "" {
		if strings.EqualFold(isNullable, "no") || strings.EqualFold(isNullable, "n") {
			column.ColumnNotNull = true
		}
	}

	if GetStringValue(data["isNotNull"]) == "1" || GetStringValue(data["ISNOTNULL"]) == "1" {
		column.ColumnNotNull = true
	}
	var columnType string
	if data["columnType"] != nil {
		columnType = data["columnType"].(string)
	}
	if data["COLUMNTYPE"] != nil {
		columnType = data["COLUMNTYPE"].(string)
	}
	if column.ColumnDataType == "" {
		if strings.Contains(columnType, "(") {
			column.ColumnDataType = columnType[:strings.Index(columnType, "(")]
		} else {
			column.ColumnDataType = columnType
		}
	}
	if strings.Contains(column.ColumnDataType, "(") {
		column.ColumnDataType = column.ColumnDataType[:strings.Index(column.ColumnDataType, "(")]
	}
	dataLength := GetStringValue(data["DATA_LENGTH"])
	if dataLength != "" && dataLength != "0" {
		column.ColumnLength, err = StringToInt(dataLength)
		if err != nil {
			return
		}
	}
	dataPrecision := GetStringValue(data["DATA_PRECISION"])
	if dataPrecision != "" && dataPrecision != "0" {
		column.ColumnLength, err = StringToInt(dataPrecision)
		if err != nil {
			return
		}
	}
	dataScale := GetStringValue(data["DATA_SCALE"])
	if dataScale != "" && dataScale != "0" {
		column.ColumnDecimal, err = StringToInt(dataScale)
		if err != nil {
			return
		}
	}
	characterMaximumLength := GetStringValue(data["CHARACTER_MAXIMUM_LENGTH"])
	if characterMaximumLength != "" && characterMaximumLength != "0" {
		column.ColumnLength, err = StringToInt(characterMaximumLength)
		if err != nil {
			return
		}
	}
	numericPrecision := GetStringValue(data["NUMERIC_PRECISION"])
	if numericPrecision != "" && numericPrecision != "0" {
		column.ColumnLength, err = StringToInt(numericPrecision)
		if err != nil {
			return
		}
	}
	numericScale := GetStringValue(data["NUMERIC_SCALE"])
	if numericScale != "" && numericScale != "0" {
		column.ColumnDecimal, err = StringToInt(numericScale)
		if err != nil {
			return
		}
	}
	datetimePrecision := GetStringValue(data["DATETIME_PRECISION"])
	if datetimePrecision != "" && datetimePrecision != "0" {
		column.ColumnLength, err = StringToInt(datetimePrecision)
		if err != nil {
			return
		}
	}
	columnTypeInfo, err := this_.GetColumnTypeInfo(column.ColumnDataType)
	if err != nil {
		bs, _ = json.Marshal(data)
		err = errors.New("ColumnModel error column data:" + string(bs) + ",error:" + err.Error())
		return
	}
	column.ColumnDataType = columnTypeInfo.Name

	if column.ColumnDefault != "" {
		if strings.HasSuffix(column.ColumnDefault, "::"+column.ColumnDataType) {
			column.ColumnDefault = strings.TrimSuffix(column.ColumnDefault, "::"+column.ColumnDataType)
		}
		column.ColumnDefault = strings.TrimLeft(column.ColumnDefault, "'")
		column.ColumnDefault = strings.TrimRight(column.ColumnDefault, "'")
		column.ColumnDefault = strings.TrimLeft(column.ColumnDefault, "\"")
		column.ColumnDefault = strings.TrimRight(column.ColumnDefault, "\"")
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

func (this_ *mappingDialect) ColumnUpdateSql(param *ParamModel, ownerName string, tableName string, oldColumn *ColumnModel, column *ColumnModel) (sqlList []string, err error) {
	if oldColumn.ColumnName == "" {
		oldColumn.ColumnName = column.ColumnName
	}
	columnTypePack, err := this_.ColumnTypePack(column)
	if err != nil {
		return
	}
	columnDefaultPack, err := this_.ColumnDefaultPack(param, column)
	if err != nil {
		return
	}
	data := map[string]string{
		"ownerName":         ownerName,
		"tableName":         tableName,
		"oldColumnName":     oldColumn.ColumnName,
		"columnTypePack":    columnTypePack,
		"columnDefaultPack": columnDefaultPack,
	}
	var sqlList_ []string
	var hasChangeName bool
	if oldColumn.ColumnName != column.ColumnName {
		hasChangeName = true
	}
	var hasChangeComment bool
	if oldColumn.ColumnComment != column.ColumnComment {
		hasChangeComment = true
	}
	var hasChangeAfter bool
	if oldColumn.ColumnAfterColumn != column.ColumnAfterColumn {
		hasChangeAfter = true
	}
	if !this_.ColumnUpdateHasRename {
		if hasChangeName {
			sqlList_, err = this_.ColumnRenameSql(param, ownerName, tableName, oldColumn.ColumnName, column.ColumnName)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
			hasChangeName = false
		}
	}
	if !this_.ColumnUpdateHasComment {
		if hasChangeComment {
			sqlList_, err = this_.ColumnCommentSql(param, ownerName, tableName, column.ColumnName, column.ColumnComment)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
			hasChangeComment = false
		}
	}
	if !this_.ColumnUpdateHasAfter {
		if hasChangeAfter {
			sqlList_, err = this_.ColumnAfterSql(param, ownerName, tableName, column.ColumnName, column.ColumnComment)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
			hasChangeAfter = false
		}
	}
	if hasChangeName || hasChangeComment || hasChangeAfter ||
		oldColumn.ColumnDataType != column.ColumnDataType ||
		oldColumn.ColumnLength != column.ColumnLength ||
		oldColumn.ColumnDecimal != column.ColumnDecimal ||
		oldColumn.ColumnDefault != column.ColumnDefault {
		sqlList_, err = this_.FormatSql(this_.ColumnUpdate, param,
			oldColumn,
			column,
			data,
		)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
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

func (this_ *mappingDialect) ColumnRenameSql(param *ParamModel, ownerName string, tableName string, oldColumnName string, columnName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnRename, param,
		map[string]string{
			"ownerName":     ownerName,
			"tableName":     tableName,
			"oldColumnName": oldColumnName,
			"columnName":    columnName,
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

func (this_ *mappingDialect) ColumnAfterSql(param *ParamModel, ownerName string, tableName string, columnName string, columnAfterColumn string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnAfter, param,
		map[string]string{
			"ownerName":         ownerName,
			"tableName":         tableName,
			"columnName":        columnName,
			"columnAfterColumn": columnAfterColumn,
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

	if GetStringValue(data["UNIQUENESS"]) == "UNIQUE" {
		index.IndexType = "unique"
	}
	if GetStringValue(data["NON_UNIQUE"]) == "0" {
		index.IndexType = "unique"
	}
	if GetStringValue(data["isUnique"]) == "1" {
		index.IndexType = "unique"
	}
	if GetStringValue(data["UNIQUENESS"]) == "UNIQUE" {
		index.IndexType = "unique"
	}
	indexTypeInfo, err := this_.GetIndexTypeInfo(index.IndexType)
	if err != nil {
		//fmt.Println(data)
		return
	}
	index.IndexType = indexTypeInfo.Name

	return
}

func (this_ *mappingDialect) IndexAddSql(param *ParamModel, ownerName string, tableName string, index *IndexModel) (sqlList []string, err error) {

	indexType, err := this_.IndexTypeFormat(index)
	if err != nil {
		return
	}
	indexType = strings.TrimSpace(indexType)
	if indexType == "" {
		return
	}

	indexName, err := this_.IndexNameFormat(param, ownerName, tableName, index)
	if err != nil {
		return
	}
	indexName = strings.TrimSpace(indexName)
	if indexName == "" {
		indexName = index.IndexName
	}
	sqlList, err = this_.FormatSql(this_.IndexAdd, param,
		index,
		map[string]string{
			"ownerName": ownerName,
			"tableName": tableName,
			"indexName": indexName,
			"indexType": indexType,
		},
	)
	if err != nil {
		return
	}
	//fmt.Println("index add sql:", sqlList)
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
