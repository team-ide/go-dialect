package dialect

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type SqlMapping struct {
	dialectType *Type

	columnTypeInfoList      []*ColumnTypeInfo
	columnTypeInfoCache     map[string]*ColumnTypeInfo
	columnTypeInfoCacheLock sync.Mutex

	OwnersSelect string
	OwnerSelect  string
	OwnerCreate  string
	OwnerDelete  string

	TablesSelect string
	TableSelect  string
	TableCreate  string
	TableDelete  string
	TableComment string
	TableRename  string

	ColumnsSelect string
	ColumnSelect  string
	ColumnAdd     string
	ColumnDelete  string
	ColumnComment string
	ColumnRename  string
	ColumnUpdate  string

	PrimaryKeysSelect string
	PrimaryKeyAdd     string
	PrimaryKeyDelete  string

	IndexesSelect   string
	IndexAdd        string
	IndexDelete     string
	IndexNameFormat string

	OwnerNamePackChar  string
	TableNamePackChar  string
	ColumnNamePackChar string

	MethodCache map[string]interface{}
}

func (this_ *SqlMapping) DialectType() (dialectType *Type) {
	dialectType = this_.dialectType
	return
}

func (this_ *SqlMapping) GetColumnTypeInfos() (columnTypeInfoList []*ColumnTypeInfo) {
	columnTypeInfoList = this_.columnTypeInfoList
	return
}

func (this_ *SqlMapping) AddColumnTypeInfo(columnTypeInfo *ColumnTypeInfo) {
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

func (this_ *SqlMapping) GetColumnTypeInfo(typeName string) (columnTypeInfo *ColumnTypeInfo, err error) {
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

func (this_ *SqlMapping) FormatColumnType(column *ColumnModel) (columnType string, err error) {
	columnTypeInfo, err := this_.GetColumnTypeInfo(column.ColumnType)
	if err != nil {
		return
	}
	columnType = columnTypeInfo.TypeFormat
	lStr := ""
	dStr := ""
	if column.ColumnLength >= 0 {
		lStr = strconv.Itoa(column.ColumnLength)
	}
	if column.ColumnDecimal >= 0 {
		dStr = strconv.Itoa(column.ColumnDecimal)
	}
	if column.ColumnLength == 0 && column.ColumnDecimal == 0 {
		lStr = ""
		dStr = ""
	}
	columnType = strings.ReplaceAll(columnType, "$l", lStr)
	columnType = strings.ReplaceAll(columnType, "$d", dStr)
	columnType = strings.ReplaceAll(columnType, " ", "")
	columnType = strings.ReplaceAll(columnType, ",)", ")")
	columnType = strings.TrimSuffix(columnType, "(,)")
	columnType = strings.TrimSuffix(columnType, "()")
	return
}
