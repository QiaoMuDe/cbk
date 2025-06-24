// cmd.go
package cmd

import (
	"cbk/pkg/globals"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/verman"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// 定义全局颜色渲染器
var CL = colorlib.NewColorLib()

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
			db.Close()
		}
	}()

	// 初始化数据目录
	if initDataDirErr := initDataDir(); initDataDirErr != nil {
		return fmt.Errorf("初始化数据目录失败: %w", initDataDirErr)
	}

	// 获取命令行参数
	args := qflag.Args()

	// 检查是否有子命令
	if len(args) == 0 {
		qflag.PrintHelp()
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
		if _, execErr := db.Exec(globals.InitSql); execErr != nil {
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
	case listCmd.LongName(), listCmd.ShortName(): // list命令
		// 执行list命令的逻辑
		if err := listCmdMain(db); err != nil {
			return fmt.Errorf("列出项目列表失败: %v", err)
		}
		return nil
	case runCmd.LongName(), runCmd.ShortName(): // run命令
		// 执行run命令的逻辑
		if err := runCmdMain(db); err != nil {
			return fmt.Errorf("执行备份任务失败: %v", err)
		}
		return nil
	case addCmd.LongName(), addCmd.ShortName(): // add命令
		// 执行add命令的逻辑
		if err := addCmdMain(db); err != nil {
			return fmt.Errorf("添加项目失败: %v", err)
		}
		return nil
	case deleteCmd.LongName(), deleteCmd.ShortName(): // delete命令
		// 执行delete命令的逻辑
		if err := deleteCmdMain(db); err != nil {
			return fmt.Errorf("删除项目失败: %v", err)
		}
		return nil
	case editCmd.LongName(), editCmd.ShortName(): // edit命令
		// 执行edit命令的逻辑
		if err := editCmdMain(db); err != nil {
			return fmt.Errorf("编辑项目失败: %v", err)
		}
		return nil
	case logCmd.LongName(), logCmd.ShortName(): // log命令
		// 执行log命令的逻辑
		if err := logCmdMain(db, 1, logLimit.Get()); err != nil {
			return fmt.Errorf("查看日志失败: %v", err)
		}
		return nil
	case showCmd.LongName(), showCmd.ShortName(): // show命令
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
		v := verman.Get()
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
		if err := initCmdMain(*initType, db); err != nil {
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
