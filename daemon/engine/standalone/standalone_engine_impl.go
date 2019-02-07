package standalone

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/songxinjianqwe/scheduler/common"
	"os/exec"
	"sync"
	"time"
)

type StandAloneEngine struct {
	tasks sync.Map // key为task id，类型为string；value为task实例的指针
}

func (this *StandAloneEngine) Submit(task *common.Task) error {
	_, loaded := this.tasks.LoadOrStore(task.Id, task)
	if loaded {
		return errors.New("task id existed")
	}
	taskFunc := func(now time.Time) {
		cmd := exec.Command("/bin/bash", "-c", task.Script)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Error(err)
			task.Results = append(task.Results, err.Error())
		} else {
			log.Info(out.String())
			task.Results = append(task.Results, out.String())
		}
	}
	switch task.TaskType {
	case "delay":
		timer := time.NewTimer(task.Time)
		<-timer.C
		taskFunc(time.Now())
		timer.Stop()
	case "cron":
		ticker := time.NewTicker(task.Time)
		// 即使handlerFunc是在新的goroutine中运行的，但这里是阻塞循环，必须放在另一个goroutine中运行，否则请求无法返回
		go func() {
			for now := range ticker.C {
				taskFunc(now)
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
