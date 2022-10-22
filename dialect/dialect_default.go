package dialect

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

func NewDefaultDialect() *DefaultDialect {

	return &DefaultDialect{
		columnTypeInfoCache: make(map[string]*ColumnTypeInfo),
	}
}

type DefaultDialect struct {
	columnTypeInfoList      []*ColumnTypeInfo
	columnTypeInfoCache     map[string]*ColumnTypeInfo
	columnTypeInfoCacheLock sync.Mutex
}

func (this_ *DefaultDialect) DialectType() (dialectType *Type) {
	dialectType = DefaultType
	return
}

func (this_ *DefaultDialect) GetColumnTypeInfos() (columnTypeInfoList []*ColumnTypeInfo) {
	columnTypeInfoList = this_.columnTypeInfoList
	return
}

func (this_ *DefaultDialect) AddColumnTypeInfo(columnTypeInfo *ColumnTypeInfo) {
	this_.columnTypeInfoCacheLock.Lock()
	defer this_.columnTypeInfoCacheLock.Unlock()

	key := strings.ToLower(columnTypeInfo.Name)
	find := this_.columnTypeInfoCache[key]
	this_.columnTypeInfoCache[key] = columnTypeInfo
	if find == nil {
		this_.columnTypeInfoList = append(this_.columnTypeInfoList, columnTypeInfo)
	} else {
		var list = this_.columnTypeInfoList
		var newList []*ColumnTypeInfo
		for _, one := range list {
			if one == find {
				newList = append(newList, columnTypeInfo)
			} else {
				newList = append(newList, one)
			}
		}
		this_.columnTypeInfoList = newList
	}

	return
}
func (this_ *DefaultDialect) GetColumnTypeInfo(typeName string) (columnTypeInfo *ColumnTypeInfo, err error) {
	this_.columnTypeInfoCacheLock.Lock()
	defer this_.columnTypeInfoCacheLock.Unlock()

	key := strings.ToLower(typeName)
	columnTypeInfo = this_.columnTypeInfoCache[key]
	if columnTypeInfo == nil {
		err = errors.New("dialect [" + this_.DialectType().Name + "] not support type [" + typeName + "]")
		return
	}
	return
}
func (this_ *DefaultDialect) FormatColumnType(typeName string, length, decimal int) (columnType string, err error) {
	columnTypeInfo, err := this_.GetColumnTypeInfo(typeName)
	if err != nil {
		return
	}
	columnType = columnTypeInfo.FormatColumnType(length, decimal)
	return
}
func (this_ *DefaultDialect) ToColumnTypeInfo(columnType string) (columnTypeInfo *ColumnTypeInfo, length, decimal int, err error) {
	typeName := columnType
	if strings.Contains(columnType, "(") {
		typeName = columnType[0:strings.Index(columnType, "(")]
		lengthStr := columnType[strings.Index(columnType, "(")+1 : strings.Index(columnType, ")")]
		if strings.Contains(lengthStr, ",") {
			length, _ = strconv.Atoi(lengthStr[0:strings.Index(lengthStr, ",")])
			decimal, _ = strconv.Atoi(lengthStr[strings.Index(lengthStr, ",")+1:])
		} else {
			length, _ = strconv.Atoi(lengthStr)
		}
	}
	columnTypeInfo, err = this_.GetColumnTypeInfo(typeName)
	if err != nil {
		return
	}
	return
}

