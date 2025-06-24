package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/tools"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jmoiron/sqlx"
)

// initCmdMain 自动补全主逻辑
func initCmdMain(t string, db *sqlx.DB) error {
	// 检查自动补全类型是否为空
	if t == "" {
		return fmt.Errorf("请指定生成的类型, 例如: 'cbk init -type [bash|b|addtask|at|onekey|ok]'")
	}

	switch t {
	case "bash":
		// 检查是否为Linux或Mac系统
		if runtime.GOOS != "linux" {
			return fmt.Errorf("自动补全 'bash' 仅在Linux上受支持")
		}

		// 打印自动补全脚本
		fmt.Println(globals.BashCompletion)
		return nil
	case "b":
		// 检查是否为Linux或Mac系统
		if runtime.GOOS != "linux" {
			return fmt.Errorf("自动补全 'bash' 仅在Linux上受支持")
		}

		// 打印自动补全脚本
		fmt.Println(globals.BashCompletion)
		return nil
	case "addtask":
		// 检查当前目录是否存在add_task.yaml文件
		if _, err := tools.CheckPath("add_task.yaml"); err == nil {
			return fmt.Errorf("当前目录已存在add_task.yaml文件")
		}

		// 写入 AddTaskTemplate 的内容到当前目录下的add_task.yaml文件中
		if err := os.WriteFile("add_task.yaml", []byte(globals.AddTaskTemplate), 0644); err != nil {
			return fmt.Errorf("写入配置文件失败: %w", err)
		}

		// 打印提示信息
		CL.PrintOk("add_task.yaml配置文件已创建, 请根据需要修改后运行 'cbk add -f add_task.yaml' 命令添加备份任务")
		return nil
	case "at":
		// 检查当前目录是否存在add_task.yaml文件
		if _, err := tools.CheckPath("add_task.yaml"); err == nil {
			return fmt.Errorf("当前目录已存在add_task.yaml文件")
		}

		// 写入 AddTaskTemplate 的内容到当前目录下的add_task.yaml文件中
		if err := os.WriteFile("add_task.yaml", []byte(globals.AddTaskTemplate), 0644); err != nil {
			return fmt.Errorf("写入配置文件失败: %w", err)
		}

		// 打印提示信息
		CL.PrintOk("add_task.yaml配置文件已创建, 请根据需要修改后运行 'cbk add -f add_task.yaml' 命令添加备份任务")
		return nil
	case "onekey":
		// 生成一键脚本
		if err := initShellScript(db); err != nil {
			return fmt.Errorf("生成一键脚本失败: %w", err)
		}
		return nil
	case "ok":
		// 生成一键脚本
		if err := initShellScript(db); err != nil {
			return fmt.Errorf("生成一键脚本失败: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("未知的类型: %s", t)
	}
}

// 一键生成自动运行任务脚本
func initShellScript(db *sqlx.DB) error {
	// 检查输出路径不为空的时候检查是否存在
	if initOut.Get() != "" {
		// 检查是否已经存在同名文件
		if _, err := tools.CheckPath(initOut.Get()); err == nil {
			return fmt.Errorf("指定的输出路径 %s 已存在同名文件, 请指定其他路径", initOut.Get())
		}
	}

	// 获取输出路径
	var outPutPath string
	if initOut.Get() != "" {
		outPutPath = initOut.Get()
	} else {
		// 获取用户家目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户家目录失败: %w", err)
		}

		// 拼接输出路径
		outPutPath = filepath.Join(homeDir, globals.OneKeyScriptName)
	}

	// 构建查询sql
	querySql := "SELECT task_id FROM backup_tasks;"

	// 执行查询
	var taskIDs globals.BackupTasks
	if err := db.Select(&taskIDs, querySql); err != nil {
		return fmt.Errorf("查询任务ID失败: %w", err)
	}

	// 检查是否存在任务
	if len(taskIDs) == 0 {
		return fmt.Errorf("未找到任何任务")
	}

	// 获取当前操作系统
	osName := runtime.GOOS

	// 根据操作系统构建脚本内容
	switch osName {
	case "windows":
		// 检查输出路径是否为.bat文件
		if filepath.Ext(outPutPath) != ".bat" {
			outPutPath += ".bat"
		}

		// 打开文件
		file, err := os.OpenFile(outPutPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("打开文件失败: %w", err)
		}
		defer file.Close()

		// 写入脚本内容
		if _, err := file.WriteString("@echo off\n"); err != nil {
			return fmt.Errorf("写入脚本内容失败: %w", err)
		}

		// 循环写入任务ID对应的脚本内容
		for _, task := range taskIDs {
			if _, err := file.WriteString(fmt.Sprintf("start cbk run -id %d\n", task.TaskID)); err != nil {
				return fmt.Errorf("写入脚本内容失败: %w", err)
			}
		}

		// 打印提示信息
		CL.PrintOkf("已生成到 %s 路径下\n", outPutPath)
	case "linux", "darwin":
		// 检查输出路径是否为.sh文件
		if filepath.Ext(outPutPath) != ".sh" {
			outPutPath += ".sh"
		}

		// 打开文件
		file, err := os.OpenFile(outPutPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("打开文件失败: %w", err)
		}
		defer file.Close()

		// 写入脚本内容
		if _, err := file.WriteString("#!/bin/bash\n"); err != nil {
			return fmt.Errorf("写入脚本内容失败: %w", err)
		}

		// 循环写入任务ID对应的脚本内容
		for _, task := range taskIDs {
			if _, err := file.WriteString(fmt.Sprintf("( cbk run -id %d ) &\n pid%d=$!\n", task.TaskID, task.TaskID)); err != nil {
				return fmt.Errorf("写入脚本内容失败: %w", err)
			}
		}

		// 构建等待命令
		waitCmd := "wait"
		for _, task := range taskIDs {
			waitCmd += fmt.Sprintf(" $pid%d", task.TaskID)
		}

		// 写入等待命令
		if _, err := file.WriteString(waitCmd + "\n"); err != nil {
			return fmt.Errorf("写入脚本内容失败: %w", err)
		}

		// 打印提示信息
		CL.PrintOkf("已生成到 %s 路径下\n", outPutPath)
	default:
		return fmt.Errorf("不支持的操作系统: %s, 仅支持Windows, Linux, MacOS", osName)
	}

	return nil
}
