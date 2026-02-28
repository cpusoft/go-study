package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
// -------------------------- 测试工具函数 --------------------------

// 模拟结果保存函数（测试用）
func mockResultSaveFunc(ctx context.Context, group *TaskGroup) error {
	// 测试时可控制是否返回错误
	if group.GroupID == "test_save_fail" {
		return errors.New("模拟保存失败")
	}
	return nil
}

// 生成测试用的框架配置（缩短周期和超时，便于测试）
func getTestConfig() TaskFrameworkConfig {
	return TaskFrameworkConfig{
		CycleInterval:     1 * time.Second, // 测试用1秒周期
		ExecuteTimeout:    2 * time.Second, // 测试用2秒执行超时
		MaxFailCount:      2,               // 最大失败2次
		MaxWaitCycleCount: 3,               // 最多等3个周期
		LockCycleCount:    2,               // 锁定2个周期
	}
}

// 等待组状态变为指定值（带超时，避免死等）
func waitGroupState(t *testing.T, fw *TaskFramework, ctx context.Context, groupID string, targetState TaskState, timeout time.Duration) TaskState {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		state, err := fw.GetGroupState(ctx, groupID)
		if err == nil && state == targetState {
			return state
		}
		time.Sleep(100 * time.Millisecond)
	}
	// 超时后返回当前状态
	state, _ := fw.GetGroupState(ctx, groupID)
	return state
}

// -------------------------- 核心测试用例 --------------------------

// TestNewTaskFramework 测试框架初始化
func TestNewTaskFramework(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// 测试默认配置
	config := TaskFrameworkConfig{} // 空配置
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	assert.NotNil(t, fw)
	assert.Equal(t, 30*time.Minute, fw.config.CycleInterval)
	assert.Equal(t, 10*time.Minute, fw.config.ExecuteTimeout)
	assert.Equal(t, 3, fw.config.MaxFailCount)
	assert.Equal(t, 10, fw.config.MaxWaitCycleCount)
	assert.Equal(t, 10, fw.config.LockCycleCount)

	// 测试自定义配置
	customConfig := getTestConfig()
	fw2 := NewTaskFramework(ctx, customConfig, mockResultSaveFunc)
	defer fw2.Stop(ctx)

	assert.Equal(t, 1*time.Second, fw2.config.CycleInterval)
	assert.Equal(t, 2*time.Second, fw2.config.ExecuteTimeout)
	assert.Equal(t, 2, fw2.config.MaxFailCount)
	assert.Equal(t, 3, fw2.config.MaxWaitCycleCount)
	assert.Equal(t, 2, fw2.config.LockCycleCount)
}

// TestAddTasksFromExternal 测试外部注入任务（下周期执行）
func TestAddTasksFromExternal(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 1. 测试正常添加
	taskFuncs := map[string]func(ctx context.Context) error{
		"ext_task1": func(ctx context.Context) error { return nil },
		"ext_task2": func(ctx context.Context) error { return nil },
	}
	groupID, err := fw.AddTasksFromExternal(ctx, []string{"ext_task1", "ext_task2"}, taskFuncs)
	assert.NoError(t, err)
	assert.NotEmpty(t, groupID)

	// 验证组状态（下周期执行）
	state, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, state)

	// 验证活跃Key
	fw.mu.RLock()
	_, exist1 := fw.activeTaskKeys["ext_task1"]
	_, exist2 := fw.activeTaskKeys["ext_task2"]
	fw.mu.RUnlock()
	assert.True(t, exist1)
	assert.True(t, exist2)

	// 2. 测试排重（重复Key）
	groupID2, err := fw.AddTasksFromExternal(ctx, []string{"ext_task1", "ext_task3"}, taskFuncs)
	assert.NoError(t, err)
	assert.NotEmpty(t, groupID2)

	// 验证新组仅包含不重复的Key
	fw.mu.RLock()
	group2 := fw.groups[groupID2]
	fw.mu.RUnlock()
	assert.Len(t, group2.Tasks, 1)
	assert.Equal(t, "ext_task3", group2.Tasks[0].Key)

	// 3. 测试全重复Key
	groupID3, err := fw.AddTasksFromExternal(ctx, []string{"ext_task1", "ext_task2"}, taskFuncs)
	assert.Error(t, err)
	assert.Equal(t, "所有新任务Key均已存在，无需创建新组", err.Error())
	assert.Empty(t, groupID3)
}

