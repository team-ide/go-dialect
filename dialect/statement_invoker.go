package dialect

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func (this_ *AbstractSqlStatement) Format(context map[string]interface{}) (text string, err error) {
	text += this_.Content

	var oneText string
	for _, one := range this_.Children {
		oneText, err = one.Format(context)
		if err != nil {
			return
		}
		text += oneText
	}
	return
}

func isTrue(value interface{}) bool {
	res, _ := strconv.ParseBool(GetStringValue(value))
	return res
}

func (this_ *IfSqlStatement) Format(context map[string]interface{}) (text string, err error) {
	//text += this_.Content

	if this_.ConditionExpression == nil {
		err = errors.New("if statement expression is null")
		return
	}
	var invoked bool
	var checkOk interface{}
	var oneText string
	checkOk, err = this_.ConditionExpression.GetValue(context)
	if err != nil {
		return
	}
	if isTrue(checkOk) {
		invoked = true
		for _, one := range this_.Children {
			oneText, err = one.Format(context)
			if err != nil {
				return
			}
			text += oneText
		}
	}
	if !invoked {
		for _, one := range this_.ElseIfs {
			if one.ConditionExpression == nil {
				err = errors.New("else if statement expression is null")
				return
			}
			checkOk, err = this_.ConditionExpression.GetValue(context)
			if err != nil {
				return
			}
			if !isTrue(checkOk) {
				continue
			}
			invoked = true
			oneText, err = one.Format(context)
			if err != nil {
				return
			}
			text += oneText
			break
		}
		if !invoked {
			if this_.Else != nil {
				oneText, err = this_.Else.Format(context)
				if err != nil {
					return
				}
				text += oneText
			}
		}
	}
	return
}

func (this_ *ExpressionStatement) GetValue(context map[string]interface{}) (res interface{}, err error) {
	//text += this_.Content

	var data interface{}
	for _, one := range this_.Children {
		data, err = GetStatementValue(one, context)
		if err != nil {
			return
		}
	}
	res = data
	return
}

func GetStatementValue(sqlStatement SqlStatement, context map[string]interface{}) (res interface{}, err error) {

	var data interface{}

	switch statement := sqlStatement.(type) {
	case *ExpressionFuncStatement:
		data, err = statement.GetValue(context)
		if err != nil {
			return
		}
		break
	case *ExpressionIdentifierStatement:
		data, err = statement.GetValue(context)
		if err != nil {
			return
		}
		break
	case *ExpressionStringStatement:
		data, err = statement.GetValue(context)
		if err != nil {
			return
		}
		break
	case *ExpressionNumberStatement:
		data, err = statement.GetValue(context)
		if err != nil {
			return
		}
		break
	case *ExpressionBracketsStatement:
		data, err = statement.GetValue(context)
		if err != nil {
			return
		}
		break
	default:
		err = errors.New("Statement type [" + reflect.TypeOf(statement).String() + "] not support")
		return
	}
	res = data
	return
}

func (this_ *ExpressionIdentifierStatement) GetValue(context map[string]interface{}) (res interface{}, err error) {
	res, ok := context[this_.Identifier]
	if !ok {
		err = errors.New("identifier [" + this_.Identifier + "] not define")
		return
	}
	return
}

func (this_ *ExpressionStringStatement) GetValue(context map[string]interface{}) (res interface{}, err error) {
	res = this_.Value
	return
}

func (this_ *ExpressionNumberStatement) GetValue(context map[string]interface{}) (res interface{}, err error) {
	res = this_.Value
	return
}

func (this_ *ExpressionFuncStatement) GetValue(context map[string]interface{}) (res interface{}, err error) {
	find, ok := context[this_.Func]
	if !ok {
		err = errors.New("func [" + this_.Func + "] not define")
		return
	}
	method, ok := find.(reflect.Value)
	if !ok {
		err = errors.New("func [" + this_.Func + "] can not to reflect.Method")
		return
	}
	var values []reflect.Value
	var v interface{}
	for _, arg := range this_.Children {
		v, err = GetStatementValue(arg, context)
		if err != nil {
			return
		}
		values = append(values, reflect.ValueOf(v))
		fmt.Println("ExpressionFuncStatement GetValue arg:", arg, ",value:", v)
	}
	fmt.Println("ExpressionFuncStatement GetValue args:", this_.Args, ",values:", values)
	methodResults := method.Call(values)
	if len(methodResults) > 0 {
		for _, methodResult := range methodResults {
			switch obj := methodResult.Interface().(type) {
			case error:
				err = obj
				break
			default:
				res = methodResult.Interface()
				break
			}
		}
	}
	return
}

func (this_ *ExpressionBracketsStatement) GetValue(context map[string]interface{}) (res interface{}, err error) {

	return
}
