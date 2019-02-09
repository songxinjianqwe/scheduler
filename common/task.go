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
	// task id, unique
	Id string `json:"id"`
	// task type, enum{delay,cron}
	TaskType string `json:"taskType"`
	// execute time
	Time time.Duration `json:"time"`
	// shell script
	Script string `json:"script"`
	// results
	Results []TaskResult `json:"results"`
	// status, enum
	Status TaskStatus `json:"status"`
	// last status updated time
	LastStatusUpdated time.Time `json:"lastStatusUpdated"`
	// task version, increase when data(Results,Status,LastStatusUpdated) changed.
	Version   int64        `json:"version"`
	timer     *time.Timer  `json:"-"`
	ticker    *time.Ticker `json:"-"`
	lock      *sync.Mutex  `json:"-"`
	watchCond *sync.Cond   `json:"-"`
}

func NewTask(id string, taskType string, dueTime time.Duration, script string) *Task {
	task := Task{}
	task.Id = id
	task.TaskType = taskType
	task.Time = dueTime
	task.Script = script
	task.Status = ToBeExecuted
	task.LastStatusUpdated = time.Now()
	task.Version = 0
	return &task
}

// 将lock和watch condition延迟初始化，在客户端NewTask时这两个是用不到的
func (this *Task) PopulateTaskTimer(timer *time.Timer) {
	this.timer = timer
	this.lock = new(sync.Mutex)
	this.watchCond = sync.NewCond(this.lock)
}

func (this *Task) PopulateTaskTicker(ticker *time.Ticker) {
	this.ticker = ticker
	this.lock = new(sync.Mutex)
	this.watchCond = sync.NewCond(this.lock)
}

// 这里需要考虑线程安全问题!
// 执行永远是在同一个goroutine中执行
// 但是stop会在另一个goroutine中执行
// 它们会并发地修改Task的内部状态
// 最好是能在task粒度上加一个互斥锁
// 另外是内存可见性问题：get能否读到最新的对象值？Go中没有volatile，要想保证读到最新的值，需要加锁或者atomic等
func (this *Task) Execute() {
	// 其实如果真的停止Timer或者Ticker，那么Execute是不会被执行的
	// 极端情况下，刚刚Stop，Execute就到期开始执行了，此时需要double check一下
	// 理论上Stop先执行，会加锁，然后停止计时器
	// 如果在这段时间的同时，计时器到期，开始Execute，会因为获取不到锁而阻塞，所以
	if this.checkIfAlreadyStoppedThenUpdateStatusAtomically() {
		// 这里是先检查是否已经停止，如果是，则退出
		// 如果不是，则将状态更新为执行中
		// 这个操作是原子执行的
		// 1) Execute，check获取锁，Stop等待，check完毕，状态更新为执行中，Stop拿到锁，发现状态为执行中，则停止失败
		// 2）Stop获取锁，Execute等待，状态更新为已停止，Execute拿到锁，发现状态为已停止，则退出执行
		return
	}

	// 非原子执行，这段代码可能会耗时过长，不需要加锁，不会修改状态
	cmd := exec.Command("/bin/bash", "-c", this.Script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	var result string
	if err != nil {
		log.Error(err)
		result = err.Error()
	} else {
		log.Infof("Task stdout: %s", out.String())
		result = out.String()
	}

	// 修改状态，原子执行
	this.appendResultAndUpdateStatusAtomically(result)
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
// @Atomically
func (this *Task) Stop() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.Status == Stopped {
		return fmt.Errorf("任务已经被停止")
	}
	if this.TaskType == "delay" {
		// 此时task的状态有可能是三种情况，只有ToBeExecuted状态下我们才能去stop
		if this.Status == ToBeExecuted {
			this.timer.Stop()
			this.updateStatus(Stopped)
			this.increaseVersionAndSignalListeners()
		} else {
			// 此时为Executing或Executed状态，无法停止
			return fmt.Errorf("任务状态为%s, 本次执行无法停止", this.Status.String())
		}
	} else if this.TaskType == "cron" {
		// 虽然Executing状态无法停止，但是可以让它这次执行完毕，在这里置状态为Stopped，正在执行的goroutine会检测到，然后不会把
		// 状态改回WaitForNextExecution
		isExecuting := this.Status == Executing
		this.ticker.Stop()
		this.updateStatus(Stopped)
		this.increaseVersionAndSignalListeners()
		if isExecuting {
			return fmt.Errorf("任务状态为%s, 本次执行无法停止, 但已停止后续执行", this.Status.String())
		}
	}
	return nil
}

// 返回一份拷贝
func (this *Task) GetLatest(watch bool, version int64) (Task, error) {
	// 如果是非watch，或者是watch，但版本不同，则都返回最新值
	this.lock.Lock()
	defer this.lock.Unlock()
	if !watch || this.Version != version {
		return *this, nil
	}
	for this.Version == version {
		// 此时版本相同，需要阻塞等待
		this.watchCond.Wait()
	}
	return *this, nil
}

// @Atomically
func (this *Task) checkIfAlreadyStoppedThenUpdateStatusAtomically() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.Status == Stopped {
		log.Error("任务[%s]已经被停止，不再执行", this.Id)
		return true
	}
	this.updateStatus(Executing)
	this.increaseVersionAndSignalListeners()
	return false
}

// @Atomically
func (this *Task) appendResultAndUpdateStatusAtomically(result string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Results = append(this.Results, NewTaskResult(time.Now(), result))
	// 执行完毕
	if this.TaskType == "delay" {
		this.updateStatus(Executed)
	} else if this.TaskType == "cron" {
		// 如果任务在执行时被停止，则发起停止命令的goroutine会将状态置为Stopped，此时就保留Stopped状态，不会将状态置为WaitForNextExecution
		// 如果任务正常执行，则在执行一次后，将状态置为WaitForNextExecution
		if this.Status != Stopped {
			this.updateStatus(WaitForNextExecution)
		}
	}
	this.increaseVersionAndSignalListeners()
}

// 必须放在lock里面调用
func (this *Task) updateStatus(status TaskStatus) {
	this.Status = status
	this.LastStatusUpdated = time.Now()
}

// 必须放在lock里面调用
func (this *Task) increaseVersionAndSignalListeners() {
	this.Version++
	this.watchCond.Broadcast()
}

func (this *Task) PrintMe() {
	fmt.Printf("Id: %s\n", this.Id)
	fmt.Printf("TaskType: %s\n", this.TaskType)
	fmt.Printf("Time: %s\n", this.Time)
	fmt.Printf("Script: %s\n", this.Script)
	fmt.Printf("Status: %s\n", this.Status.String())
	fmt.Printf("LastStatusUpdated: %s\n", this.LastStatusUpdated)
	fmt.Println("Results: ")
	for index, result := range this.Results {
		fmt.Printf("[%d]%s\n", index, result.Timestamp)
		fmt.Printf("%s\n", result.Result)
	}
}
