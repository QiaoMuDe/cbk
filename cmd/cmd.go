// cmd.go
package cmd

import (
	"cbk/pkg/tools"
	"cbk/pkg/version"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gitee.com/MM-Q/colorlib"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jmoiron/sqlx"
)

// 定义全局颜色渲染器
var CL = colorlib.NewColorLib()

//go:embed help/help.txt
var HelpText string

// 定义子命令及其参数
var (
	// 子命令：list
	listCmd = flag.NewFlagSet("list", flag.ExitOnError)

	// 子命令：run
	runCmd = flag.NewFlagSet("run", flag.ExitOnError)
	runID  = runCmd.Int("id", 0, "任务ID")

	// 子命令：add
	addCmd    = flag.NewFlagSet("add", flag.ExitOnError)
	addName   = addCmd.String("n", "", "任务名")
	addTarget = addCmd.String("t", "", "目标目录")
	addBackup = addCmd.String("b", "", "备份目录(指定任务存放的父目录路径)")
	addKeep   = addCmd.Int("k", 3, "保留数量")

	// 子命令：delete
	deleteCmd  = flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID   = deleteCmd.Int("id", 0, "任务ID")
	deleteName = deleteCmd.String("n", "", "任务名")
	deleteDirF = deleteCmd.Bool("d", false, "在删除任务时，是否同时删除备份文件。若启用此选项，备份文件将被一同删除")

	// 子命令：edit
	editCmd  = flag.NewFlagSet("edit", flag.ExitOnError)
	editName = editCmd.String("n", "", "任务名")
	editID   = editCmd.Int("id", 0, "任务ID")
	editKeep = editCmd.Int("k", 3, "保留数量")

	// 子命令：log
	logCmd   = flag.NewFlagSet("log", flag.ExitOnError)
	logLimit = logCmd.Int("l", 10, "显示的行数")

	// 子命令：show
	showCmd = flag.NewFlagSet("show", flag.ExitOnError)
	showID  = showCmd.Int("id", 0, "任务ID")

	// 子命令：unpack
	unpackCmd       = flag.NewFlagSet("unpack", flag.ExitOnError)
	unpackID        = unpackCmd.Int("id", 0, "任务ID")
	unpackVersionID = unpackCmd.String("v", "", "指定解压的版本ID")
	unpackOutput    = unpackCmd.String("o", ".", "指定输出的路径(默认当前目录)")
	//unpackForce     = unpackCmd.Bool("f", false, "表示强制覆盖")

	// 子命令：version
	versionCmd = flag.NewFlagSet("version", flag.ExitOnError)

	// 子命令：help
	helpCmd = flag.NewFlagSet("help", flag.ExitOnError)
)

