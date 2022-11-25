package dialect

func appendSqliteSql(mapping *SqlMapping) {
	// 库或所属者 相关 SQL
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
	mapping.OwnerCreate = ``
	mapping.OwnerDelete = ``

	// 表 相关 SQL
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
	mapping.TableComment = ``
	mapping.TableRename = `
ALTER TABLE [{ownerName}.]{oldTableName} RENAME AS {newTableName}
`
	mapping.TableDelete = `
DROP TABLE IF EXISTS [{ownerName}.]{tableName}
`

	// 字段 相关 SQL
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
	mapping.ColumnRename = `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} RENAME COLUMN {oldColumnNamePack} TO {columnNamePack}
`
	mapping.ColumnComment = ``
	mapping.ColumnUpdateHasRename = false
	mapping.ColumnUpdateHasComment = false
	mapping.ColumnUpdateHasAfter = false
	mapping.ColumnUpdate = `
`

	// 主键 相关 SQL
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

	// 索引 相关 SQL
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
}
