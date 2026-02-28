package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// -------------------------- 核心枚举与结构体定义 --------------------------

// TaskState 任务组状态枚举
type TaskState string

const (
	TaskStatePendingImmediate TaskState = "待执行立即执行"  // 立即执行
	TaskStatePendingCycle     TaskState = "待执行下周期执行" // 下周期执行
	TaskStateRunning          TaskState = "执行中"      // 执行中
	TaskStateSuccess          TaskState = "执行成功"     // 执行成功
	TaskStateFailed           TaskState = "执行失败"     // 执行失败
)

// Task 单个任务结构体
type Task struct {
	Key         string                          // 任务唯一标识（用于排重）
	FailCount   int                             // 失败次数
	FailReason  string                          // 失败原因
	ExecuteFunc func(ctx context.Context) error // 任务执行函数（抽象执行逻辑）
}

// TaskGroup 任务组结构体
type TaskGroup struct {
	GroupID        string       // 组唯一标识
	Tasks          []*Task      // 组内任务列表
	State          TaskState    // 组状态
	CreatedAt      time.Time    // 创建时间
	StartTime      time.Time    // 开始执行时间
	WaitCycleCount int          // 等待周期数（用于超过10周期的判断）
	TimeoutTimer   *time.Timer  // 10分钟执行超时定时器
	mu             sync.RWMutex // 状态保护锁
}

// TaskFrameworkConfig 框架配置
type TaskFrameworkConfig struct {
	CycleInterval     time.Duration // 周期间隔（默认30分钟）
	ExecuteTimeout    time.Duration // 单组执行超时（默认10分钟）
	MaxFailCount      int           // 最大失败重试次数
	MaxWaitCycleCount int           // 最大等待周期数（默认10）
	LockCycleCount    int           // 失败锁定周期数（默认10）
}

// TaskFramework 周期状态机任务框架核心结构体
type TaskFramework struct {
	config          TaskFrameworkConfig                               // 配置
	groups          map[string]*TaskGroup                             // 所有任务组（key: GroupID）
	activeTaskKeys  map[string]struct{}                               // 活跃任务Key（排重用：待处理、执行中、本周期成功）
	cycleTimer      *time.Ticker                                      // 周期定时器
	resultSaveFunc  func(ctx context.Context, group *TaskGroup) error // 结果保存接口函数
	mu              sync.RWMutex                                      // 全局锁
	executingGroups sync.Map                                          // 标记正在执行的组（避免重复执行）
}

// -------------------------- 框架初始化 --------------------------

// NewTaskFramework 创建新的任务框架实例
func NewTaskFramework(ctx context.Context, config TaskFrameworkConfig,
	resultSaveFunc func(ctx context.Context, group *TaskGroup) error) *TaskFramework {
	// 设置默认配置
	if config.CycleInterval == 0 {
		config.CycleInterval = 30 * time.Minute
	}
	if config.ExecuteTimeout == 0 {
		config.ExecuteTimeout = 10 * time.Minute
	}
	if config.MaxFailCount == 0 {
		config.MaxFailCount = 3 // 默认最大失败3次
	}
	if config.MaxWaitCycleCount == 0 {
		config.MaxWaitCycleCount = 10 // 默认最多等10个周期
	}
	if config.LockCycleCount == 0 {
		config.LockCycleCount = 10 // 默认锁定10个周期
	}

	fw := &TaskFramework{
		config:         config,
		groups:         make(map[string]*TaskGroup),
		activeTaskKeys: make(map[string]struct{}),
		resultSaveFunc: resultSaveFunc,
	}

	// 启动周期定时器
	fw.cycleTimer = time.NewTicker(config.CycleInterval)
	go fw.cycleHandler(ctx)

	return fw
}

// -------------------------- 周期处理核心逻辑 --------------------------
// processCycle 周期处理核心逻辑（补充最大等待周期处理）
func (fw *TaskFramework) processCycle(ctx context.Context) {
	fw.mu.RLock()
	// 1. 处理锁定周期到期的失败组（改为待执行下周期）
	for _, group := range fw.groups {
		group.mu.Lock()
		if group.State == TaskStateFailed {
			group.WaitCycleCount++
			if group.WaitCycleCount >= fw.config.LockCycleCount {
				group.State = TaskStatePendingCycle
				group.WaitCycleCount = 0
			}
		}
		group.mu.Unlock()
	}

	// 2. 处理最大等待周期（核心新增：超期则标记为失败）
	for _, group := range fw.groups {
		group.mu.Lock()
		// 仅处理「下周期执行/执行中」的组，且等待周期超上限
		if (group.State == TaskStatePendingCycle || group.State == TaskStateRunning) &&
			group.WaitCycleCount >= fw.config.MaxWaitCycleCount {
			group.State = TaskStateFailed // 超最大等待周期 → 执行失败
			group.WaitCycleCount = 0
		}
		group.mu.Unlock()
	}

	// 3. 处理非失败组的「下周期转立即执行」
	for _, group := range fw.groups {
		group.mu.Lock()
		if group.State == TaskStatePendingCycle && group.State != TaskStateFailed && group.State != TaskStateRunning {
			group.State = TaskStatePendingImmediate
		}
		group.mu.Unlock()
	}
	fw.mu.RUnlock()

	// 4. 处理待执行立即执行的组
	fw.processPendingImmediateGroups(ctx)

	// 5. 清理本周期成功的任务Key
	fw.mu.Lock()
	for key := range fw.activeTaskKeys {
		exist := false
		for _, group := range fw.groups {
			group.mu.RLock()
			if group.State == TaskStatePendingImmediate ||
				group.State == TaskStatePendingCycle ||
				group.State == TaskStateRunning {
				for _, task := range group.Tasks {
					if task.Key == key {
						exist = true
						break
					}
				}
			}
			group.mu.RUnlock()
			if exist {
				break
			}
		}
		if !exist {
			delete(fw.activeTaskKeys, key)
		}
	}
	fw.mu.Unlock()
}