// 定义子命令的执行逻辑
func ExecuteCommands(db *sqlx.DB, args []string) error {
	switch args[0] {
	case "list":
		// 解析list命令的参数
		listCmd.Parse(args[1:])
		// 执行list命令的逻辑
		if err := listCmdMain(db); err != nil {
			return fmt.Errorf("列出项目列表失败: %v", err)
		}
		return nil
	case "l":
		// 解析list命令的参数
		listCmd.Parse(args[1:])
		// 执行list命令的逻辑
		if err := listCmdMain(db); err != nil {
			return fmt.Errorf("列出项目列表失败: %v", err)
		}
		return nil
	case "run":
		// 解析run命令的参数
		runCmd.Parse(args[1:])
		// 执行run命令的逻辑
		if err := runCmdMain(db); err != nil {
			return fmt.Errorf("执行备份任务失败: %v", err)
		}
		return nil
	case "r":
		// 解析run命令的参数
		runCmd.Parse(args[1:])
		// 执行run命令的逻辑
		if err := runCmdMain(db); err != nil {
			return fmt.Errorf("执行备份任务失败: %v", err)
		}
		return nil
	case "add":
		// 解析add命令的参数
		addCmd.Parse(args[1:])
		// 执行add命令的逻辑
		if err := addCmdMain(db); err != nil {
			return fmt.Errorf("添加项目失败: %v", err)
		}
		return nil
	case "a":
		// 解析add命令的参数
		addCmd.Parse(args[1:])
		// 执行add命令的逻辑
		if err := addCmdMain(db); err != nil {
			return fmt.Errorf("添加项目失败: %v", err)
		}
		return nil
	case "delete":
		// 解析delete命令的参数
		deleteCmd.Parse(args[1:])
		// 执行delete命令的逻辑
		if err := deleteCmdMain(db); err != nil {
			return fmt.Errorf("删除项目失败: %v", err)
		}
		return nil
	case "d":
		// 解析delete命令的参数
		deleteCmd.Parse(args[1:])
		// 执行delete命令的逻辑
		if err := deleteCmdMain(db); err != nil {
			return fmt.Errorf("删除项目失败: %v", err)
		}
		return nil
	case "edit":
		// 解析edit命令的参数
		editCmd.Parse(args[1:])
		// 执行edit命令的逻辑
		if err := editCmdMain(db); err != nil {
			return fmt.Errorf("编辑项目失败: %v", err)
		}
		return nil
	case "e":
		// 解析edit命令的参数
		editCmd.Parse(args[1:])
		// 执行edit命令的逻辑
		if err := editCmdMain(db); err != nil {
			return fmt.Errorf("编辑项目失败: %v", err)
		}
		return nil
	case "log":
		// 解析log命令的参数
		logCmd.Parse(args[1:])
		// 执行log命令的逻辑
		if err := logCmdMain(db, 1, *logLimit); err != nil {
			return fmt.Errorf("查看日志失败: %v", err)
		}
		return nil
	case "show":
		// 解析show命令的参数
		showCmd.Parse(args[1:])
		// 执行show命令的逻辑
		if err := showCmdMain(db); err != nil {
			return fmt.Errorf("查看指定备份任务的信息失败: %v", err)
		}
		return nil
	case "s":
		// 解析show命令的参数
		showCmd.Parse(args[1:])
		// 执行show命令的逻辑
		if err := showCmdMain(db); err != nil {
			return fmt.Errorf("查看指定备份任务的信息失败: %v", err)
		}
		return nil
	case "unpack":
		// 解析unpack命令的参数
		unpackCmd.Parse(args[1:])
		// 执行unpack命令的逻辑
		if err := unpackCmdMain(db); err != nil {
			return fmt.Errorf("解压备份任务失败: %v", err)
		}
		return nil
	case "u":
		// 解析unpack命令的参数
		unpackCmd.Parse(args[1:])
		// 执行unpack命令的逻辑
		if err := unpackCmdMain(db); err != nil {
			return fmt.Errorf("解压备份任务失败: %v", err)
		}
		return nil
	// 打印版本信息
	case "version":
		versionCmd.Parse(args[1:])
		v := version.Get()
		if versionInfo, err := v.SprintVersion("text"); err != nil {
			CL.PrintError(err)
			os.Exit(1)
		} else {
			CL.Green(versionInfo)
		}
		return nil
	// 打印帮助信息
	case "help":
		// 解析help命令的参数
		helpCmd.Parse(args[1:])
		// 执行help命令的逻辑
		fmt.Println(HelpText)
		return nil
	// 未知命令
	default:
		return fmt.Errorf("未知命令: %s", args[0])
	}
}