// TestAddTasksFromSuccess 测试从成功任务添加新任务（立即执行）
func TestAddTasksFromSuccess(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 第一步：先创建外部注入的任务组（初始状态：待执行下周期执行）
	successTaskFuncs := map[string]func(ctx context.Context) error{
		"success_task": func(ctx context.Context) error {
			t.Log("执行初始任务：success_task")
			return nil
		},
	}
	successGroupID, err := fw.AddTasksFromExternal(ctx, []string{"success_task"}, successTaskFuncs)
	assert.NoError(t, err)
	assert.NotEmpty(t, successGroupID)

	// 验证初始状态：待执行下周期执行
	initialState, err := fw.GetGroupState(ctx, successGroupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, initialState)

	// 等待：至少2个周期（确保下周期转立即执行 + 执行完成）
	// 周期是1秒，等待3秒确保周期触发两次
	time.Sleep(3 * time.Second)

	// 验证初始组执行成功（核心断言）
	state := waitGroupState(t, fw, ctx, successGroupID, TaskStateSuccess, 5*time.Second)
	assert.Equal(t, TaskStateSuccess, state, "初始任务组未执行成功，当前状态：%s", state)

	// 第二步：从成功组添加新任务（立即执行）
	fw.mu.RLock()
	successGroup, exist := fw.groups[successGroupID]
	fw.mu.RUnlock()
	assert.True(t, exist, "成功任务组不存在")

	newTaskFuncs := map[string]func(ctx context.Context) error{
		"recur_task1": func(ctx context.Context) error {
			t.Log("执行递归任务：recur_task1")
			return nil
		},
		"recur_task2": func(ctx context.Context) error {
			t.Log("执行递归任务：recur_task2")
			return nil
		},
	}
	recurGroupID, err := fw.AddTasksFromSuccess(ctx, successGroup, []string{"recur_task1", "recur_task2"}, newTaskFuncs)
	assert.NoError(t, err, "添加递归任务失败：%v", err)
	assert.NotEmpty(t, recurGroupID, "递归任务组ID为空")

	// 验证原成功组改为下周期执行
	stateAfter, err := fw.GetGroupState(ctx, successGroupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, stateAfter, "原成功组状态未改为下周期执行，当前状态：%s", stateAfter)

	// 等待新组立即执行成功（最多等5秒）
	recurState := waitGroupState(t, fw, ctx, recurGroupID, TaskStateSuccess, 5*time.Second)
	assert.Equal(t, TaskStateSuccess, recurState, "递归任务组未立即执行成功，当前状态：%s", recurState)

	// 测试从非成功组添加（预期失败）
	// 重新创建一个未成功的组（仅添加，不执行）
	nonSuccessGroupID, err := fw.AddTasksFromExternal(ctx, []string{"non_success_task"}, map[string]func(ctx context.Context) error{
		"non_success_task": func(ctx context.Context) error { return nil },
	})
	assert.NoError(t, err)

	fw.mu.RLock()
	nonSuccessGroup, exist := fw.groups[nonSuccessGroupID]
	fw.mu.RUnlock()
	assert.True(t, exist)

	// 尝试从非成功组添加任务（预期失败）
	_, err = fw.AddTasksFromSuccess(ctx, nonSuccessGroup, []string{"test_fail"}, newTaskFuncs)
	assert.Error(t, err)
	assert.Equal(t, "仅能从执行成功的任务组添加新任务", err.Error())
}

// TestTaskExecute_Success 测试任务组执行成功
func TestTaskExecute_Success(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 添加立即执行的任务组（直接改状态，跳过外部注入）
	groupID := fmt.Sprintf("test_success_%d", time.Now().UnixNano())
	task := &Task{
		Key:         "success_task",
		FailCount:   0,
		FailReason:  "",
		ExecuteFunc: func(ctx context.Context) error { return nil },
	}
	group := &TaskGroup{
		GroupID:        groupID,
		Tasks:          []*Task{task},
		State:          TaskStatePendingImmediate,
		CreatedAt:      time.Now(),
		StartTime:      time.Time{},
		WaitCycleCount: 0,
	}

	fw.mu.Lock()
	fw.groups[groupID] = group
	fw.activeTaskKeys["success_task"] = struct{}{}
	fw.mu.Unlock()

	// 触发立即执行
	go fw.processPendingImmediateGroups(ctx)

	// 等待执行成功（最多等2秒）
	state := waitGroupState(t, fw, ctx, groupID, TaskStateSuccess, 2*time.Second)
	assert.Equal(t, TaskStateSuccess, state)

	// 验证任务无失败
	fw.mu.RLock()
	taskAfter := fw.groups[groupID].Tasks[0]
	fw.mu.RUnlock()
	assert.Equal(t, 0, taskAfter.FailCount)
	assert.Empty(t, taskAfter.FailReason)
}