// cycleHandler 周期处理入口（30分钟触发一次）
func (fw *TaskFramework) cycleHandler(ctx context.Context) {
	for range fw.cycleTimer.C {
		fw.processCycle(ctx) // 调用抽离的核心逻辑
	}
}

// processPendingImmediateGroups 处理所有待执行立即执行的组（独立函数，供周期和立即执行调用）
func (fw *TaskFramework) processPendingImmediateGroups(ctx context.Context) {
	fw.mu.RLock()
	// 筛选所有待执行立即执行的组
	pendingGroups := make([]*TaskGroup, 0)
	for _, group := range fw.groups {
		group.mu.RLock()
		if group.State == TaskStatePendingImmediate {
			// 先改为执行中（避免重复处理）
			group.State = TaskStateRunning
			group.StartTime = time.Now()
			pendingGroups = append(pendingGroups, group)
		}
		group.mu.RUnlock()
	}
	fw.mu.RUnlock()

	// 异步执行每个待处理组（避免阻塞）
	var wg sync.WaitGroup
	for _, group := range pendingGroups {
		// 检查是否正在执行，避免重复
		if _, loaded := fw.executingGroups.LoadOrStore(group.GroupID, struct{}{}); loaded {
			continue
		}
		wg.Add(1)
		go func(g *TaskGroup) {
			defer func() {
				wg.Done()
				fw.executingGroups.Delete(g.GroupID) // 执行完成后移除标记
			}()
			fw.executeGroup(ctx, g)
		}(group)
	}
	wg.Wait()
}

// executeTask 执行单个任务（带超时控制）
func (fw *TaskFramework) executeTask(ctx context.Context, task *Task) error {
	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, fw.config.ExecuteTimeout)
	defer cancel() // 确保上下文被取消，释放资源

	// 用通道接收任务执行结果
	resultChan := make(chan error, 1)
	go func() {
		defer func() {
			// 捕获panic，避免任务崩溃导致通道阻塞
			if r := recover(); r != nil {
				resultChan <- fmt.Errorf("任务执行panic: %v", r)
			}
		}()
		resultChan <- task.ExecuteFunc(timeoutCtx)
	}()

	// 等待任务执行完成或超时
	select {
	case err := <-resultChan:
		return err
	case <-timeoutCtx.Done():
		// 超时：返回超时错误，标记任务失败
		return fmt.Errorf("任务执行超时（%v）", fw.config.ExecuteTimeout)
	}
}

// executeGroup 执行任务组（核心逻辑：补全结果保存+精准失败判定）
func (fw *TaskFramework) executeGroup(ctx context.Context, group *TaskGroup) {
	group.mu.Lock()
	group.State = TaskStateRunning
	group.StartTime = time.Now()
	group.mu.Unlock()

	allSuccess := true
	for _, task := range group.Tasks {
		err := fw.executeTask(ctx, task)
		if err != nil {
			allSuccess = false
			task.FailCount++
			task.FailReason = err.Error()
		}
	}

	// 处理执行结果 + 结果保存逻辑（核心修复）
	group.mu.Lock()
	defer group.mu.Unlock()

	// 第一步：执行结果保存
	saveErr := fw.resultSaveFunc(ctx, group)

	// 第二步：判定最终状态（优先级：保存失败 > 任务执行失败 > 执行成功）
	if saveErr != nil {
		// 保存失败：无论任务是否成功，均标记为执行失败
		group.State = TaskStateFailed
		group.WaitCycleCount = 0
		// 补充保存失败的原因到第一个任务
		if len(group.Tasks) > 0 {
			group.Tasks[0].FailCount++
			group.Tasks[0].FailReason = fmt.Sprintf("任务执行成功但保存结果失败：%v", saveErr)
		}
	} else if !allSuccess {
		// 任务执行失败：判断是否超最大失败次数
		if group.Tasks[0].FailCount >= fw.config.MaxFailCount {
			group.State = TaskStateFailed
			group.WaitCycleCount = 0
		} else {
			// 未超上限：下周期重试
			group.State = TaskStatePendingCycle
			group.WaitCycleCount = 0
		}
	} else {
		// 任务执行成功且保存成功
		group.State = TaskStateSuccess
		group.WaitCycleCount = 0
	}

	// 移除执行中标记
	fw.executingGroups.Delete(group.GroupID)
}

