package dialect

import (
	"fmt"
	"strings"
	"sync"
)

type SqlMapping struct {
	dialectType *Type

	CanAppendOwnerName bool

	columnTypeInfoList      []*ColumnTypeInfo
	columnTypeInfoCache     map[string]*ColumnTypeInfo
	columnTypeInfoCacheLock sync.Mutex

	indexTypeInfoList      []*IndexTypeInfo
	indexTypeInfoCache     map[string]*IndexTypeInfo
	indexTypeInfoCacheLock sync.Mutex

	OwnersSelect string
	OwnerSelect  string
	OwnerCreate  string
	OwnerDelete  string

	TablesSelect                string
	TableSelect                 string
	TableCreate                 string
	TableCreateColumn           string
	TableCreateColumnHasComment bool
	TableCreatePrimaryKey       string
	TableDelete                 string
	TableComment                string
	TableRename                 string

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
	IndexNameMaxLen int

	OwnerNamePackChar  string
	TableNamePackChar  string
	ColumnNamePackChar string
	SqlValuePackChar   string
	SqlValueEscapeChar string

	MethodCache map[string]interface{}
}

func (this_ *SqlMapping) DialectType() (dialectType *Type) {
	dialectType = this_.dialectType
	return
}

func (this_ *SqlMapping) GenDemoTable() (table *TableModel) {
	table = &TableModel{
		TableName:    "TABLE_DEMO",
		TableComment: "TABLE_DEMO_comment",
	}
	columnTypeInfos := this_.GetColumnTypeInfos()
	var lastIndexColumnIndex int
	for i, columnTypeInfo := range columnTypeInfos {
		column := &ColumnModel{}
		column.ColumnName = fmt.Sprintf("column_%d", i)
		column.ColumnDataType = columnTypeInfo.Name
		column.ColumnLength = 5
		column.ColumnDecimal = 2
		column.ColumnComment = fmt.Sprintf("column_%d-comment", i)

		if columnTypeInfo.IsEnum {
			column.ColumnEnums = append(column.ColumnEnums, "option1")
			column.ColumnEnums = append(column.ColumnEnums, "option2")
		}

		if i%3 == 0 {
			column.ColumnNotNull = true
		}
		table.AddColumn(column)
		if len(table.PrimaryKeys) > 2 {
			continue
		}
		lastIndexColumnIndex = i
		if !strings.EqualFold(columnTypeInfo.Name, "text") &&
			!strings.EqualFold(columnTypeInfo.Name, "blob") {
			table.PrimaryKeys = append(table.PrimaryKeys, column.ColumnName)
		}
	}
	indexTypeInfos := this_.GetIndexTypeInfos()
	for _, indexTypeInfo := range indexTypeInfos {
		index := &IndexModel{}
		for i, column := range table.ColumnList {
			if i <= lastIndexColumnIndex {
				continue
			}
			if len(indexTypeInfo.OnlySupportDataTypes) > 0 {
				if StringsIndex(indexTypeInfo.OnlySupportDataTypes, strings.ToUpper(column.ColumnDataType)) < 0 {
					continue
				}
			}
			if StringsIndex(indexTypeInfo.NotSupportDataTypes, strings.ToUpper(column.ColumnDataType)) >= 0 {
				continue
			}
			lastIndexColumnIndex = i
			index.ColumnNames = append(index.ColumnNames, column.ColumnName)
			if len(index.ColumnNames) >= 1 {
				break
			}
		}
		if len(index.ColumnNames) == 0 {
			continue
		}
		index.IndexType = indexTypeInfo.Name
		table.AddIndex(index)
	}
	return
}
