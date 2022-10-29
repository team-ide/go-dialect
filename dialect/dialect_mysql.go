package dialect

import (
	"errors"
	"strings"
)

func NewMysqlDialect() *MysqlDialect {

	res := &MysqlDialect{
		DefaultDialect: NewDefaultDialect(MysqlType),
	}
	res.init()
	return res
}

type MysqlDialect struct {
	*DefaultDialect
}

func (this_ *MysqlDialect) init() {
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

func (this_ *MysqlDialect) PackOwner(ownerName string) string {
	return packingName("`", ownerName)
}

func (this_ *MysqlDialect) PackTable(tableName string) string {
	return packingName("`", tableName)
}

func (this_ *MysqlDialect) PackColumn(columnName string) string {
	return packingName("`", columnName)
}

func (this_ *MysqlDialect) PackColumns(columnNames []string) string {
	return packingNames("`", columnNames)
}

func (this_ *MysqlDialect) PackValue(column *ColumnModel, value interface{}) string {
	var columnTypeInfo *ColumnTypeInfo
	if column != nil {
		if strings.EqualFold(column.Type, "ENUM") {
			if value == nil || value == "" {
				if !column.NotNull {
					return "NULL"
				}
			}
		}
		columnTypeInfo, _ = this_.GetColumnTypeInfo(column.Type)
	}
	return packingValue(columnTypeInfo, `'`, `\`, value)
}

func (this_ *MysqlDialect) FormatColumnType(column *ColumnModel) (columnType string, err error) {
	if strings.EqualFold(column.Type, "SET") || strings.EqualFold(column.Type, "ENUM") {
		columnType = column.Type + "(" + packingValues("'", column.Defaults) + ")"
		return
	}
	columnTypeInfo, err := this_.GetColumnTypeInfo(column.Type)
	if err != nil {
		return
	}
	columnType = columnTypeInfo.FormatColumnType(column.Length, column.Decimal)
	return
}
func (this_ *MysqlDialect) FormatDefaultValue(column *ColumnModel) (defaultValue string) {
	defaultValue = "DEFAULT "
	if column.DefaultCurrentTimestamp || column.OnUpdateCurrentTimestamp {
		if column.DefaultCurrentTimestamp {
			defaultValue += " CURRENT_TIMESTAMP "
		}
		if column.OnUpdateCurrentTimestamp {
			defaultValue += " ON UPDATE CURRENT_TIMESTAMP "
		}
		return
	}
	defaultValue += this_.PackValue(nil, column.Default)
	return
}
func (this_ *MysqlDialect) OwnerModel(data map[string]interface{}) (owner *OwnerModel, err error) {
	if data == nil {
		return
	}
	owner = &OwnerModel{}
	if data["SCHEMA_NAME"] != nil {
		owner.Name = data["SCHEMA_NAME"].(string)
	}
	if data["DEFAULT_CHARACTER_SET_NAME"] != nil {
		owner.CharacterSetName = data["DEFAULT_CHARACTER_SET_NAME"].(string)
	}
	if data["DEFAULT_COLLATION_NAME"] != nil {
		owner.CollationName = data["DEFAULT_COLLATION_NAME"].(string)
	}
	return
}
func (this_ *MysqlDialect) OwnersSelectSql() (sql string, err error) {
	sql = `SELECT SCHEMA_NAME,DEFAULT_CHARACTER_SET_NAME,DEFAULT_COLLATION_NAME from information_schema.schemata ORDER BY SCHEMA_NAME`
	return
}
func (this_ *MysqlDialect) OwnerSelectSql(ownerName string) (sql string, err error) {
	sql = `SELECT SCHEMA_NAME,DEFAULT_CHARACTER_SET_NAME,DEFAULT_COLLATION_NAME from information_schema.schemata `
	sql += ` WHERE SCHEMA_NAME='` + ownerName + `'`
	return
}
func (this_ *MysqlDialect) OwnerCreateSql(owner *OwnerModel) (sqlList []string, err error) {
	var sql string
	sql = `CREATE DATABASE ` + this_.PackOwner(owner.Name)
	if owner.CharacterSetName != "" {
		sql += ` CHARACTER SET ` + owner.CharacterSetName
	}
	if owner.CollationName != "" {
		sql += ` COLLATE '` + owner.CollationName + "'"
	}

	sqlList = append(sqlList, sql)
	return
}
func (this_ *MysqlDialect) OwnerDeleteSql(ownerName string) (sqlList []string, err error) {
	var sql string
	sql = `DROP DATABASE IF EXISTS ` + this_.PackOwner(ownerName)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *MysqlDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	if data == nil {
		return
	}
	table = &TableModel{}
	if data["TABLE_NAME"] != nil {
		table.Name = data["TABLE_NAME"].(string)
	}
	if data["TABLE_COMMENT"] != nil {
		table.Comment = data["TABLE_COMMENT"].(string)
	}
	if data["TABLE_SCHEMA"] != nil {
		table.OwnerName = data["TABLE_SCHEMA"].(string)
	}
	return
}
func (this_ *MysqlDialect) TablesSelectSql(ownerName string) (sql string, err error) {
	sql = `SELECT TABLE_NAME,TABLE_COMMENT,TABLE_SCHEMA from information_schema.tables `
	if ownerName != "" {
		sql += `WHERE TABLE_SCHEMA='` + ownerName + `' `
	}
	sql += `ORDER BY TABLE_NAME`
	return
}
func (this_ *MysqlDialect) TableSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT TABLE_NAME,TABLE_COMMENT,TABLE_SCHEMA from information_schema.tables `
	sql += `WHERE 1=1 `
	if ownerName != "" {
		sql += `AND TABLE_SCHEMA='` + ownerName + `' `
	}
	sql += `AND TABLE_NAME='` + tableName + `' `
	sql += `ORDER BY TABLE_NAME`
	return
}
func (this_ *MysqlDialect) TableCreateSql(ownerName string, table *TableModel) (sqlList []string, err error) {
	sqlList = []string{}

	createTableSql := `CREATE TABLE `

	if ownerName != "" {
		createTableSql += this_.PackOwner(ownerName) + "."
	}
	createTableSql += this_.PackTable(table.Name)

	createTableSql += `(`
	createTableSql += "\n"
	primaryKeys := ""
	for _, column := range table.ColumnList {
		var columnSql = this_.PackColumn(column.Name)
		var columnType string
		columnType, err = this_.FormatColumnType(column)
		if err != nil {
			return
		}

		columnSql += " " + columnType

		if column.CharacterSetName != "" {
			columnSql += ` CHARACTER SET ` + column.CharacterSetName
		}
		if column.Default != "" {
			columnSql += " " + this_.FormatDefaultValue(column)
		}
		if column.NotNull {
			columnSql += ` NOT NULL`
		}
		if column.Comment != "" {
			columnSql += " COMMENT " + this_.PackValue(nil, column.Comment)
		}

		if column.PrimaryKey {
			primaryKeys += "" + column.Name + ","
		}
		createTableSql += "\t" + columnSql
		createTableSql += ",\n"
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + this_.PackColumns(strings.Split(primaryKeys, ",")) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`
	if table.CharacterSetName != "" {
		createTableSql += ` DEFAULT CHARSET ` + table.CharacterSetName
	}

	sqlList = append(sqlList, createTableSql)

	var sqlList_ []string
	// 添加注释
	if table.Comment != "" {
		sqlList_, err = this_.TableCommentSql(ownerName, table.Name, table.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	for _, one := range table.IndexList {
		sqlList_, err = this_.IndexAddSql(ownerName, table.Name, one)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)

	}
	return
}
func (this_ *MysqlDialect) TableCommentSql(ownerName string, tableName string, comment string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)
	sql += " COMMENT " + this_.PackValue(nil, comment)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *MysqlDialect) TableRenameSql(ownerName string, oldTableName string, newTableName string) (sqlList []string, err error) {
	sql := "ALTER TABLE  "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(oldTableName)
	sql += " RENAME AS  "
	sql += "" + this_.PackTable(newTableName)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *MysqlDialect) TableDeleteSql(ownerName string, tableName string) (sqlList []string, err error) {
	var sql string
	sql = `DROP TABLE IF EXISTS `

	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += this_.PackTable(tableName)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *MysqlDialect) ColumnModel(data map[string]interface{}) (column *ColumnModel, err error) {
	if data == nil {
		return
	}
	column = &ColumnModel{}
	if data["COLUMN_NAME"] != nil {
		column.Name = data["COLUMN_NAME"].(string)
	}
	if data["COLUMN_COMMENT"] != nil {
		column.Comment = data["COLUMN_COMMENT"].(string)
	}
	if data["COLUMN_DEFAULT"] != nil {
		column.Default = GetStringValue(data["COLUMN_DEFAULT"])
		if strings.Contains(strings.ToLower(column.Default), "current_timestamp") {
			column.DefaultCurrentTimestamp = true
		} else if strings.Contains(strings.ToLower(column.Default), "0000-00-00 00:00:00") {
			column.DefaultCurrentTimestamp = true
		}
	}
	if data["TABLE_NAME"] != nil {
		column.TableName = data["TABLE_NAME"].(string)
	}
	if data["TABLE_SCHEMA"] != nil {
		column.OwnerName = data["TABLE_SCHEMA"].(string)
	}
	if data["CHARACTER_SET_NAME"] != nil {
		column.CharacterSetName = data["CHARACTER_SET_NAME"].(string)
	}
	if data["EXTRA"] != nil {
		column.Extra = data["EXTRA"].(string)
		if strings.Contains(strings.ToLower(column.Extra), "on update current_timestamp") {
			column.OnUpdateCurrentTimestamp = true
		}
	}

	if GetStringValue(data["IS_NULLABLE"]) == "NO" {
		column.NotNull = true
	}
	var columnTypeInfo *ColumnTypeInfo
	if data["COLUMN_TYPE"] != nil {
		columnType := data["COLUMN_TYPE"].(string)
		columnTypeInfo, column.Length, column.Decimal, err = this_.ToColumnTypeInfo(columnType)
		if err != nil {
			return
		}
		column.Type = columnTypeInfo.Name

		if strings.EqualFold(column.Type, "SET") || strings.EqualFold(column.Type, "ENUM") {
			if strings.Contains(columnType, "(") {
				setStr := columnType[strings.Index(columnType, "(")+1 : strings.Index(columnType, ")")]
				setStr = strings.ReplaceAll(setStr, "'", "")
				column.Defaults = strings.Split(setStr, ",")
			}
		}
		dataType := data["DATA_TYPE"].(string)
		if !strings.EqualFold(dataType, column.Type) {
			err = errors.New("column type [" + columnType + "] not eq data type [" + dataType + "]")
			return
		}
	}
	return
}
func (this_ *MysqlDialect) ColumnsSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT COLUMN_NAME,COLUMN_COMMENT,COLUMN_DEFAULT,EXTRA,TABLE_NAME,TABLE_SCHEMA,CHARACTER_SET_NAME,IS_NULLABLE,COLUMN_TYPE,DATA_TYPE from information_schema.columns `
	sql += `WHERE 1=1 `
	if ownerName != "" {
		sql += `AND TABLE_SCHEMA='` + ownerName + `' `
	}
	sql += `AND TABLE_NAME='` + tableName + `' `
	return
}
func (this_ *MysqlDialect) ColumnAddSql(ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.FormatColumnType(column)
	if err != nil {
		return
	}

	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)
	sql += " ADD COLUMN " + this_.PackColumn(column.Name)
	sql += " " + columnType
	if column.Default != "" {
		sql += " " + this_.FormatDefaultValue(column)
	}
	if column.NotNull {
		sql += " NOT NULL"
	}
	sql += " COMMENT " + this_.PackValue(nil, column.Comment)
	if column.BeforeColumn != "" {
		sql += " AFTER " + this_.PackColumn(column.BeforeColumn)
	}

	sqlList = append(sqlList, sql)
	return
}
func (this_ *MysqlDialect) ColumnUpdateSql(ownerName string, tableName string, oldColumn *ColumnModel, newColumn *ColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.FormatColumnType(newColumn)
	if err != nil {
		return
	}

	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	sql += " CHANGE COLUMN"
	sql += " " + this_.PackColumn(oldColumn.Name)
	sql += " " + this_.PackColumn(newColumn.Name)
	sql += " " + columnType

	if newColumn.Default == "" {
		sql += " DEFAULT NULL"
	} else {
		sql += " " + this_.FormatDefaultValue(newColumn)
	}
	if newColumn.NotNull {
		sql += " NOT NULL"
	}
	sql += " COMMENT " + this_.PackValue(nil, newColumn.Comment)
	if newColumn.BeforeColumn != "" {
		sql += " AFTER " + this_.PackColumn(newColumn.BeforeColumn)
	}

	sqlList = append(sqlList, sql)

	return
}
func (this_ *MysqlDialect) ColumnDeleteSql(ownerName string, tableName string, columnName string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += this_.PackTable(tableName)

	sql += ` DROP COLUMN `
	sql += this_.PackColumn(columnName)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *MysqlDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
	if data == nil {
		return
	}
	primaryKey = &PrimaryKeyModel{}
	if data["COLUMN_NAME"] != nil {
		primaryKey.ColumnName = data["COLUMN_NAME"].(string)
	}
	if data["TABLE_NAME"] != nil {
		primaryKey.TableName = data["TABLE_NAME"].(string)
	}
	if data["TABLE_SCHEMA"] != nil {
		primaryKey.OwnerName = data["TABLE_SCHEMA"].(string)
	}
	return
}
func (this_ *MysqlDialect) PrimaryKeysSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT k.COLUMN_NAME,t.TABLE_NAME,t.TABLE_SCHEMA from information_schema.table_constraints t `
	sql += `JOIN information_schema.key_column_usage k USING (CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME) `
	sql += `WHERE 1=1 `
	if ownerName != "" {
		sql += `AND t.TABLE_SCHEMA='` + ownerName + `' `
	}
	sql += `AND t.TABLE_NAME='` + tableName + `' `
	sql += `AND t.CONSTRAINT_TYPE='PRIMARY KEY' `
	return
}
func (this_ *MysqlDialect) PrimaryKeyAddSql(ownerName string, tableName string, primaryKeys []string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	sql += ` ADD PRIMARY KEY `

	sql += "(" + this_.PackColumns(primaryKeys) + ")"

	sqlList = append(sqlList, sql)
	return
}
func (this_ *MysqlDialect) PrimaryKeyDeleteSql(ownerName string, tableName string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	sql += ` DROP PRIMARY KEY `

	sqlList = append(sqlList, sql)
	return
}

