package dialect

import "strings"

func (this_ *mappingDialect) OwnerNamePack(param *ParamModel, ownerName string) string {
	char := this_.OwnerNamePackChar
	if param != nil {
		if param.OwnerNamePack != nil && !*param.OwnerNamePack {
			char = ""
		} else if param.OwnerNamePackChar != nil {
			char = *param.OwnerNamePackChar
		}
	}
	return packingName(char, ownerName)
}

func (this_ *mappingDialect) TableNamePack(param *ParamModel, tableName string) string {
	char := this_.TableNamePackChar
	if param != nil {
		if param.TableNamePack != nil && !*param.TableNamePack {
			char = ""
		} else if param.TableNamePackChar != nil {
			char = *param.TableNamePackChar
		}
	}
	return packingName(char, tableName)
}

func (this_ *mappingDialect) ColumnNamePack(param *ParamModel, columnName string) string {
	char := this_.ColumnNamePackChar
	if param != nil {
		if param.ColumnNamePack != nil && !*param.ColumnNamePack {
			char = ""
		} else if param.ColumnNamePackChar != nil {
			char = *param.ColumnNamePackChar
		}
	}
	return packingName(char, columnName)
}

func (this_ *mappingDialect) ColumnNamesPack(param *ParamModel, columnNames []string) string {
	char := this_.ColumnNamePackChar
	if param != nil {
		if param.ColumnNamePack != nil && !*param.ColumnNamePack {
			char = ""
		} else if param.ColumnNamePackChar != nil {
			char = *param.ColumnNamePackChar
		}
	}
	return packingNames(char, columnNames)
}

func (this_ *mappingDialect) SqlValuePack(param *ParamModel, column *ColumnModel, value interface{}) string {
	var columnTypeInfo *ColumnTypeInfo
	if column != nil {
		//columnTypeInfo, _ = this_.GetColumnTypeInfo(column.Type)
	}
	return packingValue(columnTypeInfo, `'`, `'`, value)
}

func (this_ *mappingDialect) ColumnNamesStrPack(param *ParamModel, columnNamesStr string) string {
	return this_.ColumnNamesPack(param, strings.Split(columnNamesStr, ","))
}

func (this_ *mappingDialect) IsSqlEnd(sqlStr string) (isSqlEnd bool) {
	if !strings.HasSuffix(strings.TrimSpace(sqlStr), ";") {
		return
	}
	cacheKey := UUID()
	sqlCache := sqlStr
	sqlCache = strings.ReplaceAll(sqlCache, `''`, `|-`+cacheKey+`-|`)
	sqlCache = strings.ReplaceAll(sqlCache, `""`, `|--`+cacheKey+`--|`)

	var inStringLevel int
	var inStringPack byte
	var thisChar byte
	var lastChar byte

	var stringPackChars = []byte{'"', '\''}
	for i := 0; i < len(sqlCache); i++ {
		thisChar = sqlCache[i]
		if i > 0 {
			lastChar = sqlCache[i-1]
		}

		// inStringLevel == 0 表示 不在 字符串 包装 中
		if thisChar == ';' && inStringLevel == 0 {
		} else {
			packCharIndex := BytesIndex(stringPackChars, thisChar)
			if packCharIndex >= 0 {
				// inStringLevel == 0 表示 不在 字符串 包装 中
				if inStringLevel == 0 {
					inStringPack = stringPackChars[packCharIndex]
					// 字符串包装层级 +1
					inStringLevel++
				} else {
					if thisChar != inStringPack {
					} else if lastChar == '\\' { // 如果有转义符号 类似 “\'”，“\"”
					} else if lastChar == inStringPack {
						// 如果 前一个字符 与字符串包装字符一致
					} else {
						// 字符串包装层级 -1
						inStringLevel--
					}
				}
			}
		}

	}
	isSqlEnd = inStringLevel == 0
	return
}
func (this_ *mappingDialect) SqlSplit(sqlStr string) (sqlList []string) {
	cacheKey := UUID()
	sqlCache := sqlStr
	sqlCache = strings.ReplaceAll(sqlCache, `''`, `|-`+cacheKey+`-|`)
	sqlCache = strings.ReplaceAll(sqlCache, `""`, `|--`+cacheKey+`--|`)

	var list []string
	var beg int

	var inStringLevel int
	var inStringPack byte
	var thisChar byte
	var lastChar byte

	var stringPackChars = []byte{'"', '\''}
	for i := 0; i < len(sqlCache); i++ {
		thisChar = sqlCache[i]
		if i > 0 {
			lastChar = sqlCache[i-1]
		}

		// inStringLevel == 0 表示 不在 字符串 包装 中
		if thisChar == ';' && inStringLevel == 0 {
			if i > 0 {
				list = append(list, sqlCache[beg:i])
			}
			beg = i + 1
		} else {
			packCharIndex := BytesIndex(stringPackChars, thisChar)
			if packCharIndex >= 0 {
				// inStringLevel == 0 表示 不在 字符串 包装 中
				if inStringLevel == 0 {
					inStringPack = stringPackChars[packCharIndex]
					// 字符串包装层级 +1
					inStringLevel++
				} else {
					if thisChar != inStringPack {
					} else if lastChar == '\\' { // 如果有转义符号 类似 “\'”，“\"”
					} else if lastChar == inStringPack {
						// 如果 前一个字符 与字符串包装字符一致
					} else {
						// 字符串包装层级 -1
						inStringLevel--
					}
				}
			}
		}

	}
	list = append(list, sqlCache[beg:])
	for _, sqlOne := range list {
		sqlOne = strings.TrimSpace(sqlOne)
		if sqlOne == "" {
			continue
		}
		sqlOne = strings.ReplaceAll(sqlOne, `|-`+cacheKey+`-|`, `''`)
		sqlOne = strings.ReplaceAll(sqlOne, `|--`+cacheKey+`--|`, `""`)
		sqlList = append(sqlList, sqlOne)
	}
	return
}
