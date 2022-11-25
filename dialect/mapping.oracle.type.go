package dialect

var (
	oracleColumnTypeList []*ColumnTypeInfo
)

func appendOracleColumnType(columnType *ColumnTypeInfo) {
	oracleColumnTypeList = append(oracleColumnTypeList, columnType)
}
func init() {

}
