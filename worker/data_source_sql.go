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

func NewDataSourceSql(param *DataSourceParam) (res DataSource) {
	res = &dataSourceSql{
		DataSourceParam: param,
	}
	return
}

type dataSourceSql struct {
	*DataSourceParam
	saveFile *os.File
	isStop   bool
}

func (this_ *dataSourceSql) Stop() {
	this_.isStop = true
}

func (this_ *dataSourceSql) ReadStart() (err error) {
	return
}
func (this_ *dataSourceSql) ReadEnd() (err error) {
	return
}
func (this_ *dataSourceSql) Read(columnList []*dialect.ColumnModel, onRead func(data *DataSourceData) (err error)) (err error) {
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
	var sqlInfo string
	for {
		if this_.isStop {
			return
		}
		line, err = buf.ReadString('\n')
		if line != "" {
			sqlInfo += line
			if this_.Dia.IsSqlEnd(sqlInfo) {
				sqlInfo = strings.TrimSpace(sqlInfo)
				if sqlInfo != "" {
					err = onRead(&DataSourceData{
						HasSql: true,
						Sql:    sqlInfo,
					})
					if err != nil {
						return
					}
				}
				sqlInfo = ""
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
	sqlInfo = strings.TrimSpace(sqlInfo)
	if sqlInfo != "" {
		err = onRead(&DataSourceData{
			HasSql: true,
			Sql:    sqlInfo,
		})
		if err != nil {
			return
		}
	}
	return
}

func (this_ *dataSourceSql) WriteStart() (err error) {

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
func (this_ *dataSourceSql) WriteEnd() (err error) {
	if this_.saveFile != nil {
		err = this_.saveFile.Close()
		return
	}
	return
}

func (this_ *dataSourceSql) WriteHeader(columnList []*dialect.ColumnModel) (err error) {

	return
}
func (this_ *dataSourceSql) Write(data *DataSourceData) (err error) {
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
	if data.HasSql {
		_, err = this_.saveFile.WriteString(data.Sql + ";\n")
		if err != nil {
			return
		}
	}
	return
}
