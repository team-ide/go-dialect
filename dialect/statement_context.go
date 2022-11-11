package dialect

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

func NewStatementContext() (res *StatementContext) {
	res = &StatementContext{
		dataCache:   make(map[string]interface{}),
		methodCache: make(map[string]*MethodInfo),
	}
	return
}

type MethodInfo struct {
	name          string
	method        reflect.Value
	methodType    reflect.Type
	ins           []reflect.Type
	outs          []reflect.Type
	outErrorIndex *int
}

func (this_ *MethodInfo) Call(inValues []interface{}) (outValues []interface{}, err error) {
	if len(inValues) != len(this_.ins) {
		err = errors.New(fmt.Sprintf("func [%s] ins len is [%d] by inValues len is [%d]", this_.name, len(this_.ins), len(inValues)))
		return
	}
	var callValues []reflect.Value
	for i, inValue := range inValues {
		in := this_.ins[i]
		inValueType := reflect.TypeOf(inValue)
		if in.Kind() != inValueType.Kind() {
			// TODO

		}
		callValue := reflect.ValueOf(inValue)
		callValues = append(callValues, callValue)
	}

	callResults := this_.method.Call(callValues)
	for i, callResult := range callResults {
		data := callResult.Interface()
		outValues = append(outValues, data)
		if this_.outErrorIndex != nil && *this_.outErrorIndex == i {
			dataErr, ok := data.(error)
			if ok {
				err = dataErr
			}
		}
	}

	return
}

type StatementContext struct {
	dataCache       map[string]interface{}
	dataCacheLock   sync.Mutex
	methodCache     map[string]*MethodInfo
	methodCacheLock sync.Mutex
}

func (this_ *StatementContext) GetData(name string) (value interface{}, find bool) {
	this_.dataCacheLock.Lock()
	defer this_.dataCacheLock.Unlock()

	value, find = this_.dataCache[name]
	//fmt.Println("GetData name:", name, ",value:", value)
	return
}

func (this_ *StatementContext) SetData(name string, value interface{}) *StatementContext {
	this_.dataCacheLock.Lock()
	defer this_.dataCacheLock.Unlock()

	this_.dataCache[name] = value
	return this_
}

func (this_ *StatementContext) SetJSONData(data interface{}) (err error) {
	if data == nil {
		return
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	if len(bs) == 0 {
		return
	}
	dataMap := map[string]interface{}{}
	err = json.Unmarshal(bs, &dataMap)
	if err != nil {
		return
	}
	return
}

func (this_ *StatementContext) SetMapData(data map[string]interface{}) *StatementContext {
	if data == nil {
		return this_
	}
	this_.dataCacheLock.Lock()
	defer this_.dataCacheLock.Unlock()

	for name, value := range data {
		this_.dataCache[name] = value
	}

	return this_
}

func (this_ *StatementContext) SetDataIfAbsent(name string, value interface{}) *StatementContext {
	this_.dataCacheLock.Lock()
	defer this_.dataCacheLock.Unlock()

	if this_.dataCache[name] == nil {
		this_.dataCache[name] = value
	}
	return this_
}

func (this_ *StatementContext) GetMethod(name string) (method *MethodInfo, find bool) {
	this_.methodCacheLock.Lock()
	defer this_.methodCacheLock.Unlock()

	method, find = this_.methodCache[name]
	return
}

func (this_ *StatementContext) setMethod(name string, method *MethodInfo) *StatementContext {
	this_.methodCacheLock.Lock()
	defer this_.methodCacheLock.Unlock()

	this_.methodCache[name] = method
	return this_
}

func (this_ *StatementContext) AddMethod(name string, methodFunc interface{}) *StatementContext {
	method := reflect.ValueOf(methodFunc)
	info := &MethodInfo{
		name:       name,
		method:     method,
		methodType: method.Type(),
	}
	for i := 0; i < info.methodType.NumIn(); i++ {
		info.ins = append(info.ins, info.methodType.In(i))
	}
	for i := 0; i < info.methodType.NumOut(); i++ {
		out := info.methodType.Out(i)
		switch out.String() {
		case "error":
			var i_ = new(int)
			*i_ = i
			info.outErrorIndex = i_
			break
		}
		info.outs = append(info.outs, out)
	}
	this_.setMethod(name, info)
	return this_
}
