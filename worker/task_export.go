package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"os"
	"strings"
)

func NewTaskExport(db *sql.DB, dia dialect.Dialect, targetDialect dialect.Dialect, taskExportParam *TaskExportParam) (res *taskExport) {
	if targetDialect == nil {
		targetDialect = dia
	}
	if taskExportParam.DataSourceType == nil {
		taskExportParam.DataSourceType = DataSourceTypeSql
	}
	task := &Task{
		dia:        dia,
		db:         db,
		onProgress: taskExportParam.OnProgress,
	}
	res = &taskExport{
		Task:            task,
		TaskExportParam: taskExportParam,
		targetDialect:   targetDialect,
	}
	task.do = res.do
	return
}

type TaskExportParam struct {
	Owners         []*TaskExportOwner `json:"owners"`
	SkipOwnerNames []string           `json:"skipOwnerNames"`

	DataSourceType  *DataSourceType `json:"dataSourceType"`
	BatchNumber     int             `json:"batchNumber"`
	ExportStruct    bool            `json:"exportStruct"`
	ExportData      bool            `json:"exportData"`
	ExportBatchSql  bool            `json:"exportBatchSql"`
	ErrorContinue   bool            `json:"errorContinue"`
	Dir             string          `json:"dir"`
	AppendOwnerName bool            `json:"appendOwnerName"`

	FormatIndexName func(ownerName string, tableName string, index *dialect.IndexModel) string `json:"-"`
	OnProgress      func(progress *TaskProgress)                                               `json:"-"`
}

type TaskExportOwner struct {
	SourceName     string             `json:"sourceName"`
	TargetName     string             `json:"targetName"`
	SkipTableNames []string           `json:"skipTableNames"`
	Tables         []*TaskExportTable `json:"tables"`
}

type TaskExportTable struct {
	SourceName string              `json:"sourceName"`
	TargetName string              `json:"targetName"`
	Columns    []*TaskExportColumn `json:"columns"`
}

type TaskExportColumn struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
	Value      string `json:"value"`
}

type taskExport struct {
	*Task
	*TaskExportParam
	targetDialect dialect.Dialect
}

