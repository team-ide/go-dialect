package dialect

import "strings"

// Mysql 数据库 字段类型
var mysqlColumnTypeList = []*ColumnTypeInfo{
	{Name: `TINYINT`, Format: `TINYINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<3)||(columnPrecision>0&&columnPrecision<3))`, `INT1`, `BOOL`, `BOOLEAN`}, IsNumber: true, IsInteger: true, Comment: `1 Bytes 范围（有符号）(-128，127) 范围（无符号）(0，255) 小整数值`},
	{Name: `SMALLINT`, Format: `SMALLINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<6)||(columnPrecision>0&&columnPrecision<6))`, `INT2`}, IsNumber: true, IsInteger: true, Comment: `2 Bytes 范围（有符号）(-32 768，32 767) 范围（无符号）(0，65 535)  大整数值`},
	{Name: `MEDIUMINT`, Format: `MEDIUMINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<9)||(columnPrecision>0&&columnPrecision<9))`}, IsNumber: true, IsInteger: true, Comment: `3 Bytes 范围（有符号）(-8 388 608，8 388 607) 范围（无符号）(0，16 777 215)  大整数值`},
	{Name: `INT`, Format: `INT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<11)||(columnPrecision>0&&columnPrecision<11))`, `INT4`}, IsNumber: true, IsInteger: true, Comment: `4 Bytes 范围（有符号）(-2 147 483 648，2 147 483 647) 范围（无符号）(0，4 294 967 295)  大整数值`},
	{Name: `INTEGER`, Format: `INTEGER($l)`, IsNumber: true, IsInteger: true, Comment: `同上`},
	{Name: `BIGINT`, Format: `BIGINT($l)`, Matches: []string{`NUMBER&&columnScale==0`, `INT8`}, IsNumber: true, IsInteger: true, Comment: `8 Bytes 范围（有符号）(-9,223,372,036,854,775,808，9 223 372 036 854 775 807) 范围（无符号）(0，18 446 744 073 709 551 615)  极大整数值`},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`FLOAT4`}, IsNumber: true, IsFloat: true, Comment: `4 Bytes 范围（有符号）(-3.402 823 466 E+38，-1.175 494 351 E-38)，0，(1.175 494 351 E-38，3.402 823 466 351 E+38) 范围（无符号）0，(1.175 494 351 E-38，3.402 823 466 E+38)  单精度 浮点数值`},
	{Name: `DOUBLE`, Format: `DOUBLE`, Matches: []string{`FLOAT8`, `DOUBLE PRECISION`}, IsNumber: true, IsFloat: true, Comment: `8 Bytes 范围（有符号）(-1.797 693 134 862 315 7 E+308，-2.225 073 858 507 201 4 E-308)，0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308)范围（无符号）0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308)  双精度 浮点数值`},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, Matches: []string{`NUMBER`, `REAL`, `NUMERIC`}, IsNumber: true, IsFloat: true, Comment: `对DECIMAL(M,D) ，如果M>D，为M+2 否则为D+2 小数值`},
	{Name: `DEC`, Format: `DEC($p, $s)`, IsNumber: true, IsFloat: true, Comment: `同上`},
	{Name: `BIT`, Format: `BIT($p)`, Comment: `位字段类型。M 表示每个值的位数，范围为 1～64。如果 M 被省略，默认值为 1。如果为 BIT(M) 列分配的值的长度小于 M 位，在值的左边用 0 填充。例如，为 BIT(6) 列分配一个值 b'101'，其效果与分配 b'000101' 相同`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`NCHAR`, `CHARACTER`}, IsString: true, Comment: `0-255 bytes 定长字符串 固定长度非二进制字符串 M 字节，1<=M<=255`},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`VARCHAR2`, `NVARCHAR2`, `BPCHAR`}, IsString: true, Comment: `0-65535 bytes	变长字符串 变长非二进制字符串 L+1字节，在此，L< = M和 1<=M<=255`},
	{Name: `TINYTEXT`, Format: `TINYTEXT`, IsString: true, Comment: `0-255 bytes 短文本字符串`},
	{Name: `TEXT`, Format: `TEXT`, Matches: []string{`CLOB`, `NCLOB`, `ROWID`, `UROWID`}, IsString: true, Comment: `0-65 535 bytes 长文本数据`},
	{Name: `MEDIUMTEXT`, Format: `MEDIUMTEXT`, Matches: []string{`LONGVARCHAR`}, IsString: true, Comment: `0-16 777 215 bytes	中等长度文本数据`},
	{Name: `LONGTEXT`, Format: `LONGTEXT`, IsString: true, Comment: `0-4 294 967 295 bytes	极大文本数据`},
	{Name: `BINARY`, Format: `BINARY($l)`, IsBytes: true, Comment: `0-255 bytes 不超过 255 个字符的二进制字符串 固定长度二进制字符串 M 字节`},
	{Name: `VARBINARY`, Format: `VARBINARY($l)`, IsBytes: true, Comment: `可变长度二进制字符串 M+1 字节`},
	{Name: `TINYBLOB`, Format: `TINYBLOB`, IsBytes: true, Comment: `非常小的BLOB L+1 字节，在此，L<2^8`},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BFILE`, `RAW`, `BYTEA`}, IsBytes: true, Comment: `0-65 535 bytes 二进制形式的长文本数据 小 BLOB L+2 字节，在此，L<2^16`},
	{Name: `MEDIUMBLOB`, Format: `MEDIUMBLOB`, IsBytes: true, Comment: `0-16 777 215 bytes	二进制形式的中等长度文本数据 中等大小的BLOB	L+3 字节，在此，L<2^24`},
	{Name: `LONGBLOB`, Format: `LONGBLOB`, Matches: []string{`LONG RAW`}, IsBytes: true, Comment: `0-4 294 967 295 bytes	二进制形式的极大文本数据 非常大的BLOB	L+4 字节，在此，L<2^32`},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true, Comment: `3 bytes '-838:59:59'/'838:59:59' HH:MM:SS 时间值或持续时间`},
	{Name: `TIME`, Format: `TIME`, IsDateTime: true, Comment: `1 bytes 1901/2155 YYYY 年份值`},
	{Name: `YEAR`, Format: `YEAR`, IsDateTime: true, Comment: `8 bytes '1000-01-01 00:00:00' 到 '9999-12-31 23:59:59' YYYY-MM-DD hh:mm:ss 混合日期和时间值`},
	{Name: `DATETIME`, Format: `DATETIME`, Matches: []string{`DATETIME WITH TIME ZONE`}, IsDateTime: true, Comment: `4 bytes '1970-01-01 00:00:01' UTC 到 '2038-01-19 03:14:07' UTC 结束时间是第 2147483647 秒，北京时间 2038-1-19 11:14:07，格林尼治时间 2038年1月19日 凌晨 03:14:07 YYYY-MM-DD hh:mm:ss 混合日期和时间值，时间戳`,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			if strings.Contains(strings.ToLower(column.ColumnExtra), "on update current_timestamp") {
				columnDefaultPack += " ON UPDATE CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`TIMESTAMP WITH TIME ZONE`, `TIMESTAMP WITH LOCAL TIME ZONE`, `INTERVAL YEAR TO MONTH`, `INTERVAL DAY TO SECOND`, `TIME WITH TIME ZONE`, `TIMESTAMP WITHOUT TIME ZONE`}, IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			if strings.Contains(strings.ToLower(column.ColumnExtra), "on update current_timestamp") {
				columnDefaultPack += " ON UPDATE CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `ENUM`, Format: `ENUM`, IsEnum: true, Comment: `枚举类型，只能有一个枚举字符串值 1或2个字节，取决于枚举值的数目 (最大值为65535)`,
		FullColumnByColumnType: func(columnType string, column *ColumnModel) (err error) {
			if strings.Contains(columnType, "(") {
				setStr := columnType[strings.Index(columnType, "(")+1 : strings.Index(columnType, ")")]
				setStr = strings.ReplaceAll(setStr, "'", "")
				column.ColumnEnums = strings.Split(setStr, ",")
			}
			return
		},
	},
	{Name: `SET`, Format: `SET`, IsEnum: true, Comment: `一个设置，字符串对象可以有零个或 多个SET成员 1、2、3、4或8个字节，取决于集合 成员的数量（最多64个成员）`,
		FullColumnByColumnType: func(columnType string, column *ColumnModel) (err error) {
			if strings.Contains(columnType, "(") {
				setStr := columnType[strings.Index(columnType, "(")+1 : strings.Index(columnType, ")")]
				setStr = strings.ReplaceAll(setStr, "'", "")
				column.ColumnEnums = strings.Split(setStr, ",")
			}
			return
		},
	},
}

// Oracle 数据库 字段类型
var oracleColumnTypeList = []*ColumnTypeInfo{
	{Name: `INTEGER`, Format: `INTEGER`, Matches: []string{`TINYINT`, `SMALLINT`, `MEDIUMINT`, `INT`, `BIGINT`, `BIT&&columnLength==1||columnPrecision==1`, `INT1`, `INT2`, `INT4`, `BOOL`, `BOOLEAN`}, IsNumber: true, IsInteger: true, Comment: `INTEGER是NUMBER的子类型，它等同于NUMBER（38,0），用来存储整数。若插入、更新的数值有小数，则会被四舍五入`},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`DOUBLE`, `FLOAT4`, `FLOAT8`, `DOUBLE PRECISION`}, IsNumber: true, IsFloat: true, Comment: `FLOAT类型也是NUMBER的子类型。

