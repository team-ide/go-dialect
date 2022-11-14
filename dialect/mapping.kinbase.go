package dialect

func NewMappingKinBase() (mapping *SqlMapping) {

	// http://www.yaotu.net/biancheng/21946.html
	// https://www.modb.pro/db/442114

	mapping = NewMappingOracle()
	mapping.dialectType = TypeKinBase

	mapping.OwnerCreate = `
CREATE USER {ownerName} WITH PASSWORD {sqlValuePack(ownerPassword)};
CREATE SCHEMA {ownerName} AUTHORIZATION {ownerName};
`
	mapping.OwnerDelete = `
DROP SCHEMA {ownerName} CASCADE;
DROP OWNED BY {ownerName} CASCADE;
`
	mapping.OwnerNamePackChar = ""
	return
}
