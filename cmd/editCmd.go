package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

// editCmdMain 编辑任务
func editCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *editID == -1 {
		return fmt.Errorf("编辑任务时, 必须指定任务ID")
	}

	// 检查所有的参数是否都没指定
	if *editName == "" && *editRetentionCount == -1 && *editRetentionDays == -1 && *editNoCompression == "" && *editNewDirName == "" {
		CL.PrintWarn("未指定任何参数, 任务将不会被修改")
		return nil
	}

	// 查询任务信息
	editSql := "select task_name, retention_count, retention_days, backup_directory, no_compression from backup_tasks where task_id =?"
	var task globals.BackupTask
	if err := db.Get(&task, editSql, *editID); err == sql.ErrNoRows {
		return fmt.Errorf("任务ID不存在 %d", *editID)
	} else if err != nil {
		return fmt.Errorf("查询任务失败: %w, SQL: %s, ID: %d", err, editSql, *editID)
	}

	// 如果指定了-n参数, 则更新任务名
	if *editName != "" {
		task.TaskName = *editName
	}

	// 如果指定了-c参数, 则更新保留数量
	if *editRetentionCount != -1 {
		task.RetentionCount = *editRetentionCount
	}

	// 如果指定了-d参数, 则更新保留天数
	if *editRetentionDays != -1 { // 保留天数默认为0
		task.RetentionDays = *editRetentionDays
	}

	// 如果指定了-nc参数, 则更新是否禁用压缩
	if *editNoCompression != "" {
		// 检查如果不是true或false则报错
		if *editNoCompression != "true" && *editNoCompression != "false" {
			return fmt.Errorf("参数 -nc 只能是 true 或 false")
		}

		// 根据参数值更新NoCompression字段
		if *editNoCompression == "true" {
			task.NoCompression = 1 // 1 表示禁用压缩
		} else {
			task.NoCompression = 0 // 0表示启用压缩
		}
	}

	// 如果指定了-bn参数, 则更新备份目录
	var oldDirName, rootPath, newDirName string
	if *editNewDirName != "" {
		newDirName = *editNewDirName

		// 检查备份目录名是否非法字符
		if tools.ContainsSpecialChars(newDirName) {
			return fmt.Errorf("备份目录名含非法字符, 请重试")
		}

		rootPath = filepath.Dir(task.BackupDirectory)    // 获取备份目录的根路径
		oldDirName = filepath.Base(task.BackupDirectory) // 获取备份目录的旧名称

		// 重命名备份目录
		if err := tools.RenameBackupDirectory(rootPath, oldDirName, newDirName); err != nil {
			return err
		}

		// 更新备份目录路径
		task.BackupDirectory = filepath.Join(rootPath, newDirName)
	}

	// 更新任务
	updateSql := "update backup_tasks set task_name = ?, retention_count = ? , retention_days = ?, backup_directory = ?, no_compression = ? where task_id = ?"

	// 更新任务SQL
	if _, err := db.Exec(updateSql, task.TaskName, task.RetentionCount, task.RetentionDays, task.BackupDirectory, task.NoCompression, *editID); err != nil {
		// 更新任务失败
		if *editNewDirName != "" {
			if err := tools.RenameBackupDirectory(rootPath, newDirName, oldDirName); err != nil {
				return fmt.Errorf("更新任务失败且恢复备份目录失败: %w", err)
			}
			CL.PrintOkf("更新任务失败, 已恢复备份目录: %s", filepath.Join(rootPath, oldDirName))
		}
		return fmt.Errorf("更新任务失败: %w, SQL: %s, ID: %d", err, updateSql, *editID)
	}

	// 打印成功信息
	CL.PrintOk("更新成功!")
	if *editName != "" {
		CL.PrintOkf("任务ID %d 的任务名已更新为: %s", *editID, task.TaskName)
	}
	if *editRetentionCount != -1 {
		CL.PrintOkf("任务ID %d 的保留数量已更新为: %d", *editID, task.RetentionCount)
	}
	if *editRetentionDays != -1 {
		CL.PrintOkf("任务ID %d 的保留天数已更新为: %d", *editID, task.RetentionDays)
	}
	if *editNewDirName != "" {
		CL.PrintOkf("任务ID %d 的备份目录已更新为: %s", *editID, task.BackupDirectory)
	}
	if *editNoCompression != "" {
		if task.NoCompression == 1 {
			CL.PrintOkf("任务ID %d 的压缩状态已更新为: 禁用", *editID)
		} else {
			CL.PrintOkf("任务ID %d 的压缩状态已更新为: 启用", *editID)
		}
	}

	return nil
}