// add命令的执行逻辑
func addCmdMain(db *sqlx.DB) error {
	// 检查任务名是否为空
	if *addName == "" {
		return fmt.Errorf("任务名不能为空")
	}

	// 检查目标目录是否为空
	if *addTarget == "" {
		return fmt.Errorf("目标目录不能为空")
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

	// 如果备份目录为空, 则使用默认值
	if *addBackup == "" {
		tempHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %w", err)
		}
		*addBackup = filepath.Join(tempHome, ".cbk", "data", *addName)

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

		// 构建备份目录的绝对路径
		*addBackup = filepath.Join(*addBackup, *addName)

		// 检查备份目录是否存在
		if _, err := tools.CheckPath(*addBackup); err != nil {
			if err := os.MkdirAll(*addBackup, 0755); err != nil {
				return fmt.Errorf("备份目录创建失败: %w", err)
			}
		}
	}

	// 扩展目标目录为绝对路径
	AbsAddTarget, err := filepath.Abs(*addTarget)
	if err != nil {
		return fmt.Errorf("获取目标目录绝对路径失败: %w", err)
	}

	// 插入新任务到数据库
	insertSql := "insert into backup_tasks(task_name, target_directory, backup_directory, retention_count) values(?, ?, ?, ?)"
	if _, err := db.Exec(insertSql, *addName, AbsAddTarget, *addBackup, *addKeep); err != nil {
		return fmt.Errorf("插入任务失败: %w", err)
	}

	// 打印成功信息
	CL.PrintSuccessf("任务添加成功: %s", *addName)
	return nil
}

// delete命令的执行逻辑
func deleteCmdMain(db *sqlx.DB) error {
	// 如果deleteName不为空, 则根据任务名删除任务
	if *deleteName != "" {
		// 获取备份存放目录
		var backupDir string
		backupDirSql := "select backup_directory from backup_tasks where task_name = ?"
		if err := db.Get(&backupDir, backupDirSql, *deleteName); err == sql.ErrNoRows {
			return fmt.Errorf("任务名不存在 %s", *deleteName)
		} else if err != nil {
			return fmt.Errorf("获取备份存放目录失败: %w", err)
		}

		// 如果deleteDir为true, 则删除备份存放目录
		if *deleteDirF {
			// 删除备份存放目录(如果存在)
			if _, err := tools.CheckPath(backupDir); err == nil {
				if err := os.RemoveAll(backupDir); err != nil {
					return fmt.Errorf("删除备份存放目录失败: %w", err)
				} else {
					// 打印成功信息
					CL.PrintSuccessf("备份存放目录删除成功: %s", backupDir)
				}
			}
		} else {
			CL.PrintWarningf("请在稍后，手动删除备份存放目录: %s", backupDir)
		}

		// 删除任务
		deleteSql := "delete from backup_tasks where task_name = ?"
		if _, err := db.Exec(deleteSql, *deleteName); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}

		// 删除备份记录
		deleteBackupSql := "delete from backup_records where task_name = ?"
		if _, err := db.Exec(deleteBackupSql, *deleteName); err != nil {
			return fmt.Errorf("删除备份记录失败: %w", err)
		}

		// 打印成功信息
		CL.PrintSuccessf("任务删除成功: %s", *deleteName)

		return nil
	}

	// 如果deleteID不为0, 则根据任务ID删除任务
	if *deleteID != 0 {
		// 获取备份存放目录
		var backupDir string
		backupDirSql := "select backup_directory from backup_tasks where task_id = ?"
		if err := db.Get(&backupDir, backupDirSql, *deleteID); err == sql.ErrNoRows {
			return fmt.Errorf("任务ID不存在 %d", *deleteID)
		} else if err != nil {
			return fmt.Errorf("获取备份存放目录失败: %w", err)
		}

		// 如果deleteDir为true, 则删除备份存放目录
		if *deleteDirF {
			// 删除备份存放目录(如果存在)
			if _, err := tools.CheckPath(backupDir); err == nil {
				if err := os.RemoveAll(backupDir); err != nil {
					return fmt.Errorf("删除备份存放目录失败: %w", err)
				} else {
					// 打印成功信息
					CL.PrintSuccessf("备份存放目录删除成功: %s", backupDir)
				}
			}
		} else {
			CL.PrintWarningf("请在稍后，手动删除备份存放目录: %s", backupDir)
		}

		// 删除任务
		deleteSql := "delete from backup_tasks where task_id = ?"
		if _, err := db.Exec(deleteSql, *deleteID); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}

		// 删除备份记录
		deleteBackupSql := "delete from backup_records where task_id = ?"
		if _, err := db.Exec(deleteBackupSql, *deleteID); err != nil {
			return fmt.Errorf("删除备份记录失败: %w", err)
		}

		// 打印成功信息
		CL.PrintSuccessf("任务ID删除成功: %d", *deleteID)

		return nil
	}

	return fmt.Errorf("删除任务时, 必须指定任务名或任务ID")
}

