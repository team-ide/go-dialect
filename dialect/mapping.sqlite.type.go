package dialect

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
