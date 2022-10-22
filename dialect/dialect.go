package dialect

import (
	"encoding/json"
	"fmt"
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

	DatabaseModel(data map[string]interface{}) (database *DatabaseModel, err error)
	DatabasesSelectSql() (sql string, err error)
	DatabaseCreateSql(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error)
	DatabaseDeleteSql(param *GenerateParam, databaseName string) (sqlList []string, err error)

	TableModel(data map[string]interface{}) (table *TableModel, err error)
	TablesSelectSql(databaseName string) (sql string, err error)
	TableSelectSql(databaseName string, tableName string) (sql string, err error)
	TableCreateSql(param *GenerateParam, databaseName string, table *TableModel) (sqlList []string, err error)
	TableCommentSql(param *GenerateParam, databaseName string, tableName string, comment string) (sqlList []string, err error)
	TableDeleteSql(param *GenerateParam, databaseName string, tableName string) (sqlList []string, err error)

	ColumnModel(data map[string]interface{}) (table *ColumnModel, err error)
	ColumnsSelectSql(databaseName string, tableName string) (sql string, err error)
	ColumnSelectSql(databaseName string, tableName string, columnName string) (sql string, err error)
	ColumnAddSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error)
	ColumnUpdateSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error)
	ColumnRenameSql(param *GenerateParam, databaseName string, tableName string, column *ColumnModel) (sqlList []string, err error)
	ColumnDeleteSql(param *GenerateParam, databaseName string, tableName string, columnName string) (sqlList []string, err error)

	IndexModel(data map[string]interface{}) (index *IndexModel, err error)
	IndexesSelectSql(databaseName string, tableName string) (sql string, err error)
	IndexSelectSql(databaseName string, tableName string, indexName string) (sql string, err error)
	IndexAddSql(param *GenerateParam, databaseName string, tableName string, index *IndexModel) (sqlList []string, err error)
	IndexUpdateSql(param *GenerateParam, databaseName string, tableName string, index *IndexModel) (sqlList []string, err error)
	IndexDeleteSql(param *GenerateParam, databaseName string, tableName string, indexName string) (sqlList []string, err error)
	IndexRenameSql(param *GenerateParam, databaseName string, tableName string, indexName string, rename string) (sqlList []string, err error)

	PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error)
	PrimaryKeysSelectSql(databaseName string, tableName string) (sql string, err error)
	PrimaryKeyAddSql(param *GenerateParam, databaseName string, tableName string, primaryKeys []string) (sqlList []string, err error)
	PrimaryKeyDeleteSql(param *GenerateParam, databaseName string, tableName string, primaryKeys []string) (sqlList []string, err error)
}

var (
	Default     = NewDefaultDialect()
	DefaultType = &Type{Name: "default"}
	Mysql       = NewMysqlDialect()
	MysqlType   = &Type{Name: "mysql"}
)

type GenerateParam struct {
	DatabaseType             string `json:"databaseType" column:"databaseType"`
	GenerateDatabase         bool   `json:"generateDatabase" column:"generateDatabase"`
	AppendDatabase           bool   `json:"appendDatabase" column:"appendDatabase"`
	CharacterSet             string `json:"characterSet" column:"characterSet"`
	Collate                  string `json:"collate" column:"collate"`
	DatabasePackingCharacter string `json:"databasePackingCharacter" column:"databasePackingCharacter"`
	TablePackingCharacter    string `json:"tablePackingCharacter" column:"tablePackingCharacter"`
	ColumnPackingCharacter   string `json:"columnPackingCharacter" column:"columnPackingCharacter"`
	StringPackingCharacter   string `json:"stringPackingCharacter" column:"stringPackingCharacter"`
	AppendSqlValue           bool   `json:"appendSqlValue" column:"appendSqlValue"`
	DateFunction             string `json:"dateFunction" column:"dateFunction"`
	OpenTransaction          bool   `json:"openTransaction"`
	ErrorContinue            bool   `json:"errorContinue"`
}

func (param *GenerateParam) PackingCharacterDatabase(value string) string {
	return param.packingCharacterDatabase(value)
}

func (param *GenerateParam) packingCharacterDatabase(value string) string {
	if param.DatabasePackingCharacter == "" {
		return value
	}
	return param.DatabasePackingCharacter + value + param.DatabasePackingCharacter
}

func (param *GenerateParam) PackingCharacterTable(value string) string {
	return param.packingCharacterTable(value)
}

func (param *GenerateParam) packingCharacterTable(value string) string {
	if param.TablePackingCharacter == "" {
		return value
	}
	return param.TablePackingCharacter + value + param.TablePackingCharacter
}

func (param *GenerateParam) PackingCharacterColumn(value string) string {
	return param.packingCharacterColumn(value)
}

func (param *GenerateParam) packingCharacterColumn(value string) string {
	if param.ColumnPackingCharacter == "" {
		return value
	}
	value = strings.ReplaceAll(value, `""`, "")
	value = strings.ReplaceAll(value, `'`, "")
	value = strings.ReplaceAll(value, "`", "")
	return param.ColumnPackingCharacter + value + param.ColumnPackingCharacter
}

func (param *GenerateParam) PackingCharacterColumns(value string) string {
	return param.packingCharacterColumns(value)
}

func (param *GenerateParam) packingCharacterColumns(columns string) string {
	if param.ColumnPackingCharacter == "" {
		return columns
	}
	res := ""
	columnList := strings.Split(columns, ",")

	for _, column := range columnList {
		res += param.packingCharacterColumn(column) + ","
	}
	res = strings.TrimSuffix(res, ",")
	return res
}

func (param *GenerateParam) PackingCharacterColumnStringValue(dia Dialect, tableColumn *ColumnModel, value interface{}) string {
	return param.packingCharacterColumnStringValue(dia, tableColumn, value)
}

func (param *GenerateParam) packingCharacterColumnStringValue(dia Dialect, tableColumn *ColumnModel, value interface{}) string {
	var formatColumnValue = param.formatColumnValue(dia, tableColumn, value)
	if formatColumnValue == nil {
		return "NULL"
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
		if param.DateFunction != "" {
			return strings.ReplaceAll(param.DateFunction, "$value", valueString)
		}
		break
	case string:
		valueString = v
		break
	case []byte:
		valueString = string(v)
	default:
		newValue, _ := json.Marshal(value)
		valueString = string(newValue)
		break
	}
	if param.StringPackingCharacter == "" {
		return valueString
	}
	return formatStringValue(param.StringPackingCharacter, valueString)
}

func (param *GenerateParam) FormatColumnValue(dia Dialect, tableColumn *ColumnModel, value interface{}) interface{} {
	return param.formatColumnValue(dia, tableColumn, value)
}
func (param *GenerateParam) formatColumnValue(dia Dialect, tableColumn *ColumnModel, value interface{}) interface{} {

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

func GetStringValue(value interface{}) string {

	var valueString string
	switch v := value.(type) {
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
		break
	case string:
		valueString = v
		break
	case []byte:
		valueString = string(v)
	default:
		newValue, _ := json.Marshal(value)
		valueString = string(newValue)
		break
	}
	return valueString
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
