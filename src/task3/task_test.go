package task

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// 测试1：全部成功
func TestExecuteGroup_AllSuccess(t *testing.T) {
	psm := NewPeriodicStateMachine(ModeExternal)
	defer psm.Stop()

	var execCount int32 = 0

	// 无延迟成功任务
	task1 := &Task{
		Key: "task1",
		ExecuteFunc: func(ctx context.Context) error {
			atomic.AddInt32(&execCount, 1)
			return nil
		},
	}
	task2 := &Task{
		Key: "task2",
		ExecuteFunc: func(ctx context.Context) error {
			atomic.AddInt32(&execCount, 1)
			return nil
		},
	}

	err := psm.AddTasks("group1", []*Task{task1, task2})
	if err != nil {
		t.Fatalf("添加任务失败：%v", err)
	}

	psm.scheduleCycle()

	status := psm.GetGroupStatus("group1")
	if status != StatusSucceeded {
		t.Errorf("期望组状态为%s，实际为%s", StatusSucceeded, status)
	}

	if atomic.LoadInt32(&execCount) != 2 {
		t.Errorf("期望任务执行2次，实际执行%d次", execCount)
	}

	succ1, _ := psm.GetTaskCount("group1", "task1")
	succ2, _ := psm.GetTaskCount("group1", "task2")
	if succ1 != 1 || succ2 != 1 {
		t.Errorf("任务成功次数错误：task1=%d, task2=%d", succ1, succ2)
	}
}

// 测试2：单个任务失败
func TestExecuteGroup_OneFail(t *testing.T) {
	psm := NewPeriodicStateMachine(ModeExternal)
	defer psm.Stop()

	task1 := &Task{Key: "task1", ExecuteFunc: func(ctx context.Context) error { return nil }}
	task2 := &Task{Key: "task2", ExecuteFunc: func(ctx context.Context) error { return errors.New("测试失败") }}

	err := psm.AddTasks("group2", []*Task{task1, task2})
	if err != nil {
		t.Fatal(err)
	}

	psm.scheduleCycle()

	status := psm.GetGroupStatus("group2")
	if status != StatusFailed {
		t.Errorf("期望组状态为%s，实际为%s", StatusFailed, status)
	}

	_, fail1 := psm.GetTaskCount("group2", "task1")
	_, fail2 := psm.GetTaskCount("group2", "task2")
	if fail1 == 0 || fail2 == 0 {
		t.Error("失败任务未正确标记失败次数")
	}
}

// 测试3：超时未完成
func TestExecuteGroup_Timeout(t *testing.T) {
	psm := NewPeriodicStateMachine(ModeExternal)
	defer psm.Stop()
	// 测试时覆盖超时时间为1秒
	psm.SetExecTimeout(1 * time.Second)

	// 长时间任务（1.5秒），超过框架内1秒的超时时间
	var execFinish atomic.Bool
	var isFirstExec atomic.Bool
	isFirstExec.Store(true)

	task := &Task{
		Key: "task3",
		ExecuteFunc: func(ctx context.Context) error {
			if isFirstExec.Load() {
				isFirstExec.Store(false)
				// 第一次执行：等待1.5秒（超时）
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(2 * time.Second):
					execFinish.Store(true)
					return nil
				}
			} else {
				// 第二次执行：等待0.5秒（超时内完成）
				time.Sleep(1 * time.Second)
				execFinish.Store(true)
				return nil
			}
		},
	}

	err := psm.AddTasks("group3", []*Task{task})
	if err != nil {
		t.Fatal(err)
	}

	// 第一次调度（任务超时）
	psm.scheduleCycle()

	// 验证第一次调度后状态
	status := psm.GetGroupStatus("group3")
	if status != StatusRunning {
		t.Errorf("期望组状态为%s，实际为%s", StatusRunning, status)
	}
	cycles := psm.GetGroupWaitCycles("group3")
	if cycles != 1 {
		t.Errorf("期望等待周期数为1，实际为%d", cycles)
	}

	// 第二次调度（任务完成）
	psm.scheduleCycle()
	time.Sleep(600 * time.Millisecond)

	// 验证任务完成
	if !execFinish.Load() {
		t.Error("任务未按预期执行完成")
	}
	// 验证最终状态
	status = psm.GetGroupStatus("group3")
	if status != StatusSucceeded {
		t.Errorf("期望组状态为%s，实际为%s", StatusSucceeded, status)
	}
}

