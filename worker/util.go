package worker

import (
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"reflect"
)

func GetRefValue(bean interface{}) reflect.Value {
	if IsPtr(bean) {
		return reflect.ValueOf(bean).Elem()
	}
	return reflect.ValueOf(bean)
}

func GetRefType(bean interface{}) reflect.Type {
	if IsPtr(bean) {
		return reflect.TypeOf(bean).Elem()
	}
	return reflect.TypeOf(bean)
}

func IsPtr(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

func GetSqlValueCache(columnTypes []*sql.ColumnType) (cache []interface{}) {
	cache = make([]interface{}, len(columnTypes)) //临时存储每行数据
	for index, _ := range cache {
		columnType := columnTypes[index]
		ct := columnType.ScanType()
		if ct == nil {
			cache[index] = new(sql.NullString)
			continue
		}
		//println("GetSqlValueCache type [" + columnType.ScanType().String() + "] columnName [" + columnType.Name() + "] databaseType [" + columnType.DatabaseTypeName() + "]")
		switch ct.String() {
		case "sql.NullString":
			cache[index] = new(sql.NullString)
			break
		case "sql.NullBool":
			cache[index] = new(sql.NullBool)
			break
		case "sql.NullByte":
			cache[index] = new(sql.NullByte)
			break
		case "sql.NullInt16":
			cache[index] = new(sql.NullInt16)
			break
		case "sql.NullInt32":
			cache[index] = new(sql.NullInt32)
			break
		case "sql.NullInt64":
			cache[index] = new(sql.NullInt64)
			break
		case "sql.NullFloat64":
			cache[index] = new(sql.NullFloat64)
			break
		case "sql.NullTime":
			cache[index] = new(sql.NullTime)
			break
		case "sql.RawBytes":
			cache[index] = new(sql.RawBytes)
			break
		default:
			cache[index] = new(interface{})
			//panic("GetSqlValueCache type [" + columnType.ScanType().String() + "] columnName [" + columnType.Name() + "] databaseType [" + columnType.DatabaseTypeName() + "] not support")
			break
		}
	}
	return
}

func GetSqlValue(columnType *sql.ColumnType, data interface{}) (value interface{}) {
	if data == nil {
		return
	}
	vOf := reflect.ValueOf(data)
	if vOf.Kind() == reflect.Ptr {
		if vOf.IsNil() {
			return nil
		}
		return GetSqlValue(columnType, vOf.Elem().Interface())
	}
	switch v := data.(type) {
	case sql.NullString:
		value = (v).String
		break
	case sql.NullBool:
		value = (v).Bool
		break
	case sql.NullByte:
		value = (v).Byte
		break
	case sql.NullFloat64:
		value = (v).Float64
		break
	case sql.NullInt16:
		value = (v).Int16
		break
	case sql.NullInt32:
		value = (v).Int32
		break
	case sql.NullInt64:
		value = (v).Int64
		break
	case sql.NullTime:
		value = (v).Time
		break
	case sql.RawBytes:
		value = string(v)
		break
	case []uint8:
		value = string(v)
		break
		break
	default:
		numberV, is := dialect.GetGoDrorNumberValue(data)
		if is {
			value = numberV
			break
		}
		value = v
		//panic("GetSqlValue data [" + fmt.Sprint(data) + "] name [" + columnType.Name() + "] databaseType [" + columnType.DatabaseTypeName() + "] not support")
		break
	}
	return
}
