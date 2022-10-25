package dialect

import (
	"errors"
	"strings"
)

func NewSqliteDialect() *SqliteDialect {

	res := &SqliteDialect{
		DefaultDialect: NewDefaultDialect(SqliteType),
	}
	res.init()
	return res
}

type SqliteDialect struct {
	*DefaultDialect
}

func (this_ *SqliteDialect) init() {
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

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIT", TypeFormat: "BIT($l)", HasLength: false, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "TINYINT($l)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "SMALLINT($l)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "MEDIUMINT($l)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT", TypeFormat: "INT($l)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "INTEGER($l)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "BIGINT($l)", HasLength: true, IsNumber: true})

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

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "FLOAT($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "DOUBLE($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})

	/**
	DECIMAL。浮点数类型和定点数类型都可以用（M，N）来表示。其中，M称为精度，表示总共的位数；N称为标度，表示小数的位数.DECIMAL若不指定精度则默认为(10,0)
	不论是定点数还是浮点数类型，如果用户指定的精度超出精度范围，则会四舍五入
	*/

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DEC", TypeFormat: "DEC($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "DOUBLE($l, $d)", HasLength: true, IsNumber: true})

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

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

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "YEAR", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIME", TypeFormat: "TIME", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATETIME", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/
	/**
	字符串类型指CHAR、VARCHAR、BINARY、VARBINARY、BLOB、TEXT、ENUM和SET。该节描述了这些类型如何工作以及如何在查询中使用这些类型

	注意：char(n) 和 varchar(n) 中括号中 n 代表字符的个数，并不代表字节个数，比如 CHAR(30) 就可以存储 30 个字符。
	CHAR 和 VARCHAR 类型类似，但它们保存和检索的方式不同。它们的最大长度和是否尾部空格被保留等方面也不同。在存储或检索过程中不进行大小写转换。
	BINARY 和 VARBINARY 类似于 CHAR 和 VARCHAR，不同的是它们包含二进制字符串而不要非二进制字符串。也就是说，它们包含字节字符串而不是字符字符串。这说明它们没有字符集，并且排序和比较基于列值字节的数值值。
	BLOB 是一个二进制大对象，可以容纳可变数量的数据。有 4 种 BLOB 类型：TINYBLOB、BLOB、MEDIUMBLOB 和 LONGBLOB。它们区别在于可容纳存储范围不同。
	有 4 种 TEXT 类型：TINYTEXT、TEXT、MEDIUMTEXT 和 LONGTEXT。对应的这 4 种 BLOB 类型，可存储的最大长度不同，可根据实际情况选择。
	*/

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "TINYTEXT($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "TEXT($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "MEDIUMTEXT($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "LONGTEXT", HasLength: false, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "ENUM($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "TINYBLOB($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "MEDIUMBLOB($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "LONGBLOB", HasLength: false, IsString: true})

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})

	this_.AddFuncTypeInfo(&FuncTypeInfo{Name: "md5", Format: "md5"})
}

