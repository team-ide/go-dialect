package dialect

import "fmt"

func NewMappingPostgresql() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypePostgresql,
	}

	mapping.PackPageSql = func(selectSql string, pageSize int, pageNo int) (pageSql string) {
		pageSql = selectSql + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, pageSize*(pageNo-1))
		return
	}
	return
}