// 测试4：最大周期失败
func TestExecuteGroup_MaxCycle(t *testing.T) {
	psm := NewPeriodicStateMachine(ModeExternal)
	defer psm.Stop()
	psm.SetExecTimeout(1 * time.Second)

	// 永不完成的任务
	task := &Task{
		Key: "task4",
		ExecuteFunc: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}

	err := psm.AddTasks("group4", []*Task{task})
	if err != nil {
		t.Fatal(err)
	}

	// 执行10次调度
	for i := 0; i < 10; i++ {
		psm.scheduleCycle()
		time.Sleep(100 * time.Millisecond)
	}

	// 第11次调度，触发最大周期失败
	psm.scheduleCycle()

	// 验证状态
	status := psm.GetGroupStatus("group4")
	if status != StatusFailed {
		t.Errorf("期望组状态为%s，实际为%s", StatusFailed, status)
	}

	// 验证失败原因
	psm.mu.Lock()
	taskFailReason := psm.groupMap["group4"].Tasks["task4"].FailReason
	psm.mu.Unlock()
	if taskFailReason != "超过10个30分钟周期，强制标记失败" {
		t.Errorf("期望失败原因为%s，实际为%s", "超过10个30分钟周期，强制标记失败", taskFailReason)
	}
}

// 测试5：任务排重
func TestAddTasks_Duplicate(t *testing.T) {
	psm := NewPeriodicStateMachine(ModeExternal)
	defer psm.Stop()

	task1 := &Task{Key: "task5", ExecuteFunc: func(ctx context.Context) error { return nil }}
	err := psm.AddTasks("group5", []*Task{task1})
	if err != nil {
		t.Fatal(err)
	}

	task2 := &Task{Key: "task5", ExecuteFunc: func(ctx context.Context) error { return nil }}
	task3 := &Task{Key: "task6", ExecuteFunc: func(ctx context.Context) error { return nil }}
	err = psm.AddTasks("group6", []*Task{task2, task3})
	if err != nil {
		t.Fatal(err)
	}

	psm.mu.Lock()
	group := psm.groupMap["group6"]
	psm.mu.Unlock()

	if len(group.Tasks) != 1 {
		t.Errorf("期望组内1个任务，实际%d个", len(group.Tasks))
	}
	if _, exists := group.Tasks["task5"]; exists {
		t.Error("重复任务未被过滤")
	}
	if _, exists := group.Tasks["task6"]; !exists {
		t.Error("新任务未被添加")
	}
}
func TestAll(t *testing.T) {
	psm := NewPeriodicStateMachine(ModeExternal)

	psm.AddBanKey("ban_group_1")

	// 1分钟内完成的任务
	task1 := &Task{
		Key: "task_1",
		ExecuteFunc: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Minute):
				fmt.Println("执行任务1：成功")
				return nil
			}
		},
	}

	// 15分钟完成的任务
	task2 := &Task{
		Key: "task_2",
		ExecuteFunc: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(15 * time.Minute):
				fmt.Println("执行任务2：成功")
				return nil
			}
		},
	}

	// 执行失败的任务
	task3 := &Task{
		Key: "task_3",
		ExecuteFunc: func(ctx context.Context) error {
			fmt.Println("执行任务3：失败")
			return errors.New("任务3执行失败")
		},
	}

	err := psm.AddTasks("test_group", []*Task{task1, task2, task3})
	if err != nil {
		fmt.Printf("添加任务失败：%v\n", err)
	}

	go psm.StartScheduler()

	// 缩短模拟运行时间为10秒
	time.Sleep(10 * time.Second)

	psm.Stop()

	psm.mu.Lock()
	for groupKey, group := range psm.groupMap {
		fmt.Printf("组 %s 状态：%s，已等待周期数：%d\n", groupKey, group.Status, group.WaitCycles)
		for taskKey, task := range group.Tasks {
			fmt.Printf("  任务 %s：成功次数=%d，失败次数=%d，失败原因=%s，是否完成=%t\n",
				taskKey, task.SuccessCount, task.FailCount, task.FailReason, task.Done)
		}
	}
	psm.mu.Unlock()
}
