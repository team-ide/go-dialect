package dialect

var (
	shenTongColumnTypeList []*ColumnTypeInfo
)

func appendShenTongColumnType(columnType *ColumnTypeInfo) {
	shenTongColumnTypeList = append(shenTongColumnTypeList, columnType)
}
func init() {

}
