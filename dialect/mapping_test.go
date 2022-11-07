package dialect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestMappingSqlTemplate(t *testing.T) {
	content := `
{ if EqualFold(indexType, 'UNIQUE') {}
ALTER TABLE [{ownerName}.]{tableName} ADD UNIQUE {indexName} ({columnNames}) [COMMENT {columnComment}]
{ } else if EqualFold(indexType, 'FULLTEXT') {}
ALTER TABLE [{ownerName}.]{tableName} ADD FULLTEXT {indexName} ({columnNames}) [COMMENT {columnComment}]
{ } else if indexType == '' {}
ALTER TABLE [{ownerName}.]{tableName} ADD INDEX {indexName} ({columnNames}) [COMMENT {columnComment}]
{ } else {}
ALTER TABLE [{ownerName}.]{tableName} ADD {indexType} {indexName} ({columnNames}) [COMMENT {columnComment}]
{ }}
`

	sqlStatement, err := sqlStatementParse(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("sql-template:", sqlStatement.GetTemplate())
	fmt.Println("sql-statement:", sqlStatement.GetTemplate())
	testOut(sqlStatement, 0)
}

func testOut(sqlStatement SqlStatement, tab int) {
	if sqlStatement == nil {
		return
	}

	cr := *sqlStatement.GetChildren()
	*sqlStatement.GetChildren() = []SqlStatement{}
	bs, _ := json.Marshal(sqlStatement)
	for i := 0; i < tab; i++ {
		fmt.Print("\t")
	}
	fmt.Print(reflect.TypeOf(sqlStatement).String() + ":")
	*sqlStatement.GetChildren() = cr
	fmt.Println("", string(bs))
	switch data := sqlStatement.(type) {
	case *IfSqlStatement:
		for i := 0; i < tab+1; i++ {
			fmt.Print("\t")
		}
		fmt.Println("if condition:")
		testOut(data.ConditionExpression, tab+2)
	case *ElseIfSqlStatement:
		for i := 0; i < tab+1; i++ {
			fmt.Print("\t")
		}
		fmt.Println("else if condition:")
		testOut(data.ConditionExpression, tab+2)
	}
	if sqlStatement.GetChildren() != nil {
		for _, node := range *sqlStatement.GetChildren() {
			testOut(node, tab+1)
		}
	}
	switch stat := sqlStatement.(type) {
	case *IfSqlStatement:
		for _, one := range stat.ElseIfs {
			testOut(one, tab)
		}
		if stat.Else != nil {
			testOut(stat.Else, tab)
		}
		break
	}

}
func TestSqlStatementParser(t *testing.T) {
	content := `
{ if EqualFold(indexType, 'UNIQUE}') }
ALTER TABLE [{ownerName}.]{tableName} ADD UNIQUE {indexName} ({columnNames}) [COMMENT {columnComment}]
{ else if EqualFold(indexType, 'FULLTEXT') }
ALTER TABLE [{ownerName}.]{tableName} ADD FULLTEXT {indexName} ({columnNames}) [COMMENT {columnComment}]
{ else if indexType == '' }
ALTER TABLE [{ownerName}.]{tableName} ADD INDEX {indexName} ({columnNames}) [COMMENT {columnComment}]
{ else }
ALTER TABLE [{ownerName}.]{tableName} ADD {indexType} {indexName} ({columnNames}) [COMMENT {columnComment}]
{ }
`

	sqlStatement, err := sqlStatementParse(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("sql-template:", sqlStatement.GetTemplate())
	bs, _ := json.Marshal(sqlStatement)
	fmt.Println("sql-statement:", string(bs))
	testOut(sqlStatement, 0)
	context := map[string]interface{}{}
	text, err := sqlStatement.Invoke(context)
	if err != nil {
		panic(err)
	}
	fmt.Println("sql-text:", text)
}

func TestMappingSql(t *testing.T) {

	mappingSql, err := ParseMapping(mappingMySql)
	if err != nil {
		panic(err)
	}

	for key, value := range mappingSql.SqlTemplates {
		fmt.Println(*key, ":sql-template:", value.GetTemplate())
		bs, _ := json.Marshal(value.Children)
		fmt.Println(*key, ":children:", string(bs))
	}
}
