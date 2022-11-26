package dialect

var (
	gBaseColumnTypeList []*ColumnTypeInfo
)

func appendGBaseColumnType(columnType *ColumnTypeInfo) {
	gBaseColumnTypeList = append(gBaseColumnTypeList, columnType)
}

func init() {
	appendGBaseColumnType(&ColumnTypeInfo{Name: "INT", Format: "INT", IsNumber: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER", IsNumber: true})

	appendGBaseColumnType(&ColumnTypeInfo{Name: "DECIMAL", Format: "DECIMAL($l, $d)", IsNumber: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "SERIAL", Format: "SERIAL", IsNumber: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "FLOAT", Format: "FLOAT", IsNumber: true})

	appendGBaseColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "NCHAR", Format: "NCHAR($l)", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "LVARCHAR", Format: "LVARCHAR($l)", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "NVARCHAR", Format: "NVARCHAR($l)", IsString: true})

	appendGBaseColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "DATETIME", Format: "DATETIME", IsDateTime: true})

	appendGBaseColumnType(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "BYTE", Format: "BYTE", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT", IsString: true})
	appendGBaseColumnType(&ColumnTypeInfo{Name: "BOOLEAN", Format: "BOOLEAN", IsString: true})
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
