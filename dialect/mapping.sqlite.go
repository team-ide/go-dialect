package dialect

func NewMappingSqlite() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeSqlite,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendSqliteSql(mapping)

	for _, one := range sqliteColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range sqliteIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}
