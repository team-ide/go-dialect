package dialect

// Mysql 数据库 SQL
func appendMysqlSql(mapping *SqlMapping) {

	mapping.OwnersSelect = `

SELECT
    SCHEMA_NAME ownerName,
    DEFAULT_CHARACTER_SET_NAME ownerCharacterSetName,
    DEFAULT_COLLATION_NAME ownerCollationName
FROM information_schema.schemata
ORDER BY SCHEMA_NAME
`

	mapping.OwnerSelect = `

SELECT
    SCHEMA_NAME ownerName,
    DEFAULT_CHARACTER_SET_NAME ownerCharacterSetName,
    DEFAULT_COLLATION_NAME ownerCollationName
FROM information_schema.schemata
WHERE SCHEMA_NAME={sqlValuePack(ownerName)}
`

	mapping.OwnerCreate = `

CREATE DATABASE [IF NOT EXISTS] {ownerNamePack}
[CHARACTER SET {ownerCharacterSetName}]
[COLLATE {ownerCollationName}]
`

	mapping.OwnerDelete = `

DROP DATABASE IF EXISTS {ownerNamePack}
`

	mapping.TablesSelect = `

SELECT
    TABLE_NAME tableName,
    TABLE_COMMENT tableComment,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
ORDER BY TABLE_NAME
`

	mapping.TableSelect = `

SELECT
    TABLE_NAME tableName,
    TABLE_COMMENT tableComment,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
`

	mapping.TableCreate = `

CREATE TABLE [{ownerNamePack}.]{tableNamePack}(
{ tableCreateColumnContent }
{ tableCreatePrimaryKeyContent }
)[CHARACTER SET {tableCharacterSetName}]
`

	mapping.TableCreateColumn = `

	{columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}]
`
	mapping.TableCreateColumnHasComment = true

	mapping.TableCreatePrimaryKey = `

PRIMARY KEY ({primaryKeysPack})
`

	mapping.TableDelete = `

DROP TABLE IF EXISTS [{ownerNamePack}.]{tableNamePack}
`

	mapping.TableComment = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} COMMENT {sqlValuePack(tableComment)}`

	mapping.TableRename = `

ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME AS {tableNamePack}
`

	mapping.ColumnsSelect = `

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
    COLUMN_TYPE columnType,
    DATETIME_PRECISION DATETIME_PRECISION,
    NUMERIC_PRECISION NUMERIC_PRECISION,
    NUMERIC_SCALE NUMERIC_SCALE,
    CHARACTER_MAXIMUM_LENGTH CHARACTER_MAXIMUM_LENGTH
FROM information_schema.columns
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
`

	mapping.ColumnSelect = `

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
`

	mapping.ColumnAdd = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} ADD COLUMN {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}] [AFTER {columnAfterColumnPack}]
`

	mapping.ColumnDelete = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP COLUMN {columnNamePack}
`

	mapping.ColumnComment = `
`

	mapping.ColumnRename = `
`

	mapping.ColumnUpdate = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} CHANGE COLUMN {oldColumnNamePack} {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}] [AFTER {columnAfterColumnPack}]
`
	mapping.ColumnUpdateHasComment = true
	mapping.ColumnUpdateHasRename = true
	mapping.ColumnUpdateHasAfter = true

	mapping.ColumnAfter = `
`

	mapping.PrimaryKeysSelect = `

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
`

	mapping.PrimaryKeyAdd = `

ALTER TABLE [{ownerName}.]{tableName} ADD PRIMARY KEY ({columnNamesPack})
`

	mapping.PrimaryKeyDelete = `

ALTER TABLE [{ownerName}.]{tableName} DROP PRIMARY KEY
`

	mapping.IndexesSelect = `

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
`

	mapping.IndexAdd = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} ADD {indexType} [{indexNamePack}] ({columnNamesPack}) [COMMENT {sqlValuePack(indexComment)}]
`

	mapping.IndexDelete = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP INDEX {indexNamePack}
`

	mapping.IndexNamePack = `
`
}

