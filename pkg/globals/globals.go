package globals

import (
	"os"
	"path/filepath"
)

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
	RetentionDays   int    `db:"retention_days"`   // 保留天数
	NoCompression   int    `db:"no_compression"`   // 是否禁用压缩(默认启用压缩, 0 表示启用压缩, 1 表示禁用压缩)
	ExcludeRules    string `db:"exclude_rules"`    // 排除规则
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

// 定义任务配置的结构体
type TaskConfig struct {
	Task Task `yaml:"task"`
}

// 定义任务的结构体
type Task struct {
	Name          string    `yaml:"name"`            // 任务名
	Target        string    `yaml:"target"`          // 目标目录
	Backup        string    `yaml:"backup"`          // 备份目录
	Retention     Retention `yaml:"retention"`       // 保留策略
	BackupDirName string    `yaml:"backup_dir_name"` // 备份目录名
	NoCompression int       `yaml:"no_compression"`  // 是否禁用压缩(默认启用压缩, 0 表示启用压缩, 1 表示禁用压缩)
	ExcludeRules  string    `yaml:"exclude_rules"`   // 排除规则
}

// 定义保留策略的结构体
type Retention struct {
	Count int `yaml:"count"` // 保留数量
	Days  int `yaml:"days"`  // 保留天数(配置为 0 表示不限制天数)
}

// 定义排除函数的类型
// ExcludeFunc 是一个函数类型，用于排除文件或目录
// path: 文件或目录的路径
// info: 文件或目录的信息
// 返回值: 如果文件或目录应该被排除, 返回 true, 否则返回 false
type ExcludeFunc func(path string, info os.FileInfo) bool

// 全局排除函数变量
var (
	// NoExcludeFunc 是一个空的排除函数，表示不排除任何文件或目录
	NoExcludeFunc ExcludeFunc = func(path string, info os.FileInfo) bool {
		return false // 不排除任何文件或目录
	}
)
