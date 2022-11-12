package dialect

func NewMappingOracle() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeOracle,
	}

	return
}
