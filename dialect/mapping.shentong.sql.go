package dialect

func appendShenTongSql(mapping *SqlMapping) {

	appendOracleSql(mapping)

	mapping.OwnerCreate = `
CREATE USER {ownerName} WITH PASSWORD {sqlValuePack(ownerPassword)};
`
	mapping.OwnerDelete = `
DROP USER {ownerName} cascade;
`
}
