package dialect

func NewMappingMysql() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeMysql,

		OwnerNamePackChar:  "`",
		TableNamePackChar:  "`",
		ColumnNamePackChar: "`",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "\\",

		// 库或所属者 相关 SQL
		OwnersSelect: `
SELECT
    SCHEMA_NAME ownerName,
    DEFAULT_CHARACTER_SET_NAME ownerCharacterSetName,
    DEFAULT_COLLATION_NAME ownerCollationName
FROM information_schema.schemata
ORDER BY SCHEMA_NAME
`,
		OwnerSelect: `
SELECT
    SCHEMA_NAME ownerName,
    DEFAULT_CHARACTER_SET_NAME ownerCharacterSetName,
    DEFAULT_COLLATION_NAME ownerCollationName
FROM information_schema.schemata
WHERE SCHEMA_NAME={sqlValuePack(ownerName)}
`,
		OwnerCreate: `
CREATE DATABASE [IF NOT EXISTS] {ownerNamePack}
[CHARACTER SET {ownerCharacterSetName}]
[COLLATE {ownerCollationName}]
`,
		OwnerDelete: `
DROP DATABASE IF EXISTS {ownerNamePack}
`,

		// 表 相关 SQL
		TablesSelect: `
SELECT
    TABLE_NAME tableName,
    TABLE_COMMENT tableComment,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
ORDER BY TABLE_NAME
`,
		TableSelect: `
SELECT
    TABLE_NAME tableName,
    TABLE_COMMENT tableComment,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
`,
		TableCreate: `
CREATE TABLE [{ownerNamePack}.]{tableNamePack}(
{ tableCreateColumnContent }
{ tableCreatePrimaryKeyContent }
)[CHARACTER SET {tableCharacterSetName}]
`,
		TableCreateColumnHasComment: true,
		TableCreateColumn: `
	{columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}]
`,
		TableCreatePrimaryKey: `
PRIMARY KEY ({primaryKeysPack})
`,
		TableComment: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} COMMENT {sqlValuePack(tableComment)}
`,
		TableRename: `
ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME AS {newTableNamePack}
`,
		TableDelete: `
DROP TABLE IF EXISTS [{ownerNamePack}.]{tableNamePack}
`,

		// 字段 相关 SQL
		ColumnsSelect: `
SELECT
    COLUMN_NAME columnName,
    COLUMN_COMMENT columnComment,
    COLUMN_DEFAULT columnDefault,
    EXTRA columnExtra,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    CHARACTER_SET_NAME columnCharacterSetName,
    IS_NULLABLE isNullable,
    DATA_TYPE columnDataType,
    COLUMN_TYPE columnType
FROM information_schema.columns
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
`,
		ColumnSelect: `
SELECT
    COLUMN_NAME columnName,
    COLUMN_COMMENT columnComment,
    COLUMN_DEFAULT columnDefault,
    EXTRA columnExtra,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    CHARACTER_SET_NAME columnCharacterSetName,
    IS_NULLABLE isNullable,
    DATA_TYPE columnDataType,
    COLUMN_TYPE columnType
FROM information_schema.columns
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
  AND COLUMN_NAME={sqlValuePack(columnName)}
`,
		ColumnAdd: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} ADD COLUMN {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}] [AFTER {columnAfterColumnPack}]
`,
		ColumnDelete: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP COLUMN {columnNamePack}
`,
		ColumnComment: `
`,
		ColumnRename: `
`,
		ColumnUpdateHasRename:  true,
		ColumnUpdateHasComment: true,
		ColumnUpdateHasAfter:   true,
		ColumnUpdate: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} CHANGE COLUMN {oldColumnNamePack} {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}] [AFTER {columnAfterColumnPack}]
`,

		// 主键 相关 SQL
		PrimaryKeysSelect: `
SELECT
    t2.COLUMN_NAME columnName,
    t1.TABLE_NAME tableName,
    t1.TABLE_SCHEMA ownerName
FROM information_schema.table_constraints t1
LEFT JOIN information_schema.key_column_usage t2 
ON (t2.CONSTRAINT_NAME=t1.CONSTRAINT_NAME AND t2.TABLE_SCHEMA=t1.TABLE_SCHEMA AND t2.TABLE_NAME=t1.TABLE_NAME)
WHERE t1.TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND t1.TABLE_NAME={sqlValuePack(tableName)}
  AND t1.CONSTRAINT_TYPE='PRIMARY KEY'
`,
		PrimaryKeyAdd: `
ALTER TABLE [{ownerName}.]{tableName} ADD PRIMARY KEY ({columnNamesPack})
`,
		PrimaryKeyDelete: `
ALTER TABLE [{ownerName}.]{tableName} DROP PRIMARY KEY
`,

		// 索引 相关 SQL
		IndexesSelect: `
SELECT
    t1.INDEX_NAME indexName,
    t1.COLUMN_NAME columnName,
    t1.INDEX_COMMENT indexComment,
    t1.NON_UNIQUE nonUnique,
    t1.TABLE_NAME tableName,
    t1.TABLE_SCHEMA ownerName,
    t2.CONSTRAINT_TYPE
FROM information_schema.statistics t1
LEFT JOIN information_schema.table_constraints t2 
ON (t2.CONSTRAINT_NAME=t1.INDEX_NAME AND t2.TABLE_SCHEMA=t1.TABLE_SCHEMA AND t2.TABLE_NAME=t1.TABLE_NAME)
WHERE t1.TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND t1.TABLE_NAME={sqlValuePack(tableName)}
  AND (t2.CONSTRAINT_TYPE !='PRIMARY KEY' OR t2.CONSTRAINT_TYPE = '' OR t2.CONSTRAINT_TYPE IS NULL)
`,
		IndexAdd: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} ADD {indexType} [{indexNamePack}] ({columnNamesPack}) [COMMENT {sqlValuePack(indexComment)}]
`,
		IndexDelete: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP INDEX {indexNamePack}
`,
	}

	AppendMysqlColumnType(mapping)
	AppendMysqlIndexType(mapping)

	return
}

func AppendMysqlColumnType(mapping *SqlMapping) {
	for _, columnType := range mysqlColumnTypeList {
		mapping.AddColumnTypeInfo(columnType)
	}

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

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ANYDATA", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ROWID", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NCHAR", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_DIM_ARRAY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_TOPO_GEOMETRY_LAYER_ARRAY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_GEOMETRY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_NUMBER_ARRAY", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG", Format: "DECIMAL($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG RAW", Format: "DECIMAL($l, $d)", IsNumber: true, IsExtend: true})
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

func AppendMysqlIndexType(mapping *SqlMapping) {

	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"TEXT"},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"TEXT"},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX"})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"TEXT"},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", OnlySupportDataTypes: []string{"CHAR", "VARCHAR", "TEXT"}})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", OnlySupportDataTypes: []string{"GEOMETRY", "POINT", "LINESTRING", "POLYGON"}})

}
