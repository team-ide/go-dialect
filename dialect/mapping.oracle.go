package dialect

import (
	"strconv"
	"strings"
)

func NewMappingOracle() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeOracle,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",

		// 库或所属者 相关 SQL
		OwnersSelect: `
SELECT 
	USERNAME ownerName
FROM DBA_USERS 
ORDER BY USERNAME
`,
		OwnerSelect: `
SELECT 
	USERNAME ownerName
FROM DBA_USERS 
WHERE USERNAME={sqlValuePack(ownerName)}
`,
		OwnerCreate: `
CREATE USER {ownerName} IDENTIFIED BY {doubleQuotationMarksPack(ownerPassword)};
GRANT dba,resource,connect TO {ownerName};
`,
		OwnerDelete: `
DROP USER {ownerName} CASCADE
`,

		// 表 相关 SQL
		TablesSelect: `
SELECT 
	TABLE_NAME tableName,
	OWNER ownerName
FROM ALL_TABLES
WHERE OWNER={sqlValuePack(ownerName)}
`,
		TableSelect: `
SELECT 
	TABLE_NAME tableName,
	OWNER ownerName
FROM ALL_TABLES
WHERE OWNER={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
`,
		TableCreate: `
CREATE TABLE [{ownerNamePack}.]{tableNamePack}(
{ tableCreateColumnContent }
{ tableCreatePrimaryKeyContent }
)
`,
		TableCreateColumn: `
	{columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`,
		TableCreatePrimaryKey: `
PRIMARY KEY ({primaryKeysPack})
`,
		TableComment: `
COMMENT ON TABLE [{ownerNamePack}.]{tableNamePack} IS {sqlValuePack(tableComment)}
`,
		TableRename: `
ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME TO {tableNamePack}
`,
		TableDelete: `
DROP TABLE [{ownerNamePack}.]{tableNamePack}
`,

		// 字段 相关 SQL
		ColumnsSelect: `
SELECT 
	t.COLUMN_NAME columnName,
	t.DATA_DEFAULT columnDefault,
	t.CHARACTER_SET_NAME columnCharacterSetName,
	t.NULLABLE isNullable,
	t.DATA_TYPE columnDataType,
	t.DATA_LENGTH,
	t.DATA_PRECISION,
	t.DATA_SCALE,
	tc.COMMENTS columnComment,
	t.TABLE_NAME tableName,
	t.OWNER ownerName
FROM ALL_TAB_COLUMNS t
LEFT JOIN ALL_COL_COMMENTS tc ON(tc.OWNER=t.OWNER AND tc.TABLE_NAME=t.TABLE_NAME AND tc.COLUMN_NAME=t.COLUMN_NAME)
WHERE t.OWNER={sqlValuePack(ownerName)}
    AND t.TABLE_NAME={sqlValuePack(tableName)}
`,
		ColumnSelect: `
SELECT 
	t.COLUMN_NAME columnName,
	t.DATA_DEFAULT columnDefault,
	t.CHARACTER_SET_NAME columnCharacterSetName,
	t.NULLABLE isNullable,
	t.DATA_TYPE columnDataType,
	t.DATA_LENGTH,
	t.DATA_PRECISION,
	t.DATA_SCALE,
	tc.COMMENTS columnComment,
	t.TABLE_NAME tableName,
	t.OWNER ownerName
FROM ALL_TAB_COLUMNS t
LEFT JOIN ALL_COL_COMMENTS tc ON(tc.OWNER=t.OWNER AND tc.TABLE_NAME=t.TABLE_NAME AND tc.COLUMN_NAME=t.COLUMN_NAME)
WHERE t.OWNER={sqlValuePack(ownerName)}
    AND t.TABLE_NAME={sqlValuePack(tableName)}
    AND t.COLUMN_NAME={sqlValuePack(columnName)}
`,
		ColumnAdd: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} ADD {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`,
		ColumnComment: `
COMMENT ON COLUMN [{ownerNamePack}.]{tableNamePack}.{columnNamePack} IS {sqlValuePack(columnComment)}
`,
		ColumnDelete: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP COLUMN {columnNamePack}
`,
		ColumnRename: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} RENAME COLUMN {oldColumnNamePack} {columnNamePack}
