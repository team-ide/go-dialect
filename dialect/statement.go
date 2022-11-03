package dialect

import (
	"errors"
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
	curIndex     int
	curStr       string
	curStatement SqlStatement
	curParent    SqlStatement
	curRow       int
	curCol       int
	bracketLevel int // “{}” 层级
	braceLevel   int // “[]” 层级
}

func (this_ *SqlStatementParser) reset() {

	this_.contentLen = len(this_.content)
	this_.curRow = 0
	this_.curCol = 0
	this_.curIndex = 0
	this_.curStr = ""
	this_.curParent = nil
	this_.curStatement = nil
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
	switch char {
	case "\n":
		this_.curRow++
		this_.curCol = 0
		break
	case "[":
		this_.bracketLevel++
		break
	case "]":
		this_.bracketLevel--
		break
	case "{":
		str := strings.TrimSpace(this_.curStr)
		if strings.HasPrefix(str, "if ") ||
			strings.HasPrefix(strings.ReplaceAll(str, " ", ""), "if(") {
			if strings.HasSuffix(str, "{") {

			} else {

			}
		}
		this_.braceLevel++
		break
	case "}":
		this_.braceLevel--
		break
	default:
		this_.curStr += char
		break

	}

	return
}

func (this_ *SqlStatementParser) next() (err error) {
	this_.curIndex++
	err = this_.parseStr()
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
}
