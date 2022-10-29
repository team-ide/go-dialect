package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"strings"
)

func NewTaskSync(sourceDB *sql.DB, sourceDialect dialect.Dialect, targetDb *sql.DB, targetDialect dialect.Dialect, taskSyncParam *TaskSyncParam) (res *taskSync) {
	task := &Task{
		dia: sourceDialect,
		db:  sourceDB,
	}
	res = &taskSync{
		Task:          task,
		targetDialect: targetDialect,
		targetDb:      targetDb,
		TaskSyncParam: taskSyncParam,
	}
	task.do = res.do
	return
}

type TaskSyncParam struct {
	Owners []*TaskSyncOwner `json:"owners"`

	BatchNumber     int  `json:"batchNumber"`
	SyncStructure   bool `json:"syncStructure"`
	SyncData        bool `json:"syncData"`
	ContinueIsError bool `json:"continueIsError"`
	FormatIndexName func(ownerName string, tableName string, index *dialect.IndexModel) string
}

type TaskSyncOwner struct {
	SourceName     string           `json:"sourceName"`
	TargetName     string           `json:"targetName"`
	SkipTableNames []string         `json:"skipTableNames"`
	Tables         []*TaskSyncTable `json:"tables"`
}

type TaskSyncTable struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
}

type taskSync struct {
	*Task
	*TaskSyncParam
	targetDialect dialect.Dialect
	targetDb      *sql.DB
}

func (this_ *taskSync) do() (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()

	if len(this_.Owners) == 0 {
		return
	}
	for _, owner := range this_.Owners {
		err = this_.syncOwner(owner)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskSync) syncOwner(owner *TaskSyncOwner) (err error) {
	progress := &TaskProgress{
		Title: "同步[" + owner.SourceName + "]",
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

	ownerOne, err := OwnerSelect(this_.db, this_.dia, owner.SourceName)
	if err != nil {
		return
	}
	if ownerOne == nil {
		err = errors.New("source db owner [" + owner.SourceName + "] is not exist")
		return
	}

	tables := owner.Tables

	if len(tables) == 0 {
		var list []*dialect.TableModel
		list, err = TablesSelect(this_.db, this_.dia, owner.SourceName)
		if err != nil {
			return
		}
		for _, one := range list {
			tables = append(tables, &TaskSyncTable{
				SourceName: one.Name,
			})
		}
	}

	for _, table := range tables {

		if len(owner.SkipTableNames) > 0 {
			var skip bool
			for _, skipTableName := range owner.SkipTableNames {
				if strings.EqualFold(table.SourceName, skipTableName) {
					skip = true
				}
			}
			if skip {
				continue
			}
		}

		err = this_.syncTable(owner.SourceName, table.SourceName, owner.TargetName, table.TargetName)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskSync) syncTable(sourceOwnerName string, sourceTableName string, targetOwnerName string, targetTableName string) (err error) {
	if targetOwnerName == "" {
		targetOwnerName = sourceOwnerName
	}
	if targetTableName == "" {
		targetTableName = sourceTableName
	}
	progress := &TaskProgress{
		Title: "同步[" + sourceOwnerName + "." + sourceTableName + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
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

	newTableDetail, err := TableDetail(this_.db, this_.dia, sourceOwnerName, sourceTableName)
	if err != nil {
		return
	}
	if newTableDetail == nil {
		err = errors.New("source db table [" + sourceOwnerName + "." + sourceTableName + "] is not exist")
		return
	}
	oldTableDetail, err := TableDetail(this_.targetDb, this_.targetDialect, targetOwnerName, targetTableName)
	if err != nil {
		return
	}
	if oldTableDetail == nil {
		oldTableDetail = &dialect.TableModel{}
	}
	oldTableDetail.OwnerName = targetOwnerName
	oldTableDetail.Name = targetTableName
	if this_.SyncStructure {
		err = this_.syncTableSyncStructure(newTableDetail, oldTableDetail)
		if err != nil {
			return
		}
	}
	if this_.SyncData {
		err = this_.syncTableSyncData(newTableDetail, oldTableDetail)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskSync) syncTableSyncStructure(newTableDetail *dialect.TableModel, oldTableDetail *dialect.TableModel) (err error) {

	progress := &TaskProgress{
		Title: "同步表结构[" + newTableDetail.OwnerName + "." + newTableDetail.Name + "] 到 [" + oldTableDetail.OwnerName + "." + oldTableDetail.Name + "]",
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

	if this_.FormatIndexName != nil {
		for _, index := range newTableDetail.IndexList {
			index.Name = this_.FormatIndexName(oldTableDetail.OwnerName, oldTableDetail.Name, index)
		}
	}

	if len(oldTableDetail.ColumnList) == 0 {
		newTableDetail.Name = oldTableDetail.Name
		err = TableCreate(this_.targetDb, this_.targetDialect, oldTableDetail.OwnerName, newTableDetail)
		if err != nil {
			return
		}
		return
	} else {
		err = TableUpdate(this_.targetDb, this_.targetDialect, oldTableDetail, this_.dia, newTableDetail)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskSync) syncTableSyncData(newTableDetail *dialect.TableModel, oldTableDetail *dialect.TableModel) (err error) {

	progress := &TaskProgress{
		Title: "同步表数据[" + newTableDetail.OwnerName + "." + newTableDetail.Name + "] 到 [" + oldTableDetail.OwnerName + "." + oldTableDetail.Name + "]",
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

	selectSqlInfo := "SELECT "
	var columnNames []string
	for _, one := range newTableDetail.ColumnList {
		columnNames = append(columnNames, one.Name)
	}
	selectSqlInfo += this_.dia.PackColumns(columnNames)
	selectSqlInfo += " FROM "
	if newTableDetail.OwnerName != "" {
		selectSqlInfo += this_.dia.PackOwner(newTableDetail.OwnerName) + "."
	}
	selectSqlInfo += this_.dia.PackTable(newTableDetail.Name)

	list, err := DoQuery(this_.db, selectSqlInfo)
	if err != nil {
		return
	}
	var columnList = newTableDetail.ColumnList
	if len(oldTableDetail.ColumnList) > 0 {
		columnList = oldTableDetail.ColumnList
	}
	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}

	dataListGroup := SplitArrayMap(list, batchNumber)

	for _, dataList := range dataListGroup {
		err = this_.insertDataList(dataList, oldTableDetail.OwnerName, oldTableDetail.Name, columnList)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskSync) insertDataList(dataList []map[string]interface{}, targetOwnerName string, targetTableName string, columnList []*dialect.ColumnModel) (err error) {

	progress := &TaskProgress{
		Title: "插入数据[" + targetOwnerName + "." + targetTableName + "]",
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

	_, sqlList, err := InsertDataListSql(this_.targetDialect, targetOwnerName, targetTableName, columnList, dataList)
	if err != nil {
		return
	}
	var errSql string
	errSql, err = DoExec(this_.targetDb, sqlList)
	if err != nil {
		if errSql != "" {
			err = errors.New("sql:" + errSql + " exec error," + err.Error())
		}
		return
	}
	return
}
