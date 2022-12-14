package worker

import "strings"

type DataSourceType struct {
	Name       string `json:"name"`
	New        func(param *DataSourceParam) (dataSource DataSource)
	FileSuffix string `json:"fileSuffix"`
}

var (
	DataSourceTypeSql = &DataSourceType{
		Name:       "sql",
		FileSuffix: "sql",
		New:        NewDataSourceSql,
	}
	DataSourceTypeExcel = &DataSourceType{
		Name:       "excel",
		FileSuffix: "xlsx",
		New:        NewDataSourceExcel,
	}
	DataSourceTypeText = &DataSourceType{
		Name:       "text",
		FileSuffix: "txt",
		New:        NewDataSourceText,
	}
	DataSourceTypeCsv = &DataSourceType{
		Name:       "csv",
		FileSuffix: "csv",
		New:        NewDataSourceCsv,
	}

	DataSourceTypeList = []*DataSourceType{
		DataSourceTypeSql,
		DataSourceTypeExcel,
		DataSourceTypeText,
		DataSourceTypeCsv,
	}
)

func GetDataSource(dataSourceType string) (res *DataSourceType) {
	switch strings.ToLower(dataSourceType) {
	case "sql":
		res = DataSourceTypeSql
		break
	case "txt", "text":
		res = DataSourceTypeText
		break
	case "excel":
		res = DataSourceTypeExcel
		break
	case "csv":
		res = DataSourceTypeCsv
		break
	}
	return
}
