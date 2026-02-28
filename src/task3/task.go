package task

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 任务状态枚举
type TaskGroupStatus string

const (
	StatusPending TaskGroupStatus = "待执行"
	StatusRunning TaskGroupStatus = "执行中"
	StatusDone    TaskGroupStatus = "已完成"

// StatusSucceeded TaskGroupStatus = "执行成功"
// StatusFailed    TaskGroupStatus = "执行失败"
)

// 添加任务模式
type AddTaskMode int

const (
	ModeRecursive AddTaskMode = iota
	ModeExternal
)

// 单个任务
type Task struct {
	Key    string
	Result string // success/fail

	FailTime   time.Time
	FailCount  int
	FailResult string

	SuccessTime   time.Time
	SuccessCount  int
	SuccessResult string

	ExecuteFunc func(ctx context.Context) error // 支持ctx终止的执行函数
	Done        bool
}

// 任务组
type TaskGroup struct {
	GroupKey          string
	Status            TaskGroupStatus
	Tasks             map[string]*Task
	TaskWaitCycles    int
	TaskMaxWaitCycles int // 3个周期（3*30分钟/周期）后强制标记失败
}

// 任务执行结果（新增：收集单个任务执行结果）
type taskExecResult struct {
	task      *Task
	isTimeout bool
	isFail    bool
}

// 状态机核心
type PeriodicStateMachine struct {
	taskGroupMap map[string]*TaskGroup
	taskGroupMut sync.RWMutex

	banMap map[string]struct{}
	banMut sync.RWMutex

	mode            AddTaskMode
	cycleticker     *time.Ticker
	maxWaitCycles   int
	taskExecTimeout time.Duration // 单个任务执行超时时间
}

// 初始化（新增execTimeout参数，默认10分钟）
func NewPeriodicStateMachine(mode AddTaskMode) *PeriodicStateMachine {
	return &PeriodicStateMachine{
		taskGroupMap:    make(map[string]*TaskGroup),
		banMap:          make(map[string]struct{}),
		mode:            mode,
		maxWaitCycles:   3,
		taskExecTimeout: 3 * 30 * time.Minute, // 单个任务最长90分钟
		cycleticker:     time.NewTicker(30 * time.Minute),
	}
}

// ///////////////////
// ban
// //////////////////
func (psm *PeriodicStateMachine) AddBanKey(groupKey string) {
	psm.banMut.Lock()
	defer psm.banMut.Unlock()
	psm.banMap[groupKey] = struct{}{}
}

func (psm *PeriodicStateMachine) RemoveBanKey(groupKey string) {
	psm.banMut.Lock()
	defer psm.banMut.Unlock()
	delete(psm.banMap, groupKey)
}

func (psm *PeriodicStateMachine) isBanned(groupKey string) bool {
	psm.banMut.RLock()
	defer psm.banMut.RUnlock()
	_, exists := psm.banMap[groupKey]
	return exists
}

// ///////////////////
// 任务添加与执行
// //////////////////
func (psm *PeriodicStateMachine) AddTasks(groupKey string, tasks []*Task) error {
	if psm.isBanned(groupKey) {
		return fmt.Errorf("组 %s 被禁止执行", groupKey)
	}

	psm.taskGroupMut.Lock()
	defer psm.taskGroupMut.Unlock()

	// 排重逻辑
	uniqueTasks := make(map[string]*Task)
	for _, t := range tasks {
		t.Done = false
		if !psm.isTaskDuplicatedWithinMut(t.Key) {
			uniqueTasks[t.Key] = t
		} else {
			fmt.Printf("任务 %s 重复，已过滤\n", t.Key)
		}
	}

	if len(uniqueTasks) == 0 {
		return errors.New("所有任务均为重复，无可用任务添加")
	}

	// 创建/更新任务组
	if _, exists := psm.taskGroupMap[groupKey]; !exists {
		var status = StatusPending
		if psm.mode == ModeRecursive {
			status = StatusRunning
		}
		psm.taskGroupMap[groupKey] = &TaskGroup{
			GroupKey:          groupKey,
			Status:            status,
			Tasks:             uniqueTasks,
			TaskMaxWaitCycles: psm.maxWaitCycles,
		}
		// 递归模式立即执行（异步）
		if psm.mode == ModeRecursive {
			// 异步执行，调用方不等待
			go psm.executeGroup(psm.taskGroupMap[groupKey])
		}
	}

	return nil
}