// listCmdMain 查询并打印任务列表
func listCmdMain(db *sqlx.DB) error {
	// 查询所有任务
	querySql := "SELECT task_id, task_name, target_directory, backup_directory, retention_count FROM backup_tasks;"

	var tasks []struct {
		TaskID          int    `db:"task_id"`          // 任务ID
		TaskName        string `db:"task_name"`        // 任务名
		TargetDirectory string `db:"target_directory"` // 目标目录
		BackupDirectory string `db:"backup_directory"` // 备份目录
		RetentionCount  int    `db:"retention_count"`  // 保留数量
	}

	if err := db.Select(&tasks, querySql); err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout) // 使用标准输出作为输出目标
	t.AppendHeader(table.Row{"ID", "任务名", "保留数量", "目标目录", "备份目录"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ID", Align: text.AlignCenter},
		{Name: "任务名", Align: text.AlignLeft},
		{Name: "保留数量", Align: text.AlignCenter},
		{Name: "目标目录", Align: text.AlignLeft},
		{Name: "备份目录", Align: text.AlignLeft},
	})

	// 添加数据行
	for _, task := range tasks {
		t.AppendRow(table.Row{
			task.TaskID,
			task.TaskName,
			task.RetentionCount,
			task.TargetDirectory,
			task.BackupDirectory,
		})
	}

	// 设置表格样式
	// t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ID", Align: text.AlignCenter},   // 居中对齐
		{Name: "任务名", Align: text.AlignLeft},    // 左对齐
		{Name: "保留数量", Align: text.AlignCenter}, // 居中对齐
		{Name: "目标目录", Align: text.AlignLeft},   // 左对齐
		{Name: "备份目录", Align: text.AlignLeft},   // 左对齐
	})
	// t.SetColumnConfigs([]table.ColumnConfig{
	// 	{Name: "ID", WidthMax: 4},
	// 	{Name: "任务名", WidthMax: 10},
	// 	{Name: "保留数量", WidthMax: 10},
	// 	{Name: "目标目录", WidthMax: 30},
	// 	{Name: "备份目录", WidthMax: 30},
	// })

	// 渲染表格
	t.Render()

	return nil
}

