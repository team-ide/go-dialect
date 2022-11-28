package dialect

import "strings"

func NewMappingSqlite() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeSqlite,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	appendSqliteSql(mapping)

	for _, one := range sqliteColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range sqliteIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}

func AppendSqliteColumnType(mapping *SqlMapping) {

	// mysql 数据类型转换
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIT", Format: "BIT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYINT", Format: "TINYINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SMALLINT", Format: "SMALLINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMINT", Format: "MEDIUMINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT", Format: "INT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIGINT", Format: "BIGINT($l)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT", Format: "FLOAT($l, $d)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE", Format: "DOUBLE($l, $d)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DEC", Format: "DEC($l, $d)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DECIMAL", Format: "DECIMAL($l, $d)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "YEAR", Format: "YEAR", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIME", Format: "TIME", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATETIME", Format: "DATETIME", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			//if strings.Contains(strings.ToLower(column.ColumnExtra), "on update current_timestamp") {
			//	columnDefaultPack += " ON UPDATE CURRENT_TIMESTAMP"
			//}
			return
		}})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
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

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYTEXT", Format: "TINYTEXT", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMTEXT", Format: "MEDIUMTEXT", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGTEXT", Format: "LONGTEXT", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ENUM", Format: "ENUM", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYBLOB", Format: "TINYBLOB", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMBLOB", Format: "MEDIUMBLOB", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGBLOB", Format: "LONGBLOB", IsString: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SET", Format: "TEXT", IsString: true})

	// 浮点数
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REAL", Format: "REAL", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMERIC", IsNumber: true})

	// oracle
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR2", Format: "VARCHAR2($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "RAW", Format: "RAW($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NVARCHAR2", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NCLOB", Format: "NCLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XMLTYPE", Format: "XMLTYPE($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NROWID", Format: "NROWID($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BFILE", Format: "BFILE($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ANYDATA", Format: "ANYDATA($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ROWID", Format: "ROWID($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NCHAR", Format: "NCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_DIM_ARRAY", Format: "SDO_DIM_ARRAY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_TOPO_GEOMETRY_LAYER_ARRAY", Format: "SDO_TOPO_GEOMETRY_LAYER_ARRAY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_GEOMETRY", Format: "SDO_GEOMETRY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_NUMBER_ARRAY", Format: "SDO_NUMBER_ARRAY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG", Format: "LONG", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG RAW", Format: "LONG RAW", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "UNDEFINED", Format: "UNDEFINED", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MLSLABEL", Format: "MLSLABEL($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "WRI$_REPT_ABSTRACT_T", Format: "WRI$_REPT_ABSTRACT_T($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "RE$NV_LIST", Format: "RE$NV_LIST($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_AGENT", Format: "AQ$_AGENT($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTERVAL DAY", Format: "INTERVAL DAY($l)", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DBMS_DBFS_CONTENT_PROPERTIES_T", Format: "DBMS_DBFS_CONTENT_PROPERTIES_T($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SCHEDULER$_EVENT_INFO", Format: "SCHEDULER$_EVENT_INFO($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SCHEDULER$_REMOTE_DB_JOB_INFO", Format: "SCHEDULER$_REMOTE_DB_JOB_INFO($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SCHEDULER_FILEWATCHER_RESULT", Format: "SCHEDULER_FILEWATCHER_RESULT($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ALERT_TYPE", Format: "ALERT_TYPE($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "HSBLKNAMLST", Format: "HSBLKNAMLST($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_EVENT_MESSAGE", Format: "AQ$_EVENT_MESSAGE($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_NOTIFY_MSG", Format: "AQ$_NOTIFY_MSG($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "KUPC$_MESSAGE", Format: "KUPC$_MESSAGE($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS$RLBTYP", Format: "SYS$RLBTYP($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_SIG_PROP", Format: "AQ$_SIG_PROP($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_SUBSCRIBERS", Format: "AQ$_SUBSCRIBERS($l)", IsString: true, IsExtend: true})

	// ShenTong
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BPCHAR", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT1", Format: "NUMBER(1)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT2", Format: "NUMBER(2)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT2", Format: "NUMBER(2)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT4", Format: "NUMBER(4)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT4", Format: "NUMBER(4)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT8", Format: "NUMBER(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT8", Format: "NUMBER(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT2", Format: "NUMBER(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT4", Format: "NUMBER(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT8", Format: "NUMBER(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_FLOAT8", Format: "NUMBER(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BOOL", Format: "NUMBER(1)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "OIDVECTOR", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT2VECTOR", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BFILE", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_ACLITEM", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMPTZ", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_TEXT", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_OID", Format: "CLOB", IsString: true, IsExtend: true})

	// 金仓
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP WITHOUT TIME ZONE", Format: "TIMESTAMP", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER VARYING", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTEA", Format: "BLOB($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "OID", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NAME", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ARRAY", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP WITH TIME ZONE", Format: "TIMESTAMP", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGROLE", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGCLASS", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGPROC", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BOOLEAN", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE PRECISION", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS_LSN", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTERVAL", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "\"CHAR\"", Format: "CLOB", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS_NODE_TREE", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "JSON", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ANYARRAY", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INET", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ABSTIME", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XID", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TDEKEY", Format: "CLOB", IsString: true, IsExtend: true})

	// 达梦
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARBINARY", Format: "VARBINARY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BINARY", Format: "BINARY($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTE", Format: "BYTE($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLASS234882065", Format: "CLASS234882065", IsString: true, IsExtend: true})

}
