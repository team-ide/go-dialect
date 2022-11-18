package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"sync"
)

type Task struct {
	TaskId     string `json:"taskId"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
	UseTime    int64  `json:"useTime"`
	Error      string `json:"error"`
	PanicError string `json:"panicError"`
	IsEnd      bool   `json:"isEnd"`
	IsStop     bool   `json:"isStop"`

	DataCount    int `json:"dataCount"`
	ReadyCount   int `json:"readyDataCount"`
	SuccessCount int `json:"successCount"`
	ErrorCount   int `json:"errorCount"`

	countLock sync.Mutex

	Extend           interface{}     `json:"extend"`
	Errors           []string        `json:"errors"`
	TaskProgressList []*TaskProgress `json:"taskProgressList"`

	onProgress func(progress *TaskProgress)
	dia        dialect.Dialect
	db         *sql.DB
	do         func() (err error)
	Param      *dialect.ParamModel
}

type TaskProgress struct {
	Title string   `json:"title"`
	Infos []string `json:"infos"`
	Error string   `json:"error"`
}

var (
	taskCache     = make(map[string]*Task)
	taskCacheLock sync.Mutex
)

func GetTask(taskId string) (task *Task) {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	task = taskCache[taskId]
	return
}

func addTask(task *Task) {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	task.TaskId = dialect.UUID()
	taskCache[task.TaskId] = task
	return
}

func StopTask(taskId string) {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	task := taskCache[taskId]
	if task != nil {
		task.Stop()
	}
	return
}

func (this_ *Task) Start() (err error) {
	this_.IsStop = false
	addTask(this_)

	defer func() {
		if e := recover(); e != nil {
			this_.PanicError = fmt.Sprint(e)
			this_.Error = this_.PanicError
		}
		if err != nil {
			this_.Error = err.Error()
		}
		this_.EndTime = NowTime()
		this_.UseTime = this_.EndTime - this_.StartTime
		this_.IsEnd = true
	}()

	this_.StartTime = NowTime()
	if this_.do == nil {
		err = errors.New("has nothing to do")
		return
	}
	err = this_.do()
	if err != nil {
		return
	}
	return
}

func (this_ *Task) addProgress(progress *TaskProgress) {
	this_.TaskProgressList = append(this_.TaskProgressList, progress)
	if this_.onProgress != nil {
		this_.onProgress(progress)
	}
	return
}

func (this_ *Task) addError(err string) {
	this_.Errors = append(this_.Errors, err)
	return
}

func (this_ *Task) Stop() {
	this_.IsStop = true
}

func (this_ *Task) dataCountIncr() {
	this_.countIncr(&this_.DataCount)
}

func (this_ *Task) readyCountIncr() {
	this_.countIncr(&this_.ReadyCount)
}

func (this_ *Task) successCountIncr() {
	this_.countIncr(&this_.SuccessCount)
}

func (this_ *Task) errorCountIncr() {
	this_.countIncr(&this_.ErrorCount)
}

func (this_ *Task) countIncr(count *int) {
	this_.countLock.Lock()
	defer this_.countLock.Unlock()
	*count++
	return
}