// editCmdMain 编辑任务
func editCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *editID == 0 {
		return fmt.Errorf("编辑任务时, 必须指定任务ID")
	}

	// 编辑任务
	editSql := "select task_name, retention_count from backup_tasks where task_id =?"
	var task struct {
		TaskName       string `db:"task_name"`       // 任务名
		RetentionCount int    `db:"retention_count"` // 保留数量
	}
	if err := db.Get(&task, editSql, *editID); err == sql.ErrNoRows {
		return fmt.Errorf("任务ID不存在 %d", *editID)
	} else if err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 如果指定了-n参数, 则更新任务名
	if *editName != "" {
		task.TaskName = *editName
	}

	// 如果指定了-k参数, 则更新保留数量
	if *editKeep != 0 {
		task.RetentionCount = *editKeep
	}

	// 更新任务
	updateSql := "update backup_tasks set task_name = ?, retention_count = ? where task_id = ?"
	if _, err := db.Exec(updateSql, task.TaskName, task.RetentionCount, *editID); err != nil {
		return fmt.Errorf("更新任务失败: %w", err)
	}

	// 如果更新任务名, 则更新备份目录
	if *editName != "" {
		// 获取备份存放目录
		var backupDir string
		backupDirSql := "select backup_directory from backup_tasks where task_id =?"
		if err := db.Get(&backupDir, backupDirSql, *editID); err != nil {
			return fmt.Errorf("获取备份存放目录失败: %w", err)
		}

		// 如果存在则，重命名备份存放目录
		// if _, err := tools.CheckPath(backupDir); err == nil {
		// 	if err := os.Rename(backupDir, filepath.Join(filepath.Dir(backupDir), *editName)); err != nil {
		// 		return fmt.Errorf("重命名备份存放目录失败: %w", err)
		// 	}
		// }
		CL.PrintWarningf("请在稍后，手动重命名备份存放目录: %s -> %s", backupDir, filepath.Join(filepath.Dir(backupDir), *editName))

		// 更新备份存放目录
		updateBackupDirSql := "update backup_tasks set backup_directory = ? where task_id = ?"
		if _, err := db.Exec(updateBackupDirSql, filepath.Join(filepath.Dir(backupDir), *editName), *editID); err != nil {
			return fmt.Errorf("更新备份存放目录失败: %w", err)
		}

		// 打印成功信息
		CL.PrintSuccessf("备份存放目录已自动跟随任务名修改为: %s", filepath.Join(filepath.Dir(backupDir), *editName))
	}

	// 打印成功信息
	CL.PrintSuccess("更新成功!")

	// 查询并打印分别打印任务ID的当前任务信息
	querySql := "select task_name, retention_count from backup_tasks where task_id =?"
	if err := db.Get(&task, querySql, *editID); err != nil {
		return fmt.Errorf("查询更新后的任务失败: %w", err)
	}

	// 打印任务信息
	CL.PrintSuccessf("任务ID %d 的当前任务名为: %s", *editID, task.TaskName)
	CL.PrintSuccessf("任务ID %d 的当前保留数量为: %d", *editID, task.RetentionCount)

	return nil
}

// runCmdMain 运行备份任务
func runCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *runID == 0 {
		return fmt.Errorf("运行备份任务时, 必须指定任务ID")
	}

	// 获取任务信息
	var task struct {
		TaskName        string `db:"task_name"`        // 任务名
		TargetDirectory string `db:"target_directory"` // 目标目录
		BackupDirectory string `db:"backup_directory"` // 备份目录
		RetentionCount  int    `db:"retention_count"`  // 保留数量
	}
	querySql := "select task_name, target_directory, backup_directory, retention_count from backup_tasks where task_id =?"
	if err := db.Get(&task, querySql, *runID); err == sql.ErrNoRows {
		return fmt.Errorf("任务ID不存在 %d", *runID)
	} else if err != nil {
		return fmt.Errorf("获取任务信息失败: %w", err)
	}

	// 检查目标目录或文件是否存在
	if _, err := tools.CheckPath(task.TargetDirectory); err != nil {
		return fmt.Errorf("目标目录或文件不存在: %w", err)
	}

	// 检查备份目录是否存在
	if _, err := tools.CheckPath(task.BackupDirectory); err != nil {
		if err := os.MkdirAll(task.BackupDirectory, 0755); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}
	}

	// 构建备份文件名
	backupTime := time.Now().Format("20060102150405")
	backupFileNamePrefix := fmt.Sprintf("%s_%s", task.TaskName, backupTime)

	// 获取versionID
	versionID := tools.GenerateID(6)

	// 运行备份任务
	targetDir := filepath.Dir(task.TargetDirectory)                                 // 获取目标目录的目录部分
	targetName := filepath.Base(task.TargetDirectory)                               // 获取目标目录的最后一个部分
	backupFileNamePath := filepath.Join(task.BackupDirectory, backupFileNamePrefix) // 获取构建的备份文件路径

	// 执行备份任务
	CL.PrintSuccess("开始备份任务...")
	zipPath, err := tools.CompressFilesByOS(db, targetDir, targetName, backupFileNamePath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, data_status, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "0", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("备份任务失败: %w", err)
	}

	// 获取备份文件的后8位MD5哈希值
	backupFileMD5, err := tools.GetFileMD5Last8(zipPath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, data_status, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "0", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("获取备份文件MD5失败: %w", err)
	}

	// 获取备份文件的大小
	backupFileSize, err := tools.HumanReadableSize(zipPath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, data_status, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "0", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("获取备份文件大小失败: %w", err)
	}

	// 插入备份记录
	insertSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, data_status, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?,?)"
	if _, err := db.Exec(insertSql, versionID, *runID, backupTime, task.TaskName, "true", filepath.Base(zipPath), backupFileSize, task.BackupDirectory, "1", backupFileMD5); err != nil {
		return fmt.Errorf("插入备份记录失败: %w", err)
	}

	// 获取备份目录下的以指定扩展名的文件列表
	fileExtensionSql := "select file_extension from compress_config where os_type = ?"
	var fileExtension string
	if err := db.Get(&fileExtension, fileExtensionSql, runtime.GOOS); err != nil {
		return fmt.Errorf("获取文件扩展名失败: %w", err)
	}
	zipFiles, err := tools.GetZipFiles(task.BackupDirectory, fileExtension)
	if err != nil {
		return fmt.Errorf("获取备份目录下的.zip文件失败: %w", err)
	}

	// 删除多余的备份文件
	if len(zipFiles) > task.RetentionCount {
		if err := tools.RetainLatestFiles(db, zipFiles, task.RetentionCount); err != nil {
			return fmt.Errorf("删除多余的备份文件失败: %w", err)
		}
	}

	// 打印备份信息
	CL.PrintSuccessf("备份任务 %s 完成", task.TaskName)
	CL.PrintSuccessf("备份文件: %s", zipPath)
	CL.PrintSuccessf("备份文件大小: %s", backupFileSize)
	CL.PrintSuccessf("备份文件MD5: %s", backupFileMD5)
	CL.PrintSuccessf("备份文件版本ID: %s", versionID)

	return nil
}

