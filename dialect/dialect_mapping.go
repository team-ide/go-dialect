package dialect

import (
	"errors"
	"reflect"
	"strings"
)

func NewMappingDialect(mapping *SqlMapping) (dia Dialect, err error) {
	mappingDia := &mappingDialect{
		SqlMapping: mapping,
	}

	err = mappingDia.init()
	if err != nil {
		return
	}
	dia = mappingDia
	return
}

type mappingDialect struct {
	*SqlMapping

	OwnersSelect *RootStatement
	OwnerSelect  *RootStatement
	OwnerCreate  *RootStatement
	OwnerDelete  *RootStatement

	TablesSelect          *RootStatement
	TableSelect           *RootStatement
	TableCreate           *RootStatement
	TableCreateColumn     *RootStatement
	TableCreatePrimaryKey *RootStatement
	TableDelete           *RootStatement
	TableComment          *RootStatement
	TableRename           *RootStatement

	ColumnsSelect *RootStatement
	ColumnSelect  *RootStatement
	ColumnAdd     *RootStatement
	ColumnDelete  *RootStatement
	ColumnComment *RootStatement
	ColumnRename  *RootStatement
	ColumnUpdate  *RootStatement

	PrimaryKeysSelect *RootStatement
	PrimaryKeyAdd     *RootStatement
	PrimaryKeyDelete  *RootStatement

	IndexesSelect *RootStatement
	IndexAdd      *RootStatement
	IndexDelete   *RootStatement
}

func (this_ *mappingDialect) init() (err error) {
	rootStatementType := reflect.TypeOf(&RootStatement{})

	mappingValue := reflect.ValueOf(this_.SqlMapping).Elem()
	sqlStatementValue := reflect.ValueOf(this_).Elem()
	sqlStatementType := reflect.TypeOf(this_).Elem()
	var statement *RootStatement
	for i := 0; i < sqlStatementValue.NumField(); i++ {
		fieldValue := sqlStatementValue.Field(i)
		fieldType := sqlStatementType.Field(i)
		if fieldType.Type != rootStatementType {
			continue
		}
		mappingField := mappingValue.FieldByName(fieldType.Name)
		if mappingField.Kind() == reflect.Invalid {
			err = errors.New("mapping field [" + fieldType.Name + "] is invalid")
			return
		}
		sqlTemplate := strings.TrimSpace(mappingField.String())
		if len(sqlTemplate) == 0 {
			continue
		}
		statement, err = statementParse(sqlTemplate)
		if err != nil {
			return
		}
		fieldValue.Set(reflect.ValueOf(statement))
	}
	return
}

type StatementScript struct {
	*ParamModel
	Dialect
}

func (this_ StatementScript) sqlValuePack(value interface{}) (res string) {

	res = this_.SqlValuePack(this_.ParamModel, nil, value)
	return
}

func (this_ StatementScript) columnNotNull(columnNotNull interface{}) (res string) {
	if isTrue(columnNotNull) {
		res = "NOT NULL"
	}
	return
}

func (this_ StatementScript) equalFold(arg1 interface{}, arg2 interface{}) bool {
	if arg1 == arg2 {
		return true
	}
	str1 := GetStringValue(arg1)
	str2 := GetStringValue(arg2)
	return strings.EqualFold(str1, str2)
}

func (this_ StatementScript) joins(joinList interface{}, joinObj interface{}) (res string) {
	if joinList == nil {
		return
	}
	list, ok := joinList.([]string)
	if !ok {
		objList := joinList.([]interface{})
		for _, one := range objList {
			list = append(list, GetStringValue(one))
		}
	}
	res = strings.Join(list, GetStringValue(joinObj))
	return
}

