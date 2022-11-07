package dialect

import (
	"errors"
	"strings"
)

func parseTextSqlStatement(content string, parent SqlStatement) (sqlStatements []SqlStatement, err error) {
	list, err := parseStringStatement(content, parent,
		func(thisChar string, parent SqlStatement) (sqlStatement SqlStatement) {
			if thisChar == "[" {
				sqlStatement = &IgnorableSqlStatement{
					AbstractSqlStatement: &AbstractSqlStatement{
						Parent: parent,
					},
				}
			}
			return
		},
		func(thisChar string) (isEnd bool) {
			if thisChar == "]" {
				isEnd = true
			}
			return
		},
	)
	if err != nil {
		return
	}

	var list_ []SqlStatement
	for _, one := range list {
		list_, err = parseTextExpressionStatement(*one.GetContent(), one.GetParent())
		if err != nil {
			return
		}
		switch one.(type) {
		case *IgnorableSqlStatement:
			*one.GetContent() = ""
			*one.GetChildren() = list_
			sqlStatements = append(sqlStatements, one)
		default:
			sqlStatements = append(sqlStatements, list_...)
		}
	}
	//fmt.Println(this_.Sql)
	return
}

func parseTextExpressionStatement(content string, parent SqlStatement) (sqlStatements []SqlStatement, err error) {

	list, err := parseStringStatement(content, parent,
		func(thisChar string, parent SqlStatement) (sqlStatement SqlStatement) {
			if thisChar == "{" {
				sqlStatement = &ExpressionStatement{
					AbstractSqlStatement: &AbstractSqlStatement{
						Parent: parent,
					},
				}
			}
			return
		},
		func(thisChar string) (matchStart bool) {
			matchStart = thisChar == "}"
			return
		},
	)
	if err != nil {
		return
	}
	var expressionStatement *ExpressionStatement
	for _, one := range list {
		switch one.(type) {
		case *ExpressionStatement:
			expressionStatement, err = parseExpressionStatement(*one.GetContent(), one.GetParent())
			if err != nil {
				return
			}
			sqlStatements = append(sqlStatements, expressionStatement)
		default:
			sqlStatements = append(sqlStatements, one)
		}
	}
	return
}

func parseStringStatement(content string, parent SqlStatement,
	matchStart func(thisChar string, parent SqlStatement) (sqlStatement SqlStatement),
	matchEnd func(thisChar string) (matchStart bool),
) (sqlStatements []SqlStatement, err error) {

	var level int
	var levelStatement = make(map[int]SqlStatement)
	var str string

	var inStringPack string
	var inStringLevel int
	var stringPackChars = []string{"\"", "'"}
	var lastChar string
	var thisChar string

	strList := strings.Split(content, "")
	var matchStatement SqlStatement
	for i := 0; i < len(strList); i++ {
		thisChar = strList[i]

		if i > 0 {
			lastChar = strList[i-1]
		}
		packCharIndex := StringsIndex(stringPackChars, thisChar)
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

		if inStringLevel == 0 {
			if matchStatement = matchStart(thisChar, thisParent); matchStatement != nil {
				if thisParent == nil {
					err = errors.New("sql template [" + content + "] parse match start error")
					return
				}

				if str != "" {
					if levelStatement[level] != nil {
						*levelStatement[level].GetContent() += str
					} else {
						textSqlStatement := &TextSqlStatement{
							AbstractSqlStatement: &AbstractSqlStatement{
								Parent:  thisParent,
								Content: str,
							},
						}
						*thisParentChildren = append(*thisParentChildren, textSqlStatement)
					}
				}
				*thisParentChildren = append(*thisParentChildren, matchStatement)
				level++
				levelStatement[level] = matchStatement
				str = ""

			} else if matchEnd(thisChar) {
				if thisParent == nil || levelStatement[level] == nil {
					err = errors.New("sql template [" + content + "] parse match end error")
					return
				}
				*levelStatement[level].GetContent() = str
				levelStatement[level] = nil
				level--
				str = ""
			} else {
				str += thisChar
			}
		} else {
			str += thisChar
		}
	}
	if str != "" {
		textSqlStatement := &TextSqlStatement{
			AbstractSqlStatement: &AbstractSqlStatement{
				Parent:  parent,
				Content: str,
			},
		}
		sqlStatements = append(sqlStatements, textSqlStatement)
	}
	return
}
