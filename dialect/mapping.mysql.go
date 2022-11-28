package dialect

func NewMappingMysql() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeMysql,

		OwnerNamePackChar:  "`",
		TableNamePackChar:  "`",
		ColumnNamePackChar: "`",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "\\",
	}

	appendMysqlSql(mapping)

	for _, one := range mysqlColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range mysqlIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}

func AppendMysqlColumnType(mapping *SqlMapping) {

	// sqlite
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REAL", Format: "DOUBLE($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMERIC", Format: "DECIMAL($l, $d)", IsNumber: true, IsExtend: true})

	// oracle
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR2", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMBER", Format: "DECIMAL($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLOB", Format: "LONGTEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "RAW", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NVARCHAR2", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NCLOB", Format: "LONGTEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XMLTYPE", Format: "VARCHAR($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NROWID", Format: "TEXT($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BFILE", Format: "LONGTEXT", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ANYDATA", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ROWID", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NCHAR", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_DIM_ARRAY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_TOPO_GEOMETRY_LAYER_ARRAY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_GEOMETRY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_NUMBER_ARRAY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG", Format: "LONGBLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG RAW", Format: "LONGBLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "UNDEFINED", Format: "VARCHAR", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MLSLABEL", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "WRI$_REPT_ABSTRACT_T", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "RE$NV_LIST", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_AGENT", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTERVAL DAY", Format: "DATETIME", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DBMS_DBFS_CONTENT_PROPERTIES_T", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SCHEDULER$_EVENT_INFO", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SCHEDULER$_REMOTE_DB_JOB_INFO", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SCHEDULER_FILEWATCHER_RESULT", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ALERT_TYPE", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "HSBLKNAMLST", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_EVENT_MESSAGE", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_NOTIFY_MSG", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "KUPC$_MESSAGE", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS$RLBTYP", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_SIG_PROP", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "AQ$_SUBSCRIBERS", Format: "VARCHAR($l)", IsString: true, IsExtend: true})

	// ShenTong
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BPCHAR", Format: "VARCHAR($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT1", Format: "DECIMAL(1)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT2", Format: "DECIMAL(2)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT2", Format: "DECIMAL(2)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT4", Format: "DECIMAL(4)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT4", Format: "DECIMAL(4)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT8", Format: "DECIMAL(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT8", Format: "DECIMAL(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT2", Format: "DECIMAL(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT4", Format: "DECIMAL(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT8", Format: "DECIMAL(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_FLOAT8", Format: "DECIMAL(8)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BOOL", Format: "DECIMAL(1)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "OIDVECTOR", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT2VECTOR", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BFILE", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_ACLITEM", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMPTZ", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_TEXT", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_OID", Format: "TEXT", IsString: true, IsExtend: true})
	// 金仓
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP WITHOUT TIME ZONE", Format: "TIMESTAMP", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER VARYING", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTEA", Format: "BLOB($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "OID", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NAME", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ARRAY", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP WITH TIME ZONE", Format: "TIMESTAMP", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGROLE", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGCLASS", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGPROC", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BOOLEAN", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE PRECISION", Format: "DECIMAL($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS_LSN", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTERVAL", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "\"CHAR\"", Format: "TEXT", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS_NODE_TREE", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "JSON", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ANYARRAY", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INET", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ABSTIME", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XID", Format: "TEXT", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TDEKEY", Format: "TEXT", IsString: true, IsExtend: true})

	// 达梦
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARBINARY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BINARY", Format: "DECIMAL($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTE", Format: "DECIMAL($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLASS234882065", Format: "LONGTEXT", IsString: true, IsExtend: true})

}
