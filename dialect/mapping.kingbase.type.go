package dialect

import "strings"

var (
	kingBaseColumnTypeList []*ColumnTypeInfo
)

func appendKingBaseColumnType(columnType *ColumnTypeInfo) {
	kingBaseColumnTypeList = append(kingBaseColumnTypeList, columnType)
}
func init() {
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER", IsNumber: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "FLOAT", Format: "FLOAT", IsNumber: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMERIC($l, $d)", IsNumber: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "DOUBLE", Format: "DOUBLE", IsNumber: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "DOUBLE PRECISION", Format: "DOUBLE PRECISION", IsNumber: true})

	appendKingBaseColumnType(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CHARACTER", Format: "CHARACTER($l)", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CHARACTER VARYING", Format: "CHARACTER VARYING($l)", IsString: true})

	appendKingBaseColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			return
		},
	})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "TIMESTAMP WITHOUT TIME ZONE", Format: "TIMESTAMP WITHOUT TIME ZONE", IsDateTime: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "INTERVAL", Format: "INTERVAL", IsDateTime: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "INTERVAL DAY TO SECOND", Format: "INTERVAL DAY TO SECOND", IsDateTime: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "INTERVAL YEAR TO MONTH", Format: "INTERVAL YEAR TO MONTH", IsDateTime: true})

	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "BYTEA", Format: "BYTEA", IsString: true})

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
