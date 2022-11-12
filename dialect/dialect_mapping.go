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

	TablesSelect *RootStatement
	TableSelect  *RootStatement
	TableCreate  *RootStatement
	TableDelete  *RootStatement
	TableComment *RootStatement
	TableRename  *RootStatement

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

	IndexesSelect   *RootStatement
	IndexAdd        *RootStatement
	IndexDelete     *RootStatement
	IndexNameFormat *RootStatement
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

func (this_ *mappingDialect) NewStatementContext(param *ParamModel, dataList ...interface{}) (statementContext *StatementContext, err error) {
	statementContext = NewStatementContext()
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

	ownerName, _ := statementContext.GetData("ownerName")
	if ownerName != nil && ownerName != "" {
		statementContext.SetData("ownerNamePack", this_.OwnerNamePack(param, ownerName.(string)))
	}
	tableName, _ := statementContext.GetData("tableName")
	if tableName != nil && tableName != "" {
		statementContext.SetData("tableNamePack", this_.TableNamePack(param, tableName.(string)))
	}
	oldTableName, _ := statementContext.GetData("oldTableName")
	if oldTableName != nil && oldTableName != "" {
		statementContext.SetData("oldTableNamePack", this_.TableNamePack(param, oldTableName.(string)))
	}
	newTableName, _ := statementContext.GetData("newTableName")
	if newTableName != nil && newTableName != "" {
		statementContext.SetData("newTableNamePack", this_.TableNamePack(param, newTableName.(string)))
	}
	columnName, _ := statementContext.GetData("columnName")
	if columnName != nil && columnName != "" {
		statementContext.SetData("columnNamePack", this_.ColumnNamePack(param, columnName.(string)))
	}
	oldColumnName, _ := statementContext.GetData("oldColumnName")
	if oldColumnName != nil && oldColumnName != "" {
		statementContext.SetData("oldColumnNamePack", this_.ColumnNamePack(param, oldColumnName.(string)))
	}
	newColumnName, _ := statementContext.GetData("newColumnName")
	if newColumnName != nil && newColumnName != "" {
		statementContext.SetData("newColumnNamePack", this_.ColumnNamePack(param, newColumnName.(string)))
	}
	columnNamesStr, _ := statementContext.GetData("columnNamesStr")
	if columnNamesStr != nil && columnNamesStr != "" {
		statementContext.SetData("columnNamesStrPack", this_.ColumnNamesStrPack(param, columnNamesStr.(string)))
	}
	columnNames, _ := statementContext.GetData("columnNames")
	if columnNames != nil {
		columnNamesList, ok := columnNames.([]string)
		if ok {
			statementContext.SetData("columnNamesPack", this_.ColumnNamesPack(param, columnNamesList))
		}
	}
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
	return
}
