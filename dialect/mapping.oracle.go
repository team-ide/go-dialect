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
		pageSql = selectSql
		locOrderBy := findOrderByIndex(selectSql)
		if len(locOrderBy) < 1 { //如果没有 order by,增加默认的排序
			pageSql += " ORDER BY NULL "
		}
		pageSql += " OFFSET "
		pageSql += strconv.Itoa(pageSize * (pageNo - 1))
		pageSql += " ROWS FETCH NEXT "
		pageSql += strconv.Itoa(pageSize)
		pageSql += " ROWS ONLY "
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
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLOB", Format: "CLOB", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB", IsString: true})

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
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BPCHAR", Format: "VARCHAR($l)", IsString: true, IsExtend: true})

	// 金仓
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP WITHOUT TIME ZONE", Format: "TIMESTAMP", IsDateTime: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHARACTER VARYING", Format: "VARCHAR2($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTEA", Format: "BLOB", IsString: true, IsExtend: true})

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
