package dialect

import (
	"errors"
	"reflect"
	"strconv"
)

func FormatStatements(statements []Statement, statementContext *StatementContext) (text string, err error) {

	var oneText string
	for _, one := range statements {
		oneText, err = FormatStatement(one, statementContext)
		if err != nil {
			return
		}
		text += oneText
	}
	return
}

func FormatStatement(statement_ Statement, statementContext *StatementContext) (text string, err error) {

	switch statement := statement_.(type) {
	case *IgnorableStatement:
		text, err = statement.Format(statementContext)
		break
	default:
		text, err = statement.Format(statementContext)
		break
	}
	return
}

func (this_ *AbstractStatement) Format(statementContext *StatementContext) (text string, err error) {
	text += this_.Content

	childrenText, err := FormatStatements(this_.Children, statementContext)
	if err != nil {
		return
	}
	text += childrenText
	return
}

func (this_ *IgnorableStatement) Format(statementContext *StatementContext) (text string, err error) {

	findValue, err := StatementsFindValue(this_.Children, statementContext)
	if err != nil {
		return
	}
	if !findValue {
		return
	}
	childrenText, err := FormatStatements(this_.Children, statementContext)
	if err != nil {
		return
	}
	text += childrenText
	return
}

func isTrue(value interface{}) bool {
	res, _ := strconv.ParseBool(GetStringValue(value))
	return res
}

func (this_ *ElseIfStatement) Format(statementContext *StatementContext) (text string, err error) {
	childrenText, err := FormatStatements(this_.Children, statementContext)
	if err != nil {
		return
	}
	text += childrenText
	return
}

func (this_ *ElseStatement) Format(statementContext *StatementContext) (text string, err error) {
	childrenText, err := FormatStatements(this_.Children, statementContext)
	if err != nil {
		return
	}
	text += childrenText
	return
}

func (this_ *IfStatement) Format(statementContext *StatementContext) (text string, err error) {
	//text += this_.Content

	if this_.ConditionExpression == nil {
		err = errors.New("if statement expression is null")
		return
	}
	var invoked bool
	var checkOk interface{}
	var oneText string
	checkOk, err = this_.ConditionExpression.GetValue(statementContext)
	if err != nil {
		return
	}
	if isTrue(checkOk) {
		invoked = true
		for _, one := range this_.Children {
			oneText, err = one.Format(statementContext)
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
			checkOk, err = this_.ConditionExpression.GetValue(statementContext)
			if err != nil {
				return
			}
			if !isTrue(checkOk) {
				continue
			}
			invoked = true
			oneText, err = one.Format(statementContext)
			if err != nil {
				return
			}
			text += oneText
			break
		}
		if !invoked {
			if this_.Else != nil {
				oneText, err = this_.Else.Format(statementContext)
				if err != nil {
					return
				}
				text += oneText
			}
		}
	}
	return
}

func (this_ *ExpressionStatement) GetValue(statementContext *StatementContext) (res interface{}, err error) {
	//text += this_.Content

	var data interface{}
	for _, one := range this_.Children {
		data, err = GetStatementValue(one, statementContext)
		if err != nil {
			return
		}
	}
	res = data
	return
}

func GetStatementValue(statement_ Statement, statementContext *StatementContext) (res interface{}, err error) {

	var data interface{}

	switch statement := statement_.(type) {
	case *ExpressionFuncStatement:
		data, err = statement.GetValue(statementContext)
		break
	case *ExpressionIdentifierStatement:
		data, err = statement.GetValue(statementContext)
		break
	case *ExpressionStringStatement:
		data, err = statement.GetValue(statementContext)
		break
	case *ExpressionNumberStatement:
		data, err = statement.GetValue(statementContext)
		break
	case *ExpressionBracketsStatement:
		data, err = statement.GetValue(statementContext)
		break
	default:
		err = errors.New("Statement type [" + reflect.TypeOf(statement).String() + "] not support")
		return
	}
	if err != nil {
		return
	}
	res = data
	return
}

func StatementsFindValue(statements []Statement, statementContext *StatementContext) (findValue bool, err error) {
	var data interface{}
	for _, one := range statements {
		switch statement := one.(type) {
		//case *ExpressionFuncStatement:
		//	data, err = statement.GetValue(statementContext)
		//	break
		case *ExpressionIdentifierStatement:
			data, err = statement.GetValue(statementContext)
			break
		case *ExpressionStringStatement:
			data, err = statement.GetValue(statementContext)
			break
		case *ExpressionNumberStatement:
			data, err = statement.GetValue(statementContext)
			break
		case *ExpressionBracketsStatement:
			data, err = statement.GetValue(statementContext)
			break
		}
		if err != nil {
			return
		}
		if data != nil && data != "" {
			findValue = true
			break
		}
		findValue, err = StatementsFindValue(*one.GetChildren(), statementContext)

		if err != nil {
			return
		}
		if findValue {
			break
		}
	}

	return
}

func (this_ *ExpressionIdentifierStatement) Format(statementContext *StatementContext) (text string, err error) {
	value, err := this_.GetValue(statementContext)
	if err != nil {
		return
	}
	text = GetStringValue(value)
	return
}

func (this_ *ExpressionIdentifierStatement) GetValue(statementContext *StatementContext) (res interface{}, err error) {
	res, ok := statementContext.GetData(this_.Identifier)
	if !ok {
		err = errors.New("identifier [" + this_.Identifier + "] not define")
		return
	}
	return
}

func (this_ *ExpressionStringStatement) Format(statementContext *StatementContext) (text string, err error) {
	value, err := this_.GetValue(statementContext)
	if err != nil {
		return
	}
	text = GetStringValue(value)
	return
}

func (this_ *ExpressionStringStatement) GetValue(statementContext *StatementContext) (res interface{}, err error) {
	res = this_.Value
	return
}

func (this_ *ExpressionNumberStatement) Format(statementContext *StatementContext) (text string, err error) {
	value, err := this_.GetValue(statementContext)
	if err != nil {
		return
	}
	text = GetStringValue(value)
	return
}

func (this_ *ExpressionNumberStatement) GetValue(statementContext *StatementContext) (res interface{}, err error) {
	res = this_.Value
	return
}

func (this_ *ExpressionFuncStatement) Format(statementContext *StatementContext) (text string, err error) {
	value, err := this_.GetValue(statementContext)
	if err != nil {
		return
	}
	text = GetStringValue(value)
	return
}

func (this_ *ExpressionFuncStatement) GetValue(statementContext *StatementContext) (res interface{}, err error) {
	method, ok := statementContext.GetMethod(this_.Func)
	if !ok {
		err = errors.New("func [" + this_.Func + "] not define")
		return
	}
	var values []interface{}
	var v interface{}
	for _, arg := range this_.Children {
		v, err = GetStatementValue(arg, statementContext)
		if err != nil {
			return
		}
		values = append(values, v)
	}
	methodResults, err := method.Call(values)
	if err != nil {
		return
	}
	if len(methodResults) > 0 {
		res = methodResults[0]
	}
	return
}

func (this_ *ExpressionBracketsStatement) GetValue(statementContext *StatementContext) (res interface{}, err error) {

	return
}
