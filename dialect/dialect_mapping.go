package dialect

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func NewMappingDialect(mapping *SqlMapping) (res *mappingDialect, err error) {
	res = &mappingDialect{
		SqlMapping: mapping,
	}

	err = res.init()
	if err != nil {
		res = nil
	}
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

func (this_ *mappingDialect) OwnerNamePack(param *ParamModel, ownerName string) string {
	char := this_.OwnerNamePackChar
	if param != nil {
		if param.OwnerNamePack != nil && !*param.OwnerNamePack {
			char = ""
		} else if param.OwnerNamePackChar != nil {
			char = *param.OwnerNamePackChar
		}
	}
	return packingName(char, ownerName)
}

func (this_ *mappingDialect) TableNamePack(param *ParamModel, tableName string) string {
	char := this_.TableNamePackChar
	if param != nil {
		if param.TableNamePack != nil && !*param.TableNamePack {
			char = ""
		} else if param.TableNamePackChar != nil {
			char = *param.TableNamePackChar
		}
	}
	return packingName(char, tableName)
}

func (this_ *mappingDialect) ColumnNamePack(param *ParamModel, columnName string) string {
	char := this_.ColumnNamePackChar
	if param != nil {
		if param.ColumnNamePack != nil && !*param.ColumnNamePack {
			char = ""
		} else if param.ColumnNamePackChar != nil {
			char = *param.ColumnNamePackChar
		}
	}
	return packingName(char, columnName)
}

func (this_ *mappingDialect) ColumnNamesPack(param *ParamModel, columnNames []string) string {
	char := this_.ColumnNamePackChar
	if param != nil {
		if param.ColumnNamePack != nil && !*param.ColumnNamePack {
			char = ""
		} else if param.ColumnNamePackChar != nil {
			char = *param.ColumnNamePackChar
		}
	}
	return packingNames(char, columnNames)
}

func (this_ *mappingDialect) ColumnNamesStrPack(param *ParamModel, columnNamesStr string) string {
	return this_.ColumnNamesPack(param, strings.Split(columnNamesStr, ","))
}

func (this_ *mappingDialect) SqlSplit(sqlStr string) (sqlList []string) {
	cacheKey := UUID()
	sqlCache := sqlStr
	sqlCache = strings.ReplaceAll(sqlCache, `''`, `|-`+cacheKey+`-|`)
	sqlCache = strings.ReplaceAll(sqlCache, `""`, `|--`+cacheKey+`--|`)

	var list []string
	var beg int

	var inStringLevel int
	var inStringPack byte
	var thisChar byte
	var lastChar byte

	var stringPackChars = []byte{'"', '\''}
	for i := 0; i < len(sqlCache); i++ {
		thisChar = sqlCache[i]
		if i > 0 {
			lastChar = sqlCache[i-1]
		}

		// inStringLevel == 0 表示 不在 字符串 包装 中
		if thisChar == ';' && inStringLevel == 0 {
			if i > 0 {
				list = append(list, sqlCache[beg:i])
			}
			beg = i + 1
		} else {
			packCharIndex := BytesIndex(stringPackChars, thisChar)
			if packCharIndex >= 0 {
				// inStringLevel == 0 表示 不在 字符串 包装 中
				if inStringLevel == 0 {
					inStringPack = stringPackChars[packCharIndex]
					// 字符串包装层级 +1
					inStringLevel++
				} else {
					if thisChar != inStringPack {
					} else if lastChar == '\\' { // 如果有转义符号 类似 “\'”，“\"”
					} else if lastChar == inStringPack {
						// 如果 前一个字符 与字符串包装字符一致
					} else {
						// 字符串包装层级 -1
						inStringLevel--
					}
				}
			}
		}

	}
	list = append(list, sqlCache[beg:])
	for _, sqlOne := range list {
		sqlOne = strings.TrimSpace(sqlOne)
		if sqlOne == "" {
			continue
		}
		sqlOne = strings.ReplaceAll(sqlOne, `|-`+cacheKey+`-|`, `''`)
		sqlOne = strings.ReplaceAll(sqlOne, `|--`+cacheKey+`--|`, `""`)
		sqlList = append(sqlList, sqlOne)
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

func (this_ *mappingDialect) OwnersSelectSql(param *ParamModel) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.OwnersSelect, param)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) OwnerSelectSql(param *ParamModel, ownerName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.OwnerSelect, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) OwnerModel(data map[string]interface{}) (owner *OwnerModel, err error) {
	if data == nil {
		return
	}
	owner = &OwnerModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) OwnerCreateSql(param *ParamModel, owner *OwnerModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.OwnerCreate, param, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) OwnerDeleteSql(param *ParamModel, ownerName string) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.OwnerDelete, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
	)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TablesSelectSql(param *ParamModel, ownerName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.TablesSelect, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) TableSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.TableSelect, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	if data == nil {
		return
	}
	table = &TableModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, table)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TableCreateSql(param *ParamModel, ownerName string, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableCreate, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		table)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TableRenameSql(param *ParamModel, ownerName string, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableRename, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		table)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TableDeleteSql(param *ParamModel, ownerName string, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableDelete, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		table)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) TableCommentSql(param *ParamModel, ownerName string, table *TableModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.TableComment, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		table)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnsSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.ColumnsSelect, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) ColumnSelectSql(param *ParamModel, ownerName string, tableName string, columnName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.ColumnSelect, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		&ColumnModel{
			ColumnName: columnName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) ColumnModel(data map[string]interface{}) (column *ColumnModel, err error) {
	if data == nil {
		return
	}
	column = &ColumnModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, column)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnAddSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnAdd, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		column)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnUpdateSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnUpdate, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		column)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnDeleteSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnDelete, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		column)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnRenameSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnRename, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		column)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) ColumnCommentSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.ColumnComment, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		column)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) PrimaryKeysSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.PrimaryKeysSelect, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
	if data == nil {
		return
	}
	primaryKey = &PrimaryKeyModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, primaryKey)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) PrimaryKeyAddSql(param *ParamModel, ownerName string, tableName string, primaryKey *PrimaryKeyModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.PrimaryKeyAdd, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		primaryKey)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) PrimaryKeyDeleteSql(param *ParamModel, ownerName string, tableName string, primaryKey *PrimaryKeyModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.PrimaryKeyDelete, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		primaryKey)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) IndexesSelectSql(param *ParamModel, ownerName string, tableName string) (sqlInfo string, err error) {
	sqlList, err := this_.FormatSql(this_.IndexesSelect, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
	)
	if err != nil {
		return
	}
	if len(sqlList) > 0 {
		sqlInfo = sqlList[0]
	}
	return
}

func (this_ *mappingDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
	if data == nil {
		return
	}
	index = &IndexModel{}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, index)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) IndexAddSql(param *ParamModel, ownerName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.IndexAdd, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		index)
	if err != nil {
		return
	}
	return
}

func (this_ *mappingDialect) IndexDeleteSql(param *ParamModel, ownerName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	sqlList, err = this_.FormatSql(this_.IndexDelete, param,
		&OwnerModel{
			OwnerName: ownerName,
		},
		&TableModel{
			TableName: tableName,
		},
		index)
	if err != nil {
		return
	}
	return
}
