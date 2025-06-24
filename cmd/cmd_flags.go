package cmd

import (
	"cbk/pkg/globals"
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/verman"

	_ "embed"
)

//go:embed help/help.txt
var HelpText string

//go:embed help/help_list.txt
var HelpListText string // 定义子命令: list的帮助文本

//go:embed help/help_run.txt
var HelpRunText string // 定义子命令: run的帮助文本

//go:embed help/help_add.txt
var HelpAddText string // 定义子命令: add的帮助文本

//go:embed help/help_delete.txt
var HelpDeleteText string // 定义子命令: delete的帮助文本

//go:embed help/help_edit.txt
var HelpEditText string // 定义子命令: edit的帮助文本

//go:embed help/help_log.txt
var HelpLogText string // 定义子命令: log的帮助文本

//go:embed help/help_show.txt
var HelpShowText string // 定义子命令: show的帮助文本

//go:embed help/help_zip.txt
var HelpZipText string // 定义子命令: zip的帮助文本

//go:embed help/help_unzip.txt
var HelpUnzipText string // 定义子命令: unzip的帮助文本

//go:embed help/help_unpack.txt
var HelpUnpackText string // 定义子命令: unpack的帮助文本

//go:embed help/help_clear.txt
var HelpClearText string // 定义子命令: clear的帮助文本
//go:embed help/help_init.txt
var HelpInitText string // 定义子命令: init的帮助文本

//go:embed autocomplete/bash/cbk.sh
var BashCompletion string // 定义bash补全脚本

//go:embed templates/add_task.yaml
var AddTaskTemplate string // 定义添加任务的模板文件

//go:embed help/help_export.txt
var HelpExportText string // 定义子命令: export的帮助文本

var (
	// 根命令
	versionF *qflag.BoolFlag

	// list 子命令
	listCmd        *qflag.Cmd
	listTableStyle *qflag.EnumFlag

	// run 子命令
	runCmd *qflag.Cmd
	runID  *qflag.IntFlag
	runIDS *qflag.StringFlag

	// 子命令: add
	addCmd            *qflag.Cmd
	addName           *qflag.StringFlag
	addTarget         *qflag.StringFlag
	addBackup         *qflag.StringFlag
	addRetentionCount *qflag.IntFlag
	addRetentionDays  *qflag.IntFlag
	addBackupDirName  *qflag.StringFlag
	addNoCompression  *qflag.IntFlag
	addConfig         *qflag.StringFlag
	addExcludeRules   *qflag.StringFlag

	// 子命令: delete
	deleteCmd       *qflag.Cmd
	deleteID        *qflag.IntFlag
	deleteIDS       *qflag.StringFlag
	deleteName      *qflag.StringFlag
	deleteDirF      *qflag.BoolFlag
	deleteVersionID *qflag.StringFlag

	// 子命令: edit
	editCmd            *qflag.Cmd
	editID             *qflag.IntFlag
	editIDS            *qflag.StringFlag
	editName           *qflag.StringFlag
	editRetentionCount *qflag.IntFlag
	editRetentionDays  *qflag.IntFlag
	editNewDirName     *qflag.StringFlag
	editNoCompression  *qflag.IntFlag
	editExcludeRules   *qflag.StringFlag

	// 子命令: log
	logCmd        *qflag.Cmd
	logLimit      *qflag.IntFlag
	logView       *qflag.BoolFlag
	logTableStyle *qflag.EnumFlag

	// 子命令: show
	showCmd        *qflag.Cmd
	showID         *qflag.IntFlag
	showView       *qflag.BoolFlag
	showTableStyle *qflag.EnumFlag

	// 子命令: unpack
	unpackCmd       *qflag.Cmd
	unpackID        *qflag.IntFlag
	unpackVersionID *qflag.StringFlag
	unpackOutput    *qflag.StringFlag

	// 子命令: zip
	zipCmd           *qflag.Cmd
	zipOutput        *qflag.StringFlag
	zipTarget        *qflag.StringFlag
	zipNoCompression *qflag.IntFlag
	zipExcludeRules  *qflag.StringFlag

	// 子命令: unzip
	unzipCmd       *qflag.Cmd
	unzipFile      *qflag.StringFlag
	unzipOutputDir *qflag.StringFlag
)

