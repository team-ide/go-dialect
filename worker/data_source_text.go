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

func NewDataSourceText(param *DataSourceParam) (res DataSource) {
	res = &dataSourceText{
		DataSourceParam: param,
	}
	return
}

type dataSourceText struct {
	*DataSourceParam
	saveFile *os.File
	isStop   bool
}

func (this_ *dataSourceText) Stop() {
	this_.isStop = true
}

func (this_ *dataSourceText) ReadStart() (err error) {
	return
}
func (this_ *dataSourceText) ReadEnd() (err error) {
	return
}
func (this_ *dataSourceText) Read(columnList []*dialect.ColumnModel, onRead func(data *DataSourceData) (err error)) (err error) {
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
	separator := this_.GetTextSeparator()
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

func GetColumnNames(columnList []*dialect.ColumnModel) (columnNames []string) {
	for _, column := range columnList {
		columnNames = append(columnNames, column.Name)
	}
	return
}

func readRow(rowInfo string, separator string, columnList []*dialect.ColumnModel, onRead func(data *DataSourceData) (err error)) (err error) {
	calls := strings.Split(rowInfo, separator)
	data := make(map[string]interface{})
	if len(calls) != len(columnList) {
		err = errors.New("row [" + rowInfo + "] can not to column names [" + strings.Join(GetColumnNames(columnList), ",") + "]")
		return
	}
	for i, column := range columnList {
		v := calls[i]
		if !column.NotNull && v == "" {
			continue
		}
		if strings.EqualFold(column.Type, "timestamp") {
			if v == "" {
				continue
			}
		}

		data[column.Name] = v
	}
	err = onRead(&DataSourceData{
		HasData: true,
		Data:    data,
	})
	return
}

func (this_ *dataSourceText) WriteStart() (err error) {

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
func (this_ *dataSourceText) WriteEnd() (err error) {
	if this_.saveFile != nil {
		err = this_.saveFile.Close()
		return
	}
	return
}
func (this_ *dataSourceText) Write(data *DataSourceData) (err error) {
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

	_, err = this_.saveFile.WriteString(strings.Join(valueList, this_.GetTextSeparator()) + "\n")
	if err != nil {
		return
	}

	return
}

func isRowEnd(rowInfo string) (isEnd bool) {

	var inStringLevel int
	var inStringPack byte
	var thisChar byte
	var lastChar byte

	var stringPackChars = []byte{'"', '\''}
	for i := 0; i < len(rowInfo); i++ {
		thisChar = rowInfo[i]
		if i > 0 {
			lastChar = rowInfo[i-1]
		}

		// inStringLevel == 0 表示 不在 字符串 包装 中
		if thisChar == ';' && inStringLevel == 0 {
		} else {
			packCharIndex := dialect.BytesIndex(stringPackChars, thisChar)
			if packCharIndex >= 0 {
				// inStringLevel == 0 表示 不在 字符串 包装 中
				if inStringLevel == 0 {
					inStringPack = stringPackChars[packCharIndex]
					// 字符串包装层级 +1
					inStringLevel++
				} else {
					// 如果有转义符号 类似 “\'”，“\"”
					if lastChar == '\\' {
					} else if lastChar == inStringPack {
						// 如果 前一个字符 与字符串包装字符一致
						inStringLevel--
					} else {
						// 字符串包装层级 -1
						inStringLevel--
					}
				}
			}
		}

	}
	isEnd = inStringLevel == 0
	return
}
