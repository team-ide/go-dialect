package dialect

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
