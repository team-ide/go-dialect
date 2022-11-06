package dialect

import (
	"errors"
	"strings"
)

func sqlStatementParse(content string) (sqlStatement *RootSqlStatement, err error) {
	parser := &SqlStatementParser{
		content: strings.Split(content, ""),
	}
	sqlStatement, err = parser.parse()
	return
}

type SqlStatementParser struct {
	content       []string
	contentLen    int
	curIndex      int // 当前索引
	curRowStart   int // 当前行开始索引
	curStr        string
	curStatement  SqlStatement
	curParent     SqlStatement
	curRow        int // 当前行
	curCol        int // 当前列
	bracketLevel  int // “{}” 层级
	braceLevel    int // “[]” 层级
	inStringLevel int
	inStringPack  string

	curIf     *IfSqlStatement
	curElseIf *ElseIfSqlStatement
	curElse   *ElseSqlStatement
}

func (this_ *SqlStatementParser) reset() {

	this_.contentLen = len(this_.content)
	this_.curRow = 0
	this_.curCol = 0
	this_.curIndex = 0
	this_.curRowStart = 0
	this_.curStr = ""
	this_.curParent = nil
	this_.curStatement = nil
	this_.bracketLevel = 0
	this_.braceLevel = 0

	this_.inStringLevel = 0
	this_.inStringPack = ""

	this_.curIf = nil
	this_.curElseIf = nil
	this_.curElse = nil
}

func (this_ *SqlStatementParser) parse() (sqlStatement *RootSqlStatement, err error) {
	sqlStatement = &RootSqlStatement{
		AbstractSqlStatement: &AbstractSqlStatement{},
	}
	this_.reset()
	this_.curParent = sqlStatement
	err = this_.parseStr()
	if err != nil {
		return
	}
	return
}

func (this_ *SqlStatementParser) parseStr() (err error) {
	if this_.contentLen <= this_.curIndex {
		return
	}
	var char = this_.content[this_.curIndex]
	var lastChar string
	if this_.curIndex > 0 {
		lastChar = this_.content[this_.curIndex-1]
	}
	this_.curStr += char

	var stringPackChars = []string{"\"", "'"}
	packCharIndex := StringsIndex(stringPackChars, char)
	if packCharIndex >= 0 {
		// inStringLevel == 0 表示 不在 字符串 包装 中
		if this_.inStringLevel == 0 {
			this_.inStringPack = stringPackChars[packCharIndex]
			// 字符串包装层级 +1
			this_.inStringLevel++
		} else {
			// 如果有转义符号 类似 “\'”，“\"”
			if lastChar == "\\" {
			} else if lastChar == this_.inStringPack {
				// 如果 前一个字符 与字符串包装字符一致
				this_.inStringLevel--
			} else {
				// 字符串包装层级 -1
				this_.inStringLevel--
			}
		}
	}

	if this_.inStringLevel == 0 {
		err = this_.checkStatement()
		if err != nil {
			return
		}
	}
	if char == "\n" {
		this_.curRow++
		this_.curCol = 0
	} else {
		this_.curCol++
	}
	this_.curIndex++
	err = this_.parseStr()
	return
}

func (this_ *SqlStatementParser) checkStatement() (err error) {
	var startIndex int
	var endIndex int
	defer func() {
		if startIndex != endIndex {
			this_.curStr = ""
		}
	}()

	switch this_.curStatement.(type) {
	case *IfSqlStatement:
		if startIndex, endIndex, err = this_.checkElseIfStatement(); err != nil || startIndex != endIndex {
			return
		}
		if startIndex, endIndex, err = this_.checkElseStatement(); err != nil || startIndex != endIndex {
			return
		}
		if startIndex, endIndex, err = this_.checkIfEndStatement(); err != nil || startIndex != endIndex {
			return
		}
		break
	case *ElseIfSqlStatement:
		if startIndex, endIndex, err = this_.checkElseIfStatement(); err != nil || startIndex != endIndex {
			return
		}
		if startIndex, endIndex, err = this_.checkElseStatement(); err != nil || startIndex != endIndex {
			return
		}
		if startIndex, endIndex, err = this_.checkIfEndStatement(); err != nil || startIndex != endIndex {
			return
		}
		break
	case *ElseSqlStatement:
		if startIndex, endIndex, err = this_.checkIfEndStatement(); err != nil || startIndex != endIndex {
			return
		}
		break
	default:
		if startIndex, endIndex, err = this_.checkIfStatement(); err != nil || startIndex != endIndex {
			return
		}
	}

	return
}

