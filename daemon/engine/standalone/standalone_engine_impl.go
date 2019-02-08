package standalone

import (
	"errors"
	"github.com/songxinjianqwe/scheduler/common"
	"sync"
	"time"
)

type StandAloneEngine struct {
	tasks sync.Map // key为task id，类型为string；value为task实例的指针
}

func (this *StandAloneEngine) Stop(id string) error {
	task, ok := this.tasks.Load(id)
	if !ok {
		return errors.New("task id not existed")
	}
	return (task.(*common.Task)).Stop()
}

func (this *StandAloneEngine) Delete(id string) error {
	// 先stop
	// 再delete
	task, ok := this.tasks.Load(id)
	if !ok {
		return errors.New("task id not existed")
	}
	(task.(*common.Task)).Stop()
	// 删除之后，有可能还会更新一次Result，但是无所谓，之后会回收掉
	this.tasks.Delete(id)
	return nil
}

func (this *StandAloneEngine) Submit(task *common.Task) error {
	_, loaded := this.tasks.LoadOrStore(task.Id, task)
	if loaded {
		return errors.New("task id existed")
	}

	switch task.TaskType {
	case "delay":
		timer := time.NewTimer(task.Time)
		task.Timer = timer
		go func() {
			<-timer.C
			task.Execute()
		}()
	case "cron":
		ticker := time.NewTicker(task.Time)
		task.Ticker = ticker
		// 即使handlerFunc是在新的goroutine中运行的，但这里是阻塞循环，必须放在另一个goroutine中运行，否则请求无法返回
		go func() {
			for range ticker.C {
				task.Execute()
			}
		}()
	}
	return nil
}

func (this *StandAloneEngine) Get(id string) (*common.Task, error) {
	value, ok := this.tasks.Load(id)
	if !ok {
		return nil, errors.New("task id not existed")
	}
	return value.(*common.Task), nil
}

func (this *StandAloneEngine) List() ([]*common.Task, error) {
	var tasks []*common.Task
	this.tasks.Range(func(key, value interface{}) bool {
		tasks = append(tasks, value.(*common.Task))
		return true
	})
	if tasks == nil {
		tasks = make([]*common.Task, 0)
	}
	return tasks, nil
}
