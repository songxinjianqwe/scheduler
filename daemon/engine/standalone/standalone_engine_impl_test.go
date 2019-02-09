package standalone

import (
	"github.com/satori/go.uuid"
	"github.com/songxinjianqwe/scheduler/common"
	"testing"
	"time"
)

var standAloneEngine = NewStandAloneEngine()

//************************************************************************************************
// DELAY TASK
//************************************************************************************************

/**
测试：延时任务正常执行的情况
*/
func TestDelayTaskExecuteNormally(t *testing.T) {
	id := "test_delay_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	delayTime := time.Second * 2
	task := common.NewTask(
		id,
		"delay",
		delayTime,
		"echo "+result,
	)
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	time.Sleep(delayTime + time.Second)
	assertStatusAndSameResult(t, task, common.Executed, result, 1)
}

/**
测试：在延时任务执行前将其Stop掉的情况
*/
func TestDelayTaskStopBeforeExecution(t *testing.T) {
	id := "test_delay_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	delayTime := time.Second * 2
	task := common.NewTask(
		id,
		"delay",
		delayTime,
		"echo "+result,
	)
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	time.Sleep(delayTime - time.Second)
	err = standAloneEngine.Stop(id)
	if err != nil {
		t.Fatalf("停止任务失败: %#v", task)
	}
	assertStatusAndSameResult(t, task, common.Stopped, "", 0)
}

/**
测试：在延时任务开始执行后将其Stop掉的情况
*/
func TestDelayTaskStopAfterStartExecution(t *testing.T) {
	id := "test_delay_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	delayTime := time.Second * 2
	executeTime := time.Second * 2
	task := common.NewTask(
		id,
		"delay",
		delayTime,
		"sleep 2s;echo "+result,
	)
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	time.Sleep(delayTime + executeTime - time.Second)
	assertStatus(t, task, common.Executing)
	err = standAloneEngine.Stop(id)
	if err == nil {
		t.Fatalf("此处应该停止任务失败，因为任务已经在执行了: %#v", task)
	}
	time.Sleep(time.Second * 2)
	assertStatusAndSameResult(t, task, common.Executed, result, 1)
}

/**
测试：在延时任务执行完毕后将其Stop掉的情况
*/
func TestDelayTaskStopAfterEndExecution(t *testing.T) {
	id := "test_delay_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	delayTime := time.Second * 2
	task := common.NewTask(
		id,
		"delay",
		delayTime,
		"echo "+result,
	)
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	time.Sleep(delayTime + time.Second)
	assertStatus(t, task, common.Executed)
	err = standAloneEngine.Stop(id)
	if err == nil {
		t.Fatalf("此处应该停止任务失败，因为任务已经在执行了: %#v", task)
	}
	assertStatusAndSameResult(t, task, common.Executed, result, 1)
}

//************************************************************************************************
// CRON TASK
//************************************************************************************************
/**
测试：定时任务正常执行两次
*/
func TestCronTaskExecuteNormally(t *testing.T) {
	id := "test_cron_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	cycleTime := time.Second * 2
	cycle := 2
	task := common.NewTask(
		id,
		"cron",
		cycleTime,
		"echo "+result,
	)
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	// 表示任务执行了cycle次
	time.Sleep(cycleTime*time.Duration(cycle) + time.Second)
	assertStatusAndSameResult(t, task, common.WaitForNextExecution, result, cycle)
}

/**
测试：在定时任务开始执行前Stop
*/
func TestCronTaskStopBeforeStartExecution(t *testing.T) {
	id := "test_cron_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	cycleTime := time.Second * 2
	task := common.NewTask(
		id,
		"cron",
		cycleTime,
		"echo "+result,
	)
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	time.Sleep(cycleTime - time.Second)
	err = standAloneEngine.Stop(id)
	if err != nil {
		t.Fatalf("停止任务失败: %#v", task)
	}
	assertStatusAndSameResult(t, task, common.Stopped, "", 0)
}

/**
测试：在定时任务开始执行后Stop
注意：scheduleAtFixedRate！计算真正的执行时间是 cycleTime + (cycle -1) * MAX(cycleTime, executionTime)
*/
func TestCronTaskAfterStartExecution(t *testing.T) {
	id := "test_cron_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	cycleTime := time.Second * 2
	cycle := 2
	executeTime := time.Second * 3
	task := common.NewTask(
		id,
		"cron",
		cycleTime,
		"sleep 3s;echo "+result,
	)
	t.Logf("提交任务: %s", time.Now())
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	// 表示任务执行了cycle次
	// 这里应该是5+1=6s，此时已经是第二次任务开始执行，即处于cmd.Run的时候
	time.Sleep(lastCycleStartTime(cycle, cycleTime, executeTime) + time.Second)
	assertStatus(t, task, common.Executing)
	err = standAloneEngine.Stop(id)
	if err == nil {
		t.Fatalf("此处应该停止任务失败，因为任务已经在执行了: %#v", task)
	}
	// 等待最后一次执行完毕
	time.Sleep(executeTime + time.Second)
	// 此时应该是两次均执行完毕，然后状态为Stopped
	assertStatusAndSameResult(t, task, common.Stopped, result, 2)
}

