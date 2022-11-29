package dialect

func NewMappingDM() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeDM,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendDmSql(mapping)

	for _, one := range dmColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range dmIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}
