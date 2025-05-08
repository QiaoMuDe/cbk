package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// runCmdMain 运行备份任务
func runCmdMain(db *sqlx.DB) error {
	// 存储任务ID的切片
	var ids []int

	// 如果指定了多个任务ID, 则执行多任务模式
	if *runIDS != "" {
		// 解析多个任务ID
		for _, idStr := range strings.Split(*runIDS, ",") {
			// 检查解析的任务ID是否为空
			if idStr == "" {
				CL.PrintErr("任务ID不能为空")
				continue
			}

			// 检查解析的任务ID是否包含特殊字符
			if tools.ContainsSpecialChars(idStr) {
				CL.PrintErrf("任务ID包含危险字符: %s", idStr)
				continue
			}

			// 将字符串转换为整数
			id, err := strconv.Atoi(idStr)
			if err != nil {
				CL.PrintErrf("无效的任务ID: %s", idStr)
				continue
			}

			// 添加任务ID到切片中
			ids = append(ids, id)
		}

		// 执行任务
		if err := runTask(db, ids); err != nil {
			return fmt.Errorf("运行任务失败: %w", err)
		}

		return nil
	}

	// 如果指定了单个任务ID, 则执行单任务模式
	if *runID != 0 {
		// 添加单个任务ID到切片中
		ids = append(ids, *runID)

		// 执行任务
		if err := runTask(db, ids); err != nil {
			return fmt.Errorf("运行任务失败: %w", err)
		}

		return nil
	}

	// 检查任务ID是否指定
	if *runID == 0 || *runIDS == "" {
		return fmt.Errorf("运行备份任务时, 必须指定任务ID或任务ID列表, 使用 -id 或 -ids 指定, 例如: -id 1 或 -ids '1,2,3'")
	}

	return nil
}

// runTask 执行备份任务
// 参数:
// - db: 数据库连接
// - ids: 任务ID切片
// 返回值:
// - error: 错误信息
func runTask(db *sqlx.DB, ids []int) error {
	// 构建存储查询任务信息的结构体
	var task globals.BackupTask

	// 构建查询任务信息的SQL语句
	querySql := "select task_name, target_directory, backup_directory, retention_count, retention_days, no_compression, exclude_rules from backup_tasks where task_id =?"

	// 构建失败记录的SQL语句
	errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// 构建插入备份记录的SQL语句
	insertSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// 循环处理每个任务ID
	for _, id := range ids {
		// 查询任务信息
		if err := db.Get(&task, querySql, id); err == sql.ErrNoRows {
			CL.PrintErrf("任务ID不存在 %d", id)
			continue
		} else if err != nil {
			CL.PrintErrf("获取任务信息失败: %v", err)
			continue
		}

		// 检查目标目录或文件是否存在
		if _, err := tools.CheckPath(task.TargetDirectory); err != nil {
			CL.PrintErrf("目标目录或文件不存在: %v", err)
			continue
		}

		// 检查备份目录是否存在
		if err := tools.EnsureDirExists(task.BackupDirectory); err != nil {
			CL.PrintErrf("备份目录创建失败: %v", err)
			continue
		}

		// 打印提示信息
		CL.PrintOkf("备份任务 [%s] 已启动，正在运行中……", task.TaskName)

		// 构建备份文件名
		backupTime := time.Now().Format("20060102150405")
		backupFileNamePrefix := fmt.Sprintf("%s_%s", task.TaskName, backupTime)

		// 获取排除函数
		var excludeFunc globals.ExcludeFunc
		if task.ExcludeRules != "none" {
			var err error
			if excludeFunc, err = tools.ParseExclude(task.ExcludeRules); err != nil {
				CL.PrintErrf("解析任务ID %d 的排除规则失败: %v", id, err)
				continue
			}
		} else {
			excludeFunc = globals.NoExcludeFunc // 默认不进行过滤
		}

		// 获取versionID
		versionID := tools.GenerateID(6)

		// 运行备份任务
		targetDir := filepath.Dir(task.TargetDirectory)                                 // 获取目标目录的目录部分
		targetName := filepath.Base(task.TargetDirectory)                               // 获取目标目录的最后一个部分
		backupFileNamePath := filepath.Join(task.BackupDirectory, backupFileNamePrefix) // 获取构建的备份文件路径

		// 执行备份任务
		zipPath, err := tools.CreateZipFromOSPaths(db, targetDir, targetName, backupFileNamePath, task.NoCompression, excludeFunc)
		if err != nil {
			// 插入备份记录
			if _, execErr := db.Exec(errorSql, versionID, id, backupTime, task.TaskName, "false", "-", "-", "-", "-"); execErr != nil {
				CL.PrintErrf("插入备份记录失败: %v", execErr)
				continue
			}
			CL.PrintErrf("备份 %s 任务失败: %v", task.TaskName, err)
			continue
		}

		// 获取备份文件的后8位MD5哈希值
		backupFileMD5, err := tools.GetFileMD5Last8(zipPath)
		if err != nil {
			// 插入备份记录
			if _, execErr := db.Exec(errorSql, versionID, id, backupTime, task.TaskName, "false", "-", "-", "-", "-"); execErr != nil {
				CL.PrintErrf("插入备份记录失败: %v", execErr)
				continue
			}
			CL.PrintErrf("获取备份文件MD5失败: %v", err)
			continue
		}

		// 获取备份文件的大小
		backupFileSize, err := tools.HumanReadableSize(zipPath)
		if err != nil {
			// 插入备份记录
			if _, execErr := db.Exec(errorSql, versionID, id, backupTime, task.TaskName, "false", "-", "-", "-", "-"); execErr != nil {
				CL.PrintErrf("插入备份记录失败: %v", execErr)
				continue
			}
			CL.PrintErrf("获取备份文件大小失败: %v", err)
			continue
		}

		// 插入备份记录
		if _, execErr := db.Exec(insertSql, versionID, id, backupTime, task.TaskName, "true", filepath.Base(zipPath), backupFileSize, task.BackupDirectory, backupFileMD5); execErr != nil {
			CL.PrintErrf("插入备份记录失败: %v", execErr)
			continue
		}

		// 获取备份目录下的以指定扩展名的文件列表
		zipFiles, err := tools.GetZipFiles(task.BackupDirectory, ".zip")
		if err != nil {
			CL.PrintErrf("获取备份目录下的.zip文件失败: %v", err)
			continue
		}

		// 删除多余的备份文件
		if len(zipFiles) > task.RetentionCount {
			if err := tools.RetainLatestFiles(db, zipFiles, task.RetentionCount, task.RetentionDays); err != nil {
				CL.PrintErrf("删除多余的备份文件失败: %v", err)
				continue
			}
		}

		// 打印成功信息
		CL.PrintOkf(`备份 %s 成功!`, task.TaskName)
	}

	return nil
}
