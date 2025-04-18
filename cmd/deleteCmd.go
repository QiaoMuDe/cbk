package cmd

import (
	"cbk/pkg/tools"
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

// delete命令的执行逻辑
func deleteCmdMain(db *sqlx.DB) error {
	// 如果版本ID不为空, 但是任务ID为0, 则返回错误
	if *deleteVersionID != "" && *deleteID == 0 {
		return fmt.Errorf("指定了版本ID, 但是未指定任务ID")
	}

	// 检查是否同时指定了任务ID和任务名
	if *deleteID != 0 && *deleteName != "" {
		return fmt.Errorf("不能同时指定任务ID和任务名")
	}

	// 检查是否没有指定任务ID和任务名
	if *deleteName == "" && *deleteID == 0 {
		return fmt.Errorf("删除任务时, 必须指定任务名或任务ID")
	}

	// 共用的删除备份目录逻辑
	deleteBackupDir := func(backupDir string) error {
		if *deleteDirF {
			if _, err := tools.CheckPath(backupDir); err == nil {
				if err := os.RemoveAll(backupDir); err != nil {
					return fmt.Errorf("删除备份存放目录失败: %w", err)
				}
				CL.PrintOkf("备份存放目录删除成功: %s", backupDir)
			} else {
				CL.PrintWarnf("请在稍后，手动删除备份存放目录: %s", backupDir)
			}
		}
		return nil
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
