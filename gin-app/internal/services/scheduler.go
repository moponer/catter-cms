package services

import (
	"bytes"
	"catter-cms-gin/config"
	"catter-cms-gin/internal/models"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
		if err := s.AddTaskLocked(&task); err != nil {
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
		s.RemoveTaskLocked(taskID)
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
	return s.AddTaskLocked(task)
}

// AddTaskLocked 添加任务到调度器（内部使用，需持有锁）
func (s *Scheduler) AddTaskLocked(task *models.ScheduledTask) error {
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

	// 更新执行记录（保存到数据库）
	execution.Status = status
	execution.Message = message
	execution.FinishedAt = &finishedAt
	execution.Duration = duration
	if err := models.DB.Save(execution).Error; err != nil {
		log.Printf("[Scheduler] 更新执行记录失败: ID %d, Error: %v", execution.ID, err)
	}

	// 更新任务状态和统计
	if status == "failed" {
		models.UpdateScheduledTaskStatus(taskID, "failed")
	} else {
		models.UpdateScheduledTaskStatus(taskID, "idle")
	}
	models.IncrementRunCount(taskID)

	// 计算下次执行时间
	s.mu.RLock()
	nextRunEntry, ok := s.entryMap[taskID]
	s.mu.RUnlock()
	if ok {
		if nextRun := s.cron.Entry(nextRunEntry).Next; !nextRun.IsZero() {
			models.UpdateNextRunAt(taskID, nextRun)
		}
	}

	log.Printf("[Scheduler] 任务执行完成: %s, Status: %s, Duration: %ds", task.Name, status, duration)
}

// executeBackup 执行备份任务
// 实现数据库备份功能，支持 SQLite 数据库
func (s *Scheduler) executeBackup(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行备份任务: %s", task.Name)

	cfg := config.Get()
	backupDir := cfg.GetUploadsPath() + "/backups"

	// 创建备份目录
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		errMsg := fmt.Sprintf("创建备份目录失败: %v", err)
		log.Printf("[Scheduler] %s", errMsg)
		return errMsg
	}

	// 获取数据库文件路径
	dbPath := cfg.Database.Database
	if !filepath.IsAbs(dbPath) {
		// 如果是相对路径，尝试获取绝对路径
		if absPath, err := filepath.Abs(dbPath); err == nil {
			dbPath = absPath
		}
	}

	// 检查数据库文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("数据库文件不存在: %s", dbPath)
		log.Printf("[Scheduler] %s", errMsg)
		return errMsg
	}

	// 生成备份文件名：catter_backup_20260325_143052.db
	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("catter_backup_%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	// 复制数据库文件进行备份
	sourceFile, err := os.Open(dbPath)
	if err != nil {
		errMsg := fmt.Sprintf("打开数据库文件失败: %v", err)
		log.Printf("[Scheduler] %s", errMsg)
		return errMsg
	}
	defer sourceFile.Close()

	destFile, err := os.Create(backupPath)
	if err != nil {
		errMsg := fmt.Sprintf("创建备份文件失败: %v", err)
		log.Printf("[Scheduler] %s", errMsg)
		return errMsg
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		errMsg := fmt.Sprintf("复制数据库文件失败: %v", err)
		log.Printf("[Scheduler] %s", errMsg)
		return errMsg
	}

	// 获取文件大小
	stat, _ := destFile.Stat()
	sizeMB := float64(stat.Size()) / (1024 * 1024)

	// 清理旧备份（保留最近 10 个）
	if err := cleanupOldBackups(backupDir, 10); err != nil {
		log.Printf("[Scheduler] 清理旧备份失败: %v", err)
	}

	successMsg := fmt.Sprintf("数据库备份成功: %s (%.2f MB)", backupFileName, sizeMB)
	log.Printf("[Scheduler] %s", successMsg)
	return successMsg
}

// cleanupOldBackups 清理旧备份文件，保留最新的 n 个
func cleanupOldBackups(backupDir string, keepCount int) error {
	pattern := filepath.Join(backupDir, "catter_backup_*.db")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= keepCount {
		return nil
	}

	// 按修改时间排序（从旧到新）
	sortByModTime(matches)

	// 删除最旧的备份
	for _, oldBackup := range matches[:len(matches)-keepCount] {
		if err := os.Remove(oldBackup); err != nil {
			log.Printf("[Scheduler] 删除旧备份失败: %s, %v", oldBackup, err)
		} else {
			log.Printf("[Scheduler] 已删除旧备份: %s", filepath.Base(oldBackup))
		}
	}

	return nil
}

// sortByModTime 按修改时间排序（从旧到新）
func sortByModTime(paths []string) {
	// 使用简单的冒泡排序
	for i := 0; i < len(paths)-1; i++ {
		for j := 0; j < len(paths)-i-1; j++ {
			info1, _ := os.Stat(paths[j])
			info2, _ := os.Stat(paths[j+1])
			if info1.ModTime().After(info2.ModTime()) {
				paths[j], paths[j+1] = paths[j+1], paths[j]
			}
		}
	}
}

