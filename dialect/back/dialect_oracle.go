package back

import (
	"fmt"
	"strings"
)

func NewOracleDialect() *OracleDialect {

	res := &OracleDialect{
		DefaultDialect: NewDefaultDialect(OracleType),
	}
	res.init()
	return res
}

type OracleDialect struct {
	*DefaultDialect
}

func (this_ *OracleDialect) init() {
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
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "VARCHAR2(4000)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "CLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "CLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLOB", TypeFormat: "CLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "VARCHAR2(50)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "BLOB", HasLength: true, IsString: true})

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "SET", TypeFormat: "VARCHAR2(50)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "XMLTYPE", TypeFormat: "XMLTYPE($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "RAW", TypeFormat: "RAW($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "NVARCHAR2", TypeFormat: "NVARCHAR2($l)", HasLength: true, IsString: true})

	this_.AddFuncTypeInfo(&FuncTypeInfo{Name: "md5", Format: "md5"})
}

func (this_ *OracleDialect) OwnerModel(data map[string]interface{}) (owner *OwnerModel, err error) {
	if data == nil {
		return
	}
	owner = &OwnerModel{}
	if data["USERNAME"] != nil {
		owner.Name = data["USERNAME"].(string)
	}
	return
}
func (this_ *OracleDialect) OwnersSelectSql() (sql string, err error) {
	sql = `SELECT USERNAME FROM DBA_USERS ORDER BY USERNAME`
	return
}
func (this_ *OracleDialect) OwnerSelectSql(ownerName string) (sql string, err error) {
	sql = `SELECT USERNAME FROM DBA_USERS `
	sql += ` WHERE USERNAME='` + ownerName + `'`
	return
}
func (this_ *OracleDialect) OwnerChangeSql(ownerName string) (sql string, err error) {
	return
}
func (this_ *OracleDialect) OwnerCreateSql(owner *OwnerModel) (sqlList []string, err error) {
	sql := ``
	sql = `CREATE USER ` + owner.Name + ` IDENTIFIED BY "` + owner.Password + `" `
	sqlList = append(sqlList, sql)

	sql = `GRANT dba,resource,connect TO ` + owner.Name + ` `
	sqlList = append(sqlList, sql)
	return
}

func (this_ *OracleDialect) TableModel(data map[string]interface{}) (table *TableModel, err error) {
	if data == nil {
		return
	}
	table = &TableModel{}
	if data["TABLE_NAME"] != nil {
		table.Name = data["TABLE_NAME"].(string)
	}
	if data["OWNER"] != nil {
		table.OwnerName = data["OWNER"].(string)
	}
	return
}
func (this_ *OracleDialect) TablesSelectSql(ownerName string) (sql string, err error) {
	sql = `SELECT TABLE_NAME,OWNER FROM ALL_TABLES  `
	if ownerName != "" {
		sql += `WHERE OWNER ='` + ownerName + `' `
	}
	sql += `ORDER BY TABLE_NAME`
	return
}
func (this_ *OracleDialect) TableSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT TABLE_NAME,OWNER FROM ALL_TABLES `
	sql += `WHERE 1=1 `
	if ownerName != "" {
		sql += `AND OWNER='` + ownerName + `' `
	}
	sql += `AND TABLE_NAME='` + tableName + `' `
	sql += `ORDER BY TABLE_NAME`
	return
}

func (this_ *OracleDialect) ColumnModel(data map[string]interface{}) (column *ColumnModel, err error) {
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
func (this_ *OracleDialect) ColumnsSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT t.COLUMN_NAME,t.DATA_DEFAULT,t.TABLE_NAME,t.CHARACTER_SET_NAME,t.NULLABLE,t.DATA_TYPE,t.DATA_LENGTH,t.DATA_PRECISION,t.DATA_SCALE,tc.COMMENTS from ALL_TAB_COLUMNS t `
	sql += "LEFT JOIN ALL_COL_COMMENTS tc ON(tc.OWNER=t.OWNER AND tc.TABLE_NAME=t.TABLE_NAME AND tc.COLUMN_NAME=t.COLUMN_NAME)"
	sql += `WHERE 1=1 `
	if ownerName != "" {
		sql += `AND t.OWNER='` + ownerName + `' `
	}
	sql += `AND t.TABLE_NAME='` + tableName + `' `
	return
}

