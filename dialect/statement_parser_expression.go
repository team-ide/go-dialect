package dialect

import (
	"errors"
	"strings"
)

func parseExpressionStatement(content string, parent SqlStatement) (expressionStatement *ExpressionStatement, err error) {
	content = strings.TrimSpace(content)

	var sqlStatements []SqlStatement

	var level int
	var levelStatement = make(map[int]SqlStatement)
	var str string

	var inStringPack string
	var inStringLevel int
	var stringPackChars = []string{"\"", "'"}
	var lastChar string
	var thisChar string

	strList := strings.Split(content, "")
	for i := 0; i < len(strList); i++ {
		thisChar = strList[i]

		if i > 0 {
			lastChar = strList[i-1]
		}
		packCharIndex := StringsIndex(stringPackChars, thisChar)
		var isStringEnd bool
		if packCharIndex >= 0 {
			// inStringLevel == 0 表示 不在 字符串 包装 中
			if inStringLevel == 0 {
				inStringPack = stringPackChars[packCharIndex]
				// 字符串包装层级 +1
				inStringLevel++
			} else {
				// 如果有转义符号 类似 “\'”，“\"”
				if lastChar == "\\" {
				} else if lastChar == inStringPack {
					// 如果 前一个字符 与字符串包装字符一致
					inStringLevel--
				} else {
					// 字符串包装层级 -1
					inStringLevel--
				}
			}
			if inStringLevel == 0 {
				isStringEnd = true
			}
		}
		var thisParentChildren *[]SqlStatement
		var thisParent = parent
		if levelStatement[level] == nil {
			thisParentChildren = &sqlStatements
		} else {
			thisParent = levelStatement[level].GetParent()
			if thisParent == parent {
				thisParentChildren = &sqlStatements
			} else {
				thisParentChildren = levelStatement[level].GetParent().GetChildren()
			}
		}

		if isStringEnd {
			stringValue := str
			stringValue = strings.TrimSuffix(stringValue, stringPackChars[packCharIndex])
			stringValue = strings.TrimPrefix(stringValue, stringPackChars[packCharIndex])
			stringStatement := &ExpressionStringStatement{
				Value: stringValue,
				AbstractSqlStatement: &AbstractSqlStatement{
					Parent:  thisParent,
					Content: str,
				},
			}
			if levelStatement[level] != nil {
				*levelStatement[level].GetChildren() = append(*levelStatement[level].GetChildren(), stringStatement)
			} else {
				*thisParentChildren = append(*thisParentChildren, stringStatement)
			}
			str = ""
		} else if inStringLevel == 0 {
			if thisChar == "(" {
				if thisParent == nil {
					err = errors.New("sql template [" + content + "] parse match start error")
					return
				}
				var statement SqlStatement
				if str != "" {
					statement = &ExpressionFuncStatement{
						Func: str,
						AbstractSqlStatement: &AbstractSqlStatement{
							Parent:  thisParent,
							Content: str,
						},
					}
				} else {
					statement = &ExpressionBracketsStatement{
						AbstractSqlStatement: &AbstractSqlStatement{
							Parent: thisParent,
						},
					}
				}
				*thisParentChildren = append(*thisParentChildren, statement)
				level++
				levelStatement[level] = statement
				str = ""
			} else if thisChar == ")" {
				if thisParent == nil || levelStatement[level] == nil {
					err = errors.New("sql template [" + content + "] parse match end error")
					return
				}
				if str != "" {
					statement := &ExpressionIdentifierStatement{
						Identifier: str,
						AbstractSqlStatement: &AbstractSqlStatement{
							Parent:  thisParent,
							Content: str,
						},
					}
					*levelStatement[level].GetChildren() = append(*levelStatement[level].GetChildren(), statement)
				}
				levelStatement[level] = nil
				level--
				str = ""
			} else if thisChar == "," {
				if thisParent == nil || levelStatement[level] == nil {
					err = errors.New("sql template [" + content + "] parse match end error")
					return
				}
				statement := &ExpressionIdentifierStatement{
					Identifier: str,
					AbstractSqlStatement: &AbstractSqlStatement{
						Parent:  thisParent,
						Content: str,
					},
				}
				*levelStatement[level].GetChildren() = append(*levelStatement[level].GetChildren(), statement)
				str = ""
			} else {
				str += thisChar
			}
		} else {
			str += thisChar
		}
	}
	if str != "" {
		statement := &ExpressionIdentifierStatement{
			Identifier: str,
			AbstractSqlStatement: &AbstractSqlStatement{
				Parent:  parent,
				Content: str,
			},
		}
		sqlStatements = append(sqlStatements, statement)
	}
	expressionStatement = &ExpressionStatement{
		AbstractSqlStatement: &AbstractSqlStatement{
			Children: sqlStatements,
		},
	}
	return
}
