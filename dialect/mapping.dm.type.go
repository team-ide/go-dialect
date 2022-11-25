package dialect

var (
	dmColumnTypeList []*ColumnTypeInfo
)

func appendDmColumnType(columnType *ColumnTypeInfo) {
	dmColumnTypeList = append(dmColumnTypeList, columnType)
}
func init() {

}
