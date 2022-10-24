package dialect

type InsertModel struct {
	DatabaseName string          `json:"databaseName"`
	TableName    string          `json:"tableName"`
	Columns      []string        `json:"columns"`
	Rows         [][]*ValueModel `json:"rows"`
}

type ValueModel struct {
	Type  ValueType `json:"type"`
	Value string    `json:"value"`
}

type ValueType string

var (
	ValueTypeString ValueType = "string"
	ValueTypeNumber ValueType = "number"
	ValueTypeFunc   ValueType = "func"
)