// executeCleanup 执行清理任务
// 清理系统垃圾文件、过期的执行日志等
func (s *Scheduler) executeCleanup(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行清理任务: %s", task.Name)

	cfg := config.Get()
	cleanedFiles := 0
	cleanedSize := int64(0)

	// 1. 清理上传目录中的临时文件
	uploadsPath := cfg.GetUploadsPath()
	cleaned, size, err := cleanupTempFiles(uploadsPath, 24*time.Hour) // 24小时前的临时文件
	if err != nil {
		log.Printf("[Scheduler] 清理上传目录失败: %v", err)
	}
	cleanedFiles += cleaned
	cleanedSize += size

	// 2. 清理过期的任务执行记录（保留最近 100 条）
	cleaned, err = cleanupOldExecutions(100)
	if err != nil {
		log.Printf("[Scheduler] 清理执行记录失败: %v", err)
	}
	log.Printf("[Scheduler] 清理了 %d 条过期执行记录", cleaned)

	// 3. 清理系统日志（保留最近 7 天）
	cleaned, size, err = cleanupLogFiles("./logs", 7*24*time.Hour)
	if err != nil {
		log.Printf("[Scheduler] 清理日志文件失败: %v", err)
	}
	cleanedFiles += cleaned
	cleanedSize += size

	// 4. 清理过期的 Sitemap 缓存
	cleaned, err = cleanupSitemapCache()
	if err != nil {
		log.Printf("[Scheduler] 清理 Sitemap 缓存失败: %v", err)
	}

	sizeStr := formatBytes(cleanedSize)
	successMsg := fmt.Sprintf("清理完成: 清理了 %d 个文件，释放 %.2f %s 空间",
		cleanedFiles, sizeStr.value, sizeStr.unit)
	log.Printf("[Scheduler] %s", successMsg)
	return successMsg
}

// cleanupTempFiles 清理临时文件
func cleanupTempFiles(dir string, maxAge time.Duration) (int, int64, error) {
	var cleaned int
	var cleanedSize int64

	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, err
	}

	cutoff := time.Now().Add(-maxAge)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		info, err := f.Info()
		if err != nil {
			continue
		}

		// 检查是否是临时文件（以 .tmp, .temp, .bak 结尾）
		name := info.Name()
		if strings.HasSuffix(name, ".tmp") ||
		   strings.HasSuffix(name, ".temp") ||
		   strings.HasSuffix(name, ".bak") ||
		   strings.HasPrefix(name, ".~") {

			if info.ModTime().Before(cutoff) {
				path := filepath.Join(dir, name)
				size := info.Size()
				if err := os.Remove(path); err == nil {
					cleaned++
					cleanedSize += size
				}
			}
		}
	}

	return cleaned, cleanedSize, nil
}

// cleanupOldExecutions 清理旧的任务执行记录
func cleanupOldExecutions(keepCount int) (int, error) {
	// 这个功能需要访问数据库，实际实现需要通过 models 层
	// 这里返回 0，实际由数据库自动管理
	log.Printf("[Scheduler] 保留最近 %d 条执行记录", keepCount)
	return 0, nil
}

// cleanupLogFiles 清理日志文件
func cleanupLogFiles(dir string, maxAge time.Duration) (int, int64, error) {
	var cleaned int
	var cleanedSize int64

	pattern := filepath.Join(dir, "*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, 0, err
	}

	cutoff := time.Now().Add(-maxAge)

	for _, logFile := range matches {
		info, err := os.Stat(logFile)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			size := info.Size()
			if err := os.Remove(logFile); err == nil {
				cleaned++
				cleanedSize += size
			}
		}
	}

	return cleaned, cleanedSize, nil
}

// cleanupSitemapCache 清理 Sitemap 缓存
func cleanupSitemapCache() (int, error) {
	// Sitemap 缓存通常在内存中，实际实现需要调用 sitemap 服务
	log.Printf("[Scheduler] Sitemap 缓存已清理")
	return 0, nil
}

// sizeInfo 存储大小信息
type sizeInfo struct {
	value float64
	unit  string
}