Float(n)，数n指示位的精度，可以存储的值的数目。n值的范围可以从 1 到 126。若要从二进制转换为十进制的精度，请将n乘以 0.30103。要从十进制转换为二进制的精度，请用3.32193乘小数精度。126位二进制精度的最大值是大约相当于38位小数精度`},
	{Name: `NUMBER`, Format: `NUMBER($p, $s)`, Matches: []string{`DECIMAL`, `DEC`, `REAL`, `NUMERIC`, `INT8`}, IsNumber: true, Comment: `NUMBER(P,S)是最常见的数字类型，可以存放数据范围为10^130~10^126（不包含此值)，需要1~22字节(BYTE)不等的存储空间。

P 是Precison的英文缩写，即精度缩写，表示有效数字的位数，最多不能超过38个有效数字。

S是Scale的英文缩写，可以使用的范围为-84~127。Scale为正数时，表示从小数点到最低有效数字的位数，它为负数时，表示从最大有效数字到小数点的位数。`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`CHARACTER`}, IsString: true, Comment: `定长字符串，会用空格填充来达到其最大长度。非NULL的CHAR（12）总是包含12字节信息。CHAR字段最多可以存储2000字节的信息。如果创建表时，不指定CHAR长度，则默认为1。另外你可以指定它存储字节或字符，例如 CHAR(12 BYTYE)、CHAR(12 CHAR)。一般来说默认是存储字节`},
	{Name: `NCHAR`, Format: `NCHAR($l)`, IsString: true, Comment: `一个包含UNICODE格式数据的定长字符串。NCHAR字段最多可以存储2000字节的信息，它的最大长度取决于国家字符集`},
	{Name: `VARCHAR2`, Format: `VARCHAR2($l)`, Matches: []string{`VARCHAR`, `BPCHAR`}, IsString: true, Comment: `变长字符串，与CHAR类型不同，它不会使用空格填充至最大长度。VARCHAR2最多可以存储4000字节的信息`},
	{Name: `NVARCHAR2`, Format: `NVARCHAR2($l)`, IsString: true, Comment: `一个包含UNICODE格式数据的变长字符串，NVARCHAR2最多可以存储4000字节的信息`},
	{Name: `CLOB`, Format: `CLOB`, Matches: []string{`TINYTEXT`, `TEXT`, `MEDIUMTEXT`, `LONGTEXT`, `ENUM`, `SET`, `LONGVARCHAR`}, IsString: true, Comment: `CLOB存储单字节和多字节字符数据。支持固定宽度和可变宽度的字符集。CLOB对象可以存储最多 (4 gigabytes-1) * (database block size) 大小的字符`},
	{Name: `NCLOB`, Format: `NCLOB`, Matches: []string{`BINARY`, `VARBINARY`, `TINYBLOB`, `MEDIUMBLOB`, `LONGBLOB`}, IsString: true, Comment: `NCLOB存储UNICODE类型的数据，支持固定宽度和可变宽度的字符集，NCLOB对象可以存储最多(4 gigabytes-1) * (database block size)大小的文本数据`},
	{Name: `RAW`, Format: `RAW($l)`, IsString: true, Comment: `用于存储二进制或字符类型数据，变长二进制数据类型，这说明采用这种数据类型存储的数据不会发生字符集转换。这种类型最多可以存储2000字节的信息，建议使用 BLOB 来代替它`},
	{Name: `ROWID`, Format: `ROWID`, IsString: true, Comment: `ROWID是一种特殊的列类型，称之为伪列（pseudocolumn）。ROWID伪列在SQL SELECT语句中可以像普通列那样被访问。ROWID表示行的地址，ROWID伪列用ROWID数据类型定义。Oracle数据库中每行都有一个伪列。

ROWID与磁盘驱动的特定位置有关，因此，ROWID是获得行的最快方法。但是，行的ROWID会随着卸载和重载数据库而发生变化，因此建议不要在事务中使用ROWID伪列的值。例如，一旦当前应用已经使用完记录，就没有理由保存行的ROWID。不能通过任何SQL语句来设置标准的ROWID伪列的值。

