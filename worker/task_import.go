package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"path/filepath"
	"strings"
)

func NewTaskImport(db *sql.DB, dia dialect.Dialect, newDb func(owner *TaskImportOwner) (db *sql.DB, err error), taskImportParam *TaskImportParam) (res *taskImport) {
	if taskImportParam.DataSourceType == nil {
		taskImportParam.DataSourceType = DataSourceTypeSql
	}
	task := &Task{
		dia:        dia,
		db:         db,
		onProgress: taskImportParam.OnProgress,
	}
	res = &taskImport{
		Task:            task,
		TaskImportParam: taskImportParam,
		newDb:           newDb,
	}
	task.do = res.do
	return
}

type TaskImportParam struct {
	Owners []*TaskImportOwner `json:"owners"`

	DataSourceType        *DataSourceType `json:"dataSourceType"`
	BatchNumber           int             `json:"batchNumber"`
	OwnerCreateIfNotExist bool            `json:"ownerCreateIfNotExist"`
	ErrorContinue         bool            `json:"errorContinue"`

	FormatIndexName func(ownerName string, tableName string, index *dialect.IndexModel) string `json:"-"`
	OnProgress      func(progress *TaskProgress)
}

type TaskImportOwner struct {
	Name           string             `json:"name"`
	Path           string             `json:"path"`
	SkipTableNames []string           `json:"skipTableNames"`
	Tables         []*TaskImportTable `json:"tables"`
	Username       string             `json:"username"`
	Password       string             `json:"password"`
}

type TaskImportTable struct {
	Name    string              `json:"name"`
	Path    string              `json:"path"`
	Columns []*TaskImportColumn `json:"columns"`
}

type TaskImportColumn struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type taskImport struct {
	*Task
	*TaskImportParam `json:"-"`
	newDb            func(owner *TaskImportOwner) (db *sql.DB, err error)
}

