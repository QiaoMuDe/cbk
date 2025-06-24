package globals

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/v6/table"
)

//go:embed init.sql
var InitSql string // 初始化SQL语句

const (
	CbkHomeDir       = ".cbk"       // 数据目录
	CbkDBFile        = "cbk.db"     // 数据库文件
	CbkDataDir       = "data"       // 数据目录
	OneKeyScriptName = "onekey_bak" // 一键脚本名称
)

// 数据库文件路径
var CbkDbPath = filepath.Join(CbkHomeDir, CbkDBFile)

// 定义任务表结构体
type BackupTask struct {
	TaskID          int    `db:"task_id"`          // 任务ID
	TaskName        string `db:"task_name"`        // 任务名
	TargetDirectory string `db:"target_directory"` // 目标目录
	BackupDirectory string `db:"backup_directory"` // 备份目录
	RetentionCount  int    `db:"retention_count"`  // 保留数量
	RetentionDays   int    `db:"retention_days"`   // 保留天数
	NoCompression   int    `db:"no_compression"`   // 是否禁用压缩(默认启用压缩, 0 表示启用压缩, 1 表示禁用压缩)
	ExcludeRules    string `db:"exclude_rules"`    // 排除规则
}

// 定义任务表结构体切片
type BackupTasks []BackupTask

// 定义备份记录表结构体
type BackupRecord struct {
	VersionID      string `db:"version_id"`       // 版本ID
	TaskID         int    `db:"task_id"`          // 任务ID
	Timestamp      string `db:"timestamp"`        // 时间戳
	TaskName       string `db:"task_name"`        // 任务名
	BackupStatus   string `db:"backup_status"`    // 备份状态
	BackupFileName string `db:"backup_file_name"` // 备份文件名
	BackupSize     string `db:"backup_size"`      // 备份文件大小
	BackupPath     string `db:"backup_path"`      // 备份文件路径
	VersionHash    string `db:"version_hash"`     // 版本哈希
}

// 定义备份记录表结构体切片
type BackupRecords []BackupRecord

// 定义任务配置的结构体
type TaskConfig struct {
	Task Task `yaml:"task"`
}

// 定义任务的结构体
type Task struct {
	Name          string    `yaml:"name"`            // 任务名
	Target        string    `yaml:"target"`          // 目标目录
	Backup        string    `yaml:"backup"`          // 备份目录
	Retention     Retention `yaml:"retention"`       // 保留策略
	BackupDirName string    `yaml:"backup_dir_name"` // 备份目录名
	NoCompression int       `yaml:"no_compression"`  // 是否禁用压缩(默认启用压缩, 0 表示启用压缩, 1 表示禁用压缩)
	ExcludeRules  string    `yaml:"exclude_rules"`   // 排除规则
}

// 定义保留策略的结构体
type Retention struct {
	Count int `yaml:"count"` // 保留数量
	Days  int `yaml:"days"`  // 保留天数(配置为 0 表示不限制天数)
}

// 定义排除函数的类型
// ExcludeFunc 是一个函数类型，用于排除文件或目录
// path: 文件或目录的路径
// info: 文件或目录的信息
// 返回值: 如果文件或目录应该被排除, 返回 true, 否则返回 false
type ExcludeFunc func(path string, info os.FileInfo) bool

// 全局排除函数变量
var (
	// NoExcludeFunc 是一个空的排除函数，表示不排除任何文件或目录
	NoExcludeFunc ExcludeFunc = func(path string, info os.FileInfo) bool {
		return false // 不排除任何文件或目录
	}
)

// 定义存放表格样式的MAP
var (
	TableStyle = map[string]table.Style{
		"df":   table.StyleDefault,       // 默认样式
		"bd":   table.StyleBold,          // 加粗样式
		"cb":   table.StyleColoredBright, // 亮色样式
		"cd":   table.StyleColoredDark,   // 暗色样式
		"de":   table.StyleDouble,        // 双边框样式
		"lt":   table.StyleLight,         // 浅色样式
		"ro":   table.StyleRounded,       // 圆角样式
		"none": StyleNone,                // 禁用样式
	}
)

// 定义禁用样式
var StyleNone = table.Style{
	Box: table.BoxStyle{
		PaddingLeft:      " ", // 左边框
		PaddingRight:     " ", // 右边框
		MiddleHorizontal: " ", // 水平线
		MiddleVertical:   " ", // 垂直线
		TopLeft:          " ", // 左上角
		TopRight:         " ", // 右上角
		BottomLeft:       " ", // 左下角
		BottomRight:      " ", // 右下角
	},
}

// 定义存放表格样式的切片
var TableStyleList = []string{"df", "bd", "cb", "cd", "de", "lt", "ro", "none"}

// 表格样式的帮助信息
var TableStyleHelp = `支持的表格样式：
        "df":    默认样式
	"bd":    加粗样式
	"cb":    亮色样式
	"cd":    暗色样式
	"de":    双边框样
	"lt":    浅色样式
	"ro":    圆角样式
	"none":  禁用样式`
