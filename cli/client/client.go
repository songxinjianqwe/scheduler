package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/songxinjianqwe/scheduler/common"
	"io/ioutil"
	"net/http"
)

const (
	ServerAddr              = "http://localhost:8865/api"
)

type SchedulerClient struct {
	httpClient *http.Client
}

/**
	构造一个Client实例
 */
func NewClient() (c *SchedulerClient, err error) {
	schedulerClient := &SchedulerClient{}
	schedulerClient.httpClient = &http.Client{}
	return schedulerClient, nil
}

/**
	获取当前的所有任务
 */
func (this *SchedulerClient) List() ([]common.Task, error) {
	request, _ := http.NewRequest(http.MethodGet, ServerAddr+"/tasks", nil)
	response, err := this.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var tasks []common.Task
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = make([]common.Task, 0)
	}
	return tasks, nil
}

/**
	提交一个任务
 */
func (this *SchedulerClient) Submit(task *common.Task) error {
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}
	request, _ := http.NewRequest(http.MethodPost, ServerAddr+"/tasks", bytes.NewBuffer(taskBytes))
	response, err := this.httpClient.Do(request)
	defer response.Body.Close()
	return err
}

/**
	获取一个任务当前的执行情况
 */
func (this *SchedulerClient) Get(id string, watch bool) (*common.Task, error) {
	request, _ := http.NewRequest(http.MethodGet, ServerAddr + "/tasks/" + id, nil)
	response, err := this.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}
	var task common.Task
	return &task, json.Unmarshal(body, &task)
}
