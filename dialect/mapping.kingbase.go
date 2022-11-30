package dialect

import (
	"fmt"
	"strconv"
	"strings"
)

func NewMappingKingBase() (mapping *SqlMapping) {

	// http://www.yaotu.net/biancheng/21946.html
	// https://www.modb.pro/db/442114
	// https://help.kingbase.com.cn/v8/index.html
	mapping = &SqlMapping{
		dialectType: TypeKingBase,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendKingBaseSql(mapping)

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

	for _, one := range kingBaseColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range kingBaseIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}

var (
	kingBaseIndexTypeList []*IndexTypeInfo
)

func appendKingBaseIndexType(indexType *IndexTypeInfo) {
	kingBaseIndexTypeList = append(kingBaseIndexTypeList, indexType)
}

func init() {
	appendKingBaseIndexType(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendKingBaseIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendKingBaseIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendKingBaseIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	appendKingBaseIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	appendKingBaseIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
