// main.go
package main

import (
	"cbk/cmd"
	"cbk/models"
	"cbk/pkg/version"
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/colorlib"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 定义全局颜色渲染器
var CL = colorlib.NewColorLib()

func main() {
	// 如果db文件不存在，则创建
	if _, err := os.Stat("test/backup.db"); os.IsNotExist(err) {
		initmain()
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
		v.PrintVersion("simple")
		return
	}
	// 打印更详细的版本信息
	if *vvFlag {
		v := version.Get()
		v.PrintVersion("text")
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
	err := cmd.ExecuteCommands(args)
	if err != nil {
		CL.PrintError(err)
		os.Exit(1)
	}
}

func initmain() {
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open("test/backup.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移表结构
	db.AutoMigrate(&models.BackupRecord{}, &models.BackupTask{})
}
