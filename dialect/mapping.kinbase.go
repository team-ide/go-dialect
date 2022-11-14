package dialect

func NewMappingKinBase() (mapping *SqlMapping) {
	mapping = NewMappingOracle()
	mapping.dialectType = TypeKinBase

	mapping.dialectType = TypeShenTong
	mapping.OwnerCreate = `
CREATE USER {ownerName} WITH PASSWORD {sqlValuePack(ownerPassword)};
`
	mapping.OwnerDelete = `
DROP USER {ownerName} cascade;
`
	mapping.OwnerNamePackChar = ""
	return
}
