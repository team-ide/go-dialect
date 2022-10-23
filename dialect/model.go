package dialect

type DatabaseModel struct {
	Name                 string `json:"name,omitempty"`
	Comment              string `json:"comment,omitempty"`
	CatalogName          string `json:"catalogName,omitempty"`
	DefaultCharacterSet  string `json:"defaultCharacterSet,omitempty"`
	DefaultCollationName string `json:"defaultCollationName,omitempty"`
	SchemaName           string `json:"schemaName,omitempty"`
	Error                string `json:"error,omitempty"`
}

type TableModel struct {
	Name       string         `json:"name,omitempty"`
	Comment    string         `json:"comment,omitempty"`
	ColumnList []*ColumnModel `json:"columnList,omitempty"`
	IndexList  []*IndexModel  `json:"indexList,omitempty"`

	TableCatalog string `json:"tableCatalog,omitempty"`
	TableSchema  string `json:"tableSchema,omitempty"`
	TableType    string `json:"tableType,omitempty"`
	Sql          string `json:"sql,omitempty"`
	Error        string `json:"error,omitempty"`
}

func (this_ *TableModel) FindColumnByName(name string) *ColumnModel {
	if len(this_.ColumnList) > 0 {
		for _, one := range this_.ColumnList {
			if one.Name == name {
				return one
			}
		}
	}
	return nil
}
func (this_ *TableModel) FindColumnByOldName(oldName string) *ColumnModel {
	if len(this_.ColumnList) > 0 {
		for _, one := range this_.ColumnList {
			if one.OldName == oldName {
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
func (this_ *TableModel) FindIndexByOldName(oldName string) *IndexModel {
	if len(this_.IndexList) > 0 {
		for _, one := range this_.IndexList {
			if one.OldName == oldName {
				return one
			}
		}
	}
	return nil
}

func (this_ *TableModel) AddPrimaryKey(models ...*PrimaryKeyModel) {
	for _, model := range models {
		find := this_.FindColumnByName(model.ColumnName)
		if find != nil {
			find.PrimaryKey = true
		}
	}
}

func (this_ *TableModel) AddIndex(models ...*IndexModel) {
	for _, model := range models {
		find := this_.FindIndexByName(model.ColumnName)
		if find != nil {
			find.Columns = append(find.Columns, model.ColumnName)
		} else {
			this_.IndexList = append(this_.IndexList, model)
		}
	}
}

type ColumnModel struct {
	Name             string      `json:"name,omitempty"`
	Comment          string      `json:"comment,omitempty"`
	Type             string      `json:"type,omitempty"`
	Length           int         `json:"length,omitempty"`
	Decimal          int         `json:"decimal,omitempty"`
	PrimaryKey       bool        `json:"primaryKey,omitempty"`
	NotNull          bool        `json:"notNull,omitempty"`
	Default          string      `json:"default,omitempty"`
	OldName          string      `json:"oldName,omitempty"`
	OldComment       string      `json:"oldComment,omitempty"`
	OldType          string      `json:"oldType,omitempty"`
	OldLength        int         `json:"oldLength,omitempty"`
	OldDecimal       int         `json:"oldDecimal,omitempty"`
	OldPrimaryKey    bool        `json:"oldPrimaryKey,omitempty"`
	OldNotNull       bool        `json:"oldNotNull,omitempty"`
	OldDefault       interface{} `json:"oldDefault,omitempty"`
	BeforeColumn     string      `json:"beforeColumn,omitempty"`
	Deleted          bool        `json:"deleted,omitempty"`
	CharacterSetName string      `json:"characterSetName,omitempty"`
	TableCatalog     string      `json:"tableCatalog,omitempty"`
	TableSchema      string      `json:"tableSchema,omitempty"`
	TableName        string      `json:"tableName,omitempty"`
	Error            string      `json:"error,omitempty"`
}

type PrimaryKeyModel struct {
	Columns      []string `json:"columns,omitempty"`
	ColumnName   string   `json:"columnName,omitempty"`
	TableCatalog string   `json:"tableCatalog,omitempty"`
	TableSchema  string   `json:"tableSchema,omitempty"`
	TableName    string   `json:"tableName,omitempty"`
	Error        string   `json:"error,omitempty"`
}

type IndexModel struct {
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	ColumnName   string   `json:"columnName,omitempty"`
	Columns      []string `json:"columns,omitempty"`
	Comment      string   `json:"comment,omitempty"`
	OldName      string   `json:"oldName,omitempty"`
	OldComment   string   `json:"oldComment,omitempty"`
	OldType      string   `json:"oldType,omitempty"`
	OldColumns   []string `json:"oldColumns,omitempty"`
	Deleted      bool     `json:"deleted,omitempty"`
	TableCatalog string   `json:"tableCatalog,omitempty"`
	TableSchema  string   `json:"tableSchema,omitempty"`
	TableName    string   `json:"tableName,omitempty"`
	Error        string   `json:"error,omitempty"`
}