`,
		ColumnUpdate: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} CHANGE COLUMN {columnNamePack} {columnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] {columnNotNull(columnNotNull)} [COMMENT {columnComment}] [AFTER {columnAfter}]
`,

		// 主键 相关 SQL
		PrimaryKeysSelect: `
SELECT 
    t1.COLUMN_NAME columnName,
    t2.TABLE_NAME tableName,
    t2.OWNER ownerName
FROM ALL_CONS_COLUMNS t1
LEFT JOIN ALL_CONSTRAINTS t2 ON (t2.CONSTRAINT_NAME = t1.CONSTRAINT_NAME)
WHERE t2.OWNER={sqlValuePack(ownerName)}
	AND t2.TABLE_NAME={sqlValuePack(tableName)}
	AND t2.CONSTRAINT_TYPE = 'P'
`,
		PrimaryKeyAdd: `
ALTER TABLE [{ownerName}.]{tableName} ADD PRIMARY KEY ({columnNames})
`,
		PrimaryKeyDelete: `
ALTER TABLE [{ownerName}.]{tableName} DROP PRIMARY KEY
`,

		// 索引 相关 SQL
		IndexesSelect: `
SELECT 
    t1.INDEX_NAME indexName,
    t1.COLUMN_NAME columnName,
    t1.TABLE_OWNER ownerName,
    t1.TABLE_NAME tableName,
    t2.UNIQUENESS 
FROM ALL_IND_COLUMNS t1
LEFT JOIN ALL_INDEXES t2 ON (t2.INDEX_NAME = t1.INDEX_NAME)
LEFT JOIN ALL_CONSTRAINTS t3 ON (t3.CONSTRAINT_NAME = t1.INDEX_NAME)
WHERE t1.TABLE_OWNER={sqlValuePack(ownerName)}
	AND t1.TABLE_NAME={sqlValuePack(tableName)}
	AND (t3.CONSTRAINT_TYPE !='P' OR t3.CONSTRAINT_TYPE = '' OR t3.CONSTRAINT_TYPE IS NULL)
`,

		IndexNameMaxLen: 30,
		IndexAdd: `
CREATE {indexType} [{indexNamePack}] ON [{ownerNamePack}.]{tableNamePack} ({columnNamesPack})
`,
		IndexDelete: `
DROP INDEX {indexNamePack}
`,
	}

	AppendOracleColumnType(mapping)
	AppendOracleIndexType(mapping)

	mapping.PackPageSql = func(selectSql string, pageSize int, pageNo int) (pageSql string) {
		pageSql = `SELECT * FROM(SELECT ROWNUM rn,t.* FROM(` + selectSql + `) t WHERE ROWNUM <=` + strconv.Itoa(pageSize*pageNo) + ")"
		pageSql += `WHERE rn>=` + strconv.Itoa(pageSize*(pageNo-1)+1)
		return
	}
	mapping.ReplaceSqlVariable = func(sqlInfo string, args []interface{}) (variableSql string) {
		strList := strings.Split(sqlInfo, "?")
		if len(strList) < 1 {
			variableSql = sqlInfo
			return
		}
		variableSql = strList[0]
		for i := 1; i < len(strList); i++ {
			variableSql += ":" + strconv.Itoa(i)
			variableSql += strList[i]
		}
		return
	}
	return
}

func AppendOracleColumnType(mapping *SqlMapping) {

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR2", Format: "VARCHAR2($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})

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
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "RAW", Format: "RAW($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NVARCHAR2", Format: "NVARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NCLOB", Format: "NCLOB($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XMLTYPE", Format: "XMLTYPE($l)", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ANYDATA", Format: "ANYDATA($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ROWID", Format: "ROWID($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NCHAR", Format: "NCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_DIM_ARRAY", Format: "SDO_DIM_ARRAY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_TOPO_GEOMETRY_LAYER_ARRAY", Format: "SDO_TOPO_GEOMETRY_LAYER_ARRAY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_GEOMETRY", Format: "SDO_GEOMETRY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SDO_NUMBER_ARRAY", Format: "SDO_NUMBER_ARRAY($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG", Format: "LONG", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONG RAW", Format: "LONG RAW", IsNumber: true, IsExtend: true})
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

	// mysql
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIT", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYINT", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SMALLINT", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMINT", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTEGER", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIGINT", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DEC", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DECIMAL", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "YEAR", Format: "DATE", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIME", Format: "DATE", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATETIME", Format: "DATE", IsDateTime: true, IsExtend: true,
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

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYTEXT", Format: "VARCHAR2(1000)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TEXT", Format: "VARCHAR2(4000)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMTEXT", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGTEXT", Format: "CLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ENUM", Format: "VARCHAR2(50)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYBLOB", Format: "BLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMBLOB", Format: "BLOB", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGBLOB", Format: "BLOB", IsString: true, IsExtend: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SET", Format: "VARCHAR2(50)", IsString: true, IsExtend: true})

	// sqlite
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REAL", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMBER($l, $d)", IsNumber: true, IsExtend: true})

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
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER VARYING", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTEA", Format: "BLOB", IsString: true, IsExtend: true})

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
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARBINARY", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BINARY", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTE", Format: "NUMBER($l)", IsNumber: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLASS234882065", Format: "CLOB", IsString: true, IsExtend: true})

}

func AppendOracleIndexType(mapping *SqlMapping) {

	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "INDEX", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "NORMAL", Format: "INDEX",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
		NotSupportDataTypes: []string{"CLOB", "BLOB"},
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			indexTypeFormat = "UNIQUE INDEX"
			return
		},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "FULLTEXT", Format: "FULLTEXT", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "SPATIAL", Format: "SPATIAL", IsExtend: true,
		IndexTypeFormat: func(index *IndexModel) (indexTypeFormat string, err error) {
			return
		},
	})
}