// 定义子命令及其参数
var (
	// 子命令: version
	versionCmd = flag.NewFlagSet("version", flag.ExitOnError)

	// 子命令: help
	helpCmd = flag.NewFlagSet("help", flag.ExitOnError)

	// 子命令: clear
	clearCmd     = flag.NewFlagSet("clear", flag.ExitOnError)
	clearConfirm = clearCmd.Bool("confirm", false, "确认是否执行清空数据操作")

	// 子命令: init
	initCmd  = flag.NewFlagSet("complete", flag.ExitOnError)
	initType = initCmd.String("t", "", "指定要生成的配置类型, 可选值: bash, addtask")
	initOut  = initCmd.String("o", "", "指定输出文件路径")

	// 子命令: export
	exportCmd = flag.NewFlagSet("export", flag.ExitOnError)
	exportID  = exportCmd.Int("id", 0, "指定要导出的任务ID")
	exportAll = exportCmd.Bool("all", false, "导出所有任务")
)

func init() {
	// 根命令
	qflag.SetUseChinese(true) // 设置使用中文
	qflag.SetDescription("命令行备份任务管理工具, 用于管理备份任务，包括添加、运行、编辑、删除备份任务，查看任务日志，显示任务详情等")
	versionF = qflag.Bool("version", "v", false, "显示版本信息")

	// 子命令: list
	listCmd = qflag.NewCmd("list", "ls", flag.ExitOnError)
	listCmd.SetUseChinese(true)                    // 设置使用中文
	listCmd.SetModuleHelps(globals.TableStyleHelp) // 设置模块帮助信息
	listCmd.SetDescription("列出所有备份任务的概览信息")
	listTableStyle = listCmd.Enum("table-style", "ts", "df", "指定表格样式", globals.TableStyleList)

	// 子命令: run
	runCmd = qflag.NewCmd("run", "r", flag.ExitOnError)
	runCmd.SetUseChinese(true) // 设置使用中文
	runCmd.AddNote("任务ID: 任务ID是必需的, 用于标识要运行的备份任务")
	runCmd.AddNote("任务配置: 备份任务的配置 (如目标路径、备份路径、保留数量等) 在任务创建时已经设置, 运行任务时将按照这些配置执行")
	runCmd.SetDescription("运行指定的备份任务")
	runID = runCmd.Int("id", "", 0, "任务ID")
	runIDS = runCmd.String("ids", "", "", "任务ID列表, 多个ID用逗号分隔")

	// 子命令: add
	addCmd = qflag.NewCmd("add", "a", flag.ExitOnError)
	addCmd.SetUseChinese(true) // 设置使用中文
	addCmd.SetDescription("添加一个新的备份任务")
	addCmd.AddExample(qflag.ExampleInfo{Description: "添加一个名为“任务5”的备份任务, 目标目录为“/home/user/documents”, 不排除任何文件或文件夹", Usage: "cbk add -n '任务5' -t '/home/user/documents' -ex 'none'"})
	addCmd.AddExample(qflag.ExampleInfo{Description: "使用指定的YAML配置文件添加一个任务", Usage: "cbk add -f /path/to/add_task.yaml"})
	addName = addCmd.String("name", "n", "", "任务名")
	addTarget = addCmd.String("target", "t", "", "目标目录路径")
	addBackup = addCmd.String("backup", "b", "", "备份存放路径(默认: 用户主目录/.cbk/data/[项目名]/")
	addRetentionCount = addCmd.Int("count", "c", 3, "保留数量")
	addRetentionDays = addCmd.Int("day", "d", 0, "保留天数")
	addBackupDirName = addCmd.String("bak-name", "bn", "", "备份目录名(默认: 目标目录名)")
	addNoCompression = addCmd.Int("nc", "", 0, "是否禁用压缩(0: 启用压缩, 1: 禁用压缩)")
	addConfig = addCmd.String("cfg", "f", "", "指定YAML格式的配置文件路径, 用于批量添加任务(格式参考: add_task.yaml)")
	addExcludeRules = addCmd.String("ex", "", "none", "指定要排除的目录名、文件名、扩展名, 用于排除备份文件, 支持通配符模式")

	// 子命令: delete
	deleteCmd = qflag.NewCmd("delete", "d", flag.ExitOnError)
	deleteCmd.SetUseChinese(true) // 设置使用中文
	deleteCmd.AddNote("如果启用了 `-d` 选项，还会同时删除与该任务相关的备份文件")
	deleteCmd.SetDescription("删除指定的备份任务")
	deleteID = deleteCmd.Int("id", "", 0, "任务ID")
	deleteIDS = deleteCmd.String("ids", "", "", "任务ID列表, 多个ID用逗号分隔")
	deleteName = deleteCmd.String("name", "n", "", "任务名")
	deleteDirF = deleteCmd.Bool("d", "", false, "在删除任务时，是否同时删除备份文件。若启用此选项，备份文件将被一同删除")
	deleteVersionID = deleteCmd.String("vid", "v", "", "指定要删除的备份版本ID")

	// 子命令: edit
	editCmd = qflag.NewCmd("edit", "e", flag.ExitOnError)
	editCmd.SetUseChinese(true) // 设置使用中文
	editCmd.SetDescription("编辑指定备份任务的配置信息")
	editID = editCmd.Int("id", "", -1, "指定要编辑的备份任务ID")
	editIDS = editCmd.String("ids", "", "", "指定要编辑的备份任务ID列表, 多个ID用逗号分隔")
	editName = editCmd.String("name", "n", "", "指定新的任务名。如果未指定，则任务名保持不变")
	editRetentionCount = editCmd.Int("count", "c", -1, "指定备份文件的保留数量。如果未指定，则保留数量保持不变")
	editRetentionDays = editCmd.Int("day", "d", -1, "指定备份文件的保留天数。如果未指定，则保留天数保持不变")
	editNewDirName = editCmd.String("backup-name", "bn", "", "指定新的备份目录名。如果未指定，则备份目录名保持不变")
	editNoCompression = editCmd.Int("nc", "", -1, "是否禁用压缩(0: 启用压缩, 1: 禁用压缩, -1: 不修改)")
	editExcludeRules = editCmd.String("ex", "", "", "指定要排除的目录名、文件名、扩展名, 用于排除备份文件, 支持通配符模式")

	// 子命令: log
	logCmd = qflag.NewCmd("log", "l", flag.ExitOnError)
	logCmd.SetUseChinese(true) // 设置使用中文
	logCmd.SetDescription("查看备份任务的日志信息")
	logLimit = logCmd.Int("limit", "l", 10, "显示的行数")
	logView = logCmd.Bool("view", "v", false, "是否显示详细日志")
	logTableStyle = logCmd.Enum("table-style", "ts", "df", "指定表格样式", globals.TableStyleList)

	// 子命令: show
	showCmd = qflag.NewCmd("show", "s", flag.ExitOnError)
	showCmd.SetUseChinese(true) // 设置使用中文
	showCmd.SetDescription("查看指定备份任务的元数据信息")
	showID = showCmd.Int("id", "", 0, "任务ID")
	showView = showCmd.Bool("view", "v", false, "是否显示详细信息")
	showTableStyle = showCmd.Enum("table-style", "ts", "df", "指定表格样式", globals.TableStyleList)

	// 子命令: unpack
	unpackCmd = qflag.NewCmd("unpack", "up", flag.ExitOnError)
	unpackCmd.SetUseChinese(true) // 设置使用中文
	unpackCmd.SetDescription("根据指定的任务ID解压指定版本的备份文件")
	unpackID = unpackCmd.Int("id", "", 0, "任务ID")
	unpackVersionID = unpackCmd.String("vid", "v", "", "指定解压的版本ID")
	unpackOutput = unpackCmd.String("output", "o", ".", "指定输出的路径(默认当前目录)")

	// 子命令: zip
	zipCmd = qflag.NewCmd("zip", "z", flag.ExitOnError)
	zipCmd.SetUseChinese(true) // 设置使用中文
	zipCmd.SetDescription("将指定的目标路径打包为一个压缩文件")
	zipOutput = zipCmd.String("output", "o", "未命名.zip", "指定输出的压缩包名(默认: 未命名.zip)")
	zipTarget = zipCmd.String("target", "t", "", "指定要打包的目标路径")
	zipNoCompression = zipCmd.Int("nc", "", 0, "是否禁用压缩(0: 启用压缩, 1: 禁用压缩)")
	zipExcludeRules = zipCmd.String("ex", "", "none", "指定要排除的目录名、文件名、扩展名, 用于排除备份文件, 支持通配符模式")

	// 子命令: unzip
	unzipCmd = qflag.NewCmd("unzip", "uz", flag.ExitOnError)
	unzipCmd.SetUseChinese(true) // 设置使用中文
	unpackCmd.SetDescription("解压指定的压缩文件到目标路径。如果未指定目标路径，则解压到当前目录")
	unzipFile = unzipCmd.String("file", "f", "", "指定要解压的压缩文件名")
	unzipOutputDir = unzipCmd.String("dir", "d", ".", "指定解压的目标路径。如果未指定，则解压到当前目录")

	// 添加子命令
	if addErr := qflag.AddSubCmd(listCmd, runCmd, addCmd, deleteCmd, editCmd, logCmd, showCmd, unpackCmd, zipCmd, unzipCmd); addErr != nil {
		fmt.Printf("err: 添加子命令失败: %s\n", addErr)
		os.Exit(1)
	}

	// 解析标志
	if parseErr := qflag.Parse(); parseErr != nil {
		fmt.Printf("err: %s\n", parseErr)
		os.Exit(1)
	}

	// 检查版本标志
	if versionF.Get() {
		v := verman.Get()
		fmt.Println(v.AppName, v.GitVersion)
		os.Exit(0)
	}
}
