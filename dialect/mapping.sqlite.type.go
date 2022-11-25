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
