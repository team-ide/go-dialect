package dialect

type ParamModel struct {
	OwnerNamePack      *bool   `json:"ownerNamePack,omitempty"`
	OwnerNamePackChar  *string `json:"ownerNamePackChar,omitempty"`
	TableNamePack      *bool   `json:"tableNamePack,omitempty"`
	TableNamePackChar  *string `json:"tableNamePackChar,omitempty"`
	ColumnNamePack     *bool   `json:"columnNamePack,omitempty"`
	ColumnNamePackChar *string `json:"columnNamePackChar,omitempty"`

	CustomData map[string]interface{} `json:"customData,omitempty"`
}

type OwnerModel struct {
	OwnerName             string `json:"ownerName,omitempty"`
	OwnerComment          string `json:"ownerComment,omitempty"`
	OwnerPassword         string `json:"ownerPassword,omitempty"`
	OwnerCharacterSetName string `json:"ownerCharacterSetName,omitempty"`
	OwnerCollationName    string `json:"ownerCollationName,omitempty"`

	Error string `json:"error,omitempty"`
}

type TableModel struct {
	TableName    string         `json:"tableName,omitempty"`
	TableComment string         `json:"tableComment,omitempty"`
	ColumnList   []*ColumnModel `json:"columnList,omitempty"`
	IndexList    []*IndexModel  `json:"indexList,omitempty"`

	TableCharacterSetName string `json:"tableCharacterSetName,omitempty"`

	OwnerName string `json:"ownerName,omitempty"`

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
			if one.Name == name {
				return one
			}
		}
	}
	return nil
}

func (this_ *TableModel) AddIndex(models ...*IndexModel) {

	for _, model := range models {
		find := this_.FindIndexByName(model.Name)
		columnNames := model.Columns
		if model.ColumnName != "" && StringsIndex(columnNames, model.ColumnName) < 0 {
			columnNames = append(columnNames, model.ColumnName)
		}
		if find != nil {
			for _, columnName := range columnNames {
				if StringsIndex(find.Columns, columnName) < 0 {
					find.Columns = append(find.Columns, columnName)
				}
			}
		} else {
			model.Columns = columnNames
			this_.IndexList = append(this_.IndexList, model)
		}
	}
}

type ColumnModel struct {
	ColumnName             string `json:"columnName,omitempty"`
	ColumnComment          string `json:"columnComment,omitempty"`
	ColumnType             string `json:"columnType,omitempty"`
	ColumnLength           int    `json:"columnLength,omitempty"`
	ColumnDecimal          int    `json:"columnDecimal,omitempty"`
	ColumnNotNull          bool   `json:"columnNotNull,omitempty"`
	ColumnDefault          string `json:"columnDefault,omitempty"`
	ColumnBeforeColumn     string `json:"columnBeforeColumn,omitempty"`
	ColumnCharacterSetName string `json:"columnCharacterSetName,omitempty"`

	ColumnDefaults                 []string `json:"columnDefaults,omitempty"`
	ColumnDefaultCurrentTimestamp  bool     `json:"columnDefaultCurrentTimestamp"`
	ColumnOnUpdateCurrentTimestamp bool     `json:"columnOnUpdateCurrentTimestamp"`
	ColumnExtra                    string   `json:"columnExtra,omitempty"`
	OwnerName                      string   `json:"ownerName,omitempty"`
	TableName                      string   `json:"tableName,omitempty"`

	Error string `json:"error,omitempty"`
}

type PrimaryKeyModel struct {
	Columns    []string `json:"columns,omitempty"`
	ColumnName string   `json:"columnName,omitempty"`

	OwnerName string `json:"ownerName,omitempty"`
	TableName string `json:"tableName,omitempty"`
	Error     string `json:"error,omitempty"`
}

type IndexModel struct {
	Name       string   `json:"name,omitempty"`
	Type       string   `json:"type,omitempty"`
	ColumnName string   `json:"columnName,omitempty"`
	Columns    []string `json:"columns,omitempty"`
	Comment    string   `json:"comment,omitempty"`

	OwnerName string `json:"ownerName,omitempty"`
	TableName string `json:"tableName,omitempty"`
	Error     string `json:"error,omitempty"`
}
