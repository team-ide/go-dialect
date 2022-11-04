package dialect

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func GetSqlStatement(content string) (sqlStatement *RootSqlStatement, err error) {
	content = strings.TrimSpace(content)

	sqlStatement = &RootSqlStatement{
		AbstractSqlStatement: &AbstractSqlStatement{},
	}

	var inBracketsLevel int
	var thisStr string
	var last SqlStatement
	var lastParent SqlStatement = sqlStatement
	strList := strings.Split(content, "")
	for i := 0; i < len(strList); i++ {
		thisStr = strList[i]
		if thisStr == "[" {
			if lastParent == nil {
				err = errors.New("sql template [" + content + "] parse error")
				return
			}
			inBracketsLevel++
			ignorableSqlStatement := &IgnorableSqlStatement{
				AbstractSqlStatement: &AbstractSqlStatement{
					Parent: lastParent,
				},
			}
			last = ignorableSqlStatement
			*last.GetParent().GetChildren() = append(*last.GetParent().GetChildren(), ignorableSqlStatement)
			lastParent = ignorableSqlStatement
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

	//fmt.Println(this_.Sql)
	return
}

func sqlStatementParser(content string) (sqlStatement *RootSqlStatement, err error) {
	parser := &SqlStatementParser{
		content: strings.Split(content, ""),
	}
	sqlStatement, err = parser.parse()
	return
}

type SqlStatementParser struct {
	content      []string
	contentLen   int
	curIndex     int //当前索引
	curRowStart  int //当前行开始索引
	curStr       string
	curStatement SqlStatement
	curParent    SqlStatement
	curIf        *IfSqlStatement
	curElseIf    *ElseIfSqlStatement
	curElse      *ElseSqlStatement
	curRow       int // 当前行
	curCol       int // 当前列
	bracketLevel int // “{}” 层级
	braceLevel   int // “[]” 层级
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
	this_.curIf = nil
	this_.curElseIf = nil
	this_.curElse = nil
	this_.bracketLevel = 0
	this_.braceLevel = 0
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
	this_.curStr += char
	err = this_.checkStatement()
	if err != nil {
		return
	}
	if char == "\n" {
		this_.curRow++
		this_.curCol = 0
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
	if startIndex, endIndex, err = this_.checkIfStatement(); err != nil || startIndex != endIndex {
		return
	}
	if this_.curIf != nil {
		if startIndex, endIndex, err = this_.checkElseIfStatement(); err != nil || startIndex != endIndex {
			return
		}
		if startIndex, endIndex, err = this_.checkElseStatement(); err != nil || startIndex != endIndex {
			return
		}
		if startIndex, endIndex, err = this_.checkIfEndStatement(); err != nil || startIndex != endIndex {
			return
		}
	}

	return
}

func getIfCondition(str string) (condition string, err error) {
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
	return
}
func (this_ *SqlStatementParser) checkIfStatement() (startIndex int, endIndex int, err error) {
	var reg *regexp.Regexp
	var finds [][]int
	var curStr = this_.curStr
	// 匹配格式如：
	// { if xxx }
	// { if (xxx) }
	reg, err = regexp.Compile("\\{\\s*if[\\s\\(]+.*\\}")
	if reg == nil || err != nil {
		return
	}
	finds = reg.FindAllStringIndex(curStr, -1)
	if len(finds) > 0 {

		startIndex = finds[0][0]
		endIndex = finds[0][1]
		text := curStr[0:startIndex]
		if text != "" {
			textSqlStatement := &TextSqlStatement{
				AbstractSqlStatement: &AbstractSqlStatement{
					Content: text,
					Parent:  this_.curParent,
				},
			}
			*this_.curParent.GetChildren() = append(*this_.curParent.GetChildren(), textSqlStatement)
		}
		ifSqlStatement := &IfSqlStatement{
			AbstractSqlStatement: &AbstractSqlStatement{
				Content: curStr[startIndex:endIndex],
				Parent:  this_.curParent,
			},
		}
		ifSqlStatement.Condition, err = getIfCondition(ifSqlStatement.Content)
		if err != nil {
			return
		}

		*this_.curParent.GetChildren() = append(*this_.curParent.GetChildren(), ifSqlStatement)
		this_.curStatement = ifSqlStatement
		this_.curParent = ifSqlStatement
		this_.curIf = ifSqlStatement
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
	reg, err = regexp.Compile("\\{\\s*else\\s*if[\\s\\(]+.+\\s*\\}")
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
			textSqlStatement := &TextSqlStatement{
				AbstractSqlStatement: &AbstractSqlStatement{
					Content: text,
					Parent:  this_.curParent,
				},
			}
			if this_.curElseIf != nil {
				this_.curElseIf.Children = append(this_.curElseIf.Children, textSqlStatement)
			} else {
				this_.curIf.Children = append(this_.curIf.Children, textSqlStatement)
				fmt.Println("curIf Children:", len(this_.curIf.Children))
			}
		}

		elseIfSqlStatement := &ElseIfSqlStatement{
			AbstractSqlStatement: &AbstractSqlStatement{
				Content: curStr[startIndex:endIndex],
				Parent:  this_.curIf,
			},
			If:    this_.curIf,
			Index: len(this_.curIf.ElseIfs),
		}
		elseIfSqlStatement.Condition, err = getIfCondition(elseIfSqlStatement.Content)
		if err != nil {
			return
		}
		this_.curIf.ElseIfs = append(this_.curIf.ElseIfs, elseIfSqlStatement)
		this_.curStatement = elseIfSqlStatement
		this_.curParent = elseIfSqlStatement
		this_.curElseIf = elseIfSqlStatement
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
	reg, err = regexp.Compile("\\{\\s*else\\s*\\}")
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
			textSqlStatement := &TextSqlStatement{
				AbstractSqlStatement: &AbstractSqlStatement{
					Content: text,
					Parent:  this_.curParent,
				},
			}
			if this_.curElseIf != nil {
				this_.curElseIf.Children = append(this_.curElseIf.Children, textSqlStatement)
			} else {
				this_.curIf.Children = append(this_.curIf.Children, textSqlStatement)
			}
		}

		elseSqlStatement := &ElseSqlStatement{
			AbstractSqlStatement: &AbstractSqlStatement{
				Content: curStr[startIndex:endIndex],
				Parent:  this_.curIf,
			},
			If: this_.curIf,
		}

		this_.curIf.Else = elseSqlStatement
		this_.curStatement = elseSqlStatement
		this_.curParent = elseSqlStatement
		this_.curElse = elseSqlStatement
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
	reg, err = regexp.Compile("\\{\\s*\\}")
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
			textSqlStatement := &TextSqlStatement{
				AbstractSqlStatement: &AbstractSqlStatement{
					Content: text,
					Parent:  this_.curParent,
				},
			}
			if this_.curElse != nil {
				this_.curElse.Children = append(this_.curElse.Children, textSqlStatement)
			} else if this_.curElseIf != nil {
				this_.curElseIf.Children = append(this_.curElseIf.Children, textSqlStatement)
			} else {
				this_.curIf.Children = append(this_.curIf.Children, textSqlStatement)
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

type SqlStatement interface {
	GetTemplate() (template string)
	GetContent() (content *string)
	GetParent() (parent SqlStatement)
	GetChildren() (children *[]SqlStatement)
}

type AbstractSqlStatement struct {
	Content  string         `json:"content,omitempty"`
	Children []SqlStatement `json:"children,omitempty"`
	Parent   SqlStatement   `json:"-"`
}

func (this_ *AbstractSqlStatement) GetContent() (content *string) {
	content = &this_.Content
	return
}

func (this_ *AbstractSqlStatement) GetParent() (parent SqlStatement) {
	parent = this_.Parent
	return
}
func (this_ *AbstractSqlStatement) GetChildren() (children *[]SqlStatement) {
	children = &this_.Children
	return
}
func (this_ *AbstractSqlStatement) GetTemplate() (template string) {
	template += this_.Content
	for _, one := range this_.Children {
		template += one.GetTemplate()
	}
	return
}

type RootSqlStatement struct {
	*AbstractSqlStatement
}

type IgnorableSqlStatement struct {
	*AbstractSqlStatement
}

func (this_ *IgnorableSqlStatement) GetTemplate() (template string) {
	template += "["

	template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	template += "]"
	return
}

type TextSqlStatement struct {
	*AbstractSqlStatement
}

type IfSqlStatement struct {
	*AbstractSqlStatement
	Condition string                `json:"condition"`
	ElseIfs   []*ElseIfSqlStatement `json:"elseIfs"`
	Else      *ElseSqlStatement     `json:"else"`
}

func (this_ *IfSqlStatement) GetTemplate() (template string) {
	template += "if " + this_.Condition + " {\n"

	//template += this_.Content
	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	if len(this_.ElseIfs) > 0 || this_.Else != nil {
		template += "\n}"
	} else {
		template += "\n}\n"
	}
	for _, one := range this_.ElseIfs {
		template += one.GetTemplate()
	}

	if this_.Else != nil {
		template += this_.Else.GetTemplate()
	}

	return
}

type ElseIfSqlStatement struct {
	*AbstractSqlStatement
	Condition string          `json:"condition"`
	If        *IfSqlStatement `json:"-"`
	Index     int             `json:"index"`
}

func (this_ *ElseIfSqlStatement) GetTemplate() (template string) {
	template += " else if " + this_.Condition + " {\n"

	//template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	if this_.IsEndElseIf() {
		if this_.If.Else != nil {
			template += "\n}"
		} else {
			template += "\n}\n"
		}

	} else {
		template += "\n}"
	}
	return
}
func (this_ *ElseIfSqlStatement) IsEndElseIf() (isEnd bool) {
	if this_.Index == len(this_.If.ElseIfs)-1 {
		isEnd = true
	}
	return
}

type ElseSqlStatement struct {
	*AbstractSqlStatement
	If *IfSqlStatement `json:"-"`
}

func (this_ *ElseSqlStatement) GetTemplate() (template string) {
	template += " else {\n"

	//template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	template += "\n}\n"
	return
}
