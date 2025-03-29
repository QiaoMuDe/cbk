package models

import (
	"time"
)

// 备份记录表
type BackupRecord struct {
	VersionID      int       `gorm:"primaryKey;column:version_id"`     // 版本ID，主键
	TaskID         int       `gorm:"foreignKey:TaskID;column:task_id"` // 任务ID，外键
	Timestamp      time.Time `gorm:"column:timestamp"`                 // 时间戳
	TaskName       string    `gorm:"column:task_name"`                 // 任务名
	BackupStatus   bool      `gorm:"column:backup_status"`             // 备份状态
	BackupFileName string    `gorm:"column:backup_file_name"`          // 备份文件名
	BackupSize     string    `gorm:"column:backup_size"`               // 备份大小
	BackupPath     string    `gorm:"column:backup_path"`               // 备份路径
	VersionHash    string    `gorm:"column:version_hash"`              // 版本hash
}

// 存储备份任务表
type BackupTask struct {
	TaskID      int    `gorm:"primaryKey;autoIncrement;column:task_id"` // 任务ID，主键，自动增长
	TaskName    string `gorm:"column:task_name"`                        // 任务名
	TargetDir   string `gorm:"column:target_dir"`                       // 目标目录
	BackupDir   string `gorm:"column:backup_dir"`                       // 备份目录
	RetainCount int    `gorm:"default:3;column:retain_count"`           // 保留数量，默认值为 3
}