列或变量可以定义成ROWID数据类型，但是Oracle不能保证该列或变量的值是一个有效的ROWID`},
	{Name: `UROWID`, Format: `UROWID`, IsString: true, Comment: `UROWID，它用于表，是行主键的一个表示，基于主键生成。UROWID与ROWID的区别就是UROWID可以表示各种ROWID，使用较安全。一般是索引组织表在使用UROWID`},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BIT&&columnLength>1||columnPrecision>1`, `BYTEA`}, IsBytes: true, Comment: `BLOB存储非结构化的二进制数据大对象，它可以被认为是没有字符集语义的比特流，一般是图像、声音、视频等文件。BLOB对象最多存储(4 gigabytes-1) * (database block size)的二进制数据`},
	{Name: `BFILE`, Format: `BFILE`, IsBytes: true, Comment: `二进制文件，存储在数据库外的系统文件，只读的，数据库会将该文件当二进制文件处理`},
	{Name: `LONG RAW`, Format: `LONG RAW`, IsBytes: true, Comment: `LONG RAW类型，能存储2GB的原始二进制数据（不用进行字符集转换的数据）。建议使用BLOB来代替它`},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true, Comment: `DATE是最常用的数据类型，日期数据类型存储日期和时间信息。虽然可以用字符或数字类型表示日期和时间信息，但是日期数据类型具有特殊关联的属性。为每个日期值，Oracle 存储以下信息： 世纪、 年、 月、 日期、 小时、 分钟和秒。一般占用7个字节的存储空间`},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`TIME`, `YEAR`, `DATETIME`}, IsDateTime: true, Comment: `TIMESTAMP是一个7字节或12字节的定宽日期/时间数据类型，是DATE类型的扩展类型。它与DATE数据类型不同，因为TIMESTAMP可以包含小数秒，带小数秒的TIMESTAMP在小数点右边最多可以保留9位`,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `TIMESTAMP WITH TIME ZONE`, Format: `TIMESTAMP WITH TIME ZONE`, Matches: []string{`DATETIME WITH TIME ZONE`, `TIME WITH TIME ZONE`, `TIMESTAMP WITHOUT TIME ZONE`}, IsDateTime: true, Comment: `和TIMESTAMP一样，只不过可以在设置时候指定时区`},
	{Name: `TIMESTAMP WITH LOCAL TIME ZONE`, Format: `TIMESTAMP WITH LOCAL TIME ZONE`, IsDateTime: true},
}

// 达梦 数据库 字段类型
var dmColumnTypeList = []*ColumnTypeInfo{
	{Name: `BIT`, Format: `BIT`, IsNumber: true, IsInteger: true, Comment: `BIT 类型用于存储整数数据 1、0 或 NULL`},
	{Name: `INTEGER`, Format: `INTEGER`, Matches: []string{`MEDIUMINT`, `INT4`}, IsNumber: true, IsInteger: true, Comment: `用于存储有符号整数，精度为 10，标度为 0。取值范围为：-2147483648(-2^31)～+2147483647(2^31-1)`},
	{Name: `INT`, Format: `INT`, IsNumber: true, IsInteger: true, Comment: `同上`},
	{Name: `PLS_INTEGER`, Format: `PLS_INTEGER`, IsNumber: true, IsInteger: true, Comment: `同上`},
	{Name: `BIGINT`, Format: `BIGINT`, Matches: []string{`INT8`}, IsNumber: true, IsInteger: true, Comment: `用于存储有符号整数，精度为 19，标度为 0。取值范围为：-9223372036854775808(-2^63)～+9223372036854775807(2^63-1)`},
	{Name: `TINYINT`, Format: `TINYINT`, Matches: []string{`INT1`, `BOOL`, `BOOLEAN`}, IsNumber: true, IsInteger: true, Comment: `储有符号整数，精度为 3，标度为 0。取值范围为：-128～+127`},
	{Name: `SMALLINT`, Format: `SMALLINT`, Matches: []string{`INT2`}, IsNumber: true, IsInteger: true, Comment: `用于存储有符号整数，精度为 5，标度为 0。取值范围为：-32768(-2^15)~ +32767(2^15-1)`},
	{Name: `REAL`, Format: `REAL`, IsNumber: true, IsFloat: true, Comment: `REAL 是带二进制的浮点数，但它不能由用户指定使用的精度，系统指定其二进制精度为 24，十进制精度为 7。取值范围-3.4E+38～3.4E + 38`},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`FLOAT4`}, IsNumber: true, IsFloat: true, Comment: `FLOAT 是带二进制精度的浮点数，精度最大不超过 53，如省略精度，则二进制精度为 53，十进制精度为 15。取值范围为-1.7E+308～1.7E+308`},
	{Name: `DOUBLE`, Format: `DOUBLE`, Matches: []string{`FLOAT8`}, IsNumber: true, IsFloat: true, Comment: `同 FLOAT 相似，精度最大不超过 53`},
	{Name: `DOUBLE PRECISION`, Format: `DOUBLE PRECISION`, IsNumber: true, IsFloat: true, Comment: `该类型指明双精度浮点数，其二进制精度为 53，十进制精度为 15。取值范围-1.7E+308 ～1.7E+308`},
	{Name: `NUMERIC`, Format: `NUMERIC($p, $s)`, IsNumber: true, Comment: `精度是一个无符号整数，定义了总的数字数，精度范围是 1~38 ，标度定义了小数点右边的数字位数，一个数的标度不应大于其精度`},
	{Name: `NUMBER`, Format: `NUMBER($p, $s)`, IsNumber: true, Comment: `同上`},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, Matches: []string{`REAL`}, IsNumber: true, Comment: `同上`},
	{Name: `DEC`, Format: `DEC($p, $s)`, IsNumber: true, Comment: `同上`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`NCHAR`}, IsString: true, Comment: `定长字符串，最大长度由数据库页面大小决定，具体可参考《DM8_SQL 语言使用手册》1.4.1 节。长度不足时，自动填充空格`},
	{Name: `CHARACTER`, Format: `CHARACTER($l)`, IsString: true, Comment: `同上`},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`VARCHAR2`, `NVARCHAR2`, `BPCHAR`}, IsString: true, Comment: `可变长字符串，最大长度由数据库页面大小决定`},
	{Name: `TEXT`, Format: `TEXT`, IsString: true, Comment: `变长字符串类型，其字符串的长度最大为 100G-1，可用于存储长的文本串`},
	{Name: `CLOB`, Format: `CLOB`, Matches: []string{`TINYTEXT`, `MEDIUMTEXT`, `LONGTEXT`, `SET`, `ENUM`, `NCLOB`, `ROWID`, `UROWID`}, IsString: true, Comment: `CLOB 类型用于指明变长的字符串，长度最大为 100G-1 字节`},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BINARY`, `VARBINARY`, `TINYBLOB`, `MEDIUMBLOB`, `LONGBLOB`, `RAW`, `LONG RAW`, `BYTEA`}, IsBytes: true, Comment: `BLOB 类型用于指明变长的二进制大对象，长度最大为 100G-1 字节`},
	{Name: `BFILE`, Format: `BFILE`, IsBytes: true, Comment: `BFILE 用于指明存储在操作系统中的二进制文件，文件存储在操作系统而非数据库中，仅能进行只读访问`},
	{Name: `BINARY`, Format: `BINARY($l)`, IsBytes: true, Comment: `BINARY 数据类型指定定长二进制数据。缺省长度为 1 个字节`},
	{Name: `VARBINARY`, Format: `VARBINARY($l)`, IsBytes: true, Comment: `VARBINARY 数据类型指定变长二进制数据，用法类似 BINARY 数据类型，可以指定一个正整数作为数据长度。缺省长度为 8188 个字节，最大长度由数据库页面大小决定`},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true, Comment: `DATE 类型包括年、月、日信息，定义了'-4712-01-01'和'9999-12-31'之间任何一个有效的格里高利日期`},
	{Name: `TIME`, Format: `TIME`, IsDateTime: true, Comment: `IME 类型包括时、分、秒信息，定义了一个在'00:00:00.000000'和'23:59:59.999999'之间的有效时间。TIME 类型的小数秒精度规定了秒字段中小数点后面的位数，取值范围为 0～6，如果未定义，缺省精度为 0`},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`YEAR`}, IsDateTime: true, Comment: `TIMESTAMP/DATETIME 类型包括年、月、日、时、分、秒信息，定义了一个在'-4712-01-0100:00:00.000000'和'9999-12-31 23:59:59.999999'之间的有效格里高利日期时间。小数秒精度规定了秒字段中小数点后面的位数，取值范围为 0～6，如果未定义，缺省精度为 6`,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `DATETIME`, Format: `DATETIME`, IsDateTime: true, Comment: `同上`,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `DATETIME WITH TIME ZONE`, Format: `DATETIME WITH TIME ZONE`, IsDateTime: true},
	{Name: `TIME WITH TIME ZONE`, Format: `TIME WITH TIME ZONE`, Matches: []string{`TIMESTAMP WITHOUT TIME ZONE`}, IsDateTime: true, Comment: `描述一个带时区的 TIME 值，其定义是在 TIME 类型的后面加上时区信息。时区部分的实质是 INTERVAL HOUR TO MINUTE 类型，取值范围：-12:59 与 +14:00 之间。例如：TIME '09:10:21 +8:00'`},
	{Name: `TIMESTAMP WITH TIME ZONE`, Format: `TIMESTAMP WITH TIME ZONE`, IsDateTime: true, Comment: `描述一个带时区的 TIMESTAMP 值，其定义是在 TIMESTAMP 类型的后面加上时区信息。时区部分的实质是 INTERVAL HOUR TO MINUTE 类型，取值范围：-12:59 与 +14:00 之间。例如：’2009-10-11 19:03:05.0000 -02:10’`},
	{Name: `TIMESTAMP WITH LOCAL TIME ZONE`, Format: `TIMESTAMP WITH LOCAL TIME ZONE`, IsDateTime: true, Comment: `描述一个本地时区的 TIMESTAMP 值，能够将标准时区类型 TIMESTAMP WITH TIME ZONE 类型转化为本地时区类型，如果插入的值没有指定时区，则默认为本地时区。`},
}

// 金仓 数据库 字段类型
var kingBaseColumnTypeList = []*ColumnTypeInfo{
	{Name: `TINYINT`, Format: `TINYINT`, Matches: []string{`BIT&&columnLength==1||columnPrecision==1`}, IsNumber: true, IsInteger: true, Comment: `有符号整数，取值范围 -128 ~ +127`},
	{Name: `SMALLINT`, Format: `SMALLINT`, IsNumber: true, IsInteger: true, Comment: `有符号整数，取值范围 -32768 ~ +32767`},
	{Name: `INTEGER`, Format: `INTEGER`, Matches: []string{`MEDIUMINT`}, IsNumber: true, IsInteger: true, Comment: `有符号整数，取值范围 -2147483648~ +2147483647`},
	{Name: `INT`, Format: `INT`, Matches: []string{`INT1`, `INT2`, `INT4`}, IsNumber: true, IsInteger: true, Comment: `同上`},
	{Name: `BIGINT`, Format: `BIGINT`, Matches: []string{`INT8`}, IsNumber: true, IsInteger: true, Comment: `有符号整数，取值范围 -9223372036854775808~ +9223372036854775807`},
	{Name: `SMALLSERIAL`, Format: `SMALLSERIAL`, IsNumber: true, IsInteger: true, Comment: `相当于创建一个SMALLINT列`},
	{Name: `SERIAL`, Format: `SERIAL`, IsNumber: true, IsInteger: true, Comment: `相当于创建一个INT列`},
	{Name: `BIGSERIAL`, Format: `BIGSERIAL`, IsNumber: true, IsInteger: true, Comment: `相当于创建一个BIGINT列`},
	{Name: `REAL`, Format: `REAL`, IsNumber: true, IsFloat: true, Comment: `范围在 -1E+37 到 +1E+37 之间，精度至少是 6 位小数`},
	{Name: `DOUBLE PRECISION`, Format: `DOUBLE PRECISION`, IsNumber: true, IsFloat: true, Comment: `范围在 -1E+37 到 +1E+37 之间，精度至少是15位小数`},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`DOUBLE`, `FLOAT4`, `FLOAT8`}, IsNumber: true, IsFloat: true, Comment: `当p取值为1-24时，与REAL相同。当p取值为25-53时，与DOUBLE PRECISION相同。 没有指定精度时，与DOUBLE PRECISION相同`},
	{Name: `NUMERIC`, Format: `NUMERIC($p, $s)`, Matches: []string{`REAL`}, IsNumber: true, Comment: `存储0 以及绝对值为[1.0 x 10-130, 1.0 x 10126)的正、负定点数。 在算术运算中，如果超出范围，KingbaseE报错。