func (psm *PeriodicStateMachine) isTaskDuplicatedWithinMut(taskKey string) bool {
	for _, g := range psm.taskGroupMap {
		if g.Status == StatusPending || g.Status == StatusRunning {
			if _, exists := g.Tasks[taskKey]; exists {
				return true
			}
		}
	}
	return false
}

func (psm *PeriodicStateMachine) markSuccessToPendingWithinMut() {

	for _, g := range psm.taskGroupMap {
		if g.Status == StatusSucceeded {
			g.Status = StatusPending
		}
	}
}

// -------------------------- 调整：executeSingleTask 异步执行+返回结果 --------------------------
// executeSingleTask 异步执行单个任务，通过channel返回执行结果
// 入参：task-待执行任务，execTimeout-单次执行超时时间，resultChan-结果返回通道，wg-等待组
func (psm *PeriodicStateMachine) executeSingleTask(task *Task, taskExecTimeout time.Duration, resultChan chan<- taskExecResult, wg *sync.WaitGroup) {
	defer wg.Done() // 任务执行完成，等待组计数-1

	result := taskExecResult{task: task}

	// 创建超时ctx
	ctx, cancel := context.WithTimeout(context.Background(), taskExecTimeout)
	defer cancel()

	// 执行任务（同步，内部逻辑不变）
	err := task.ExecuteFunc(ctx)

	// 判断执行结果
	if ctx.Err() != nil {
		// 任务超时未完成
		result.isTimeout = true
		task.Done = false
		resultChan <- result
		return
	}

	// 任务在超时内完成
	task.Done = true
	if err != nil {
		// 执行失败
		result.isFail = true
		task.FailResult = err.Error()
		task.FailCount++
	} else {
		// 执行成功
		task.SuccessCount++
		task.SuccessTime = time.Now()
		task.Result = "成功"
	}

	resultChan <- result
}

// -------------------------- 调整：executeGroup 完全异步执行 --------------------------
func (psm *PeriodicStateMachine) executeGroup(group *TaskGroup) {
	// 异步执行，外层包裹goroutine，调用方立即返回
	go func() {
		// 加锁保护组状态修改（异步执行仍需保证线程安全）
		psm.taskGroupMut.Lock()
		defer psm.taskGroupMut.Unlock()

		// 前置检查：被禁止/已完成的组不执行
		if psm.isBanned(group.GroupKey) || group.Status == StatusPending {
			return
		}

		// 状态切换为执行中（初始化周期）
		if group.Status == StatusPending {
			group.Status = StatusRunning
			group.TaskWaitCycles = 0
		}

		// 检查最大等待周期：超过则标记失败
		if group.TaskWaitCycles >= group.TaskMaxWaitCycles {
			psm.markGroupFail(group, "超过10个30分钟周期，强制标记失败")
			return
		}

		// 筛选未完成任务
		unfinishedTasks := psm.filterUnfinishedTasks(group)
		if len(unfinishedTasks) == 0 {
			// 所有任务已完成：更新组状态
			psm.updateGroupStatus(group)
			return
		}

		// 批量异步执行未完成任务，等待所有结果
		var (
			hasFail         bool                                              // 是否有任务执行失败
			hasTimeout      bool                                              // 是否有任务超时
			taskExecTimeout = psm.taskExecTimeout                             // 单次任务超时时间
			wg              sync.WaitGroup                                    // 等待所有任务完成
			resultChan      = make(chan taskExecResult, len(unfinishedTasks)) // 结果通道（带缓冲，避免阻塞）
		)

		// 启动所有任务异步执行
		for _, task := range unfinishedTasks {
			wg.Add(1)
			go psm.executeSingleTask(task, taskExecTimeout, resultChan, &wg)
		}

		// 等待所有任务完成后关闭通道
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// 收集所有任务执行结果
		for result := range resultChan {
			if result.isTimeout {
				hasTimeout = true
			}
			if result.isFail {
				hasFail = true
			}
		}

		// 根据执行结果更新组状态
		if hasTimeout {
			// 有任务超时：保持执行中，累加周期
			group.Status = StatusRunning
			group.TaskWaitCycles++
		} else {
			// 所有任务在超时内完成：根据是否失败更新状态
			if hasFail {
				psm.markGroupPartialFail(group) // 部分任务失败，标记整组失败
			} else {
				group.Status = StatusSucceeded // 全部成功
			}
			group.TaskWaitCycles = 0 // 重置周期
		}
	}()
}

