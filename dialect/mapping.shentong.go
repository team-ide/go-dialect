package dialect

func NewMappingShenTong() (mapping *SqlMapping) {
	// https://blog.csdn.net/asd051377305/article/details/108766792
	
	mapping = NewMappingOracle()
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
