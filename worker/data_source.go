package worker

import "github.com/team-ide/go-dialect/dialect"

type DataSource interface {
	Stop()
	ReadStart() (err error)
	Read(columnList []*dialect.ColumnModel, onRead func(data *DataSourceData) (err error)) (err error)
	ReadEnd() (err error)
	WriteStart() (err error)
	Write(data *DataSourceData) (err error)
	WriteEnd() (err error)
	WriteHeader(columnList []*dialect.ColumnModel) (err error)
}

type DataSourceData struct {
	HasSql     bool
	Sql        string
	HasData    bool
	Data       map[string]interface{}
	ColumnList []*dialect.ColumnModel
}

type DataSourceParam struct {
	Path       string
	Separator  string
	SheetIndex int
	StartRow   int
	SheetName  string
	Linefeed   string
	TitleList  []string
	Dia        dialect.Dialect
}

func (this_ *DataSourceParam) GetTextSeparator() string {
	if this_.Separator != "" {
		return this_.Separator
	}
	return "|:-:|"
}

func (this_ *DataSourceParam) GetCsvSeparator() string {
	if this_.Separator != "" {
		return this_.Separator
	}
	return ","
}

func (this_ *DataSourceParam) GetLinefeed() string {
	if this_.Linefeed != "" {
		return this_.Linefeed
	}
	return "|:-n-:|"
}
