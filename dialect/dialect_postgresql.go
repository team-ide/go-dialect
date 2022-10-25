package dialect

import (
	"strings"
)

func NewPostgresqlDialect() *PostgresqlDialect {

	res := &PostgresqlDialect{
		DefaultDialect: NewDefaultDialect(PostgresqlType),
	}
	res.init()
	return res
}

type PostgresqlDialect struct {
	*DefaultDialect
}

func (this_ *PostgresqlDialect) init() {
	/** 数值类型 **/

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIT", TypeFormat: "NUMBER($l, $d)", HasLength: false, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 小数 **/

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DEC", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "DATE", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIME", TypeFormat: "DATE", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATE", IsDateTime: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR2($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARCHAR2", TypeFormat: "VARCHAR2($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "VARCHAR2($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "VARCHAR2($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "CLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "CLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLOB", TypeFormat: "CLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XMLTYPE", TypeFormat: "XMLTYPE($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "RAW", TypeFormat: "RAW($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NVARCHAR2", TypeFormat: "NVARCHAR2($l)", HasLength: true, IsString: true})

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NUMERIC", TypeFormat: "NUMERIC($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "OID", TypeFormat: "OID($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NAME", TypeFormat: "NAME($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BOOL", TypeFormat: "BOOL($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT1", TypeFormat: "INT1($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT2", TypeFormat: "INT2($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT4", TypeFormat: "INT4($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT8", TypeFormat: "INT8($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS_LSN", TypeFormat: "SYS_LSN($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGCLASS", TypeFormat: "REGCLASS($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TIMESTAMPTZ", TypeFormat: "TIMESTAMPTZ($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_TEXT", TypeFormat: "_TEXT", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "JSON", TypeFormat: "JSON", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SYS_NODE_TREE", TypeFormat: "SYS_NODE_TREE", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "character_data", TypeFormat: "character_data", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "yes_or_no", TypeFormat: "yes_or_no", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "cardinal_number", TypeFormat: "cardinal_number", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INTERVAL", TypeFormat: "INTERVAL($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "REGPROC", TypeFormat: "REGPROC($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_ACLITEM", TypeFormat: "_ACLITEM", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT4", TypeFormat: "FLOAT4($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "FLOAT8", TypeFormat: "FLOAT8($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XID", TypeFormat: "XID($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TDEKEY", TypeFormat: "TDEKEY($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT2", TypeFormat: "_INT2", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_INT4", TypeFormat: "_INT4", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_OID", TypeFormat: "_OID", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "INT2VECTOR", TypeFormat: "INT2VECTOR", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "OIDVECTOR", TypeFormat: "OIDVECTOR", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTEA", TypeFormat: "BYTEA", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_CHAR", TypeFormat: "_CHAR", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_FLOAT4", TypeFormat: "_FLOAT4", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "_FLOAT8", TypeFormat: "_FLOAT8", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ANYARRAY", TypeFormat: "ANYARRAY", HasLength: true, IsString: true})

	// Postgresql

	this_.AddFuncTypeInfo(&FuncTypeInfo{Name: "md5", Format: "md5"})
}

func (this_ *PostgresqlDialect) DatabaseModel(data map[string]interface{}) (database *DatabaseModel, err error) {
	if data == nil {
		return
	}
	database = &DatabaseModel{}
	if data["USERNAME"] != nil {
		database.Name = data["USERNAME"].(string)
	}
	return
}
func (this_ *PostgresqlDialect) DatabasesSelectSql() (sql string, err error) {
	sql = `SELECT * FROM dba_users ORDER BY USERNAME`
	return
}

func (this_ *PostgresqlDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	if data == nil {
		return
	}
	table = &TableModel{}
	if data["TABLE_NAME"] != nil {
		table.Name = data["TABLE_NAME"].(string)
	}
	return
}
func (this_ *PostgresqlDialect) TablesSelectSql(databaseName string) (sql string, err error) {
	sql = `SELECT * FROM all_tables  `
	if databaseName != "" {
		sql += `WHERE OWNER ='` + databaseName + `' `
	}
	sql += `ORDER BY TABLE_NAME`
	return
}
func (this_ *PostgresqlDialect) TableSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `SELECT * FROM all_tables `
	sql += `WHERE 1=1 `
	if databaseName != "" {
		sql += `AND owner='` + databaseName + `' `
	}
	sql += `AND TABLE_NAME='` + tableName + `' `
	sql += `ORDER BY TABLE_NAME`
	return
}

func (this_ *PostgresqlDialect) ColumnModel(data map[string]interface{}) (column *ColumnModel, err error) {
	if data == nil {
		return
	}
	column = &ColumnModel{}
	if data["COLUMN_NAME"] != nil {
		column.Name = data["COLUMN_NAME"].(string)
	}
	if data["COMMENTS"] != nil {
		column.Comment = data["COMMENTS"].(string)
	}
	if data["DATA_DEFAULT"] != nil {
		column.Default = GetStringValue(data["DATA_DEFAULT"])
	}
	if data["TABLE_NAME"] != nil {
		column.TableName = data["TABLE_NAME"].(string)
	}
	if data["CHARACTER_SET_NAME"] != nil {
		column.CharacterSetName = data["CHARACTER_SET_NAME"].(string)
	}

	if GetStringValue(data["NULLABLE"]) == "N" {
		column.NotNull = true
	}
	var columnTypeInfo *ColumnTypeInfo
	if data["DATA_TYPE"] != nil {
		dataType := data["DATA_TYPE"].(string)
		if strings.Contains(dataType, "(") {
			dataType = dataType[:strings.Index(dataType, "(")]
		}
		columnTypeInfo, err = this_.GetColumnTypeInfo(dataType)
		if err != nil {
			return
		}
		column.Type = columnTypeInfo.Name

		//bs, _ := json.Marshal(data)
		//println("data:", string(bs))
		dataLength := GetStringValue(data["DATA_LENGTH"])
		if dataLength != "" && dataLength != "0" {
			column.Length, err = StringToInt(dataLength)
			if err != nil {
				return
			}
		}
		dataPrecision := GetStringValue(data["DATA_PRECISION"])
		if dataPrecision != "" && dataPrecision != "0" {
			column.Length, err = StringToInt(dataPrecision)
			if err != nil {
				return
			}
		}
		dataScale := GetStringValue(data["DATA_SCALE"])
		if dataScale != "" && dataScale != "0" {
			column.Decimal, err = StringToInt(dataScale)
			if err != nil {
				return
			}
		}
	}
	return
}
func (this_ *PostgresqlDialect) ColumnsSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `SELECT t.*,tc.COMMENTS from all_tab_columns t `
	sql += "LEFT JOIN all_col_comments tc ON(tc.OWNER=t.OWNER AND tc.TABLE_NAME=t.TABLE_NAME AND tc.COLUMN_NAME=t.COLUMN_NAME)"
	sql += `WHERE 1=1 `
	if databaseName != "" {
		sql += `AND t.OWNER='` + databaseName + `' `
	}
	sql += `AND t.TABLE_NAME='` + tableName + `' `
	return
}
func (this_ *PostgresqlDialect) ColumnUpdateSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error) {

	var sqlList_ []string

	if column.OldName != "" && column.OldName != column.Name {
		sqlList_, err = this_.columnRenameSql(param, databaseName, tableName, column.OldName, column.Name)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	if column.Comment != column.OldComment {
		sqlList_, err = this_.ColumnCommentSql(param, databaseName, tableName, column.Name, column.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}
	return
}

func (this_ *PostgresqlDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
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
		primaryKey.TableSchema = data["TABLE_SCHEMA"].(string)
	}
	if data["TABLE_CATALOG"] != nil {
		primaryKey.TableCatalog = data["TABLE_CATALOG"].(string)
	}
	return
}
func (this_ *PostgresqlDialect) PrimaryKeysSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `SELECT cu.* FROM all_cons_columns cu, all_constraints au `
	sql += `WHERE cu.constraint_name = au.constraint_name and au.constraint_type = 'P' `
	if databaseName != "" {
		sql += `AND au.OWNER='` + databaseName + `' `
	}
	sql += `AND au.TABLE_NAME='` + tableName + `' `
	return
}

func (this_ *PostgresqlDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
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
	if GetStringValue(data["UNIQUENESS"]) == "UNIQUE" {
		index.Type = "unique"
	}
	if data["TABLE_NAME"] != nil {
		index.TableName = data["TABLE_NAME"].(string)
	}
	return
}
func (this_ *PostgresqlDialect) IndexesSelectSql(databaseName string, tableName string) (sql string, err error) {
	sql = `SELECT t.*,i.index_type,i.UNIQUENESS FROM all_ind_columns t,all_indexes i  `
	sql += `WHERE t.index_name = i.index_name `
	if databaseName != "" {
		sql += `AND t.TABLE_OWNER='` + databaseName + `' `
	}
	sql += `AND t.TABLE_NAME='` + tableName + `' `
	sql += `AND t.COLUMN_NAME NOT IN( `
	sql += `SELECT cu.COLUMN_NAME FROM all_cons_columns cu, all_constraints au `
	sql += `WHERE cu.constraint_name = au.constraint_name and au.constraint_type = 'P' `
	if databaseName != "" {
		sql += `AND au.OWNER='` + databaseName + `' `
	}
	sql += `AND au.TABLE_NAME='` + tableName + `' `

	sql += ") "
	return
}