func parseIgnorableSqlStatement(content string) (sqlStatements []SqlStatement, err error) {
	rootSqlStatement := &RootSqlStatement{
		AbstractSqlStatement: &AbstractSqlStatement{},
	}
	var inBracketsLevel int
	var thisStr string
	var last SqlStatement
	var lastParent SqlStatement = rootSqlStatement
	strList := strings.Split(content, "")
	for i := 0; i < len(strList); i++ {
		thisStr = strList[i]
		if thisStr == "[" {
			if lastParent == nil {
				err = errors.New("sql template [" + content + "] parse error")
				return
			}
			inBracketsLevel++
			sqlStatement := &IgnorableSqlStatement{
				AbstractSqlStatement: &AbstractSqlStatement{
					Parent: lastParent,
				},
			}
			last = sqlStatement
			*last.GetParent().GetChildren() = append(*last.GetParent().GetChildren(), sqlStatement)
			lastParent = sqlStatement
		} else if thisStr == "]" {
			if last == nil || inBracketsLevel == 0 {
				err = errors.New("sql template [" + content + "] parse error, has more “[”")
				return
			}
			inBracketsLevel--
			lastParent = lastParent.GetParent()
			last = nil
		} else {
			if last == nil {
				textSqlStatement := &TextSqlStatement{
					AbstractSqlStatement: &AbstractSqlStatement{
						Parent: lastParent,
					},
				}
				last = textSqlStatement
				*last.GetParent().GetChildren() = append(*last.GetParent().GetChildren(), textSqlStatement)
			}
			*last.GetContent() += thisStr
		}

	}
	var list []SqlStatement
	for _, one := range rootSqlStatement.Children {
		list, err = parseExpressionStatement(*one.GetContent())
		if err != nil {
			return
		}
		switch one.(type) {
		case *IgnorableSqlStatement:
			*one.GetContent() = ""
			*one.GetChildren() = list
			sqlStatements = append(sqlStatements, one)
		default:
			sqlStatements = append(sqlStatements, list...)
		}
	}
	//fmt.Println(this_.Sql)
	return
}

func parseExpressionStatement(content string) (sqlStatements []SqlStatement, err error) {
	rootSqlStatement := &RootSqlStatement{
		AbstractSqlStatement: &AbstractSqlStatement{},
	}
	var inBracketsLevel int
	var thisStr string
	var last SqlStatement
	var lastParent SqlStatement = rootSqlStatement
	strList := strings.Split(content, "")
	for i := 0; i < len(strList); i++ {
		thisStr = strList[i]
		if thisStr == "{" {
			if lastParent == nil {
				err = errors.New("sql template [" + content + "] parse error")
				return
			}
			inBracketsLevel++
			sqlStatement := &ExpressionStatement{
				AbstractSqlStatement: &AbstractSqlStatement{
					Parent: lastParent,
				},
			}
			last = sqlStatement
			*last.GetParent().GetChildren() = append(*last.GetParent().GetChildren(), sqlStatement)
			lastParent = sqlStatement
		} else if thisStr == "}" {
			if last == nil || inBracketsLevel == 0 {
				err = errors.New("sql template [" + content + "] parse error, has more “[”")
				return
			}
			inBracketsLevel--
			lastParent = lastParent.GetParent()
			last = nil
		} else {
			if last == nil {
				textSqlStatement := &TextSqlStatement{
					AbstractSqlStatement: &AbstractSqlStatement{
						Parent: lastParent,
					},
				}
				last = textSqlStatement
				*last.GetParent().GetChildren() = append(*last.GetParent().GetChildren(), textSqlStatement)
			}
			*last.GetContent() += thisStr
		}

	}
	sqlStatements = rootSqlStatement.Children
	//fmt.Println(this_.Sql)
	return
}
