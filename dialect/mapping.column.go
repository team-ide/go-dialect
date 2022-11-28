package dialect

import (
	"encoding/json"
	"errors"
	"regexp"
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

		var list = this_.columnTypeInfoList
		for _, one := range list {
			if len(one.Matches) == 0 {
				continue
			}
			var matched = false
			//fmt.Println("typeName:", typeName, ",MatchName:", one.Name, ",matches:", one.Matches)
			for _, match := range one.Matches {
				if match == strings.ToUpper(typeName) {
					matched = true
					break
				}
				//fmt.Println("typeName:", typeName, ",match:", match)
				if strings.Contains(match, "&") {

					match = strings.ReplaceAll(match, "&&", "&")
					ss := strings.Split(match, "&")
					for _, s := range ss {
						s = strings.TrimSpace(s)
						if s == "" {
							continue
						}
						if strings.Contains(s, ">") || strings.Contains(s, "=") || strings.Contains(s, "<") {

						} else {
							if !regexp.MustCompile(s).MatchString(strings.ToUpper(typeName)) {
								matched = false
								break
							}
						}
					}
				} else {
					if regexp.MustCompile(match).MatchString(strings.ToUpper(typeName)) {
						matched = true
						break
					}
				}

			}
			if matched {
				columnTypeInfo = one
				break
			}
		}
	}

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
	if strings.Contains(columnTypePack, "(") {
		beforeStr := columnTypePack[0:strings.Index(columnTypePack, "(")]
		endStr := columnTypePack[strings.Index(columnTypePack, "("):]
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
		endStr = strings.ReplaceAll(endStr, "$l", lStr)
		endStr = strings.ReplaceAll(endStr, "$d", dStr)
		endStr = strings.ReplaceAll(endStr, " )", ")")
		endStr = strings.ReplaceAll(endStr, " ,", ",")
		endStr = strings.ReplaceAll(endStr, ",)", ")")
		endStr = strings.TrimSuffix(endStr, "(,)")
		endStr = strings.TrimSuffix(endStr, "()")
		columnTypePack = beforeStr + endStr
	}
	return
}
