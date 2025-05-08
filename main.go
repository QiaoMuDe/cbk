// main.go
package main

import (
	"cbk/cmd"
	"cbk/pkg/tools"
	"os"
)

// 入口点
func main() {
	// 在返回时捕获处理
	defer func() {
		if err := recover(); err != nil {
			tools.CL.PrintErrf("程序发生异常: %v", err)
			tools.CL.PrintErr("程序退出")
			os.Exit(1)
		}
	}()

	// 运行程序
	if err := cmd.AppRun(); err != nil {
		tools.CL.PrintErrf("程序发生错误: %v", err)
		tools.CL.PrintErr("程序退出")
		os.Exit(1)
	}

	// 程序正常退出
	os.Exit(0)
}
