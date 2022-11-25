package dialect

import "strings"

var (
	dmColumnTypeList []*ColumnTypeInfo
)

func appendDmColumnType(columnType *ColumnTypeInfo) {
	dmColumnTypeList = append(dmColumnTypeList, columnType)
}
func init() {
	appendDmColumnType(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})

	appendDmColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "VARCHAR2", Format: "VARCHAR2($l)", IsString: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})

	appendDmColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			//if strings.Contains(strings.ToLower(column.ColumnExtra), "on update current_timestamp") {
			//	columnDefaultPack += " ON UPDATE CURRENT_TIMESTAMP"
			//}
			return
		},
	})

	// dm
	appendDmColumnType(&ColumnTypeInfo{Name: "VARBINARY", Format: "VARBINARY($l)", IsString: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "BINARY", Format: "BINARY($l)", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "BYTE", Format: "BYTE($l)", IsNumber: true})
	appendDmColumnType(&ColumnTypeInfo{Name: "CLASS234882065", Format: "CLASS234882065", IsString: true})
}