func (this_ *OracleDialect) PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error) {
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
	if data["OWNER"] != nil {
		primaryKey.OwnerName = data["OWNER"].(string)
	}
	return
}
func (this_ *OracleDialect) PrimaryKeysSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT cu.COLUMN_NAME,au.TABLE_NAME,au.OWNER FROM ALL_CONS_COLUMNS cu, ALL_CONSTRAINTS au `
	sql += `WHERE cu.CONSTRAINT_NAME = au.CONSTRAINT_NAME and au.CONSTRAINT_TYPE = 'P' `
	if ownerName != "" {
		sql += `AND au.OWNER='` + ownerName + `' `
	}
	sql += `AND au.TABLE_NAME='` + tableName + `' `
	return
}

func (this_ *OracleDialect) IndexModel(data map[string]interface{}) (index *IndexModel, err error) {
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
	if data["TABLE_OWNER"] != nil {
		index.OwnerName = data["TABLE_OWNER"].(string)
	}
	return
}
func (this_ *OracleDialect) IndexesSelectSql(ownerName string, tableName string) (sql string, err error) {
	sql = `SELECT t.INDEX_NAME,t.COLUMN_NAME,t.TABLE_OWNER,t.TABLE_NAME,i.INDEX_TYPE,i.UNIQUENESS FROM ALL_IND_COLUMNS t,ALL_INDEXES i  `
	sql += `WHERE t.INDEX_NAME = i.INDEX_NAME `
	if ownerName != "" {
		sql += `AND t.TABLE_OWNER='` + ownerName + `' `
	}
	sql += `AND t.TABLE_NAME='` + tableName + `' `
	sql += `AND t.COLUMN_NAME NOT IN( `
	sql += `SELECT cu.COLUMN_NAME FROM ALL_CONS_COLUMNS cu, ALL_CONSTRAINTS au `
	sql += `WHERE cu.CONSTRAINT_NAME = au.CONSTRAINT_NAME and au.CONSTRAINT_TYPE = 'P' `
	if ownerName != "" {
		sql += `AND au.OWNER='` + ownerName + `' `
	}
	sql += `AND au.TABLE_NAME='` + tableName + `' `

	sql += ") "
	return
}

func (this_ *OracleDialect) InsertDataListSql(ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}) (sqlList []string, batchSqlList []string, err error) {
	var batchSql = "INSERT ALL "
	var columnNames []string
	for _, one := range columnList {
		columnNames = append(columnNames, one.Name)
	}
	for index, data := range dataList {
		var columnList_ []string
		var values = "("
		for _, column := range columnList {
			str := this_.PackValueForSql(column, data[column.Name])
			if strings.EqualFold(str, "null") {
				continue
			}
			columnList_ = append(columnList_, column.Name)
			if column.Type == "TIMESTAMP" {
				str = `TO_TIMESTAMP(` + str + `,'yyyy-MM-dd HH24:mi:ss.ff6')`
			} else {
				if column.NotNull {
					if str == `''` {
						str = `' '`
					}
				}
				if len(str) > 1000 {
					key := fmt.Sprintf("%s_%s_%s_%d", ownerName, tableName, column.Name, index)
					batchSqlList = append(batchSqlList, `DECLARE `+key+` CLOB :=`+str)
					str = `:` + key
				}
			}
			values += str + ", "
		}
		values = strings.TrimSuffix(values, ", ")
		values += ")"

		insertSqlInfo := "INSERT INTO "
		if ownerName != "" {
			insertSqlInfo += this_.PackOwner(ownerName) + "."
		}
		insertSqlInfo += this_.PackTable(tableName)
		insertSqlInfo += " ("
		insertSqlInfo += this_.PackColumns(columnList_)
		insertSqlInfo += ") VALUES "

		sqlList = append(sqlList, insertSqlInfo+values)

		batchOne := "INTO "
		if ownerName != "" {
			batchOne += this_.PackOwner(ownerName) + "."
		}
		batchOne += this_.PackTable(tableName)
		batchOne += " ("
		batchOne += this_.PackColumns(columnList_)
		batchOne += ") VALUES " + values

		batchSql += "\n" + batchOne
	}
	batchSql += "\n" + `SELECT 1 FROM DUAL`
	batchSqlList = append(batchSqlList, batchSql)
	return
}
