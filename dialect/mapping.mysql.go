package dialect

func NewMappingMysql() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeMysql,

		OwnerNamePackChar:  "`",
		TableNamePackChar:  "`",
		ColumnNamePackChar: "`",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "\\",
	}

	appendMysqlSql(mapping)

	for _, one := range mysqlColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range mysqlIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}
