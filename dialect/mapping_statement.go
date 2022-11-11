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

func (this_ *SqlMappingStatement) NewStatementContext(param *ParamModel, dataList ...interface{}) (statementContext *StatementContext, err error) {
	statementContext = NewStatementContext()
	if this_.SqlMapping != nil {
		if this_.SqlMapping.MethodCache != nil {
			for name, method := range this_.SqlMapping.MethodCache {
				statementContext.AddMethod(name, method)
			}
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

func (this_ *SqlMappingStatement) FormatSql(statement Statement, param *ParamModel, dataList ...interface{}) (sqlList []string, err error) {
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

func (this_ *SqlMappingStatement) OwnersSelectSql(param *ParamModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.OwnersSelect, param)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) OwnerSelectSql(param *ParamModel, owner *OwnerModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.OwnerSelect, param, owner)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) OwnerCreateSql(param *ParamModel, owner *OwnerModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.OwnerCreate, param, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) OwnerDeleteSql(param *ParamModel, owner *OwnerModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.OwnerDelete, param, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) TablesSelectSql(param *ParamModel, owner *OwnerModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.TablesSelect, param, owner)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) TableSelectSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.TableSelect, param, owner, table)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) TableCreateSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableCreate, param, owner, table)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) TableRenameSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableRename, param, owner, table)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) TableDeleteSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableDelete, param, owner, table)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) TableCommentSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableComment, param, owner, table)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) ColumnsSelectSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.ColumnsSelect, param, owner, table)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) ColumnSelectSql(param *ParamModel, owner *OwnerModel, table *TableModel, column *ColumnModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.ColumnSelect, param, owner, table, column)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) ColumnAddSql(param *ParamModel, owner *OwnerModel, table *TableModel, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnAdd, param, owner, table, column)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) ColumnUpdateSql(param *ParamModel, owner *OwnerModel, table *TableModel, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnUpdate, param, owner, table, column)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) ColumnDeleteSql(param *ParamModel, owner *OwnerModel, table *TableModel, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnDelete, param, owner, table, column)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) ColumnRenameSql(param *ParamModel, owner *OwnerModel, table *TableModel, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnRename, param, owner, table, column)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) ColumnCommentSql(param *ParamModel, owner *OwnerModel, table *TableModel, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnComment, param, owner, table, column)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) PrimaryKeysSelectSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.PrimaryKeysSelect, param, owner, table)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) PrimaryKeyAddSql(param *ParamModel, owner *OwnerModel, table *TableModel, primaryKey *PrimaryKeyModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.PrimaryKeyAdd, param, owner, table, primaryKey)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) PrimaryKeyDeleteSql(param *ParamModel, owner *OwnerModel, table *TableModel, primaryKey *PrimaryKeyModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.PrimaryKeyDelete, param, owner, table, primaryKey)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) IndexesSelectSql(param *ParamModel, owner *OwnerModel, table *TableModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.IndexesSelect, param, owner, table)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *SqlMappingStatement) IndexAddSql(param *ParamModel, owner *OwnerModel, table *TableModel, index *IndexModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.IndexAdd, param, owner, table, index)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlMappingStatement) IndexDeleteSql(param *ParamModel, owner *OwnerModel, table *TableModel, index *IndexModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.IndexDelete, param, owner, table, index)
	if err != nil {
		return
	}
	return
}
