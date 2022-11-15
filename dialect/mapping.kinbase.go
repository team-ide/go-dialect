package dialect

func NewMappingKinBase() (mapping *SqlMapping) {

	// http://www.yaotu.net/biancheng/21946.html
	// https://www.modb.pro/db/442114
	mapping = NewMappingOracle()
	mapping.dialectType = TypeKinBase

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
CREATE SCHEMA {ownerName} 
`
	mapping.OwnerDelete = `
DROP SCHEMA {ownerName} CASCADE
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
    t1.CONSTRAINT_NAME indexName,
    t2.COLUMN_NAME columnName,
    t1.TABLE_NAME tableName,
    t1.TABLE_SCHEMA ownerName
FROM information_schema.table_constraints t1
LEFT JOIN information_schema.key_column_usage t2 
ON (t2.CONSTRAINT_NAME=t1.CONSTRAINT_NAME AND t2.TABLE_SCHEMA=t1.TABLE_SCHEMA AND t2.TABLE_NAME=t1.TABLE_NAME)
WHERE t1.TABLE_SCHEMA={sqlValuePack(ownerName)}
  AND t1.TABLE_NAME={sqlValuePack(tableName)}
  AND (t1.CONSTRAINT_TYPE !='PRIMARY KEY' OR t1.CONSTRAINT_TYPE = '' OR t1.CONSTRAINT_TYPE IS NULL)
`
	return
}