// logCmdMain 函数，支持分页查询
// 参数：
//
//	db - 数据库连接
//	page - 页码
//	pageSize - 每页记录数
//
// 返回值：
//
//	error - 如果发生错误，返回错误信息；否则返回 nil
func logCmdMain(db *sqlx.DB, page, pageSize int) error {
	// 验证分页参数
	if page < 1 {
		return fmt.Errorf("页码必须从1开始")
	}
	if pageSize <= 0 {
		return fmt.Errorf("每页记录数必须大于0")
	}

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 定义查询语句
	querySql := `
		SELECT version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash
		FROM backup_records
		where data_status = 1
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?;
	`

	// 定义结构体来接收查询结果
	var records []struct {
		VersionID      string `db:"version_id"`       // 版本ID
		TaskID         int    `db:"task_id"`          // 任务ID
		Timestamp      string `db:"timestamp"`        // 时间戳
		TaskName       string `db:"task_name"`        // 任务名
		BackupStatus   string `db:"backup_status"`    // 备份状态
		BackupFileName string `db:"backup_file_name"` // 备份文件名
		BackupSize     string `db:"backup_size"`      // 备份文件大小
		BackupPath     string `db:"backup_path"`      // 备份文件路径
		VersionHash    string `db:"version_hash"`     // 版本哈希
	}

	// 执行查询
	if err := db.Select(&records, querySql, pageSize, offset); err == sql.ErrNoRows {
		return fmt.Errorf("没有查询到任何记录")
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	}

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(log.Writer())
	t.AppendHeader(table.Row{"版本ID", "任务ID", "时间戳", "任务名", "备份状态", "备份文件名", "备份文件大小", "备份存放目录", "版本哈希"})

	// 遍历查询结果，将数据添加到表格中
	for _, record := range records {
		// 将时间戳转换为时间对象并格式化为易读格式
		timestamp, err := time.Parse("20060102150405", record.Timestamp)
		if err != nil {
			return fmt.Errorf("解析时间戳失败: %w", err)
		}
		formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

		// 将数据添加到表格中
		t.AppendRow(table.Row{
			record.VersionID,
			record.TaskID,
			formattedTimestamp,
			record.TaskName,
			record.BackupStatus,
			record.BackupFileName,
			record.BackupSize,
			record.BackupPath,
			record.VersionHash,
		})
	}

	// 设置表格样式
	//t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "版本ID", WidthMax: 10},
		{Name: "任务ID", WidthMax: 10},
		{Name: "时间戳", WidthMax: 20},
		{Name: "任务名", WidthMax: 20},
		{Name: "备份状态", WidthMax: 10},
		{Name: "备份文件名", WidthMax: 20},
		{Name: "备份文件大小", WidthMax: 10},
		{Name: "备份存放目录", WidthMax: 30},
		{Name: "版本哈希", WidthMax: 20},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "版本ID", Align: text.AlignCenter},
		{Name: "任务ID", Align: text.AlignCenter},
		{Name: "时间戳", Align: text.AlignLeft},
		{Name: "任务名", Align: text.AlignLeft},
		{Name: "备份状态", Align: text.AlignCenter},
		{Name: "备份文件名", Align: text.AlignLeft},
		{Name: "备份文件大小", Align: text.AlignCenter},
		{Name: "备份存放目录", Align: text.AlignLeft},
		{Name: "版本哈希", Align: text.AlignCenter},
	})

	// 打印表格
	t.Render()

	return nil
}

