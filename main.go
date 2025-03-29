// main.go
package main

import (
	"cbk/cmd"
	"cbk/pkg/version"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/colorlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// 定义全局颜色渲染器
var CL = colorlib.NewColorLib()

//go:embed sql/init.sql
var initSql string // 初始化SQL语句

const (
	cbkHomeDir = ".cbk"   // 数据目录
	cbkDBFile  = "cbk.db" // 数据库文件
	cbkDataDir = "data"   // 数据目录
)

func main() {
	// 初始化数据库
	_, err := initDB()
	if err != nil {
		CL.PrintErrorf("初始化数据库失败: %v", err)
		os.Exit(1)
	}

	// 初始化数据目录
	if err := initDataDir(); err != nil {
		CL.PrintErrorf("初始化数据目录失败: %v", err)
		os.Exit(1)
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
		CL.Green(v.SprintVersion("simple"))
		return
	}
	// 打印更详细的版本信息
	if *vvFlag {
		v := version.Get()
		CL.Green(v.SprintVersion("text"))
		return
	}
	// 打印帮助信息
	if *hFlag || *helpFlag {
		fmt.Println("显示帮助信息")
		return
	}

	// 获取命令行参数
	args := flag.Args()

	// 检查是否有子命令
	if len(args) == 0 {
		flag.PrintDefaults()
		return
	}

	// 执行子命令
	err = cmd.ExecuteCommands(args)
	if err != nil {
		CL.PrintError(err)
		os.Exit(1)
	}
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
	dbDir := filepath.Join(homeDir, cbkHomeDir)

	// 检查数据库目录是否存在, 如果不存在, 则创建
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("检查数据库目录失败: %w", err)
	}

	// 构造数据库文件路径
	dbPath := filepath.Join(dbDir, cbkDBFile)

	// 检查数据库文件是否存在, 如果不存在, 则创建并初始化
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// 连接数据库
		db, err := sqlx.Connect("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("连接数据库失败: %w", err)
		}

		// 执行初始化SQL语句
		if _, err := db.Exec(initSql); err != nil {
			return nil, fmt.Errorf("执行初始化SQL语句失败: %w", err)
		}

		// 直接返回连接
		return db, nil
	} else if err != nil {
		return nil, fmt.Errorf("连接到数据库文件失败: %w", err)
	}

	// 存在数据库文件, 直接连接
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
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
	dataDir := filepath.Join(homeDir, cbkHomeDir, cbkDataDir)

	// 检查数据目录是否存在, 如果不存在, 则创建
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("创建数据目录失败: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("检查数据目录失败: %w", err)
	}

	return nil
}
