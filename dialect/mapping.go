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

	OwnerNamePackChar  string
	TableNamePackChar  string
	ColumnNamePackChar string

	MethodCache map[string]interface{}
}
