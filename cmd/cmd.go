// cmd.go
package cmd

import (
	"cbk/pkg/globals"
	"cbk/pkg/version"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/colorlib"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jmoiron/sqlx"
)

// 定义全局颜色渲染器
var CL = colorlib.NewColorLib()

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
	initType = initCmd.String("type", "", "指定要生成的配置类型, 可选值: bash, addtask")

	// 子命令: export
	exportCmd = flag.NewFlagSet("export", flag.ExitOnError)
	exportID  = exportCmd.Int("id", 0, "指定要导出的任务ID")
	exportAll = exportCmd.Bool("all", false, "导出所有任务")
)

// 初始化子命令的帮助信息
func init() {
	// 初始化list命令的帮助信息
	listCmd.Usage = func() {
		fmt.Println(HelpListText)
	}

	// 初始化run命令的帮助信息
	runCmd.Usage = func() {
		fmt.Println(HelpRunText)
	}

	// 初始化add命令的帮助信息
	addCmd.Usage = func() {
		fmt.Println(HelpAddText)
	}

	// 初始化delete命令的帮助信息
	deleteCmd.Usage = func() {
		fmt.Println(HelpDeleteText)
	}

	// 初始化edit命令的帮助信息
	editCmd.Usage = func() {
		fmt.Println(HelpEditText)
	}

	// 初始化log命令的帮助信息
	logCmd.Usage = func() {
		fmt.Println(HelpLogText)
	}

	// 初始化show命令的帮助信息
	showCmd.Usage = func() {
		fmt.Println(HelpShowText)
	}

	// 初始化unpack命令的帮助信息
	unpackCmd.Usage = func() {
		fmt.Println(HelpUnpackText)
	}

	// 初始化zip命令的帮助信息
	zipCmd.Usage = func() {
		fmt.Println(HelpZipText)
	}

	// 初始化unzip命令的帮助信息
	unzipCmd.Usage = func() {
		fmt.Println(HelpUnzipText)
	}

	// 初始化clear命令的帮助信息
	clearCmd.Usage = func() {
		fmt.Println(HelpClearText)
	}

	// 初始化init命令的帮助信息
	initCmd.Usage = func() {
		fmt.Println(HelpInitText)
	}

	// 初始化export命令的帮助信息
	exportCmd.Usage = func() {
		fmt.Println(HelpExportText)
	}
}

// 程序运行入口
func AppRun() error {
	// 初始化数据库
	db, initDBErr := initDB()
	if initDBErr != nil {
		return fmt.Errorf("初始化数据库失败: %w", initDBErr)
	}

	// 在返回时关闭数据库连接
	defer func() {
		// 检查数据库是否打开，如果打开则关闭
		if db != nil {
			if closeErr := db.Close(); closeErr != nil {
				CL.PrintErrf("关闭数据库连接失败: %v", closeErr)
			}
		}
	}()

	// 初始化数据目录
	if initDataDirErr := initDataDir(); initDataDirErr != nil {
		return fmt.Errorf("初始化数据目录失败: %w", initDataDirErr)
	}

	// 主标志
	vFlag := flag.Bool("v", false, "显示版本信息")
	vvFlag := flag.Bool("vv", false, "显示更详细的版本信息")
	hFlag := flag.Bool("h", false, "显示帮助信息")
	helpFlag := flag.Bool("help", false, "显示帮助信息")

	// 解析主标志
	flag.Parse()

	// 打印版本信息
	if *vFlag {
		v := version.Get()
		if versionInfo, newErr := v.SprintVersion("simple"); newErr != nil {
			return fmt.Errorf("获取版本信息失败: %w", newErr)
		} else {
			CL.Green(versionInfo)
		}
		return nil
	}

	// 打印更详细的版本信息
	if *vvFlag {
		v := version.Get()
		if versionInfo, newErr := v.SprintVersion("text"); newErr != nil {
			return fmt.Errorf("获取版本信息失败: %w", newErr)
		} else {
			CL.Green(versionInfo)
		}
		return nil
	}

	// 打印帮助信息
	if *hFlag || *helpFlag {
		fmt.Println(HelpText)
		return nil
	}

	// 获取命令行参数
	args := flag.Args()

	// 检查是否有子命令
	if len(args) == 0 {
		flag.PrintDefaults()
		return nil
	}

	// 执行子命令
	if execCmdErr := executeCommands(db, args); execCmdErr != nil {
		return fmt.Errorf("执行子命令失败: %w", execCmdErr)
	}

	return nil
}

