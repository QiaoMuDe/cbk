package cmd

import (
	"cbk/pkg/globals"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
	var tasks globals.BackupTasks

	// 执行查询备份任务的SQL语句
	if err := db.Select(&tasks, querySql); err != nil {
		return fmt.Errorf("查询备份任务失败: %w", err)
	}

	// 检查是否存在备份任务, 如果存在则清理备份存放目录
	if len(tasks) > 0 {
		// 遍历清理备份存放目录
		for _, task := range tasks {
			if err := os.RemoveAll(task.BackupDirectory); err != nil {
				CL.PrintErrorf("清理备份存放目录失败: %s", task.BackupDirectory)
				CL.PrintWarnf("请手动删除备份存放目录: %s", task.BackupDirectory)
				continue
			}
			CL.PrintOkf("清理备份存放目录成功: %s", task.BackupDirectory)
		}
	}

	// 构建数据库文件路径
	var useHomeDir string
	useHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %w", err)
	}
	dbPath := filepath.Join(useHomeDir, globals.CbkDbPath)

	// 检查数据库连接是否打开，如果打开则关闭
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %w", err)
		}
	}

	// 清理数据库文件
	if err := os.Remove(dbPath); err != nil {
		return fmt.Errorf("清理数据库文件失败: %w", err)
	}

	return nil
}
