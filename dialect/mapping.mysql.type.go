package dialect

var (
	mysqlIndexTypeList []*IndexTypeInfo
)

func appendMysqlIndexType(indexType *IndexTypeInfo) {
	mysqlIndexTypeList = append(mysqlIndexTypeList, indexType)
}

func init() {
	appendMysqlIndexType(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"TEXT"},
	})
	appendMysqlIndexType(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"TEXT"},
	})
	appendMysqlIndexType(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX"})
	appendMysqlIndexType(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"TEXT"},
	})
	appendMysqlIndexType(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", OnlySupportDataTypes: []string{"CHAR", "VARCHAR", "TEXT"}})
	appendMysqlIndexType(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", OnlySupportDataTypes: []string{"GEOMETRY", "POINT", "LINESTRING", "POLYGON"}})
}
