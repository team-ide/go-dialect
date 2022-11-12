package dialect

func NewMappingMysql() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeMysql,
		// 库或所属者 相关 SQL
		OwnersSelect: `
SELECT
    SCHEMA_NAME name,
    DEFAULT_CHARACTER_SET_NAME characterSetName,
    DEFAULT_COLLATION_NAME collationName
FROM information_schema.schemata
ORDER BY SCHEMA_NAME
`,
		OwnerSelect: `
SELECT
    SCHEMA_NAME name,
    DEFAULT_CHARACTER_SET_NAME characterSetName,
    DEFAULT_COLLATION_NAME collationName
FROM information_schema.schemata
WHERE SCHEMA_NAME='{ownerName}'
`,
		OwnerCreate: `
CREATE DATABASE [IF NOT EXISTS] {ownerName}
[[DEFAULT] CHARACTER SET {characterSetName}]
[[DEFAULT] COLLATE {collationName}]
`,
		OwnerDelete: `
DROP DATABASE IF EXISTS {ownerName}
`,

		// 表 相关 SQL
		TablesSelect: `
SELECT
    TABLE_NAME name,
    TABLE_COMMENT comment,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA='{ownerName}'
ORDER BY TABLE_NAME
`,
		TableSelect: `
SELECT
  TABLE_NAME name,
  TABLE_COMMENT comment,
  TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA='{ownerName}'
  AND TABLE_NAME='{tableName}'
`,
		TableCreate: `
CREATE TABLE [{ownerName}.]{tableName}(
    { for column in columnList }
	{column.columnName} {column.columnType} [CHARACTER SET {column.characterSetName}] [DEFAULT {column.default}] [NOT NULL][,]
	{ }
)[CHARACTER SET {characterSetName}] [COMMENT {tableComment}]
`,
		TableComment: `
ALTER TABLE [{ownerName}.]{tableName} COMMENT '{tableComment}'
`,
		TableRename: `
ALTER TABLE [{ownerName}.]{oldTableName} RENAME AS {newTableName}
`,
		TableDelete: `
DROP TABLE IF EXISTS [{ownerName}.]{tableName}
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
    CHARACTER_SET_NAME characterSetName,
    IS_NULLABLE isNullable,
    COLUMN_TYPE columnType,
    DATA_TYPE dataType
FROM information_schema.columns
WHERE TABLE_SCHEMA='{ownerName}'
  AND TABLE_NAME='{tableName}'
`,
		ColumnSelect: `
SELECT
    COLUMN_NAME columnName,
    COLUMN_COMMENT columnComment,
    COLUMN_DEFAULT columnDefault,
    EXTRA columnExtra,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    CHARACTER_SET_NAME characterSetName,
    IS_NULLABLE isNullable,
    COLUMN_TYPE columnType,
    DATA_TYPE dataType
FROM information_schema.columns
WHERE TABLE_SCHEMA='{ownerName}'
  AND TABLE_NAME='{tableName}'
  AND TABLE_NAME='{tableName}'
  AND COLUMN_NAME='{columnName}'
`,
		ColumnAdd: `
ALTER TABLE [{ownerName}.]{tableName} ADD COLUMN {columnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] [NOT NULL] [COMMENT {columnComment}]
`,
		ColumnComment: `
ALTER TABLE [{ownerName}.]{tableName} CHANGE COLUMN {columnName} {columnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] [NOT NULL] [COMMENT {columnComment}]
`,
		ColumnDelete: `
ALTER TABLE [{ownerName}.]{tableName} DROP COLUMN {columnName}
`,
		ColumnRename: `
ALTER TABLE [{ownerName}.]{tableName} CHANGE COLUMN {oldColumnName} {newColumnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] [NOT NULL] [COMMENT {columnComment}]
`,
		ColumnUpdate: `
ALTER TABLE [{ownerName}.]{tableName} CHANGE COLUMN {columnName} {columnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] [NOT NULL] [COMMENT {columnComment}] [AFTER {columnAfter}]
`,

		// 主键 相关 SQL
		PrimaryKeysSelect: `
SELECT
    key_column_usage.COLUMN_NAME columnName,
    table_constraints.TABLE_NAME tableName,
    table_constraints.TABLE_SCHEMA ownerName
FROM information_schema.table_constraints
JOIN information_schema.key_column_usage USING (CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME)
WHERE table_constraints.TABLE_SCHEMA='{ownerName}'
  AND table_constraints.TABLE_NAME='{tableName}'
  AND table_constraints.CONSTRAINT_TYPE='PRIMARY KEY'
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
    INDEX_NAME indexName,
    COLUMN_NAME columnName,
    INDEX_COMMENT indexComment,
    NON_UNIQUE nonUnique,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName
FROM information_schema.statistics
WHERE TABLE_SCHEMA='{ownerName}'
  AND TABLE_NAME='{tableName}'
  AND TABLE_NAME='{tableName}'
  AND INDEX_NAME NOT IN(
    SELECT table_constraints.CONSTRAINT_NAME
    FROM information_schema.table_constraints
    JOIN information_schema.key_column_usage USING (CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME)
    WHERE table_constraints.TABLE_SCHEMA='{ownerName}'
      AND table_constraints.TABLE_NAME='{tableName}'
      AND table_constraints.CONSTRAINT_TYPE='PRIMARY KEY'
)
`,
		IndexAdd: `
ALTER TABLE [{ownerName}.]{tableName} ADD [PRIMARY KEY | UNIQUE | FULLTEXT | INDEX] {indexName} ({columnNames}) [COMMENT {columnComment}]
`,
		IndexDelete: `
ALTER TABLE [{ownerName}.]{tableName} DROP INDEX
`,
	}

	return
}
