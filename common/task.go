package common

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
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
	Timer             *time.Timer    `json:"-"`
	Ticker            *time.Ticker   `json:"-"`
}

func NewTask(id string, taskType string, time time.Duration, script string) *Task {
	task := Task{}
	task.Id = id
	task.TaskType = taskType
	task.Time = time
	task.Script = script
	task.Status = ToBeExecuted
	return &task
}

func (this *Task) Execute() {
	if this.Status == Stopped {
		return
	}
	// 开始执行
	this.Status = Executing
	this.LastStatusUpdated = time.Now()
	// 执行
	now := time.Now()
	cmd := exec.Command("/bin/bash", "-c", this.Script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
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
	}else if this.TaskType == "cron" {
		if this.Status != Stopped {
			this.Status = WaitForNextExecution
			this.LastStatusUpdated = time.Now()
		}
	}
}

func (this *Task) Stop() error {
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
		executing := this.Status == Executing
		this.Ticker.Stop()
		this.Status = Stopped
		this.LastStatusUpdated = time.Now()
		if executing {
			return fmt.Errorf("任务状态为%s, 本次执行无法停止, 但已停止后续执行", this.Status.String())
		}
	}
	return nil
}
