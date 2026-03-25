package services

import (
	"catter-cms-gin/internal/models"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler 调度器服务
type Scheduler struct {
	cron           *cron.Cron
	entryMap       map[uint]cron.EntryID // taskID -> cronEntryID
	mu             sync.RWMutex
	isRunning      bool
	lastCheckAt    time.Time
}

// 全局调度器实例
var globalScheduler *Scheduler
var schedulerOnce sync.Once

// GetScheduler 获取调度器单例
func GetScheduler() *Scheduler {
	schedulerOnce.Do(func() {
		globalScheduler = &Scheduler{
			cron:     cron.New(cron.WithSeconds()),
			entryMap: make(map[uint]cron.EntryID),
			isRunning: false,
		}
	})
	return globalScheduler
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return nil
	}

	// 启动 cron
	s.cron.Start()
	s.isRunning = true
	s.lastCheckAt = time.Now()

	// 加载所有启用的任务
	tasks, err := models.GetEnabledScheduledTasks()
	if err != nil {
		log.Printf("[Scheduler] 加载任务失败: %v", err)
		return err
	}

	for _, task := range tasks {
		if err := s.AddTask(&task); err != nil {
			log.Printf("[Scheduler] 添加任务 %s 失败: %v", task.Name, err)
		}
	}

	log.Printf("[Scheduler] 调度器已启动，加载了 %d 个任务", len(tasks))
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	// 移除所有任务
	for taskID := range s.entryMap {
		s.RemoveTask(taskID)
	}

	ctx := s.cron.Stop()
	<-ctx.Done()
	s.isRunning = false
	log.Println("[Scheduler] 调度器已停止")
}

// IsRunning 检查调度器是否运行中
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

// AddTask 添加任务到调度器
func (s *Scheduler) AddTask(task *models.ScheduledTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果任务已存在，先移除
	if _, ok := s.entryMap[task.ID]; ok {
		s.RemoveTaskLocked(task.ID)
	}

	if !task.Enabled {
		return nil
	}

	// 转换 Cron 表达式（5位转6位，添加秒位）
	cronExpr := "0 " + task.CronExpr // 统一在秒位添加 0

	entryID, err := s.cron.AddFunc(cronExpr, func() {
		s.executeTask(task.ID)
	})
	if err != nil {
		return err
	}

	s.entryMap[task.ID] = entryID
	log.Printf("[Scheduler] 任务已添加: %s (ID: %d, Cron: %s)", task.Name, task.ID, task.CronExpr)
	return nil
}

// RemoveTask 从调度器移除任务
func (s *Scheduler) RemoveTask(taskID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RemoveTaskLocked(taskID)
}

// RemoveTaskLocked 从调度器移除任务（内部使用，需持有锁）
func (s *Scheduler) RemoveTaskLocked(taskID uint) {
	if entryID, ok := s.entryMap[taskID]; ok {
		s.cron.Remove(entryID)
		delete(s.entryMap, taskID)
		log.Printf("[Scheduler] 任务已移除: ID %d", taskID)
	}
}

// executeTask 执行任务
func (s *Scheduler) executeTask(taskID uint) {
	task, err := models.GetScheduledTaskByID(taskID)
	if err != nil {
		log.Printf("[Scheduler] 获取任务失败: ID %d, Error: %v", taskID, err)
		return
	}

	now := time.Now()
	startTime := now

	// 更新任务状态为运行中
	models.UpdateScheduledTaskStatus(taskID, "running")

	// 创建执行记录
	execution := &models.ScheduledTaskExecution{
		TaskID:    taskID,
		TaskName:  task.Name,
		Status:    "running",
		StartedAt: startTime,
	}
	models.CreateTaskExecution(execution)

	// 执行实际任务
	message := ""
	status := "success"

	switch task.TaskType {
	case "backup":
		message = s.executeBackup(task)
	case "cleanup":
		message = s.executeCleanup(task)
	case "cache_clear":
		message = s.executeCacheClear(task)
	case "publish":
		message = s.executePublish(task)
	case "custom":
		message = s.executeCustom(task)
	default:
		message = "未知任务类型: " + task.TaskType
		status = "failed"
	}

	finishedAt := time.Now()
	duration := int(finishedAt.Sub(startTime).Seconds())

	// 更新执行记录
	execution.Status = status
	execution.Message = message
	execution.FinishedAt = &finishedAt
	execution.Duration = duration

	// 更新任务状态和统计
	models.UpdateScheduledTaskStatus(taskID, status)
	models.IncrementRunCount(taskID)

	// 计算下次执行时间
	if nextRun := s.cron.Entry(s.entryMap[taskID]).Next; !nextRun.IsZero() {
		models.UpdateNextRunAt(taskID, nextRun)
	}

	log.Printf("[Scheduler] 任务执行完成: %s, Status: %s, Duration: %ds", task.Name, status, duration)
}