/**
测试：在定时任务开始执行后Stop
注意，如果executeTime比cycleTime，那么是非常难捕捉到WaitForNextExecution状态的
cycleTime和executionTime要保证有2s间隔，否则无法用1s来捕捉状态变更间隙（WaitForNextExecution和Executing）
*/
func TestCronTaskAfterEndExecution(t *testing.T) {
	id := "test_cron_" + uuid.Must(uuid.NewV4()).String()
	result := "1"
	cycleTime := time.Second * 4
	cycle := 2
	executeTime := time.Second * 2
	task := common.NewTask(
		id,
		"cron",
		cycleTime,
		"sleep 2s;echo "+result,
	)
	t.Logf("提交任务: %s", time.Now())
	err := standAloneEngine.Submit(*task)
	if err != nil {
		t.Fatalf("提交任务失败: %#v", task)
	}
	// 等第二次执行完，第三次开始执行前，这里应该睡眠是10(即4 + 2 + 4*(2-1))+1=11s
	time.Sleep(lastCycleEndTime(cycle, cycleTime, executeTime) + time.Second)
	assertStatus(t, task, common.WaitForNextExecution)
	err = standAloneEngine.Stop(id)
	if err != nil {
		t.Fatalf("停止任务失败: %s", err.Error())
	}
	time.Sleep(executeTime + time.Second)
	// 此时应该是两次均执行完毕，然后状态为Stopped
	assertStatusAndSameResult(t, task, common.Stopped, result, 2)
}

//************************************************************************************************
// COMMON
//************************************************************************************************

func assertStatus(t *testing.T, task *common.Task, status common.TaskStatus) {
	latestTask, err := standAloneEngine.Get(task.Id, false, 0)
	t.Logf("最新任务实例: %#v", latestTask)
	if err != nil {
		t.Fatalf("获取任务失败: %#v", task)
	}
	if latestTask.Status != status {
		t.Fatalf("任务状态错误: 预期: %s，实际：%s", status, latestTask.Status)
	}
}

func assertStatusAndSameResult(t *testing.T, task *common.Task, status common.TaskStatus, result string, resultSize int) {
	latestTask, err := standAloneEngine.Get(task.Id, false, 0)
	t.Logf("最新任务实例: %#v", latestTask)
	if err != nil {
		t.Fatalf("获取任务失败: %#v", task)
	}
	if latestTask.Status != status {
		t.Fatalf("任务状态错误: %#v", latestTask)
	}
	if len(latestTask.Results) != resultSize {
		t.Fatalf("任务结果数量错误: 预期: %d，实际: %d", resultSize, len(latestTask.Results))
	}
	for _, taskResult := range latestTask.Results {
		if taskResult.Result != (result + "\n") {
			t.Fatalf("任务结果错误: %#v", latestTask)
		}
	}
}

/**
返回的是最后一次任务的开始执行时间，状态变成executing一般是lastCycleStartTime+1s（如果执行时间多于1s的话）
示例1：
	cycle=3,cycleTime=2s,executionTime=3s
	那么第3次任务开始的时间为2s+3s+(3-2)*3s=8s
示例2：
	cycle=3,cycleTime=4s,executionTime=2s
	那么第3次任务开始的时间为4s+2s+(3-2)*4=12s
*/
func lastCycleStartTime(cycle int, cycleTime time.Duration, executionTime time.Duration) time.Duration {
	if cycle == 1 {
		return time.Duration(0)
	}
	return cycleTime + executionTime + time.Duration(cycle-2)*Max(executionTime, cycleTime)
}

/**
返回的是最后一次任务的执行完成时间，状态变成executing一般是lastCycleStartTime+1s（如果执行时间多于1s的话）
示例1：
	cycle=3,cycleTime=2s,executionTime=3s
	那么第3次任务结束的时间为2s+3s+(4-2)*3s=11s
示例2：
	cycle=3,cycleTime=4s,executionTime=2s
	那么第3次任务开始的时间为4s+2s+(4-2)*4=16s
*/
func lastCycleEndTime(cycle int, cycleTime time.Duration, executionTime time.Duration) time.Duration {
	return lastCycleStartTime(cycle+1, cycleTime, executionTime)
}

func Max(d1 time.Duration, d2 time.Duration) time.Duration {
	if d1 > d2 {
		return d1
	}
	return d2
}
