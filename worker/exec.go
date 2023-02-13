package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func DoExec(db *sql.DB, sqlInfo string, args []interface{}) (result sql.Result, err error) {
	if len(sqlInfo) == 0 {
		return
	}
	resultList, _, _, err := DoExecs(db, []string{sqlInfo}, [][]interface{}{args})
	if err != nil {
		return
	}
	if len(resultList) > 0 {
		result = resultList[0]
	}
	return
}

func DoOwnerExecs(dia dialect.Dialect, db *sql.DB, ownerName string, sqlList []string, argsList [][]interface{}) (resultList []sql.Result, errSql string, errArgs []interface{}, err error) {
	sqlListSize := len(sqlList)
	if sqlListSize == 0 {
		return
	}
	if len(argsList) == 0 {
		argsList = make([][]interface{}, sqlListSize)
	}
	argsListSize := len(argsList)
	if sqlListSize != argsListSize {
		err = errors.New(fmt.Sprintf("sqlList size is [%d] but argsList size is [%d]", sqlListSize, argsListSize))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if ownerName != "" {
		switch dia.DialectType() {
		case dialect.TypeMysql:
			_, _ = tx.Exec(" USE " + ownerName)
			break
		case dialect.TypeOracle:
			_, _ = tx.Exec("ALTER SESSION SET CURRENT_SCHEMA=" + ownerName)
			break
		case dialect.TypeGBase:
			_, _ = tx.Exec("database " + ownerName)
			break
		}
	}
	var result sql.Result
	for i := 0; i < sqlListSize; i++ {
		sqlInfo := sqlList[i]
		args := argsList[i]
		if strings.TrimSpace(sqlInfo) == "" {
			continue
		}
		result, err = tx.Exec(sqlInfo, args...)
		if err != nil {
			errSql = sqlInfo
			errArgs = args
			return
		}
		resultList = append(resultList, result)
	}

	return
}

func DoExecs(db *sql.DB, sqlList []string, argsList [][]interface{}) (resultList []sql.Result, errSql string, errArgs []interface{}, err error) {
	sqlListSize := len(sqlList)
	if sqlListSize == 0 {
		return
	}
	if len(argsList) == 0 {
		argsList = make([][]interface{}, sqlListSize)
	}
	argsListSize := len(argsList)
	if sqlListSize != argsListSize {
		err = errors.New(fmt.Sprintf("sqlList size is [%d] but argsList size is [%d]", sqlListSize, argsListSize))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	var result sql.Result
	for i := 0; i < sqlListSize; i++ {
		sqlInfo := sqlList[i]
		args := argsList[i]
		if strings.TrimSpace(sqlInfo) == "" {
			continue
		}
		result, err = tx.Exec(sqlInfo, args...)
		if err != nil {
			errSql = sqlInfo
			errArgs = args
			return
		}
		resultList = append(resultList, result)
	}

	return
}

func DoQuery(db *sql.DB, sqlInfo string, args []interface{}) (list []map[string]interface{}, err error) {
	_, _, list, err = DoQueryWithColumnTypes(db, sqlInfo, args)
	if err != nil {
		return
	}
	return
}

func DoQueryOne(db *sql.DB, sqlInfo string, args []interface{}) (data map[string]interface{}, err error) {
	_, _, list, err := DoQueryWithColumnTypes(db, sqlInfo, args)
	if err != nil {
		return
	}
	if len(list) > 0 {
		data = list[0]
		if len(list) > 1 {
			err = errors.New("has more rows by query one")
			return
		}
	}
	return
}

func DoQueryStructs(db *sql.DB, sqlInfo string, args []interface{}, list interface{}) (err error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	listVOf := reflect.ValueOf(list).Elem()
	listStrType := GetListStructType(list)
	for rows.Next() {
		var values []interface{}
		for range columnTypes {
			values = append(values, new(interface{}))
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}

		item := make(map[string]interface{})
		for index, data := range values {
			item[columnTypes[index].Name()] = GetSqlValue(columnTypes[index], data)
		}
		listStrValue := reflect.New(listStrType)
		SetStructColumnValues(item, listStrValue.Elem())
		listVOf = reflect.Append(listVOf, listStrValue)
	}
	reflect.ValueOf(list).Elem().Set(listVOf)
	return
}

func DoQueryStruct(db *sql.DB, sqlInfo string, args []interface{}, str interface{}) (find bool, err error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	strVOf := reflect.ValueOf(str)

	var isBase bool
	switch str.(type) {
	case *int, *int8, *int16, *int32, *int64, *float32, *float64:
		isBase = true
		break
	}
	for rows.Next() {
		if find {
			err = errors.New("has more rows by query one")
			return
		}
		find = true
		var values []interface{}
		if isBase {
			values = []interface{}{str}
		} else {
			for range columnTypes {
				values = append(values, new(interface{}))
			}
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}
		if isBase {
			continue
		}
		item := make(map[string]interface{})
		for index, data := range values {
			item[columnTypes[index].Name()] = GetSqlValue(columnTypes[index], data)
		}
		SetStructColumnValues(item, strVOf.Elem())
	}
	return
}
func DoQueryWithColumnTypes(db *sql.DB, sqlInfo string, args []interface{}) (columns []string, columnTypes []*sql.ColumnType, list []map[string]interface{}, err error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	columnTypes, err = rows.ColumnTypes()
	if err != nil {
		return
	}
	for rows.Next() {
		var values []interface{}
		for range columnTypes {
			values = append(values, new(interface{}))
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}
		item := make(map[string]interface{})
		for index, data := range values {
			item[columns[index]] = GetSqlValue(columnTypes[index], data)
		}
		list = append(list, item)
	}

	return
}

var (
	structFieldMapCache  = map[reflect.Type]map[string]reflect.StructField{}
	structColumnMapCache = map[reflect.Type]map[string]reflect.StructField{}
	structMapLock        sync.Mutex
)

func getStructColumn(tOf reflect.Type) (structFieldMap map[string]reflect.StructField, structColumnMap map[string]reflect.StructField) {
	structMapLock.Lock()
	defer structMapLock.Unlock()
	structFieldMap, ok := structFieldMapCache[tOf]
	structColumnMap = structColumnMapCache[tOf]
	if ok {
		//fmt.Println("find from cache")
		return
	}
	structFieldMap = map[string]reflect.StructField{}
	structColumnMap = map[string]reflect.StructField{}
	for i := 0; i < tOf.NumField(); i++ {
		field := tOf.Field(i)
		structFieldMap[field.Name] = field
		str := field.Tag.Get("column")
		if str != "" && str != "-" {
			ss := strings.Split(str, ",")
			structColumnMap[ss[0]] = field
		} else {
			str = field.Tag.Get("json")
			if str != "" && str != "-" {
				ss := strings.Split(str, ",")
				structColumnMap[ss[0]] = field
			}
		}
	}
	structFieldMapCache[tOf] = structFieldMap
	structColumnMapCache[tOf] = structColumnMap
	return
}
func SetStructColumnValues(columnValueMap map[string]interface{}, strValue reflect.Value) {
	if len(columnValueMap) == 0 {
		return
	}
	tOf := strValue.Type()

	_, structColumnMap := getStructColumn(tOf)

	for columnName, columnValue := range columnValueMap {
		field, find := structColumnMap[columnName]
		if !find {
			field, find = structColumnMap[columnName]
		}
		if !find {
			continue
		}
		valueTypeOf := reflect.TypeOf(columnValue)
		columnValueType := ""
		fieldType := field.Type.String()
		if valueTypeOf != nil {
			columnValueType = valueTypeOf.String()
		}
		if columnValueType != fieldType {
			switch fieldType {
			case "string":
				columnValue = dialect.GetStringValue(columnValue)
				break
			case "int8", "int16", "int32", "int64", "int":
				str := dialect.GetStringValue(columnValue)
				var num int64
				if str != "" {
					num, _ = dialect.StringToInt64(str)
				}
				if fieldType == "int8" {
					columnValue = int8(num)
				} else if fieldType == "int16" {
					columnValue = int16(num)
				} else if fieldType == "int32" {
					columnValue = int32(num)
				} else if fieldType == "int64" {
					columnValue = num
				} else if fieldType == "int" {
					columnValue = int(num)
				}
				break
			case "float32", "float64":
				str := dialect.GetStringValue(columnValue)
				var num float64
				if str != "" {
					num, _ = strconv.ParseFloat(str, 64)
				}
				if fieldType == "float32" {
					columnValue = float32(num)
				} else if fieldType == "float64" {
					columnValue = num
				}
				break
			case "time.Time":
				if columnValue == nil || columnValue == 0 {
					columnValue = time.Time{}
					break
				}
				valueOf := reflect.ValueOf(columnValue)
				if valueOf.IsNil() || valueOf.IsZero() {
					columnValue = time.Time{}
				}
				break
			}
		}

		valueOf := reflect.ValueOf(columnValue)
		strValue.FieldByName(field.Name).Set(valueOf)
	}
	return
}

func GetListStructType(list interface{}) reflect.Type {
	vOf := reflect.ValueOf(list)
	if vOf.Kind() == reflect.Ptr {
		return GetListStructType(vOf.Elem().Interface())
	}
	tOf := reflect.TypeOf(list).Elem()
	if tOf.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
		tOf = tOf.Elem()
	}
	return tOf
}

func DoQueryCount(db *sql.DB, sqlInfo string, args []interface{}) (count int, err error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}

	return
}

