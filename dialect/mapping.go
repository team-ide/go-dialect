package dialect

import "strings"

type SqlMapping struct {
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

type SqlMappingStatement struct {
	SqlMapping   *SqlMapping
	OwnersSelect *RootStatement
	OwnerSelect  *RootStatement
	OwnerCreate  *RootStatement
	OwnerDelete  *RootStatement

	TablesSelect *RootStatement
	TableSelect  *RootStatement
	TableCreate  *RootStatement
	TableDelete  *RootStatement
	TableComment *RootStatement
	TableRename  *RootStatement

	ColumnsSelect *RootStatement
	ColumnSelect  *RootStatement
	ColumnAdd     *RootStatement
	ColumnDelete  *RootStatement
	ColumnComment *RootStatement
	ColumnRename  *RootStatement
	ColumnUpdate  *RootStatement

	PrimaryKeysSelect *RootStatement
	PrimaryKeyAdd     *RootStatement
	PrimaryKeyDelete  *RootStatement

	IndexesSelect   *RootStatement
	IndexAdd        *RootStatement
	IndexDelete     *RootStatement
	IndexNameFormat *RootStatement
}

func (this_ *SqlMappingStatement) OwnerNamePack(param *ParamModel, ownerName string) string {
	char := ""
	if this_.SqlMapping != nil {
		char = this_.SqlMapping.OwnerNamePackChar
	}
	if param != nil {
		if param.OwnerNamePack != nil && !*param.OwnerNamePack {
			char = ""
		} else if param.OwnerNamePackChar != nil {
			char = *param.OwnerNamePackChar
		}
	}
	return packingName(char, ownerName)
}

func (this_ *SqlMappingStatement) TableNamePack(param *ParamModel, tableName string) string {
	char := ""
	if this_.SqlMapping != nil {
		char = this_.SqlMapping.TableNamePackChar
	}
	if param != nil {
		if param.TableNamePack != nil && !*param.TableNamePack {
			char = ""
		} else if param.TableNamePackChar != nil {
			char = *param.TableNamePackChar
		}
	}
	return packingName(char, tableName)
}

func (this_ *SqlMappingStatement) ColumnNamePack(param *ParamModel, columnName string) string {
	char := ""
	if this_.SqlMapping != nil {
		char = this_.SqlMapping.ColumnNamePackChar
	}
	if param != nil {
		if param.ColumnNamePack != nil && !*param.ColumnNamePack {
			char = ""
		} else if param.ColumnNamePackChar != nil {
			char = *param.ColumnNamePackChar
		}
	}
	return packingName(char, columnName)
}

func (this_ *SqlMappingStatement) ColumnNamesPack(param *ParamModel, columnNames []string) string {
	char := ""
	if this_.SqlMapping != nil {
		char = this_.SqlMapping.ColumnNamePackChar
	}
	if param != nil {
		if param.ColumnNamePack != nil && !*param.ColumnNamePack {
			char = ""
		} else if param.ColumnNamePackChar != nil {
			char = *param.ColumnNamePackChar
		}
	}
	return packingNames(char, columnNames)
}

func (this_ *SqlMappingStatement) ColumnNamesStrPack(param *ParamModel, columnNamesStr string) string {
	return this_.ColumnNamesPack(param, strings.Split(columnNamesStr, ","))
}

func (this_ *SqlMappingStatement) SqlSplit(sqlStr string) (sqlList []string) {
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