// TestTaskExecute_Fail_Retry 测试任务失败重试（未超最大次数）
func TestTaskExecute_Fail_Retry(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 模拟失败任务（失败1次，未超最大次数2）
	failCount := 0
	taskFunc := func(ctx context.Context) error {
		failCount++
		return errors.New("模拟任务失败")
	}

	// 添加立即执行组
	groupID := fmt.Sprintf("test_fail_retry_%d", time.Now().UnixNano())
	group := &TaskGroup{
		GroupID:        groupID,
		Tasks:          []*Task{{Key: "fail_task", ExecuteFunc: taskFunc}},
		State:          TaskStatePendingImmediate,
		CreatedAt:      time.Now(),
		StartTime:      time.Time{},
		WaitCycleCount: 0,
	}
	fw.mu.Lock()
	fw.groups[groupID] = group
	fw.activeTaskKeys["fail_task"] = struct{}{}
	fw.mu.Unlock()

	// 触发执行
	go fw.processPendingImmediateGroups(ctx)

	// 等待状态改为下周期执行（最多等2秒）
	state := waitGroupState(t, fw, ctx, groupID, TaskStatePendingCycle, 2*time.Second)
	assert.Equal(t, TaskStatePendingCycle, state)

	// 验证失败次数
	fw.mu.RLock()
	task := fw.groups[groupID].Tasks[0]
	fw.mu.RUnlock()
	assert.Equal(t, 1, task.FailCount)
	assert.Equal(t, "模拟任务失败", task.FailReason)
}

// TestTaskExecute_Timeout 测试任务执行超时
func TestTaskExecute_Timeout(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	config.ExecuteTimeout = 1 * time.Second // 1秒超时
	config.MaxFailCount = 3                 // 最大失败次数3（避免单次超时超上限）
	config.MaxWaitCycleCount = 2            // 缩短最大等待周期，便于测试
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 模拟超时任务（仅执行1次，确保超时失败）
	executed := false
	taskFunc := func(ctx context.Context) error {
		if executed {
			return nil
		}
		executed = true
		// 强制返回超时错误，无需等待
		return fmt.Errorf("任务执行超时（%v）", config.ExecuteTimeout)
	}

	// 添加立即执行组
	groupID := fmt.Sprintf("test_timeout_%d", time.Now().UnixNano())
	group := &TaskGroup{
		GroupID:        groupID,
		Tasks:          []*Task{{Key: "timeout_task", ExecuteFunc: taskFunc}},
		State:          TaskStatePendingImmediate,
		CreatedAt:      time.Now(),
		StartTime:      time.Time{},
		WaitCycleCount: 0, // 初始等待周期0
	}
	// 标记为执行中，阻止重复执行
	fw.executingGroups.Store(groupID, struct{}{})
	fw.mu.Lock()
	fw.groups[groupID] = group
	fw.activeTaskKeys["timeout_task"] = struct{}{}
	fw.mu.Unlock()

	// 1. 手动执行任务组，验证超时后转为「下周期执行」
	fw.executeGroup(ctx, group)
	state, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, state, "超时后状态应为下周期执行，当前：%s", state)

	// 验证失败次数=1
	fw.mu.RLock()
	failCount := fw.groups[groupID].Tasks[0].FailCount
	fw.mu.RUnlock()
	assert.Equal(t, 1, failCount, "超时后失败次数应为1，当前：%d", failCount)

	// 2. 累加等待周期到最大上限（2）
	fw.mu.RLock()
	for _, g := range fw.groups {
		if g.GroupID == groupID {
			g.mu.Lock()
			g.WaitCycleCount = config.MaxWaitCycleCount // 设为2（最大等待周期）
			g.mu.Unlock()
		}
	}
	fw.mu.RUnlock()

	// 3. 手动触发周期处理，触发「最大等待周期」逻辑
	fw.processCycle(ctx)
	time.Sleep(100 * time.Millisecond)

	// 4. 验证最终状态为「执行失败」
	finalState, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStateFailed, finalState, "超过最大等待周期后应为执行失败，当前：%s", finalState)

	// 清理执行标记
	fw.executingGroups.Delete(groupID)
}

