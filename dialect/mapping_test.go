package dialect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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
ALTER TABLE [{ownerName}.]{tableName} ADD UNIQUE {indexName} ({columnNames}) [COMMENT {indexComment}]
{ else if EqualFold(indexType, 'FULLTEXT') }
ALTER TABLE [{ownerName}.]{tableName} ADD FULLTEXT {indexName} ({columnNames}) [COMMENT {indexComment}]
{ else if indexType == '' }
ALTER TABLE [{ownerName}.]{tableName} ADD INDEX {indexName} ({columnNames}) [COMMENT {indexComment}]
{ else }
ALTER TABLE [{ownerName}.]{tableName} ADD {indexType} {indexName} ({columnNames}) [COMMENT {indexComment}]
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
	context["ownerName"] = "库名"
	context["tableName"] = "表名"
	context["indexName"] = "索引名称"
	context["columnNames"] = "字段1,字段2"
	context["indexComment"] = "索引注释"
	context["indexType"] = "索引注释"

	context["EqualFold"] = reflect.ValueOf(StringEqualFold)
	text, err := sqlStatement.Format(context)
	if err != nil {
		panic(err)
	}
	fmt.Println("sql-text:", text)
}

func StringEqualFold(arg1 interface{}, arg2 interface{}) (equal bool) {
	str1 := GetStringValue(arg1)
	str2 := GetStringValue(arg2)
	equal = strings.EqualFold(str1, str2)
	return
}

func TestSplitOperator(t *testing.T) {
	res, err := splitOperator("a+a-s/d")
	if err != nil {
		panic(err)
	}
	for _, one := range res {
		fmt.Println(one)
	}

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

func method1() string {
	return "a"
}

type methodObject struct {
}

func (this_ *methodObject) method1() string {
	return "a"
}
func TestMethod(t *testing.T) {
	methodValue := reflect.ValueOf(method1)
	res := methodValue.Call([]reflect.Value{})
	println("call method1 result:", res)

	obj := &methodObject{}
	methodValue = reflect.ValueOf(obj.method1)
	res = methodValue.Call([]reflect.Value{})
	println("call methodObject method1 result:", res)
}
