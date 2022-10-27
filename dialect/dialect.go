package dialect

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Type struct {
	Name string `json:"name"`
}

type Dialect interface {
	DialectType() (dialectType *Type)
	GetColumnTypeInfos() (columnTypeInfoList []*ColumnTypeInfo)
	GetColumnTypeInfo(typeName string) (columnTypeInfo *ColumnTypeInfo, err error)
	FormatColumnType(typeName string, length, decimal int) (columnType string, err error)
	ToColumnTypeInfo(columnType string) (columnTypeInfo *ColumnTypeInfo, length, decimal int, err error)

	OwnerModel(data map[string]interface{}) (database *OwnerModel, err error)
	OwnersSelectSql() (sql string, err error)
	OwnerCreateSql(param *GenerateParam, database *OwnerModel) (sqlList []string, err error)
	OwnerDeleteSql(param *GenerateParam, ownerName string) (sqlList []string, err error)

	TableModel(data map[string]interface{}) (table *TableModel, err error)
	TablesSelectSql(ownerName string) (sql string, err error)
	TableSelectSql(ownerName string, tableName string) (sql string, err error)
	TableCreateSql(param *GenerateParam, ownerName string, table *TableModel) (sqlList []string, err error)
	TableCommentSql(param *GenerateParam, ownerName string, tableName string, comment string) (sqlList []string, err error)
	TableDeleteSql(param *GenerateParam, ownerName string, tableName string) (sqlList []string, err error)

	ColumnModel(data map[string]interface{}) (table *ColumnModel, err error)
	ColumnsSelectSql(ownerName string, tableName string) (sql string, err error)
	ColumnSelectSql(ownerName string, tableName string, columnName string) (sql string, err error)
	ColumnAddSql(param *GenerateParam, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error)
	ColumnCommentSql(param *GenerateParam, ownerName string, tableName string, columnName string, comment string) (sqlList []string, err error)
	ColumnUpdateSql(param *GenerateParam, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error)
	ColumnDeleteSql(param *GenerateParam, ownerName string, tableName string, columnName string) (sqlList []string, err error)

	PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error)
	PrimaryKeysSelectSql(ownerName string, tableName string) (sql string, err error)
	PrimaryKeyAddSql(param *GenerateParam, ownerName string, tableName string, primaryKeys []string) (sqlList []string, err error)
	PrimaryKeyDeleteSql(param *GenerateParam, ownerName string, tableName string) (sqlList []string, err error)

	IndexModel(data map[string]interface{}) (index *IndexModel, err error)
	IndexesSelectSql(ownerName string, tableName string) (sql string, err error)
	IndexAddSql(param *GenerateParam, ownerName string, tableName string, index *IndexModel) (sqlList []string, err error)
	IndexUpdateSql(param *GenerateParam, ownerName string, tableName string, index *IndexModel) (sqlList []string, err error)
	IndexDeleteSql(param *GenerateParam, ownerName string, tableName string, indexName string) (sqlList []string, err error)

	InsertSql(param *GenerateParam, insert *InsertModel) (sqlList []string, err error)
}

var (
	Mysql          = NewMysqlDialect()
	MysqlType      = &Type{Name: "mysql"}
	Sqlite         = NewSqliteDialect()
	SqliteType     = &Type{Name: "sqlite"}
	Oracle         = NewOracleDialect()
	OracleType     = &Type{Name: "oracle"}
	DaMen          = NewDaMenDialect()
	DaMenType      = &Type{Name: "damen"}
	KinBase        = NewKinBaseDialect()
	KinBaseType    = &Type{Name: "kinbase"}
	ShenTong       = NewShenTongDialect()
	ShenTongType   = &Type{Name: "shentong"}
	Postgresql     = NewPostgresqlDialect()
	PostgresqlType = &Type{Name: "postgresql"}
)

type GenerateParam struct {
	AppendOwner            bool   `json:"appendOwner"`
	CharacterSetName       string `json:"characterSetName"`
	CollationName          string `json:"collationName"`
	OwnerPackingCharacter  string `json:"ownerPackingCharacter"`
	TablePackingCharacter  string `json:"tablePackingCharacter"`
	ColumnPackingCharacter string `json:"columnPackingCharacter"`
	StringPackingCharacter string `json:"stringPackingCharacter"`
	AppendSqlValue         bool   `json:"appendSqlValue"`
	DateFunction           string `json:"dateFunction"`
	OpenTransaction        bool   `json:"openTransaction"`
	ErrorContinue          bool   `json:"errorContinue"`
}

func (this_ *GenerateParam) PackingCharacterOwner(value string) string {
	return this_.packingCharacterOwner(value)
}

func (this_ *GenerateParam) packingCharacterOwner(value string) string {
	if this_.OwnerPackingCharacter == "" {
		return value
	}
	return this_.OwnerPackingCharacter + value + this_.OwnerPackingCharacter
}

func (this_ *GenerateParam) PackingCharacterTable(value string) string {
	return this_.packingCharacterTable(value)
}

func (this_ *GenerateParam) packingCharacterTable(value string) string {
	if this_.TablePackingCharacter == "" {
		return value
	}
	return this_.TablePackingCharacter + value + this_.TablePackingCharacter
}

func (this_ *GenerateParam) PackingCharacterColumn(value string) string {
	return this_.packingCharacterColumn(value)
}

