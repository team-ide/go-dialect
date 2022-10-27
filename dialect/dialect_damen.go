package dialect

func NewDaMenDialect() *DaMenDialect {

	dialect := NewOracleDialect()
	dialect.dialectType = DaMenType

	res := &DaMenDialect{
		OracleDialect: dialect,
	}
	res.init()
	return res
}

type DaMenDialect struct {
	*OracleDialect
}

func (this_ *DaMenDialect) init() {
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

	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "CLASS234882065", TypeFormat: "CLASS234882065($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "VARBINARY", TypeFormat: "VARBINARY($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BYTE", TypeFormat: "BYTE($l)", HasLength: true, IsString: true})
	this_.AddColumnTypeInfo(&ColumnTypeInfo{Name: "BINARY", TypeFormat: "BINARY($l)", HasLength: true, IsString: true})

	this_.AddFuncTypeInfo(&FuncTypeInfo{Name: "md5", Format: "md5"})
}
