package dialect

func NewMappingPostgresql() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypePostgresql,
	}

	return
}
