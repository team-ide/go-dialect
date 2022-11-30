package dialect

func NewMappingDM() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeDM,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendDmSql(mapping)

	for _, one := range dmColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range dmIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
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
