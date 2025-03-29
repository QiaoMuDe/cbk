// cmd.go
package cmd

import (
	"cbk/pkg/version"
	"flag"
	"fmt"

	"gitee.com/MM-Q/colorlib"
)

// 定义全局颜色渲染器
var CL = colorlib.NewColorLib()

// 定义子命令及其参数
var (
	// // 子命令：list
	// listCmd = flag.NewFlagSet("list", flag.ExitOnError)

	// // 子命令：run
	// runCmd = flag.NewFlagSet("run", flag.ExitOnError)

	// // 子命令：add
	// addCmd = flag.NewFlagSet("add", flag.ExitOnError)
	// addName = addCmd.String("name", "", "任务名")
	// addNameShort = addCmd.String("n", "", "任务名")
	// addTarget = addCmd.String("target", "", "目标目录")
	// addTargetShort = addCmd.String("t", "", "目标目录")
	// addBackup = addCmd.String("backup", "", "备份目录")
	// addBackupShort = addCmd.String("b", "", "备份目录")
	// addKeep = addCmd.Int("keep", 3, "保留数量")
	// addKeepShort = addCmd.Int("k", 3, "保留数量")

	// // 子命令：delete
	// deleteCmd = flag.NewFlagSet("delete", flag.ExitOnError)

	// // 子命令：edit
	// editCmd = flag.NewFlagSet("edit", flag.ExitOnError)
	// editName = editCmd.String("name", "", "任务名")
	// editNameShort = editCmd.String("n", "", "任务名")
	// editID = editCmd.Int("id", 0, "任务ID")
	// editTarget = editCmd.String("target", "", "目标目录")
	// editTargetShort = editCmd.String("t", "", "目标目录")
	// editBackup = editCmd.String("backup", "", "备份目录")
	// editBackupShort = editCmd.String("b", "", "备份目录")
	// editKeep = editCmd.Int("keep", 3, "保留数量")
	// editKeepShort = editCmd.Int("k", 3, "保留数量")

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
func ExecuteCommands(args []string) error {
	switch args[0] {
	case "list":
		fmt.Println("列出项目列表")
		fmt.Println(args[1:])
	case "l":
		fmt.Println("列出项目列表")
		fmt.Println(args[1:])
	case "run":
		fmt.Printf("执行备份任务: %s\n", args[1])
		fmt.Println(args[1:])
	case "r":
		fmt.Printf("执行备份任务: %s\n", args[1])
		fmt.Println(args[1:])
	case "add":
		fmt.Printf("添加备份任务: %s, %s, %s, %s\n", args[1], args[2], args[3], args[4])
		fmt.Println(args[1:])
	case "a":
		fmt.Printf("添加备份任务: %s, %s, %s, %s\n", args[1], args[2], args[3], args[4])
		fmt.Println(args[1:])
	case "delete":
		fmt.Printf("删除备份任务: %s\n", args[1])
		fmt.Println(args[1:])
	case "d":
		fmt.Printf("删除备份任务: %s\n", args[1])
		fmt.Println(args[1:])
	case "edit":
		fmt.Printf("编辑备份任务: %s, %s, %s, %s\n", args[1], args[2], args[3], args[4])
		fmt.Println(args[1:])
	case "e":
		fmt.Printf("编辑备份任务: %s, %s, %s, %s\n", args[1], args[2], args[3], args[4])
		fmt.Println(args[1:])
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
