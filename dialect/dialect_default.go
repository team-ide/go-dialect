package dialect

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

func NewDefaultDialect(dialectType *Type) *DefaultDialect {

	return &DefaultDialect{
		columnTypeInfoCache: make(map[string]*ColumnTypeInfo),
		funcTypeInfoCache:   make(map[string]*FuncTypeInfo),
		dialectType:         dialectType,
	}
}

type DefaultDialect struct {
	columnTypeInfoList      []*ColumnTypeInfo
	columnTypeInfoCache     map[string]*ColumnTypeInfo
	columnTypeInfoCacheLock sync.Mutex

	funcTypeInfoList      []*FuncTypeInfo
	funcTypeInfoCache     map[string]*FuncTypeInfo
	funcTypeInfoCacheLock sync.Mutex

	dialectType *Type
}

func (this_ *DefaultDialect) DialectType() (dialectType *Type) {
	dialectType = this_.dialectType
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
func (this_ *DefaultDialect) FormatColumnType(column *ColumnModel) (columnType string, err error) {
	columnTypeInfo, err := this_.GetColumnTypeInfo(column.Type)
	if err != nil {
		return
	}
	columnType = columnTypeInfo.FormatColumnType(column.Length, column.Decimal)
	return
}
func (this_ *DefaultDialect) FormatDefaultValue(column *ColumnModel) (defaultValue string) {
	defaultValue = "DEFAULT " + this_.PackValue(nil, column.Default)
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

func (this_ *DefaultDialect) GeFuncTypeInfos() (funcTypeInfoList []*FuncTypeInfo) {
	funcTypeInfoList = this_.funcTypeInfoList
	return
}

func (this_ *DefaultDialect) AddFuncTypeInfo(funcTypeInfo *FuncTypeInfo) {
	this_.funcTypeInfoCacheLock.Lock()
	defer this_.funcTypeInfoCacheLock.Unlock()

	key := strings.ToLower(funcTypeInfo.Name)
	find := this_.funcTypeInfoCache[key]
	this_.funcTypeInfoCache[key] = funcTypeInfo
	if find == nil {
		this_.funcTypeInfoList = append(this_.funcTypeInfoList, funcTypeInfo)
	} else {
		var list = this_.funcTypeInfoList
		var newList []*FuncTypeInfo
		for _, one := range list {
			if one == find {
				newList = append(newList, funcTypeInfo)
			} else {
				newList = append(newList, one)
			}
		}
		this_.funcTypeInfoList = newList
	}

	return
}
func (this_ *DefaultDialect) GetFuncTypeInfo(funcName string) (funcTypeInfo *FuncTypeInfo, err error) {
	this_.funcTypeInfoCacheLock.Lock()
	defer this_.funcTypeInfoCacheLock.Unlock()

	key := strings.ToLower(funcName)
	funcTypeInfo = this_.funcTypeInfoCache[key]
	if funcTypeInfo == nil {
		err = errors.New("dialect [" + this_.DialectType().Name + "] not support func [" + funcName + "]")
		return
	}
	return
}
func (this_ *DefaultDialect) FormatFunc(funcStr string) (res string, err error) {
	funcName := funcStr[:strings.Index(funcStr, "(")]

	funcTypeInfo, err := this_.GetFuncTypeInfo(funcName)
	if err != nil {
		return
	}
	res = funcTypeInfo.Format + funcStr[strings.Index(funcStr, "("):]
	return
}

func (this_ *DefaultDialect) PackOwner(ownerName string) string {
	return packingName(`"`, ownerName)
}

func (this_ *DefaultDialect) PackTable(tableName string) string {
	return packingName(`"`, tableName)
}

func (this_ *DefaultDialect) PackColumn(columnName string) string {
	return packingName(`"`, columnName)
}

func (this_ *DefaultDialect) PackColumns(columnNames []string) string {
	return packingNames(`"`, columnNames)
}

func (this_ *DefaultDialect) PackValue(column *ColumnModel, value interface{}) string {
	var columnTypeInfo *ColumnTypeInfo
	if column != nil {
		columnTypeInfo, _ = this_.GetColumnTypeInfo(column.Type)
	}
	return packingValue(columnTypeInfo, `'`, `'`, value)
}

func (this_ *DefaultDialect) OwnerModel(data map[string]interface{}) (owner *OwnerModel, err error) {
	return
}
func (this_ *DefaultDialect) OwnersSelectSql() (sql string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support owner select sql")
	return
}
func (this_ *DefaultDialect) OwnerCreateSql(owner *OwnerModel) (sqlList []string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support owner create sql")

	return
}
func (this_ *DefaultDialect) OwnerDeleteSql(ownerName string) (sqlList []string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support owner delete sql")
	return
}

func (this_ *DefaultDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	return
}
func (this_ *DefaultDialect) TablesSelectSql(ownerName string) (sql string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support table select")
	return
}
func (this_ *DefaultDialect) TableSelectSql(ownerName string, tableName string) (sql string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support table select")
	return
}
func (this_ *DefaultDialect) TableCreateSql(ownerName string, table *TableModel) (sqlList []string, err error) {

	createTableSql := `CREATE TABLE `

	if ownerName != "" {
		createTableSql += this_.PackOwner(ownerName) + "."
	}
	createTableSql += this_.PackTable(table.Name)

	createTableSql += `(`
	createTableSql += "\n"
	primaryKeys := ""
	if len(table.ColumnList) > 0 {
		for _, column := range table.ColumnList {
			var columnSql = this_.PackColumn(column.Name)

			var columnType string
			columnType, err = this_.FormatColumnType(column)
			if err != nil {
				return
			}
			columnSql += " " + columnType

			if column.Default != "" {
				columnSql += ` ` + this_.FormatDefaultValue(column)
			}
			if column.NotNull {
				columnSql += ` NOT NULL`
			}

			if column.PrimaryKey {
				primaryKeys += "" + column.Name + ","
			}
			createTableSql += "\t" + columnSql + ",\n"
		}
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + this_.PackColumns(strings.Split(primaryKeys, ",")) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`

	sqlList = append(sqlList, createTableSql)

	// 添加注释
	if table.Comment != "" {
		var sqlList_ []string
		sqlList_, err = this_.TableCommentSql(ownerName, table.Name, table.Comment)
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
			sqlList_, err = this_.ColumnCommentSql(ownerName, table.Name, one.Name, one.Comment)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}

	if len(table.IndexList) > 0 {
		for _, one := range table.IndexList {
			var sqlList_ []string
			sqlList_, err = this_.IndexAddSql(ownerName, table.Name, one)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}
	return
}
func (this_ *DefaultDialect) TableCommentSql(ownerName string, tableName string, comment string) (sqlList []string, err error) {
	sql := "COMMENT ON TABLE  "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)
	sql += " IS " + this_.PackValue(nil, comment)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) TableRenameSql(ownerName string, oldTableName string, newTableName string) (sqlList []string, err error) {
	sql := "ALTER TABLE  "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(oldTableName)
	sql += " RENAME TO  "
	sql += "" + this_.PackTable(newTableName)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) TableDeleteSql(ownerName string, tableName string) (sqlList []string, err error) {
	var sql string
	sql = `DROP TABLE `
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += this_.PackTable(tableName)
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DefaultDialect) ColumnModel(data map[string]interface{}) (table *ColumnModel, err error) {
	return
}
func (this_ *DefaultDialect) ColumnsSelectSql(ownerName string, tableName string) (sql string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support columns select")
	return
}
func (this_ *DefaultDialect) ColumnSelectSql(ownerName string, tableName string, columnName string) (sql string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support column select")
	return
}
func (this_ *DefaultDialect) ColumnAddSql(ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.FormatColumnType(column)
	if err != nil {
		return
	}

	var sql string
	sql = `ALTER TABLE `

	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += this_.PackTable(tableName)

	sql += ` ADD `
	sql += this_.PackColumn(column.Name)
	sql += ` ` + columnType + ``
	if column.Default != "" {
		sql += ` ` + this_.FormatDefaultValue(column)
	}
	if column.NotNull {
		sql += ` NOT NULL`
	}
	sql += ``

	sqlList = append(sqlList, sql)

	if column.Comment != "" {
		var sqlList_ []string
		sqlList_, err = this_.ColumnCommentSql(ownerName, tableName, column.Name, column.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	return
}
func (this_ *DefaultDialect) ColumnCommentSql(ownerName string, tableName string, columnName string, comment string) (sqlList []string, err error) {
	sql := "COMMENT ON COLUMN "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)
	sql += "." + this_.PackColumn(columnName)
	sql += " IS " + this_.PackValue(nil, comment)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) ColumnRenameSql(ownerName string, tableName string, oldColumnName string, newColumnName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += this_.PackTable(tableName)

	sql += ` RENAME COLUMN `
	sql += this_.PackColumn(oldColumnName)
	sql += ` TO `
	sql += this_.PackColumn(newColumnName)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) ColumnUpdateSql(ownerName string, tableName string, oldColumn *ColumnModel, newColumn *ColumnModel) (sqlList []string, err error) {

	if oldColumn.Name != newColumn.Name {
		var sqlList_ []string
		sqlList_, err = this_.ColumnRenameSql(ownerName, tableName, oldColumn.Name, newColumn.Name)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	if oldColumn.Type != newColumn.Type ||
		oldColumn.Length != newColumn.Length ||
		oldColumn.Decimal != newColumn.Decimal ||
		oldColumn.Default != newColumn.Default ||
		oldColumn.NotNull != newColumn.NotNull ||
		oldColumn.BeforeColumn != newColumn.BeforeColumn {
		var columnType string
		columnType, err = this_.FormatColumnType(newColumn)
		if err != nil {
			return
		}

		var sql string
		sql = `ALTER TABLE `

		if ownerName != "" {
			sql += this_.PackOwner(ownerName) + "."
		}
		sql += this_.PackTable(tableName)

		sql += ` MODIFY `
		sql += this_.PackColumn(newColumn.Name)
		sql += ` ` + columnType + ``
		if newColumn.Default == "" {
			sql += " DEFAULT NULL"
		} else {
			sql += " " + this_.FormatDefaultValue(newColumn)
		}
		if newColumn.NotNull {
			sql += ` NOT NULL`
		}
		sql += ``

		sqlList = append(sqlList, sql)
	}

	return
}
func (this_ *DefaultDialect) ColumnDeleteSql(ownerName string, tableName string, columnName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += this_.PackTable(tableName)

	sql += ` DROP COLUMN `
	sql += this_.PackColumn(columnName)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DefaultDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
	return
}
func (this_ *DefaultDialect) PrimaryKeysSelectSql(ownerName string, tableName string) (sql string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support primaryKeys select")
	return
}
func (this_ *DefaultDialect) PrimaryKeyAddSql(ownerName string, tableName string, primaryKeys []string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	sql += ` ADD PRIMARY KEY `

	sql += "(" + this_.PackColumns(primaryKeys) + ")"

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) PrimaryKeyDeleteSql(ownerName string, tableName string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	sql += ` DROP PRIMARY KEY `

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DefaultDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
	return
}
func (this_ *DefaultDialect) IndexesSelectSql(ownerName string, tableName string) (sql string, err error) {
	err = errors.New("dialect [" + this_.DialectType().Name + "] not support indexes select")
	return
}
func (this_ *DefaultDialect) IndexAddSql(ownerName string, tableName string, index *IndexModel) (sqlList []string, err error) {
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

	sql += " " + this_.PackColumn(index.Name) + ""

	sql += " ON "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	sql += "(" + this_.PackColumns(index.Columns) + ")"

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) IndexRenameSql(ownerName string, tableName string, oldIndexName string, newIndexName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER INDEX `

	sql += this_.PackColumn(oldIndexName)
	sql += ` RENAME TO `
	sql += this_.PackColumn(newIndexName)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DefaultDialect) IndexDeleteSql(ownerName string, tableName string, indexName string) (sqlList []string, err error) {
	sql := "DROP INDEX "
	sql += "" + this_.PackColumn(indexName)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DefaultDialect) InsertSql(insert *InsertModel) (sqlList []string, err error) {

	sql := "INSERT INTO "
	if insert.OwnerName != "" {
		sql += this_.PackOwner(insert.OwnerName) + "."
	}
	sql += "" + this_.PackTable(insert.TableName)

	sql += "(" + this_.PackColumns(insert.Columns) + ")"
	sql += ` VALUES `

	for rowIndex, row := range insert.Rows {
		if rowIndex > 0 {
			sql += `, `
		}
		sql += `( `

		for valueIndex, value := range row {
			if valueIndex > 0 {
				sql += `, `
			}
			switch value.Type {
			case ValueTypeString:
				sql += this_.PackValue(nil, value.Value)
				break
			case ValueTypeNumber:
				sql += value.Value
				break
			case ValueTypeFunc:

				var funcStr = value.Value
				funcStr, err = this_.FormatFunc(funcStr)
				if err != nil {
					return
				}
				sql += funcStr
				break
			}
		}

		sql += `) `
	}

	sqlList = append(sqlList, sql)
	return
}