precision表示精度，是整个数中有效位的总数，也就是小数点两边的位数。取值范围为 1~1000。 scale表示标度，是小数部分的数字位数，也就是小数点右边的部分。取值范围为0~1000。

使用该数据类型时，最好指定定点数的小数位数和精度，以便在输入时进行额外的完整性检查。 指定小数位数和精度不会强制所有值都达到固定长度。如果某个值超过精度，KingbaseES将返回错误。如果某个值超过标度，KingbaseES会对其进行四舍五入。

也可以使用NUMERIC(precision) 类型，即标度为0的定点数，即NUMERIC(precision, 0)

也可以直接使用NUMERIC类型，缺省精度和标度，指定KingbaseES数值的最大精度和标度。 考虑到移植性，在使用时最好是显式声明精度和标度。
————————————————
版权声明：本文为CSDN博主「沉舟侧畔千帆过_」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/arthemis_14/article/details/125843469`},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, Matches: []string{`DEC`}, IsNumber: true, Comment: `同上`},
	{Name: `NUMBER`, Format: `NUMBER($p, $s)`, IsNumber: true, Comment: `同上`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`NCHAR`, `CHARACTER`}, IsString: true},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`VARCHAR2`, `NVARCHAR2`, `BPCHAR`}, IsString: true},
	{Name: `CHARACTER`, Format: `CHARACTER($l)`, Matches: []string{`CHARACTER VARYING`}, IsString: true},
	{Name: `CLOB`, Format: `CLOB`, Matches: []string{`TINYTEXT`, `MEDIUMTEXT`, `LONGTEXT`, `SET`, `ENUM`, `NCLOB`, `ROWID`, `UROWID`, `LONGVARCHAR`}, IsString: true},
	{Name: `TEXT`, Format: `TEXT`, IsString: true},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BINARY`, `VARBINARY`, `TINYBLOB`, `MEDIUMBLOB`, `LONGBLOB`, `BFILE`, `RAW`, `LONG RAW`, `BIT&&columnLength>1||columnPrecision>1`}, IsBytes: true},
	{Name: `BYTEA`, Format: `BYTEA`, IsBytes: true},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`YEAR`, `DATETIME`, `TIME`}, IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `TIMESTAMP WITHOUT TIME ZONE`, Format: `TIMESTAMP WITHOUT TIME ZONE`, Matches: []string{`TIMESTAMP WITH TIME ZONE`, `TIMESTAMP WITH LOCAL TIME ZONE`, `DATETIME WITH TIME ZONE`, `TIME WITH TIME ZONE`}, IsDateTime: true},
	{Name: `BOOL`, Format: `BOOL`, IsBoolean: true, Comment: `布尔数据类型：TRUE 和 FALSE。DMSQL 程序的布尔类型和 INT 类型可以相互转化。如果变量或方法返回的类型是布尔类型，则返回值为 0 或 1。TRUE 和非 0 值的返回值为 1，FALSE 和 0 值返回为 0`},
	{Name: `BOOLEAN`, Format: `BOOLEAN`, IsBoolean: true, Comment: `同上`},
}

