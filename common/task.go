package common

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"sync"
	"time"
)

type Task struct {
	Id                string        `json:"id"`
	TaskType          string        `json:"taskType"`
	Time              time.Duration `json:"time"`
	Script            string        `json:"script"`
	Results           []*TaskResult `json:"results"`
	Status            TaskStatus    `json:"status"`
	LastStatusUpdated time.Time     `json:"lastStatusUpdated"`
	Timer             *time.Timer   `json:"-"`
	Ticker            *time.Ticker  `json:"-"`
	lock              *sync.Mutex   `json:"-"`
}

func NewTask(id string, taskType string, time time.Duration, script string) *Task {
	task := Task{}
	task.Id = id
	task.TaskType = taskType
	task.Time = time
	task.Script = script
	task.Status = ToBeExecuted
	task.lock = new(sync.Mutex)
	return &task
}

// 这里是否需要考虑线程安全问题？
// 执行永远是在同一个goroutine中执行
// 但是stop和delete会在另一个goroutine中执行
// 它们会并发地修改Task的内部状态
// 最好是能在task粒度上加一个互斥锁
// 另外是内存可见性问题：get能否读到最新的对象值？
func (this *Task) Execute() {
	this.lock.Lock()
	// ~~~ATOMIC BLOCK BEGIN
	// 其实如果真的停止Timer或者Ticker，那么Execute是不会被执行的
	// 极端情况下，刚刚Stop，Execute就到期开始执行了，此时需要double check一下
	if this.Status == Stopped {
		log.Warn("任务[%s]已经被停止，不再执行", this.Id)
		this.lock.Unlock()
		return
	}
	// 开始执行
	this.Status = Executing
	this.LastStatusUpdated = time.Now()
	this.lock.Unlock()
	// ~~~ATOMIC BLOCK END

	// 执行，这段代码可能会耗时过长，不要阻塞Stop
	now := time.Now()
	cmd := exec.Command("/bin/bash", "-c", this.Script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	// ~~~ATOMIC BLOCK BEGIN
	this.lock.Lock()
	if err != nil {
		log.Error(err)
		this.Results = append(this.Results, NewTaskResult(now, err.Error()))
	} else {
		log.Infof("Task stdout: %s", out.String())
		this.Results = append(this.Results, NewTaskResult(now, out.String()))
	}
	// 执行完毕
	if this.TaskType == "delay" {
		this.Status = Executed
		this.LastStatusUpdated = time.Now()
		this.Timer.Stop()
	} else if this.TaskType == "cron" {
		// 如果任务在执行时被停止，则发起停止命令的goroutine会将状态置为Stopped，此时就保留Stopped状态，不会将状态置为WaitForNextExecution
		// 如果任务正常执行，则在执行一次后，将状态置为WaitForNextExecution
		if this.Status != Stopped {
			this.Status = WaitForNextExecution
			this.LastStatusUpdated = time.Now()
		}
	}
	this.lock.Unlock()
	// ~~~ATOMIC BLOCK END
}

/**
	如果是延迟任务，则：
		1) 可以在任务开始前停止，则任务不会被执行，终态为Stopped
		2）如果任务已经开始执行，则无法停止，且报错，终态为Executed
		3）如果任务已经执行完毕，则无法停止，且报错，终态为Executed
	如果是定时任务，则：
		1）可以在任务开始前停止，则任务一次都不会被执行，终态为Stopped
		2）如果在某一次任务开始执行后停止，则任务本次执行不会被中止，且会保存本次执行结果，终态为Stopped
		3）如果在某一个任务任务执行后，下一次任务执行前停止，则下次执行会被跳过，终态为Stopped
 */
func (this *Task) Stop() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.Status == Stopped {
		return fmt.Errorf("任务已经被停止")
	}
	if this.TaskType == "delay" {
		// 此时task的状态有可能是三种情况，只有ToBeExecuted状态下我们才能去stop
		if this.Status == ToBeExecuted {
			this.Timer.Stop()
			this.Status = Stopped
			this.LastStatusUpdated = time.Now()
		} else {
			// 此时为Executing或Executed状态，无法停止
			return fmt.Errorf("任务状态为%s, 本次执行无法停止", this.Status.String())
		}
	} else if this.TaskType == "cron" {
		// 虽然Executing状态无法停止，但是可以让它这次执行完毕，在这里置状态为Stopped，正在执行的goroutine会检测到，然后不会把
		// 状态改回WaitForNextExecution
		isExecuting := this.Status == Executing
		this.Ticker.Stop()
		this.Status = Stopped
		this.LastStatusUpdated = time.Now()
		if isExecuting {
			return fmt.Errorf("任务状态为%s, 本次执行无法停止, 但已停止后续执行", this.Status.String())
		}
	}
	return nil
}
