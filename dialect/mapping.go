package dialect

import (
	"errors"
	"fmt"
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

	TablesSelect          string
	TableSelect           string
	TableCreate           string
	TableCreateColumn     string
	TableCreatePrimaryKey string
	TableDelete           string
	TableComment          string
	TableRename           string

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
	SqlValuePackChar   string
	SqlValueEscapeChar string

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

	if this_.columnTypeInfoCache == nil {
		this_.columnTypeInfoCache = make(map[string]*ColumnTypeInfo)
	}

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
	if typeName == "" {
		err = errors.New("dialect [" + this_.DialectType().Name + "] GetColumnTypeInfo column type name is null")
		return
	}
	this_.columnTypeInfoCacheLock.Lock()
	defer this_.columnTypeInfoCacheLock.Unlock()

	if this_.columnTypeInfoCache == nil {
		this_.columnTypeInfoCache = make(map[string]*ColumnTypeInfo)
	}

	key := strings.ToLower(typeName)
	columnTypeInfo = this_.columnTypeInfoCache[key]
	if columnTypeInfo == nil {
		err = errors.New("dialect [" + this_.DialectType().Name + "] GetColumnTypeInfo not support column type name [" + typeName + "]")
		fmt.Println(err)
		return
	}
	return
}

func (this_ *SqlMapping) ColumnTypePack(column *ColumnModel) (columnTypePack string, err error) {
	columnTypeInfo, err := this_.GetColumnTypeInfo(column.ColumnDataType)
	if err != nil {
		return
	}
	if columnTypeInfo.ColumnTypePack != nil {
		columnTypePack, err = columnTypeInfo.ColumnTypePack(column)
		return
	}
	columnTypePack = columnTypeInfo.Format
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
	columnTypePack = strings.ReplaceAll(columnTypePack, "$l", lStr)
	columnTypePack = strings.ReplaceAll(columnTypePack, "$d", dStr)
	columnTypePack = strings.ReplaceAll(columnTypePack, " ", "")
	columnTypePack = strings.ReplaceAll(columnTypePack, ",)", ")")
	columnTypePack = strings.TrimSuffix(columnTypePack, "(,)")
	columnTypePack = strings.TrimSuffix(columnTypePack, "()")
	return
}
