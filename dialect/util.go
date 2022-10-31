package dialect

import (
	"database/sql"
	"github.com/google/uuid"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StringToInt(str string) (res int, err error) {
	i64, err := StringToInt64(str)
	if err != nil {
		return
	}
	res = int(i64)
	return
}

func StringToInt64(str string) (res int64, err error) {
	if str == "null" {
		return
	}
	res, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	return
}

// StringsIndex Returns the index position of the string val in array
func StringsIndex(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// BytesIndex Returns the index position of the string val in array
func BytesIndex(array []byte, val byte) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

func ReplaceStringByRegex(str, rule, replace string) string {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return str
	}
	return reg.ReplaceAllString(str, replace)
}

func GetStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	vOf := reflect.ValueOf(value)
	if vOf.Kind() == reflect.Ptr {
		if vOf.IsNil() {
			return ""
		}
		return GetStringValue(vOf.Elem().Interface())
	}
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
			return ""
		}
		return v.Format("2006-01-02 15:04:05")
	case string:
		return v
	case []byte:
		return string(v)
	case sql.NullString:
		return v.String
	default:
		baseValue, isBaseType := GetBaseTypeValue(value)
		if isBaseType {
			return GetStringValue(baseValue)
		}
		panic("value type [" + reflect.TypeOf(value).String() + "] not support")
	}
	return ""
}

func GetBaseTypeValue(data interface{}) (res interface{}, is bool) {
	if data == nil {
		return
	}
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		if dataValue.IsNil() {
			return
		}
		return GetBaseTypeValue(dataValue.Elem().Interface())
	}

	is = true
	kindName := reflect.TypeOf(data).Kind().String()
	switch kindName {
	case "string":
		res = dataValue.String()
		break
	case "int":
		res = int(dataValue.Int())
		break
	case "int8":
		res = int8(dataValue.Int())
		break
	case "int16":
		res = int16(dataValue.Int())
		break
	case "int32":
		res = int32(dataValue.Int())
		break
	case "int64":
		res = dataValue.Int()
		break
	case "float32":
		res = float32(dataValue.Float())
		break
	case "float64":
		res = dataValue.Float()
		break
	case "bool":
		res = dataValue.Bool()
		break
	case "[]bytes":
		res = dataValue.Bytes()
		break
	case "uint":
		res = uint(dataValue.Uint())
		break
	case "uint8":
		res = uint8(dataValue.Uint())
		break
	case "uint16":
		res = uint16(dataValue.Uint())
		break
	case "uint32":
		res = uint32(dataValue.Uint())
		break
	case "uint64":
		res = dataValue.Uint()
		break
	default:
		is = false
		res = dataValue.Interface()
		break
	}
	return
}

// UUID 生成UUID
func UUID() (res string) {
	res = uuid.NewString()
	res = strings.ReplaceAll(res, "-", "")
	return
}
