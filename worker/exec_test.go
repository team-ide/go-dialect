package worker

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-driver/db_mysql"
	"reflect"
	"testing"
	"time"
)

func TestDoQuery(t *testing.T) {
	db, err := db_mysql.Open(db_mysql.GetDSN("root", "123456", "127.0.0.1", 3306, ""))
	if err != nil {
		panic(err)
	}
	var list []*QueryStruct
	err = DoQueryStructs(db, `select user as a from mysql.user`, []interface{}{}, &list)
	if err != nil {
		panic(err)
	}
	for _, one := range list {
		bs, _ := json.Marshal(one)
		fmt.Println(string(bs))
	}
}
func TestDoQueryOne(t *testing.T) {
	db, err := db_mysql.Open(db_mysql.GetDSN("root", "123456", "127.0.0.1", 3306, ""))
	if err != nil {
		panic(err)
	}
	one := &QueryStruct{}
	_, err = DoQueryStruct(db, `select user as a,1 b,0 c,now() deleteTime from mysql.user where user='mysql.sys'`, []interface{}{}, one)
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(one)
	fmt.Println(string(bs))

	data, err := DoQueryOne(db, `select user as a from mysql.user where user='mysql.sys'`, []interface{}{})
	if err != nil {
		panic(err)
	}
	bs, _ = json.Marshal(data)
	fmt.Println(string(bs))
}
func TestDoQueryPage(t *testing.T) {
	page := NewPage()
	page.PageNo = 3
	page.PageSize = 3
	db, err := db_mysql.Open(db_mysql.GetDSN("root", "123456", "127.0.0.1", 3306, ""))
	if err != nil {
		panic(err)
	}
	dia, err := dialect.NewDialect("mysql")
	list, err := DoQueryPage(db, dia, `select user as a from mysql.user `, []interface{}{}, page)
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(page)
	fmt.Println(string(bs))
	for _, one := range list {
		bs, _ = json.Marshal(one)
		fmt.Println(string(bs))
	}
	var dataList []*QueryStruct

	err = DoQueryPageStructs(db, dia, `select user as a from mysql.user `, []interface{}{}, page, &dataList)
	if err != nil {
		panic(err)
	}
	bs, _ = json.Marshal(page)
	fmt.Println(string(bs))
	for _, one := range dataList {
		bs, _ = json.Marshal(one)
		fmt.Println(string(bs))
	}
}

func TestBean(t *testing.T) {
	//var list []*QueryBean
	//
	//listBeanType := GetListBeanType(list)
	//listBeanValue := reflect.New(listBeanType)
	//listBeanValue.Elem().FieldByName("A").Set(reflect.ValueOf("xxx"))
	//newBean := listBeanValue.Interface() // 调用反射创建对象
	//fmt.Println(listBeanType.String())
	//fmt.Println(listBeanValue)
	//fmt.Println(newBean)
	//list = append(list, newBean.(*QueryBean))
	//fmt.Println(list)

	data := &QueryStruct{}
	reflect.ValueOf(data).Elem().FieldByName("DeleteTime").Set(reflect.ValueOf(time.Time{}))
	bs, _ := json.Marshal(data)
	fmt.Println(string(bs))
}

type QueryStruct struct {
	A          string    `json:"a"`
	B          int8      `json:"b"`
	C          int8      `json:"c"`
	DeleteTime time.Time `json:"deleteTime,omitempty"`
}
