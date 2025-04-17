package cmd

import (
	"cbk/pkg/tools"
	"database/sql"
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
	AbsAddTarget, err := filepath.Abs(filepath.Clean(*addTarget))
	if err != nil {
		return fmt.Errorf("获取目标目录绝对路径失败: %w", err)
	}

	// 如果备份目录名为空, 则获取目标目录的basename作为存放备份的目录名
	var bakDirName string
	if *addBackupDirName == "" {
		bakDirName = filepath.Base(AbsAddTarget)
	} else {
		// 使用指定的备份目录名
		bakDirName = *addBackupDirName
	}

	// 如果备份目录为空, 则使用默认值路径，格式为: /home/username/.cbk/data/xxx
	if *addBackup == "" {
		// 获取用户主目录
		tempHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %w", err)
		}
		// 构建备份目录的绝对路径, 格式为: /home/username/.cbk/data/bakDirName
		*addBackup = filepath.Join(tempHome, ".cbk", "data", bakDirName)

		// 检查备份目录是否存在
		if _, err := tools.CheckPath(*addBackup); err != nil {
			if err := os.MkdirAll(*addBackup, 0755); err != nil {
				return fmt.Errorf("备份目录创建失败: %w", err)
			}
		}
	} else {
		// 检查备份目录是否存在
		if _, err := tools.CheckPath(*addBackup); err != nil {
			return fmt.Errorf("备份目录不存在: %w", err)
		}

		// 检查备份目录是否为绝对路径
		if !filepath.IsAbs(*addBackup) {
			var err error
			*addBackup, err = filepath.Abs(*addBackup)
			if err != nil {
				return fmt.Errorf("获取备份目录绝对路径失败: %w", err)
			}
		}

		// 构建自定义备份目录的绝对路径, 格式为: /path/to/bakDirName
		*addBackup = filepath.Join(*addBackup, bakDirName)

		// 检查备份目录是否存在
		if _, err := tools.CheckPath(*addBackup); err != nil {
			if err := os.MkdirAll(*addBackup, 0755); err != nil {
				return fmt.Errorf("备份目录创建失败: %w", err)
			}
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
	insertSql := "insert into backup_tasks(task_name, target_directory, backup_directory, retention_count, no_compression) values(?, ?, ?, ?, ?)"
	if _, err := db.Exec(insertSql, *addName, AbsAddTarget, *addBackup, *addKeep, noCompression); err != nil {
		return fmt.Errorf("插入任务失败: %w", err)
	}

	// 打印成功信息
	CL.PrintOkf("任务添加成功: %s", *addName)
	return nil
}

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