func (this_ *taskImport) do() (err error) {

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
		success, err = this_.importOwner(owner)
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

func (this_ *taskImport) importOwner(owner *TaskImportOwner) (success bool, err error) {
	progress := &TaskProgress{
		Title: "导入[" + owner.Name + "]",
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

	ownerName := owner.Name

	//
	if this_.OwnerCreateIfNotExist {
		var ownerOne *dialect.OwnerModel
		ownerOne, err = OwnerSelect(this_.db, this_.dia, this_.Param, ownerName)
		if err != nil {
			return
		}
		if ownerOne == nil {
			this_.addProgress(&TaskProgress{
				Title: "导入[" + owner.Name + "] 不存在，创建",
			})
			_, err = OwnerCreate(this_.db, this_.dia, this_.Param, &dialect.OwnerModel{
				OwnerName:             ownerName,
				OwnerPassword:         owner.Password,
				OwnerCharacterSetName: "utf8mb4",
			})
			if err != nil {
				return
			}
		} else {
			this_.addProgress(&TaskProgress{
				Title: "导入[" + owner.Name + "] 存在",
			})
		}
	}

	workDb, err := this_.newDb(owner)
	if err != nil {
		return
	}

	if owner.Path != "" {
		var exists bool
		exists, err = PathExists(owner.Path)
		if err != nil {
			return
		}
		if !exists {
			err = errors.New("import [" + ownerName + "] path [" + owner.Path + "] not exists.")
			return
		}
	}

	if this_.DataSourceType == DataSourceTypeSql {
		if owner.Path != "" {
			var isDir bool
			isDir, err = PathIsDir(owner.Path)
			if err != nil {
				return
			}
			if !isDir {
				_, err = this_.importSql(workDb, ownerName, owner.Path)
				if err != nil {
					return
				}
			}
		}

	}
	tables := owner.Tables
	this_.countIncr(&this_.TableCount, len(tables))

	for _, table := range tables {
		tableName := table.Name
		if len(owner.SkipTableNames) > 0 {
			var skip bool
			for _, skipTableName := range owner.SkipTableNames {
				if strings.EqualFold(tableName, skipTableName) {
					skip = true
				}
			}
			if skip {
				this_.countIncr(&this_.TableSuccessCount, 1)
				continue
			}
		}
		var success_ bool
		success_, err = this_.importTable(workDb, owner.Name, tableName, owner.Path, table.Path, table.Columns)
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

func (this_ *taskImport) importSql(workDb *sql.DB, ownerName string, path string) (success bool, err error) {
	progress := &TaskProgress{
		Title: "导入[" + ownerName + "]",
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

	param := &DataSourceParam{
		Path:      path,
		SheetName: ownerName,
		Dia:       this_.dia,
	}
	ownerDataSource := this_.DataSourceType.New(param)
	err = ownerDataSource.ReadStart()
	if err != nil {
		return
	}
	defer func() {
		_ = ownerDataSource.ReadEnd()
	}()
	err = ownerDataSource.Read(nil, func(data *DataSourceData) (err error) {

		if this_.IsStop {
			return
		}

		if data.HasSql {
			this_.countIncr(&this_.DataCount, 1)

			var result sql.Result
			result, err = DoExec(workDb, data.Sql, nil)
			if err != nil {
				this_.countIncr(&this_.DataErrorCount, 1)
				if !this_.ErrorContinue {
					err = errors.New("sql:" + data.Sql + " exec error," + err.Error())
					return
				}
				err = nil
			}
			rowsAffected, _ := result.RowsAffected()
			this_.countIncr(&this_.DataSuccessCount, int(rowsAffected))
		}
		return
	})
	if err != nil {
		return
	}
	success = true
	return
}

func (this_ *taskImport) importTable(workDb *sql.DB, ownerName string, tableName string, ownerPath string, tablePath string, columns []*TaskImportColumn) (success bool, err error) {

	progress := &TaskProgress{
		Title: "导入[" + ownerName + "." + tableName + "]",
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

	if tablePath == "" {
		if ownerPath == "" {
			err = errors.New("import [" + ownerName + "." + tableName + "] table path is empty.")
			return
		}
		tablePath = ownerPath + string(filepath.Separator) + tableName + "." + this_.DataSourceType.FileSuffix

		var exists bool
		exists, err = PathExists(tablePath)
		if err != nil {
			return
		}
		if !exists {
			err = errors.New("import [" + ownerName + "." + tableName + "] path [" + tablePath + "] not exists.")
			return
		}
		var isDir bool
		isDir, err = PathIsDir(tablePath)
		if err != nil {
			return
		}
		if isDir {
			err = errors.New("import [" + ownerName + "." + tableName + "] path [" + tablePath + "] is dir.")
			return
		}
	}

	if this_.DataSourceType == DataSourceTypeSql {
		_, err = this_.importSql(workDb, ownerName, tablePath)
		if err != nil {
			return
		}
	} else {
		var tableDetail *dialect.TableModel
		tableDetail, err = TableDetail(workDb, this_.dia, this_.Param, ownerName, tableName, false)
		if err != nil {
			return
		}
		if tableDetail == nil {
			err = errors.New("source db table [" + ownerName + "." + tableName + "] is not exist")
			return
		}

		var tableDataSource DataSource
		param := &DataSourceParam{
			Path:      tablePath,
			SheetName: tableName,
			Dia:       this_.dia,
		}
		tableDataSource = this_.DataSourceType.New(param)
		err = tableDataSource.ReadStart()
		if err != nil {
			return
		}
		defer func() {
			_ = tableDataSource.ReadEnd()
		}()

		//if this_.ImportStructure {
		//	err = this_.exportTableStructure(ownerDataSource, tableDataSource, tableDetail, targetOwnerName, targetTableName)
		//	if err != nil {
		//		return
		//	}
		//}
		err = this_.importTableData(workDb, tableDataSource, tableDetail, ownerName, tableName, columns)
		if err != nil {
			return
		}
	}

	success = true
	return
}
func (this_ *taskImport) importTableData(workDb *sql.DB, tableDataSource DataSource, tableDetail *dialect.TableModel, targetOwnerName string, targetTableName string, columns []*TaskImportColumn) (err error) {

	progress := &TaskProgress{
		Title: "导入表数据[" + tableDetail.OwnerName + "." + tableDetail.TableName + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
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

	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}
	var columnCache = make(map[string]*dialect.ColumnModel)
	for _, one := range tableDetail.ColumnList {
		columnCache[one.ColumnName] = one
	}
	var newColumnList = tableDetail.ColumnList
	if len(columns) > 0 {
		newColumnList = []*dialect.ColumnModel{}
		for _, one := range columns {
			column := columnCache[one.Name]
			newColumn := &dialect.ColumnModel{}
			newColumn.ColumnName = one.Name
			if column != nil {
				newColumn.ColumnDataType = column.ColumnDataType
				newColumn.ColumnDefault = column.ColumnDefault
				newColumn.ColumnLength = column.ColumnLength
				newColumn.ColumnDecimal = column.ColumnDecimal
			}
			newColumnList = append(newColumnList, newColumn)
		}
	}

	var dataList []map[string]interface{}

	err = tableDataSource.Read(newColumnList, func(data *DataSourceData) (err error) {

		if this_.IsStop {
			return
		}
		if data.HasData && data.Data != nil {
			dataList = append(dataList, data.Data)
			this_.countIncr(&this_.DataCount, 1)
			if len(dataList) >= batchNumber {
				err = this_.importDataList(workDb, dataList, targetOwnerName, targetTableName, newColumnList)
				dataList = make([]map[string]interface{}, 0)
				if err != nil {
					return
				}
			}

		}
		return
	})
	if err != nil {
		return
	}
	if len(dataList) >= 0 {

		if this_.IsStop {
			return
		}

		err = this_.importDataList(workDb, dataList, targetOwnerName, targetTableName, newColumnList)
		dataList = make([]map[string]interface{}, 0)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskImport) importDataList(workDb *sql.DB, dataList []map[string]interface{}, ownerName string, tableName string, columnList []*dialect.ColumnModel) (err error) {

	dataListCount := len(dataList)
	progress := &TaskProgress{
		Title: "插入数据[" + ownerName + "." + tableName + "]",
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
			this_.countIncr(&this_.DataErrorCount, dataListCount)
		} else {
			this_.countIncr(&this_.DataSuccessCount, dataListCount)
		}

		if this_.ErrorContinue {
			err = nil
		}
	}()

	this_.addProgress(progress)

	if this_.IsStop {
		return
	}

	this_.countIncr(&this_.DataReadyCount, dataListCount)

	var newColumnList []*dialect.ColumnModel
	for _, one := range columnList {
		if one.ColumnName != "" {
			newColumnList = append(newColumnList, one)
		}
	}

	_, sqlList, err := this_.dia.InsertDataListSql(this_.Param, ownerName, tableName, newColumnList, dataList)
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
	return
}