func (this_ *SqliteDialect) DatabaseCreateSql(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {

	return
}
func (this_ *SqliteDialect) DatabaseDeleteSql(param *GenerateParam, databaseName string) (sqlList []string, err error) {

	return
}
func (this_ *SqliteDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	if data == nil {
		return
	}
	table = &TableModel{}
	if data["name"] != nil {
		table.Name = data["name"].(string)
	}
	if data["sql"] != nil {
		table.Sql = data["sql"].(string)
	}
	return
}
func (this_ *SqliteDialect) TablesSelectSql(databaseName string) (sql string, err error) {
	sql = `SELECT * FROM sqlite_master WHERE type ='table' `
	sql += `ORDER BY name`
	return
}
func (this_ *SqliteDialect) TableSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `SELECT * FROM sqlite_master WHERE type ='table' `
	sql += `AND name='` + tableName + `' `
	sql += `ORDER BY name`
	return
}
func (this_ *SqliteDialect) TableCreateSql(param *GenerateParam, databaseName string, table *TableModel) (sqlList []string, err error) {

	createTableSql := `CREATE TABLE `

	if param.AppendDatabase && databaseName != "" {
		createTableSql += param.packingCharacterDatabase(databaseName) + "."
	}
	createTableSql += param.packingCharacterTable(table.Name)

	createTableSql += `(`
	createTableSql += "\n"
	primaryKeys := ""
	if len(table.ColumnList) > 0 {
		for _, column := range table.ColumnList {
			var columnSql = param.packingCharacterColumn(column.Name)

			var columnType string
			columnType, err = this_.FormatColumnType(column.Type, column.Length, column.Decimal)
			if err != nil {
				return
			}
			columnSql += " " + columnType

			if column.Default != "" {
				columnSql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
			}
			if column.NotNull {
				columnSql += ` NOT NULL`
			}

			if column.PrimaryKey {
				primaryKeys += "" + column.Name + ","
			}
			createTableSql += "\t" + columnSql + ",\n"
		}
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + param.packingCharacterColumns(primaryKeys) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`

	sqlList = append(sqlList, createTableSql)

	if len(table.IndexList) > 0 {
		for _, one := range table.IndexList {
			var sqlList_ []string
			sqlList_, err = this_.IndexAddSql(param, databaseName, table.Name, one)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}
	return
}
func (this_ *SqliteDialect) TableCommentSql(param *GenerateParam, databaseName string, tableName string, comment string) (sqlList []string, err error) {

	return
}
func (this_ *SqliteDialect) TableDeleteSql(param *GenerateParam, databaseName string, tableName string) (sqlList []string, err error) {
	var sql string
	sql = `DROP TABLE `

	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *SqliteDialect) ColumnModel(data map[string]interface{}) (column *ColumnModel, err error) {
	if data == nil {
		return
	}
	column = &ColumnModel{}
	if data["name"] != nil {
		column.Name = data["name"].(string)
	}
	if data["dflt_value"] != nil {
		column.Default = GetStringValue(data["dflt_value"])
	}
	if GetStringValue(data["dflt_value"]) == "1" {
		column.PrimaryKey = true
	}
	if GetStringValue(data["notnull"]) == "1" {
		column.NotNull = true
	}

	var columnTypeInfo *ColumnTypeInfo
	if data["type"] != nil {
		columnType := data["type"].(string)
		columnTypeInfo, column.Length, column.Decimal, err = this_.ToColumnTypeInfo(columnType)
		if err != nil {
			return
		}
		column.Type = columnTypeInfo.Name
	}
	return
}
func (this_ *SqliteDialect) ColumnsSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `select * from pragma_table_info("` + tableName + `") as t_i `
	return
}
func (this_ *SqliteDialect) ColumnAddSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.FormatColumnType(column.Type, column.Length, column.Decimal)
	if err != nil {
		return
	}

	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)

	sql += ` ADD COLUMN `
	sql += param.packingCharacterColumn(column.Name)
	sql += ` ` + columnType + ``
	if column.Default != "" {
		sql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
	}
	if column.NotNull {
		sql += ` NOT NULL`
	}
	sql += ``

	sqlList = append(sqlList, sql)

	return
}
func (this_ *SqliteDialect) ColumnCommentSql(param *GenerateParam, databaseName string, tableName string, columnName string, comment string) (sqlList []string, err error) {

	return
}
func (this_ *SqliteDialect) columnRenameSql(param *GenerateParam, databaseName string, tableName string, oldName string, newName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)

	sql += ` RENAME COLUMN `
	sql += param.packingCharacterColumn(oldName)
	sql += ` TO `
	sql += param.packingCharacterColumn(newName)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *SqliteDialect) ColumnUpdateSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error) {

	return
}
func (this_ *SqliteDialect) ColumnDeleteSql(param *GenerateParam, databaseName string, tableName string, columnName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += param.packingCharacterTable(tableName)

	sql += ` DROP COLUMN `
	sql += param.packingCharacterColumn(columnName)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *SqliteDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
	if data == nil {
		return
	}
	primaryKey = &PrimaryKeyModel{}
	if data["name"] != nil {
		primaryKey.ColumnName = data["name"].(string)
	}
	return
}
func (this_ *SqliteDialect) PrimaryKeysSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `select * from pragma_table_info("` + tableName + `") as t_i where t_i.pk=1 `
	return
}

func (this_ *SqliteDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
	if data == nil {
		return
	}
	index = &IndexModel{}
	if data["name"] != nil {
		index.Name = data["name"].(string)
	}
	if data["name"] != nil {
		index.ColumnName = data["name"].(string)
	}
	if GetStringValue(data["unique"]) == "1" {
		index.Type = "unique"
	}
	return
}
func (this_ *SqliteDialect) IndexesSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `select * from pragma_index_list("` + tableName + `") as t_i where origin!="pk"  `
	return
}
func (this_ *SqliteDialect) IndexAddSql(param *GenerateParam, databaseName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	sql := "CREATE "
	switch strings.ToUpper(index.Type) {
	case "UNIQUE":
		sql += "UNIQUE INDEX"
	case "":
		sql += "INDEX"
	default:
		err = errors.New("dialect [" + this_.DialectType().Name + "] not support index type [" + index.Type + "]")
		return
	}

	sql += " " + param.packingCharacterColumn(index.Name) + ""

	sql += " ON "
	if param.AppendDatabase && databaseName != "" {
		sql += param.packingCharacterDatabase(databaseName) + "."
	}
	sql += "" + param.packingCharacterTable(tableName)

	sql += "(" + param.packingCharacterColumns(strings.Join(index.Columns, ",")) + ")"

	sqlList = append(sqlList, sql)
	return
}
