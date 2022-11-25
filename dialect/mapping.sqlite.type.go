package dialect

var (
	sqliteColumnTypeList []*ColumnTypeInfo
)

func appendSqliteColumnType(columnType *ColumnTypeInfo) {
	sqliteColumnTypeList = append(sqliteColumnTypeList, columnType)
}
func init() {

}
