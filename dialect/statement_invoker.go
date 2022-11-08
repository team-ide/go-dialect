package dialect

import (
	"errors"
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
func isTrue(value string) bool {
	res, _ := strconv.ParseBool(value)
	return res
}

func (this_ *IfSqlStatement) Format(context map[string]interface{}) (text string, err error) {
	//text += this_.Content

	if this_.ConditionExpression == nil {
		err = errors.New("if statement expression is null")
		return
	}
	var invoked bool
	var checkOk string
	var oneText string
	checkOk, err = this_.ConditionExpression.Invoke(context)
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
			checkOk, err = this_.ConditionExpression.Invoke(context)
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

func (this_ *ExpressionStatement) Invoke(context map[string]interface{}) (res string, err error) {
	//text += this_.Content

	res = "true"
	return
}

func (this_ *ExpressionIdentifierStatement) Format(context map[string]interface{}) (res string, err error) {
	v, ok := context[this_.Identifier]
	if !ok {
		err = errors.New("identifier [" + this_.Identifier + "] not define")
		return
	}
	res = GetStringValue(v)
	return
}
