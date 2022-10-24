package dialect

import (
	"reflect"
	"strconv"
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

func GetGoDrorNumberValue(data interface{}) (res string, is bool) {
	if data == nil {
		return
	}
	vType := reflect.TypeOf(data)
	if vType.String() == "godror.Number" {
		is = true
		method, find := vType.MethodByName("String")
		if find {
			var args = make([]reflect.Value, 0)
			args = append(args, reflect.ValueOf(data))
			resList := method.Func.Call(args)
			if len(resList) > 0 {
				res = resList[0].Interface().(string)
			}
		}
	}
	return
}
