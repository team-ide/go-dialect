package dialect

import (
	"fmt"
	"strconv"
	"strings"
)

func NewMappingShenTong() (mapping *SqlMapping) {
	// https://blog.csdn.net/asd051377305/article/details/108766792

	mapping = &SqlMapping{
		dialectType: TypeShenTong,

		OwnerNamePackChar:  "",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendShenTongSql(mapping)

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
			variableSql += ":" + strconv.Itoa(i)
			variableSql += strList[i]
		}
		return
	}

	for _, one := range shenTongColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range shenTongIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}
