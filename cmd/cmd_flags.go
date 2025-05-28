package cmd

import (
	"flag"

	_ "embed"

	"github.com/jedib0t/go-pretty/v6/table"
)

//go:embed help/help.txt
var HelpText string

// 定义存放表格样式的MAP
var (
	TableStyle = map[string]table.Style{
		"default":     table.StyleDefault,       // 默认样式
		"bold":        table.StyleBold,          // 加粗样式
		"colorbright": table.StyleColoredBright, // 亮色样式
		"colordark":   table.StyleColoredDark,   // 暗色样式
		"double":      table.StyleDouble,        // 双边框样式
		"light":       table.StyleLight,         // 浅色样式
		"rounded":     table.StyleRounded,       // 圆角样式
		"bd":          table.StyleBold,          // 加粗样式
		"cb":          table.StyleColoredBright, // 亮色样式
		"cd":          table.StyleColoredDark,   // 暗色样式
		"de":          table.StyleDouble,        // 双边框样式
		"lt":          table.StyleLight,         // 浅色样式
		"ro":          table.StyleRounded,       // 圆角样式
	}
)

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

//go:embed sql/init.sql
var initSql string // 初始化SQL语句

// 定义子命令及其参数
var (
	// 子命令: list
	listCmd          = flag.NewFlagSet("list", flag.ExitOnError)
	listTableStyle   = listCmd.String("ts", "default", "表格样式(default, bold, colorbright, colordark, double, light, rounded, bd, cb, cd, de, lt, ro)")
	listNoTable      = listCmd.Bool("no-table", false, "是否禁用表格输出")
	listNoTableShort = listCmd.Bool("nt", false, "是否禁用表格输出")

	// 子命令: run
	runCmd = flag.NewFlagSet("run", flag.ExitOnError)
	runID  = runCmd.Int("id", 0, "任务ID")
	runIDS = runCmd.String("ids", "", "任务ID列表, 多个ID用逗号分隔")

	// 子命令: add
	addCmd            = flag.NewFlagSet("add", flag.ExitOnError)
	addName           = addCmd.String("n", "", "任务名")
	addTarget         = addCmd.String("t", "", "目标目录路径")
	addBackup         = addCmd.String("b", "", "备份存放路径(默认: 用户主目录/.cbk/data/[项目名]/")
	addRetentionCount = addCmd.Int("c", 3, "保留数量")
	addRetentionDays  = addCmd.Int("d", 0, "保留天数")
	addBackupDirName  = addCmd.String("bn", "", "备份目录名(默认: 目标目录名)")
	addNoCompression  = addCmd.Int("nc", 0, "是否禁用压缩(0: 启用压缩, 1: 禁用压缩)")
	addConfig         = addCmd.String("f", "", "指定YAML格式的配置文件路径, 用于批量添加任务(格式参考: add_task.yaml)")
	addExcludeRules   = addCmd.String("ex", "none", "指定要排除的目录名、文件名、扩展名, 用于排除备份文件, 支持通配符模式(默认为none, 不排除任何文件)")

	// 子命令: delete
	deleteCmd       = flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID        = deleteCmd.Int("id", 0, "任务ID")
	deleteIDS       = deleteCmd.String("ids", "", "任务ID列表, 多个ID用逗号分隔")
	deleteName      = deleteCmd.String("n", "", "任务名")
	deleteDirF      = deleteCmd.Bool("d", false, "在删除任务时，是否同时删除备份文件。若启用此选项，备份文件将被一同删除")
	deleteVersionID = deleteCmd.String("v", "", "指定要删除的备份版本ID")

	// 子命令: edit
	editCmd            = flag.NewFlagSet("edit", flag.ExitOnError)
	editID             = editCmd.Int("id", -1, "指定要编辑的备份任务ID")
	editIDS            = editCmd.String("ids", "", "指定要编辑的备份任务ID列表, 多个ID用逗号分隔")
	editName           = editCmd.String("n", "", "指定新的任务名。如果未指定，则任务名保持不变")
	editRetentionCount = editCmd.Int("c", -1, "指定备份文件的保留数量。如果未指定，则保留数量保持不变")
	editRetentionDays  = editCmd.Int("d", -1, "指定备份文件的保留天数。如果未指定，则保留天数保持不变")
	editNewDirName     = editCmd.String("bn", "", "指定新的备份目录名。如果未指定，则备份目录名保持不变")
	editNoCompression  = editCmd.Int("nc", -1, "是否禁用压缩(0: 启用压缩, 1: 禁用压缩, -1: 不修改)")
	editExcludeRules   = editCmd.String("ex", "", "指定要排除的目录名、文件名、扩展名, 用于排除备份文件, 支持通配符模式")

	// 子命令: log
	logCmd          = flag.NewFlagSet("log", flag.ExitOnError)
	logLimit        = logCmd.Int("l", 10, "显示的行数")
	logView         = logCmd.Bool("v", false, "是否显示详细日志")
	logTableStyle   = logCmd.String("ts", "default", "表格样式(default, bold, colorbright, colordark, double, light, rounded, bd, cb, cd, de, lt, ro)")
	logNoTable      = logCmd.Bool("no-table", false, "是否禁用表格输出")
	logNoTableShort = logCmd.Bool("nt", false, "是否禁用表格输出")

	// 子命令: show
	showCmd          = flag.NewFlagSet("show", flag.ExitOnError)
	showID           = showCmd.Int("id", 0, "任务ID")
	showView         = showCmd.Bool("v", false, "是否显示详细信息")
	showTableStyle   = showCmd.String("ts", "default", "表格样式(default, bold, colorbright, colordark, double, light, rounded, bd, cb, cd, de, lt, ro)")
	showNoTable      = showCmd.Bool("no-table", false, "是否禁用表格输出")
	showNoTableShort = showCmd.Bool("nt", false, "是否禁用表格输出")

	// 子命令: unpack
	unpackCmd       = flag.NewFlagSet("unpack", flag.ExitOnError)
	unpackID        = unpackCmd.Int("id", 0, "任务ID")
	unpackVersionID = unpackCmd.String("v", "", "指定解压的版本ID")
	unpackOutput    = unpackCmd.String("o", ".", "指定输出的路径(默认当前目录)")

	// 子命令: zip
	zipCmd           = flag.NewFlagSet("zip", flag.ExitOnError)
	zipOutput        = zipCmd.String("o", "未命名.zip", "指定输出的压缩包名(默认: 未命名.zip)")
	zipTarget        = zipCmd.String("t", "", "指定要打包的目标路径")
	zipNoCompression = zipCmd.Int("nc", 0, "是否禁用压缩（默认启用压缩）")
	zipExcludeRules  = zipCmd.String("ex", "none", "指定要排除的目录名、文件名、扩展名, 用于排除备份文件, 支持通配符模式")

	// 子命令: unzip
	unzipCmd       = flag.NewFlagSet("unzip", flag.ExitOnError)
	unzipFile      = unzipCmd.String("f", "", "指定要解压的压缩文件名")
	unzipOutputDir = unzipCmd.String("d", ".", "指定解压的目标路径。如果未指定，则解压到当前目录")

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