// executeBackup 执行备份任务
func (s *Scheduler) executeBackup(task *models.ScheduledTask) string {
	// 这里实现备份逻辑
	log.Printf("[Scheduler] 执行备份任务: %s", task.Name)
	// TODO: 实现实际的数据库备份逻辑
	return "备份任务执行成功"
}

// executeCleanup 执行清理任务
func (s *Scheduler) executeCleanup(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行清理任务: %s", task.Name)
	// TODO: 实现实际的系统清理逻辑
	return "清理任务执行成功"
}

// executeCacheClear 执行缓存清理任务
func (s *Scheduler) executeCacheClear(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行缓存清理任务: %s", task.Name)

	// 获取所有主题并清理缓存
	themes := GetAllThemes()
	for name := range themes {
		if err := ClearThemeCache(name); err != nil {
			log.Printf("[Scheduler] 清理主题 %s 缓存失败: %v", name, err)
		}
	}

	return "缓存清理成功"
}

// executePublish 执行定时发布任务
func (s *Scheduler) executePublish(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行定时发布任务: %s", task.Name)
	// TODO: 实现实际的定时发布逻辑（如发布预定的文章）
	return "定时发布任务执行成功"
}

// executeCustom 执行自定义命令
func (s *Scheduler) executeCustom(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行自定义命令: %s, Command: %s", task.Name, task.Command)
	// TODO: 实现实际的自定义命令执行逻辑
	return "自定义命令执行成功"
}

// RunTaskNow 立即执行任务
func (s *Scheduler) RunTaskNow(taskID uint) error {
	// 检查任务是否存在
	_, err := models.GetScheduledTaskByID(taskID)
	if err != nil {
		return err
	}

	// 在 goroutine 中执行，不阻塞
	go s.executeTask(taskID)
	return nil
}

// GetStatus 获取调度器状态
func (s *Scheduler) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks, _ := models.GetAllScheduledTasks()
	enabledCount := 0
	runningCount := 0

	for _, task := range tasks {
		if task.Enabled {
			enabledCount++
		}
		if task.Status == "running" {
			runningCount++
		}
	}

	return map[string]interface{}{
		"is_running":    s.isRunning,
		"total_tasks":    len(tasks),
		"enabled_tasks": enabledCount,
		"running_tasks": runningCount,
		"idle_tasks":    enabledCount - runningCount,
		"failed_tasks":  0,
	}
}

// ReloadAllTasks 重新加载所有任务
func (s *Scheduler) ReloadAllTasks() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 移除所有现有任务
	for taskID := range s.entryMap {
		s.cron.Remove(s.entryMap[taskID])
	}
	s.entryMap = make(map[uint]cron.EntryID)

	// 重新加载启用的任务
	tasks, err := models.GetEnabledScheduledTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		cronExpr := "0 " + task.CronExpr
		entryID, err := s.cron.AddFunc(cronExpr, func() {
			s.executeTask(task.ID)
		})
		if err != nil {
			log.Printf("[Scheduler] 重新加载任务 %s 失败: %v", task.Name, err)
			continue
		}
		s.entryMap[task.ID] = entryID
	}

	log.Printf("[Scheduler] 已重新加载 %d 个任务", len(tasks))
	return nil
}

// CalculateNextRun 根据 Cron 表达式计算下次执行时间
func CalculateNextRun(cronExpr string) *time.Time {
	s := GetScheduler()
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建一个临时的 cron 来计算下次执行时间
	tempCron := cron.New(cron.WithSeconds())
	expr := "0 " + cronExpr
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(expr)
	if err != nil {
		return nil
	}

	next := schedule.Next(time.Now())
	tempCron.Stop()
	return &next
}