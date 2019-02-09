package common

import "time"

type TaskResult struct {
	Timestamp time.Time
	Result    string
}

func NewTaskResult(timestamp time.Time, result string) TaskResult {
	taskResult := TaskResult{}
	taskResult.Result = result
	taskResult.Timestamp = timestamp
	return taskResult
}
