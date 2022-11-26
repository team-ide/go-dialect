package dialect

import "strings"

var (
	shenTongColumnTypeList []*ColumnTypeInfo
)

func appendShenTongColumnType(columnType *ColumnTypeInfo) {
	shenTongColumnTypeList = append(shenTongColumnTypeList, columnType)
}
func init() {
	// -128 到 127
	appendShenTongColumnType(&ColumnTypeInfo{Name: "TINYINT", Format: "TINYINT", IsNumber: true})
	// -2^31 到 2^31-1
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT", Format: "INT", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT4", Format: "INT4", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMERIC($l, $d)", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "DECIMAL", Format: "DECIMAL", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "SERIAL", Format: "SERIAL", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "BPCHAR", Format: "BPCHAR($l)", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "BINARY", Format: "BINARY($l)", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "VARBINARY", Format: "VARBINARY($l)", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "TIME", Format: "TIME", IsDateTime: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			return
		},
	})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT1", Format: "INT1", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT2", Format: "INT2", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT8", Format: "INT8", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "_INT1", Format: "_INT1", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_INT2", Format: "_INT2", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_INT4", Format: "_INT4", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_INT8", Format: "_INT8", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "FLOAT4", Format: "FLOAT4", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "FLOAT8", Format: "FLOAT8", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "_FLOAT4", Format: "_FLOAT4", IsNumber: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_FLOAT8", Format: "_FLOAT8", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "BOOL", Format: "BOOL", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "BOOLEAN", Format: "BOOLEAN", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "BFILE", Format: "BFILE", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_TEXT", Format: "_TEXT", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "ACLITEM", Format: "ACLITEM", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_ACLITEM", Format: "_ACLITEM", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "OID", Format: "OID", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_OID", Format: "_OID", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "TIMESTAMPTZ", Format: "TIMESTAMPTZ", IsString: true})

}

var (
	shenTongIndexTypeList []*IndexTypeInfo
)

func appendShenTongIndexType(indexType *IndexTypeInfo) {
	shenTongIndexTypeList = append(shenTongIndexTypeList, indexType)
}

func init() {
	appendShenTongIndexType(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendShenTongIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendShenTongIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendShenTongIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	appendShenTongIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT",
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	appendShenTongIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL",
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
