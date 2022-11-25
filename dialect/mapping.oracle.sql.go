package dialect

func appendOracleSql(mapping *SqlMapping) {
	// 库或所属者 相关 SQL
	mapping.OwnersSelect = `
SELECT 
	USERNAME ownerName
FROM DBA_USERS 
ORDER BY USERNAME
`
	mapping.OwnerSelect = `
SELECT 
	USERNAME ownerName
FROM DBA_USERS 
WHERE USERNAME={sqlValuePack(ownerName)}
`
	mapping.OwnerCreate = `
CREATE USER {ownerName} IDENTIFIED BY {doubleQuotationMarksPack(ownerPassword)};
GRANT dba,resource,connect TO {ownerName};
`
	mapping.OwnerDelete = `
DROP USER {ownerName} CASCADE
`

	// 表 相关 SQL
	mapping.TablesSelect = `
SELECT 
	TABLE_NAME tableName,
	OWNER ownerName
FROM ALL_TABLES
WHERE OWNER={sqlValuePack(ownerName)}
`
	mapping.TableSelect = `
SELECT 
	TABLE_NAME tableName,
	OWNER ownerName
FROM ALL_TABLES
WHERE OWNER={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
`
	mapping.TableCreate = `
CREATE TABLE [{ownerNamePack}.]{tableNamePack}(
{ tableCreateColumnContent }
{ tableCreatePrimaryKeyContent }
)
`
	mapping.TableCreateColumn = `
	{columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`
	mapping.TableCreatePrimaryKey = `
PRIMARY KEY ({primaryKeysPack})
`
	mapping.TableComment = `
COMMENT ON TABLE [{ownerNamePack}.]{tableNamePack} IS {sqlValuePack(tableComment)}
`
	mapping.TableRename = `
ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME TO {tableNamePack}
`
	mapping.TableDelete = `
DROP TABLE [{ownerNamePack}.]{tableNamePack}
`

	// 字段 相关 SQL
	mapping.ColumnsSelect = `
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
`
	mapping.ColumnSelect = `
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
`
	mapping.ColumnAdd = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} ADD {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`
	mapping.ColumnDelete = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP COLUMN {columnNamePack}
`
	mapping.ColumnComment = `
COMMENT ON COLUMN [{ownerNamePack}.]{tableNamePack}.{columnNamePack} IS {sqlValuePack(columnComment)}
`
	mapping.ColumnRename = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} RENAME COLUMN {oldColumnNamePack} TO {columnNamePack}
`
	mapping.ColumnUpdateHasRename = false
	mapping.ColumnUpdateHasComment = false
	mapping.ColumnUpdateHasAfter = false
	mapping.ColumnUpdate = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} MODIFY {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`

	// 主键 相关 SQL
	mapping.PrimaryKeysSelect = `
SELECT 
    t1.COLUMN_NAME columnName,
    t2.TABLE_NAME tableName,
    t2.OWNER ownerName
FROM ALL_CONS_COLUMNS t1
LEFT JOIN ALL_CONSTRAINTS t2 ON (t2.CONSTRAINT_NAME = t1.CONSTRAINT_NAME)
WHERE t2.OWNER={sqlValuePack(ownerName)}
	AND t2.TABLE_NAME={sqlValuePack(tableName)}
	AND t2.CONSTRAINT_TYPE = 'P'
`
	mapping.PrimaryKeyAdd = `
ALTER TABLE [{ownerName}.]{tableName} ADD PRIMARY KEY ({columnNamesPack})
`
	mapping.PrimaryKeyDelete = `
ALTER TABLE [{ownerName}.]{tableName} DROP PRIMARY KEY
`

	// 索引 相关 SQL
	mapping.IndexesSelect = `
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
`

	mapping.IndexNameMaxLen = 30
	mapping.IndexAdd = `
CREATE {indexType} [{indexNamePack}] ON [{ownerNamePack}.]{tableNamePack} ({columnNamesPack})
`
	mapping.IndexDelete = `
DROP INDEX {indexNamePack}
`
}
