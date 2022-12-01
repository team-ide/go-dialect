package worker

import (
	"database/sql"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"os"
	"reflect"
	"time"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathIsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// NowTime 获取当前时间戳
func NowTime() int64 {
	return GetTime(Now())
}

// GetTime 获取当前时间戳
func GetTime(time time.Time) int64 {
	return time.UnixNano() / 1e6
}

// Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

func GetSqlValueCache(columnTypes []*sql.ColumnType) (cache []interface{}) {
	cache = make([]interface{}, len(columnTypes)) //临时存储每行数据
	for index, _ := range cache {
		cache[index] = new(interface{})
		//columnType := columnTypes[index]
		//ct := columnType.ScanType()
		//if ct == nil {
		//	cache[index] = new(sql.NullString)
		//	continue
		//}
		//println("GetSqlValueCache type [" + columnType.ScanType().String() + "] columnName [" + columnType.Name() + "] databaseType [" + columnType.DatabaseTypeName() + "]")
		//switch ct.String() {
		//case "sql.NullString":
		//	cache[index] = new(sql.NullString)
		//	break
		//case "sql.NullBool":
		//	cache[index] = new(sql.NullBool)
		//	break
		//case "sql.NullByte":
		//	cache[index] = new(sql.NullByte)
		//	break
		//case "sql.NullInt16":
		//	cache[index] = new(sql.NullInt16)
		//	break
		//case "sql.NullInt32":
		//	cache[index] = new(sql.NullInt32)
		//	break
		//case "sql.NullInt64":
		//	cache[index] = new(sql.NullInt64)
		//	break
		//case "sql.NullFloat64":
		//	cache[index] = new(sql.NullFloat64)
		//	break
		//case "sql.NullTime":
		//	cache[index] = new(sql.NullTime)
		//	break
		//case "sql.RawBytes":
		//	cache[index] = new(sql.RawBytes)
		//	break
		//default:
		//	cache[index] = new(interface{})
		//	//panic("GetSqlValueCache type [" + columnType.ScanType().String() + "] columnName [" + columnType.Name() + "] databaseType [" + columnType.DatabaseTypeName() + "] not support")
		//	break
		//}
	}
	return
}

func GetSqlValue(columnType *sql.ColumnType, data interface{}) (value interface{}) {
	if data == nil {
		return
	}
	typeName := reflect.TypeOf(data).String()
	if typeName == "*dm.DmClob" {
		typeV := reflect.ValueOf(data)
		method := typeV.MethodByName("GetLength")
		vs := method.Call([]reflect.Value{})
		if vs[1].Interface() == nil {
			length := vs[0].Int()
			method = typeV.MethodByName("ReadString")
			vs = method.Call([]reflect.Value{reflect.ValueOf(0), reflect.ValueOf(int(length))})
			value = vs[0].String()
		}
		return
	} else if typeName == "godror.Number" {
		typeV := reflect.ValueOf(data)
		method := typeV.MethodByName("String")
		vs := method.Call([]reflect.Value{})
		value = vs[0].String()
		return
	}
	vOf := reflect.ValueOf(data)
	if vOf.Kind() == reflect.Ptr {
		if vOf.IsNil() {
			return nil
		}
		return GetSqlValue(columnType, vOf.Elem().Interface())
	}
	//if columnType.Name() == "NESTING_EVENT_TYPE" {
	//	fmt.Println("NESTING_EVENT_TYPE value type", reflect.TypeOf(data).String(), " value is ", data)
	//}
	switch v := data.(type) {
	case sql.NullString:
		if !v.Valid {
			return nil
		}
		value = (v).String
		break
	case sql.NullBool:
		if !v.Valid {
			return nil
		}
		value = (v).Bool
		break
	case sql.NullByte:
		if !v.Valid {
			return nil
		}
		value = (v).Byte
		break
	case sql.NullFloat64:
		if !v.Valid {
			return nil
		}
		value = (v).Float64
		break
	case sql.NullInt16:
		if !v.Valid {
			return nil
		}
		value = (v).Int16
		break
	case sql.NullInt32:
		if !v.Valid {
			return nil
		}
		value = (v).Int32
		break
	case sql.NullInt64:
		if !v.Valid {
			return nil
		}
		value = (v).Int64
		break
	case sql.NullTime:
		if !v.Valid {
			return nil
		}
		value = (v).Time
		break
	case sql.RawBytes:
		value = string(v)
		break
	case []uint8:
		value = string(v)
		break
	case string, int, int8, int16, int32, int64, float32, float64, bool, uint, uint8, uint16, uint32, uint64:
		value = v
		break
	case time.Time:
		value = v
		break
	default:
		baseValue, isBaseType := dialect.GetBaseTypeValue(value)
		if isBaseType {
			value = baseValue
			return
		}
		value = v
		panic("GetSqlValue data [" + fmt.Sprint(data) + "] data type [" + reflect.TypeOf(data).String() + "] name [" + columnType.Name() + "] databaseType [" + columnType.DatabaseTypeName() + "] not support")
		break
	}
	return
}

// SplitArrayMap 分割数组，根据传入的数组和分割大小，将数组分割为大小等于指定大小的多个数组，如果不够分，则最后一个数组元素小于其他数组
func SplitArrayMap(arr []map[string]interface{}, num int) [][]map[string]interface{} {
	max := len(arr)
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		if max == 0 {
			return [][]map[string]interface{}{}
		}
		return [][]map[string]interface{}{arr}
	}
	//获取应该数组分割为多少份
	var quantity int
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]map[string]interface{}, 0)
	//声明分割数组的截止下标
	var start, end, i int
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}
