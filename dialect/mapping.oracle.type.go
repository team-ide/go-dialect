package dialect

import "strings"

var (
	oracleColumnTypeList []*ColumnTypeInfo
)

func appendOracleColumnType(columnType *ColumnTypeInfo) {
	oracleColumnTypeList = append(oracleColumnTypeList, columnType)
}
func init() {
	appendOracleColumnType(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER", IsNumber: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "FLOAT", Format: "FLOAT", IsNumber: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "BINARY_FLOAT", Format: "BINARY_FLOAT", IsNumber: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "BINARY_DOUBLE", Format: "BINARY_DOUBLE", IsNumber: true})

	appendOracleColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "NCHAR", Format: "NCHAR($l)", IsString: true})

	appendOracleColumnType(&ColumnTypeInfo{Name: "VARCHAR2", Format: "VARCHAR2($l)", IsString: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "NVARCHAR2", Format: "NVARCHAR2($l)", IsString: true})

	appendOracleColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			//if strings.Contains(strings.ToLower(column.ColumnExtra), "on update current_timestamp") {
			//	columnDefaultPack += " ON UPDATE CURRENT_TIMESTAMP"
			//}
			return
		},
	})
	appendOracleColumnType(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB", IsString: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "NCLOB", Format: "NCLOB", IsString: true})

	appendOracleColumnType(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "BFILE", Format: "BFILE", IsString: true})

	appendOracleColumnType(&ColumnTypeInfo{Name: "ROWID", Format: "ROWID", IsString: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "UROWID", Format: "UROWID", IsString: true})
	//
	appendOracleColumnType(&ColumnTypeInfo{Name: "RAW", Format: "RAW($l)", IsString: true})
	appendOracleColumnType(&ColumnTypeInfo{Name: "LONG", Format: "LONG", IsString: true})

}

var (
	oracleIndexTypeList []*IndexTypeInfo
)

func appendOracleIndexType(indexType *IndexTypeInfo) {
	oracleIndexTypeList = append(oracleIndexTypeList, indexType)
}

func init() {
	appendOracleIndexType(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendOracleIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendOracleIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	appendOracleIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	appendOracleIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	appendOracleIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
