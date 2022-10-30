package worker

import (
	"errors"
	"fmt"
	"os"
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

func (this_ *dataSourceSql) open() (file *os.File, err error) {
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	file, err = os.Open(this_.Path)
	if err != nil {
		err = errors.New("sql [" + this_.Path + "] open error, " + err.Error())
		return
	}
	return
}

func (this_ *dataSourceSql) create() (file *os.File, err error) {
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	file, err = os.Create(this_.Path)
	if err != nil {
		return
	}
	return
}

func (this_ *dataSourceSql) ReadStart() (err error) {
	return
}
func (this_ *dataSourceSql) ReadEnd() (err error) {
	return
}
func (this_ *dataSourceSql) Read(onRead func(data *DataSourceData) (err error)) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	bs, err := os.ReadFile(this_.Path)
	if err != nil {
		return
	}
	sqlInfo := string(bs)
	sqlList := SplitSqlList(sqlInfo)
	for _, sqlOne := range sqlList {
		err = onRead(&DataSourceData{
			HasSql: true,
			Sql:    sqlOne,
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
func (this_ *dataSourceSql) Write(data *DataSourceData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	if this_.saveFile != nil {
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
