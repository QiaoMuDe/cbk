package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
)

// runCmdMain 运行备份任务
func runCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *runID == 0 {
		return fmt.Errorf("运行备份任务时, 必须指定任务ID")
	}

	// 构建存储查询任务信息的结构体
	var task globals.BackupTask

	// 构建查询任务信息的SQL语句
	querySql := "select task_name, target_directory, backup_directory, retention_count, no_compression from backup_tasks where task_id =?"

	// 查询任务信息
	if err := db.Get(&task, querySql, *runID); err == sql.ErrNoRows {
		return fmt.Errorf("任务ID不存在 %d", *runID)
	} else if err != nil {
		return fmt.Errorf("获取任务信息失败: %w", err)
	}

	// 检查目标目录或文件是否存在
	if _, err := tools.CheckPath(task.TargetDirectory); err != nil {
		return fmt.Errorf("目标目录或文件不存在: %w", err)
	}

	// 检查备份目录是否存在
	if _, err := tools.CheckPath(task.BackupDirectory); err != nil {
		if err := os.MkdirAll(task.BackupDirectory, 0755); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}
	}

	// 打印提示信息
	CL.PrintOkf("备份任务 [%s] 已启动，正在运行中……", task.TaskName)

	// 构建备份文件名
	backupTime := time.Now().Format("20060102150405")
	backupFileNamePrefix := fmt.Sprintf("%s_%s", task.TaskName, backupTime)

	// 获取versionID
	versionID := tools.GenerateID(6)

	// 运行备份任务
	targetDir := filepath.Dir(task.TargetDirectory)                                 // 获取目标目录的目录部分
	targetName := filepath.Base(task.TargetDirectory)                               // 获取目标目录的最后一个部分
	backupFileNamePath := filepath.Join(task.BackupDirectory, backupFileNamePrefix) // 获取构建的备份文件路径

	// 执行备份任务
	zipPath, err := tools.CreateZipFromOSPaths(db, targetDir, targetName, backupFileNamePath, task.NoCompression)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("备份任务失败: %w", err)
	}

	// 获取备份文件的后8位MD5哈希值
	backupFileMD5, err := tools.GetFileMD5Last8(zipPath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("获取备份文件MD5失败: %w", err)
	}

	// 获取备份文件的大小
	backupFileSize, err := tools.HumanReadableSize(zipPath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("获取备份文件大小失败: %w", err)
	}

	// 插入备份记录
	insertSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if _, err := db.Exec(insertSql, versionID, *runID, backupTime, task.TaskName, "true", filepath.Base(zipPath), backupFileSize, task.BackupDirectory, backupFileMD5); err != nil {
		return fmt.Errorf("插入备份记录失败: %w", err)
	}

	// 获取备份目录下的以指定扩展名的文件列表
	zipFiles, err := tools.GetZipFiles(task.BackupDirectory, ".zip")
	if err != nil {
		return fmt.Errorf("获取备份目录下的.zip文件失败: %w", err)
	}

	// 删除多余的备份文件
	if len(zipFiles) > task.RetentionCount {
		if err := tools.RetainLatestFiles(db, zipFiles, task.RetentionCount); err != nil {
			return fmt.Errorf("删除多余的备份文件失败: %w", err)
		}
	}

	// 打印成功信息
	CL.PrintOk(`备份成功!`)

	return nil
}
