package worker

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/team-ide/go-dialect/dialect"
	"strings"
)

func NewDataSourceExcel(param *DataSourceParam) (res DataSource) {
	res = &dataSourceExcel{
		DataSourceParam: param,
	}
	return
}

type dataSourceExcel struct {
	*DataSourceParam
	xlsxFForRead  *xlsx.File
	xlsxFForWrite *xlsx.File
	sheetForWrite *xlsx.Sheet
	isStop        bool
}

func (this_ *dataSourceExcel) Stop() {
	this_.isStop = true
}

func (this_ *dataSourceExcel) ReadStart() (err error) {
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	this_.xlsxFForRead, err = xlsx.OpenFile(this_.Path)
	if err != nil {
		err = errors.New("excel [" + this_.Path + "] open error, " + err.Error())
		return
	}
	return
}
func (this_ *dataSourceExcel) ReadEnd() (err error) {
	if this_.xlsxFForRead != nil {
	}
	return
}
func (this_ *dataSourceExcel) Read(columnList []*dialect.ColumnModel, onRead func(data *DataSourceData) (err error)) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()

	if this_.xlsxFForRead == nil {
		err = this_.ReadStart()
		if err != nil {
			return
		}
	}

	sheets := this_.xlsxFForRead.Sheets
	if this_.SheetIndex >= 0 {
		if len(sheets) < this_.SheetIndex+1 {
			err = errors.New("excel [" + this_.Path + "] sheets len is [" + fmt.Sprint(len(sheets)) + "]")
			return
		}
	}

	startRow := this_.StartRow - 1
	if startRow < 0 {
		startRow = 0
	}
	for index, sheet := range sheets {
		if this_.SheetIndex >= 0 {
			if index != this_.SheetIndex {
				continue
			}
		}
		if this_.isStop {
			return
		}
		maxRow := sheet.MaxRow

		for rowIndex := startRow; rowIndex < maxRow; rowIndex++ {

			if this_.isStop {
				return
			}

			row := sheet.Rows[rowIndex]

			var data = map[string]interface{}{}

			for cellIndex, column := range columnList {
				if cellIndex >= len(row.Cells) {
					break
				}
				cell := row.Cells[cellIndex]
				var v = cell.String()
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
			if err != nil {
				return
			}
		}
	}
	return
}

func (this_ *dataSourceExcel) save() (err error) {
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	err = this_.xlsxFForWrite.Save(this_.Path)
	if err != nil {
		err = errors.New("excel [" + this_.Path + "] save error, " + err.Error())
		return
	}
	return
}

func (this_ *dataSourceExcel) WriteStart() (err error) {
	this_.xlsxFForWrite = xlsx.NewFile()

	sheetName := this_.SheetName
	if len(sheetName) > 31 {
		sheetName = sheetName[0:30]
	}
	this_.sheetForWrite, err = this_.xlsxFForWrite.AddSheet(sheetName)
	if err != nil {
		err = errors.New("excel [" + this_.Path + "] add shell [" + this_.SheetName + "] error, " + err.Error())
		return
	}

	if len(this_.TitleList) > 0 {
		var valueList []interface{}
		for _, title := range this_.TitleList {
			valueList = append(valueList, title)
		}
		sheetWrite(this_.sheetForWrite, valueList)
	}
	err = this_.save()
	if err != nil {
		return
	}
	return
}
func (this_ *dataSourceExcel) WriteEnd() (err error) {
	err = this_.save()
	if err != nil {
		return
	}
	return
}
func (this_ *dataSourceExcel) Write(data *DataSourceData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()

	if this_.xlsxFForWrite == nil {
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
	var valueList []interface{}
	for _, column := range data.ColumnList {
		valueList = append(valueList, data.Data[column.Name])
	}
	sheetWrite(this_.sheetForWrite, valueList)
	return
}

func sheetWrite(sheet *xlsx.Sheet, valueList []interface{}) {
	row := sheet.AddRow()
	for _, value := range valueList {
		str := dialect.GetStringValue(value)
		call := row.AddCell()
		call.SetValue(str)
	}
}
