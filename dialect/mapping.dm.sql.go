package dialect

func appendDmSql(mapping *SqlMapping) {

	appendOracleSql(mapping)

	mapping.OwnerCreate = `
CREATE USER {doubleQuotationMarksPack(ownerName)} IDENTIFIED BY {doubleQuotationMarksPack(ownerPassword)};
GRANT DBA TO {doubleQuotationMarksPack(ownerName)};
`

}