// TestResultSaveFail 测试结果保存失败（视为执行失败）
func TestResultSaveFail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	// 自定义保存失败的函数
	failSaveFunc := func(ctx context.Context, group *TaskGroup) error {
		return errors.New("模拟保存失败")
	}
	fw := NewTaskFramework(ctx, config, failSaveFunc)
	defer fw.Stop(ctx)

	// 添加结果保存失败的任务组
	groupID := "test_save_fail"
	group := &TaskGroup{
		GroupID:        groupID,
		Tasks:          []*Task{{Key: "save_fail_task", ExecuteFunc: func(ctx context.Context) error { return nil }}},
		State:          TaskStatePendingImmediate,
		CreatedAt:      time.Now(),
		StartTime:      time.Time{},
		WaitCycleCount: 0,
	}
	// 标记为执行中，阻止重复执行
	fw.executingGroups.Store(groupID, struct{}{})
	fw.mu.Lock()
	fw.groups[groupID] = group
	fw.activeTaskKeys["save_fail_task"] = struct{}{}
	fw.mu.Unlock()

	// 手动执行任务组
	fw.executeGroup(ctx, group)

	// 验证状态为执行失败
	state, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStateFailed, state, "保存失败后状态应为执行失败，当前：%s", state)

	// 验证失败原因包含保存失败信息
	fw.mu.RLock()
	task := fw.groups[groupID].Tasks[0]
	fw.mu.RUnlock()
	assert.Contains(t, task.FailReason, "任务执行成功但保存结果失败", "失败原因不含保存失败：%s", task.FailReason)
	assert.Contains(t, task.FailReason, "模拟保存失败", "失败原因不含具体错误：%s", task.FailReason)
	assert.Equal(t, 1, task.FailCount, "保存失败后失败次数应为1，当前：%d", task.FailCount)

	// 清理执行标记
	fw.executingGroups.Delete(groupID)
}

// TestConcurrentExecute 测试并发执行（组间异步，并发安全）
func TestConcurrentExecute(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 批量添加10个立即执行组
	var wg sync.WaitGroup
	groupIDs := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			groupID := fmt.Sprintf("concurrent_group_%d", idx)
			group := &TaskGroup{
				GroupID:        groupID,
				Tasks:          []*Task{{Key: fmt.Sprintf("concurrent_task_%d", idx), ExecuteFunc: func(ctx context.Context) error { time.Sleep(100 * time.Millisecond); return nil }}},
				State:          TaskStatePendingImmediate,
				CreatedAt:      time.Now(),
				StartTime:      time.Time{},
				WaitCycleCount: 0,
			}
			fw.mu.Lock()
			fw.groups[groupID] = group
			fw.activeTaskKeys[fmt.Sprintf("concurrent_task_%d", idx)] = struct{}{}
			fw.mu.Unlock()
			groupIDs = append(groupIDs, groupID)
		}(i)
	}
	wg.Wait()

	// 触发并发执行
	go fw.processPendingImmediateGroups(ctx)

	// 等待所有组执行成功（最多等3秒）
	for _, gid := range groupIDs {
		state := waitGroupState(t, fw, ctx, gid, TaskStateSuccess, 3*time.Second)
		assert.Equal(t, TaskStateSuccess, state, "组 %s 执行失败，当前状态：%s", gid, state)
	}

	// 验证无重复执行（executingGroups清理）
	fw.executingGroups.Range(func(key, value any) bool {
		assert.Fail(t, "executingGroups未清理，存在残留组：%s", key)
		return true
	})
}

// TestGetGroupState 测试获取任务组状态
func TestGetGroupState(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 测试不存在的组
	_, err := fw.GetGroupState(ctx, "non_exist_group")
	assert.Error(t, err)
	assert.Equal(t, "任务组不存在", err.Error())

	// 测试存在的组
	groupID := "test_get_state"
	group := &TaskGroup{
		GroupID: groupID,
		State:   TaskStatePendingCycle,
	}
	fw.mu.Lock()
	fw.groups[groupID] = group
	fw.mu.Unlock()

	state, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, state)
}

