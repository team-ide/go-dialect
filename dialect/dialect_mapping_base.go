package dialect

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

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

func (this_ *mappingDialect) OwnerTablePack(param *ParamModel, ownerName string, tableName string) string {
	if this_.SqlMapping.OwnerTablePack != nil {
		return this_.SqlMapping.OwnerTablePack(param, ownerName, tableName)
	}
	var res string
	if ownerName != "" {
		res += this_.OwnerNamePack(param, ownerName) + "."
	}
	if tableName != "" {
		res += this_.TableNamePack(param, tableName)
	}
	return res
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
	char := this_.SqlValuePackChar
	escapeChar := this_.SqlValueEscapeChar
	if param != nil {
		if param.SqlValuePackChar != nil {
			char = *param.SqlValuePackChar
		}
		if param.SqlValueEscapeChar != nil {
			escapeChar = *param.SqlValueEscapeChar
		}
	}
	var columnTypeInfo *ColumnTypeInfo
	if column != nil {
		columnTypeInfo, _ = this_.GetColumnTypeInfo(column)
	}
	return packingValue(column, columnTypeInfo, char, escapeChar, value)
}

func (this_ *mappingDialect) ColumnDefaultPack(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
	var columnTypeInfo *ColumnTypeInfo
	if column != nil {
		columnTypeInfo, err = this_.GetColumnTypeInfo(column)
		if err != nil {
			bs, _ := json.Marshal(column)
			err = errors.New("ColumnDefaultPack error column:" + string(bs) + ",error:" + err.Error())
		}
	}
	if columnTypeInfo != nil && columnTypeInfo.ColumnDefaultPack != nil {
		columnDefaultPack, err = columnTypeInfo.ColumnDefaultPack(param, column)
		return
	}
	if column.ColumnDefault == "" {
		return
	}
	columnDefaultPack = this_.SqlValuePack(param, column, column.ColumnDefault)
	return
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

func (this_ *mappingDialect) InsertSql(param *ParamModel, insert *InsertModel) (sqlList []string, err error) {

	sql := "INSERT INTO "
	sql += this_.OwnerTablePack(param, insert.OwnerName, insert.TableName)

	sql += "(" + this_.ColumnNamesPack(param, insert.Columns) + ")"
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
				sql += this_.SqlValuePack(param, nil, value.Value)
				break
			case ValueTypeNumber:
				sql += value.Value
				break
			case ValueTypeFunc:

				var funcStr = value.Value
				//funcStr, err = this_.FormatFunc(funcStr)
				//if err != nil {
				//	return
				//}
				sql += funcStr
				break
			}
		}

		sql += `) `
	}

	sqlList = append(sqlList, sql)
	return
}
func (this_ *mappingDialect) InsertDataListSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}) (sqlList []string, batchSqlList []string, err error) {
	var batchSqlCache = make(map[string]string)
	var batchSqlIndexCache = make(map[string]int)
	var columnNames []string
	for _, one := range columnList {
		columnNames = append(columnNames, one.ColumnName)
	}
	for _, data := range dataList {
		var columnList_ []string
		var values = "("
		for _, column := range columnList {
			str := this_.SqlValuePack(param, column, data[column.ColumnName])
			if strings.EqualFold(str, "null") {
				continue
			}
			columnList_ = append(columnList_, column.ColumnName)
			values += str + ", "
		}
		values = strings.TrimSuffix(values, ", ")
		values += ")"

		insertSqlInfo := "INSERT INTO "
		insertSqlInfo += this_.OwnerTablePack(param, ownerName, tableName)
		insertSqlInfo += " ("
		insertSqlInfo += this_.ColumnNamesPack(param, columnList_)
		insertSqlInfo += ") VALUES "

		sqlList = append(sqlList, insertSqlInfo+values)

		key := strings.Join(columnList_, ",")
		find, ok := batchSqlCache[key]
		if ok {
			find += ",\n" + values
			batchSqlCache[key] = find
			batchSqlList[batchSqlIndexCache[key]] = find
		} else {
			find = insertSqlInfo + "\n" + values
			batchSqlIndexCache[key] = len(batchSqlCache)
			batchSqlCache[key] = find
			batchSqlList = append(batchSqlList, find)
		}
	}
	return
}
func (this_ *mappingDialect) PackPageSql(selectSql string, pageSize int, pageNo int) (pageSql string) {
	if this_.SqlMapping.PackPageSql != nil {
		return this_.SqlMapping.PackPageSql(selectSql, pageSize, pageNo)
	}
	pageSql = selectSql + fmt.Sprintf(" LIMIT %d,%d", pageSize*(pageNo-1), pageSize)
	return
}
func (this_ *mappingDialect) ReplaceSqlVariable(sqlInfo string, args []interface{}) (variableSql string) {
	if this_.SqlMapping.ReplaceSqlVariable != nil {
		return this_.SqlMapping.ReplaceSqlVariable(sqlInfo, args)
	}
	variableSql = sqlInfo
	return
}
