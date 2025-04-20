package globals

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	CbkHomeDir = ".cbk"   // 数据目录
	CbkDBFile  = "cbk.db" // 数据库文件
	CbkDataDir = "data"   // 数据目录
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
}

// 定义保留策略的结构体
type Retention struct {
	Count int `yaml:"count"` // 保留数量
	Days  int `yaml:"days"`  // 保留天数(配置为 0 表示不限制天数)
}

// 定义过滤函数的类型
// FilterFunc 是一个函数类型，用于过滤文件或目录
// path: 文件或目录的路径
// info: 文件或目录的信息
// 返回值: 如果文件或目录应该被过滤掉, 返回 true, 否则返回 false
type FilterFunc func(path string, info os.FileInfo) bool

// 全局过滤函数变量
var (
	// NoFilter 是一个不过滤任何文件的过滤函数
	// 参数:
	//   path - 文件路径(未使用)
	//   info - 文件信息(未使用)
	// 返回值:
	//   bool - 总是返回false，表示不过滤任何文件
	NoFilter FilterFunc = func(path string, info os.FileInfo) bool {
		return false // 默认不跳过任何文件
	}

	// FilterDirectories 创建一个过滤函数，用于过滤指定目录名的目录
	// 参数:
	//   dirNames - 需要过滤的目录名列表
	// 返回值:
	//   FilterFunc - 过滤函数，当目录名匹配列表中的任一项时返回true
	FilterDirectories = func(dirNames []string) FilterFunc {
		return func(path string, info os.FileInfo) bool {
			// 只处理目录
			if info.IsDir() {
				// 遍历所有需要过滤的目录名
				for _, dirName := range dirNames {
					// 检查当前目录名是否匹配
					if filepath.Base(path) == dirName {
						return true // 匹配则返回true表示需要过滤
					}
				}
			}
			return false // 不匹配则返回false表示不需要过滤
		}
	}

	// FilterFileExtensions 创建一个过滤函数，用于过滤指定扩展名的文件
	// 参数:
	//   exts - 需要过滤的文件扩展名列表(如: [".txt", ".log"])
	// 返回值:
	//   FilterFunc - 过滤函数，当文件扩展名匹配列表中的任一扩展名时返回true
	FilterFileExtensions = func(exts []string) FilterFunc {
		return func(path string, info os.FileInfo) bool {
			// 跳过目录，只处理文件
			if !info.IsDir() {
				// 获取文件扩展名(包含点)
				ext := filepath.Ext(path)
				// 遍历所有需要过滤的扩展名
				for _, extToFilter := range exts {
					// 检查当前文件扩展名是否匹配
					if ext == extToFilter {
						return true // 匹配则返回true表示需要过滤
					}
				}
			}
			return false // 不匹配则返回false表示不需要过滤
		}
	}

	// FilterWithRegex 根据正则表达式模式创建一个过滤函数
	// 参数:
	//   pattern - 用于匹配的正则表达式字符串
	// 返回值:
	//   FilterFunc - 符合正则表达式的文件/路径将返回true
	//   error - 如果正则表达式编译失败则返回错误
	FilterWithRegex = func(pattern string) (FilterFunc, error) {
		// 编译正则表达式
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("编译正则表达式失败: %w", err)
		}

		// 返回过滤函数
		return func(path string, info os.FileInfo) bool {
			// 使用编译好的正则表达式匹配路径
			return regex.MatchString(path)
		}, nil
	}

	// FilterFileNames 创建一个过滤函数，用于过滤指定文件名的文件
	// 参数:
	//   fileNames - 需要过滤的文件名列表
	// 返回值:
	//   FilterFunc - 过滤函数，当文件名匹配列表中的任一文件名时返回true
	FilterFileNames = func(fileNames []string) FilterFunc {
		return func(path string, info os.FileInfo) bool {
			// 跳过目录，只处理文件
			if !info.IsDir() {
				// 遍历所有需要过滤的文件名
				for _, fileName := range fileNames {
					// 检查当前文件名是否匹配
					if filepath.Base(path) == fileName {
						return true // 匹配则返回true表示需要过滤
					}
				}
			}
			return false // 不匹配则返回false表示不需要过滤
		}
	}
)
