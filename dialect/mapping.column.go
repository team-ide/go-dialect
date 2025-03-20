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

func (this_ *SqlMapping) GetColumnTypeInfo(column *ColumnModel) (columnTypeInfo *ColumnTypeInfo, err error) {
	if column == nil || column.ColumnDataType == "" {
		err = errors.New("dialect [" + this_.DialectType().Name + "] GetColumnTypeInfo column data type is null")
		return
	}
	this_.columnTypeInfoCacheLock.Lock()
	defer this_.columnTypeInfoCacheLock.Unlock()

	if this_.columnTypeInfoCache == nil {
		this_.columnTypeInfoCache = make(map[string]*ColumnTypeInfo)
	}
	columnDataType := column.ColumnDataType

	key := strings.ToLower(columnDataType)
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
				var dataTypeMatch bool
				if match.DataType != "" {
					if strings.EqualFold(columnDataType, match.DataType) {
						dataTypeMatch = true
					}
					if regexp.MustCompile(match.DataType).MatchString(strings.ToUpper(columnDataType)) {
						dataTypeMatch = true
					}
				}
				if match.Match == nil {
					matched = dataTypeMatch
					if matched {
						break
					}
					continue
				}
				if match.Match != nil && (match.DataType == "" || dataTypeMatch) {
					if !match.Match(column.ColumnLength, column.ColumnPrecision, column.ColumnScale, column.ColumnDataType, column.ColumnDefault) {
						continue
					}
					if match.SetScript != "" {
						setValues := strings.Split(match.SetScript, ",")
						for _, setValueStr := range setValues {
							if !strings.Contains(setValueStr, "=") {
								continue
							}
							setName := setValueStr[0:strings.Index(setValueStr, "=")]
							setValue := setValueStr[strings.Index(setValueStr, "=")+1:]
							setName = strings.TrimSpace(setName)
							setValue = strings.TrimSpace(setValue)
							if strings.EqualFold(setName, "columnLength") {
								column.ColumnLength, err = StringToInt(setValue)
								if err != nil {
									err = errors.New("set value [" + setValue + "] error:" + err.Error())
									return
								}
							} else if strings.EqualFold(setName, "columnPrecision") {
								column.ColumnPrecision, err = StringToInt(setValue)
								if err != nil {
									err = errors.New("set value [" + setValue + "] error:" + err.Error())
									return
								}
							} else if strings.EqualFold(setName, "columnScale") {
								column.ColumnScale, err = StringToInt(setValue)
								if err != nil {
									err = errors.New("set value [" + setValue + "] error:" + err.Error())
									return
								}
							}
						}
					}
					matched = true
				}
			}
			if matched {
				columnTypeInfo = one
				break
			}
		}
	}

	if columnTypeInfo == nil {
		columnTypeInfo = &ColumnTypeInfo{
			Name:     column.ColumnDataType,
			Format:   column.ColumnDataType,
			IsExtend: true,
		}
		this_.columnTypeInfoCache[key] = columnTypeInfo
		//err = errors.New("dialect [" + this_.DialectType().Name + "] GetColumnTypeInfo not support column type name [" + column.ColumnDataType + "]")
		//return
	}
	return
}

func (this_ *SqlMapping) ColumnTypePack(column *ColumnModel) (columnTypePack string, err error) {
	columnTypeInfo, err := this_.GetColumnTypeInfo(column)
	if err != nil {
		bs, _ := json.Marshal(column)
		err = errors.New("ColumnTypePack error column:" + string(bs) + ",error:" + err.Error())
		return
	}
	if columnTypeInfo.ColumnTypePack != nil {
		columnTypePack, err = columnTypeInfo.ColumnTypePack(column)
		if err != nil {
			return
		}
	} else {
		if columnTypeInfo.IsEnum {
			enums := column.ColumnEnums
			if len(enums) == 0 {
				enums = []string{""}
			}
			columnTypePack = columnTypeInfo.Name + "(" + packingValues("'", enums) + ")"
		} else {
			columnTypePack = columnTypeInfo.Format
			if strings.Contains(columnTypePack, "(") {
				beforeStr := columnTypePack[0:strings.Index(columnTypePack, "(")]
				endStr := columnTypePack[strings.Index(columnTypePack, "("):]
				lStr := ""
				pStr := ""
				sStr := ""
				if column.ColumnLength >= 0 {
					lStr = strconv.Itoa(column.ColumnLength)
				}
				if column.ColumnPrecision >= 0 {
					pStr = strconv.Itoa(column.ColumnPrecision)
				}
				if column.ColumnScale >= 0 {
					sStr = strconv.Itoa(column.ColumnScale)
				}
				if pStr == "0" && sStr == "0" {
					pStr = ""
					sStr = ""
				}
				endStr = strings.ReplaceAll(endStr, "$l", lStr)
				endStr = strings.ReplaceAll(endStr, "$p", pStr)
				endStr = strings.ReplaceAll(endStr, "$s", sStr)
				columnTypePack = beforeStr + endStr
			}
		}
	}
	columnTypePack = strings.ReplaceAll(columnTypePack, " )", ")")
	columnTypePack = strings.ReplaceAll(columnTypePack, " ,", ",")
	columnTypePack = strings.ReplaceAll(columnTypePack, ",)", ")")
	columnTypePack = strings.ReplaceAll(columnTypePack, "(,)", "")
	columnTypePack = strings.ReplaceAll(columnTypePack, "()", "")

	return
}
