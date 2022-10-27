package worker

import (
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/vitess/sqlparser"
	"io"
	"strings"
)

func NewConvertParser(srcSql string, dest dialect.Dialect, param *dialect.GenerateParam) *ConvertParser {
	if param == nil {
		param = &dialect.GenerateParam{}
	}
	return &ConvertParser{
		srcSql: srcSql,
		dest:   dest,
		param:  param,
	}
}

type ConvertParser struct {
	srcSql  string
	dest    dialect.Dialect
	param   *dialect.GenerateParam
	sqlList []string
	destSql string
}

func (this_ *ConvertParser) GetSrcSql() (srcSql string) {
	srcSql = this_.srcSql
	return
}

func (this_ *ConvertParser) GetSqlList() (sqlList []string) {
	sqlList = this_.sqlList
	return
}

func (this_ *ConvertParser) GetDestSql() (destSql string) {
	destSql = this_.destSql
	return
}

func GetStatements(sqlInfo string) (stmts []sqlparser.Statement, err error) {
	tokens := sqlparser.NewStringTokenizer(sqlInfo)
	for {
		var stmt sqlparser.Statement
		stmt, err = sqlparser.ParseNext(tokens)

		if err != nil {
			if err != io.EOF {
				return
			}
			err = nil
		}
		if stmt == nil {
			break
		}
		stmts = append(stmts, stmt)
	}
	return
}

func (this_ *ConvertParser) Parse() (err error) {
	stmts, err := GetStatements(this_.srcSql)
	if err != nil {
		return
	}

	for _, stmt := range stmts {
		err = this_.parse(stmt)
		if err != nil {
			return
		}
	}
	destSql := ""
	for _, sqlOne := range this_.sqlList {
		destSql += sqlOne + ";\n"
	}
	this_.destSql = destSql
	return
}
func (this_ *ConvertParser) parse(stmt_ sqlparser.Statement) (err error) {
	switch stmt := stmt_.(type) {
	case *sqlparser.Select:
		return parseSelect(stmt)
	case *sqlparser.Insert:
		insert, err := parseInsert(stmt)
		if err != nil {
			return err
		}
		sqlList, err := this_.dest.InsertSql(this_.param, insert)
		if err != nil {
			return err
		}
		this_.sqlList = append(this_.sqlList, sqlList...)
	case *sqlparser.Update:
		return parseUpdate(stmt)
	case *sqlparser.Delete:
		return parseDelete(stmt)
	case *sqlparser.CreateTable:
		databaseName, table, err := parseCreateTable(stmt)
		if err != nil {
			return err
		}
		sqlList, err := this_.dest.TableCreateSql(this_.param, databaseName, table)
		if err != nil {
			return err
		}
		this_.sqlList = append(this_.sqlList, sqlList...)
	}
	return
}
func parseCreateTable(stmt *sqlparser.CreateTable) (databaseName string, table *dialect.TableModel, err error) {
	table = &dialect.TableModel{}
	databaseName = stmt.GetTable().Qualifier.String()
	table.Name = stmt.GetTable().Name.String()

	tableSpec := stmt.GetTableSpec()
	if tableSpec == nil {
		return
	}
	for _, tableSpecColumn := range tableSpec.Columns {
		if tableSpecColumn == nil {
			continue
		}
		column := &dialect.ColumnModel{}
		column.Name = tableSpecColumn.Name.String()
		column.Type = tableSpecColumn.Type.Type
		if tableSpecColumn.Type.Length != nil {
			column.Length, _ = dialect.StringToInt(tableSpecColumn.Type.Length.Val)
		}
		if tableSpecColumn.Type.Scale != nil {
			column.Decimal, _ = dialect.StringToInt(tableSpecColumn.Type.Scale.Val)
		}
		if tableSpecColumn.Type.Options != nil {
			if tableSpecColumn.Type.Options.Comment != nil {
				column.Comment = tableSpecColumn.Type.Options.Comment.Val
			}
			if tableSpecColumn.Type.Options.Null != nil {
				column.NotNull = true
			}
			if tableSpecColumn.Type.Options.Default != nil {
				buf := sqlparser.NewTrackedBuffer(nil)
				tableSpecColumn.Type.Options.Default.Format(buf)
				column.Default = buf.String()
				column.Default = strings.TrimLeft(column.Default, "'")
				column.Default = strings.TrimRight(column.Default, "'")
				column.Default = strings.TrimLeft(column.Default, "\"")
				column.Default = strings.TrimRight(column.Default, "\"")
			}
		}
		table.AddColumn(column)
	}

	for _, tableSpecIndex := range tableSpec.Indexes {
		if tableSpecIndex == nil || tableSpecIndex.Info == nil {
			continue
		}
		if tableSpecIndex.Info.Primary {
			for _, indexColumn := range tableSpecIndex.Columns {
				table.AddPrimaryKey(&dialect.PrimaryKeyModel{ColumnName: indexColumn.Column.String()})
			}
			continue
		}

		for _, indexColumn := range tableSpecIndex.Columns {
			index := &dialect.IndexModel{}
			index.Name = tableSpecIndex.Info.Name.String()
			index.ColumnName = indexColumn.Column.String()
			for _, option := range tableSpecIndex.Options {
				if option == nil || option.Value == nil {
					continue
				}
				if option.Name == "COMMENT" {
					index.Comment = option.Value.Val
				}
			}
			if tableSpecIndex.Info.Unique {
				index.Type = "unique"
			}

			table.AddIndex(index)
		}
	}
	for _, option := range tableSpec.Options {
		if option == nil || option.Value == nil {
			continue
		}
		if option.Name == "COMMENT" {
			table.Comment = option.Value.Val
		}
	}

	//var bs []byte
	//bs, _ = json.MarshalIndent(stmt, "", "  ")
	//fmt.Println(string(bs))
	//
	//buf := sqlparser.NewTrackedBuffer(nil)
	//stmt.Format(buf)
	//sqlInfo := buf.String()
	//bs, _ = json.MarshalIndent(table, "", "  ")
	//fmt.Println("createTable sql:")
	//fmt.Println(sqlInfo)
	//fmt.Println("createTable table:")
	//fmt.Println(string(bs))
	return
}

