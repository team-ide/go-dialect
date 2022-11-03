package dialect

import (
	"encoding/json"
	"fmt"
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

	sqlStatement, err := GetSqlStatement(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("sql-template:", sqlStatement.GetTemplate())
	for _, node := range *sqlStatement.GetChildren() {
		fmt.Println("children-sql-template:", node.GetTemplate())
		bs, _ := json.Marshal(node)
		fmt.Println("children:", string(bs))
	}
}
func TestSqlStatementParser(t *testing.T) {
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

	sqlStatement, err := sqlStatementParser(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("sql-template:", sqlStatement.GetTemplate())
	for _, node := range *sqlStatement.GetChildren() {
		bs, _ := json.Marshal(node)
		fmt.Println("children:", string(bs))
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
