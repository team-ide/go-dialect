package dialect

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

func ParseMapping(content string) (mappingSql *MappingSql, err error) {
	mappingParser := &MappingParser{content: content}
	mappingSql, err = mappingParser.Parse()
	if err != nil {
		return
	}
	return
}

type MappingSql struct {
	Content      string
	SqlTemplates map[*MappingSqlType]*MappingSqlTemplate
}

type MappingSqlTemplate struct {
	Content string            `json:"content"`
	Root    *MappingStatement `json:"root"`
}
type MappingStatement struct {
	Content     string              `json:"content"`
	Statements  []*MappingStatement `json:"statements"`
	Parent      *MappingStatement   `json:"-"`
	HasBrackets bool                `json:"hasBrackets"`
}

func (this_ *MappingSqlTemplate) parse() (err error) {
	content := strings.TrimSpace(this_.Content)

	this_.Root = &MappingStatement{}

	var inBracketsLevel int
	var thisChar byte
	var lastStatement *MappingStatement

	for i := 0; i < len(content); i++ {
		thisChar = content[i]
		if thisChar == '[' {
			inBracketsLevel++
			statement := &MappingStatement{
				HasBrackets: true,
			}
			if lastStatement == nil {
				statement.Parent = this_.Root
			} else {
				statement.Parent = lastStatement
			}
			lastStatement = statement
			lastStatement.Parent.Statements = append(lastStatement.Parent.Statements, statement)
		} else if thisChar == ']' {
			if lastStatement == nil || inBracketsLevel == 0 {
				err = errors.New("sql template [" + content + "] parse error, has more “[”")
				return
			}
			inBracketsLevel--
			lastStatement = lastStatement.Parent
		} else {
			if lastStatement == nil {
				statement := &MappingStatement{
					Parent: this_.Root,
				}
				lastStatement = statement
				this_.Root.Statements = append(this_.Root.Statements, statement)
			}
			lastStatement.Content += string(thisChar)
		}

	}

	//fmt.Println(this_.Sql)
	return
}

type MappingSqlType string

var (
	MappingSqlTypeOwnerCreate  = appendMappingSqlType("OwnerCreateSql")
	MappingSqlTypeOwnersSelect = appendMappingSqlType("OwnersSelectSql")
	MappingSqlTypeOwnerSelect  = appendMappingSqlType("OwnerSelectSql")
	MappingSqlTypeOwnerDelete  = appendMappingSqlType("OwnerDeleteSql")

	MappingSqlTypeTableCreate       = appendMappingSqlType("TableCreateSql")
	MappingSqlTypeTableCreateColumn = appendMappingSqlType("TableCreateColumnSql")
	MappingSqlTypeTablesSelect      = appendMappingSqlType("TablesSelectSql")
	MappingSqlTypeTableSelect       = appendMappingSqlType("TableSelectSql")
	MappingSqlTypeTableDelete       = appendMappingSqlType("TableDeleteSql")
	MappingSqlTypeTableComment      = appendMappingSqlType("TableCommentSql")
	MappingSqlTypeTableRename       = appendMappingSqlType("TableRenameSql")

	MappingSqlTypeColumnAdd     = appendMappingSqlType("ColumnAddSql")
	MappingSqlTypeColumnsSelect = appendMappingSqlType("ColumnsSelectSql")
	MappingSqlTypeColumnDelete  = appendMappingSqlType("ColumnDeleteSql")
	MappingSqlTypeColumnComment = appendMappingSqlType("ColumnCommentSql")
	MappingSqlTypeColumnRename  = appendMappingSqlType("ColumnRenameSql")
	MappingSqlTypeColumnUpdate  = appendMappingSqlType("ColumnUpdateSql")

	MappingSqlTypePrimaryKeyAdd     = appendMappingSqlType("PrimaryKeyAddSql")
	MappingSqlTypePrimaryKeysSelect = appendMappingSqlType("PrimaryKeysSelectSql")
	MappingSqlTypePrimaryKeyDelete  = appendMappingSqlType("PrimaryKeyDeleteSql")

	MappingSqlTypeIndexAdd      = appendMappingSqlType("IndexAddSql")
	MappingSqlTypeIndexesSelect = appendMappingSqlType("IndexesSelectSql")
	MappingSqlTypeIndexDelete   = appendMappingSqlType("IndexDeleteSql")

	MappingSqlTypes []*MappingSqlType
)

func (this_ *MappingSqlType) isStart(line string) bool {
	line = strings.ReplaceAll(line, " ", "")
	line = strings.ToLower(line)
	if strings.Contains(line, strings.ToLower(string(*this_)+"start")) {
		return true
	}
	return false
}
func (this_ *MappingSqlType) isEnd(line string) bool {
	line = strings.ReplaceAll(line, " ", "")
	line = strings.ToLower(line)
	if strings.Contains(line, strings.ToLower(string(*this_)+"end")) {
		return true
	}
	return false
}

func appendMappingSqlType(mappingSqlType string) *MappingSqlType {
	res := MappingSqlType(mappingSqlType)
	MappingSqlTypes = append(MappingSqlTypes, &res)
	return &res
}

type MappingParser struct {
	content string
}

func (this_ *MappingParser) Parse() (mappingSql *MappingSql, err error) {
	mappingSql = &MappingSql{
		SqlTemplates: make(map[*MappingSqlType]*MappingSqlTemplate),
	}
	reader := strings.NewReader(this_.content)
	buf := bufio.NewReader(reader)
	var line string
	var lastMappingSqlType *MappingSqlType
	var lastSqlTemplate *MappingSqlTemplate
	for {
		line, err = buf.ReadString('\n')
		if err != nil && err != io.EOF {
			err = errors.New("sql template read error," + err.Error())
			return
		}
		if lastMappingSqlType != nil {
			if lastMappingSqlType.isEnd(line) {
				err = lastSqlTemplate.parse()
				if err != nil {
					err = errors.New("sql template parse error," + err.Error())
					return
				}
				lastMappingSqlType = nil
				lastSqlTemplate = nil
				continue
			}
		}
		var isStart bool
		for _, one := range MappingSqlTypes {
			if isStart = one.isStart(line); isStart {
				lastMappingSqlType = one
				lastSqlTemplate = &MappingSqlTemplate{}
				mappingSql.SqlTemplates[lastMappingSqlType] = lastSqlTemplate
				break
			}
		}
		if isStart {
			continue
		}
		if lastSqlTemplate != nil {
			lastSqlTemplate.Content = lastSqlTemplate.Content + line
		}
		if err == io.EOF { //读取结束，会报EOF
			err = nil
			break
		}
	}
	return
}
