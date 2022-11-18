package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"strings"
)

func NewTaskSync(sourceDB *sql.DB, sourceDialect dialect.Dialect, targetDb *sql.DB, targetDialect dialect.Dialect, newDb func(ownerName string) (db *sql.DB, err error), taskSyncParam *TaskSyncParam) (res *taskSync) {
	task := &Task{
		dia:        sourceDialect,
		db:         sourceDB,
		onProgress: taskSyncParam.OnProgress,
	}
	res = &taskSync{
		Task:          task,
		targetDialect: targetDialect,
		targetDb:      targetDb,
		TaskSyncParam: taskSyncParam,
		newDb:         newDb,
	}
	task.do = res.do
	return
}

type TaskSyncParam struct {
	Owners []*TaskSyncOwner `json:"owners"`

	BatchNumber           int    `json:"batchNumber"`
	SyncStruct            bool   `json:"syncStruct"`
	SyncData              bool   `json:"syncData"`
	ContinueIsError       bool   `json:"continueIsError"`
	OwnerCreateIfNotExist bool   `json:"ownerCreateIfNotExist"`
	OwnerCreatePassword   string `json:"ownerCreatePassword"`

	FormatIndexName func(ownerName string, tableName string, index *dialect.IndexModel) string `json:"-"`
	OnProgress      func(progress *TaskProgress)                                               `json:"-"`
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
	newDb         func(ownerName string) (db *sql.DB, err error)
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

	ownerOne, err := OwnerSelect(this_.db, this_.dia, this_.Param, owner.SourceName)
	if err != nil {
		//fmt.Println("task sync syncOwner OwnerSelect owner:", owner.SourceName, " error:", err.Error())
		return
	}
	if ownerOne == nil {
		err = errors.New("source db owner [" + owner.SourceName + "] is not exist")
		return
	}

	tables := owner.Tables

	if len(tables) == 0 {
		var list []*dialect.TableModel
		list, err = TablesSelect(this_.db, this_.dia, this_.Param, owner.SourceName)
		if err != nil {
			//fmt.Println("task sync syncOwner TablesSelect owner:", owner.SourceName, " error:", err.Error())
			return
		}
		progress.Infos = append(progress.Infos, fmt.Sprintf("owner[%s] table size[%d]", owner.SourceName, len(list)))
		for _, one := range list {
			tables = append(tables, &TaskSyncTable{
				SourceName: one.TableName,
			})
		}
	}
	targetOwnerName := owner.TargetName
	if targetOwnerName == "" {
		targetOwnerName = owner.SourceName
	}

	targetOwnerOne, err := OwnerSelect(this_.targetDb, this_.targetDialect, this_.Param, targetOwnerName)
	if err != nil {
		return
	}
	if targetOwnerOne == nil {
		if !this_.OwnerCreateIfNotExist {
			err = errors.New("target db owner [" + targetOwnerName + "] is not exist")
			return
		}
		this_.addProgress(&TaskProgress{
			Title: "同步[" + targetOwnerName + "] 不存在，创建",
		})
		_, err = OwnerCreate(this_.targetDb, this_.targetDialect, this_.Param, &dialect.OwnerModel{
			OwnerName:             targetOwnerName,
			OwnerPassword:         this_.OwnerCreatePassword,
			OwnerCharacterSetName: "utf8mb4",
		})
		if err != nil {
			return
		}
	}

	workDb, err := this_.newDb(targetOwnerName)
	if err != nil {
		return
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

		err = this_.syncTable(workDb, owner.SourceName, table.SourceName, owner.TargetName, table.TargetName)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskSync) syncTable(workDb *sql.DB, sourceOwnerName string, sourceTableName string, targetOwnerName string, targetTableName string) (err error) {
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

	newTableDetail, err := TableDetail(this_.db, this_.dia, this_.Param, sourceOwnerName, sourceTableName, false)
	if err != nil {
		return
	}
	if newTableDetail == nil {
		err = errors.New("source db table [" + sourceOwnerName + "." + sourceTableName + "] is not exist")
		return
	}
	oldTableDetail, err := TableDetail(this_.targetDb, this_.targetDialect, this_.Param, targetOwnerName, targetTableName, false)
	if err != nil {
		return
	}
	if oldTableDetail == nil {
		oldTableDetail = &dialect.TableModel{}
	}
	oldTableDetail.OwnerName = targetOwnerName
	oldTableDetail.TableName = targetTableName
	if this_.SyncStruct {
		err = this_.syncTableSyncStructure(workDb, newTableDetail, oldTableDetail)
		if err != nil {
			return
		}
	}
	if this_.SyncData {
		err = this_.syncTableSyncData(workDb, newTableDetail, oldTableDetail)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskSync) syncTableSyncStructure(workDb *sql.DB, newTableDetail *dialect.TableModel, oldTableDetail *dialect.TableModel) (err error) {

	progress := &TaskProgress{
		Title: "同步表结构[" + newTableDetail.OwnerName + "." + newTableDetail.TableName + "] 到 [" + oldTableDetail.OwnerName + "." + oldTableDetail.TableName + "]",
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
			index.IndexName = this_.FormatIndexName(oldTableDetail.OwnerName, oldTableDetail.TableName, index)
		}
	}
	if len(oldTableDetail.ColumnList) == 0 {
		newTableDetail.TableName = oldTableDetail.TableName
		err = TableCreate(workDb, this_.targetDialect, this_.Param, oldTableDetail.OwnerName, newTableDetail)
		if err != nil {
			return
		}
		return
	} else {
		err = TableUpdate(workDb, this_.targetDialect, oldTableDetail, this_.dia, newTableDetail)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskSync) syncTableSyncData(workDb *sql.DB, newTableDetail *dialect.TableModel, oldTableDetail *dialect.TableModel) (err error) {

	progress := &TaskProgress{
		Title: "同步表数据[" + newTableDetail.OwnerName + "." + newTableDetail.TableName + "] 到 [" + oldTableDetail.OwnerName + "." + oldTableDetail.TableName + "]",
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
		columnNames = append(columnNames, one.ColumnName)
	}
	selectSqlInfo += this_.dia.ColumnNamesPack(this_.Param, columnNames)
	selectSqlInfo += " FROM "
	if newTableDetail.OwnerName != "" {
		selectSqlInfo += this_.dia.OwnerNamePack(this_.Param, newTableDetail.OwnerName) + "."
	}
	selectSqlInfo += this_.dia.TableNamePack(this_.Param, newTableDetail.TableName)

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
		err = this_.insertDataList(workDb, dataList, oldTableDetail.OwnerName, oldTableDetail.TableName, columnList)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskSync) insertDataList(workDb *sql.DB, dataList []map[string]interface{}, targetOwnerName string, targetTableName string, columnList []*dialect.ColumnModel) (err error) {

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

	_, sqlList, err := this_.targetDialect.InsertDataListSql(this_.Param, targetOwnerName, targetTableName, columnList, dataList)
	if err != nil {
		return
	}
	var errSql string
	errSql, err = DoExec(workDb, sqlList)
	if err != nil {
		if errSql != "" {
			err = errors.New("sql:" + errSql + " exec error," + err.Error())
		}
		return
	}
	return
}