// TestStop 测试框架停止
func TestStop(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)

	// 停止框架
	fw.Stop(ctx)

	// 验证定时器停止（通过触发周期处理，无执行）
	time.Sleep(1500 * time.Millisecond)
	// 添加一个下周期执行的组，验证周期处理未触发
	groupID, _ := fw.AddTasksFromExternal(ctx, []string{"stop_task"}, map[string]func(ctx context.Context) error{"stop_task": func(ctx context.Context) error { return nil }})
	state, _ := fw.GetGroupState(ctx, groupID)
	assert.Equal(t, TaskStatePendingCycle, state) // 未执行，状态不变
}

*/
// -------------------------- 测试工具函数 --------------------------

// getTestConfig 获取测试用默认配置（缩短周期/超时，加快测试）
func getTestConfig() TaskFrameworkConfig {
	return TaskFrameworkConfig{
		CycleInterval:     1 * time.Second, // 测试用1秒周期
		ExecuteTimeout:    1 * time.Second, // 测试用1秒执行超时
		MaxFailCount:      2,               // 最大失败2次
		MaxWaitCycleCount: 2,               // 最大等待2个周期
		LockCycleCount:    2,               // 失败锁定2个周期
	}
}

// mockResultSaveFunc 模拟结果保存函数（默认成功）
func mockResultSaveFunc(ctx context.Context, group *TaskGroup) error {
	return nil
}

// mockFailSaveFunc 模拟结果保存失败函数
func mockFailSaveFunc(ctx context.Context, group *TaskGroup) error {
	return errors.New("模拟保存失败")
}

// waitGroupState 等待任务组状态达到目标（带超时）
func waitGroupState(t *testing.T, fw *TaskFramework, ctx context.Context, groupID string, targetState TaskState, timeout time.Duration) TaskState {
	start := time.Now()
	for time.Since(start) < timeout {
		state, err := fw.GetGroupState(ctx, groupID)
		if err == nil && state == targetState {
			return state
		}
		time.Sleep(100 * time.Millisecond)
	}
	// 超时返回当前状态
	state, _ := fw.GetGroupState(ctx, groupID)
	return state
}

// -------------------------- 核心功能测试 --------------------------

// TestTaskFramework_BasicSuccess 测试基础成功场景
func TestTaskFramework_BasicSuccess(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 模拟成功任务
	executed := false
	successFunc := func(ctx context.Context) error {
		executed = true
		t.Log("成功任务执行")
		return nil
	}

	// 添加立即执行组
	groupID, err := fw.AddTasksFromExternal(ctx, []string{"success_task"}, map[string]func(ctx context.Context) error{
		"success_task": successFunc,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, groupID)

	// 手动触发周期处理（转为立即执行）
	fw.processCycle(ctx)
	time.Sleep(500 * time.Millisecond)

	// 验证任务执行成功
	state := waitGroupState(t, fw, ctx, groupID, TaskStateSuccess, 2*time.Second)
	assert.Equal(t, TaskStateSuccess, state)
	assert.True(t, executed, "任务未执行")
}

// TestTaskFramework_FailMax 测试失败超最大次数（锁定）
func TestTaskFramework_FailMax(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 模拟失败任务
	failCount := 0
	failFunc := func(ctx context.Context) error {
		failCount++
		t.Logf("失败任务执行，次数：%d", failCount)
		return errors.New("模拟失败")
	}

	// 添加立即执行组
	groupID, err := fw.AddTasksFromExternal(ctx, []string{"fail_task"}, map[string]func(ctx context.Context) error{
		"fail_task": failFunc,
	})
	assert.NoError(t, err)

	// 手动触发2次执行（达到最大失败次数）
	fw.processCycle(ctx)
	time.Sleep(1 * time.Second)
	fw.processCycle(ctx)
	time.Sleep(1 * time.Second)

	// 验证状态为执行失败（锁定）
	state := waitGroupState(t, fw, ctx, groupID, TaskStateFailed, 2*time.Second)
	assert.Equal(t, TaskStateFailed, state)
	assert.GreaterOrEqual(t, failCount, config.MaxFailCount)

	// 验证锁定周期到期后转为下周期执行
	fw.processCycle(ctx) // 累加锁定计数1
	time.Sleep(100 * time.Millisecond)
	fw.processCycle(ctx) // 累加锁定计数2（解锁）
	time.Sleep(100 * time.Millisecond)

	unlockState, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, unlockState)
}

