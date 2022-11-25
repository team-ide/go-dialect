package dialect

var (
	sqliteColumnTypeList []*ColumnTypeInfo
)

func appendSqliteColumnType(columnType *ColumnTypeInfo) {
	sqliteColumnTypeList = append(sqliteColumnTypeList, columnType)
}
func init() {
	appendSqliteColumnType(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER($l)", IsNumber: true})
	appendSqliteColumnType(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT($l)", IsString: true})
	appendSqliteColumnType(&ColumnTypeInfo{Name: "NONE", Format: "NONE", IsString: true})
	appendSqliteColumnType(&ColumnTypeInfo{Name: "REAL", Format: "REAL($l, $d)", IsNumber: true})
	appendSqliteColumnType(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMERIC($l, $d)", IsNumber: true, IsDateTime: true})

}

var (
	sqliteIndexTypeList []*IndexTypeInfo
)

func appendSqliteIndexType(indexType *IndexTypeInfo) {
	sqliteIndexTypeList = append(sqliteIndexTypeList, indexType)
}

func init() {

	appendSqliteIndexType(&IndexTypeInfo{Name: "", Format: "INDEX"})
	appendSqliteIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX"})
	appendSqliteIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX", IsExtend: true})
	appendSqliteIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	appendSqliteIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	appendSqliteIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