// Oracle 数据库 SQL
func appendOracleSql(mapping *SqlMapping) {

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

	mapping.TablesSelect = `

SELECT 
	TABLE_NAME tableName,
	OWNER ownerName
FROM ALL_TABLES
WHERE OWNER={sqlValuePack(ownerName)}
ORDER BY TABLE_NAME `

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

	mapping.TableDelete = `

DROP TABLE [{ownerNamePack}.]{tableNamePack}
`

	mapping.TableComment = `

COMMENT ON TABLE [{ownerNamePack}.]{tableNamePack} IS {sqlValuePack(tableComment)}
`

	mapping.TableRename = `

ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME TO {tableNamePack}
`

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

	mapping.ColumnUpdate = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} MODIFY {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`

	mapping.ColumnAfter = `
`

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

	mapping.IndexAdd = `

CREATE {indexType} [{indexNamePack}] ON [{ownerNamePack}.]{tableNamePack} ({columnNamesPack})
`

	mapping.IndexDelete = `

DROP INDEX {indexNamePack}
`

	mapping.IndexNamePack = `
`
}

// 达梦 数据库 SQL
func appendDmSql(mapping *SqlMapping) {

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

CREATE USER {doubleQuotationMarksPack(ownerName)} IDENTIFIED BY {doubleQuotationMarksPack(ownerPassword)};
GRANT DBA TO {doubleQuotationMarksPack(ownerName)};
`

	mapping.OwnerDelete = `

DROP USER {ownerName} CASCADE
`

	mapping.TablesSelect = `

SELECT 
	TABLE_NAME tableName,
	OWNER ownerName
FROM ALL_TABLES
WHERE OWNER={sqlValuePack(ownerName)}
ORDER BY TABLE_NAME 
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

	mapping.TableDelete = `

DROP TABLE [{ownerNamePack}.]{tableNamePack}
`

	mapping.TableComment = `

COMMENT ON TABLE [{ownerNamePack}.]{tableNamePack} IS {sqlValuePack(tableComment)}
`

	mapping.TableRename = `

ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME TO {tableNamePack}
`

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

	mapping.ColumnUpdate = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} MODIFY {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`

	mapping.ColumnAfter = `
`

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

	mapping.IndexAdd = `

CREATE {indexType} [{indexNamePack}] ON [{ownerNamePack}.]{tableNamePack} ({columnNamesPack})
`

	mapping.IndexDelete = `

DROP INDEX {indexNamePack}
`

	mapping.IndexNamePack = `
`
}

// 金仓 数据库 SQL
func appendKingBaseSql(mapping *SqlMapping) {

	mapping.OwnersSelect = `

SELECT
    SCHEMA_NAME ownerName
FROM information_schema.schemata
ORDER BY SCHEMA_NAME
`

	mapping.OwnerSelect = `

SELECT
    SCHEMA_NAME ownerName
FROM information_schema.schemata
WHERE SCHEMA_NAME={sqlValuePack(ownerName)}
`

	mapping.OwnerCreate = `

CREATE USER {ownerName} WITH PASSWORD {sqlValuePack(ownerPassword)};
CREATE SCHEMA {ownerName};
GRANT USAGE ON SCHEMA {ownerName} TO {ownerName};
GRANT ALL ON SCHEMA {ownerName} TO {ownerName};
GRANT ALL ON ALL TABLES IN SCHEMA {ownerName} TO {ownerName};
`

	mapping.OwnerDelete = `

DROP SCHEMA IF EXISTS {ownerName} CASCADE;
DROP USER IF EXISTS {ownerName};
`

	mapping.TablesSelect = `

SELECT
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
ORDER BY TABLE_NAME
`

	mapping.TableSelect = `

SELECT
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
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

	mapping.TableDelete = `

DROP TABLE [{ownerNamePack}.]{tableNamePack}
`

	mapping.TableComment = `

COMMENT ON TABLE [{ownerNamePack}.]{tableNamePack} IS {sqlValuePack(tableComment)}
`

	mapping.TableRename = `

ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME TO {tableNamePack}
`

	mapping.ColumnsSelect = `


SELECT
    COLUMN_NAME columnName,
    COLUMN_DEFAULT columnDefault,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    IS_NULLABLE isNullable,
    DATA_TYPE columnDataType,
    DATETIME_PRECISION DATETIME_PRECISION,
    NUMERIC_PRECISION NUMERIC_PRECISION,
    NUMERIC_SCALE NUMERIC_SCALE,
    CHARACTER_MAXIMUM_LENGTH CHARACTER_MAXIMUM_LENGTH
FROM information_schema.columns
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}

`

	mapping.ColumnSelect = `
SELECT
    COLUMN_NAME columnName,
    COLUMN_DEFAULT columnDefault,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    IS_NULLABLE isNullable,
    DATA_TYPE columnDataType,
    NUMERIC_PRECISION NUMERIC_PRECISION,
    NUMERIC_SCALE NUMERIC_SCALE,
    CHARACTER_MAXIMUM_LENGTH CHARACTER_MAXIMUM_LENGTH
FROM information_schema.columns
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
  AND COLUMN_NAME={sqlValuePack(columnName)}
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

	mapping.ColumnUpdate = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} ALTER COLUMN {columnNamePack} TYPE {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}

`

	mapping.ColumnAfter = `
`

	mapping.PrimaryKeysSelect = `

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
`

	mapping.PrimaryKeyAdd = `

ALTER TABLE [{ownerName}.]{tableName} ADD PRIMARY KEY ({columnNamesPack})
`

	mapping.PrimaryKeyDelete = `

ALTER TABLE [{ownerName}.]{tableName} DROP PRIMARY KEY
`

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
WHERE t1.INDEX_NAME IN(
    SELECT INDEXNAME 
    FROM SYS_CATALOG.sys_indexes 
    WHERE SCHEMANAME={sqlValuePack(ownerName)}
        AND TABLENAME={sqlValuePack(tableName)}
) 
	AND t1.INDEX_NAME NOT IN(
    SELECT
    t1.CONSTRAINT_NAME
FROM information_schema.table_constraints t1
WHERE t1.TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND t1.TABLE_NAME={sqlValuePack(tableName)}
  AND t1.CONSTRAINT_TYPE='PRIMARY KEY'
)
`

	mapping.IndexAdd = `

CREATE {indexType} [{indexNamePack}] ON [{ownerNamePack}.]{tableNamePack} ({columnNamesPack})
`

	mapping.IndexDelete = `

DROP INDEX {indexNamePack}
`

	mapping.IndexNamePack = `
`
}

