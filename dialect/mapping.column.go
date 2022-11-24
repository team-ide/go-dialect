package dialect

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func (this_ *SqlMapping) GetColumnTypeInfos() (columnTypeInfoList []*ColumnTypeInfo) {
	list := this_.columnTypeInfoList
	for _, one := range list {
		if one.IsExtend {
			continue
		}
		columnTypeInfoList = append(columnTypeInfoList, one)
	}
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
		return
	}
	return
}

func (this_ *SqlMapping) ColumnTypePack(column *ColumnModel) (columnTypePack string, err error) {
	columnTypeInfo, err := this_.GetColumnTypeInfo(column.ColumnDataType)
	if err != nil {
		bs, _ := json.Marshal(column)
		err = errors.New("ColumnTypePack error column:" + string(bs) + ",error:" + err.Error())
		return
	}
	if columnTypeInfo.ColumnTypePack != nil {
		columnTypePack, err = columnTypeInfo.ColumnTypePack(column)
		return
	}
	if columnTypeInfo.IsEnum {
		enums := column.ColumnEnums
		if len(enums) == 0 {
			enums = []string{""}
		}
		columnTypePack = columnTypeInfo.Name + "(" + packingValues("'", enums) + ")"
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