func (this_ *MysqlDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
	if data == nil {
		return
	}
	index = &IndexModel{}
	if data["INDEX_NAME"] != nil {
		index.Name = data["INDEX_NAME"].(string)
	}
	if data["COLUMN_NAME"] != nil {
		index.ColumnName = data["COLUMN_NAME"].(string)
	}
	if data["INDEX_COMMENT"] != nil {
		index.Comment = data["INDEX_COMMENT"].(string)
	}
	if GetStringValue(data["NON_UNIQUE"]) == "0" {
		index.Type = "unique"
	}
	if data["TABLE_NAME"] != nil {
		index.TableName = data["TABLE_NAME"].(string)
	}
	if data["TABLE_SCHEMA"] != nil {
		index.OwnerName = data["TABLE_SCHEMA"].(string)
	}
	return
}
func (this_ *MysqlDialect) IndexesSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT INDEX_NAME,COLUMN_NAME,INDEX_COMMENT,NON_UNIQUE,TABLE_NAME,TABLE_SCHEMA from information_schema.statistics `
	sql += `WHERE 1=1 `
	if ownerName != "" {
		sql += `AND TABLE_SCHEMA='` + ownerName + `' `
	}
	sql += `AND TABLE_NAME='` + tableName + `' `
	sql += `AND INDEX_NAME NOT IN(`
	sql += `SELECT t.CONSTRAINT_NAME from information_schema.table_constraints t `
	sql += `JOIN information_schema.key_column_usage k USING (CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME) `
	sql += `WHERE 1=1 `
	if ownerName != "" {
		sql += `AND t.TABLE_SCHEMA='` + ownerName + `' `
	}
	sql += `AND t.TABLE_NAME='` + tableName + `' `
	sql += `AND t.CONSTRAINT_TYPE='PRIMARY KEY' `
	sql += `) `
	return
}

// ALTER TABLE [${ownerName}.]${tableName} ADD ${'PRIMARY KEY' | 'UNIQUE' | 'INDEX'} ${indexName} (${columns}) [COMMENT '${indexCommend}']

func (this_ *MysqlDialect) IndexAddSql(ownerName string, tableName string, index *IndexModel) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	switch strings.ToUpper(index.Type) {
	case "PRIMARY":
		sql += " ADD PRIMARY KEY "
	case "UNIQUE":
		sql += " ADD UNIQUE "
	case "FULLTEXT":
		sql += " ADD FULLTEXT "
	case "":
		sql += " ADD INDEX "
	default:
		err = errors.New("dialect [" + this_.DialectType().Name + "] not support index type [" + index.Type + "]")
		return
	}
	if index.Name != "" {
		sql += "" + this_.PackColumn(index.Name) + " "
	}
	if len(index.Columns) > 0 {
		sql += "(" + this_.PackColumns(index.Columns) + ")"
	}

	if index.Comment != "" {
		sql += " COMMENT " + this_.PackValue(nil, index.Comment)
	}

	sqlList = append(sqlList, sql)
	return
}

// ALTER TABLE [${ownerName}.]${tableName} DROP INDEX ${indexName}

func (this_ *MysqlDialect) IndexDeleteSql(ownerName string, tableName string, indexName string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if ownerName != "" {
		sql += this_.PackOwner(ownerName) + "."
	}
	sql += "" + this_.PackTable(tableName)

	sql += ` DROP INDEX `
	sql += "" + this_.PackColumn(indexName)

	sqlList = append(sqlList, sql)
	return
}
