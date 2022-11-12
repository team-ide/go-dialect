package dialect

func NewMappingSqlite() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeSqlite,
	}

	return
}
