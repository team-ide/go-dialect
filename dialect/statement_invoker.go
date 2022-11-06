package dialect

func (this_ *AbstractSqlStatement) Invoke(context map[string]interface{}) (text string, err error) {
	text += this_.Content

	var oneText string
	for _, one := range this_.Children {
		oneText, err = one.Invoke(context)
		if err != nil {
			return
		}
		text += oneText
	}
	return
}
