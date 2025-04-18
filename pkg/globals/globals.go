package globals

import "path/filepath"

const (
	CbkHomeDir = ".cbk"   // 数据目录
	CbkDBFile  = "cbk.db" // 数据库文件
	CbkDataDir = "data"   // 数据目录
)

// 数据库文件路径
var CbkDbPath = filepath.Join(CbkHomeDir, CbkDBFile)

// 定义任务表结构体
type BackupTask struct {
	TaskID          int    `db:"task_id"`          // 任务ID
	TaskName        string `db:"task_name"`        // 任务名
	TargetDirectory string `db:"target_directory"` // 目标目录
	BackupDirectory string `db:"backup_directory"` // 备份目录
	RetentionCount  int    `db:"retention_count"`  // 保留数量
	NoCompression   int    `db:"no_compression"`   // 是否禁用压缩(默认启用压缩, 0 表示启用压缩, 1 表示禁用压缩)
}

// 定义任务表结构体切片
type BackupTasks []BackupTask

// 定义备份记录表结构体
type BackupRecord struct {
	VersionID      string `db:"version_id"`       // 版本ID
	TaskID         int    `db:"task_id"`          // 任务ID
	Timestamp      string `db:"timestamp"`        // 时间戳
	TaskName       string `db:"task_name"`        // 任务名
	BackupStatus   string `db:"backup_status"`    // 备份状态
	BackupFileName string `db:"backup_file_name"` // 备份文件名
	BackupSize     string `db:"backup_size"`      // 备份文件大小
	BackupPath     string `db:"backup_path"`      // 备份文件路径
	VersionHash    string `db:"version_hash"`     // 版本哈希
}

// 定义备份记录表结构体切片
type BackupRecords []BackupRecord
