package dialect

import "strings"

var (
	shenTongColumnTypeList []*ColumnTypeInfo
)

func appendShenTongColumnType(columnType *ColumnTypeInfo) {
	shenTongColumnTypeList = append(shenTongColumnTypeList, columnType)
}
func init() {

	appendShenTongColumnType(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "VARCHAR2", Format: "VARCHAR2($l)", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
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

	// 神通
	appendShenTongColumnType(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB($l)", IsString: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})

	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT1", Format: "INT1(1)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT2", Format: "INT2(2)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_INT2", Format: "INT2(2)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT4", Format: "INT4(4)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_INT4", Format: "INT4(4)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT8", Format: "INT8(8)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_INT8", Format: "INT8(8)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "FLOAT2", Format: "FLOAT2(8)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "FLOAT4", Format: "FLOAT4(8)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "FLOAT8", Format: "FLOAT8(8)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_FLOAT8", Format: "FLOAT8(8)", IsNumber: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "OIDVECTOR", Format: "OIDVECTOR", IsString: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "BOOL", Format: "BOOL(1)", IsString: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "INT2VECTOR", Format: "INT2VECTOR", IsString: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "BFILE", Format: "BFILE", IsString: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_ACLITEM", Format: "_ACLITEM", IsString: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "TIMESTAMPTZ", Format: "TIMESTAMPTZ($l)", IsString: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_TEXT", Format: "_TEXT", IsString: true, IsExtend: true})
	appendShenTongColumnType(&ColumnTypeInfo{Name: "_OID", Format: "_OID", IsString: true, IsExtend: true})
}
