package dialect

import (
	"errors"
	"regexp"
	"strings"
)

func getIfCondition(str string, parent SqlStatement) (condition string, expressionStatement *ExpressionStatement, err error) {
	condition = strings.TrimSpace(str)
	condition = strings.TrimPrefix(condition, "{")
	condition = strings.TrimSuffix(condition, "}")
	condition = strings.TrimSpace(condition)
	var reg *regexp.Regexp
	reg, err = regexp.Compile("(else\\s+)?if")
	if reg == nil || err != nil {
		return
	}
	condition = reg.ReplaceAllString(condition, "")
	condition = strings.TrimSpace(condition)

	expressionStatement, err = parseExpressionStatement(condition, parent)
	if err != nil {
		return
	}

	return
}
func (this_ *SqlStatementParser) checkIfStatement() (startIndex int, endIndex int, err error) {
	var reg *regexp.Regexp
	var finds [][]int
	var curStr = this_.curStr
	// 匹配格式如：
	// { if xxx }
	// { if (xxx) }
	reg, err = regexp.Compile("\\{\\s*if[\\s\\(]+.*\\}$")
	if reg == nil || err != nil {
		return
	}
	finds = reg.FindAllStringIndex(curStr, -1)
	if len(finds) > 0 {

		startIndex = finds[0][0]
		endIndex = finds[0][1]
		text := curStr[0:startIndex]
		if text != "" {
			var sqlStatements []SqlStatement
			sqlStatements, err = parseTextSqlStatement(text, this_.curParent)
			if err != nil {
				return
			}
			*this_.curParent.GetChildren() = append(*this_.curParent.GetChildren(), sqlStatements...)
		}
		statement := &IfSqlStatement{
			AbstractSqlStatement: &AbstractSqlStatement{
				Content: curStr[startIndex:endIndex],
				Parent:  this_.curParent,
			},
		}
		statement.Condition, statement.ConditionExpression, err = getIfCondition(statement.Content, statement.Parent)
		if err != nil {
			return
		}

		*this_.curParent.GetChildren() = append(*this_.curParent.GetChildren(), statement)
		this_.curStatement = statement
		this_.curParent = statement
		this_.curIf = statement
		this_.curElseIf = nil
		this_.curElse = nil
		return
	}
	return
}

func (this_ *SqlStatementParser) checkElseIfStatement() (startIndex int, endIndex int, err error) {
	var reg *regexp.Regexp
	var finds [][]int
	var curStr = this_.curStr

	// 匹配格式如：
	// { else if xxx }
	// { else if (xxx) }
	reg, err = regexp.Compile("\\{\\s*else\\s*if[\\s\\(]+.+\\s*\\}$")
	if reg == nil || err != nil {
		return
	}
	finds = reg.FindAllStringIndex(curStr, -1)
	if len(finds) > 0 {
		if this_.curIf == nil {
			err = errors.New("sql template [" + curStr + "] parse error, not find ”if“")
			return
		}
		startIndex = finds[0][0]
		endIndex = finds[0][1]
		text := curStr[0:startIndex]
		if text != "" {
			p := this_.curParent
			if this_.curElseIf != nil {
				p = this_.curElseIf
			} else {
				p = this_.curIf
			}
			var sqlStatements []SqlStatement
			sqlStatements, err = parseTextSqlStatement(text, p)
			if err != nil {
				return
			}
			if this_.curElseIf != nil {
				this_.curElseIf.Children = append(this_.curElseIf.Children, sqlStatements...)
			} else {
				this_.curIf.Children = append(this_.curIf.Children, sqlStatements...)
			}
		}

		statement := &ElseIfSqlStatement{
			AbstractSqlStatement: &AbstractSqlStatement{
				Content: curStr[startIndex:endIndex],
				Parent:  this_.curIf,
			},
			If:    this_.curIf,
			Index: len(this_.curIf.ElseIfs),
		}
		statement.Condition, statement.ConditionExpression, err = getIfCondition(statement.Content, statement.Parent)
		if err != nil {
			return
		}
		this_.curIf.ElseIfs = append(this_.curIf.ElseIfs, statement)
		this_.curStatement = statement
		this_.curParent = statement
		this_.curElseIf = statement
		this_.curElse = nil
		return
	}
	return
}
func (this_ *SqlStatementParser) checkElseStatement() (startIndex int, endIndex int, err error) {
	var reg *regexp.Regexp
	var finds [][]int
	var curStr = this_.curStr
	// 匹配格式如：
	// { else }
	reg, err = regexp.Compile("\\{\\s*else\\s*\\}$")
	if reg == nil || err != nil {
		return
	}
	finds = reg.FindAllStringIndex(curStr, -1)
	if len(finds) > 0 {
		if this_.curIf == nil {
			err = errors.New("sql template [" + curStr + "] parse error, not find ”if“")
			return
		}

		startIndex = finds[0][0]
		endIndex = finds[0][1]
		text := curStr[0:startIndex]
		if text != "" {
			p := this_.curParent
			if this_.curElseIf != nil {
				p = this_.curElseIf
			} else {
				p = this_.curIf
			}
			var sqlStatements []SqlStatement
			sqlStatements, err = parseTextSqlStatement(text, p)
			if err != nil {
				return
			}
			if this_.curElseIf != nil {
				this_.curElseIf.Children = append(this_.curElseIf.Children, sqlStatements...)
			} else {
				this_.curIf.Children = append(this_.curIf.Children, sqlStatements...)
			}
		}

		statement := &ElseSqlStatement{
			AbstractSqlStatement: &AbstractSqlStatement{
				Content: curStr[startIndex:endIndex],
				Parent:  this_.curIf,
			},
			If: this_.curIf,
		}

		this_.curIf.Else = statement
		this_.curStatement = statement
		this_.curParent = statement
		this_.curElse = statement
		return
	}
	return
}
func (this_ *SqlStatementParser) checkIfEndStatement() (startIndex int, endIndex int, err error) {
	var reg *regexp.Regexp
	var finds [][]int
	var curStr = this_.curStr
	// 匹配格式如：
	// { }
	reg, err = regexp.Compile("\\{\\s*\\}$")
	if reg == nil || err != nil {
		return
	}
	finds = reg.FindAllStringIndex(curStr, -1)
	if len(finds) > 0 {
		if this_.curIf == nil {
			err = errors.New("sql template [" + curStr + "] parse error, not find ”if“")
			return
		}

		startIndex = finds[0][0]
		endIndex = finds[0][1]
		text := curStr[0:startIndex]
		if text != "" {
			p := this_.curParent
			if this_.curElse != nil {
				p = this_.curElse
			} else if this_.curElseIf != nil {
				p = this_.curElseIf
			} else {
				p = this_.curIf
			}
			var sqlStatements []SqlStatement
			sqlStatements, err = parseTextSqlStatement(text, p)
			if err != nil {
				return
			}
			if this_.curElse != nil {
				this_.curElse.Children = append(this_.curElse.Children, sqlStatements...)
			} else if this_.curElseIf != nil {
				this_.curElseIf.Children = append(this_.curElseIf.Children, sqlStatements...)
			} else {
				this_.curIf.Children = append(this_.curIf.Children, sqlStatements...)
			}
		}

		this_.curStatement = nil
		this_.curParent = this_.curIf.Parent
		this_.curIf = nil
		this_.curElseIf = nil
		this_.curElse = nil
		return
	}

	return
}
