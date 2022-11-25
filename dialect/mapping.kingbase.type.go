package dialect

var (
	kingBaseColumnTypeList []*ColumnTypeInfo
)

func appendKingBaseColumnType(columnType *ColumnTypeInfo) {
	kingBaseColumnTypeList = append(kingBaseColumnTypeList, columnType)
}
func init() {

}
