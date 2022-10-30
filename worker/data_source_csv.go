package worker

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"io"
	"os"
	"strings"
)

func NewDataSourceCsv(param *DataSourceParam) (res DataSource) {
	res = &dataSourceCsv{
		DataSourceParam: param,
	}
	return
}

type dataSourceCsv struct {
	*DataSourceParam
	saveFile *os.File
	isStop   bool
}

func (this_ *dataSourceCsv) Stop() {
	this_.isStop = true
}

func (this_ *dataSourceCsv) ReadStart() (err error) {
	return
}
func (this_ *dataSourceCsv) ReadEnd() (err error) {
	return
}
func (this_ *dataSourceCsv) Read(columnList []*dialect.ColumnModel, onRead func(data *DataSourceData) (err error)) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	f, err := os.Open(this_.Path)
	if err != nil {
		return
	}
	buf := bufio.NewReader(f)
	var line string
	var rowInfo string
	separator := this_.GetCsvSeparator()
	columnLength := len(columnList)
	if columnLength == 0 {
		err = errors.New("column is null")
		return
	}
	for {
		if this_.isStop {
			return
		}
		line, err = buf.ReadString('\n')
		if line != "" {
			if rowInfo == "" {
				rowInfo = line
			} else {
				rowInfo += line
			}

			ss := strings.Split(strings.TrimSpace(rowInfo), separator)
			if len(ss) > columnLength {
				err = errors.New("row [" + rowInfo + "] can not to column names [" + strings.Join(GetColumnNames(columnList), ",") + "]")
				return
			}
			if len(ss) == columnLength {
				rowInfo = strings.TrimSpace(rowInfo)
				if rowInfo != "" {
					rowInfo = strings.ReplaceAll(rowInfo, this_.GetLinefeed(), "\n")
					err = readRow(rowInfo, separator, columnList, onRead)
					if err != nil {
						return
					}
				}
				rowInfo = ""
			}
		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				err = nil
			}
			break
		}
	}
	if err != nil {
		return
	}
	rowInfo = strings.TrimSpace(rowInfo)
	if rowInfo != "" {
		err = readRow(rowInfo, separator, columnList, onRead)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *dataSourceCsv) WriteStart() (err error) {

	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}

	this_.saveFile, err = os.Create(this_.Path)
	if err != nil {
		return
	}
	return
}
func (this_ *dataSourceCsv) WriteEnd() (err error) {
	if this_.saveFile != nil {
		err = this_.saveFile.Close()
		return
	}
	return
}
func (this_ *dataSourceCsv) Write(data *DataSourceData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	if this_.saveFile == nil {
		err = this_.WriteStart()
		if err != nil {
			return
		}
	}

	if this_.isStop {
		return
	}
	columnList := data.ColumnList
	if data.Data == nil || columnList == nil {
		return
	}
	var valueList []string
	for _, column := range data.ColumnList {
		str := dialect.GetStringValue(data.Data[column.Name])
		str = strings.ReplaceAll(str, "\r\n", this_.GetLinefeed())
		str = strings.ReplaceAll(str, "\n", this_.GetLinefeed())
		str = strings.ReplaceAll(str, "\r", this_.GetLinefeed())
		valueList = append(valueList, str)
	}

	_, err = this_.saveFile.WriteString(strings.Join(valueList, this_.GetCsvSeparator()) + "\n")
	if err != nil {
		return
	}

	return
}
