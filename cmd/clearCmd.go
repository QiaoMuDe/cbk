package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

// clearCmdMain 清除数据主逻辑
func clearCmdMain(db *sqlx.DB) error {
	// 检查是否确认
	if !*clearConfirm {
		return fmt.Errorf("请使用 -confirm 参数确认清除操作")
	}

	CL.PrintWarn("即将清空整个数据库和备份存放目录，这将删除所有备份任务和相关数据, 撤回可在三秒内按Ctrl+C退出")
	time.Sleep(3 * time.Second) // 等待3秒

	// 构建查询备份任务的SQL语句
	querySql := "SELECT backup_directory FROM backup_tasks;"
	// 定义存储查询结果的结构体
	var tasks []struct {
		BackupDirectory string `db:"backup_directory"` // 备份目录
	}

	// 执行查询备份任务的SQL语句
	if err := db.Select(&tasks, querySql); err != nil {
		return fmt.Errorf("查询备份任务失败: %w", err)
	}

	// 检查是否存在备份任务
	if len(tasks) == 0 {
		return fmt.Errorf("未找到任何备份任务")
	}

	// 构建清除备份任务的SQL语句
	clearBackupTasksSql := "delete from backup_tasks;"

	// 构建清空备份记录的SQL语句
	clearBackupRecordsSql := "delete from backup_records;"

	// 执行清除备份任务的SQL语句
	if _, err := db.Exec(clearBackupTasksSql); err != nil {
		return fmt.Errorf("清除备份任务失败: %w", err)
	}

	// 执行清空备份记录的SQL语句
	if _, err := db.Exec(clearBackupRecordsSql); err != nil {
		return fmt.Errorf("清空备份记录失败: %w", err)
	}

	// 遍历清理备份存放目录
	for _, task := range tasks {
		if err := os.RemoveAll(task.BackupDirectory); err != nil {
			CL.PrintErrorf("清理备份存放目录失败: %s", task.BackupDirectory)
			CL.PrintWarnf("请手动删除备份存放目录: %s", task.BackupDirectory)
			continue
		}
		CL.PrintOkf("清理备份存放目录成功: %s", task.BackupDirectory)
	}

	return nil
}
