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
		dia: dia,
		db:  db,
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
	Owners []*TaskExportOwner `json:"owners"`

	DataSourceType  *DataSourceType `json:"dataSourceType"`
	BatchNumber     int             `json:"batchNumber"`
	ExportStructure bool            `json:"exportStructure"`
	ExportData      bool            `json:"exportData"`
	ExportBatchSql  bool            `json:"exportBatchSql"`
	ContinueIsError bool            `json:"continueIsError"`
	Dir             string          `json:"dir"`
	FormatIndexName func(ownerName string, tableName string, index *dialect.IndexModel) string
}

type TaskExportOwner struct {
	SourceName     string             `json:"sourceName"`
	TargetName     string             `json:"targetName"`
	SkipTableNames []string           `json:"skipTableNames"`
	Tables         []*TaskExportTable `json:"tables"`
}

type TaskExportTable struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
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
	if len(this_.Owners) == 0 {
		return
	}
	for _, owner := range this_.Owners {
		err = this_.exportOwner(owner)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskExport) getFileName(name string) (fileName string, err error) {
	var exist bool
	if this_.Dir != "" {
		exist, err = PathExists(this_.Dir)
		if err != nil {
			return
		}
		if !exist {
			err = os.MkdirAll(this_.Dir, 0777)
			if err != nil {
				return
			}
		}
		fileName = this_.Dir + string(os.PathSeparator)
	}
	fileName += name
	return
}

func (this_ *taskExport) exportOwner(owner *TaskExportOwner) (err error) {
	progress := &TaskProgress{
		Title: "导出[" + owner.SourceName + "]",
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
			tables = append(tables, &TaskExportTable{
				SourceName: one.Name,
			})
		}
	}

	ownerName := owner.TargetName
	if ownerName == "" {
		ownerName = owner.SourceName
	}

	var ownerDataSource DataSource
	if this_.DataSourceType.OwnerFileName != nil {
		fileName := this_.DataSourceType.OwnerFileName(ownerName)
		fileName, err = this_.getFileName(fileName)
		if err != nil {
			return
		}
		param := &DataSourceParam{
			Path:      fileName,
			SheetName: ownerName,
		}
		ownerDataSource = this_.DataSourceType.New(param)
		err = ownerDataSource.WriteStart()
		go func() {
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
				continue
			}
		}

		err = this_.exportTable(ownerDataSource, owner.SourceName, table.SourceName, owner.TargetName, table.TargetName)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskExport) exportTable(ownerDataSource DataSource, sourceOwnerName string, sourceTableName string, targetOwnerName string, targetTableName string) (err error) {
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
		}

		if this_.ContinueIsError {
			err = nil
		}
	}()

	this_.addProgress(progress)

	tableDetail, err := TableDetail(this_.db, this_.dia, sourceOwnerName, sourceTableName)
	if err != nil {
		return
	}
	if tableDetail == nil {
		err = errors.New("source db table [" + sourceOwnerName + "." + sourceTableName + "] is not exist")
		return
	}

	var tableDataSource DataSource
	if this_.DataSourceType.TableFileName != nil {
		fileName := this_.DataSourceType.TableFileName(targetOwnerName, targetTableName)
		fileName, err = this_.getFileName(fileName)
		if err != nil {
			return
		}
		param := &DataSourceParam{
			Path:      fileName,
			SheetName: targetTableName,
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

	if this_.ExportStructure {
		err = this_.exportTableStructure(ownerDataSource, tableDataSource, tableDetail, targetOwnerName, targetTableName)
		if err != nil {
			return
		}
	}
	if this_.ExportData {
		err = this_.exportTableData(ownerDataSource, tableDataSource, tableDetail, targetOwnerName, targetTableName)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskExport) exportTableStructure(ownerDataSource DataSource, tableDataSource DataSource, tableDetail *dialect.TableModel, targetOwnerName string, targetTableName string) (err error) {

	progress := &TaskProgress{
		Title: "导出表结构[" + tableDetail.OwnerName + "." + tableDetail.Name + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
	}
	var oldOwnerName = tableDetail.OwnerName
	var oldTableName = tableDetail.Name
	tableDetail.OwnerName = targetOwnerName
	tableDetail.Name = targetTableName
	defer func() {
		tableDetail.OwnerName = oldOwnerName
		tableDetail.Name = oldTableName

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
		for _, index := range tableDetail.IndexList {
			index.Name = this_.FormatIndexName(tableDetail.OwnerName, tableDetail.Name, index)
		}
	}

	// 导出结构体

	lines, err := this_.targetDialect.TableCreateSql(tableDetail.OwnerName, tableDetail)

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

func (this_ *taskExport) exportTableData(ownerDataSource DataSource, tableDataSource DataSource, tableDetail *dialect.TableModel, targetOwnerName string, targetTableName string) (err error) {

	progress := &TaskProgress{
		Title: "导出表数据[" + tableDetail.OwnerName + "." + tableDetail.Name + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
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
	for _, one := range tableDetail.ColumnList {
		columnNames = append(columnNames, one.Name)
	}
	selectSqlInfo += this_.dia.PackColumns(columnNames)
	selectSqlInfo += " FROM "
	if tableDetail.OwnerName != "" {
		selectSqlInfo += this_.dia.PackOwner(tableDetail.OwnerName) + "."
	}
	selectSqlInfo += this_.dia.PackTable(tableDetail.Name)

	list, err := DoQuery(this_.db, selectSqlInfo)
	if err != nil {
		return
	}

	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}

	dataListGroup := SplitArrayMap(list, batchNumber)

	for _, dataList := range dataListGroup {
		err = this_.exportDataList(ownerDataSource, tableDataSource, dataList, targetOwnerName, targetTableName, tableDetail.ColumnList)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskExport) exportDataList(ownerDataSource DataSource, tableDataSource DataSource, dataList []map[string]interface{}, targetOwnerName string, targetTableName string, columnList []*dialect.ColumnModel) (err error) {

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

	sqlList, batchSqlList, err := InsertDataListSql(this_.targetDialect, targetOwnerName, targetTableName, columnList, dataList)
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

	return
}
