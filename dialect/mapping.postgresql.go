package dialect

import (
	"fmt"
	"strconv"
	"strings"
)

func NewMappingPostgresql() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypePostgresql,
	}

	mapping.PackPageSql = func(selectSql string, pageSize int, pageNo int) (pageSql string) {
		pageSql = selectSql + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, pageSize*(pageNo-1))
		return
	}
	mapping.ReplaceSqlVariable = func(sqlInfo string, args []interface{}) (variableSql string) {
		strList := strings.Split(sqlInfo, "?")
		if len(strList) < 1 {
			variableSql = sqlInfo
			return
		}
		variableSql = strList[0]
		for i := 1; i < len(strList); i++ {
			variableSql += "$" + strconv.Itoa(i)
			variableSql += strList[i]
		}
		return
	}
	return
}
