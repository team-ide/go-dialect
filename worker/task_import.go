package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"os"
	"strings"
)

func NewTaskImport(db *sql.DB, dia dialect.Dialect, taskImportParam *TaskImportParam) (res *taskImport) {
	if taskImportParam.DataSourceType == nil {
		taskImportParam.DataSourceType = DataSourceTypeSql
	}
	task := &Task{
		dia: dia,
		db:  db,
	}
	res = &taskImport{
		Task:            task,
		TaskImportParam: taskImportParam,
	}
	task.do = res.do
	return
}

type TaskImportParam struct {
	Owners []*TaskImportOwner `json:"owners"`

	DataSourceType  *DataSourceType `json:"dataSourceType"`
	BatchNumber     int             `json:"batchNumber"`
	ContinueIsError bool            `json:"continueIsError"`
	FormatIndexName func(ownerName string, tableName string, index *dialect.IndexModel) string
}

type TaskImportOwner struct {
	Name           string             `json:"name"`
	Dir            string             `json:"dir"`
	Path           string             `json:"path"`
	SkipTableNames []string           `json:"skipTableNames"`
	Tables         []*TaskImportTable `json:"tables"`
}

type TaskImportTable struct {
	SourceName string `json:"sourceName"`
	TargetName string `json:"targetName"`
}

type taskImport struct {
	*Task
	*TaskImportParam
}

func (this_ *taskImport) do() (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	if len(this_.Owners) == 0 {
		return
	}
	for _, owner := range this_.Owners {
		err = this_.importOwner(owner)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *taskImport) getFileName(dir string, name string) (fileName string, err error) {
	var exist bool

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

func (this_ *taskImport) importOwner(owner *TaskImportOwner) (err error) {
	progress := &TaskProgress{
		Title: "导入[" + owner.Name + "]",
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

	ownerName := owner.Name

	var ownerDataSource DataSource
	if !this_.DataSourceType.OwnerIsDir {
		fileName := owner.Path
		fileName, err = this_.getFileName("", fileName)
		if err != nil {
			return
		}
		param := &DataSourceParam{
			Path:      fileName,
			SheetName: ownerName,
			Dia:       this_.dia,
		}
		ownerDataSource = this_.DataSourceType.New(param)
		err = ownerDataSource.ReadStart()
		if err != nil {
			return
		}
		defer func() {
			_ = ownerDataSource.ReadEnd()
		}()
		err = ownerDataSource.Read(nil, func(data *DataSourceData) (err error) {
			if data.HasSql {
				_, err = DoExec(this_.db, []string{data.Sql})
				if err != nil {
					err = errors.New("sql:" + data.Sql + " exec error," + err.Error())
					return
				}
			}
			return
		})
		if err != nil {
			return
		}
	} else {
		dir := owner.Dir
		var ds []os.DirEntry
		ds, err = os.ReadDir(dir)
		if err != nil {
			return
		}

		for _, d := range ds {
			path := dir + string(os.PathSeparator) + d.Name()
			var f os.FileInfo
			f, err = os.Lstat(path)
			if err != nil {
				return
			}
			if f.IsDir() {
				continue
			}
			tableName := d.Name()
			if strings.Index(tableName, ".") > 0 {
				tableName = tableName[0:strings.Index(tableName, ".")]
			}
			if len(owner.SkipTableNames) > 0 {
				var skip bool
				for _, skipTableName := range owner.SkipTableNames {
					if strings.EqualFold(tableName, skipTableName) {
						skip = true
					}
				}
				if skip {
					continue
				}
			}
			err = this_.importTable(owner.Name, path, tableName, tableName)
			if err != nil {
				return
			}

		}

	}

	return
}

func (this_ *taskImport) importTable(ownerName string, path string, sourceTableName string, targetTableName string) (err error) {
	if targetTableName == "" {
		targetTableName = sourceTableName
	}
	progress := &TaskProgress{
		Title: "导入[" + ownerName + "." + sourceTableName + "] 到 [" + ownerName + "." + targetTableName + "]",
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

	tableDetail, err := TableDetail(this_.db, this_.dia, ownerName, targetTableName)
	if err != nil {
		return
	}
	if tableDetail == nil {
		err = errors.New("source db table [" + ownerName + "." + targetTableName + "] is not exist")
		return
	}

	var tableDataSource DataSource
	if this_.DataSourceType.OwnerIsDir {
		param := &DataSourceParam{
			Path:      path,
			SheetName: targetTableName,
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
	}

	//if this_.ImportStructure {
	//	err = this_.exportTableStructure(ownerDataSource, tableDataSource, tableDetail, targetOwnerName, targetTableName)
	//	if err != nil {
	//		return
	//	}
	//}
	err = this_.importTableData(tableDataSource, tableDetail, ownerName, targetTableName)
	if err != nil {
		return
	}
	return
}
func (this_ *taskImport) importTableData(tableDataSource DataSource, tableDetail *dialect.TableModel, targetOwnerName string, targetTableName string) (err error) {

	progress := &TaskProgress{
		Title: "导入表数据[" + tableDetail.OwnerName + "." + tableDetail.Name + "] 到 [" + targetOwnerName + "." + targetTableName + "]",
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
	err = tableDataSource.Read(tableDetail.ColumnList, func(data *DataSourceData) (err error) {
		if data.HasData && data.Data != nil {
			dataList = append(dataList, data.Data)
			if len(dataList) >= batchNumber {
				err = this_.importDataList(dataList, targetOwnerName, targetTableName, tableDetail.ColumnList)
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
		err = this_.importDataList(dataList, targetOwnerName, targetTableName, tableDetail.ColumnList)
		dataList = make([]map[string]interface{}, 0)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *taskImport) importDataList(dataList []map[string]interface{}, ownerName string, tableName string, columnList []*dialect.ColumnModel) (err error) {

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

	_, sqlList, err := InsertDataListSql(this_.dia, ownerName, tableName, columnList, dataList)
	if err != nil {
		return
	}
	var errSql string
	errSql, err = DoExec(this_.db, sqlList)
	if err != nil {
		if errSql != "" {
			err = errors.New("sql:" + errSql + " exec error," + err.Error())
		}
		return
	}
	return
}