// -------------------------- 辅助方法（无修改） --------------------------
// filterUnfinishedTasks 筛选组内未完成的任务
func (psm *PeriodicStateMachine) filterUnfinishedTasks(group *TaskGroup) []*Task {
	var unfinished []*Task
	for _, task := range group.Tasks {
		if !task.Done {
			unfinished = append(unfinished, task)
		}
	}
	return unfinished
}

// markGroupPartialFail 标记组内部分任务失败（整组失败）
func (psm *PeriodicStateMachine) markGroupPartialFail(group *TaskGroup) {
	group.Status = StatusFailed
	// 为无失败原因的任务补充默认原因
	for _, task := range group.Tasks {
		if task.FailResult == "" {
			task.FailResult = "组内存在失败任务，标记为失败"
			task.FailCount++
		}
	}
}

func (psm *PeriodicStateMachine) updateGroupStatus(group *TaskGroup) {
	allSuccess := true
	for _, t := range group.Tasks {
		if t.FailCount > 0 {
			allSuccess = false
			break
		}
	}

	if allSuccess {
		group.Status = StatusSucceeded
	} else {
		psm.markGroupPartialFail(group)
	}
}

func (psm *PeriodicStateMachine) markGroupFail(group *TaskGroup, reason string) {
	group.Status = StatusFailed
	for _, t := range group.Tasks {
		t.FailResult = reason
		t.FailCount++
		t.Done = true
	}
}

// -------------------------- 调度方法（调整为异步执行） --------------------------
func (psm *PeriodicStateMachine) scheduleCycle() {
	psm.taskGroupMut.RLock()
	// 复制当前任务组，避免锁持有时间过长（异步执行特性）
	groups := make([]*TaskGroup, 0, len(psm.taskGroupMap))
	for _, g := range psm.taskGroupMap {
		if g.Status == StatusPending || g.Status == StatusRunning {
			if !psm.isBanned(g.GroupKey) {
				groups = append(groups, g)
			}
		}
	}
	psm.taskGroupMut.RUnlock()

	// 异步执行所有符合条件的任务组
	for _, g := range groups {
		psm.executeGroup(g)
	}
}

func (psm *PeriodicStateMachine) StartScheduler() {
	// 立即异步执行一次调度
	go psm.scheduleCycle()

	// 周期性异步执行
	go func() {
		for range psm.cycleticker.C {
			psm.scheduleCycle()
		}
	}()
}

func (psm *PeriodicStateMachine) Stop() {
	psm.cycleticker.Stop()
}

// -------------------------- 测试辅助方法（调整锁类型） --------------------------
// SetExecTimeout 用于测试时覆盖单个任务超时时间
func (psm *PeriodicStateMachine) SetExecTimeout(timeout time.Duration) {
	psm.taskGroupMut.Lock()
	defer psm.taskGroupMut.Unlock()
	psm.taskExecTimeout = timeout
}

func (psm *PeriodicStateMachine) GetGroupStatus(groupKey string) TaskGroupStatus {
	psm.taskGroupMut.RLock()
	defer psm.taskGroupMut.RUnlock()
	if g, ok := psm.taskGroupMap[groupKey]; ok {
		return g.Status
	}
	return ""
}

func (psm *PeriodicStateMachine) GetTaskCount(groupKey, taskKey string) (int, int) {
	psm.taskGroupMut.RLock()
	defer psm.taskGroupMut.RUnlock()
	if g, ok := psm.taskGroupMap[groupKey]; ok {
		if t, ok := g.Tasks[taskKey]; ok {
			return t.SuccessCount, t.FailCount
		}
	}
	return 0, 0
}

func (psm *PeriodicStateMachine) GetGroupWaitCycles(groupKey string) int {
	psm.taskGroupMut.RLock()
	defer psm.taskGroupMut.RUnlock()
	if g, ok := psm.taskGroupMap[groupKey]; ok {
		return g.TaskWaitCycles
	}
	return 0
}
