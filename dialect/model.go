package dialect

type OwnerModel struct {
	Name             string `json:"name,omitempty"`
	Comment          string `json:"comment,omitempty"`
	Password         string `json:"password,omitempty"`
	CharacterSetName string `json:"characterSetName,omitempty"`
	CollationName    string `json:"collationName,omitempty"`
	Error            string `json:"error,omitempty"`
}

type TableModel struct {
	Name       string         `json:"name,omitempty"`
	Comment    string         `json:"comment,omitempty"`
	ColumnList []*ColumnModel `json:"columnList,omitempty"`
	IndexList  []*IndexModel  `json:"indexList,omitempty"`

	CharacterSetName string `json:"characterSetName,omitempty"`

	OwnerName string `json:"ownerName,omitempty"`
	Sql       string `json:"sql,omitempty"`
	Error     string `json:"error,omitempty"`
}

func (this_ *TableModel) AddColumn(column *ColumnModel) *ColumnModel {
	this_.ColumnList = append(this_.ColumnList, column)
	return nil
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
	Name             string `json:"name,omitempty"`
	Comment          string `json:"comment,omitempty"`
	Type             string `json:"type,omitempty"`
	Length           int    `json:"length,omitempty"`
	Decimal          int    `json:"decimal,omitempty"`
	PrimaryKey       bool   `json:"primaryKey,omitempty"`
	NotNull          bool   `json:"notNull,omitempty"`
	Default          string `json:"default,omitempty"`
	BeforeColumn     string `json:"beforeColumn,omitempty"`
	CharacterSetName string `json:"characterSetName,omitempty"`

	Defaults                 []string `json:"defaults,omitempty"`
	DefaultCurrentTimestamp  bool     `json:"defaultCurrentTimestamp"`
	OnUpdateCurrentTimestamp bool     `json:"onUpdateCurrentTimestamp"`
	Extra                    string   `json:"extra,omitempty"`
	OwnerName                string   `json:"ownerName,omitempty"`
	TableName                string   `json:"tableName,omitempty"`
	Error                    string   `json:"error,omitempty"`
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
