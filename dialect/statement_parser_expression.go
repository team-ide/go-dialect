package dialect

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	operators      = []string{"+", "-", "*", "/", "<=", ">=", "==", "<", ">"}
	matchOperators = []string{"\\+", "\\-", "\\*", "/", "<=", ">=", "==", "<", ">"}
)

func isOperator(str string) bool {
	return StringsIndex(operators, str) >= 0
}

func splitOperator(content string) (res []string, err error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return
	}
	reg := regexp.MustCompile("[(" + strings.Join(matchOperators, ")(") + ")]+")
	matches := reg.FindAllStringIndex(content, -1)
	if len(matches) == 0 {
		res = append(res, content)
		return
	}
	lastIndex := 0
	for _, match := range matches {
		res = append(res, content[lastIndex:match[0]])
		res = append(res, content[match[0]:match[1]])
		lastIndex = match[1]
	}
	if len(content) > lastIndex {
		res = append(res, content[lastIndex:])
	}

	return
}

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

	processStr := func(str string, statements *[]SqlStatement, thisParent SqlStatement) (err error) {
		if str == "" {
			return
		}
		var splitOperatorValues []string
		splitOperatorValues, err = splitOperator(str)
		if err != nil {
			return
		}
		for _, one := range splitOperatorValues {
			if one == "" {
				continue
			}
			if isOperator(one) {
				statement := &ExpressionOperatorStatement{
					Operator: one,
					AbstractSqlStatement: &AbstractSqlStatement{
						Parent:  thisParent,
						Content: one,
					},
				}
				*statements = append(*statements, statement)
			} else {
				number, e := strconv.ParseFloat(one, 64)
				if e != nil {
					statement := &ExpressionIdentifierStatement{
						Identifier: one,
						AbstractSqlStatement: &AbstractSqlStatement{
							Parent:  thisParent,
							Content: one,
						},
					}
					*statements = append(*statements, statement)
				} else {
					statement := &ExpressionNumberStatement{
						Value: number,
						AbstractSqlStatement: &AbstractSqlStatement{
							Parent:  thisParent,
							Content: one,
						},
					}
					*statements = append(*statements, statement)
				}

			}
		}
		return
	}
	for i := 0; i < len(strList); i++ {
		thisChar = strList[i]

		if i > 0 {
			lastChar = strList[i-1]
		}
		packCharIndex := StringsIndex(stringPackChars, thisChar)
		var isStringEnd bool
		var isStringStart bool
		if packCharIndex >= 0 {
			// inStringLevel == 0 表示 不在 字符串 包装 中
			if inStringLevel == 0 {
				inStringPack = stringPackChars[packCharIndex]
				// 字符串包装层级 +1
				inStringLevel++
				isStringStart = true
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
				if inStringLevel == 0 {
					isStringEnd = true
				}
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

		if isStringStart {
			if levelStatement[level] != nil {
				err = processStr(str, levelStatement[level].GetChildren(), thisParent)
				if err != nil {
					return
				}
			} else {
				err = processStr(str, thisParentChildren, thisParent)
				if err != nil {
					return
				}
			}
			str = ""

		} else if isStringEnd {
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
				var splitOperatorValues []string
				splitOperatorValues, err = splitOperator(str)
				if err != nil {
					return
				}

				for i, one := range splitOperatorValues {
					if one == "" {
						continue
					}
					if i < len(splitOperatorValues)-1 || isOperator(one) {
						err = processStr(one, thisParentChildren, thisParent)
						if err != nil {
							return
						}
					} else {
						statement = &ExpressionFuncStatement{
							Func: one,
							AbstractSqlStatement: &AbstractSqlStatement{
								Parent:  thisParent,
								Content: one,
							},
						}
					}
				}
				if statement == nil {
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
				err = processStr(str, levelStatement[level].GetChildren(), thisParent)
				if err != nil {
					return
				}
				levelStatement[level] = nil
				level--
				str = ""
			} else if thisChar == "," {
				if thisParent == nil || levelStatement[level] == nil {
					err = errors.New("sql template [" + content + "] parse match end error")
					return
				}
				err = processStr(str, levelStatement[level].GetChildren(), thisParent)
				if err != nil {
					return
				}
				str = ""
			} else {
				str += thisChar
			}
		} else {
			str += thisChar
		}
	}
	if str != "" {
		err = processStr(str, &sqlStatements, parent)
		if err != nil {
			return
		}
	}
	expressionStatement = &ExpressionStatement{
		AbstractSqlStatement: &AbstractSqlStatement{
			Children: sqlStatements,
		},
	}
	return
}
