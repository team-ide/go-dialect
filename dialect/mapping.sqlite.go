package dialect

import "strings"

func NewMappingSqlite() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType:        TypeSqlite,
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
		ColumnComment: `
ALTER TABLE [{ownerName}.]{tableName} CHANGE COLUMN {columnName} {columnName} {columnType} [CHARACTER SET {characterSetName}] [DEFAULT {columnDefault}] {columnNotNull(columnNotNull)} [COMMENT {columnComment}]
`,
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
ALTER TABLE [{ownerName}.]{tableName} ADD [PRIMARY KEY | UNIQUE | FULLTEXT | INDEX] {indexName} ({columnNames}) [COMMENT {columnComment}]
`,
		IndexDelete: `
ALTER TABLE [{ownerName}.]{tableName} DROP INDEX
`,
	}

	/** 数值类型 **/
	/**
	MySQL 支持所有标准 SQL 数值数据类型。
	这些类型包括严格数值数据类型(INTEGER、SMALLINT、DECIMAL 和 NUMERIC)，以及近似数值数据类型(FLOAT、REAL 和 DOUBLE PRECISION)。
	关键字INT是INTEGER的同义词，关键字DEC是DECIMAL的同义词。
	BIT数据类型保存位字段值，并且支持 MyISAM、MEMORY、InnoDB 和 BDB表。
	作为 SQL 标准的扩展，MySQL 也支持整数类型 TINYINT、MEDIUMINT 和 BIGINT。下面的表显示了需要的每个整数类型的存储和范围。

	如果不设置长度，会有默认的长度
	长度代表了显示的最大宽度，如果不够会用0在左边填充，但必须搭配zerofill 使用！
	例如：
	INT(7) 括号中7不是指范围，范围是由数据类型决定的，只是代表显示结果的宽度
	*/

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIT", Format: "BIT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYINT", Format: "TINYINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SMALLINT", Format: "SMALLINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMINT", Format: "MEDIUMINT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT", Format: "INT($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER($l)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIGINT", Format: "BIGINT($l)", IsNumber: true})

	/** 小数 **/

	/**
	M：整数部位+小数部位
	D：小数部位
	如果超过范围，则插入临界值
	M和D都可以省略
	如果是DECIMAL，则M默认为10，D默认为0
	如果是FLOAT和DOUBLE，则会根据插入的数值的精度来决定精度
	定点型的精确度较高，如果要求插入数值的精度较高如货币运算等则考虑使用
	原则：所选择的类型越简单越好，能保存数值的类型越小越好
	*/

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT", Format: "FLOAT($l, $d)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE", Format: "DOUBLE($l, $d)", IsNumber: true})

	/**
	DECIMAL。浮点数类型和定点数类型都可以用（M，N）来表示。其中，M称为精度，表示总共的位数；N称为标度，表示小数的位数.DECIMAL若不指定精度则默认为(10,0)
	不论是定点数还是浮点数类型，如果用户指定的精度超出精度范围，则会四舍五入
	*/

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DEC", Format: "DEC($l, $d)", IsNumber: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DECIMAL", Format: "DOUBLE($l, $d)", IsNumber: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMBER", Format: "NUMBER($l, $d)", IsNumber: true})

	/** 日期/时间类型 **/
	/**
	表示时间值的日期和时间类型为DATETIME、DATE、TIMESTAMP、TIME和YEAR。
	每个时间类型有一个有效值范围和一个"零"值，当指定不合法的MySQL不能表示的值时使用"零"值。
	TIMESTAMP类型有专有的自动更新特性，将在后面描述。
	DATE:
	（1）以‘YYYY-MM-DD’或者‘YYYYMMDD’字符串格式表示的日期，取值范围为‘1000-01-01’～‘9999-12-3’。例如，输入‘2012-12-31’或者‘20121231’，插入数据库的日期都为2012-12-31。
	（2）以‘YY-MM-DD’或者‘YYMMDD’字符串格式表示的日期，在这里YY表示两位的年值。包含两位年值的日期会令人模糊，因为不知道世纪。MySQL使用以下规则解释两位年值：‘00～69’范围的年值转换为‘2000～2069’；‘70～99’范围的年值转换为‘1970～1999’。例如，输入‘12-12-31’，插入数据库的日期为2012-12-31；输入‘981231’，插入数据的日期为1998-12-31。
	（3）以YY-MM-DD或者YYMMDD数字格式表示的日期，与前面相似，00~69范围的年值转换为2000～2069，70～99范围的年值转换为1970～1999。例如，输入12-12-31插入数据库的日期为2012-12-31；输入981231，插入数据的日期为1998-12-31
	*/

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "YEAR", Format: "YEAR", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIME", Format: "TIME", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATETIME", Format: "DATETIME", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			return
		},
	})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
			return
		},
	})

	/** 字符串类型 **/
	/**
	字符串类型指CHAR、VARCHAR、BINARY、VARBINARY、BLOB、TEXT、ENUM和SET。该节描述了这些类型如何工作以及如何在查询中使用这些类型

	注意：char(n) 和 varchar(n) 中括号中 n 代表字符的个数，并不代表字节个数，比如 CHAR(30) 就可以存储 30 个字符。
	CHAR 和 VARCHAR 类型类似，但它们保存和检索的方式不同。它们的最大长度和是否尾部空格被保留等方面也不同。在存储或检索过程中不进行大小写转换。
	BINARY 和 VARBINARY 类似于 CHAR 和 VARCHAR，不同的是它们包含二进制字符串而不要非二进制字符串。也就是说，它们包含字节字符串而不是字符字符串。这说明它们没有字符集，并且排序和比较基于列值字节的数值值。
	BLOB 是一个二进制大对象，可以容纳可变数量的数据。有 4 种 BLOB 类型：TINYBLOB、BLOB、MEDIUMBLOB 和 LONGBLOB。它们区别在于可容纳存储范围不同。
	有 4 种 TEXT 类型：TINYTEXT、TEXT、MEDIUMTEXT 和 LONGTEXT。对应的这 4 种 BLOB 类型，可存储的最大长度不同，可根据实际情况选择。
	*/

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYTEXT", Format: "TINYTEXT($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMTEXT", Format: "MEDIUMTEXT($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGTEXT", Format: "LONGTEXT", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYBLOB", Format: "TINYBLOB($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", Format: "BLOB($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMBLOB", Format: "MEDIUMBLOB($l)", IsString: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGBLOB", Format: "LONGBLOB", IsString: true})

	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ENUM", Format: "VARCHAR($l)", IsString: true, IsExtend: true})
	mapping.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SET", Format: "VARCHAR($l)", IsString: true, IsExtend: true})

	return
}