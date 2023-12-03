package dialect

type ParamModel struct {
	OwnerNamePack      *bool   `json:"ownerNamePack"`
	OwnerNamePackChar  *string `json:"ownerNamePackChar"`
	TableNamePack      *bool   `json:"tableNamePack"`
	TableNamePackChar  *string `json:"tableNamePackChar"`
	ColumnNamePack     *bool   `json:"columnNamePack"`
	ColumnNamePackChar *string `json:"columnNamePackChar"`
	SqlValuePackChar   *string `json:"sqlValuePackChar"`
	SqlValueEscapeChar *string `json:"sqlValueEscapeChar"`

	AppendSqlValue *bool `json:"appendSqlValue"`

	CustomData map[string]interface{} `json:"customData"`
}

type OwnerModel struct {
	OwnerName             string `json:"ownerName"`
	OwnerComment          string `json:"ownerComment,omitempty"`
	OwnerPassword         string `json:"ownerPassword,omitempty"`
	OwnerCharacterSetName string `json:"ownerCharacterSetName,omitempty"`
	OwnerCollationName    string `json:"ownerCollationName,omitempty"`

	Error string `json:"error,omitempty"`
}

type TableModel struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`

	ColumnList            []*ColumnModel `json:"columnList"`
	IndexList             []*IndexModel  `json:"indexList"`
	PrimaryKeys           []string       `json:"primaryKeys"`
	TableCharacterSetName string         `json:"tableCharacterSetName,omitempty"`
	TableCollationName    string         `json:"tableCollationName,omitempty"`

	OwnerName string `json:"ownerName"`

	Sql   string `json:"sql,omitempty"`
	Error string `json:"error,omitempty"`
}

func (this_ *TableModel) AddColumn(column *ColumnModel) *ColumnModel {
	this_.ColumnList = append(this_.ColumnList, column)
	return nil
}
func (this_ *TableModel) FindColumnByName(name string) *ColumnModel {
	if len(this_.ColumnList) > 0 {
		for _, one := range this_.ColumnList {
			if one.ColumnName == name {
				return one
			}
		}
	}
	return nil
}

func (this_ *TableModel) FindIndexByName(name string) *IndexModel {
	if len(this_.IndexList) > 0 {
		for _, one := range this_.IndexList {
			if one.IndexName == name {
				return one
			}
		}
	}
	return nil
}

func (this_ *TableModel) AddPrimaryKey(models ...*PrimaryKeyModel) {

	for _, model := range models {
		if StringsIndex(this_.PrimaryKeys, model.ColumnName) >= 0 {
			continue
		}
		this_.PrimaryKeys = append(this_.PrimaryKeys, model.ColumnName)
		find := this_.FindColumnByName(model.ColumnName)
		if find != nil {
			find.PrimaryKey = true
		}
	}
}
func (this_ *TableModel) AddIndex(models ...*IndexModel) {

	for _, model := range models {
		var find *IndexModel
		if model.IndexName != "" {
			find = this_.FindIndexByName(model.IndexName)
		}
		columnNames := model.ColumnNames
		if model.ColumnName != "" && StringsIndex(columnNames, model.ColumnName) < 0 {
			columnNames = append(columnNames, model.ColumnName)
		}
		if find != nil {
			for _, columnName := range columnNames {
				if StringsIndex(find.ColumnNames, columnName) < 0 {
					find.ColumnNames = append(find.ColumnNames, columnName)
				}
			}
		} else {
			model.ColumnNames = columnNames
			this_.IndexList = append(this_.IndexList, model)
		}
	}
}

type ColumnModel struct {
	ColumnName     string `json:"columnName"`
	ColumnComment  string `json:"columnComment"`
	ColumnDataType string `json:"columnDataType"`
	//ColumnType             string `json:"columnType"`
	ColumnLength           int    `json:"columnLength"`
	ColumnPrecision        int    `json:"columnPrecision"`
	ColumnScale            int    `json:"columnScale"`
	ColumnNotNull          bool   `json:"columnNotNull"`
	ColumnDefault          string `json:"columnDefault"`
	ColumnAfterColumn      string `json:"columnAfterColumn"`
	ColumnCharacterSetName string `json:"columnCharacterSetName"`

	PrimaryKey bool `json:"primaryKey"`

	ColumnEnums []string `json:"columnEnums"`
	ColumnExtra string   `json:"columnExtra"`
	OwnerName   string   `json:"ownerName"`
	TableName   string   `json:"tableName"`

	Error string `json:"error,omitempty"`
}

type ColumnTypeInfo struct {
	Name         string `json:"name,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Format       string `json:"format,omitempty"`
	MinLength    *int   `json:"minLength"`
	MaxLength    *int   `json:"maxLength"`
	MinPrecision *int   `json:"minPrecision"`
	MaxPrecision *int   `json:"maxPrecision"`
	MinScale     *int   `json:"minScale"`
	MaxScale     *int   `json:"maxScale"`

	// IsNumber 如果 是 数字 数据存储 设置该属性
	IsNumber  bool `json:"isNumber,omitempty"`
	IsInteger bool `json:"isInteger,omitempty"`
	IsFloat   bool `json:"isFloat,omitempty"`

	// IsString 如果 是 字符串 数据存储 设置该属性
	IsString bool `json:"isString,omitempty"`

	// IsDateTime 如果 是 日期时间 数据存储 设置该属性
	IsDateTime bool `json:"isDateTime,omitempty"`

	// IsBytes 如果 是 流 数据存储 设置该属性
	IsBytes bool `json:"isBytes,omitempty"`

	IsBoolean bool `json:"isBoolean,omitempty"`

	// IsEnum 如果 是 枚举 数据存储 设置该属性
	IsEnum bool `json:"isEnum,omitempty"`

	// IsExtend 如果 非 当前 数据库能支持的类型 设置该属性
	IsExtend bool     `json:"isExtend,omitempty"`
	Matches  []string `json:"matches"`

	IfNotFound             bool                                                                               `json:"ifNotFound,omitempty"`
	ColumnDefaultPack      func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) `json:"-"`
	ColumnTypePack         func(column *ColumnModel) (columnTypePack string, err error)                       `json:"-"`
	SqlValuePack           func(value string) (sqlValue string)                                               `json:"-"`
	FullColumnByColumnType func(columnType string, column *ColumnModel) (err error)                           `json:"-"`
}

type IndexTypeInfo struct {
	Name   string `json:"name,omitempty"`
	Format string `json:"format,omitempty"`

	// IsExtend 如果 非 当前 数据库能支持的类型 设置该属性
	IsExtend bool `json:"isExtend,omitempty"`

	OnlySupportDataTypes []string `json:"onlySupportDataTypes"`
	NotSupportDataTypes  []string `json:"notSupportDataTypes"`

	IndexTypeFormat func(index *IndexModel) (indexTypeFormat string, err error)                                                        `json:"-"`
	IndexNameFormat func(param *ParamModel, ownerName string, tableName string, index *IndexModel) (indexNameFormat string, err error) `json:"-"`
}

type PrimaryKeyModel struct {
	ColumnName string `json:"columnName"`

	OwnerName string `json:"ownerName"`
	TableName string `json:"tableName"`
	Error     string `json:"error,omitempty"`
}

type IndexModel struct {
	IndexName    string   `json:"indexName"`
	IndexType    string   `json:"indexType"`
	ColumnName   string   `json:"columnName"`
	ColumnNames  []string `json:"columnNames"`
	IndexComment string   `json:"indexComment"`

	OwnerName string `json:"ownerName"`
	TableName string `json:"tableName"`
	Error     string `json:"error,omitempty"`
}
