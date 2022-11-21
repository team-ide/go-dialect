package dialect

func NewMappingDM() (mapping *SqlMapping) {
	mapping = NewMappingOracle()
	mapping.dialectType = TypeDM

	mapping.OwnerCreate = `
CREATE USER {doubleQuotationMarksPack(ownerName)} IDENTIFIED BY {doubleQuotationMarksPack(ownerPassword)};
GRANT DBA TO {doubleQuotationMarksPack(ownerName)};
`

	mapping.PackPageSql = nil
	mapping.ReplaceSqlVariable = nil
	return
}
