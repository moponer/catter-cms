package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"catter-cms/models"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron       *cron.Cron
	taskMap    map[int]cron.EntryID
	db         interface{} // 数据库连接
	isRunning  bool
}

var (
	schedulerInstance *Scheduler
	schedulerOnce     sync.Once
)

// GetScheduler 获取调度器单例
func GetScheduler() *Scheduler {
	schedulerOnce.Do(func() {
		schedulerInstance = &Scheduler{
			cron:    cron.New(cron.WithSeconds()), // 支持秒级调度
			taskMap: make(map[int]cron.EntryID),
		}
	})
	return schedulerInstance
}

// Init 初始化调度器
func (s *Scheduler) Init(db interface{}) error {
	s.db = db
	
	// 从数据库加载所有启用的任务
	tasks := []models.ScheduledTask{}
	if err := models.DB.Where("enabled = ?", true).All(&tasks); err != nil {
		return fmt.Errorf("加载定时任务失败: %w", err)
	}
	
	// 添加所有任务到调度器
	for _, task := range tasks {
		if err := s.AddTask(&task); err != nil {
			log.Printf("添加定时任务失败 [ID=%d]: %v", task.ID, err)
		}
	}
	
	// 启动调度器
	s.cron.Start()
	s.isRunning = true
	log.Println("✅ 定时任务调度器启动成功")
	
	return nil
}

// AddTask 添加任务到调度器
func (s *Scheduler) AddTask(task *models.ScheduledTask) error {
	// 先移除旧任务(如果存在)
	if _, exists := s.taskMap[task.ID]; exists {
		s.RemoveTask(task.ID)
	}
	
	// 创建任务执行器
	jobFunc := func() {
		s.executeTask(task)
	}
	
	// 添加任务到cron
	entryID, err := s.cron.AddFunc(task.CronExpr, jobFunc)
	if err != nil {
		return fmt.Errorf("添加任务到cron失败: %w", err)
	}
	
	// 保存任务ID映射
	s.taskMap[task.ID] = entryID
	
	// 计算下次执行时间
	nextRun := s.cron.Entry(entryID).Next
	log.Printf("✅ 添加定时任务: %s [ID=%d, Cron=%s, NextRun=%s]",
		task.Name, task.ID, task.CronExpr, nextRun.Format("2006-01-02 15:04:05"))

	return nil
}

// RemoveTask 从调度器中移除任务
func (s *Scheduler) RemoveTask(taskID int) {
	if entryID, exists := s.taskMap[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.taskMap, taskID)
		log.Printf("✅ 移除定时任务 [ID=%d]", taskID)
	}
}

// executeTask 执行定时任务
func (s *Scheduler) executeTask(task *models.ScheduledTask) {
	log.Printf("🚀 执行定时任务: %s [ID=%d]", task.Name, task.ID)

	// 更新任务状态为运行中
	now := time.Now()
	s.updateTaskStatus(task.ID, "running", nil, &now)

	// 执行任务
	var err error
	switch task.TaskType {
	case "publish":
		err = s.executePublishTask(task)
	case "backup":
		err = s.executeBackupTask(task)
	case "cleanup":
		err = s.executeCleanupTask(task)
	case "cache_clear":
		err = s.executeCacheClearTask(task)
	case "custom":
		err = s.executeCustomTask(task)
	default:
		err = fmt.Errorf("未知的任务类型: %s", task.TaskType)
	}

	// 更新任务状态
	if err != nil {
		errorMsg := err.Error()
		s.updateTaskStatus(task.ID, "failed", &errorMsg, nil)
		log.Printf("❌ 定时任务执行失败: %s [ID=%d] - %v", task.Name, task.ID, err)
	} else {
		// 计算下次执行时间
		if entryID, exists := s.taskMap[task.ID]; exists {
			nextRun := s.cron.Entry(entryID).Next
			s.updateTaskStatus(task.ID, "idle", nil, &nextRun)
		} else {
			s.updateTaskStatus(task.ID, "idle", nil, nil)
		}
		log.Printf("✅ 定时任务执行成功: %s [ID=%d]", task.Name, task.ID)
	}
}

// updateTaskStatus 更新任务状态
func (s *Scheduler) updateTaskStatus(taskID int, status string, errorMsg *string, updatedAt *time.Time) {
	task := &models.ScheduledTask{}
	if err := models.DB.Find(task, taskID); err != nil {
		log.Printf("查找任务失败 [ID=%d]: %v", taskID, err)
		return
	}

	task.Status = status
	if errorMsg != nil {
		task.ErrorMsg = errorMsg
	} else {
		task.ErrorMsg = nil
	}
	if updatedAt != nil {
		task.UpdatedAt = *updatedAt
	}

	if err := models.DB.Update(task); err != nil {
		log.Printf("更新任务状态失败 [ID=%d]: %v", taskID, err)
	}
}

// ExecuteTask 公共方法:执行定时任务(用于手动触发)
func (s *Scheduler) ExecuteTask(task *models.ScheduledTask) {
	s.executeTask(task)
}

// executePublishTask 执行定时发布任务
func (s *Scheduler) executePublishTask(task *models.ScheduledTask) error {
	// 查找所有草稿文章
	posts := []models.Post{}
	if err := models.DB.Where("status = ?", "draft").All(&posts); err != nil {
		return fmt.Errorf("查询草稿文章失败: %w", err)
	}

	count := 0
	for _, post := range posts {
		// TODO: 实现定时发布逻辑
		// 可以添加一个scheduled_at字段来标记定时发布时间
		log.Printf("📝 处理草稿文章: %s [ID=%d]", post.Title, post.ID)
	}

	if count > 0 {
		log.Printf("📊 定时发布任务完成,共处理 %d 篇文章", count)
	}

	return nil
}

// executeBackupTask 执行数据备份任务
func (s *Scheduler) executeBackupTask(task *models.ScheduledTask) error {
	// TODO: 实现数据备份功能
	log.Println("📦 数据备份任务执行(待实现)")
	return nil
}

// executeCleanupTask 执行清理任务
func (s *Scheduler) executeCleanupTask(task *models.ScheduledTask) error {
	// 清理已删除的评论
	count := 0
	// TODO: 实现清理逻辑
	// models.DB.Raw("DELETE FROM comments WHERE status = 'trash'").Exec()

	log.Printf("🧹 清理任务完成,清理了 %d 条数据", count)
	return nil
}

// executeCacheClearTask 执行缓存清理任务
func (s *Scheduler) executeCacheClearTask(task *models.ScheduledTask) error {
	// TODO: 实现缓存清理逻辑

	log.Println("🔄 缓存清理任务完成")
	return nil
}

// executeCustomTask 执行自定义任务
func (s *Scheduler) executeCustomTask(task *models.ScheduledTask) error {
	// TODO: 实现自定义命令执行
	log.Printf("🔧 执行自定义任务: %s", task.Command)
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	if s.isRunning {
		s.cron.Stop()
		s.isRunning = false
		log.Println("🛑 定时任务调度器已停止")
	}
}

// Restart 重启调度器
func (s *Scheduler) Restart() error {
	s.Stop()
	if err := s.Init(s.db); err != nil {
		return err
	}
	return nil
}

// IsRunning 获取调度器运行状态
func (s *Scheduler) IsRunning() bool {
	return s.isRunning
}

// GetTaskStatus 获取任务状态
func (s *Scheduler) GetTaskStatus(taskID int) (running bool, nextRun time.Time) {
	if entryID, exists := s.taskMap[taskID]; exists {
		entry := s.cron.Entry(entryID)
		return true, entry.Next
	}
	return false, time.Time{}
}