func (this_ *GenerateParam) packingCharacterColumn(value string) string {
	if this_.ColumnPackingCharacter == "" {
		return value
	}
	value = strings.ReplaceAll(value, `""`, "")
	value = strings.ReplaceAll(value, `'`, "")
	value = strings.ReplaceAll(value, "`", "")
	return this_.ColumnPackingCharacter + value + this_.ColumnPackingCharacter
}

func (this_ *GenerateParam) PackingCharacterColumns(value string) string {
	return this_.packingCharacterColumns(value)
}

func (this_ *GenerateParam) packingCharacterColumns(columns string) string {
	if this_.ColumnPackingCharacter == "" {
		return columns
	}
	res := ""
	columnList := strings.Split(columns, ",")

	for _, column := range columnList {
		res += this_.packingCharacterColumn(column) + ","
	}
	res = strings.TrimSuffix(res, ",")
	return res
}

func (this_ *GenerateParam) PackingCharacterColumnStringValue(dia Dialect, tableColumn *ColumnModel, value interface{}) string {
	return this_.packingCharacterColumnStringValue(dia, tableColumn, value)
}

func (this_ *GenerateParam) packingCharacterColumnStringValue(dia Dialect, tableColumn *ColumnModel, value interface{}) string {
	var formatColumnValue = this_.formatColumnValue(dia, tableColumn, value)
	if formatColumnValue == nil {
		return "NULL"
	}
	vOf := reflect.ValueOf(value)
	if vOf.Kind() == reflect.Ptr {
		if vOf.IsNil() {
			return "NULL"
		}
		return this_.packingCharacterColumnStringValue(dia, tableColumn, vOf.Elem().Interface())
	}

	baseValue, isBaseValue := GetBaseTypeValue(value)
	if isBaseValue {
		value = baseValue
	}
	var valueString string
	switch v := formatColumnValue.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case time.Time:
		if v.IsZero() {
			return "NULL"
		}
		valueString = v.Format("2006-01-02 15:04:05")
		if this_.DateFunction != "" {
			return strings.ReplaceAll(this_.DateFunction, "$value", valueString)
		}
		break
	case string:
		valueString = v
		break
	case []byte:
		valueString = string(v)
	default:
		valueString = GetStringValue(value)
		break
	}
	if this_.StringPackingCharacter == "" {
		return valueString
	}
	return formatStringValue(this_.StringPackingCharacter, valueString)
}

func (this_ *GenerateParam) FormatColumnValue(dia Dialect, tableColumn *ColumnModel, value interface{}) interface{} {
	return this_.formatColumnValue(dia, tableColumn, value)
}
func (this_ *GenerateParam) formatColumnValue(dia Dialect, tableColumn *ColumnModel, value interface{}) interface{} {

	var IsDateTime bool
	var IsNumber bool
	var Decimal int
	if tableColumn != nil {

		columnTypeInfo, err := dia.GetColumnTypeInfo(tableColumn.Type)
		if err != nil {
			fmt.Printf("GetColumnTypeInfo error %s", err)
			return value
		}
		IsDateTime = columnTypeInfo.IsDateTime
		IsNumber = columnTypeInfo.IsNumber
		Decimal = tableColumn.Decimal
	}

	if value == nil {
		return value
	}
	var stringValue = GetStringValue(value)
	if IsNumber {
		if stringValue == "" {
			return nil
		}
		if Decimal > 0 {
			f64, err := strconv.ParseFloat(stringValue, 64)
			if err != nil {
				fmt.Printf("value [%s] ParseFloat error %s", stringValue, err)
				return value
			}
			return f64
		} else {
			i64, err := strconv.ParseInt(stringValue, 10, 64)
			if err != nil {
				fmt.Printf("value [%s] ParseInt error %s", stringValue, err)
				return value
			}
			return i64
		}
	}
	if IsDateTime {
		if stringValue == "" {
			return nil
		}
		format := "2006-01-02 15:04:05.000"
		valueLen := len(stringValue)
		if valueLen >= len("2006-01-02 15:04:05.000") {
			format = "2006-01-02 15:04:05.000"
		} else if valueLen >= len("2006-01-02 15:04:05") {
			format = "2006-01-02 15:04:05"
		} else if valueLen >= len("2006-01-02 15:04") {
			format = "2006-01-02 15:04"
		} else if valueLen >= len("2006-01-02 15") {
			format = "2006-01-02 15"
		} else if valueLen >= len("2006-01-02") {
			format = "2006-01-02"
		} else if valueLen >= len("15:04:05") {
			format = "15:04:05"
		} else if valueLen >= len("15:04") {
			format = "15:04"
		} else if valueLen >= len("2006") {
			format = "2006"
		}
		timeValue, err := time.ParseInLocation(format, stringValue, time.Local)
		if err != nil {
			fmt.Printf("value [%s] to time error %s", stringValue, err)
			return value
		}
		return timeValue
	}
	return value
}

func formatStringValue(packingCharacter string, valueString string) string {
	if packingCharacter == "" {
		return valueString
	}
	ss := strings.Split(valueString, "")
	out := packingCharacter
	for _, s := range ss {
		switch s {
		case packingCharacter:
			out += "\\" + s
		case "\\":
			out += "\\" + s
		default:
			out += s
		}
	}
	out += packingCharacter
	return out
}
