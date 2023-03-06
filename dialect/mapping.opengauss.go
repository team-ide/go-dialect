package dialect

import (
	"fmt"
	"strconv"
	"strings"
)

func NewMappingOpenGauss() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeOpenGauss,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendOpenGaussSql(mapping)

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

	for _, one := range openGaussColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range openGaussIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}

var (
	openGaussIndexTypeList []*IndexTypeInfo
)

func appendOpenGaussIndexType(indexType *IndexTypeInfo) {
	openGaussIndexTypeList = append(openGaussIndexTypeList, indexType)
}

func init() {
	appendOpenGaussIndexType(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendOpenGaussIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendOpenGaussIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendOpenGaussIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	appendOpenGaussIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	appendOpenGaussIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
