package dialect

import "strconv"

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