// 初始化数据库
// 返回值:
// *sqlx.DB: 数据库连接
// error: 错误信息
func initDB() (*sqlx.DB, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户主目录失败: %w", err)
	}

	// 构造数据库目录路径
	dbDir := filepath.Join(homeDir, globals.CbkHomeDir)

	// 检查数据库目录是否存在, 如果不存在, 则创建
	if _, statErr := os.Stat(dbDir); os.IsNotExist(statErr) {
		if mkdirErr := os.MkdirAll(dbDir, 0755); mkdirErr != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %w", mkdirErr)
		}
	} else if statErr != nil {
		return nil, fmt.Errorf("检查数据库目录失败: %w", statErr)
	}

	// 构造数据库文件路径
	dbPath := filepath.Join(dbDir, globals.CbkDBFile)

	// 检查数据库文件是否存在, 如果不存在, 则创建并初始化
	if _, statErr := os.Stat(dbPath); os.IsNotExist(statErr) {
		// 连接数据库
		db, connectErr := sqlx.Connect("sqlite3", dbPath)
		if connectErr != nil {
			return nil, fmt.Errorf("连接数据库失败: %w", connectErr)
		}

		// 执行初始化SQL语句
		if _, execErr := db.Exec(initSql); execErr != nil {
			return nil, fmt.Errorf("执行初始化SQL语句失败: %w", execErr)
		}

		// 直接返回连接
		return db, nil
	} else if statErr != nil {
		return nil, fmt.Errorf("连接到数据库文件失败: %w", statErr)
	}

	// 存在数据库文件, 直接连接
	db, connectErr := sqlx.Connect("sqlite3", dbPath)
	if connectErr != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", connectErr)
	}

	return db, nil
}

// 初始化数据目录
// 返回值:
// error: 错误信息
func initDataDir() error {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %w", err)
	}

	// 构造数据目录路径
	dataDir := filepath.Join(homeDir, globals.CbkHomeDir, globals.CbkDataDir)

	// 检查数据目录是否存在, 如果不存在, 则创建
	if _, statErr := os.Stat(dataDir); os.IsNotExist(statErr) {
		if mkdirErr := os.MkdirAll(dataDir, 0755); mkdirErr != nil {
			return fmt.Errorf("创建数据目录失败: %w", mkdirErr)
		}
	} else if statErr != nil {
		return fmt.Errorf("检查数据目录失败: %w", statErr)
	}

	return nil
}

