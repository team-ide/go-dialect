package dialect

var (
	postgresqlColumnTypeList []*ColumnTypeInfo
)

func appendPostgresqlColumnType(columnType *ColumnTypeInfo) {
	postgresqlColumnTypeList = append(postgresqlColumnTypeList, columnType)
}
func init() {

}
