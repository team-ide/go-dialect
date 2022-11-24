package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"strings"
)

func NewTaskSync(sourceDB *sql.DB, sourceDialect dialect.Dialect, targetDb *sql.DB, targetDialect dialect.Dialect, newDb func(owner *TaskSyncOwner) (db *sql.DB, err error), taskSyncParam *TaskSyncParam) (res *taskSync) {
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

	BatchNumber           int  `json:"batchNumber"`
	SyncStruct            bool `json:"syncStruct"`
	SyncData              bool `json:"syncData"`
	ErrorContinue         bool `json:"errorContinue"`
	OwnerCreateIfNotExist bool `json:"ownerCreateIfNotExist"`

	FormatIndexName func(ownerName string, tableName string, index *dialect.IndexModel) string `json:"-"`
	OnProgress      func(progress *TaskProgress)                                               `json:"-"`
}

type TaskSyncOwner struct {
	SourceName     string           `json:"sourceName"`
	TargetName     string           `json:"targetName"`
	SkipTableNames []string         `json:"skipTableNames"`
	Tables         []*TaskSyncTable `json:"tables"`
	Username       string           `json:"username"`
	Password       string           `json:"password"`
}

type TaskSyncTable struct {
	SourceName string            `json:"sourceName"`
	TargetName string            `json:"targetName"`
	Columns    []*TaskSyncColumn `json:"columns"`
}

type TaskSyncColumn struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
}

type taskSync struct {
	*Task
	*TaskSyncParam
	targetDialect dialect.Dialect
	targetDb      *sql.DB
	newDb         func(owner *TaskSyncOwner) (db *sql.DB, err error)
}

func (this_ *taskSync) do() (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()

	owners := this_.Owners
	if len(owners) == 0 {
		return
	}
	this_.countIncr(&this_.OwnerCount, len(owners))
	for _, owner := range owners {
		var success bool
		success, err = this_.syncOwner(owner)
		if success {
			this_.countIncr(&this_.OwnerSuccessCount, 1)
		} else {
			this_.countIncr(&this_.OwnerErrorCount, 1)
		}
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskSync) syncOwner(owner *TaskSyncOwner) (success bool, err error) {
	progress := &TaskProgress{
		Title: "同步[" + owner.SourceName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
			if progress.OnError != nil {
				progress.OnError(err)
			}
		}

		if this_.ErrorContinue {
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
			OwnerPassword:         owner.Password,
			OwnerCharacterSetName: "utf8mb4",
		})
		if err != nil {
			return
		}
	}

	workDb, err := this_.newDb(owner)
	if err != nil {
		return
	}

	this_.countIncr(&this_.TableCount, len(tables))
	for _, table := range tables {

		if len(owner.SkipTableNames) > 0 {
			var skip bool
			for _, skipTableName := range owner.SkipTableNames {
				if strings.EqualFold(table.SourceName, skipTableName) {
					skip = true
				}
			}
			if skip {
				this_.countIncr(&this_.TableSuccessCount, 1)
				continue
			}
		}

		var success_ bool
		success_, err = this_.syncTable(workDb, owner.SourceName, table.SourceName, owner.TargetName, table.TargetName)
		if success_ {
			this_.countIncr(&this_.TableSuccessCount, 1)
		} else {
			this_.countIncr(&this_.TableErrorCount, 1)
		}
		if err != nil {
			return
		}
	}
	success = true

	return
}

func (this_ *taskSync) syncTable(workDb *sql.DB, sourceOwnerName string, sourceTableName string, targetOwnerName string, targetTableName string) (success bool, err error) {
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
			if progress.OnError != nil {
				progress.OnError(err)
			}
		}

		if this_.ErrorContinue {
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
	success = true
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
			if progress.OnError != nil {
				progress.OnError(err)
			}
		}

		if this_.ErrorContinue {
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
			if progress.OnError != nil {
				progress.OnError(err)
			}
		}

		if this_.ErrorContinue {
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

	countSql, err := dialect.FormatCountSql(selectSqlInfo)
	if err != nil {
		return
	}
	totalCount, err := DoQueryCount(this_.db, countSql, nil)
	if err != nil {
		return
	}

	this_.countIncr(&this_.DataCount, totalCount)
	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}
	var pageSize = batchNumber
	var pageNo = 1

	var dataList []map[string]interface{}
	var columnList = newTableDetail.ColumnList
	if len(oldTableDetail.ColumnList) > 0 {
		columnList = oldTableDetail.ColumnList
	}
	for {

		if this_.IsStop {
			return
		}
		pageSql := this_.dia.PackPageSql(selectSqlInfo, pageSize, pageNo)
		dataList, err = DoQuery(this_.db, pageSql, nil)
		if err != nil {
			return
		}
		pageNo += 1
		dataListCount := len(dataList)
		this_.countIncr(&this_.DataReadyCount, dataListCount)
		if dataListCount == 0 {
			break
		}
		var success bool
		success, err = this_.insertDataList(workDb, dataList, oldTableDetail.OwnerName, oldTableDetail.TableName, columnList)
		if success {
			this_.countIncr(&this_.DataSuccessCount, dataListCount)
		} else {
			this_.countIncr(&this_.DataErrorCount, dataListCount)
		}
		if err != nil {
			return
		}
		if dataListCount == 0 {
			break
		}
	}

	return
}

func (this_ *taskSync) insertDataList(workDb *sql.DB, dataList []map[string]interface{}, targetOwnerName string, targetTableName string, columnList []*dialect.ColumnModel) (success bool, err error) {

	progress := &TaskProgress{
		Title: "插入数据[" + targetOwnerName + "." + targetTableName + "]",
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			progress.Error = err.Error()
			if progress.OnError != nil {
				progress.OnError(err)
			}
		}

		if this_.ErrorContinue {
			err = nil
		}
	}()

	this_.addProgress(progress)

	_, sqlList, err := this_.targetDialect.InsertDataListSql(this_.Param, targetOwnerName, targetTableName, columnList, dataList)
	if err != nil {
		return
	}
	var errSql string
	_, errSql, _, err = DoExecs(workDb, sqlList, nil)
	if err != nil {
		if errSql != "" {
			err = errors.New("sql:" + errSql + " exec error," + err.Error())
		}
		return
	}
	success = true
	return
}