// 定义子命令的执行逻辑
func executeCommands(db *sqlx.DB, args []string) error {
	switch args[0] {
	case "list":
		// 解析list命令的参数
		if err := listCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析list命令参数失败: %v", err)
		}
		// 执行list命令的逻辑
		if err := listCmdMain(db); err != nil {
			return fmt.Errorf("列出项目列表失败: %v", err)
		}
		return nil
	case "l":
		// 解析list命令的参数
		if err := listCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析list命令参数失败: %v", err)
		}
		// 执行list命令的逻辑
		if err := listCmdMain(db); err != nil {
			return fmt.Errorf("列出项目列表失败: %v", err)
		}
		return nil
	case "run":
		// 解析run命令的参数
		if err := runCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析run命令参数失败: %v", err)
		}
		// 执行run命令的逻辑
		if err := runCmdMain(db); err != nil {
			return fmt.Errorf("执行备份任务失败: %v", err)
		}
		return nil
	case "r":
		// 解析run命令的参数
		if err := runCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析run命令参数失败: %v", err)
		}
		// 执行run命令的逻辑
		if err := runCmdMain(db); err != nil {
			return fmt.Errorf("执行备份任务失败: %v", err)
		}
		return nil
	case "add":
		// 解析add命令的参数
		if err := addCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析add命令参数失败: %v", err)
		}
		// 执行add命令的逻辑
		if err := addCmdMain(db); err != nil {
			return fmt.Errorf("添加项目失败: %v", err)
		}
		return nil
	case "a":
		// 解析add命令的参数
		if err := addCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析add命令参数失败: %v", err)
		}
		// 执行add命令的逻辑
		if err := addCmdMain(db); err != nil {
			return fmt.Errorf("添加项目失败: %v", err)
		}
		return nil
	case "delete":
		// 解析delete命令的参数
		if err := deleteCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析delete命令参数失败: %v", err)
		}
		// 执行delete命令的逻辑
		if err := deleteCmdMain(db); err != nil {
			return fmt.Errorf("删除项目失败: %v", err)
		}
		return nil
	case "d":
		// 解析delete命令的参数
		if err := deleteCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析delete命令参数失败: %v", err)
		}
		// 执行delete命令的逻辑
		if err := deleteCmdMain(db); err != nil {
			return fmt.Errorf("删除项目失败: %v", err)
		}
		return nil
	case "edit":
		// 解析edit命令的参数
		if err := editCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析edit命令参数失败: %v", err)
		}
		// 执行edit命令的逻辑
		if err := editCmdMain(db); err != nil {
			return fmt.Errorf("编辑项目失败: %v", err)
		}
		return nil
	case "e":
		// 解析edit命令的参数
		if err := editCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析edit命令参数失败: %v", err)
		}
		// 执行edit命令的逻辑
		if err := editCmdMain(db); err != nil {
			return fmt.Errorf("编辑项目失败: %v", err)
		}
		return nil
	case "log":
		// 解析log命令的参数
		if err := logCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析log命令参数失败: %v", err)
		}
		// 执行log命令的逻辑
		if err := logCmdMain(db, 1, *logLimit); err != nil {
			return fmt.Errorf("查看日志失败: %v", err)
		}
		return nil
	case "show":
		// 解析show命令的参数
		if err := showCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析show命令参数失败: %v", err)
		}
		// 执行show命令的逻辑
		if err := showCmdMain(db); err != nil {
			return fmt.Errorf("查看指定备份任务的信息失败: %v", err)
		}
		return nil
	case "s":
		// 解析show命令的参数
		if err := showCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析show命令参数失败: %v", err)
		}
		// 执行show命令的逻辑
		if err := showCmdMain(db); err != nil {
			return fmt.Errorf("查看指定备份任务的信息失败: %v", err)
		}
		return nil
	case "unpack":
		// 解析unpack命令的参数
		if err := unpackCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析unpack命令参数失败: %v", err)
		}
		// 执行unpack命令的逻辑
		if err := unpackCmdMain(db); err != nil {
			return fmt.Errorf("解压备份任务失败: %v", err)
		}
		return nil
	case "u":
		// 解析unpack命令的参数
		if err := unpackCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析unpack命令参数失败: %v", err)
		}
		// 执行unpack命令的逻辑
		if err := unpackCmdMain(db); err != nil {
			return fmt.Errorf("解压备份任务失败: %v", err)
		}
		return nil
	case "zip":
		// 解析zip命令的参数
		if err := zipCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析zip命令参数失败: %v", err)
		}
		// 执行zip命令的逻辑
		if err := zipCmdMain(); err != nil {
			return fmt.Errorf("打包ZIP文件失败: %v", err)
		}
		return nil
	case "z":
		// 解析zip命令的参数
		if err := zipCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析zip命令参数失败: %v", err)
		}
		// 执行zip命令的逻辑
		if err := zipCmdMain(); err != nil {
			return fmt.Errorf("打包ZIP文件失败: %v", err)
		}
		return nil
	case "unzip":
		// 解析unzip命令的参数
		if err := unzipCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析unzip命令参数失败: %v", err)
		}
		// 执行unzip命令的逻辑
		if err := unzipCmdMain(); err != nil {
			return fmt.Errorf("解压ZIP文件失败: %v", err)
		}
		return nil
	case "uz":
		// 解析unzip命令的参数
		if err := unzipCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析unzip命令参数失败: %v", err)
		}
		// 执行unzip命令的逻辑
		if err := unzipCmdMain(); err != nil {
			return fmt.Errorf("解压ZIP文件失败: %v", err)
		}
		return nil
	// 打印版本信息
	case "version":
		// 解析version命令的参数
		if err := versionCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析version命令参数失败: %v", err)
		}
		// 执行version命令的逻辑
		v := version.Get()
		if versionInfo, err := v.SprintVersion("text"); err != nil {
			return fmt.Errorf("获取版本信息失败: %v", err)
		} else {
			CL.Green(versionInfo)
		}
		return nil
	// 打印帮助信息
	case "help":
		// 解析help命令的参数
		if err := helpCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析help命令参数失败: %v", err)
		}

		// 如果没有指定子命令，则打印帮助信息
		if len(helpCmd.Args()) == 0 {
			return fmt.Errorf("请指定要查看帮助的命令, 例如: 'cbk help 指定命令'")
		}

		// 执行help命令的逻辑
		if err := helpCmdMain(helpCmd.Args()[0]); err != nil {
			return fmt.Errorf("打印帮助信息失败: %v", err)
		}
		return nil
	case "clear":
		// 解析clear命令的参数
		if err := clearCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析clear命令参数失败: %v", err)
		}

		// 执行clear命令的逻辑
		if err := clearCmdMain(db); err != nil {
			return fmt.Errorf("清空数据库失败: %v", err)
		}
		return nil
	case "init":
		// 解析init命令的参数
		if err := initCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析init命令参数失败: %v", err)
		}

		// 执行init命令的逻辑
		if err := initCmdMain(*initType); err != nil {
			return fmt.Errorf("生成文件失败: %v", err)
		}
		return nil
	case "export":
		// 解析export命令的参数
		if err := exportCmd.Parse(args[1:]); err != nil {
			return fmt.Errorf("解析export命令参数失败: %v", err)
		}
		// 执行export命令的逻辑
		if err := exportCmdMain(db); err != nil {
			return fmt.Errorf("导出数据库失败: %v", err)
		}
		return nil
	// 未知命令
	default:
		return fmt.Errorf("未知命令: %s", args[0])
	}
}
