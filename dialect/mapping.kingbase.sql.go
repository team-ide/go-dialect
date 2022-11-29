package dialect

func appendKingBaseSql(mapping *SqlMapping) {

	appendOracleSql(mapping)

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
	mapping.TableSelect = `
SELECT
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND TABLE_NAME={sqlValuePack(tableName)}
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
	mapping.ColumnUpdate = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} ALTER COLUMN {columnNamePack} TYPE {columnTypePack} [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
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
}
