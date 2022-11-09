package dialect

type SqlStatement interface {
	GetTemplate() (template string)
	GetContent() (content *string)
	GetParent() (parent SqlStatement)
	GetChildren() (children *[]SqlStatement)
	Format(context map[string]interface{}) (text string, err error)
}

type AbstractSqlStatement struct {
	Content  string         `json:"content,omitempty"`
	Children []SqlStatement `json:"children,omitempty"`
	Parent   SqlStatement   `json:"-"`
}

func (this_ *AbstractSqlStatement) GetContent() (content *string) {
	content = &this_.Content
	return
}

func (this_ *AbstractSqlStatement) GetParent() (parent SqlStatement) {
	parent = this_.Parent
	return
}
func (this_ *AbstractSqlStatement) GetChildren() (children *[]SqlStatement) {
	children = &this_.Children
	return
}
func (this_ *AbstractSqlStatement) GetTemplate() (template string) {
	template += this_.Content
	for _, one := range this_.Children {
		template += one.GetTemplate()
	}
	return
}

type RootSqlStatement struct {
	*AbstractSqlStatement
}

type IgnorableSqlStatement struct {
	*AbstractSqlStatement
}

func (this_ *IgnorableSqlStatement) GetTemplate() (template string) {
	template += "["

	template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	template += "]"
	return
}

type TextSqlStatement struct {
	*AbstractSqlStatement
}

type IfSqlStatement struct {
	*AbstractSqlStatement
	Condition           string                `json:"condition"`
	ConditionExpression *ExpressionStatement  `json:"conditionExpression"`
	ElseIfs             []*ElseIfSqlStatement `json:"elseIfs"`
	Else                *ElseSqlStatement     `json:"else"`
}

func (this_ *IfSqlStatement) GetTemplate() (template string) {
	template += "if " + this_.Condition + " {\n"

	//template += this_.Content
	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	if len(this_.ElseIfs) > 0 || this_.Else != nil {
		template += "\n}"
	} else {
		template += "\n}\n"
	}
	for _, one := range this_.ElseIfs {
		template += one.GetTemplate()
	}

	if this_.Else != nil {
		template += this_.Else.GetTemplate()
	}

	return
}

type ElseIfSqlStatement struct {
	*AbstractSqlStatement
	Condition           string               `json:"condition"`
	ConditionExpression *ExpressionStatement `json:"conditionExpression"`
	If                  *IfSqlStatement      `json:"-"`
	Index               int                  `json:"index"`
}

func (this_ *ElseIfSqlStatement) GetTemplate() (template string) {
	template += " else if " + this_.Condition + " {\n"

	//template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	if this_.IsEndElseIf() {
		if this_.If.Else != nil {
			template += "\n}"
		} else {
			template += "\n}\n"
		}

	} else {
		template += "\n}"
	}
	return
}
func (this_ *ElseIfSqlStatement) IsEndElseIf() (isEnd bool) {
	if this_.Index == len(this_.If.ElseIfs)-1 {
		isEnd = true
	}
	return
}

type ElseSqlStatement struct {
	*AbstractSqlStatement
	If *IfSqlStatement `json:"-"`
}

func (this_ *ElseSqlStatement) GetTemplate() (template string) {
	template += " else {\n"

	//template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	template += "\n}\n"
	return
}

type ForStatement struct {
	*AbstractSqlStatement
}

func (this_ *ForStatement) GetTemplate() (template string) {
	template += "for " + this_.Content + " {\n"

	//template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	template += "\n}\n"
	return
}

type ExpressionStatement struct {
	*AbstractSqlStatement
}

type ExpressionStringStatement struct {
	*AbstractSqlStatement
	Value string `json:"value"`
}

type ExpressionNumberStatement struct {
	*AbstractSqlStatement
	Value float64 `json:"value"`
}

type ExpressionIdentifierStatement struct {
	*AbstractSqlStatement
	Identifier string `json:"identifier"`
}

type ExpressionFuncStatement struct {
	*AbstractSqlStatement
	Func string         `json:"func"`
	Args []SqlStatement `json:"args"`
}

type ExpressionOperatorStatement struct {
	*AbstractSqlStatement
	Operator string `json:"operator"`
}

// ExpressionBracketsStatement 括号
type ExpressionBracketsStatement struct {
	*AbstractSqlStatement
}