func DoQueryPage(db *sql.DB, dia dialect.Dialect, sqlInfo string, args []interface{}, page *Page) (list []map[string]interface{}, err error) {
	if page.PageSize < 1 {
		page.PageSize = 1
	}
	if page.PageNo < 1 {
		page.PageNo = 1
	}
	pageSize := page.PageSize
	pageNo := page.PageNo

	countSql, err := dialect.FormatCountSql(sqlInfo)
	if err != nil {
		return
	}
	page.TotalCount, err = DoQueryCount(db, countSql, args)
	if err != nil {
		return
	}
	page.TotalPage = (page.TotalCount + page.PageSize - 1) / page.PageSize
	// 如果查询的页码 大于 总页码 则不查询
	if pageNo > page.TotalPage {
		return
	}
	pageSql := dia.PackPageSql(sqlInfo, pageSize, pageNo)

	list, err = DoQuery(db, pageSql, args)
	if err != nil {
		return
	}

	return
}

func DoQueryPageStructs(db *sql.DB, dia dialect.Dialect, sqlInfo string, args []interface{}, page *Page, list interface{}) (err error) {
	if page.PageSize < 1 {
		page.PageSize = 1
	}
	if page.PageNo < 1 {
		page.PageNo = 1
	}
	pageSize := page.PageSize
	pageNo := page.PageNo

	countSql, err := dialect.FormatCountSql(sqlInfo)
	if err != nil {
		return
	}
	page.TotalCount, err = DoQueryCount(db, countSql, args)
	if err != nil {
		return
	}
	page.TotalPage = (page.TotalCount + page.PageSize - 1) / page.PageSize
	// 如果查询的页码 大于 总页码 则不查询
	if pageNo > page.TotalPage {
		return
	}
	pageSql := dia.PackPageSql(sqlInfo, pageSize, pageNo)

	err = DoQueryStructs(db, pageSql, args, list)
	if err != nil {
		return
	}

	return
}

type Page struct {
	PageSize   int `json:"pageSize"`
	PageNo     int `json:"pageNo"`
	TotalCount int `json:"totalCount"`
	TotalPage  int `json:"totalPage"`
}

func NewPage() *Page {
	return &Page{
		PageSize: 1,
		PageNo:   1,
	}
}