// TestTaskFramework_Timeout 测试任务执行超时
func TestTaskFramework_Timeout(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 模拟超时任务
	timeoutFunc := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			t.Log("任务收到超时信号，退出")
			return ctx.Err()
		case <-time.After(2 * time.Second): // 超过1秒超时
			return nil
		}
	}

	// 添加立即执行组
	groupID, err := fw.AddTasksFromExternal(ctx, []string{"timeout_task"},
		map[string]func(ctx context.Context) error{
			"timeout_task": timeoutFunc,
		})
	assert.NoError(t, err)

	// 手动触发周期处理
	fw.processCycle(ctx)
	time.Sleep(1 * time.Second)

	// 验证超时后转为下周期执行
	state, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, state)

	// 验证失败次数和原因
	fw.mu.RLock()
	task := fw.groups[groupID].Tasks[0]
	fw.mu.RUnlock()
	assert.Equal(t, 1, task.FailCount)
	assert.Contains(t, task.FailReason, "任务执行超时")

	// 模拟超最大等待周期
	fw.mu.RLock()
	fw.groups[groupID].mu.Lock()
	fw.groups[groupID].WaitCycleCount = config.MaxWaitCycleCount
	fw.groups[groupID].mu.Unlock()
	fw.mu.RUnlock()

	// 触发周期处理，验证转为失败
	fw.processCycle(ctx)
	finalState, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStateFailed, finalState)
}

// TestTaskFramework_ResultSaveFail 测试结果保存失败
func TestTaskFramework_ResultSaveFail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockFailSaveFunc) // 使用保存失败函数
	defer fw.Stop(ctx)

	// 模拟执行成功但保存失败的任务
	successFunc := func(ctx context.Context) error {
		t.Log("任务执行成功")
		return nil
	}

	// 添加立即执行组
	groupID, err := fw.AddTasksFromExternal(ctx, []string{"save_fail_task"}, map[string]func(ctx context.Context) error{
		"save_fail_task": successFunc,
	})
	assert.NoError(t, err)

	// 手动触发周期处理
	fw.processCycle(ctx)
	time.Sleep(500 * time.Millisecond)

	// 验证状态为执行失败
	state, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStateFailed, state)

	// 验证失败原因
	fw.mu.RLock()
	task := fw.groups[groupID].Tasks[0]
	fw.mu.RUnlock()
	assert.Contains(t, task.FailReason, "任务执行成功但保存结果失败")
	assert.Contains(t, task.FailReason, "模拟保存失败")
}

