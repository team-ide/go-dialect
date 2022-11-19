package dialect

import "fmt"

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

	mapping.PackPageSql = func(selectSql string, pageSize int, pageNo int) (pageSql string) {
		pageSql = selectSql + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, pageSize*(pageNo-1))
		return
	}

	return
}
