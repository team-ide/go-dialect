package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
)

func NewTaskExec(db *sql.DB, dia dialect.Dialect, newDb func(ownerName string) (db *sql.DB, err error), taskExecParam *TaskExecParam) (res *taskExec) {

	task := &Task{
		dia:        dia,
		db:         db,
		onProgress: taskExecParam.OnProgress,
	}
	res = &taskExec{
		Task:          task,
		TaskExecParam: taskExecParam,
		newDb:         newDb,
	}
	task.do = res.do
	return
}

type TaskExecParam struct {
	Owners []*TaskExecOwner `json:"owners"`

	BatchNumber     int  `json:"batchNumber"`
	ContinueIsError bool `json:"continueIsError"`

	OnProgress func(progress *TaskProgress)
}

type TaskExecOwner struct {
	Name   string           `json:"name,omitempty"`
	Tables []*TaskExecTable `json:"tables,omitempty"`
}

type TaskExecTable struct {
	Name            string                   `json:"name,omitempty"`
	ColumnList      []*dialect.ColumnModel   `json:"columnList,omitempty"`
	InsertList      []map[string]interface{} `json:"insertList"`
	UpdateList      []map[string]interface{} `json:"updateList"`
	UpdateWhereList []map[string]interface{} `json:"updateWhereList"`
	DeleteList      []map[string]interface{} `json:"deleteList"`
}

type taskExec struct {
	*Task
	*TaskExecParam `json:"-"`
	newDb          func(ownerName string) (db *sql.DB, err error)
}

func (this_ *taskExec) do() (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	if len(this_.Owners) == 0 {
		return
	}
	for _, owner := range this_.Owners {
		err = this_.execOwner(owner)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskExec) execOwner(owner *TaskExecOwner) (err error) {
	progress := &TaskProgress{
		Title: "执行[" + owner.Name + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)

	for _, table := range owner.Tables {
		err = this_.execTableInsert(this_.db, owner.Name, table.Name, table.ColumnList, table.InsertList)
		if err != nil {
			return
		}
		err = this_.execTableUpdate(this_.db, owner.Name, table.Name, table.ColumnList, table.UpdateList, table.UpdateWhereList)
		if err != nil {
			return
		}
		err = this_.execTableDelete(this_.db, owner.Name, table.Name, table.ColumnList, table.DeleteList)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskExec) execTableInsert(workDb *sql.DB, ownerName string, tableName string, columnList []*dialect.ColumnModel, insertList []map[string]interface{}) (err error) {

	progress := &TaskProgress{
		Title: "导入表数据[" + ownerName + "." + tableName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)

	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}

	var dataList []map[string]interface{}
	for _, data := range insertList {
		dataList = append(dataList, data)
		if len(dataList) >= batchNumber {
			err = this_.execInsert(workDb, dataList, ownerName, tableName, columnList)
			dataList = make([]map[string]interface{}, 0)
			if err != nil {
				return
			}
		}
	}
	if err != nil {
		return
	}
	if len(dataList) >= 0 {
		err = this_.execInsert(workDb, dataList, ownerName, tableName, columnList)
		dataList = make([]map[string]interface{}, 0)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskExec) execInsert(workDb *sql.DB, dataList []map[string]interface{}, ownerName string, tableName string, columnList []*dialect.ColumnModel) (err error) {

	progress := &TaskProgress{
		Title: "插入数据[" + ownerName + "." + tableName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)

	_, _, batchSqlList, batchValuesList, err := this_.dia.DataListInsertSql(this_.Param, ownerName, tableName, columnList, dataList)
	if err != nil {
		return
	}
	var errSql string
	_, errSql, _, err = DoExecs(workDb, batchSqlList, batchValuesList)
	if err != nil {
		if errSql != "" {
			err = errors.New("sql:" + errSql + " exec error," + err.Error())
		}
		return
	}
	return
}

func (this_ *taskExec) execTableUpdate(workDb *sql.DB, ownerName string, tableName string, columnList []*dialect.ColumnModel, updateList []map[string]interface{}, updateWhereList []map[string]interface{}) (err error) {

	progress := &TaskProgress{
		Title: "导入表数据[" + ownerName + "." + tableName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)

	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}

	var updateDataList_ []map[string]interface{}
	var updateWhereList_ []map[string]interface{}
	for index := range updateList {
		updateDataList_ = append(updateDataList_, updateList[index])
		updateWhereList_ = append(updateWhereList_, updateWhereList[index])
		if len(updateDataList_) >= batchNumber {
			err = this_.execUpdate(workDb, updateDataList_, updateWhereList_, ownerName, tableName, columnList)
			updateDataList_ = make([]map[string]interface{}, 0)
			updateWhereList_ = make([]map[string]interface{}, 0)
			if err != nil {
				return
			}
		}
	}
	if err != nil {
		return
	}
	if len(updateDataList_) >= 0 {
		err = this_.execUpdate(workDb, updateDataList_, updateWhereList_, ownerName, tableName, columnList)
		updateDataList_ = make([]map[string]interface{}, 0)
		updateWhereList_ = make([]map[string]interface{}, 0)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskExec) execUpdate(workDb *sql.DB, updateList []map[string]interface{}, updateWhereList []map[string]interface{}, ownerName string, tableName string, columnList []*dialect.ColumnModel) (err error) {

	progress := &TaskProgress{
		Title: "插入数据[" + ownerName + "." + tableName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)
	sqlList, sqlValuesList, err := this_.dia.DataListUpdateSql(this_.Param, ownerName, tableName, columnList, updateList, updateWhereList)
	if err != nil {
		return
	}
	var errSql string
	_, errSql, _, err = DoExecs(workDb, sqlList, sqlValuesList)
	if err != nil {
		if errSql != "" {
			err = errors.New("sql:" + errSql + " exec error," + err.Error())
		}
		return
	}
	return
}

func (this_ *taskExec) execTableDelete(workDb *sql.DB, ownerName string, tableName string, columnList []*dialect.ColumnModel, deleteList []map[string]interface{}) (err error) {

	progress := &TaskProgress{
		Title: "导入表数据[" + ownerName + "." + tableName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)

	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}

	var dataWhereList []map[string]interface{}
	for _, data := range deleteList {
		dataWhereList = append(dataWhereList, data)
		if len(dataWhereList) >= batchNumber {
			err = this_.execDelete(workDb, dataWhereList, ownerName, tableName, columnList)
			dataWhereList = make([]map[string]interface{}, 0)
			if err != nil {
				return
			}
		}
	}
	if err != nil {
		return
	}
	if len(dataWhereList) >= 0 {
		err = this_.execDelete(workDb, dataWhereList, ownerName, tableName, columnList)
		dataWhereList = make([]map[string]interface{}, 0)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskExec) execDelete(workDb *sql.DB, dataWhereList []map[string]interface{}, ownerName string, tableName string, columnList []*dialect.ColumnModel) (err error) {

	progress := &TaskProgress{
		Title: "插入数据[" + ownerName + "." + tableName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)

	sqlList, sqlValuesList, err := this_.dia.DataListDeleteSql(this_.Param, ownerName, tableName, columnList, dataWhereList)
	if err != nil {
		return
	}
	var errSql string
	_, errSql, _, err = DoExecs(workDb, sqlList, sqlValuesList)
	if err != nil {
		if errSql != "" {
			err = errors.New("sql:" + errSql + " exec error," + err.Error())
		}
		return
	}
	return
}
