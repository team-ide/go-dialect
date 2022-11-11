package dialect

type SqlMapping struct {
	OwnersSelect string
	OwnerSelect  string
	OwnerCreate  string
	OwnerDelete  string

	TablesSelect string
	TableSelect  string
	TableCreate  string
	TableDelete  string
	TableComment string
	TableRename  string

	ColumnsSelect string
	ColumnSelect  string
	ColumnAdd     string
	ColumnDelete  string
	ColumnComment string
	ColumnRename  string
	ColumnUpdate  string

	PrimaryKeysSelect string
	PrimaryKeyAdd     string
	PrimaryKeyDelete  string

	IndexesSelect   string
	IndexAdd        string
	IndexDelete     string
	IndexNameFormat string

	MethodCache map[string]interface{}
}

type SqlMappingStatement struct {
	SqlMapping   *SqlMapping
	OwnersSelect *RootStatement
	OwnerSelect  *RootStatement
	OwnerCreate  *RootStatement
	OwnerDelete  *RootStatement

	TablesSelect *RootStatement
	TableSelect  *RootStatement
	TableCreate  *RootStatement
	TableDelete  *RootStatement
	TableComment *RootStatement
	TableRename  *RootStatement

	ColumnsSelect *RootStatement
	ColumnSelect  *RootStatement
	ColumnAdd     *RootStatement
	ColumnDelete  *RootStatement
	ColumnComment *RootStatement
	ColumnRename  *RootStatement
	ColumnUpdate  *RootStatement

	PrimaryKeysSelect *RootStatement
	PrimaryKeyAdd     *RootStatement
	PrimaryKeyDelete  *RootStatement

	IndexesSelect   *RootStatement
	IndexAdd        *RootStatement
	IndexDelete     *RootStatement
	IndexNameFormat *RootStatement
}
