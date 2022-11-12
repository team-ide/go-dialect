package dialect

func NewMappingDaMen() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeDaMen,
	}

	return
}
