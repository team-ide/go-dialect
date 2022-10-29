package worker

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/team-ide/go-dialect/dialect"
)

func NewDataSourceExcel(path string, sheetIndex int, skipRow int) (res *DataSourceExcel) {
	res = &DataSourceExcel{
		Path:       path,
		SheetIndex: sheetIndex,
		SkipRow:    skipRow,
	}
	return
}

type DataSourceExcel struct {
	Path       string `json:"path"`
	SheetIndex int
	SkipRow    int
	isStop     bool
}

func (this_ *DataSourceExcel) Stop() {
	this_.isStop = true
}
func (this_ *DataSourceExcel) Read(nameList []string, dataChan chan map[string]interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()

	xlsxF, err := this_.open()
	if err != nil {
		return
	}

	sheets := xlsxF.Sheets
	if this_.SheetIndex >= 0 {
		if len(sheets) < this_.SheetIndex+1 {
			err = errors.New("excel [" + this_.Path + "] sheets len is [" + fmt.Sprint(len(sheets)) + "]")
			return
		}
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

		if this_.SkipRow < 0 {
			this_.SkipRow = 0
		}
		for rowIndex := this_.SkipRow; rowIndex < maxRow; rowIndex++ {

			if this_.isStop {
				return
			}

			row := sheet.Rows[rowIndex]

			var data = map[string]interface{}{}

			for cellIndex, name := range nameList {
				if cellIndex >= len(row.Cells) {
					break
				}
				cell := row.Cells[cellIndex]
				var value = cell.String()
				data[name] = value
			}

			dataChan <- data
		}
	}
	return
}

func (this_ *DataSourceExcel) open() (xlsxF *xlsx.File, err error) {
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	xlsxF, err = xlsx.OpenFile(this_.Path)
	if err != nil {
		err = errors.New("excel [" + this_.Path + "] open error, " + err.Error())
		return
	}
	return
}

func (this_ *DataSourceExcel) Write(sheetName string, titles []string, nameList []string, dataChan chan map[string]interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()

	xlsxF, err := this_.open()
	if err != nil {
		return
	}
	sheet, err := xlsxF.AddSheet(sheetName)
	if err != nil {
		err = errors.New("excel [" + this_.Path + "] add shell error, " + err.Error())
		return
	}
	if len(titles) > 0 {
		var valueList []interface{}
		for _, title := range titles {
			valueList = append(valueList, title)
		}
		sheetWrite(sheet, valueList)
	}
	for {
		if this_.isStop {
			break
		}
		data, ok := <-dataChan
		if !ok {
			break
		}
		var valueList []interface{}
		for _, name := range nameList {
			valueList = append(valueList, data[name])
		}
		if this_.isStop {
			break
		}
		sheetWrite(sheet, valueList)
	}
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
