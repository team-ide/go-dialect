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

	statement, err := statementParse(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("sql-template:", statement.GetTemplate())
	fmt.Println("sql-statement:", statement.GetTemplate())
	testOut(statement, 0)
}

func testOut(statement Statement, tab int) {
	if statement == nil {
		return
	}

	cr := *statement.GetChildren()
	*statement.GetChildren() = []Statement{}
	bs, _ := json.Marshal(statement)
	for i := 0; i < tab; i++ {
		fmt.Print("\t")
	}
	fmt.Print(reflect.TypeOf(statement).String() + ":")
	*statement.GetChildren() = cr
	fmt.Println("", string(bs))
	switch data := statement.(type) {
	case *IfStatement:
		for i := 0; i < tab+1; i++ {
			fmt.Print("\t")
		}
		fmt.Println("if condition:")
		testOut(data.ConditionExpression, tab+2)
	case *ElseIfStatement:
		for i := 0; i < tab+1; i++ {
			fmt.Print("\t")
		}
		fmt.Println("else if condition:")
		testOut(data.ConditionExpression, tab+2)
	}
	if statement.GetChildren() != nil {
		for _, node := range *statement.GetChildren() {
			testOut(node, tab+1)
		}
	}
	switch stat := statement.(type) {
	case *IfStatement:
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
{ if EqualFold(indexType, 'UNIQUE') }
ALTER TABLE [{ownerName}.]{tableName} ADD UNIQUE {indexName} ({columnNames}) [COMMENT {indexComment}]
{ else if EqualFold(indexType, 'FULLTEXT') }
ALTER TABLE [{ownerName}.]{tableName} ADD FULLTEXT {indexName} ({columnNames}) [COMMENT {indexComment}]
{ else if indexType == '' }
ALTER TABLE [{ownerName}.]{tableName} ADD INDEX {indexName} ({columnNames}) [COMMENT {indexComment}]
{ else }
ALTER TABLE [{ownerName}.]{tableName} ADD {indexType} {indexName} ({columnNames}) [COMMENT {indexComment}]
{ }
`

	statement, err := statementParse(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("sql-template:", statement.GetTemplate())
	bs, _ := json.Marshal(statement)
	fmt.Println("sql-statement:", string(bs))
	testOut(statement, 0)
	statementContext := NewStatementContext()
	statementContext.SetData("ownerName", "库名")
	statementContext.SetData("tableName", "表名")
	statementContext.SetData("indexName", "索引名称")
	statementContext.SetData("columnNames", "字段1,字段2")
	statementContext.SetData("indexComment", "索引注释")
	statementContext.SetData("indexType", "uniqUe")

	statementContext.AddMethod("EqualFold", StringEqualFold)

	text, err := statement.Format(statementContext)
	if err != nil {
		panic(err)
	}
	fmt.Println("sql-text:", text)
}

func StringEqualFold(arg1 interface{}, arg2 interface{}) (equal bool, err error) {
	str1 := GetStringValue(arg1)
	str2 := GetStringValue(arg2)
	equal = strings.EqualFold(str1, str2)

	//fmt.Println("StringEqualFold str1:", str1, ",str2:", str2, ",equal:", equal)
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
	_, err := NewDialect(TypeMysql.Name)
	if err != nil {
		panic(err)
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
	methodType := reflect.ValueOf(method1)
	fmt.Println(methodType.Type().NumOut())
	methodValue := reflect.ValueOf(method1)
	res := methodValue.Call([]reflect.Value{})
	println("call method1 result:", res)

	obj := &methodObject{}
	methodValue = reflect.ValueOf(obj.method1)
	res = methodValue.Call([]reflect.Value{})
	println("call methodObject method1 result:", res)
}