// TestTaskFramework_AddFromSuccess 测试从成功组添加新任务
func TestTaskFramework_AddFromSuccess(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 1. 先创建成功组
	successFunc := func(ctx context.Context) error { return nil }
	successGroupID, err := fw.AddTasksFromExternal(ctx, []string{"success_task"}, map[string]func(ctx context.Context) error{
		"success_task": successFunc,
	})
	assert.NoError(t, err)

	// 触发执行成功
	fw.processCycle(ctx)
	time.Sleep(500 * time.Millisecond)
	successGroup := fw.groups[successGroupID]

	// 2. 从成功组添加新任务
	newFunc := func(ctx context.Context) error { return nil }
	newGroupID, err := fw.AddTasksFromSuccess(ctx, successGroup, []string{"new_task"}, map[string]func(ctx context.Context) error{
		"new_task": newFunc,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, newGroupID)

	// 验证原成功组转为下周期执行
	oldState, err := fw.GetGroupState(ctx, successGroupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, oldState)

	// 验证新组立即执行成功
	newState := waitGroupState(t, fw, ctx, newGroupID, TaskStateSuccess, 2*time.Second)
	assert.Equal(t, TaskStateSuccess, newState)

	// 测试从非成功组添加（预期失败）
	failGroupID, _ := fw.AddTasksFromExternal(ctx, []string{"fail_task"}, map[string]func(ctx context.Context) error{
		"fail_task": func(ctx context.Context) error { return errors.New("fail") },
	})
	failGroup := fw.groups[failGroupID]
	_, err = fw.AddTasksFromSuccess(ctx, failGroup, []string{"test"}, map[string]func(ctx context.Context) error{})
	assert.Error(t, err)
	assert.Equal(t, "仅能从执行成功的任务组添加新任务", err.Error())
}

// -------------------------- 边界场景测试 --------------------------

// TestTaskFramework_DuplicateKey 测试重复Key排重
func TestTaskFramework_DuplicateKey(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 先添加一个Key
	_, err := fw.AddTasksFromExternal(ctx, []string{"dup_key"}, map[string]func(ctx context.Context) error{
		"dup_key": func(ctx context.Context) error { return nil },
	})
	assert.NoError(t, err)

	// 再次添加相同Key（预期失败）
	_, err = fw.AddTasksFromExternal(ctx, []string{"dup_key"}, map[string]func(ctx context.Context) error{
		"dup_key": func(ctx context.Context) error { return nil },
	})
	assert.Error(t, err)
	assert.Equal(t, "所有新任务Key均已存在，无需创建新组", err.Error())
}

// TestTaskFramework_PanicRecover 测试任务panic恢复
func TestTaskFramework_PanicRecover(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	config := getTestConfig()
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	// 模拟panic任务
	panicFunc := func(ctx context.Context) error {
		panic("任务panic")
	}

	// 添加立即执行组
	groupID, err := fw.AddTasksFromExternal(ctx, []string{"panic_task"}, map[string]func(ctx context.Context) error{
		"panic_task": panicFunc,
	})
	assert.NoError(t, err)

	// 手动触发执行
	fw.processCycle(ctx)
	time.Sleep(500 * time.Millisecond)

	// 验证任务标记为失败
	state, err := fw.GetGroupState(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, TaskStatePendingCycle, state)

	// 验证失败原因包含panic信息
	fw.mu.RLock()
	task := fw.groups[groupID].Tasks[0]
	fw.mu.RUnlock()
	assert.Contains(t, task.FailReason, "任务执行panic: 任务panic")
}

// -------------------------- 压力测试 --------------------------

// TestTaskFramework_Stress 压力测试：并发添加/执行100个任务组
func TestTaskFramework_Stress(t *testing.T) {
	t.Skip("压力测试需手动执行，默认跳过") // 可注释该行执行
	ctx := context.Background()
	config := getTestConfig()
	config.CycleInterval = 100 * time.Millisecond // 缩短周期加快压力测试
	fw := NewTaskFramework(ctx, config, mockResultSaveFunc)
	defer fw.Stop(ctx)

	const (
		groupCount = 100 // 并发添加100个组
		taskCount  = 5   // 每个组5个任务
	)

	var wg sync.WaitGroup
	wg.Add(groupCount)

	// 并发添加任务组
	start := time.Now()
	for i := 0; i < groupCount; i++ {
		go func(idx int) {
			defer wg.Done()
			// 构建任务Key和执行函数
			taskKeys := make([]string, 0, taskCount)
			taskFuncs := make(map[string]func(ctx context.Context) error)
			for j := 0; j < taskCount; j++ {
				key := fmt.Sprintf("stress_%d_%d", idx, j)
				taskKeys = append(taskKeys, key)
				taskFuncs[key] = func(ctx context.Context) error {
					time.Sleep(10 * time.Millisecond) // 模拟任务执行耗时
					return nil
				}
			}
			// 添加外部任务组
			_, err := fw.AddTasksFromExternal(ctx, taskKeys, taskFuncs)
			assert.NoError(t, err, "添加组%d失败", idx)
		}(i)
	}
	wg.Wait()
	t.Logf("并发添加%d个组完成，耗时：%v", groupCount, time.Since(start))

	// 手动触发多次周期处理，执行所有任务
	execStart := time.Now()
	for i := 0; i < 5; i++ {
		fw.processCycle(ctx)
		time.Sleep(100 * time.Millisecond)
	}
	t.Logf("执行所有任务完成，耗时：%v", time.Since(execStart))

	// 验证所有组最终状态为成功
	successCount := 0
	fw.mu.RLock()
	for _, group := range fw.groups {
		group.mu.RLock()
		if group.State == TaskStateSuccess {
			successCount++
		}
		group.mu.RUnlock()
	}
	fw.mu.RUnlock()

	t.Logf("成功执行组数量：%d / %d", successCount, groupCount)
	assert.Equal(t, groupCount, successCount, "部分组执行失败")

	// 验证活跃Key数量正确
	fw.mu.RLock()
	activeKeyCount := len(fw.activeTaskKeys)
	fw.mu.RUnlock()
	assert.Equal(t, groupCount*taskCount, activeKeyCount, "活跃Key数量不匹配")
}
