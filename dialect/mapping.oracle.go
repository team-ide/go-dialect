package dialect

import (
	"strconv"
	"strings"
)

func NewMappingOracle() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeOracle,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	mapping.IndexNameMaxLen = 30

	appendOracleSql(mapping)

	mapping.PackPageSql = func(selectSql string, pageSize int, pageNo int) (pageSql string) {
		pageSql = `SELECT * FROM(SELECT ROWNUM rn,t.* FROM(` + selectSql + `) t WHERE ROWNUM <=` + strconv.Itoa(pageSize*pageNo) + ")"
		pageSql += `WHERE rn>=` + strconv.Itoa(pageSize*(pageNo-1)+1)
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

	for _, one := range oracleColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range oracleIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}
