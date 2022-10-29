package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
)

func NewTaskImport(dia dialect.Dialect, db *sql.DB, ownerName string, tableName string, batchNumber int, columnNameList []string, dataSource DataSource) (res *taskImport) {
	res = &taskImport{
		Task: &Task{
			dia: dia,
			db:  db,
		},
		batchNumber:    batchNumber,
		ownerName:      ownerName,
		tableName:      tableName,
		columnNameList: columnNameList,
		dataSource:     dataSource,
	}
	return
}

type taskImport struct {
	batchNumber int
	*Task
	ownerName      string
	tableName      string
	columnNameList []string
	dataSource     DataSource
}

func (this_ *taskImport) do() (err error) {

	dataChan := make(chan map[string]interface{}, 1)
	defer func() {
		this_.dataSource.Stop()
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		close(dataChan)
	}()
	batchNumber := this_.batchNumber
	if batchNumber <= 0 {
		batchNumber = 10
	}
	go func() {
		var dataList []map[string]interface{}
		for data := range dataChan {
			dataList = append(dataList, data)
			this_.dataCountIncr()
			this_.readyCountIncr()
			if len(dataList) >= batchNumber {
				err = this_.doSave(dataList)
				if err != nil {
					return
				}
			}
		}
		err = this_.doSave(dataList)
		if err != nil {
			return
		}
	}()
	err = this_.dataSource.Read(this_.columnNameList, dataChan)
	if err != nil {
		return
	}
	return
}

func (this_ *taskImport) doSave(dataList []map[string]interface{}) (err error) {
	if len(dataList) == 0 {
		return
	}

	return
}
