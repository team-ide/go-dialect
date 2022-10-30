package worker

type DataSourceType struct {
	Name          string `json:"name"`
	New           func(param *DataSourceParam) (dataSource DataSource)
	OwnerFileName func(ownerName string) (fileName string)
	TableFileName func(ownerName string, tableName string) (fileName string)
}

var (
	DataSourceTypeSql = &DataSourceType{
		Name: "sql",
		New:  NewDataSourceSql,
		OwnerFileName: func(ownerName string) (fileName string) {
			fileName = formatFileName(ownerName, "") + ".sql"
			return
		},
	}
	DataSourceTypeExcel = &DataSourceType{
		Name: "excel",
		New:  NewDataSourceExcel,
		TableFileName: func(ownerName string, tableName string) (fileName string) {
			fileName = formatFileName(ownerName, tableName) + ".xlsx"
			return
		},
	}

	DataSourceTypeList = []*DataSourceType{
		DataSourceTypeSql,
		DataSourceTypeExcel,
	}
)

func formatFileName(ownerName string, tableName string) (fileName string) {
	if ownerName != "" {
		if fileName != "" {
			fileName += "-"
		}
		fileName += ownerName
	}
	if tableName != "" {
		if fileName != "" {
			fileName += "-"
		}
		fileName += tableName
	}
	return
}
