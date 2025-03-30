// cmd.go
package cmd

import (
	"cbk/pkg/tools"
	"cbk/pkg/version"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitee.com/MM-Q/colorlib"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jmoiron/sqlx"
)

// 定义全局颜色渲染器
var CL = colorlib.NewColorLib()

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

	// 子命令：edit
	editCmd  = flag.NewFlagSet("edit", flag.ExitOnError)
	editName = editCmd.String("n", "", "任务名")
	editID   = editCmd.Int("id", 0, "任务ID")
	editKeep = editCmd.Int("k", 3, "保留数量")

	// // 子命令：log
	// logCmd = flag.NewFlagSet("log", flag.ExitOnError)

	// // 子命令：show
	// showCmd = flag.NewFlagSet("show", flag.ExitOnError)

	// // 子命令：unpack
	// unpackCmd = flag.NewFlagSet("unpack", flag.ExitOnError)
	// unpackName = unpackCmd.String("name", "", "任务名")
	// unpackNameShort = unpackCmd.String("n", "", "任务名")
	// unpackID = unpackCmd.Int("id", 0, "任务ID")
	// unpackVersion = unpackCmd.String("version", "", "指定解压的版本")
	// unpackVersionShort = unpackCmd.String("v", "", "指定解压的版本")
	// unpackOutput = unpackCmd.String("output", "", "指定输出的路径")
	// unpackOutputShort = unpackCmd.String("o", "", "指定输出的路径")
	// unpackForce = unpackCmd.Bool("force", false, "表示强制覆盖")
	// unpackForceShort = unpackCmd.Bool("f", false, "表示强制覆盖")

	// 子命令：version
	versionCmd = flag.NewFlagSet("version", flag.ExitOnError)

	// // 子命令：help
	// helpCmd = flag.NewFlagSet("help", flag.ExitOnError)
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
		fmt.Printf("查看备份日志: %s\n", args[1])
		fmt.Println(args[1:])
	case "show":
		fmt.Printf("查看指定备份任务的信息: %s\n", args[1])
		fmt.Println(args[1:])
	case "s":
		fmt.Printf("查看指定备份任务的信息: %s\n", args[1])
		fmt.Println(args[1:])
	case "unpack":
		fmt.Printf("解压备份任务: %s, %s, %s, %s, %s\n", args[1], args[2], args[3], args[4], args[5])
		fmt.Println(args[1:])
	case "u":
		fmt.Printf("解压备份任务: %s, %s, %s, %s, %s\n", args[1], args[2], args[3], args[4], args[5])
		fmt.Println(args[1:])
	// 打印版本信息
	case "version":
		versionCmd.Parse(args[1:])
		v := version.Get()
		CL.Green(v.SprintVersion("text"))
	// 打印帮助信息
	case "help":
		fmt.Println("帮助信息")
		fmt.Println(args[1:])
	// 未知命令
	default:
		return fmt.Errorf("未知命令: %s", args[0])
	}

	return nil
}

