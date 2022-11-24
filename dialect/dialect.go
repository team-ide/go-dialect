package dialect

import (
	"errors"
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
	ColumnTypePack(column *ColumnModel) (columnTypePack string, err error)
	GetIndexTypeInfos() (indexTypeInfoList []*IndexTypeInfo)
	//FormatDefaultValue(column *ColumnModel) (defaultValue string)
	//ToColumnTypeInfo(columnType string) (columnTypeInfo *ColumnTypeInfo, length, decimal int, err error)

	OwnerNamePack(param *ParamModel, ownerName string) string
	TableNamePack(param *ParamModel, tableName string) string
	ColumnNamePack(param *ParamModel, columnName string) string
	ColumnNamesPack(param *ParamModel, columnNames []string) string
	SqlValuePack(param *ParamModel, column *ColumnModel, value interface{}) string
	ColumnDefaultPack(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error)
	// IsSqlEnd 判断SQL是否以 分号 结尾
	IsSqlEnd(sqlInfo string) bool
	// SqlSplit 根据 分号 分割多条SQL
	SqlSplit(sqlInfo string) []string

	OwnerModel(data map[string]interface{}) (owner *OwnerModel, err error)
	OwnersSelectSql(param *ParamModel) (sql string, err error)
	OwnerSelectSql(param *ParamModel, ownerName string) (sql string, err error)
	OwnerCreateSql(param *ParamModel, owner *OwnerModel) (sqlList []string, err error)
	OwnerDeleteSql(param *ParamModel, ownerName string) (sqlList []string, err error)

	TableModel(data map[string]interface{}) (table *TableModel, err error)
	TablesSelectSql(param *ParamModel, ownerName string) (sql string, err error)
	TableSelectSql(param *ParamModel, ownerName string, tableName string) (sql string, err error)
	TableCreateSql(param *ParamModel, ownerName string, table *TableModel) (sqlList []string, err error)
	TableCommentSql(param *ParamModel, ownerName string, tableName string, tableComment string) (sqlList []string, err error)
	TableRenameSql(param *ParamModel, ownerName string, oldTableName string, tableName string) (sqlList []string, err error)
	TableDeleteSql(param *ParamModel, ownerName string, tableName string) (sqlList []string, err error)

	ColumnModel(data map[string]interface{}) (table *ColumnModel, err error)
	ColumnsSelectSql(param *ParamModel, ownerName string, tableName string) (sql string, err error)
	ColumnSelectSql(param *ParamModel, ownerName string, tableName string, columnName string) (sql string, err error)
	ColumnAddSql(param *ParamModel, ownerName string, tableName string, column *ColumnModel) (sqlList []string, err error)
	ColumnCommentSql(param *ParamModel, ownerName string, tableName string, columnName string, columnComment string) (sqlList []string, err error)
	ColumnUpdateSql(param *ParamModel, ownerName string, tableName string, oldColumn *ColumnModel, column *ColumnModel) (sqlList []string, err error)
	ColumnDeleteSql(param *ParamModel, ownerName string, tableName string, columnName string) (sqlList []string, err error)

	PrimaryKeyModel(data map[string]interface{}) (primaryKey *PrimaryKeyModel, err error)
	PrimaryKeysSelectSql(param *ParamModel, ownerName string, tableName string) (sql string, err error)
	PrimaryKeyAddSql(param *ParamModel, ownerName string, tableName string, columnNames []string) (sqlList []string, err error)
	PrimaryKeyDeleteSql(param *ParamModel, ownerName string, tableName string) (sqlList []string, err error)

	IndexModel(data map[string]interface{}) (index *IndexModel, err error)
	IndexesSelectSql(param *ParamModel, ownerName string, tableName string) (sql string, err error)
	IndexAddSql(param *ParamModel, ownerName string, tableName string, index *IndexModel) (sqlList []string, err error)
	IndexDeleteSql(param *ParamModel, ownerName string, tableName string, indexName string) (sqlList []string, err error)

	PackPageSql(selectSql string, pageSize int, pageNo int) (pageSql string)
	ReplaceSqlVariable(sqlInfo string, args []interface{}) (variableSql string)
	InsertSql(param *ParamModel, insert *InsertModel) (sqlList []string, err error)

	InsertDataListSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}) (sqlList []string, batchSqlList []string, err error)

	DataListInsertSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error)
	DataListUpdateSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error)
	DataListDeleteSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error)
	DataListSelectSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, whereList []*Where, orderList []*Order) (sql string, values []interface{}, err error)
}

var (
	TypeMysql      = &Type{Name: "mysql"}
	TypeSqlite     = &Type{Name: "sqlite"}
	TypeOracle     = &Type{Name: "oracle"}
	TypeDM         = &Type{Name: "dm"}
	TypeKingBase   = &Type{Name: "kingbase"}
	TypeShenTong   = &Type{Name: "shentong"}
	TypePostgresql = &Type{Name: "postgresql"}
)

func NewDialect(dialectType string) (dia Dialect, err error) {
	switch strings.ToLower(dialectType) {
	case "mysql":
		dia, err = NewMappingDialect(NewMappingMysql())
		break
	case "sqlite", "sqlite3":
		dia, err = NewMappingDialect(NewMappingSqlite())
		break
	case "dameng", "dm":
		dia, err = NewMappingDialect(NewMappingDM())
		break
	case "kingbase", "kb":
		dia, err = NewMappingDialect(NewMappingKingBase())
		break
	case "oracle":
		dia, err = NewMappingDialect(NewMappingOracle())
		break
	case "shentong", "st":
		dia, err = NewMappingDialect(NewMappingShenTong())
		break
	case "postgresql", "ps":
		dia, err = NewMappingDialect(NewMappingPostgresql())
		break
	default:
		err = errors.New("dialect type [" + dialectType + "] not support ")
		return
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

func packingValue(column *ColumnModel, columnTypeInfo *ColumnTypeInfo, packingCharacter string, escapeChar string, value interface{}) string {
	if value == nil {
		return "NULL"
	}
	vOf := reflect.ValueOf(value)
	if vOf.Kind() == reflect.Ptr {
		if vOf.IsNil() {
			return "NULL"
		}
		return packingValue(column, columnTypeInfo, packingCharacter, escapeChar, vOf.Elem().Interface())
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

	if columnTypeInfo != nil {
		if columnTypeInfo.SqlValuePack != nil {
			return columnTypeInfo.SqlValuePack(valueString)
		}
	}
	if valueString == "" && column != nil && column.ColumnNotNull && column.ColumnDefault == "" {

	} else {
		if columnTypeInfo != nil {
			if columnTypeInfo.IsNumber {
				if valueString == "" {
					return "NULL"
				}
				return valueString
			} else if columnTypeInfo.IsEnum {
				if valueString == "" {
					return "NULL"
				}
			}
		}
	}

	if packingCharacter == "" {
		return valueString
	}
	return formatStringValue(packingCharacter, escapeChar, valueString)
}

func formatStringValue(packingCharacter string, escapeChar string, valueString string) string {
	if packingCharacter == "" {
		return valueString
	}
	//valueString = strings.ReplaceAll(valueString, "\n", `\\n`)
	out := packingCharacter
	ss := strings.Split(valueString, "")
	var valueLen = len(ss)
	for i := 0; i < valueLen; i++ {
		s := ss[i]
		switch s {
		case packingCharacter:
			out += escapeChar + packingCharacter
			break
		case "\\":
			if escapeChar == "\\" {
				out += "\\"
			}
			out += "\\"
			break
		default:
			out += s
			break
		}
	}
	out += packingCharacter
	return out
}
