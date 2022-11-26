package dialect

import "strings"

var (
	dmColumnTypeList []*ColumnTypeInfo
)

func appendDmColumnType(columnType *ColumnTypeInfo) {
	dmColumnTypeList = append(dmColumnTypeList, columnType)
}

func init() {
	appendDmColumnType(&ColumnTypeInfo{Name: "INT", Format: "INT", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "DOUBLE", Format: "DOUBLE", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "FLOAT", Format: "FLOAT", IsNumber: true})

	appendDmColumnType(&ColumnTypeInfo{Name: "DECIMAL", Format: "DECIMAL($l, $d)", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "DEC", Format: "DEC($l, $d)", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMERIC($l, $d)", IsNumber: true})

	appendDmColumnType(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})

	appendDmColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "TIME", Format: "TIME", IsDateTime: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "DATETIME", Format: "DATETIME", IsDateTime: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			return
		},
	})
	appendDmColumnType(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "BINARY", Format: "BINARY", IsString: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "VARBINARY", Format: "VARBINARY", IsString: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT", IsString: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "LONGVARCHAR", Format: "LONGVARCHAR", IsString: true})
}

var (
	dmIndexTypeList []*IndexTypeInfo
)

func appendDmIndexType(indexType *IndexTypeInfo) {
	dmIndexTypeList = append(dmIndexTypeList, indexType)
}

func init() {
	appendDmIndexType(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendDmIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendDmIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendDmIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	appendDmIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	appendDmIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
