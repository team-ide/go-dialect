package dialect

import (
	"strconv"
)

func NewMappingGBase() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeGBase,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendGBaseSql(mapping)

	mapping.PackPageSql = func(selectSql string, pageSize int, pageNo int) (pageSql string) {
		pageSql = `SELECT * FROM(SELECT ROWNUM rn,t.* FROM(` + selectSql + `) t WHERE ROWNUM <=` + strconv.Itoa(pageSize*pageNo) + ")"
		pageSql += `WHERE rn>=` + strconv.Itoa(pageSize*(pageNo-1)+1)
		return
	}

	for _, one := range gBaseColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range gBaseIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	mapping.OwnerTablePack = func(param *ParamModel, ownerName string, tableName string) string {

		var res string
		if ownerName != "" {
			res += mapping.dialect.OwnerNamePack(param, ownerName) + ":"
		}
		if tableName != "" {
			res += mapping.dialect.TableNamePack(param, tableName)
		}
		return res
	}

	return
}

var (
	gBaseIndexTypeList []*IndexTypeInfo
)

func appendGBaseIndexType(indexType *IndexTypeInfo) {
	gBaseIndexTypeList = append(gBaseIndexTypeList, indexType)
}

func init() {
	appendGBaseIndexType(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendGBaseIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendGBaseIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendGBaseIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	appendGBaseIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	appendGBaseIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
