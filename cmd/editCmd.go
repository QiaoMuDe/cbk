package cmd

import (
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

// editCmdMain 编辑任务
func editCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *editID == 0 {
		return fmt.Errorf("编辑任务时, 必须指定任务ID")
	}

	// 查询任务信息
	editSql := "select task_name, retention_count,backup_directory from backup_tasks where task_id =?"
	var task struct {
		TaskName        string `db:"task_name"`        // 任务名
		RetentionCount  int    `db:"retention_count"`  // 保留数量
		BackupDirectory string `db:"backup_directory"` // 备份目录
	}
	if err := db.Get(&task, editSql, *editID); err == sql.ErrNoRows {
		return fmt.Errorf("任务ID不存在 %d", *editID)
	} else if err != nil {
		return fmt.Errorf("查询任务失败: %w, SQL: %s, ID: %d", err, editSql, *editID)
	}

	// 如果指定了-n参数, 则更新任务名
	if *editName != "" {
		task.TaskName = *editName
	}

	// 如果指定了-k参数, 则更新保留数量
	if *editKeep != 0 {
		task.RetentionCount = *editKeep
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

	// 更新任务事务
	updateSql := "update backup_tasks set task_name = ?, retention_count = ? , backup_directory = ? where task_id = ?"

	// 更新任务SQL
	if _, err := db.Exec(updateSql, task.TaskName, task.RetentionCount, task.BackupDirectory, *editID); err != nil {
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
	if *editKeep != 0 {
		CL.PrintOkf("任务ID %d 的保留数量已更新为: %d", *editID, task.RetentionCount)
	}
	if *editNewDirName != "" {
		CL.PrintOkf("任务ID %d 的备份目录已更新为: %s", *editID, task.BackupDirectory)
	}

	return nil
}
