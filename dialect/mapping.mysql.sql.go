package dialect

func appendMysqlSql(mapping *SqlMapping) {

	// 库或所属者 相关 SQL
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

	// 表 相关 SQL
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
	mapping.TableCreateColumnHasComment = true
	mapping.TableCreateColumn = `
	{columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}]
`
	mapping.TableCreatePrimaryKey = `
PRIMARY KEY ({primaryKeysPack})
`
	mapping.TableComment = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} COMMENT {sqlValuePack(tableComment)}
`
	mapping.TableRename = `
ALTER TABLE [{ownerNamePack}.]{oldTableNamePack} RENAME AS {newTableNamePack}
`
	mapping.TableDelete = `
DROP TABLE IF EXISTS [{ownerNamePack}.]{tableNamePack}
`

	// 字段 相关 SQL
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
	mapping.ColumnUpdateHasRename = true
	mapping.ColumnUpdateHasComment = true
	mapping.ColumnUpdateHasAfter = true
	mapping.ColumnUpdate = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} CHANGE COLUMN {oldColumnNamePack} {columnNamePack} {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)} [COMMENT {sqlValuePack(columnComment)}] [AFTER {columnAfterColumnPack}]
`

	// 主键 相关 SQL
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

	// 索引 相关 SQL
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
}
