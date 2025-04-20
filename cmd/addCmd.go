package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/tools"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
)

// add命令的执行逻辑
func addCmdMain(db *sqlx.DB) error {
	// 检查是否指定-f参数
	if *addConfig != "" {
		// 检查配置文件是否存在
		if _, err := tools.CheckPath(*addConfig); err != nil {
			return fmt.Errorf("指定的配置文件不存在: %w", err)
		}

		// 读取配置文件
		config, err := os.ReadFile(*addConfig)
		if err != nil {
			return fmt.Errorf("读取 %s 配置文件失败: %w", *addConfig, err)
		}

		// 解析配置文件
		var addTaskConfig globals.TaskConfig
		if err := yaml.Unmarshal(config, &addTaskConfig); err != nil {
			return fmt.Errorf("解析 %s 配置文件失败: %w", *addConfig, err)
		}

		// 添加任务
		if err := addTask(db, addTaskConfig.Task.Name, addTaskConfig.Task.Target, addTaskConfig.Task.Backup, addTaskConfig.Task.BackupDirName, addTaskConfig.Task.Retention.Count, addTaskConfig.Task.Retention.Days, addTaskConfig.Task.NoCompression, addTaskConfig.Task.ExcludeRules); err != nil {
			return fmt.Errorf("添加任务失败: %w", err)
		}

		return nil
	}

	// 如果没有指定-f参数, 则执行普通添加任务模式
	if err := addTask(db, *addName, *addTarget, *addBackup, *addBackupDirName, *addRetentionCount, *addRetentionDays, *addNoCompression, *addExcludeRules); err != nil {
		return fmt.Errorf("添加任务失败: %w", err)
	}
	return nil
}

// addTask 添加备份任务
// 参数:
// - db: 数据库连接
// - taskName: 任务名
// - targetDir: 目标目录
// - backupDir: 备份目录
// - retentionCount: 保留文件数量
// - retentionDays: 保留天数
// - noCompression: 是否禁用压缩(默认启用压缩, 0 表示启用压缩, 1 表示禁用压缩)
// - excludeRules: 排除规则
// 返回值:
// - error: 错误信息
func addTask(db *sqlx.DB, taskName string, targetDir string, backupDir string, backupDirName string, retentionCount int, retentionDays int, noCompression int, excludeRules string) error {
	// 检查任务名是否为空
	if taskName == "" {
		return fmt.Errorf("任务名不能为空")
	}

	// 检查任务名是否非法字符
	if tools.ContainsSpecialChars(taskName) {
		return fmt.Errorf("任务名含非法字符, 请重试")
	}

	// 检查目标目录是否为空
	if targetDir == "" {
		return fmt.Errorf("目标目录不能为空")
	}

	// 检查备份目录名是否非法字符
	if backupDirName != "" {
		if tools.ContainsSpecialChars(backupDirName) {
			return fmt.Errorf("备份目录名含非法字符, 请重试")
		}
	}

	// 检查保留文件数量是否合法
	if retentionCount <= 0 {
		return fmt.Errorf("保留文件数量不能小于0")
	}

	// 检查保留天数是否合法
	if retentionDays < 0 {
		return fmt.Errorf("保留天数不能小于0")
	}

	// 检查目标目录或文件是否存在
	if _, err := tools.CheckPath(targetDir); err != nil {
		return fmt.Errorf("目标目录或文件不存在: %w", err)
	}

	// 如果指定了禁用压缩, 则检查是否合法
	if *addNoCompression != 1 && *addNoCompression != 0 {
		return fmt.Errorf("-nc 参数不合法, 只能是 0(启用压缩) 或 1(禁用压缩)")
	}

	// 在数据库检查是否存在同名任务
	checkSql := "select count(*) from backup_tasks where task_name = ?"
	var count int
	if err := db.Get(&count, checkSql, taskName); err != nil {
		return fmt.Errorf("检查同名任务失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("任务名已存在, 请在更换任务名或删除已有任务后再添加")
	}

	// 扩展目标目录为绝对路径
	absTargetDir, err := filepath.Abs(filepath.Clean(targetDir))
	if err != nil {
		return fmt.Errorf("获取目标目录绝对路径失败: %w", err)
	}

	// 如果备份目录名为空, 则获取目标目录的basename作为存放备份的目录名
	if backupDirName == "" {
		backupDirName = filepath.Base(absTargetDir)
	}

	// 如果备份目录为空, 则使用默认值路径，格式为: /home/username/.cbk/data/xxx
	var absBackupDir string // 定义备份目录的绝对路径
	if backupDir == "" {
		// 获取用户主目录
		tempHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %w", err)
		}
		// 构建备份目录的绝对路径, 格式为: /home/username/.cbk/data/backupDirName
		absBackupDir = filepath.Join(tempHome, globals.CbkHomeDir, globals.CbkDataDir, backupDirName)

		// 检查备份目录是否存在
		if err := tools.EnsureDirExists(absBackupDir); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}
	} else {
		// 检查指定的备份目录是否存在并创建
		if err := tools.EnsureDirExists(backupDir); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}

		// 检查备份目录是否为绝对路径, 如果不是, 则转换为绝对路径
		if !filepath.IsAbs(backupDir) {
			var err error
			backupDir, err = filepath.Abs(backupDir)
			if err != nil {
				return fmt.Errorf("获取备份目录绝对路径失败: %w", err)
			}
		}

		// 构建自定义备份目录的绝对路径, 格式为: /path/to/backupDirName
		absBackupDir = filepath.Join(backupDir, backupDirName)

		// 检查备份目录是否存在并创建
		if err := tools.EnsureDirExists(absBackupDir); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}
	}

	// 插入新任务到数据库
	insertSql := "insert into backup_tasks(task_name, target_directory, backup_directory, retention_count, retention_days, no_compression, exclude_rules) values(?, ?, ?, ?, ?, ?, ?)"
	if _, err := db.Exec(insertSql, taskName, absTargetDir, absBackupDir, retentionCount, retentionDays, noCompression, excludeRules); err != nil {
		return fmt.Errorf("插入任务失败: %w", err)
	}

	// 打印成功信息
	CL.PrintOkf("任务添加成功: %s", taskName)
	return nil
}
