package common

import "time"

const (
	TO_BE_EXECUTED = iota
	EXECUTED
	EXECUTING
	WAIT_FOR_NEXT_EXECTION
)

type Task struct {
	Id string `json:"id"`
	TaskType string `json:"taskType"`
	Time time.Duration `json:"time"`
	Script string `json:"script"`
	Results []string `json:"results"`
	Status string
 }

func NewTask(id string, taskType string, time time.Duration, script string) *Task {
	task := Task{}
	task.Id = id
	task.TaskType = taskType
	task.Time = time
	task.Script = script
	return &task
}