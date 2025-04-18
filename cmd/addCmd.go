package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/tools"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

// add命令的执行逻辑
func addCmdMain(db *sqlx.DB) error {
	// 检查任务名是否为空
	if *addName == "" {
		return fmt.Errorf("任务名不能为空")
	}

	// 检查任务名是否非法字符
	if tools.ContainsSpecialChars(*addName) {
		return fmt.Errorf("任务名含非法字符, 请重试")
	}

	// 检查目标目录是否为空
	if *addTarget == "" {
		return fmt.Errorf("目标目录不能为空")
	}

	// 检查备份目录名是否非法字符
	if *addBackupDirName != "" {
		if tools.ContainsSpecialChars(*addBackupDirName) {
			return fmt.Errorf("备份目录名含非法字符, 请重试")
		}
	}

	// 检查目标目录或文件是否存在
	if _, err := tools.CheckPath(*addTarget); err != nil {
		return fmt.Errorf("目标目录或文件不存在: %w", err)
	}

	// 在数据库检查是否存在同名任务
	checkSql := "select count(*) from backup_tasks where task_name = ?"
	var count int
	if err := db.Get(&count, checkSql, *addName); err != nil {
		return fmt.Errorf("检查同名任务失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("任务名已存在, 请在更换任务名或删除已有任务后再添加")
	}

	// 扩展目标目录为绝对路径
	absTargetDir, err := filepath.Abs(filepath.Clean(*addTarget))
	if err != nil {
		return fmt.Errorf("获取目标目录绝对路径失败: %w", err)
	}

	// 如果备份目录名为空, 则获取目标目录的basename作为存放备份的目录名
	var bakDirName string
	if *addBackupDirName == "" {
		bakDirName = filepath.Base(absTargetDir)
	} else {
		// 使用指定的备份目录名
		bakDirName = *addBackupDirName
	}

	// 如果备份目录为空, 则使用默认值路径，格式为: /home/username/.cbk/data/xxx
	var absBackupDir string // 定义备份目录的绝对路径
	if *addBackup == "" {
		// 获取用户主目录
		tempHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %w", err)
		}
		// 构建备份目录的绝对路径, 格式为: /home/username/.cbk/data/bakDirName
		absBackupDir = filepath.Join(tempHome, globals.CbkHomeDir, globals.CbkDataDir, bakDirName)

		// 检查备份目录是否存在
		if err := tools.EnsureDirExists(absBackupDir); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}
	} else {
		// 检查指定的备份目录是否存在并创建
		if err := tools.EnsureDirExists(*addBackup); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}

		// 检查备份目录是否为绝对路径, 如果不是, 则转换为绝对路径
		if !filepath.IsAbs(*addBackup) {
			var err error
			*addBackup, err = filepath.Abs(*addBackup)
			if err != nil {
				return fmt.Errorf("获取备份目录绝对路径失败: %w", err)
			}
		}

		// 构建自定义备份目录的绝对路径, 格式为: /path/to/bakDirName
		absBackupDir = filepath.Join(*addBackup, bakDirName)

		// 检查备份目录是否存在并创建
		if err := tools.EnsureDirExists(absBackupDir); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}
	}

	// 获取是否禁用压缩
	var noCompression int
	if *addNoCompression {
		noCompression = 1 // 启用压缩
	} else {
		noCompression = 0 // 禁用压缩(默认启用压缩)
	}

	// 插入新任务到数据库
	insertSql := "insert into backup_tasks(task_name, target_directory, backup_directory, retention_count, retention_days, no_compression) values(?, ?, ?, ?, ?, ?)"
	if _, err := db.Exec(insertSql, *addName, absTargetDir, absBackupDir, *addRetentionCount, *addRetentionDays, noCompression); err != nil {
		return fmt.Errorf("插入任务失败: %w", err)
	}

	// 打印成功信息
	CL.PrintOkf("任务添加成功: %s", *addName)
	return nil
}
