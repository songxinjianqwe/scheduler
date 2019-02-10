package main

import (
	"github.com/satori/go.uuid"
	"github.com/songxinjianqwe/scheduler/cli/client"
	"github.com/songxinjianqwe/scheduler/common"
	"github.com/songxinjianqwe/scheduler/daemon/server"
	"testing"
	"time"
)

var schedulerClient *client.SchedulerClient

func TestMain(m *testing.M) {
	go server.Run()
	time.Sleep(time.Second * 2)
	schedulerClient, _ = client.NewClient()
	m.Run()
}

func TestList(t *testing.T) {
	tasks, err := schedulerClient.List()
	if err != nil {
		t.Error(err)
	}
	if len(tasks) != 0 {
		t.Errorf("任务数量必须为0: %#v", tasks)
	}
}

func TestSubmit(t *testing.T) {
	id := "test_delay_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	delayTime := time.Second * 2
	task := common.NewTask(
		id,
		"delay",
		delayTime,
		"echo "+result,
	)
	err := schedulerClient.Submit(task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	time.Sleep(delayTime + time.Second)
	tasks, err := schedulerClient.List()
	if err != nil {
		t.Error(err)
	}
	if len(tasks) != 1 {
		t.Errorf("任务数量必须为0: %#v", tasks)
	}
	existed := false
	for _, task := range tasks {
		if task.Id == id {
			existed = true
			if task.Status != common.Executed {
				t.Fatalf("任务状态错误: 预期: %s，实际：%s", common.Executed, task.Status)
			}
		}
	}
	if !existed {
		t.Errorf("List()未找到该任务: %#v", task)
	}
	latestTask, err := schedulerClient.Get(id, false, 0)
	if err != nil {
		t.Error(err)
	}
	if latestTask.Status != common.Executed {
		t.Fatalf("任务状态错误: 预期: %s，实际：%s", common.Executed, task.Status)
	}
}

func TestDelete(t *testing.T) {
	id := "test_delay_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	delayTime := time.Second * 2
	task := common.NewTask(
		id,
		"delay",
		delayTime,
		"echo "+result,
	)
	err := schedulerClient.Submit(task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	time.Sleep(delayTime + time.Second)
	latestTask, err := schedulerClient.Get(id, false, 0)
	if err != nil {
		t.Error(err)
	}
	if latestTask.Status != common.Executed {
		t.Errorf("任务状态错误: 预期: %s，实际：%s", common.Executed, task.Status)
	}
	err = schedulerClient.Delete(id)
	if err != nil {
		t.Error(err)
	}
	latestTask, err = schedulerClient.Get(id, false, 0)
	if err == nil {
		t.Errorf("预期获取不到该任务: %#v", task)
	}
}
