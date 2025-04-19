package cmd

import (
	"cbk/pkg/tools"
	"fmt"
	"os"
	"runtime"
)

// initCmdMain 自动补全主逻辑
func initCmdMain(t string) error {
	// 检查自动补全类型是否为空
	if t == "" {
		return fmt.Errorf("请指定生成的类型, 例如: 'cbk init -type [bash|addtask]'")
	}

	switch t {
	case "bash":
		// 检查是否为Linux或Mac系统
		if runtime.GOOS != "linux" {
			return fmt.Errorf("自动补全 'bash' 仅在Linux上受支持")
		}

		// 打印自动补全脚本
		fmt.Println(BashCompletion)
		return nil
	case "addtask":
		// 检查当前目录是否存在add_task.yaml文件
		if _, err := tools.CheckPath("add_task.yaml"); err == nil {
			return fmt.Errorf("当前目录已存在add_task.yaml文件")
		}

		// 写入 AddTaskTemplate 的内容到当前目录下的add_task.yaml文件中
		if err := os.WriteFile("add_task.yaml", []byte(AddTaskTemplate), 0644); err != nil {
			return fmt.Errorf("写入配置文件失败: %w", err)
		}

		// 打印提示信息
		CL.PrintOk("add_task.yaml配置文件已创建, 请根据需要修改后运行 'cbk add -f add_task.yaml' 命令添加备份任务")
		return nil
	default:
		return fmt.Errorf("未知的类型: %s", t)
	}
}