// -------------------------- 任务添加接口（支持两种方式） --------------------------

// AddTasksFromSuccess 从成功任务结果添加新任务（方式4.2：立即执行，递归处理）
// successGroup: 成功的任务组（需改为下周期执行）
// newTaskKeys: 新任务Key列表
// newTaskExecuteFunc: 新任务执行函数（按Key映射）
// 返回：新创建的任务组ID（失败返回空）
func (fw *TaskFramework) AddTasksFromSuccess(ctx context.Context, successGroup *TaskGroup,
	newTaskKeys []string, newTaskExecuteFunc map[string]func(ctx context.Context) error) (string, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// 1. 检查成功组状态（必须是执行成功）
	successGroup.mu.RLock()
	if successGroup.State != TaskStateSuccess {
		successGroup.mu.RUnlock()
		return "", errors.New("仅能从执行成功的任务组添加新任务")
	}
	successGroup.mu.RUnlock()

	// 2. 排重：过滤已存在的活跃Key
	uniqueKeys := make([]string, 0)
	for _, key := range newTaskKeys {
		if _, exist := fw.activeTaskKeys[key]; !exist {
			uniqueKeys = append(uniqueKeys, key)
		}
	}
	if len(uniqueKeys) == 0 {
		return "", errors.New("所有新任务Key均已存在，无需创建新组")
	}

	// 3. 成功组改为待执行下周期
	successGroup.mu.Lock()
	successGroup.State = TaskStatePendingCycle
	successGroup.mu.Unlock()

	// 4. 创建新任务组（立即执行）
	newGroupID := fmt.Sprintf("group_%d", time.Now().UnixNano())
	newTasks := make([]*Task, 0, len(uniqueKeys))
	for _, key := range uniqueKeys {
		newTasks = append(newTasks, &Task{
			Key:         key,
			FailCount:   0,
			FailReason:  "",
			ExecuteFunc: newTaskExecuteFunc[key],
		})
		// 加入活跃Key索引
		fw.activeTaskKeys[key] = struct{}{}
	}

	newGroup := &TaskGroup{
		GroupID:        newGroupID,
		Tasks:          newTasks,
		State:          TaskStatePendingImmediate,
		CreatedAt:      time.Now(),
		StartTime:      time.Time{},
		WaitCycleCount: 0,
	}
	fw.groups[newGroupID] = newGroup

	// 5. 关键修复：立即触发待执行立即执行组的处理（无需等周期）
	go fw.processPendingImmediateGroups(ctx)

	return newGroupID, nil
}

// AddTasksFromExternal 外部注入新任务（方式4.3：下周期执行）
// newTaskKeys: 新任务Key列表
// newTaskExecuteFunc: 新任务执行函数（按Key映射）
// 返回：新创建的任务组ID（失败返回空）
func (fw *TaskFramework) AddTasksFromExternal(ctx context.Context, newTaskKeys []string,
	newTaskExecuteFunc map[string]func(ctx context.Context) error) (string, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// 1. 排重：过滤已存在的活跃Key
	uniqueKeys := make([]string, 0)
	for _, key := range newTaskKeys {
		if _, exist := fw.activeTaskKeys[key]; !exist {
			uniqueKeys = append(uniqueKeys, key)
		}
	}
	if len(uniqueKeys) == 0 {
		return "", errors.New("所有新任务Key均已存在，无需创建新组")
	}

	// 2. 创建新任务组（下周期执行）
	newGroupID := fmt.Sprintf("group_%d", time.Now().UnixNano())
	newTasks := make([]*Task, 0, len(uniqueKeys))
	for _, key := range uniqueKeys {
		newTasks = append(newTasks, &Task{
			Key:         key,
			FailCount:   0,
			FailReason:  "",
			ExecuteFunc: newTaskExecuteFunc[key],
		})
		// 加入活跃Key索引
		fw.activeTaskKeys[key] = struct{}{}
	}

	newGroup := &TaskGroup{
		GroupID:        newGroupID,
		Tasks:          newTasks,
		State:          TaskStatePendingCycle,
		CreatedAt:      time.Now(),
		StartTime:      time.Time{},
		WaitCycleCount: 0,
	}
	fw.groups[newGroupID] = newGroup

	return newGroupID, nil
}

// -------------------------- 辅助函数 --------------------------

// GetGroupState 获取任务组状态
func (fw *TaskFramework) GetGroupState(ctx context.Context, groupID string) (TaskState, error) {
	fw.mu.RLock()
	defer fw.mu.RUnlock()
	group, exist := fw.groups[groupID]
	if !exist {
		return "", errors.New("任务组不存在")
	}
	group.mu.RLock()
	defer group.mu.RUnlock()
	return group.State, nil
}

// Stop 停止框架（关闭定时器）
func (fw *TaskFramework) Stop(ctx context.Context) {
	fw.cycleTimer.Stop()
}
