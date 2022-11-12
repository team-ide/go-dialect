package dialect

func NewMappingShenTong() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeShenTong,
	}

	return
}
