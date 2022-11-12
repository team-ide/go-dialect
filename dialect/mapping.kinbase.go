package dialect

func NewMappingKinBase() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeKinBase,
	}

	return
}
