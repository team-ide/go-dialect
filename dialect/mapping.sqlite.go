package dialect

func NewMappingSqlite() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeSqlite,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",

		// 库或所属者 相关 SQL
		OwnersSelect: `
SELECT 
	name ownerName
FROM pragma_database_list AS t_i 
ORDER BY name
`,
		OwnerSelect: `
SELECT 
	name ownerName
FROM pragma_database_list AS t_i 
WHERE name={sqlValuePack(ownerName)}
`,
		OwnerCreate: ``,
		OwnerDelete: ``,

		// 表 相关 SQL
		TablesSelect: `
SELECT 
	name tableName,
    sql 
FROM sqlite_master 
WHERE type ='table'
ORDER BY name
`,
		TableSelect: `
SELECT 
	name tableName,
    sql 
FROM sqlite_master 
WHERE type ='table'
  AND name={sqlValuePack(tableName)}
`,
		TableCreate: `
CREATE TABLE [{ownerNamePack}.]{tableNamePack}(
{ tableCreateColumnContent }
{ tableCreatePrimaryKeyContent }
)
`,
		TableCreateColumn: `
	{columnNamePack} {columnTypePack} [CHARACTER SET {columnCharacterSetName}] [DEFAULT {columnDefaultPack}] {columnNotNull(columnNotNull)}
`,
		TableCreatePrimaryKey: `
PRIMARY KEY ({primaryKeysPack})
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
	name columnName,
	dflt_value columnDefault,
	"notnull" isNotNull,
	type columnType
FROM pragma_table_info({tableNamePack}) AS t_i 
`,
		ColumnSelect: `
SELECT 
	name columnName,
	dflt_value columnDefault,
	"notnull" isNotNull,
	type columnType
FROM pragma_table_info({tableNamePack}) AS t_i 
WHERE name={sqlValuePack(columnName)}
`,
		ColumnAdd: `
ALTER TABLE [{ownerName}.]{tableName} ADD COLUMN {columnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] {columnNotNull(columnNotNull)} [COMMENT {columnComment}]
`,
		ColumnComment: ``,
		ColumnDelete: `
ALTER TABLE [{ownerName}.]{tableName} DROP COLUMN {columnName}
`,
		ColumnRename: `
ALTER TABLE [{ownerName}.]{tableName} CHANGE COLUMN {oldColumnName} {newColumnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] {columnNotNull(columnNotNull)} [COMMENT {columnComment}]
`,
		ColumnUpdate: `
ALTER TABLE [{ownerName}.]{tableName} CHANGE COLUMN {columnName} {columnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] {columnNotNull(columnNotNull)} [COMMENT {columnComment}] [AFTER {columnAfter}]
`,

		// 主键 相关 SQL
		PrimaryKeysSelect: `
SELECT 
	a.name indexName,
	b.name columnName 
FROM pragma_index_list({tableNamePack}) AS a,pragma_index_info(a.name) b 
WHERE a.origin = "pk"
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
	a.name indexName,
	a."unique",
	b.name columnName 
FROM pragma_index_list({tableNamePack}) AS a,pragma_index_info(a.name) b 
WHERE a.origin != "pk"
`,
		IndexAdd: `
CREATE {indexType} [{indexNamePack}] ON {tableNamePack}({columnNamesPack})
`,
		IndexDelete: `
ALTER TABLE [{ownerNamePack}.]{tableNamePack} DROP INDEX {indexNamePack}
`,
	}

	AppendSqliteColumnType(mapping)
	AppendSqliteIndexType(mapping)
	return
}

func AppendSqliteColumnType(mapping *SqlMapping) {

	// mysql 数据类型转换
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIT", Format: "BIT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYINT", Format: "TINYINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SMALLINT", Format: "SMALLINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMINT", Format: "MEDIUMINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT", Format: "INT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIGINT", Format: "BIGINT($l)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT", Format: "FLOAT($l, $d)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE", Format: "DOUBLE($l, $d)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DEC", Format: "DEC($l, $d)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DECIMAL", Format: "DECIMAL($l, $d)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "YEAR", Format: "YEAR", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIME", Format: "TIME", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATETIME", Format: "DATETIME", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYTEXT", Format: "TINYTEXT", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMTEXT", Format: "MEDIUMTEXT", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGTEXT", Format: "LONGTEXT", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ENUM", Format: "ENUM", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYBLOB", Format: "TINYBLOB", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMBLOB", Format: "MEDIUMBLOB", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGBLOB", Format: "LONGBLOB", IsString: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SET", Format: "TEXT", IsString: true})

	// 浮点数
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REAL", Format: "REAL", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMERIC", IsNumber: true})
}

func AppendSqliteIndexType(mapping *SqlMapping) {

	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "", Format: "INDEX"})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "INDEX", Format: "INDEX"})
	mapping.AddIndexTypeInfo(&IndexTypeInfo{Name: "UNIQUE", Format: "UNIQUE",
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