// 神通 数据库 字段类型
var shenTongColumnTypeList = []*ColumnTypeInfo{
	{Name: `TINYINT`, Format: `TINYINT`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<3)||(columnPrecision>0&&columnPrecision<3))`, `BIT&&columnLength==1||columnPrecision==1`}, IsNumber: true, IsInteger: true},
	{Name: `INT`, Format: `INT`, Matches: []string{`SMALLINT`, `MEDIUMINT`, `NUMBER&&columnScale==0&&((columnLength>0&&columnLength<11)||(columnPrecision>0&&columnPrecision<11))`}, IsNumber: true, IsInteger: true},
	{Name: `INTEGER`, Format: `INTEGER`, IsNumber: true, IsInteger: true},
	{Name: `INT1`, Format: `INT1`, IsNumber: true, IsInteger: true},
	{Name: `INT2`, Format: `INT2`, IsNumber: true, IsInteger: true},
	{Name: `INT4`, Format: `INT4`, IsNumber: true, IsInteger: true},
	{Name: `INT8`, Format: `INT8`, Matches: []string{`BIGINT`, `NUMBER&&columnScale==0`}, IsNumber: true, IsInteger: true},
	{Name: `FLOAT4`, Format: `FLOAT4`, IsNumber: true, IsFloat: true},
	{Name: `FLOAT8`, Format: `FLOAT8`, Matches: []string{`DOUBLE`, `FLOAT`, `DOUBLE PRECISION`}, IsNumber: true, IsFloat: true},
	{Name: `NUMERIC`, Format: `NUMERIC($p, $s)`, Matches: []string{`REAL`}, IsNumber: true},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, Matches: []string{`DEC`, `NUMBER`}, IsNumber: true},
	{Name: `SERIAL`, Format: `SERIAL`, IsNumber: true},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`NCHAR`, `CHARACTER`}, IsString: true},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`VARCHAR2`, `NVARCHAR2`}, IsString: true},
	{Name: `BPCHAR`, Format: `BPCHAR($l)`, IsString: true},
	{Name: `CLOB`, Format: `CLOB`, Matches: []string{`TINYTEXT`, `MEDIUMTEXT`, `LONGTEXT`, `ENUM`, `SET`, `NCLOB`, `ROWID`, `UROWID`, `LONGVARCHAR`}, IsString: true},
	{Name: `TEXT`, Format: `TEXT`, IsString: true},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BINARY`, `VARBINARY`, `TINYBLOB`, `MEDIUMBLOB`, `LONGBLOB`, `RAW`, `LONG RAW`, `BIT&&columnLength>1||columnPrecision>1`, `BYTEA`}, IsBytes: true},
	{Name: `BFILE`, Format: `BFILE`, IsBytes: true},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true},
	{Name: `TIME`, Format: `TIME`, IsDateTime: true},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`YEAR`, `DATETIME`, `TIMESTAMP WITH TIME ZONE`, `TIMESTAMP WITH LOCAL TIME ZONE`, `INTERVAL YEAR TO MONTH`, `INTERVAL DAY TO SECOND`, `DATETIME WITH TIME ZONE`, `TIME WITH TIME ZONE`, `TIMESTAMP WITHOUT TIME ZONE`}, IsDateTime: true},
	{Name: `BOOL`, Format: `BOOL`, IsBoolean: true, Comment: `布尔数据类型：TRUE 和 FALSE。DMSQL 程序的布尔类型和 INT 类型可以相互转化。如果变量或方法返回的类型是布尔类型，则返回值为 0 或 1。TRUE 和非 0 值的返回值为 1，FALSE 和 0 值返回为 0`},
	{Name: `BOOLEAN`, Format: `BOOLEAN`, IsBoolean: true, Comment: `同上`},
}

// Sqlite 数据库 字段类型
var sqliteColumnTypeList = []*ColumnTypeInfo{
	{Name: `TINYINT`, Format: `TINYINT($l)`, Matches: []string{`INT1`}, IsNumber: true, IsInteger: true},
	{Name: `SMALLINT`, Format: `SMALLINT($l)`, Matches: []string{`INT2`}, IsNumber: true, IsInteger: true},
	{Name: `MEDIUMINT`, Format: `MEDIUMINT($l)`, IsNumber: true, IsInteger: true},
	{Name: `INT`, Format: `INT($l)`, Matches: []string{`INT4`}, IsNumber: true, IsInteger: true},
	{Name: `INTEGER`, Format: `INTEGER($l)`, IsNumber: true, IsInteger: true},
	{Name: `BIGINT`, Format: `BIGINT($l)`, Matches: []string{`INT8`}, IsNumber: true, IsInteger: true},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`FLOAT4`}, IsNumber: true, IsFloat: true},
	{Name: `DOUBLE`, Format: `DOUBLE`, Matches: []string{`FLOAT8`, `DOUBLE PRECISION`}, IsNumber: true, IsFloat: true},
	{Name: `REAL`, Format: `REAL`, IsNumber: true, IsFloat: true},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, IsNumber: true, IsFloat: true},
	{Name: `DEC`, Format: `DEC($p, $s)`, IsNumber: true, IsFloat: true},
	{Name: `NUMBER`, Format: `NUMBER($p, $s)`, IsNumber: true, IsFloat: true},
	{Name: `NUMERIC`, Format: `NUMERIC($p, $s)`, IsNumber: true, IsFloat: true, Comment: `当文本数据被插入到亲缘性为NUMERIC的字段中时，如果转换操作不会导致数据信息丢失以及完全可逆，那么SQLite就会将该文本数据转换为INTEGER或REAL类型的数据，如果转换失败，SQLite仍会以TEXT方式存储该数据。对于NULL或BLOB类型的新数据，SQLite将不做任何转换，直接以NULL或BLOB的方式存储该数据。需要额外说明的是，对于浮点格式的常量文本，如"30000.0"，如果该值可以转换为INTEGER同时又不会丢失数值信息，那么SQLite就会将其转换为INTEGER的存储方式。`},
	{Name: `BIT`, Format: `BIT($l)`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`CHARACTER`}, IsString: true},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`SET`, `BPCHAR`}, IsString: true},
	{Name: `TINYTEXT`, Format: `TINYTEXT`, IsString: true},
	{Name: `TEXT`, Format: `TEXT`, IsString: true, Comment: `值是一个文本字符串，使用数据库编码（UTF-8、UTF-16BE 或 UTF-16LE）存储。`},
	{Name: `MEDIUMTEXT`, Format: `MEDIUMTEXT`, IsString: true},
	{Name: `LONGTEXT`, Format: `LONGTEXT`, Matches: []string{`LONGVARCHAR`}, IsString: true},
	{Name: `ENUM`, Format: `ENUM`, IsString: true},
	{Name: `NCHAR`, Format: `NCHAR($l)`, IsString: true},
	{Name: `VARCHAR2`, Format: `VARCHAR2($l)`, IsString: true},
	{Name: `NVARCHAR2`, Format: `NVARCHAR2($l)`, IsString: true},
	{Name: `CLOB`, Format: `CLOB($l)`, IsString: true},
	{Name: `NCLOB`, Format: `NCLOB($l)`, IsString: true},
	{Name: `RAW`, Format: `RAW($l)`, IsString: true},
	{Name: `ROWID`, Format: `ROWID`, IsString: true},
	{Name: `UROWID`, Format: `UROWID`, IsString: true},
	{Name: `BINARY`, Format: `BINARY($l)`, IsBytes: true},
	{Name: `VARBINARY`, Format: `VARBINARY($l)`, IsBytes: true},
	{Name: `TINYBLOB`, Format: `TINYBLOB`, IsBytes: true},
	{Name: `BLOB`, Format: `BLOB`, IsBytes: true},
	{Name: `MEDIUMBLOB`, Format: `MEDIUMBLOB`, IsBytes: true},
	{Name: `LONGBLOB`, Format: `LONGBLOB`, IsBytes: true},
	{Name: `BFILE`, Format: `BFILE`, IsBytes: true},
	{Name: `LONG RAW`, Format: `LONG RAW`, IsBytes: true},
	{Name: `BYTEA`, Format: `BYTEA`, IsBytes: true},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true},
	{Name: `TIME`, Format: `TIME`, IsDateTime: true},
	{Name: `YEAR`, Format: `YEAR`, IsDateTime: true},
	{Name: `DATETIME`, Format: `DATETIME`, IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`INTERVAL YEAR TO MONTH`, `INTERVAL DAY TO SECOND`}, IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `TIMESTAMP WITH TIME ZONE`, Format: `TIMESTAMP WITH TIME ZONE`, Matches: []string{`TIMESTAMP WITHOUT TIME ZONE`}, IsDateTime: true},
	{Name: `TIMESTAMP WITH LOCAL TIME ZONE`, Format: `TIMESTAMP WITH LOCAL TIME ZONE`, IsDateTime: true},
	{Name: `TIME WITH TIME ZONE`, Format: `TIME WITH TIME ZONE`, IsDateTime: true},
	{Name: `DATETIME WITH TIME ZONE`, Format: `DATETIME WITH TIME ZONE`, IsDateTime: true},
	{Name: `BOOL`, Format: `BOOL`},
	{Name: `BOOLEAN`, Format: `BOOLEAN`},
}

