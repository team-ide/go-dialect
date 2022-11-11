package dialect

import (
	"errors"
	"reflect"
	"strings"
)

var (
	rootStatementType = reflect.TypeOf(&RootStatement{})
)

func NewSqlMappingStatement(mapping *SqlMapping) (sqlMappingStatement *SqlMappingStatement, err error) {
	if mapping == nil {
		err = errors.New("sql mapping is null")
		return
	}
	sqlMappingStatement = &SqlMappingStatement{
		SqlMapping: mapping,
	}
	mappingValue := reflect.ValueOf(mapping).Elem()
	sqlStatementValue := reflect.ValueOf(sqlMappingStatement).Elem()
	sqlStatementType := reflect.TypeOf(sqlMappingStatement).Elem()
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

func (this_ *SqlMappingStatement) NewStatementContext() (statementContext *StatementContext) {
	statementContext = NewStatementContext()
	if this_.SqlMapping != nil {
		if this_.SqlMapping.MethodCache != nil {
			for name, method := range this_.SqlMapping.MethodCache {
				statementContext.AddMethod(name, method)
			}
		}
	}
	return
}

func (this_ *SqlMappingStatement) OwnersSelectSql() (sqlInfo string, err error) {
	statementContext := this_.NewStatementContext()
	sqlInfo, err = this_.OwnersSelect.Format(statementContext)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) OwnerSelectSql(ownerName string) (sqlInfo string, err error) {
	statementContext := this_.NewStatementContext()
	statementContext.SetData("ownerName", ownerName)
	sqlInfo, err = this_.OwnerSelect.Format(statementContext)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) OwnerCreateSql(owner *OwnerModel) (sqlInfo string, err error) {
	statementContext := this_.NewStatementContext()

	err = statementContext.SetJSONData(owner)
	if err != nil {
		return
	}

	sqlInfo, err = this_.OwnerCreate.Format(statementContext)
	if err != nil {
		return
	}
	return
}