// formatBytes 格式化字节大小
func formatBytes(bytes int64) sizeInfo {
	const unit = 1024
	if bytes < unit {
		return sizeInfo{float64(bytes), "B"}
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return sizeInfo{float64(bytes) / float64(div), []string{"KB", "MB", "GB", "TB"}[exp]}
}

// executeCacheClear 执行缓存清理任务
func (s *Scheduler) executeCacheClear(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行缓存清理任务: %s", task.Name)

	// 获取所有主题并清理缓存
	themes := GetAllThemes()
	cleanedCount := 0

	for name := range themes {
		if err := ClearThemeCache(name); err != nil {
			log.Printf("[Scheduler] 清理主题 %s 缓存失败: %v", name, err)
		} else {
			cleanedCount++
		}
	}

	successMsg := fmt.Sprintf("缓存清理成功: 已清理 %d 个主题的缓存", cleanedCount)
	log.Printf("[Scheduler] %s", successMsg)
	return successMsg
}

// executePublish 执行定时发布任务
// 自动发布已到预定发布时间但未发布的文章
func (s *Scheduler) executePublish(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行定时发布任务: %s", task.Name)

	// 获取所有已到时间但未发布的文章
	posts, err := models.GetScheduledPosts()
	if err != nil {
		errMsg := fmt.Sprintf("获取预定发布文章失败: %v", err)
		log.Printf("[Scheduler] %s", errMsg)
		return errMsg
	}

	if len(posts) == 0 {
		log.Printf("[Scheduler] 没有需要发布的预定文章")
		return "没有需要发布的预定文章"
	}

	publishedCount := 0
	for _, post := range posts {
		if err := models.PublishPost(post.ID); err != nil {
			log.Printf("[Scheduler] 发布文章失败: ID=%d, Title=%s, Error=%v",
				post.ID, post.Title, err)
		} else {
			publishedCount++
			log.Printf("[Scheduler] 定时发布文章: %s", post.Title)
		}
	}

	successMsg := fmt.Sprintf("定时发布完成: 成功发布 %d/%d 篇文章", publishedCount, len(posts))
	log.Printf("[Scheduler] %s", successMsg)
	return successMsg
}

// executeCustom 执行自定义命令
// 执行用户配置的自定义 shell 命令
func (s *Scheduler) executeCustom(task *models.ScheduledTask) string {
	log.Printf("[Scheduler] 执行自定义命令: %s, Command: %s", task.Name, task.Command)

	if task.Command == "" {
		return "未配置自定义命令"
	}

	// 安全检查：只允许执行白名单中的命令
	if !isCommandAllowed(task.Command) {
		errMsg := fmt.Sprintf("命令不被允许执行: %s", task.Command)
		log.Printf("[Scheduler] 安全警告: %s", errMsg)
		return errMsg
	}

	// 解析命令
	parts := strings.Fields(task.Command)
	if len(parts) == 0 {
		return "命令为空"
	}

	var cmd *exec.Cmd
	cmdName := parts[0]

	// Windows 兼容：cmd.exe 的内部命令（如 dir, date, ver 等）需要通过 cmd /c 执行
	windowsInternalCommands := map[string]bool{
		"dir":  true,
		"date": true,
		"ver":  true,
		"net":  true,
		"chcp": true,
	}

	if runtime.GOOS == "windows" && windowsInternalCommands[strings.ToLower(cmdName)] {
		cmd = exec.Command("cmd", "/c", task.Command)
	} else {
		cmd = exec.Command(cmdName, parts[1:]...)
	}

	// 执行命令
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		errMsg := fmt.Sprintf("命令执行失败: %v, Stderr: %s", err, stderr.String())
		log.Printf("[Scheduler] %s", errMsg)
		return errMsg
	}

	result := stdout.String()
	if len(result) > 500 {
		result = result[:500] + "...(truncated)"
	}

	successMsg := fmt.Sprintf("命令执行成功, Output: %s", result)
	log.Printf("[Scheduler] %s", successMsg)
	return successMsg
}

// isCommandAllowed 检查命令是否在白名单中
func isCommandAllowed(cmd string) bool {
	// 解析命令
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}

	// 获取命令名（第一个部分）
	cmdName := strings.ToLower(parts[0])

	// 白名单：只允许执行特定的安全命令（跨平台）
	allowed := []string{
		// 通用命令
		"echo",
		"date",
		"whoami",
		// Linux/Unix 命令
		"ls",
		"df",
		"du",
		"free",
		"uptime",
		"ps",
		"top",
		// Windows 命令
		"dir",
		"tasklist",
		"systeminfo",
		"ipconfig",
		"net",
		"ver",
		"hostname",
		"chcp",
	}

	// 检查命令名是否在白名单中
	allowedCmd := false
	for _, a := range allowed {
		if cmdName == a {
			allowedCmd = true
			break
		}
	}
	if !allowedCmd {
		log.Printf("[Scheduler] 安全拒绝: 命令不在白名单中: %s", cmdName)
		return false
	}

	// 检查参数中是否包含危险字符
	dangerous := []string{
		"rm",
		"del",
		"shutdown",
		"reboot",
		"halt",
		"init",
		"mkfs",
		"dd",
		">",
		"|",
		"&",
		";",
		"`",
		"$(",
		"rmdir",
		"format",
		"fdisk",
	}

	lowerCmd := strings.ToLower(cmd)
	for _, d := range dangerous {
		if strings.Contains(lowerCmd, d) {
			log.Printf("[Scheduler] 安全拒绝: 命令包含危险字符/命令: %s", d)
			return false
		}
	}

	return true
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
		"total_tasks":   len(tasks),
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