// GBase 数据库 字段类型
var gBaseColumnTypeList = []*ColumnTypeInfo{
	{Name: `INT`, Format: `INT`, IsNumber: true, IsInteger: true, Comment: `整数 -2,147,483,647 至 2,147,483,647`},
	{Name: `INTEGER`, Format: `INTEGER`, Matches: []string{`TINYINT`, `MEDIUMINT`, `INT&&columnScale==0&&((columnLength>0&&columnLength<11)||(columnPrecision>0&&columnPrecision<11))`, `BIT&&columnLength==1||columnPrecision==1`, `INT1`, `INT2`, `INT4`, `BOOL`, `BOOLEAN`}, IsNumber: true, IsInteger: true},
	{Name: `BIGINT`, Format: `BIGINT`, Matches: []string{`INT8`, `INT&&columnScale==0&&((columnLength>=11)||(columnPrecision<=11))`}, IsNumber: true, IsInteger: true},
	{Name: `SERIAL`, Format: `SERIAL`, IsNumber: true, IsInteger: true, Comment: `自增类型，默认从1开始。可以设置初始值，如：serial(n)`},
	{Name: `SERIAL8`, Format: `SERIAL8`, IsNumber: true, IsInteger: true, Comment: `自增类型，默认从1开始。可以设置初始值，如：serial(n)`},
	{Name: `SMALLINT`, Format: `SMALLINT`, IsNumber: true, IsInteger: true},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`DOUBLE`, `FLOAT4`, `FLOAT8`, `DOUBLE PRECISION`}, IsNumber: true, IsFloat: true, Comment: `双精度浮点数值 存储最多带有 16 位有效数字的双精度浮点数值`},
	{Name: `SMALLFLOAT`, Format: `SMALLFLOAT`, IsNumber: true, IsFloat: true},
	{Name: `REAL`, Format: `REAL`, IsNumber: true, IsFloat: true},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, Matches: []string{`DECIMAL`, `DEC`, `REAL`, `NUMERIC`}, IsNumber: true, IsFloat: true, Comment: `存储实数的定点小数值 在小数部分中最多 20 位有效数字，或在小数点的左边最多 32 位有效数字。`},
	{Name: `NUMERIC`, Format: `NUMERIC($p, $s)`, IsNumber: true, IsFloat: true, Comment: `DECIMAL(p,s) 的符合 ANSI 的同义词 p最大精度是38位(十进制)`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`CHARACTER`}, IsString: true},
	{Name: `NCHAR`, Format: `NCHAR($l)`, IsString: true},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`VARCHAR2`, `BPCHAR`}, IsString: true},
	{Name: `LVARCHAR`, Format: `LVARCHAR($l)`, IsString: true},
	{Name: `NVARCHAR`, Format: `NVARCHAR($l)`, Matches: []string{`NVARCHAR2`}, IsString: true},
	{Name: `TEXT`, Format: `TEXT`, IsString: true},
	{Name: `CLOB`, Format: `CLOB`, Matches: []string{`TINYTEXT`, `MEDIUMTEXT`, `LONGTEXT`, `SET`, `ENUM`, `NCLOB`, `ROWID`, `UROWID`}, IsString: true},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BINARY`, `VARBINARY`, `TINYBLOB`, `MEDIUMBLOB`, `LONGBLOB`, `RAW`, `LONG RAW`, `BYTEA`}, IsBytes: true},
	{Name: `BYTE`, Format: `BYTE`, IsBytes: true},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true, Comment: `YYYY-MM-DD 1 年 1 月 1 日直至 9999 年 12 月 31 日`},
	{Name: `DATETIME`, Format: `DATETIME`, IsDateTime: true, Comment: `（年、月、日）和每日时间（小时、分、秒和几分之一秒） 1 年至 9999 年`,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `INTERVAL`, Format: `INTERVAL`, IsDateTime: true},
	{Name: `BOOLEAN`, Format: `BOOLEAN`, IsBoolean: true},
}

// OpenGauss 数据库 字段类型
var openGaussColumnTypeList = []*ColumnTypeInfo{
	{Name: `TINYINT`, Format: `TINYINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<3)||(columnPrecision>0&&columnPrecision<3))`, `INT1`, `BOOL`, `BOOLEAN`}, IsNumber: true, IsInteger: true, Comment: `1 Bytes 范围（有符号）(-128，127) 范围（无符号）(0，255) 小整数值`},
	{Name: `SMALLINT`, Format: `SMALLINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<6)||(columnPrecision>0&&columnPrecision<6))`, `INT2`}, IsNumber: true, IsInteger: true, Comment: `2 Bytes 范围（有符号）(-32 768，32 767) 范围（无符号）(0，65 535)  大整数值`},
	{Name: `MEDIUMINT`, Format: `MEDIUMINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<9)||(columnPrecision>0&&columnPrecision<9))`}, IsNumber: true, IsInteger: true, Comment: `3 Bytes 范围（有符号）(-8 388 608，8 388 607) 范围（无符号）(0，16 777 215)  大整数值`},
	{Name: `INT`, Format: `INT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<11)||(columnPrecision>0&&columnPrecision<11))`, `INT4`}, IsNumber: true, IsInteger: true, Comment: `4 Bytes 范围（有符号）(-2 147 483 648，2 147 483 647) 范围（无符号）(0，4 294 967 295)  大整数值`},
	{Name: `INTEGER`, Format: `INTEGER($l)`, IsNumber: true, IsInteger: true, Comment: `同上`},
	{Name: `BIGINT`, Format: `BIGINT($l)`, Matches: []string{`NUMBER&&columnScale==0`, `INT8`}, IsNumber: true, IsInteger: true, Comment: `8 Bytes 范围（有符号）(-9,223,372,036,854,775,808，9 223 372 036 854 775 807) 范围（无符号）(0，18 446 744 073 709 551 615)  极大整数值`},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`FLOAT4`}, IsNumber: true, IsFloat: true, Comment: `4 Bytes 范围（有符号）(-3.402 823 466 E+38，-1.175 494 351 E-38)，0，(1.175 494 351 E-38，3.402 823 466 351 E+38) 范围（无符号）0，(1.175 494 351 E-38，3.402 823 466 E+38)  单精度 浮点数值`},
	{Name: `DOUBLE`, Format: `DOUBLE`, Matches: []string{`FLOAT8`, `DOUBLE PRECISION`}, IsNumber: true, IsFloat: true, Comment: `8 Bytes 范围（有符号）(-1.797 693 134 862 315 7 E+308，-2.225 073 858 507 201 4 E-308)，0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308)范围（无符号）0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308)  双精度 浮点数值`},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, Matches: []string{`NUMBER`, `REAL`, `NUMERIC`}, IsNumber: true, IsFloat: true, Comment: `对DECIMAL(M,D) ，如果M>D，为M+2 否则为D+2 小数值`},
	{Name: `DEC`, Format: `DEC($p, $s)`, IsNumber: true, IsFloat: true, Comment: `同上`},
	{Name: `BIT`, Format: `BIT($p)`, Comment: `位字段类型。M 表示每个值的位数，范围为 1～64。如果 M 被省略，默认值为 1。如果为 BIT(M) 列分配的值的长度小于 M 位，在值的左边用 0 填充。例如，为 BIT(6) 列分配一个值 b'101'，其效果与分配 b'000101' 相同`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`NCHAR`, `CHARACTER`}, IsString: true, Comment: `0-255 bytes 定长字符串 固定长度非二进制字符串 M 字节，1<=M<=255`},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`VARCHAR2`, `NVARCHAR2`, `BPCHAR`}, IsString: true, Comment: `0-65535 bytes	变长字符串 变长非二进制字符串 L+1字节，在此，L< = M和 1<=M<=255`},
	{Name: `TINYTEXT`, Format: `TINYTEXT`, IsString: true, Comment: `0-255 bytes 短文本字符串`},
	{Name: `TEXT`, Format: `TEXT`, Matches: []string{`CLOB`, `NCLOB`, `ROWID`, `UROWID`}, IsString: true, Comment: `0-65 535 bytes 长文本数据`},
	{Name: `MEDIUMTEXT`, Format: `MEDIUMTEXT`, Matches: []string{`LONGVARCHAR`}, IsString: true, Comment: `0-16 777 215 bytes	中等长度文本数据`},
	{Name: `LONGTEXT`, Format: `LONGTEXT`, IsString: true, Comment: `0-4 294 967 295 bytes	极大文本数据`},
	{Name: `BINARY`, Format: `BINARY($l)`, IsBytes: true, Comment: `0-255 bytes 不超过 255 个字符的二进制字符串 固定长度二进制字符串 M 字节`},
	{Name: `VARBINARY`, Format: `VARBINARY($l)`, IsBytes: true, Comment: `可变长度二进制字符串 M+1 字节`},
	{Name: `TINYBLOB`, Format: `TINYBLOB`, IsBytes: true, Comment: `非常小的BLOB L+1 字节，在此，L<2^8`},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BFILE`, `RAW`, `BYTEA`}, IsBytes: true, Comment: `0-65 535 bytes 二进制形式的长文本数据 小 BLOB L+2 字节，在此，L<2^16`},
	{Name: `MEDIUMBLOB`, Format: `MEDIUMBLOB`, IsBytes: true, Comment: `0-16 777 215 bytes	二进制形式的中等长度文本数据 中等大小的BLOB	L+3 字节，在此，L<2^24`},
	{Name: `LONGBLOB`, Format: `LONGBLOB`, Matches: []string{`LONG RAW`}, IsBytes: true, Comment: `0-4 294 967 295 bytes	二进制形式的极大文本数据 非常大的BLOB	L+4 字节，在此，L<2^32`},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true, Comment: `3 bytes '-838:59:59'/'838:59:59' HH:MM:SS 时间值或持续时间`},
	{Name: `TIME`, Format: `TIME`, IsDateTime: true, Comment: `1 bytes 1901/2155 YYYY 年份值`},
	{Name: `YEAR`, Format: `YEAR`, IsDateTime: true, Comment: `8 bytes '1000-01-01 00:00:00' 到 '9999-12-31 23:59:59' YYYY-MM-DD hh:mm:ss 混合日期和时间值`},
	{Name: `DATETIME`, Format: `DATETIME`, Matches: []string{`DATETIME WITH TIME ZONE`}, IsDateTime: true, Comment: `4 bytes '1970-01-01 00:00:01' UTC 到 '2038-01-19 03:14:07' UTC 结束时间是第 2147483647 秒，北京时间 2038-1-19 11:14:07，格林尼治时间 2038年1月19日 凌晨 03:14:07 YYYY-MM-DD hh:mm:ss 混合日期和时间值，时间戳`,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`TIMESTAMP WITH TIME ZONE`, `TIMESTAMP WITH LOCAL TIME ZONE`, `INTERVAL YEAR TO MONTH`, `INTERVAL DAY TO SECOND`, `TIME WITH TIME ZONE`, `TIMESTAMP WITHOUT TIME ZONE`}, IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `ENUM`, Format: `ENUM`, IsEnum: true, Comment: `枚举类型，只能有一个枚举字符串值 1或2个字节，取决于枚举值的数目 (最大值为65535)`},
	{Name: `SET`, Format: `SET`, IsEnum: true, Comment: `一个设置，字符串对象可以有零个或 多个SET成员 1、2、3、4或8个字节，取决于集合 成员的数量（最多64个成员）`},
}

