package cmd

import (
	"cbk/pkg/globals"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// exportCmdMain 导出备份任务的主函数
// 参数:
//   - db: *sqlx.DB, 数据库连接对象
//
// 返回值:
//   - error, 错误信息
func exportCmdMain(db *sqlx.DB) error {
	// 构建查询所有的备份任务的SQL语句
	queryAllSql := "SELECT task_name, target_directory, backup_directory, retention_count, retention_days, no_compression, exclude_rules FROM backup_tasks;"

	// 构建查询单个备份任务的SQL语句
	queryOneSql := "SELECT task_name, target_directory, backup_directory, retention_count, retention_days, no_compression, exclude_rules FROM backup_tasks WHERE task_id = ?;"

	// 定义存储查询结果的结构体切片
	var tasks globals.BackupTasks

	// 定义查询单个备份任务的结构体
	var task globals.BackupTask

	// 定义打印备份任务的cbk命令格式
	printCmd := "cbk add -n %s -bn %s -t %s -b %s -c %d -d %d -nc %d -ex %s\n"

	// 导出所有任务
	if *exportAll {
		// 执行查询所有备份任务的SQL语句
		if err := db.Select(&tasks, queryAllSql); err != nil {
			return fmt.Errorf("查询所有备份任务失败: %w", err)
		}

		// 遍历打印所有备份任务的cbk命令格式
		for _, task := range tasks {
			// 获取备份目录的名称
			bakDirName := filepath.Base(task.BackupDirectory)

			// 获取备份目录的父级目录
			parentDir := filepath.Dir(task.BackupDirectory)

			fmt.Printf(printCmd, task.TaskName, bakDirName, task.TargetDirectory, parentDir, task.RetentionCount, task.RetentionDays, task.NoCompression, task.ExcludeRules)
		}

		return nil
	}

	// 导出单个任务
	if *exportID != 0 {
		// 执行查询单个备份任务的SQL语句
		if err := db.Get(&task, queryOneSql, *exportID); err != nil {
			return fmt.Errorf("查询单个备份任务失败: %w", err)
		}

		// 打印单个备份任务的cbk命令格式
		bakDirName := filepath.Base(task.BackupDirectory) // 获取备份目录的名称
		parentDir := filepath.Dir(task.BackupDirectory)   // 获取备份目录的父级目录
		fmt.Printf(printCmd, task.TaskName, bakDirName, task.TargetDirectory, parentDir, task.RetentionCount, task.RetentionDays, task.NoCompression, task.ExcludeRules)

		return nil
	}

	return fmt.Errorf("请使用 -all 或 -id 参数指定导出的备份任务")
}
