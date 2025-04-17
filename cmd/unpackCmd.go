package cmd

import (
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

// unpackCmdMain 解压指定备份任务
func unpackCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *unpackID == 0 {
		return fmt.Errorf("解压指定备份任务时, 必须指定任务ID")
	}

	// 检查versionID是否指定
	if *unpackVersionID == "" {
		return fmt.Errorf("解压指定备份任务时, 必须指定版本ID")
	}

	// 打印提示信息
	CL.PrintOk("正在启动解压任务...")

	// 检查*unpackID是否是已存在的
	var taskCount int
	if err := db.Get(&taskCount, "SELECT COUNT(*) FROM backup_records WHERE task_id = ? AND data_status = '1';", *unpackID); err == sql.ErrNoRows {
		return fmt.Errorf("未找到指定任务ID %d 的备份记录", *unpackID)
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	} else if taskCount == 0 {
		return fmt.Errorf("未找到指定任务ID %d 的备份记录", *unpackID)
	}

	// 检查versionID是否是已存在的
	var versionCount int
	if err := db.Get(&versionCount, "SELECT COUNT(*) FROM backup_records WHERE version_id = ? AND data_status = '1';", *unpackVersionID); err == sql.ErrNoRows {
		return fmt.Errorf("未找到指定版本ID %s 的备份记录", *unpackVersionID)
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	} else if versionCount == 0 {
		return fmt.Errorf("未找到指定版本ID %s 的备份记录", *unpackVersionID)
	}

	// 构建查询sql语句
	querySql := "SELECT version_id, task_id, backup_file_name, backup_path, version_hash FROM backup_records WHERE task_id =? AND version_id =? AND data_status = '1';"

	// 定义存储查询结果的结构体
	var record struct {
		VersionID      string `db:"version_id"`       // 版本ID
		TaskID         int    `db:"task_id"`          // 任务ID
		BackupFileName string `db:"backup_file_name"` // 备份文件名
		BackupPath     string `db:"backup_path"`      // 存放路径
		VersionHash    string `db:"version_hash"`     // 版本哈希
	}

	// 执行查询
	if err := db.Get(&record, querySql, *unpackID, *unpackVersionID); err == sql.ErrNoRows {
		return fmt.Errorf("未找到指定任务ID %d 和版本ID %s 的备份记录", *unpackID, *unpackVersionID)
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	}

	// 构建备份文件路径
	backupFilePath := filepath.Join(record.BackupPath, record.BackupFileName)

	// 检查备份文件是否存在
	if _, err := tools.CheckPath(backupFilePath); err != nil {
		return fmt.Errorf("备份文件不存在: %w", err)
	}

	// 获取备份文件的后8位哈希值
	if backupFileHash, err := tools.GetFileMD5Last8(backupFilePath); err != nil {
		return fmt.Errorf("获取备份文件哈希失败: %w", err)
	} else {
		// 比较哈希值是否一致
		if backupFileHash != record.VersionHash {
			return fmt.Errorf("备份文件 %s 的版本 %s 的哈希值与记录不匹配，文件可能已损坏或被篡改。请尝试选择其他版本的备份文件重试", backupFilePath, record.VersionID)
		}
	}

	// 执行解压操作
	if unZipPath, err := tools.UncompressFilesByOS(record.BackupPath, record.BackupFileName, *unpackOutput); err != nil {
		return fmt.Errorf("解压备份文件失败: %w", err)
	} else {
		// 打印提示信息
		CL.PrintOkf("解压任务完成, 输出路径: %s", unZipPath)
		return nil
	}

}