func (this_ *mappingDialect) NewStatementContext(param *ParamModel, dataList ...interface{}) (statementContext *StatementContext, err error) {
	statementContext = NewStatementContext()

	statementScript := &StatementScript{
		ParamModel: param,
		Dialect:    this_,
	}
	statementContext.AddMethod("sqlValuePack", statementScript.sqlValuePack)
	statementContext.AddMethod("columnNotNull", statementScript.columnNotNull)
	statementContext.AddMethod("joins", statementScript.joins)
	statementContext.AddMethod("equalFold", statementScript.equalFold)

	if this_.MethodCache != nil {
		for name, method := range this_.MethodCache {
			statementContext.AddMethod(name, method)
		}
	}
	if param != nil {
		err = statementContext.SetJSONData(param)
		if err != nil {
			return
		}
		err = statementContext.SetJSONData(param.CustomData)
		if err != nil {
			return
		}
	}
	for _, data := range dataList {
		err = statementContext.SetJSONData(data)
		if err != nil {
			return
		}
	}

	ownerNamePack := ""
	ownerName, _ := statementContext.GetData("ownerName")
	if ownerName != nil && ownerName != "" {
		ownerNamePack = this_.OwnerNamePack(param, ownerName.(string))
	}
	statementContext.SetData("ownerNamePack", ownerNamePack)

	tableNamePack := ""
	tableName, _ := statementContext.GetData("tableName")
	if tableName != nil && tableName != "" {
		tableNamePack = this_.TableNamePack(param, tableName.(string))
	}
	statementContext.SetData("tableNamePack", tableNamePack)

	newTableNamePack := ""
	newTableName, _ := statementContext.GetData("newTableName")
	if newTableName != nil && newTableName != "" {
		newTableNamePack = this_.TableNamePack(param, newTableName.(string))
	}
	statementContext.SetData("newTableNamePack", newTableNamePack)

	columnNamePack := ""
	columnName, _ := statementContext.GetData("columnName")
	if columnName != nil && columnName != "" {
		columnNamePack = this_.ColumnNamePack(param, columnName.(string))
	}
	statementContext.SetData("columnNamePack", columnNamePack)

	newColumnNamePack := ""
	newColumnName, _ := statementContext.GetData("newColumnName")
	if newColumnName != nil && newColumnName != "" {
		newColumnNamePack = this_.ColumnNamePack(param, newColumnName.(string))
	}
	statementContext.SetData("newColumnNamePack", newColumnNamePack)

	columnNamesPack := ""
	columnNames, _ := statementContext.GetData("columnNames")
	if columnNames != nil {
		list := columnNames.([]interface{})
		var stringList []string
		for _, one := range list {
			stringList = append(stringList, one.(string))
		}
		columnNamesPack = this_.ColumnNamesPack(param, stringList)
	}
	statementContext.SetData("columnNamesPack", columnNamesPack)

	primaryKeysPack := ""
	primaryKeys, _ := statementContext.GetData("primaryKeys")
	if primaryKeys != nil {
		list := primaryKeys.([]interface{})
		var stringList []string
		for _, one := range list {
			stringList = append(stringList, one.(string))
		}
		primaryKeysPack = this_.ColumnNamesPack(param, stringList)
	}
	statementContext.SetData("primaryKeysPack", primaryKeysPack)

	indexNamePack := ""
	indexName, _ := statementContext.GetData("indexName")
	if indexName != nil && indexName != "" {
		indexNamePack = this_.ColumnNamePack(param, indexName.(string))
	}
	statementContext.SetData("indexNamePack", indexNamePack)

	return
}

func (this_ *mappingDialect) FormatSql(statement Statement, param *ParamModel, dataList ...interface{}) (sqlList []string, err error) {
	if statement == nil {
		return
	}
	statementContext, err := this_.NewStatementContext(param, dataList...)
	if err != nil {
		return
	}
	sqlInfo, err := statement.Format(statementContext)
	if err != nil {
		return
	}
	sqlList = this_.SqlSplit(sqlInfo)

	//fmt.Println("FormatSql sql data cache")
	//fmt.Println(statementContext.dataCache)
	//fmt.Println("FormatSql sql list")
	//for _, sqlOne := range sqlList {
	//	fmt.Println("sql:", sqlOne)
	//}
	return
}
