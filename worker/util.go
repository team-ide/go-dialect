package worker

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/team-ide/go-dialect/dialect"
	"os"
	"reflect"
	"strings"
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

func SplitSqlList(sqlInfo string) (sqlList []string) {
	var list []string
	var beg int

	var inStringLevel int
	var inStringPack byte
	var thisChar byte
	var lastChar byte

	var stringPackChars = []byte{'"', '\''}
	for i := 0; i < len(sqlInfo); i++ {
		thisChar = sqlInfo[i]
		if i > 0 {
			lastChar = sqlInfo[i-1]
		}

		// inStringLevel == 0 表示 不在 字符串 包装 中
		if thisChar == ';' && inStringLevel == 0 {
			if i > 0 {
				list = append(list, sqlInfo[beg:i])
			}
			beg = i + 1
		} else {
			packCharIndex := dialect.BytesIndex(stringPackChars, thisChar)
			if packCharIndex >= 0 {
				// inStringLevel == 0 表示 不在 字符串 包装 中
				if inStringLevel == 0 {
					inStringPack = stringPackChars[packCharIndex]
					// 字符串包装层级 +1
					inStringLevel++
				} else {
					// 如果有转义符号 类似 “\'”，“\"”
					if lastChar == '\\' {
					} else if lastChar == inStringPack {
						// 如果 前一个字符 与字符串包装字符一致
						inStringLevel--
					} else {
						// 字符串包装层级 -1
						inStringLevel--
					}
				}
			}
		}

	}
	list = append(list, sqlInfo[beg:])
	for _, sqlOne := range list {
		sqlOne = strings.TrimSpace(sqlOne)
		if sqlOne == "" {
			continue
		}
		sqlList = append(sqlList, sqlOne)
	}
	return
}

//NowTime 获取当前时间戳
func NowTime() int64 {
	return GetTime(Now())
}

//GetTime 获取当前时间戳
func GetTime(time time.Time) int64 {
	return time.UnixNano() / 1e6
}

//Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

// UUID 生成UUID
func UUID() (res string) {
	res = uuid.NewString()
	res = strings.ReplaceAll(res, "-", "")
	return
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
	default:
		baseValue, isBaseType := dialect.GetBaseTypeValue(value)
		if isBaseType {
			value = baseValue
			break
		}
		value = v
		//panic("GetSqlValue data [" + fmt.Sprint(data) + "] name [" + columnType.Name() + "] databaseType [" + columnType.DatabaseTypeName() + "] not support")
		break
	}
	return
}

//SplitArrayMap 分割数组，根据传入的数组和分割大小，将数组分割为大小等于指定大小的多个数组，如果不够分，则最后一个数组元素小于其他数组
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