func parseSelect(stmt *sqlparser.Select) (err error) {

	return
}
func parseInsert(stmt *sqlparser.Insert) (insert *dialect.InsertModel, err error) {
	insert = &dialect.InsertModel{}
	insert.OwnerName = stmt.Table.Qualifier.String()
	insert.TableName = stmt.Table.Name.String()
	for _, c := range stmt.Columns {
		name := c.CompliantName()
		insert.Columns = append(insert.Columns, name)
	}
	switch rows := stmt.Rows.(type) {
	case sqlparser.Values:
		for _, row := range rows {
			var insertRow []*dialect.ValueModel
			for _, rowV := range row {
				//fmt.Println("row v:", reflect.TypeOf(rowV).String())
				value := &dialect.ValueModel{}
				insertRow = append(insertRow, value)
				switch v := rowV.(type) {
				case *sqlparser.Literal:
					value.Value = v.Val
					if v.Type == sqlparser.StrVal {
						value.Type = dialect.ValueTypeString
					} else if v.Type == sqlparser.IntVal {
						value.Type = dialect.ValueTypeNumber
					} else {
						panic("parseInsert not support value type [" + value.Type + "]")
					}
					break
				case *sqlparser.FuncExpr:
					value.Type = dialect.ValueTypeFunc
					buf := sqlparser.NewTrackedBuffer(nil)
					v.Format(buf)
					value.Value = buf.String()
					break
				}

			}
			insert.Rows = append(insert.Rows, insertRow)
		}

	}

	//var bs []byte
	//bs, _ = json.MarshalIndent(stmt, "", "  ")
	//fmt.Println(string(bs))
	return
}
func parseUpdate(stmt *sqlparser.Update) (err error) {

	return
}
func parseDelete(stmt *sqlparser.Delete) (err error) {

	return
}

func parseCreateDatabase(stmt *sqlparser.CreateDatabase) (err error) {

	return
}
func parseDropTable(stmt *sqlparser.DropTable) (err error) {

	return
}
func parseDropDatabase(stmt *sqlparser.DropDatabase) (err error) {

	return
}
func parseDropColumn(stmt *sqlparser.DropColumn) (err error) {

	return
}
func parseDropKey(stmt *sqlparser.DropKey) (err error) {

	return
}
func parseAddColumns(stmt *sqlparser.AddColumns) (err error) {

	return
}
func parseAlterColumn(stmt *sqlparser.AlterColumn) (err error) {

	return
}
func parseRenameIndex(stmt *sqlparser.RenameIndex) (err error) {

	return
}
func parseAddIndexDefinition(stmt *sqlparser.AddIndexDefinition) (err error) {

	return
}
