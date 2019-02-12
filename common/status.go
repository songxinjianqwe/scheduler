package common

// 枚举定义开始
type TaskStatus int

// 对于delay类型的任务，只有ToBeExecuted, Executing, Executed, Stopped状态
// 对于cron类型的任务，有ToBeExecuted, Executing, WaitForNextExecution, Stopped状态
// 对于delay类型的任务的状态流转：
// 只能在ToBeExecuted时stop
// 1) ToBeExecuted -> Executing -> Executed
// 2）ToBeExecuted -> Stopped

// 对于crone类型的任务的状态流转:
// 只能在ToBeExecuted或者WaitForNextExecution时stop
// 1) ToBeExecuted -> Executing -> WaitForNextExecution
// 						 ↑		←		↓
// 2) ToBeExecuted -> Stopped
// 3) ToBeExecuted -> Executing -> WaitForNextExecution -> Stopped
const (
	ToBeExecuted         TaskStatus = iota
	Executed             TaskStatus = iota
	Executing            TaskStatus = iota
	WaitForNextExecution TaskStatus = iota
	Stopped              TaskStatus = iota
)

var statusText = map[TaskStatus]string{
	ToBeExecuted:         "ToBeExecuted",
	Executed:             "Executed",
	Executing:            "Executing",
	WaitForNextExecution: "WaitForNextExecution",
	Stopped:              "Stopped",
}

func (status TaskStatus) String() string {
	text, ok := statusText[status]
	if ok {
		return text
	}
	return "UNKNOWN"
}

// 枚举定义结束