// showCmdMain 查询指定任务ID的备份记录并以表格形式输出
func showCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *showID == 0 {
		return fmt.Errorf("查询指定备份任务时, 必须指定任务ID")
	}

	// 构建查询sql语句
	querySql := "SELECT version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash FROM backup_records WHERE task_id = ? AND data_status = '1' ORDER BY timestamp DESC"

	// 定义存储查询结果的结构体
	var records []struct {
		VersionID      string `db:"version_id"`       // 版本ID
		TaskID         int    `db:"task_id"`          // 任务ID
		Timestamp      string `db:"timestamp"`        // 时间戳
		TaskName       string `db:"task_name"`        // 任务名
		BackupStatus   string `db:"backup_status"`    // 备份状态
		BackupFileName string `db:"backup_file_name"` // 备份文件名
		BackupSize     string `db:"backup_size"`      // 备份文件大小
		BackupPath     string `db:"backup_path"`      // 备份文件路径
		VersionHash    string `db:"version_hash"`     // 版本哈希
	}

	// 执行查询
	if err := db.Select(&records, querySql, *showID); err == sql.ErrNoRows {
		return fmt.Errorf("未找到指定任务ID %d 的备份记录", *showID)
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	}

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(log.Writer())
	t.AppendHeader(table.Row{"版本ID", "任务ID", "时间戳", "任务名", "备份状态", "备份文件名", "备份文件大小", "备份文件路径", "版本哈希"})

	// 将查询结果添加到表格
	for _, record := range records {
		// 将时间戳转换为时间对象并格式化为易读格式
		timestamp, err := time.Parse("20060102150405", record.Timestamp)
		if err != nil {
			return fmt.Errorf("解析时间戳失败: %w", err)
		}
		formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

		// 将数据添加到表格中
		t.AppendRow(table.Row{
			record.VersionID,
			record.TaskID,
			formattedTimestamp, // 格式化后的时间戳
			record.TaskName,
			record.BackupStatus,
			record.BackupFileName,
			record.BackupSize,
			record.BackupPath,
			record.VersionHash,
		})
	}

	// 设置表格样式
	//t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "版本ID", WidthMax: 10},
		{Name: "任务ID", WidthMax: 10},
		{Name: "时间戳", WidthMax: 20},
		{Name: "任务名", WidthMax: 20},
		{Name: "备份状态", WidthMax: 10},
		{Name: "备份文件名", WidthMax: 20},
		{Name: "备份文件大小", WidthMax: 10},
		{Name: "备份存放目录", WidthMax: 30},
		{Name: "版本哈希", WidthMax: 20},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "版本ID", Align: text.AlignCenter},
		{Name: "任务ID", Align: text.AlignCenter},
		{Name: "时间戳", Align: text.AlignLeft},
		{Name: "任务名", Align: text.AlignLeft},
		{Name: "备份状态", Align: text.AlignCenter},
		{Name: "备份文件名", Align: text.AlignLeft},
		{Name: "备份文件大小", Align: text.AlignCenter},
		{Name: "备份存放目录", Align: text.AlignLeft},
		{Name: "版本哈希", Align: text.AlignCenter},
	})

	// 输出表格
	t.Render()

	return nil
}

func unpackCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *unpackID == 0 {
		return fmt.Errorf("解压指定备份任务时, 必须指定任务ID")
	}

	// 检查versionID是否指定
	if *unpackVersionID == "" {
		return fmt.Errorf("解压指定备份任务时, 必须指定版本ID")
	}

	// 检查*unpackID是否是已存在的
	var taskCount int
	if err := db.Get(&taskCount, "SELECT COUNT(*) FROM backup_records WHERE task_id = ? AND data_status = '1';", *unpackID); err == sql.ErrNoRows {
		return fmt.Errorf("未找到指定任务ID %d 的备份记录", *unpackID)
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	} else if taskCount == 0 {
		return fmt.Errorf("未找到指定任务ID %d 的备份记录", *unpackID)
	}

	// 检查versionID是否是已存在的
	var versionCount int
	if err := db.Get(&versionCount, "SELECT COUNT(*) FROM backup_records WHERE version_id = ? AND data_status = '1';", *unpackVersionID); err == sql.ErrNoRows {
		return fmt.Errorf("未找到指定版本ID %s 的备份记录", *unpackVersionID)
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	} else if versionCount == 0 {
		return fmt.Errorf("未找到指定版本ID %s 的备份记录", *unpackVersionID)
	}

	// 构建查询sql语句
	querySql := "SELECT version_id, task_id, backup_file_name, backup_path, version_hash FROM backup_records WHERE task_id =? AND version_id =? AND data_status = '1';"

	// 定义存储查询结果的结构体
	var record struct {
		VersionID      string `db:"version_id"`       // 版本ID
		TaskID         int    `db:"task_id"`          // 任务ID
		BackupFileName string `db:"backup_file_name"` // 备份文件名
		BackupPath     string `db:"backup_path"`      // 存放路径
		VersionHash    string `db:"version_hash"`     // 版本哈希
	}

	// 执行查询
	if err := db.Get(&record, querySql, *unpackID, *unpackVersionID); err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	} else if record.VersionID == "" {
		return fmt.Errorf("未找到指定任务ID和版本ID的备份记录")
	}

	// 构建备份文件路径
	backupFilePath := filepath.Join(record.BackupPath, record.BackupFileName)

	// 检查备份文件是否存在
	if _, err := tools.CheckPath(backupFilePath); err != nil {
		return fmt.Errorf("备份文件不存在: %w", err)
	}

	// 获取备份文件的后8位哈希值
	if backupFileHash, err := tools.GetFileMD5Last8(backupFilePath); err != nil {
		return fmt.Errorf("获取备份文件哈希失败: %w", err)
	} else if backupFileHash != record.VersionHash {
		return fmt.Errorf("备份文件 %s 的版本 %s 的哈希值与记录不匹配，文件可能已损坏或被篡改。请尝试选择其他版本的备份文件重试", backupFilePath, record.VersionID)
	}

	// 执行解压操作
	CL.PrintSuccess("开始解压备份文件...")
	if outPath, err := tools.UncompressFilesByOS(db, record.BackupPath, record.BackupFileName, *unpackOutput); err != nil {
		return fmt.Errorf("解压备份文件失败: %w", err)
	} else {
		CL.PrintSuccessf("已解压 %s -> %s", record.BackupFileName, outPath)
	}

	return nil
}
