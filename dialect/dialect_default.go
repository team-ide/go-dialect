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
	return
}
func (this_ *DefaultDialect) DatabaseDeleteSql(param *GenerateParam, databaseName string) (sqlList []string, err error) {
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
	return
}
func (this_ *DefaultDialect) TableCommentSql(param *GenerateParam, databaseName string, tableName string, comment string) (sqlList []string, err error) {
	return
}
func (this_ *DefaultDialect) TableDeleteSql(param *GenerateParam, databaseName string, tableName string) (sqlList []string, err error) {
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
	return
}
func (this_ *DefaultDialect) ColumnUpdateSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	return
}
func (this_ *DefaultDialect) ColumnRenameSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	return
}
func (this_ *DefaultDialect) ColumnDeleteSql(param *GenerateParam, databaseName string, tableName string, columnName string) (sqlList []string, err error) {
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
	return
}
func (this_ *DefaultDialect) IndexUpdateSql(param *GenerateParam, databaseName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	return
}
func (this_ *DefaultDialect) IndexDeleteSql(param *GenerateParam, databaseName string, tableName string, indexName string) (sqlList []string, err error) {
	return
}
func (this_ *DefaultDialect) IndexRenameSql(param *GenerateParam, databaseName string, tableName string, indexName string, rename string) (sqlList []string, err error) {
	return
}

func (this_ *DefaultDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
	return
}
func (this_ *DefaultDialect) PrimaryKeysSelectSql(databaseName string, tableName string) (sql string, err error) {
	return
}
func (this_ *DefaultDialect) PrimaryKeyAddSql(param *GenerateParam, databaseName string, tableName string, primaryKeys []string) (sqlList []string, err error) {
	return
}
func (this_ *DefaultDialect) PrimaryKeyDeleteSql(param *GenerateParam, databaseName string, tableName string, primaryKeys []string) (sqlList []string, err error) {
	return
}
