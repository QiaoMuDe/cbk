package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// editCmdMain 编辑任务
func editCmdMain(db *sqlx.DB) error {
	// 存储任务ID的切片
	var ids []int

	// 如果指定了多个任务ID, 则执行多任务模式
	if *editIDS != "" {
		// 解析多个任务ID
		for _, idStr := range strings.Split(*editIDS, ",") {
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

		// 编辑任务
		if err := editTask(db, ids); err != nil {
			return fmt.Errorf("编辑任务失败: %w", err)
		}

		return nil
	}

	// 如果指定了单个任务ID, 则执行单任务模式
	if *editID != -1 {
		// 添加单个任务ID到切片中
		ids = append(ids, *editID)

		// 执行任务
		if err := editTask(db, ids); err != nil {
			return fmt.Errorf("编辑任务失败: %w", err)
		}

		return nil
	}

	// 检查任务ID是否指定
	if *editID == -1 || *editIDS == "" {
		return fmt.Errorf("编辑备份任务时, 必须指定任务ID或任务ID列表, 使用 -id 或 -ids 指定, 例如: -id 1 或 -ids '1,2,3'")
	}

	return nil
}

// editTask 编辑任务
// 参数:
// - db: 数据库连接
// - ids: 任务ID切片
// 返回值:
// - error: 错误信息
func editTask(db *sqlx.DB, ids []int) error {
	// 构建存储查询任务信息的结构体
	var task globals.BackupTask

	// 查询任务信息
	editSql := "select task_name, retention_count, retention_days, backup_directory, no_compression, exclude_rules from backup_tasks where task_id =?"

	// 更新任务
	updateSql := "update backup_tasks set task_name = ?, retention_count = ? , retention_days = ?, backup_directory = ?, no_compression = ?, exclude_rules = ? where task_id = ?"

	for _, id := range ids {
		// 检查所有的参数是否都没指定
		if *editName == "" && *editRetentionCount == -1 && *editRetentionDays == -1 && *editNoCompression == -1 && *editNewDirName == "" && *editExcludeRules == "" {
			CL.PrintWarnf("在编辑 %d 时未指定任何参数, 该任务将不会被修改", id)
			continue
		}

		// 检查任务ID是否存在
		if err := db.Get(&task, editSql, id); err == sql.ErrNoRows {
			CL.PrintErrf("任务ID不存在 %d", id)
			continue
		} else if err != nil {
			CL.PrintErrf("查询任务失败: %v, SQL: %s, ID: %d", err, editSql, id)
			continue
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
		if *editNoCompression != -1 {
			// 检查如果不是true或false则报错
			if *editNoCompression != 1 && *editNoCompression != 0 {
				CL.PrintErrf("-nc 参数不合法, 只能是 0(启用压缩) 或 1(禁用压缩)")
				continue
			}

			// 根据参数值更新NoCompression字段
			task.NoCompression = *editNoCompression
		}

		// 如果指定了-bn参数, 则更新备份目录
		var oldDirName, rootPath, newDirName string
		if *editNewDirName != "" {
			newDirName = *editNewDirName

			// 检查备份目录名是否非法字符
			if tools.ContainsSpecialChars(newDirName) {
				CL.PrintErrf("备份目录名 [%s] 含非法字符, 请重试", newDirName)
				continue
			}

			rootPath = filepath.Dir(task.BackupDirectory)    // 获取备份目录的根路径
			oldDirName = filepath.Base(task.BackupDirectory) // 获取备份目录的旧名称

			// 重命名备份目录
			if err := tools.RenameBackupDirectory(rootPath, oldDirName, newDirName); err != nil {
				CL.PrintErrf("重命名备份目录失败: %v", err)
				continue
			}

			// 更新备份目录路径
			task.BackupDirectory = filepath.Join(rootPath, newDirName)
		}

		// 如果指定了-ex参数, 则更新排除规则
		if *editExcludeRules != "" {
			task.ExcludeRules = *editExcludeRules
		}

		// 更新任务SQL
		if _, err := db.Exec(updateSql, task.TaskName, task.RetentionCount, task.RetentionDays, task.BackupDirectory, task.NoCompression, task.ExcludeRules, id); err != nil {
			// 更新任务失败
			if *editNewDirName != "" {
				if err := tools.RenameBackupDirectory(rootPath, newDirName, oldDirName); err != nil {
					CL.PrintErrf("更新任务失败且恢复备份目录失败: %v", err)
					continue
				}
				CL.PrintOkf("更新任务失败, 已恢复备份目录: %s", filepath.Join(rootPath, oldDirName))
			}
			CL.PrintErrf("更新任务失败: %v, SQL: %s, ID: %d", err, updateSql, id)
			continue
		}

		// 打印成功信息
		CL.PrintOk("更新成功!")
		if *editName != "" {
			CL.PrintOkf("任务ID %d 的任务名已更新为: %s", id, task.TaskName)
		}
		if *editRetentionCount != -1 {
			CL.PrintOkf("任务ID %d 的保留数量已更新为: %d", id, task.RetentionCount)
		}
		if *editRetentionDays != -1 {
			CL.PrintOkf("任务ID %d 的保留天数已更新为: %d", id, task.RetentionDays)
		}
		if *editNewDirName != "" {
			CL.PrintOkf("任务ID %d 的备份目录已更新为: %s", id, task.BackupDirectory)
		}
		if *editNoCompression != -1 {
			if task.NoCompression == 1 {
				CL.PrintOkf("任务ID %d 的压缩状态已更新为: 禁用", id)
			} else {
				CL.PrintOkf("任务ID %d 的压缩状态已更新为: 启用", id)
			}
		}
		if *editExcludeRules != "" {
			CL.PrintOkf("任务ID %d 的排除规则已更新为: %s", id, task.ExcludeRules)
		}
	}

	return nil
}
