package dialect

var (
	postgresqlIndexTypeList []*IndexTypeInfo
)

func appendPostgresqlIndexType(indexType *IndexTypeInfo) {
	postgresqlIndexTypeList = append(postgresqlIndexTypeList, indexType)
}

func init() {

}