func (this_ *taskExport) do() (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	owners := this_.Owners
	if len(owners) == 0 {
		var list []*dialect.OwnerModel
		list, err = OwnersSelect(this_.db, this_.dia, this_.Param)
		if err != nil {
			return
		}
		for _, one := range list {
			owners = append(owners, &TaskExportOwner{
				SourceName: one.OwnerName,
			})
		}
	}

	this_.countIncr(&this_.OwnerCount, len(owners))
	for _, owner := range owners {
		if len(this_.SkipOwnerNames) > 0 {
			var skip bool
			for _, skipTableName := range this_.SkipOwnerNames {
				if strings.EqualFold(owner.SourceName, skipTableName) {
					skip = true
				}
			}
			if skip {
				this_.countIncr(&this_.OwnerSuccessCount, 1)
				continue
			}
		}
		var success bool
		success, err = this_.exportOwner(owner)
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

func (this_ *taskExport) getFileName(dir string, name string) (fileName string, err error) {
	var exist bool
	if this_.Dir != "" {
		if dir != "" {
			dir = this_.Dir + string(os.PathSeparator) + dir
		} else {
			dir = this_.Dir
		}
	}
	if dir != "" {
		exist, err = PathExists(dir)
		if err != nil {
			return
		}
		if !exist {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				return
			}
		}
		fileName = dir + string(os.PathSeparator)
	}
	fileName += name
	return
}

func (this_ *taskExport) exportOwner(owner *TaskExportOwner) (success bool, err error) {
	progress := &TaskProgress{
		Title: "导出[" + owner.SourceName + "]",
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

	if this_.IsStop {
		return
	}

	ownerOne, err := OwnerSelect(this_.db, this_.dia, this_.Param, owner.SourceName)
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
		list, err = TablesSelect(this_.db, this_.dia, this_.Param, owner.SourceName)
		if err != nil {
			return
		}
		for _, one := range list {
			tables = append(tables, &TaskExportTable{
				SourceName: one.TableName,
			})
		}
	}
	this_.countIncr(&this_.TableCount, len(tables))

	ownerName := owner.TargetName
	if ownerName == "" {
		ownerName = owner.SourceName
	}

	var ownerDataSource DataSource
	if this_.DataSourceType == DataSourceTypeSql {
		fileName := ownerName + "." + this_.DataSourceType.FileSuffix
		fileName, err = this_.getFileName("", fileName)
		if err != nil {
			return
		}
		param := &DataSourceParam{
			Path:      fileName,
			SheetName: ownerName,
			Dia:       this_.targetDialect,
		}
		ownerDataSource = this_.DataSourceType.New(param)
		err = ownerDataSource.WriteStart()
		if err != nil {
			return
		}
		defer func() {
			_ = ownerDataSource.WriteEnd()
		}()
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
				this_.countIncr(&this_.TableSuccessCount, 1)
				continue
			}
		}
		var success_ bool
		success_, err = this_.exportTable(ownerDataSource, owner.SourceName, table.SourceName, owner.TargetName, table.TargetName, table.Columns)
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

func (this_ *taskExport) exportTable(ownerDataSource DataSource, sourceOwnerName string, sourceTableName string, targetOwnerName string, targetTableName string, columns []*TaskExportColumn) (success bool, err error) {
	if targetOwnerName == "" {
		targetOwnerName = sourceOwnerName
	}
	if targetTableName == "" {
		targetTableName = sourceTableName
	}
	progress := &TaskProgress{
		Title: "导出[" + sourceOwnerName + "." + sourceTableName + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
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

	if this_.IsStop {
		return
	}
	tableDetail, err := TableDetail(this_.db, this_.dia, this_.Param, sourceOwnerName, sourceTableName, false)
	if err != nil {
		return
	}
	if tableDetail == nil {
		err = errors.New("source db table [" + sourceOwnerName + "." + sourceTableName + "] is not exist")
		return
	}

	var tableDataSource DataSource
	if this_.DataSourceType != DataSourceTypeSql {
		fileName := targetTableName + "." + this_.DataSourceType.FileSuffix
		fileName, err = this_.getFileName(targetOwnerName, fileName)
		if err != nil {
			return
		}
		param := &DataSourceParam{
			Path:      fileName,
			SheetName: targetTableName,
			Dia:       this_.targetDialect,
		}
		tableDataSource = this_.DataSourceType.New(param)
		err = tableDataSource.WriteStart()
		if err != nil {
			return
		}
		defer func() {
			_ = tableDataSource.WriteEnd()
		}()
	}

	if this_.ExportStruct {
		err = this_.exportTableStruct(ownerDataSource, tableDataSource, tableDetail, targetOwnerName, targetTableName)
		if err != nil {
			return
		}
	}
	if this_.ExportData {
		err = this_.exportTableData(ownerDataSource, tableDataSource, tableDetail, targetOwnerName, targetTableName, columns)
		if err != nil {
			return
		}
	}
	success = true
	return
}

func (this_ *taskExport) exportTableStruct(ownerDataSource DataSource, tableDataSource DataSource, tableDetail *dialect.TableModel, targetOwnerName string, targetTableName string) (err error) {

	progress := &TaskProgress{
		Title: "导出表结构[" + tableDetail.OwnerName + "." + tableDetail.TableName + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
	}
	var oldOwnerName = tableDetail.OwnerName
	var oldTableName = tableDetail.TableName
	tableDetail.OwnerName = targetOwnerName
	if this_.AppendOwnerName {
		tableDetail.OwnerName = targetOwnerName
	} else {
		tableDetail.OwnerName = ""
	}
	tableDetail.TableName = targetTableName
	defer func() {
		tableDetail.OwnerName = oldOwnerName
		tableDetail.TableName = oldTableName

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

	if this_.IsStop {
		return
	}
	if this_.FormatIndexName != nil {
		for _, index := range tableDetail.IndexList {
			index.IndexName = this_.FormatIndexName(tableDetail.OwnerName, tableDetail.TableName, index)
		}
	}

	// 导出结构体

	lines, err := this_.targetDialect.TableCreateSql(this_.Param, tableDetail.OwnerName, tableDetail)

	for _, line := range lines {
		dataSourceData := &DataSourceData{
			HasSql: true,
			Sql:    line,
		}
		if ownerDataSource != nil {
			err = ownerDataSource.Write(dataSourceData)
			if err != nil {
				return
			}
		}
		if tableDataSource != nil {
			err = tableDataSource.Write(dataSourceData)
			if err != nil {
				return
			}
		}
	}

	return
}

func (this_ *taskExport) exportTableData(ownerDataSource DataSource, tableDataSource DataSource, tableDetail *dialect.TableModel, targetOwnerName string, targetTableName string, columns []*TaskExportColumn) (err error) {

	progress := &TaskProgress{
		Title: "导出表数据[" + tableDetail.OwnerName + "." + tableDetail.TableName + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
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

	if this_.IsStop {
		return
	}

	selectSqlInfo := "SELECT "
	var columnNames []string
	for _, one := range tableDetail.ColumnList {
		columnNames = append(columnNames, one.ColumnName)
	}
	selectSqlInfo += this_.dia.ColumnNamesPack(this_.Param, columnNames)
	selectSqlInfo += " FROM "

	selectSqlInfo += this_.dia.OwnerTablePack(this_.Param, tableDetail.OwnerName, tableDetail.TableName)

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
	for {

		if this_.IsStop {
			return
		}
		pageSql := this_.dia.PackPageSql(selectSqlInfo, pageSize, pageNo)
		dataList, err = DoQuery(this_.db, pageSql, nil)
		if err != nil {
			err = errors.New("query page query sql:" + pageSql + ",error:" + err.Error())
			return
		}
		pageNo += 1
		dataListCount := len(dataList)
		this_.countIncr(&this_.DataReadyCount, dataListCount)
		if dataListCount == 0 {
			break
		}
		var success bool
		success, err = this_.exportDataList(ownerDataSource, tableDataSource, dataList, targetOwnerName, targetTableName, tableDetail.ColumnList, columns)
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

func (this_ *taskExport) exportDataList(ownerDataSource DataSource, tableDataSource DataSource, dataList []map[string]interface{}, targetOwnerName string, targetTableName string, columnList []*dialect.ColumnModel, columns []*TaskExportColumn) (success bool, err error) {

	progress := &TaskProgress{
		Title: "导出数据[" + targetOwnerName + "." + targetTableName + "]",
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

	if this_.IsStop {
		return
	}
	var columnCache = make(map[string]*dialect.ColumnModel)
	for _, one := range columnList {
		columnCache[one.ColumnName] = one
	}
	var newColumnList = columnList
	if len(columns) > 0 {
		newColumnList = []*dialect.ColumnModel{}
		for _, one := range columns {
			if one.SourceName == "" {
				continue
			}
			column := columnCache[one.SourceName]
			targetName := one.TargetName
			if targetName == "" {
				targetName = one.SourceName
			}
			newColumn := &dialect.ColumnModel{}
			newColumn.ColumnName = targetName
			if targetName != one.SourceName || one.Value != "" {
				for _, data := range dataList {
					if one.Value != "" {
						data[targetName] = one.Value
					} else {
						data[targetName] = data[one.SourceName]
					}
				}
			}
			if column != nil {
				newColumn.ColumnDataType = column.ColumnDataType
				newColumn.ColumnDefault = column.ColumnDefault
				newColumn.ColumnLength = column.ColumnLength
				newColumn.ColumnPrecision = column.ColumnPrecision
				newColumn.ColumnScale = column.ColumnScale
			}
			newColumnList = append(newColumnList, newColumn)
		}
	}
	var sqlOwner = ""
	if this_.AppendOwnerName {
		sqlOwner = targetOwnerName
	}
	this_.Param.AppendSqlValue = new(bool)
	*this_.Param.AppendSqlValue = true
	sqlList, _, batchSqlList, _, err := this_.targetDialect.DataListInsertSql(this_.Param, sqlOwner, targetTableName, newColumnList, dataList)
	if err != nil {
		return
	}

	var lines []string
	if this_.ExportBatchSql {
		lines = batchSqlList
	} else {
		lines = sqlList
	}

	for _, line := range lines {
		dataSourceData := &DataSourceData{
			HasSql: true,
			Sql:    line,
		}
		if ownerDataSource != nil {
			err = ownerDataSource.Write(dataSourceData)
			if err != nil {
				return
			}
		}
		if tableDataSource != nil {
			err = tableDataSource.Write(dataSourceData)
			if err != nil {
				return
			}
		}
	}

	for _, data := range dataList {
		dataSourceData := &DataSourceData{
			HasData:    true,
			Data:       data,
			ColumnList: columnList,
		}
		if ownerDataSource != nil {
			err = ownerDataSource.Write(dataSourceData)
			if err != nil {
				return
			}
		}
		if tableDataSource != nil {
			err = tableDataSource.Write(dataSourceData)
			if err != nil {
				return
			}
		}
	}
	success = true

	return
}