// Postgresql 数据库 字段类型
var postgresqlColumnTypeList = []*ColumnTypeInfo{
	{Name: `TINYINT`, Format: `TINYINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<3)||(columnPrecision>0&&columnPrecision<3))`, `INT1`, `BOOL`, `BOOLEAN`}, IsNumber: true, IsInteger: true, Comment: `1 Bytes 范围（有符号）(-128，127) 范围（无符号）(0，255) 小整数值`},
	{Name: `SMALLINT`, Format: `SMALLINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<6)||(columnPrecision>0&&columnPrecision<6))`, `INT2`}, IsNumber: true, IsInteger: true, Comment: `2 Bytes 范围（有符号）(-32 768，32 767) 范围（无符号）(0，65 535)  大整数值`},
	{Name: `MEDIUMINT`, Format: `MEDIUMINT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<9)||(columnPrecision>0&&columnPrecision<9))`}, IsNumber: true, IsInteger: true, Comment: `3 Bytes 范围（有符号）(-8 388 608，8 388 607) 范围（无符号）(0，16 777 215)  大整数值`},
	{Name: `INT`, Format: `INT($l)`, Matches: []string{`NUMBER&&columnScale==0&&((columnLength>0&&columnLength<11)||(columnPrecision>0&&columnPrecision<11))`, `INT4`}, IsNumber: true, IsInteger: true, Comment: `4 Bytes 范围（有符号）(-2 147 483 648，2 147 483 647) 范围（无符号）(0，4 294 967 295)  大整数值`},
	{Name: `INTEGER`, Format: `INTEGER($l)`, IsNumber: true, IsInteger: true, Comment: `同上`},
	{Name: `BIGINT`, Format: `BIGINT($l)`, Matches: []string{`NUMBER&&columnScale==0`, `INT8`}, IsNumber: true, IsInteger: true, Comment: `8 Bytes 范围（有符号）(-9,223,372,036,854,775,808，9 223 372 036 854 775 807) 范围（无符号）(0，18 446 744 073 709 551 615)  极大整数值`},
	{Name: `FLOAT`, Format: `FLOAT`, Matches: []string{`FLOAT4`}, IsNumber: true, IsFloat: true, Comment: `4 Bytes 范围（有符号）(-3.402 823 466 E+38，-1.175 494 351 E-38)，0，(1.175 494 351 E-38，3.402 823 466 351 E+38) 范围（无符号）0，(1.175 494 351 E-38，3.402 823 466 E+38)  单精度 浮点数值`},
	{Name: `DOUBLE`, Format: `DOUBLE`, Matches: []string{`FLOAT8`, `DOUBLE PRECISION`}, IsNumber: true, IsFloat: true, Comment: `8 Bytes 范围（有符号）(-1.797 693 134 862 315 7 E+308，-2.225 073 858 507 201 4 E-308)，0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308)范围（无符号）0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308)  双精度 浮点数值`},
	{Name: `DECIMAL`, Format: `DECIMAL($p, $s)`, Matches: []string{`NUMBER`, `REAL`, `NUMERIC`}, IsNumber: true, IsFloat: true, Comment: `对DECIMAL(M,D) ，如果M>D，为M+2 否则为D+2 小数值`},
	{Name: `DEC`, Format: `DEC($p, $s)`, IsNumber: true, IsFloat: true, Comment: `同上`},
	{Name: `BIT`, Format: `BIT($p)`, Comment: `位字段类型。M 表示每个值的位数，范围为 1～64。如果 M 被省略，默认值为 1。如果为 BIT(M) 列分配的值的长度小于 M 位，在值的左边用 0 填充。例如，为 BIT(6) 列分配一个值 b'101'，其效果与分配 b'000101' 相同`},
	{Name: `CHAR`, Format: `CHAR($l)`, Matches: []string{`NCHAR`, `CHARACTER`}, IsString: true, Comment: `0-255 bytes 定长字符串 固定长度非二进制字符串 M 字节，1<=M<=255`},
	{Name: `VARCHAR`, Format: `VARCHAR($l)`, Matches: []string{`VARCHAR2`, `NVARCHAR2`, `BPCHAR`}, IsString: true, Comment: `0-65535 bytes	变长字符串 变长非二进制字符串 L+1字节，在此，L< = M和 1<=M<=255`},
	{Name: `TINYTEXT`, Format: `TINYTEXT`, IsString: true, Comment: `0-255 bytes 短文本字符串`},
	{Name: `TEXT`, Format: `TEXT`, Matches: []string{`CLOB`, `NCLOB`, `ROWID`, `UROWID`}, IsString: true, Comment: `0-65 535 bytes 长文本数据`},
	{Name: `MEDIUMTEXT`, Format: `MEDIUMTEXT`, Matches: []string{`LONGVARCHAR`}, IsString: true, Comment: `0-16 777 215 bytes	中等长度文本数据`},
	{Name: `LONGTEXT`, Format: `LONGTEXT`, IsString: true, Comment: `0-4 294 967 295 bytes	极大文本数据`},
	{Name: `BINARY`, Format: `BINARY($l)`, IsBytes: true, Comment: `0-255 bytes 不超过 255 个字符的二进制字符串 固定长度二进制字符串 M 字节`},
	{Name: `VARBINARY`, Format: `VARBINARY($l)`, IsBytes: true, Comment: `可变长度二进制字符串 M+1 字节`},
	{Name: `TINYBLOB`, Format: `TINYBLOB`, IsBytes: true, Comment: `非常小的BLOB L+1 字节，在此，L<2^8`},
	{Name: `BLOB`, Format: `BLOB`, Matches: []string{`BFILE`, `RAW`, `BYTEA`}, IsBytes: true, Comment: `0-65 535 bytes 二进制形式的长文本数据 小 BLOB L+2 字节，在此，L<2^16`},
	{Name: `MEDIUMBLOB`, Format: `MEDIUMBLOB`, IsBytes: true, Comment: `0-16 777 215 bytes	二进制形式的中等长度文本数据 中等大小的BLOB	L+3 字节，在此，L<2^24`},
	{Name: `LONGBLOB`, Format: `LONGBLOB`, Matches: []string{`LONG RAW`}, IsBytes: true, Comment: `0-4 294 967 295 bytes	二进制形式的极大文本数据 非常大的BLOB	L+4 字节，在此，L<2^32`},
	{Name: `DATE`, Format: `DATE`, IsDateTime: true, Comment: `3 bytes '-838:59:59'/'838:59:59' HH:MM:SS 时间值或持续时间`},
	{Name: `TIME`, Format: `TIME`, IsDateTime: true, Comment: `1 bytes 1901/2155 YYYY 年份值`},
	{Name: `YEAR`, Format: `YEAR`, IsDateTime: true, Comment: `8 bytes '1000-01-01 00:00:00' 到 '9999-12-31 23:59:59' YYYY-MM-DD hh:mm:ss 混合日期和时间值`},
	{Name: `DATETIME`, Format: `DATETIME`, Matches: []string{`DATETIME WITH TIME ZONE`}, IsDateTime: true, Comment: `4 bytes '1970-01-01 00:00:01' UTC 到 '2038-01-19 03:14:07' UTC 结束时间是第 2147483647 秒，北京时间 2038-1-19 11:14:07，格林尼治时间 2038年1月19日 凌晨 03:14:07 YYYY-MM-DD hh:mm:ss 混合日期和时间值，时间戳`,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `TIMESTAMP`, Format: `TIMESTAMP`, Matches: []string{`TIMESTAMP WITH TIME ZONE`, `TIMESTAMP WITH LOCAL TIME ZONE`, `INTERVAL YEAR TO MONTH`, `INTERVAL DAY TO SECOND`, `TIME WITH TIME ZONE`, `TIMESTAMP WITHOUT TIME ZONE`}, IsDateTime: true,
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}

			return
		},
	},
	{Name: `ENUM`, Format: `ENUM`, IsEnum: true, Comment: `枚举类型，只能有一个枚举字符串值 1或2个字节，取决于枚举值的数目 (最大值为65535)`},
	{Name: `SET`, Format: `SET`, IsEnum: true, Comment: `一个设置，字符串对象可以有零个或 多个SET成员 1、2、3、4或8个字节，取决于集合 成员的数量（最多64个成员）`},
}

// DB2 数据库 字段类型
var db2ColumnTypeList = []*ColumnTypeInfo{
}

