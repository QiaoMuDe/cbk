package cmd

import (
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// delete命令的执行逻辑
func deleteCmdMain(db *sqlx.DB) error {
	// 存储任务ID的切片
	var ids []int

	// 如果指定了多个任务ID, 则执行多任务模式
	if *deleteIDS != "" {
		// 解析多个任务ID
		for _, idStr := range strings.Split(*deleteIDS, ",") {
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
		if err := deleteTasks(db, ids); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}

		return nil
	}

	// 如果指定了单个任务ID, 则执行单任务模式(单任务模式支持：删除任务，删除指定版本的备份)
	if *deleteID != 0 || *deleteName != "" {
		// 执行任务
		if err := deleteTask(db); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}

		return nil
	}

	// 检查必要的参数是否指定
	if *deleteID == 0 && *deleteIDS == "" && *deleteName == "" {
		return fmt.Errorf("删除备份任务时, 必须指定任务ID或者任务名, 使用-id指定任务ID或-n指定任务名称或用-ids指定多个任务ID, 例如: -ids '1,2,3'")
	}

	return nil
}

// 单ID模式删除任务
func deleteTask(db *sqlx.DB) error {
	// 检查是否指定任务ID和任务名
	if *deleteName == "" && *deleteID == 0 {
		return fmt.Errorf("必须指定要删除的任务, 请使用-id指定任务ID或-n指定任务名称")
	}

	// 如果版本ID不为空, 但是则检查是否指定了任务ID或任务名
	if *deleteVersionID != "" && (*deleteID == 0 && *deleteName == "") {
		return fmt.Errorf("删除指定版本的备份时, 必须指定任务ID或任务名, 使用-id指定任务ID或-n指定任务名称")
	}

	// 检查是否同时指定了任务ID和任务名
	if *deleteID != 0 && *deleteName != "" {
		return fmt.Errorf("不能同时使用-id和-n参数, 请选择其中一种方式指定任务")
	}

	// 根据任务名删除任务
	if *deleteName != "" && *deleteVersionID == "" {
		var backupDir string
		backupDirSql := "SELECT backup_directory FROM backup_tasks WHERE task_name = ?"
		if err := db.Get(&backupDir, backupDirSql, *deleteName); err == sql.ErrNoRows {
			return fmt.Errorf("任务名不存在: %s", *deleteName)
		} else if err != nil {
			return fmt.Errorf("获取备份存放目录失败: %w", err)
		}

		// 删除备份目录
		if err := deleteBackupDir(backupDir); err != nil {
			return err
		}

		// 删除任务和备份记录
		deleteSql := "DELETE FROM backup_tasks WHERE task_name = ?"
		if _, err := db.Exec(deleteSql, *deleteName); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}
		deleteBackupSql := "DELETE FROM backup_records WHERE task_name = ?"
		if _, err := db.Exec(deleteBackupSql, *deleteName); err != nil {
			return fmt.Errorf("删除备份记录失败: %w", err)
		}

		CL.PrintOkf("任务删除成功: %s", *deleteName)
		return nil
	}

	// 根据任务ID删除任务
	if *deleteID != 0 && *deleteVersionID == "" {
		var backupDir string
		backupDirSql := "SELECT backup_directory FROM backup_tasks WHERE task_id = ?"
		if err := db.Get(&backupDir, backupDirSql, *deleteID); err == sql.ErrNoRows {
			return fmt.Errorf("任务ID不存在: %d", *deleteID)
		} else if err != nil {
			return fmt.Errorf("获取备份存放目录失败: %w", err)
		}

		// 删除备份目录
		if err := deleteBackupDir(backupDir); err != nil {
			return err
		}

		// 删除任务和备份记录
		deleteSql := "DELETE FROM backup_tasks WHERE task_id = ?"
		if _, err := db.Exec(deleteSql, *deleteID); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}
		deleteBackupSql := "DELETE FROM backup_records WHERE task_id = ?"
		if _, err := db.Exec(deleteBackupSql, *deleteID); err != nil {
			return fmt.Errorf("删除备份记录失败: %w", err)
		}

		CL.PrintOkf("任务ID删除成功: %d", *deleteID)
		return nil
	}

	// 根据任务ID和版本ID删除备份记录
	if *deleteID != 0 && *deleteVersionID != "" {
		// 获取根据任务ID和版本ID查询备份记录
		var backupRecord struct {
			BackupPath string `db:"backup_path"`      // 备份目录
			BackupFile string `db:"backup_file_name"` // 备份文件
		}
		backupRecordSql := "select backup_path, backup_file_name from backup_records where task_id =? and version_id =?"
		if err := db.Get(&backupRecord, backupRecordSql, *deleteID, *deleteVersionID); err == sql.ErrNoRows {
			return fmt.Errorf("任务ID或版本ID不存在")
		} else if err != nil {
			return fmt.Errorf("查询备份记录失败: %w", err)
		}

		// 切换到备份目录
		if err := os.Chdir(backupRecord.BackupPath); err != nil {
			return fmt.Errorf("切换到备份目录失败: %w", err)
		}

		// 删除备份文件
		if _, err := tools.CheckPath(backupRecord.BackupFile); err == nil {
			if err := os.Remove(backupRecord.BackupFile); err != nil {
				return fmt.Errorf("删除备份文件失败: %w", err)
			}
		} else {
			CL.PrintWarnf("备份文件不存在: %s", backupRecord.BackupFile)
		}

		// 删除备份记录
		deleteBackupSql := "delete from backup_records where task_id = ? and version_id = ?"
		if _, err := db.Exec(deleteBackupSql, *deleteID, *deleteVersionID); err != nil {
			return fmt.Errorf("删除备份记录失败: %w", err)
		}

		// 打印成功信息
		CL.PrintOkf("任务ID: %d, 版本ID: %s 删除成功", *deleteID, *deleteVersionID)

		return nil
	}

	return nil
}

