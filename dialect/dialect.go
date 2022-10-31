package dialect

import (
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
	FormatColumnType(column *ColumnModel) (columnType string, err error)
	FormatDefaultValue(column *ColumnModel) (defaultValue string)
	ToColumnTypeInfo(columnType string) (columnTypeInfo *ColumnTypeInfo, length, decimal int, err error)

	PackOwner(ownerName string) string
	PackTable(tableName string) string
	PackColumn(columnName string) string
	PackColumns(columnNames []string) string
	PackValueForSql(column *ColumnModel, value interface{}) string
	// IsSqlEnd 判断SQL是否以 分号 结尾
	IsSqlEnd(sqlInfo string) bool
	// SqlSplit 根据 分号 分割多条SQL
	SqlSplit(sqlInfo string) []string

	OwnerModel(data map[string]interface{}) (owner *OwnerModel, err error)
	OwnersSelectSql() (sql string, err error)
	OwnerSelectSql(ownerName string) (sql string, err error)
	OwnerCreateSql(owner *OwnerModel) (sqlList []string, err error)
	OwnerDeleteSql(ownerName string) (sqlList []string, err error)

	TableModel(data map[string]interface{}) (table *TableModel, err error)
	TablesSelectSql(ownerName string) (sql string, err error)
	TableSelectSql(ownerName string, tableName string) (sql string, err error)
	TableCreateSql(ownerName string, table *TableModel) (sqlList []string, err error)
	TableCommentSql(ownerName string, tableName string, comment string) (sqlList []string, err error)
	TableRenameSql(ownerName string, oldTableName string, newTableName string) (sqlList []string, err error)
	TableDeleteSql(ownerName string, tableName string) (sqlList []string, err error)

	ColumnModel(data map[string]interface{}) (table *ColumnModel, err error)
	ColumnsSelectSql(ownerName string, tableName string) (sql string, err error)
	ColumnSelectSql(ownerName string, tableName string, columnName string) (sql string, err error)
	ColumnAddSql(ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error)
	ColumnCommentSql(ownerName string, tableName string, columnName string, comment string) (sqlList []string, err error)
	ColumnUpdateSql(ownerName string, tableName string, oldColumn *ColumnModel, newColumn *ColumnModel) (sqlList []string, err error)
	ColumnDeleteSql(ownerName string, tableName string, columnName string) (sqlList []string, err error)

	PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error)
	PrimaryKeysSelectSql(ownerName string, tableName string) (sql string, err error)
	PrimaryKeyAddSql(ownerName string, tableName string, primaryKeys []string) (sqlList []string, err error)
	PrimaryKeyDeleteSql(ownerName string, tableName string) (sqlList []string, err error)

	IndexModel(data map[string]interface{}) (index *IndexModel, err error)
	IndexesSelectSql(ownerName string, tableName string) (sql string, err error)
	IndexAddSql(ownerName string, tableName string, index *IndexModel) (sqlList []string, err error)
	IndexDeleteSql(ownerName string, tableName string, indexName string) (sqlList []string, err error)

	InsertSql(insert *InsertModel) (sqlList []string, err error)
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

func GetDialect(dialectType string) (dia Dialect) {
	switch strings.ToLower(dialectType) {
	case "mysql":
		dia = Mysql
		break
	case "sqlite", "sqlite3":
		dia = Sqlite
		break
	case "damen", "dm":
		dia = DaMen
		break
	case "kingbase", "kb":
		dia = KinBase
		break
	case "oracle":
		dia = Oracle
		break
	case "shentong", "st":
		dia = ShenTong
		break
	case "postgresql", "ps":
		dia = Postgresql
		break
	}
	return
}

func packingName(packingCharacter, name string) string {
	name = strings.ReplaceAll(name, `""`, "")
	name = strings.ReplaceAll(name, `'`, "")
	name = strings.ReplaceAll(name, "`", "")
	name = strings.TrimSpace(name)
	if packingCharacter == "" {
		return name
	}
	return packingCharacter + name + packingCharacter
}

func packingNames(packingCharacter string, names []string) string {
	return packingValues(packingCharacter, names)
}

func packingValues(packingCharacter string, values []string) string {

	res := ""

	for _, value := range values {
		res += packingName(packingCharacter, value) + ", "
	}
	res = strings.TrimSuffix(res, ", ")
	return res
}

func packingValue(columnTypeInfo *ColumnTypeInfo, packingCharacter string, appendCharacter string, value interface{}) string {
	if value == nil {
		return "NULL"
	}
	vOf := reflect.ValueOf(value)
	if vOf.Kind() == reflect.Ptr {
		if vOf.IsNil() {
			return "NULL"
		}
		return packingValue(columnTypeInfo, packingCharacter, appendCharacter, vOf.Elem().Interface())
	}

	baseValue, isBaseValue := GetBaseTypeValue(value)
	if isBaseValue {
		value = baseValue
	}
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
		//if this_.DateFunction != "" {
		//	return strings.ReplaceAll(this_.DateFunction, "$value", valueString)
		//}
		break
	case string:
		valueString = v
		break
	case []byte:
		valueString = string(v)
		break
	default:
		valueString = GetStringValue(value)
		break
	}

	if columnTypeInfo != nil && columnTypeInfo.IsNumber {
		if valueString == "" {
			return "NULL"
		}
		return valueString
	}

	if packingCharacter == "" {
		return valueString
	}
	return formatStringValue(packingCharacter, appendCharacter, valueString)
}

func formatStringValue(packingCharacter string, appendCharacter string, valueString string) string {
	if packingCharacter == "" {
		return valueString
	}
	//valueString = strings.ReplaceAll(valueString, "\n", `\\n`)
	out := packingCharacter
	var valueLen = len(valueString)
	for i := 0; i < valueLen; i++ {
		s := valueString[i]
		switch s {
		case packingCharacter[0]:
			out += appendCharacter + packingCharacter
			break
		case '\\':
			if appendCharacter == "\\" {
				out += "\\"
			}
			out += "\\"
			break
		default:
			out += string(s)
			break
		}
	}
	out += packingCharacter
	return out
}
