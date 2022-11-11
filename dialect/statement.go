package dialect

type Statement interface {
	GetTemplate() (template string)
	GetContent() (content *string)
	GetParent() (parent Statement)
	GetChildren() (children *[]Statement)
	Format(statementContext *StatementContext) (text string, err error)
}

type AbstractStatement struct {
	Content  string      `json:"content,omitempty"`
	Children []Statement `json:"children,omitempty"`
	Parent   Statement   `json:"-"`
}

func (this_ *AbstractStatement) GetContent() (content *string) {
	content = &this_.Content
	return
}

func (this_ *AbstractStatement) GetParent() (parent Statement) {
	parent = this_.Parent
	return
}
func (this_ *AbstractStatement) GetChildren() (children *[]Statement) {
	children = &this_.Children
	return
}
func (this_ *AbstractStatement) GetTemplate() (template string) {
	template += this_.Content
	for _, one := range this_.Children {
		template += one.GetTemplate()
	}
	return
}

type RootStatement struct {
	*AbstractStatement
}

type IgnorableStatement struct {
	*AbstractStatement
}

func (this_ *IgnorableStatement) GetTemplate() (template string) {
	template += "["

	template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	template += "]"
	return
}

type TextStatement struct {
	*AbstractStatement
}

type IfStatement struct {
	*AbstractStatement
	Condition           string               `json:"condition"`
	ConditionExpression *ExpressionStatement `json:"conditionExpression"`
	ElseIfs             []*ElseIfStatement   `json:"elseIfs"`
	Else                *ElseStatement       `json:"else"`
}

func (this_ *IfStatement) GetTemplate() (template string) {
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

type ElseIfStatement struct {
	*AbstractStatement
	Condition           string               `json:"condition"`
	ConditionExpression *ExpressionStatement `json:"conditionExpression"`
	If                  *IfStatement         `json:"-"`
	Index               int                  `json:"index"`
}

func (this_ *ElseIfStatement) GetTemplate() (template string) {
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
func (this_ *ElseIfStatement) IsEndElseIf() (isEnd bool) {
	if this_.Index == len(this_.If.ElseIfs)-1 {
		isEnd = true
	}
	return
}

type ElseStatement struct {
	*AbstractStatement
	If *IfStatement `json:"-"`
}

func (this_ *ElseStatement) GetTemplate() (template string) {
	template += " else {\n"

	//template += this_.Content

	for _, one := range this_.Children {
		template += one.GetTemplate()
	}

	template += "\n}\n"
	return
}

type ForStatement struct {
	*AbstractStatement
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
	*AbstractStatement
}

type ExpressionStringStatement struct {
	*AbstractStatement
	Value string `json:"value"`
}

type ExpressionNumberStatement struct {
	*AbstractStatement
	Value float64 `json:"value"`
}

type ExpressionIdentifierStatement struct {
	*AbstractStatement
	Identifier string `json:"identifier"`
}

type ExpressionFuncStatement struct {
	*AbstractStatement
	Func string      `json:"func"`
	Args []Statement `json:"args"`
}

type ExpressionOperatorStatement struct {
	*AbstractStatement
	Operator string `json:"operator"`
}

// ExpressionBracketsStatement 括号
type ExpressionBracketsStatement struct {
	*AbstractStatement
}