// 多ID模式删除任务
func deleteTasks(db *sqlx.DB, ids []int) error {
	// 如果版本ID不为空, 同时任务ID列表也不为空, 则返回错误
	if *deleteVersionID != "" && *deleteIDS != "" {
		return fmt.Errorf("-ids 不支持删除指定版本的备份, 请使用 -id 指定任务ID和版本ID")
	}

	// 检查是否没有指定任务ID列表
	if *deleteIDS == "" {
		return fmt.Errorf("必须指定要删除的任务, 请使用-ids指定任务ID列表, 例如: -ids '1,2,3'")
	}

	// 根据任务ID删除任务
	if *deleteIDS != "" && *deleteVersionID == "" {
		var backupDir string // 备份目录
		// 构建查询备份目录的SQL语句
		backupDirSql := "SELECT backup_directory FROM backup_tasks WHERE task_id = ?"

		for _, id := range ids {
			if err := db.Get(&backupDir, backupDirSql, id); err == sql.ErrNoRows {
				CL.PrintErrf("任务ID不存在: %d", id)
				continue
			} else if err != nil {
				CL.PrintErrf("获取备份存放目录失败: %v", err)
				continue
			}

			// 删除备份目录
			if err := deleteBackupDir(backupDir); err != nil {
				CL.PrintErrf("删除备份存放目录失败: %v", err)
				continue
			}

			// 删除任务和备份记录
			deleteSql := "DELETE FROM backup_tasks WHERE task_id = ?"
			if _, err := db.Exec(deleteSql, id); err != nil {
				CL.PrintErrf("删除任务失败: %v", err)
				continue
			}
			deleteBackupSql := "DELETE FROM backup_records WHERE task_id = ?"
			if _, err := db.Exec(deleteBackupSql, id); err != nil {
				CL.PrintErrf("删除备份记录失败: %v", err)
				continue
			}

			CL.PrintOkf("任务ID删除成功: %d", id)
		}
	}

	return nil
}

// deleteBackupDir 删除备份目录
// 参数:
// - backupDir: 备份目录
// 返回值:
// - error: 错误信息
func deleteBackupDir(backupDir string) error {
	// 检查是否设置了删除目录标志(*deleteDirF)
	if *deleteDirF {
		// 检查备份目录是否存在
		if _, err := tools.CheckPath(backupDir); err == nil {
			// 如果存在则删除整个目录
			if err := os.RemoveAll(backupDir); err != nil {
				return fmt.Errorf("删除备份存放目录失败: %w", err)
			}
			// 删除成功提示
			CL.PrintOkf("备份存放目录删除成功: %s", backupDir)
		}
	} else {
		// 没有设置删除标志时的提示
		CL.PrintWarnf("注意: 备份目录 %s 未被删除，请在删除任务后手动清理，或下次使用 -d 参数自动删除", backupDir)
	}
	return nil
}
