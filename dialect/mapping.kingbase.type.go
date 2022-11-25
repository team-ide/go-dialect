package dialect

import "strings"

var (
	kingBaseColumnTypeList []*ColumnTypeInfo
)

func appendKingBaseColumnType(columnType *ColumnTypeInfo) {
	kingBaseColumnTypeList = append(kingBaseColumnTypeList, columnType)
}
func init() {
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})

	appendKingBaseColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "VARCHAR2", Format: "VARCHAR2($l)", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})

	appendKingBaseColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
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

	// 金仓
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CHARACTER", Format: "CHARACTER($l)", IsString: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "BYTEA", Format: "BYTEA", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "TIMESTAMP WITHOUT TIME ZONE", Format: "TIMESTAMP", IsDateTime: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "CHARACTER VARYING", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})

	appendKingBaseColumnType(&ColumnTypeInfo{Name: "OID", Format: "OID", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "NAME", Format: "NAME", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "ARRAY", Format: "ARRAY", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "TIMESTAMP WITH TIME ZONE", Format: "TIMESTAMP", IsDateTime: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "REGROLE", Format: "REGROLE", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "REGCLASS", Format: "REGCLASS", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "REGPROC", Format: "REGPROC", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "BOOLEAN", Format: "BOOLEAN", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "DOUBLE PRECISION", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "SYS_LSN", Format: "SYS_LSN", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "INTERVAL", Format: "INTERVAL", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "\"CHAR\"", Format: "CLOB", IsString: true, IsExtend: true})

	appendKingBaseColumnType(&ColumnTypeInfo{Name: "SYS_NODE_TREE", Format: "SYS_NODE_TREE", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "JSON", Format: "JSON", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "ANYARRAY", Format: "ANYARRAY", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "INET", Format: "INET", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "ABSTIME", Format: "ABSTIME", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "XID", Format: "XID", IsString: true, IsExtend: true})
	appendKingBaseColumnType(&ColumnTypeInfo{Name: "TDEKEY", Format: "TDEKEY", IsString: true, IsExtend: true})
}