func (this_ *DefaultDialect) DatabaseModel(data map[string]interface{}) (database *DatabaseModel, err error) {
	return
}
func (this_ *DefaultDialect) DatabasesSelectSql() (sql string, err error) {
	return
}
func (this_ *DefaultDialect) DatabaseCreateSql(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {
	var sql string
	sql = `CREATE DATABASE `
	sql += param.packingCharacterDatabase(database.Name)

	sqlList = append(sqlList, sql)

	return
}
func (this_ *DefaultDialect) DatabaseDeleteSql(param *GenerateParam, databaseName string) (sqlList []string, err error) {
	var sql string
	sql = `DROP DATABASE `
	sql += param.packingCharacterDatabase(databaseName)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DefaultDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	return
}
func (this_ *DefaultDialect) TablesSelectSql(databaseName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) TableSelectSql(databaseName string, tableName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) TableCreateSql(param *GenerateParam, databaseName string, table *TableModel) (sqlList []string, err error) {

	createTableSql := `CREATE TABLE `

	if param.AppendDatabase && databaseName != "" {
		createTableSql += param.packingCharacterDatabase(databaseName) + "."
	}
	createTableSql += param.packingCharacterTable(table.Name)

	createTableSql += `(`
	createTableSql += "\n"
	primaryKeys := ""
	if len(table.ColumnList) > 0 {
		for _, column := range table.ColumnList {
			var columnSql = param.packingCharacterColumn(column.Name)

			var columnType string
			columnType, err = this_.FormatColumnType(column.Type, column.Length, column.Decimal)
			if err != nil {
				return
			}
			columnSql += " " + columnType

			if column.NotNull {
				columnSql += ` NOT NULL`
			}
			if column.Default != "" {
				columnSql += ` DEFAULT ` + formatStringValue("'", column.Default)
			}

			if column.PrimaryKey {
				primaryKeys += "" + column.Name + ","
			}
			createTableSql += "\t" + columnSql + ",\n"
		}
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + param.packingCharacterColumns(primaryKeys) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`

	sqlList = append(sqlList, createTableSql)

	// 添加注释
	if table.Comment != "" {
		var sqlList_ []string
		sqlList_, err = this_.TableCommentSql(param, databaseName, table.Name, table.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}
	if len(table.ColumnList) > 0 {
		for _, one := range table.ColumnList {
			if one.Comment == "" {
				continue
			}
			var sqlList_ []string
			sqlList_, err = this_.ColumnCommentSql(param, databaseName, table.Name, one.Name, one.Comment)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}

	if len(table.IndexList) > 0 {
		for _, one := range table.IndexList {
			if one.Name == "" || len(one.Columns) == 0 {
				continue
			}
			var sqlList_ []string
			sqlList_, err = this_.IndexAddSql(param, databaseName, table.Name, one)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}
	return
}
func (this_ *DefaultDialect) TableCommentSql(param *GenerateParam, databaseName string, tableName string, comment string) (sqlList []string, err error) {
	sql := "COMMENT ON TABLE  "
	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += "" + param.packingCharacterTable(tableName)
	sql += " IS " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) TableDeleteSql(param *GenerateParam, databaseName string, tableName string) (sqlList []string, err error) {
	var sql string
	sql = `DROP TABLE `
	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)
	return
}

func (this_ *DefaultDialect) ColumnModel(data map[string]interface{}) (table *ColumnModel, err error) {
	return
}
func (this_ *DefaultDialect) ColumnsSelectSql(databaseName string, tableName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) ColumnSelectSql(databaseName string, tableName string, columnName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) ColumnAddSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.FormatColumnType(column.Type, column.Length, column.Decimal)
	if err != nil {
		return
	}

	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)

	sql += ` ADD (`
	sql += param.packingCharacterColumn(column.Name)
	sql += ` ` + columnType + ``
	if column.NotNull {
		sql += ` NOT NULL`
	}
	if column.Default != "" {
		sql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
	}
	sql += `)`

	sqlList = append(sqlList, sql)

	if column.Comment != "" {
		var sqlList_ []string
		sqlList_, err = this_.ColumnCommentSql(param, databaseName, tableName, column.Name, column.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	return
}
func (this_ *DefaultDialect) ColumnCommentSql(param *GenerateParam, databaseName string, tableName string, columnName string, comment string) (sqlList []string, err error) {
	sql := "COMMENT ON COLUMN "
	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += "" + param.packingCharacterTable(tableName)
	sql += "." + param.packingCharacterColumn(columnName)
	sql += " IS " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) columnRenameSql(param *GenerateParam, databaseName string, tableName string, oldName string, newName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)

	sql += ` RENAME COLUMN `
	sql += param.packingCharacterColumn(oldName)
	sql += ` TO `
	sql += param.packingCharacterColumn(newName)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) ColumnUpdateSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.FormatColumnType(column.Type, column.Length, column.Decimal)
	if err != nil {
		return
	}

	var sqlList_ []string

	if column.OldName != "" && column.OldName != column.Name {
		sqlList_, err = this_.columnRenameSql(param, databaseName, tableName, column.OldName, column.Name)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	if column.Type != column.OldType ||
		column.Length != column.OldLength ||
		column.Decimal != column.OldDecimal ||
		column.NotNull != column.OldNotNull ||
		column.Default != column.OldDefault ||
		column.BeforeColumn != "" {
		var sql string
		sql = `ALTER TABLE `

		if param.AppendDatabase && databaseName != "" {
			sql += param.packingCharacterDatabase(databaseName) + "."
		}
		sql += param.packingCharacterTable(tableName)

		sql += ` MODIFY (`
		sql += param.packingCharacterColumn(column.Name)
		sql += ` ` + columnType + ``
		if column.NotNull {
			sql += ` NOT NULL`
		}
		if column.Default != "" {
			sql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
		}
		sql += `)`

		sqlList = append(sqlList, sql)
	}
	if column.Comment != column.OldComment {
		sqlList_, err = this_.ColumnCommentSql(param, databaseName, tableName, column.Name, column.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}
	return
}
func (this_ *DefaultDialect) ColumnDeleteSql(param *GenerateParam, databaseName string, tableName string, columnName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)

	sql += ` DROP COLUMN `
	sql += param.packingCharacterColumn(columnName)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DefaultDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
	return
}
func (this_ *DefaultDialect) PrimaryKeysSelectSql(databaseName string, tableName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) PrimaryKeyAddSql(param *GenerateParam, databaseName string, tableName string, primaryKeys []string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += "" + param.packingCharacterTable(tableName)

	sql += ` ADD PRIMARY KEY `

	sql += "(" + param.packingCharacterColumns(strings.Join(primaryKeys, ",")) + ")"

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) PrimaryKeyDeleteSql(param *GenerateParam, databaseName string, tableName string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += "" + param.packingCharacterTable(tableName)

	sql += ` DROP PRIMARY KEY `

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DefaultDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
	return
}
func (this_ *DefaultDialect) IndexesSelectSql(databaseName string, tableName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) IndexSelectSql(databaseName string, tableName string, indexName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) IndexAddSql(param *GenerateParam, databaseName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	sql := "CREATE "
	switch strings.ToUpper(index.Type) {
	case "UNIQUE":
		sql += "UNIQUE INDEX"
	case "":
		sql += "INDEX"
	default:
		err = errors.New("dialect [" + this_.DialectType().Name + "] not support index type [" + index.Type + "]")
		return
	}

	sql += " " + param.packingCharacterColumn(index.Name) + ""

	sql += " ON "
	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += "" + param.packingCharacterTable(tableName)

	sql += "(" + param.packingCharacterColumns(strings.Join(index.Columns, ",")) + ")"

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) IndexUpdateSql(param *GenerateParam, databaseName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	var sqlList_ []string

	if index.OldName != "" {
		var sql = " DROP INDEX " + param.packingCharacterColumn(index.OldName) + ""
		sqlList = append(sqlList, sql)
	} else {
		var sql = " DROP INDEX " + param.packingCharacterColumn(index.Name) + ""
		sqlList = append(sqlList, sql)
	}

	sqlList_, err = this_.IndexAddSql(param, databaseName, tableName, index)
	if err != nil {
		return
	}
	sqlList = append(sqlList, sqlList_...)
	return
}
func (this_ *DefaultDialect) IndexDeleteSql(param *GenerateParam, databaseName string, tableName string, indexName string) (sqlList []string, err error) {
	sql := "DROP INDEX "
	sql += "" + param.packingCharacterColumn(indexName)

	sqlList = append(sqlList, sql)
	return
}
