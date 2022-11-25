package dialect

var (
	postgresqlColumnTypeList []*ColumnTypeInfo
)

func appendPostgresqlColumnType(columnType *ColumnTypeInfo) {
	postgresqlColumnTypeList = append(postgresqlColumnTypeList, columnType)
}
func init() {

}

var (
	postgresqlIndexTypeList []*IndexTypeInfo
)

func appendPostgresqlIndexType(indexType *IndexTypeInfo) {
	postgresqlIndexTypeList = append(postgresqlIndexTypeList, indexType)
}

func init() {

}