// add命令的执行逻辑
func addCmdMain(db *sqlx.DB) error {
	// 检查参数是否为空
	if *addName == "" {
		return fmt.Errorf("任务名不能为空")
	}

	if *addTarget == "" {
		return fmt.Errorf("目标目录不能为空")
	}

	// 如果备份目录为空, 则使用默认值
	if *addBackup == "" {
		tempHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %w", err)
		}
		*addBackup = filepath.Join(tempHome, ".cbk", "data", *addName)
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

	// 检查备份目录是否存在
	if _, err := tools.CheckPath(*addBackup); err != nil {
		if err := os.MkdirAll(*addBackup, 0755); err != nil {
			return fmt.Errorf("备份目录创建失败: %w", err)
		}
	}

	// 插入新任务到数据库
	insertSql := "insert into backup_tasks(task_name, target_directory, backup_directory, retention_count) values(?, ?, ?, ?)"
	if _, err := db.Exec(insertSql, *addName, *addTarget, *addBackup, *addKeep); err != nil {
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
		// 检查任务是否存在
		checkSql := "select count(*) from backup_tasks where task_name =?"
		var count int
		if err := db.Get(&count, checkSql, *deleteName); err != nil {
			return fmt.Errorf("检查任务失败: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("任务名不存在 %s", *deleteName)
		}

		// 获取备份存放目录
		var backupDir string
		backupDirSql := "select backup_directory from backup_tasks where task_name = ?"
		if err := db.Get(&backupDir, backupDirSql, *deleteName); err != nil {
			return fmt.Errorf("获取备份存放目录失败: %w", err)
		}

		// 删除备份存放目录(如果存在)
		if _, err := tools.CheckPath(backupDir); err == nil {
			if err := os.RemoveAll(backupDir); err != nil {
				return fmt.Errorf("删除备份存放目录失败: %w", err)
			} else {
				// 打印成功信息
				CL.PrintSuccessf("备份存放目录删除成功: %s", backupDir)
			}
		}

		// 删除任务
		deleteSql := "delete from backup_tasks where task_name = ?"
		if _, err := db.Exec(deleteSql, *deleteName); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}

		// 打印成功信息
		CL.PrintSuccessf("任务删除成功: %s", *deleteName)

		return nil
	}

	// 如果deleteID不为0, 则根据任务ID删除任务
	if *deleteID != 0 {
		// 检查任务ID是否存在
		checkSql := "select count(*) from backup_tasks where task_id =?"
		var count int
		if err := db.Get(&count, checkSql, *deleteID); err != nil {
			return fmt.Errorf("检查任务失败: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("任务ID不存在 %d", *deleteID)
		}

		// 获取备份存放目录
		var backupDir string
		backupDirSql := "select backup_directory from backup_tasks where task_id = ?"
		if err := db.Get(&backupDir, backupDirSql, *deleteID); err != nil {
			return fmt.Errorf("获取备份存放目录失败: %w", err)
		}

		// 删除备份存放目录(如果存在)
		// if _, err := tools.CheckPath(backupDir); err == nil {
		// 	if err := os.RemoveAll(backupDir); err != nil {
		// 		return fmt.Errorf("删除备份存放目录失败: %w", err)
		// 	} else {
		// 		// 打印成功信息
		// 		CL.PrintSuccessf("备份存放目录删除成功: %s", backupDir)
		// 	}
		// }
		CL.PrintWarningf("请在稍后，手动删除备份存放目录: %s", backupDir)

		// 删除任务
		deleteSql := "delete from backup_tasks where task_id = ?"
		if _, err := db.Exec(deleteSql, *deleteID); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
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

	// 检查任务是否存在
	checkSql := "select count(*) from backup_tasks where task_id =?"
	var count int
	if err := db.Get(&count, checkSql, *editID); err != nil {
		return fmt.Errorf("检查任务失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("任务ID不存在 %d", *editID)
	}

	// 编辑任务
	editSql := "select task_name, retention_count from backup_tasks where task_id =?"
	var task struct {
		TaskName       string `db:"task_name"`       // 任务名
		RetentionCount int    `db:"retention_count"` // 保留数量
	}
	if err := db.Get(&task, editSql, *editID); err != nil {
		return fmt.Errorf("编辑任务失败: %w", err)
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

func runCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *runID == 0 {
		return fmt.Errorf("运行备份任务时, 必须指定任务ID")
	}

	// 检查任务是否存在
	checkSql := "select count(*) from backup_tasks where task_id =?"
	var count int
	if err := db.Get(&count, checkSql, *runID); err != nil {
		return fmt.Errorf("检查任务失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("任务ID不存在 %d", *runID)
	}

	// 获取任务信息
	var task struct {
		TaskName        string `db:"task_name"`        // 任务名
		TargetDirectory string `db:"target_directory"` // 目标目录
		BackupDirectory string `db:"backup_directory"` // 备份目录
		RetentionCount  int    `db:"retention_count"`  // 保留数量
	}
	querySql := "select task_name, target_directory, backup_directory, retention_count from backup_tasks where task_id =?"
	if err := db.Get(&task, querySql, *runID); err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
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
	zipPath, err := tools.CompressFilesByOS(targetDir, targetName, backupFileNamePath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("备份任务失败: %w", err)
	}

	// 获取备份文件的后8位MD5哈希值
	backupFileMD5, err := tools.GetFileMD5Last8(zipPath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("获取备份文件MD5失败: %w", err)
	}

	// 获取备份文件的大小
	backupFileSize, err := tools.HumanReadableSize(zipPath)
	if err != nil {
		errorSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		if _, err := db.Exec(errorSql, versionID, *runID, backupTime, task.TaskName, "false", "-", "-", "-", "-"); err != nil {
			return fmt.Errorf("插入备份记录失败: %w", err)
		}
		return fmt.Errorf("获取备份文件大小失败: %w", err)
	}

	// 插入备份记录
	insertSql := "insert into backup_records (version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if _, err := db.Exec(insertSql, versionID, *runID, backupTime, task.TaskName, "true", zipPath, backupFileSize, task.BackupDirectory, backupFileMD5); err != nil {
		return fmt.Errorf("插入备份记录失败: %w", err)
	}

	// 获取备份目录下的以.zip结尾的文件列表
	zipFiles, err := tools.GetZipFiles(task.BackupDirectory)
	if err != nil {
		return fmt.Errorf("获取备份目录下的.zip文件失败: %w", err)
	}

	// 删除多余的备份文件
	if len(zipFiles) > task.RetentionCount {
		if err := tools.RetainLatestFiles(zipFiles, task.RetentionCount); err != nil {
			return fmt.Errorf("删除多余的备份文件失败: %w", err)
		}
	}

	// 打印备份信息
	CL.PrintSuccessf("备份任务 %s 完成", task.TaskName)
	CL.PrintSuccessf("备份文件: %s", filepath.Join(task.BackupDirectory, zipPath))
	CL.PrintSuccessf("备份文件大小: %s", backupFileSize)
	CL.PrintSuccessf("备份文件MD5: %s", backupFileMD5)
	CL.PrintSuccessf("备份文件版本ID: %s", versionID)

	return nil
}