// 神通 数据库 SQL
func appendShenTongSql(mapping *SqlMapping) {

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

CREATE USER {ownerName} WITH PASSWORD {sqlValuePack(ownerPassword)};
`

	mapping.OwnerDelete = `

DROP USER {ownerName} CASCADE
`

	mapping.TablesSelect = `

SELECT 
	TABLE_NAME tableName,
	OWNER ownerName
FROM ALL_TABLES
WHERE OWNER={sqlValuePack(ownerName)}
ORDER BY TABLE_NAME 
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

	mapping.TableDelete = `

DROP TABLE [{ownerNamePack}.]{tableNamePack}
`

	mapping.TableComment = `

COMMENT ON TABLE [{ownerNamePack}.]{tableNamePack} IS {sqlValuePack(tableComment)}
`

	mapping.TableRename = `

ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME TO {tableNamePack}
`

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

	mapping.ColumnUpdate = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} MODIFY {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`

	mapping.ColumnAfter = `
`

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

	mapping.IndexAdd = `

CREATE {indexType} [{indexNamePack}] ON [{ownerNamePack}.]{tableNamePack} ({columnNamesPack})
`

	mapping.IndexDelete = `

DROP INDEX {indexNamePack}
`

	mapping.IndexNamePack = `
`
}

// Sqlite 数据库 SQL
func appendSqliteSql(mapping *SqlMapping) {

	mapping.OwnersSelect = `

SELECT 
	name ownerName
FROM pragma_database_list AS t_i 
ORDER BY name
`

	mapping.OwnerSelect = `

SELECT 
	name ownerName
FROM pragma_database_list AS t_i 
WHERE name={sqlValuePack(ownerName)}
`

	mapping.OwnerCreate = `
`

	mapping.OwnerDelete = `
`

	mapping.TablesSelect = `

SELECT 
	name tableName,
    sql 
FROM sqlite_master 
WHERE type ='table'
ORDER BY name
`

	mapping.TableSelect = `

SELECT 
	name tableName,
    sql 
FROM sqlite_master 
WHERE type ='table'
  AND name={sqlValuePack(tableName)}
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

	mapping.TableDelete = `

DROP TABLE IF EXISTS [{ownerName}.]{tableName}
`

	mapping.TableComment = `
`

	mapping.TableRename = `

ALTER TABLE [{ownerName}.]{oldTableName} RENAME AS {newTableName}
`

	mapping.ColumnsSelect = `

SELECT 
	name columnName,
	dflt_value columnDefault,
	"notnull" isNotNull,
	type columnType
FROM pragma_table_info({tableNamePack}) AS t_i 
`

	mapping.ColumnSelect = `

SELECT 
	name columnName,
	dflt_value columnDefault,
	"notnull" isNotNull,
	type columnType
FROM pragma_table_info({tableNamePack}) AS t_i 
WHERE name={sqlValuePack(columnName)}
`

	mapping.ColumnAdd = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} ADD COLUMN {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`

	mapping.ColumnDelete = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP COLUMN {columnNamePack}
`

	mapping.ColumnComment = `
`

	mapping.ColumnRename = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} RENAME COLUMN {oldColumnNamePack} TO {columnNamePack}
`

	mapping.ColumnUpdate = `
`

	mapping.ColumnAfter = `
`

	mapping.PrimaryKeysSelect = `

SELECT 
	a.name indexName,
	b.name columnName 
FROM pragma_index_list({tableNamePack}) AS a,pragma_index_info(a.name) b 
WHERE a.origin = "pk"
`

	mapping.PrimaryKeyAdd = `

ALTER TABLE [{ownerName}.]{tableName} ADD PRIMARY KEY ({columnNamesPack})
`

	mapping.PrimaryKeyDelete = `

ALTER TABLE [{ownerName}.]{tableName} DROP PRIMARY KEY
`

	mapping.IndexesSelect = `

SELECT 
	a.name indexName,
	a."unique" isUnique,
	b.name columnName 
FROM pragma_index_list({tableNamePack}) AS a,pragma_index_info(a.name) b 
WHERE a.origin != "pk"
`

	mapping.IndexAdd = `

CREATE {indexType} [{indexNamePack}] ON {tableNamePack}({columnNamesPack})
`

	mapping.IndexDelete = `

ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP INDEX {indexNamePack}
`

	mapping.IndexNamePack = `
`
}

// GBase 数据库 SQL
func appendGBaseSql(mapping *SqlMapping) {
}

// Postgresql 数据库 SQL
func appendPostgresqlSql(mapping *SqlMapping) {
}

// DB2 数据库 SQL
func appendDb2Sql(mapping *SqlMapping) {
}

